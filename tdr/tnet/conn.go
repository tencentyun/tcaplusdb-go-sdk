package tnet

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	log "github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/**
	@brief 回调函数
	@param [IN] url 连接地址
	@param [IN] pkg 未反序列化的传输包
	@param [IN] cbPara	回调参数
**/
type RecvCallBackFunc func(url *string, pkg *common.PKGBuffer, cbPara interface{}) error

/**
	@brief 解析包长度，用于切分
	@param [IN] buf
	@retval [IN] int 长度
**/
type ParseFunc func(buf []byte) int

const (
	Connected    = 0
	Connecting   = 1
	Disconnected = 2
	ReadErr      = 3
	WriteErr     = 4
)

/**
	@brief tcaplus api客户端
	@param [IN] appId 业务id
	@param [IN] zoneList 区列表
	@param [IN] dirUrl	dir地址
	@param [IN] initFlag 是否初始化
	@param [IN] netServer 服务管理
**/
type Conn struct {
	netConn   net.Conn
	network   string
	url       string
	ip        string
	port      string
	stat      int32

	//回调控制
	parseFunc  ParseFunc        //通过parse判断是否收到完整包
	cbFunc     RecvCallBackFunc //收到响应后会调用回调
	cbPara     interface{}      //回调参数
	timeout    time.Duration    //connect 超时时间
	createTime time.Time

	sendChan chan *Buf

	//协程控制
	closeFlag chan bool
	sync.WaitGroup
}

//url 格式为tcp://127.0.0.1:80
func ParseUrl(url *string) (network, ip, port string, err error) {
	list := strings.Split(*url, "://")
	if len(list) < 2 {
		return "", "", "", &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "url is invalid"}
	}

	network = list[0]
	addr := strings.Split(list[1], ":")
	if len(addr) < 2 {
		return "", "", "", &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "url parse network ip port fail"}
	}
	ip = addr[0]
	port = addr[1]

	return network, ip, port, nil
}

//url 格式为tcp://127.0.0.1:80
func NewConn(url string, timeout time.Duration,
	parseFunc ParseFunc, cbFunc RecvCallBackFunc, cbPara interface{}) (*Conn, error) {
	network, ip, port, err := ParseUrl(&url)
	if nil != err {
		return nil, err
	}

	cn := &Conn{
		netConn:    nil,
		url:        url,
		network:    network,
		ip:         ip,
		port:       port,
		closeFlag:  make(chan bool, 1),
		stat:       Connecting,
		parseFunc:  parseFunc,
		cbFunc:     cbFunc,
		cbPara:     cbPara,
		createTime: time.Now(),
		timeout:    timeout,
		sendChan:   make(chan *Buf, common.ConfigProcReqDepth-1),
	}
	go cn.connect()
	return cn, nil
}

func (c *Conn) connect() {
	addr := c.ip + ":" + c.port
	Conn, err := net.DialTimeout(c.network, addr, c.timeout)
	if nil != err {
		log.ERR("connect failed, %s", err.Error())
		atomic.StoreInt32(&c.stat, Disconnected)
		return
	}
	log.INFO("connect addr %s success", addr)
	c.netConn = Conn
	atomic.StoreInt32(&c.stat, Connected)
	go c.recvRoutine()
	go c.SendRoutine()
}

type Buf []byte

func (c *Conn) sendPkg(buf *Buf) {
	var err error
	if c.GetStat() != Connected || c.netConn == nil {
		log.ERR("api connect proxy stat not connected")
		err = fmt.Errorf("api connect proxy stat not connected")
	} else {
		// 设置30秒的发送超时，防止一直卡住
		c.netConn.SetWriteDeadline(common.TimeNow.Add(common.ConfigReadWriteTimeOut * time.Second))
		_, err = c.netConn.Write([]byte(*buf))
		if err != nil {
			atomic.StoreInt32(&c.stat, WriteErr)
			log.ERR("Send data failed %v, conn url %v", err.Error(), c.url)
		}
	}
}

func (c *Conn) Send(buf []byte) error{
	pkg := Buf(buf)
	select {
	case <-c.closeFlag:
		log.INFO("close flag, %s", c.url)
		return fmt.Errorf("close flag, %s", c.url)
	case c.sendChan <- &pkg:
	}
	return nil
}

func (c *Conn) SendRoutine() {
	c.Add(1)
	defer c.Done()

	proc := func(buf *Buf) {
		var length = len(c.sendChan)
		c.sendPkg(buf)
		for index := 0; index < length ; index++ {
			buf = <-c.sendChan
			c.sendPkg(buf)
		}
	}

	for {
		select {
		case <-c.closeFlag:
			for len(c.sendChan) > 0 {
				proc(<-c.sendChan)
			}
			log.INFO("close flag, %s", c.url)
			return
		case buf := <-c.sendChan:
			proc(buf)
		}
	}
}

func (c *Conn) GetStat() int32 {
	return atomic.LoadInt32(&c.stat)
}

func (c *Conn) Close() {
	c.closeFlag <- true
	atomic.StoreInt32(&c.stat, Disconnected)
	close(c.closeFlag)
	if c.netConn != nil {
		_ = c.netConn.Close()
	}
}

/*
	需求：复用内存，防止频繁gc，减少内存拷贝次数
	当前流程：复用20M缓冲区从网络里拷贝获取响应，切分响应拷贝，为每个响应分配内存并拷贝。再次分配响应内存，对响应反序列化，可能发生两次拷贝
	分析必要的分配与拷贝：1、分配内存用于从网络中读取（拷贝必要，分配内存分配几次？）
					 2、分配内存生成响应（拷贝还是复用？复用的化每次需要在从网络读取时分配内存。
							拷贝需要在解包的地方分配内存但从网络中读取的内存可以做到复用）共同问题：什么时候释放内存引发gc均不可控

	暂定方案：使用对象池分配内存，长时间运行可以减少一次内存分配，减少2次内存拷贝。

	接口调用点：从网络中读，切分出请求，对象池操作
*/

func (c *Conn) recvRoutine() {
	c.Add(1)
	defer c.Done()
	pkgManager := common.GetPKGManager(nil)
	for {
		select {
		case <-c.closeFlag:
			log.INFO("close flag, %s", c.url)
			return
		default:
			if c.GetStat() != Connected || c.netConn == nil {
				log.ERR("api connect proxy stat not connected")
				return
			}
			// 满了需要重新申请一段buffer，并将这次未处理完的buffer拷贝
			if pkgManager.BufferIsFull() {
				pkgManager = common.GetPKGManager(pkgManager)
			}
			// 读取网络报文
			n, err := pkgManager.Read(c.netConn)
			if err != nil {
				atomic.StoreInt32(&c.stat, ReadErr)
				if err == io.EOF {
					log.INFO("read close:%s, %s", err.Error(), c.url)
				} else {
					log.ERR("read err:%s, %s", err.Error(), c.url)
				}
				return
			}

			if n > 0 {
				//fmt.Println("recv:", string(buf[0:n]), "len:", n)
				for {
					//判断是否收到完整包
					// 判断buffer是否可以切分出一个完整包
					pkgSize := c.parseFunc(pkgManager.ValidBuffer())
					if pkgSize <= pkgManager.ValidLength() && pkgSize > 0 {
						// 将包从buffer中取出
						pkg := pkgManager.GetPkgBuffer(pkgSize)
						//收到完整包
						//log.DEBUG("url %s %d pkgSize %d finish recvBuffer.Len %d", c.url, n, pkgSize, len(pkg.GetData()))

						//回调处理包
						err = c.cbFunc(&c.url, pkg, c.cbPara)
						if err != nil {
							log.ERR("cbFunc err:%s, %s", err.Error(), c.url)
						}
					} else {
						//log.DEBUG("url %s %d pkgSize %d <= recvBuffer.Len %d", c.url, n, pkgSize, recvBuffer.Len())
						break
					}
				}
			}
		}
	}
}
