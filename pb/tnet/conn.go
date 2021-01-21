package tnet

import (
	"bytes"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	log "github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type RecvCallBackFunc func(url *string, buf []byte, cbPara interface{}) error
type ParseFunc func(buf []byte) int

const (
	Connected    = 0
	Connecting   = 1
	Disconnected = 2
	ReadErr      = 3
	WriteErr     = 4
)

type Conn struct {
	netConn   net.Conn
	network   string
	url       string
	ip        string
	port      string
	closeFlag chan bool
	stat      int32

	parseFunc  ParseFunc        //通过parse判断是否收到完整包
	cbFunc     RecvCallBackFunc //收到响应后会调用回调
	cbPara     interface{}      //回调参数
	timeout    time.Duration    //connect 超时时间
	createTime time.Time

	sendChan   chan *sendBuf
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
		sendChan:   make(chan *sendBuf, common.ConfigProcReqDepth-1),
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
	go c.process()
	go c.mergeSend()
}

type sendBuf struct {
	err error
	buf []byte
	sync.WaitGroup
}

func (c *Conn) send(buf []byte, ret []*sendBuf) {
	_, err := c.netConn.Write(buf)
	if err != nil {
		atomic.StoreInt32(&c.stat, WriteErr)
		log.ERR("Send data failed %v, conn url %v", err.Error(), c.url)
	}
	for _, v := range ret {
		v.err = err
		v.Done()
	}
}

func (c *Conn) Send(buf []byte) (int, error) {
	b := &sendBuf{buf: buf}
	b.Add(1)
	select {
	case <-c.closeFlag:
		log.INFO("close flag, %s", c.url)
		return 0, fmt.Errorf("close flag, %s", c.url)
	case c.sendChan <- b:
	}
	b.Wait()
	return 0, b.err
}

func (c *Conn) mergeSend() {
	var index int
	var length int
	sendBuffer := bytes.NewBuffer(make([]byte, 0, 20<<20))
	depth := common.ConfigProcReqDepth

	proc := func(buf *sendBuf) {
		bufs := make([]*sendBuf, 0, depth)
		bufs = append(bufs, buf)
		sendBuffer.Write(buf.buf)
		length = len(c.sendChan)
		index = 0
		for ; index < length; index++ {
			buf = <-c.sendChan
			bufs = append(bufs, buf)
			sendBuffer.Write(buf.buf)
			// 当数据大小大于10M时直接发送
			if sendBuffer.Len() > 10<<20 {
				break
			}
		}
		c.send(sendBuffer.Bytes(), bufs)
		sendBuffer.Reset()
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

func (c *Conn) process() {
	buf := make([]byte, 10*1024*1024)
	recvBuffer := bytes.NewBuffer(make([]byte, 0, 20*1024*1024))

	for {
		select {
		case <-c.closeFlag:
			log.INFO("close flag, %s", c.url)
			return
		default:
			n, err := c.netConn.Read(buf)
			if err != nil {
				atomic.StoreInt32(&c.stat, ReadErr)
				log.ERR("read err:%s, %s", err.Error(), c.url)
				return
			}

			if n > 0 {
				//fmt.Println("recv:", string(buf[0:n]), "len:", n)
				recvBuffer.Write(buf[0:n])
				for {
					//判断是否收到完整包
					pkgSize := c.parseFunc(recvBuffer.Bytes())
					if pkgSize <= recvBuffer.Len() && pkgSize > 0 {
						//收到完整包
						//log.DEBUG("url %s %d pkgSize %d finish recvBuffer.Len %d", c.url, n, pkgSize, recvBuffer.Len())
						pkg := make([]byte, pkgSize)
						_, err := recvBuffer.Read(pkg)
						if err != nil {
							log.ERR("recvBuffer.Read err:%s, %s", err.Error(), c.url)
						}

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
