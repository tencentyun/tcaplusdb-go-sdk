package tnet

import (
	"bufio"
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
type RecvCallBackFunc func(url *string, pkg *PKG) error

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

type Buf []byte

/**
	@brief tcaplus api客户端
	@param [IN] appId 业务id
	@param [IN] zoneList 区列表
	@param [IN] dirUrl	dir地址
	@param [IN] initFlag 是否初始化
	@param [IN] netServer 服务管理
**/
type Conn struct {
	netConn net.Conn
	network string
	url     string
	ip      string
	port    string
	stat    int32

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

	//io buf
	wrSize int
	wr     *bufio.Writer
	rd     *PKGMemory
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
	parseFunc ParseFunc, cbFunc RecvCallBackFunc, cbPara interface{}, writeIoSize int) (*Conn,
	error) {
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
		sendChan:   make(chan *Buf, common.ConfigProcConBufDepth-1),
		wrSize:     writeIoSize,
	}
	go cn.connect()
	return cn, nil
}

func (c *Conn) connect() {
	c.Add(1)
	defer c.Done()
	addr := c.ip + ":" + c.port
	Conn, err := net.DialTimeout(c.network, addr, c.timeout)
	if nil != err {
		log.ERR("connect failed, %s", err.Error())
		atomic.StoreInt32(&c.stat, Disconnected)
		return
	}
	log.INFO("connect addr %s success", addr)
	c.netConn = Conn
	c.rd = nil
	c.wr = bufio.NewWriterSize(Conn, c.wrSize)
	atomic.StoreInt32(&c.stat, Connected)
	go c.recvRoutine()
	go c.SendRoutine()
}

func (c *Conn) sendPkg(buf *Buf) {
	var err error
	if c.GetStat() != Connected || c.netConn == nil {
		log.ERR("api connect proxy stat not connected")
		err = fmt.Errorf("api connect proxy stat not connected")
	} else {
		_, err = c.wr.Write([]byte(*buf))
		if err != nil {
			atomic.StoreInt32(&c.stat, WriteErr)
			log.ERR("Send data failed %v, conn url %v", err.Error(), c.url)
		}
	}
}

func (c *Conn) Send(buf []byte) error {
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
		// 设置30秒的发送超时，防止一直卡住
		c.netConn.SetWriteDeadline(common.TimeNow.Add(ConfigReadWriteTimeOut * time.Second))
		var length = len(c.sendChan)
		c.sendPkg(buf)
		for index := 0; index < length; index++ {
			buf = <-c.sendChan
			c.sendPkg(buf)
		}
		err := c.wr.Flush()
		if err != nil {
			atomic.StoreInt32(&c.stat, WriteErr)
			log.ERR("Flush data failed %v, conn url %v", err.Error(), c.url)
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
	c.Wait()
}

/*
	接口调用点：从网络中读，切分出请求，对象池操作
*/
func (c *Conn) recvRoutine() {
	c.Add(1)
	defer c.Done()
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
			if c.rd == nil {
				c.rd = GetPKGMemory(nil)
			} else if c.rd.BufferIsFull() {
				// 满了需要重新申请一段buffer，并将这次未处理完的buffer拷贝
				c.rd = GetPKGMemory(c.rd)
			}
			// 读取网络报文
			c.netConn.SetReadDeadline(common.TimeNow.Add(ConfigReadWriteTimeOut * time.Second))
			n, err := c.rd.ReadFromNetConn(c.netConn)
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
				c.procRecvPkg()
			} else {
				//TODO长时间无包释放c.rd,释放内存
			}
		}
	}
}

func (c *Conn) procRecvPkg() {
	for {
		//判断是否收到完整包
		// 判断buffer是否可以切分出一个完整包
		pkgSize := c.parseFunc(c.rd.ValidBuffer())
		if pkgSize <= c.rd.ValidLength() && pkgSize > 0 {
			// 将包从buffer中取出
			pkg := c.rd.GetPkg(pkgSize)
			pkg.cbPara = c.cbPara
			//收到完整包回调处理包
			err := c.cbFunc(&c.url, pkg)
			if err != nil {
				log.ERR("cbFunc err:%s, %s", err.Error(), c.url)
			}
		} else {
			break
		}
	}
}
