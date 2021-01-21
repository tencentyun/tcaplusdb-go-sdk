package router

import (
	"bytes"
	"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/tnet"
	"time"
)

const (
	NotSignUp     = 0
	SignUpIng     = 1
	SignUpSuccess = 2
	SignUpFail    = 3
)

type server struct {
	appId       uint64
	zoneId      uint32
	signature   string
	proxyUrl    string
	signUpFlag  uint32
	conn        *tnet.Conn
	connectTime time.Time
	signUpTime  time.Time
	router      interface{}
	prepareStop bool //proxy准备stop

	respMsgChanList []chan []byte
	closeFlag       chan struct{}
}

func (s *server) initRecv() {
	if len(s.respMsgChanList) > 0 {
		return
	}
	s.respMsgChanList = make([]chan []byte, common.ConfigProcRespRoutineNum)
	s.closeFlag = make(chan struct{})
	for i := 0; i < common.ConfigProcRespRoutineNum; i++ {
		s.respMsgChanList[i] = make(chan []byte, common.ConfigProcRespDepth)
		go func(id int) {
			proc := func(buf []byte) {
				logger.DEBUG("recv proxy response, unpack.")
				start := time.Now()
				resp := tcaplus_protocol_cs.NewTCaplusPkg()
				if err := resp.Unpack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion, buf); err != nil {
					logger.ERR("Unpack proxy msg failed, url %s err %v", s.proxyUrl, err.Error())
					return
				}

				if time.Now().Sub(start) > 10*time.Millisecond {
					logger.WARN("unpack > 10ms data %v.", buf)
				}

				s.processRsp(resp)
			}
			ch := s.respMsgChanList[id]
			for {
				select {
				case <-s.closeFlag:
					for len(ch) > 0 {
						for i := 0; i< len(ch); i++ {
							proc(<- ch)
						}
					}
					logger.INFO("exit recv routine.")
					return
				case buf := <- ch:
					proc(buf)
				}
			}
		}(i)
	}
}

func (s *server) getSignUpStat() uint32 {
	return s.signUpFlag
}

func (s *server) update(isInRemoveList bool) {
	if !isInRemoveList {
		s.connect()
	}
}

func (s *server) isAvailable() bool {
	if s.conn == nil {
		return false
	}

	if s.prepareStop {
		return false
	}

	if s.conn.GetStat() != tnet.Connected {
		return false
	}

	if s.signUpFlag != SignUpSuccess {
		return false
	}
	return true
}

func (s *server) disConnect() {
	s.prepareStop = false
	s.signUpFlag = NotSignUp
	close(s.closeFlag)
	s.respMsgChanList = nil
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

func (s *server) send(data []byte) error {
	logger.DEBUG("send to proxy %s", s.proxyUrl)
	_, err := s.conn.Send(data)
	return err
}

func (s *server) connect() {
	s.initRecv()
	if s.conn == nil {
		//连接proxy, 3s超时
		conn, err := tnet.NewConn(s.proxyUrl, 3*time.Second, ParseProxyPkgLen, ProxyCallBackFunc, s)
		if err != nil {
			logger.ERR("new conn failed %v", err)
			return
		}
		s.conn = conn
		s.connectTime = time.Now()
		logger.DEBUG("start connect proxy %s", s.proxyUrl)
	} else {
		if s.conn.GetStat() == tnet.Connected {
			//连接成功
			//认证
			if s.signUpFlag == NotSignUp {
				s.signUpFlag = SignUpIng
				logger.INFO("start sign up proxy %v", s.proxyUrl)
				s.signUp()
			} else if s.signUpFlag == SignUpIng && time.Now().Sub(s.signUpTime).Seconds() > 3 {
				//认证超时，重新认证
				logger.ERR("sign up proxy %v timeout(3s), conn stat %v", s.proxyUrl, s.conn.GetStat())
				s.signUp()
			}
			return
		} else if s.conn.GetStat() == tnet.Connecting {
			//连接中
			return
		} else {
			//连接失败，3s重新连接
			if time.Now().Sub(s.connectTime).Seconds() < 3 {
				return
			}
			logger.ERR("connect proxy %v failed, conn stat %v, retry connect", s.proxyUrl, s.conn.GetStat())
			s.disConnect()
			conn, err := tnet.NewConn(s.proxyUrl, 3*time.Second, ParseProxyPkgLen, ProxyCallBackFunc, s)
			if err != nil {
				logger.ERR("new conn failed %v", err)
				return
			}
			s.conn = conn
			s.connectTime = time.Now()
			logger.DEBUG("start connect proxy %s", s.proxyUrl)
			return
		}
	}
}

//发送认证消息
func (s *server) signUp() {
	s.signUpTime = time.Now()
	req := tcaplus_protocol_cs.NewTCaplusPkg()
	//head
	req.Head.Magic = uint16(tcaplus_protocol_cs.TCAPLUS_PROTOCOL_MAGIC_CS)
	req.Head.Version = uint16(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	req.Head.Cmd = cmd.TcaplusApiAppSignUpReq
	req.Head.RouterInfo.AppID = int32(s.appId)
	req.Head.RouterInfo.ZoneID = int32(s.zoneId)
	req.Head.KeyInfo.Version = -1
	req.Body.Init(int64(req.Head.Cmd))

	//body
	req.Body.AppSignupReq.Signature = s.signature
	req.Body.AppSignupReq.Type = 0

	//pack
	if buf, err := req.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion); err != nil {
		logger.ERR("proxy %s signUp pack failed %v", s.proxyUrl, err.Error())
		return
	} else {
		logger.INFO("proxy %s msg:%s signUp pack len %v", s.proxyUrl, common.CsHeadVisualize(req.Head), len(buf))
		go s.conn.Send(buf)
	}
}

//发送心跳
func (s *server) sendHeartbeat() {
	req := tcaplus_protocol_cs.NewTCaplusPkg()
	//head
	req.Head.Magic = uint16(tcaplus_protocol_cs.TCAPLUS_PROTOCOL_MAGIC_CS)
	req.Head.Version = uint16(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	req.Head.Cmd = cmd.TcaplusApiHeartBeatReq
	req.Head.RouterInfo.AppID = int32(s.appId)
	req.Head.RouterInfo.ZoneID = int32(s.zoneId)
	req.Head.KeyInfo.Version = -1
	req.Body.Init(int64(req.Head.Cmd))

	//body
	req.Body.HeartBeatReq.CurTime = uint64(time.Now().Unix())
	req.Body.HeartBeatReq.ApiTimeUs = int64(time.Now().UnixNano()) / 1000

	//pack
	if buf, err := req.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion); err != nil {
		logger.ERR("proxy %s sendHeartbeat pack failed %v", s.proxyUrl, err.Error())
		return
	} else {
		logger.DEBUG("proxy %s msg:%s sendHeartbeat pack len %v", s.proxyUrl, common.CsHeadVisualize(req.Head), len(buf))
		s.conn.Send(buf)
	}
}

//proxy 准备停止的响应
func (s *server) sendStopNotifyRes(asynID uint64) {
	req := tcaplus_protocol_cs.NewTCaplusPkg()
	//head
	req.Head.Magic = uint16(tcaplus_protocol_cs.TCAPLUS_PROTOCOL_MAGIC_CS)
	req.Head.Version = uint16(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	req.Head.Cmd = cmd.TcaplusApiNotifyStopRes
	req.Head.RouterInfo.AppID = int32(s.appId)
	req.Head.RouterInfo.ZoneID = int32(s.zoneId)
	req.Head.KeyInfo.Version = -1
	req.Head.AsynID = asynID
	req.Body.Init(int64(req.Head.Cmd))
	//body empty

	//pack
	if buf, err := req.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion); err != nil {
		logger.ERR("proxy %s sendStopNotifyRes pack failed %v", s.proxyUrl, err.Error())
		return
	} else {
		logger.INFO("proxy %s msg:%s sendStopNotifyRes pack len %v", s.proxyUrl, common.CsHeadVisualize(req.Head), len(buf))
		s.conn.Send(buf)
	}
}

//判断是否收到完整proxy包
//TCaplusPkgHead = Magic(2) + Version(2) + HeadLen(4) + BodyLen(4) = 12
func ParseProxyPkgLen(buf []byte) int {
	if len(buf) >= 12 {
		headLen := int32(0)
		bodyLen := int32(0)
		binary.Read(bytes.NewReader(buf[4:8]), binary.BigEndian, &headLen)
		binary.Read(bytes.NewReader(buf[8:12]), binary.BigEndian, &bodyLen)
		return int(headLen) + int(bodyLen)
	}
	return 0
}

/*
@brief tcp回调函数，解析消息，并触发client协程进行处理
@param url 触发响应包的url
@param  buf 响应包
@param cbPara 回调参数，此处为ProxyServer
@retVal error
*/
func ProxyCallBackFunc(url *string, buf []byte, cbPara interface{}) error {
	asyncId := binary.BigEndian.Uint64(buf[12:])
	seq := binary.BigEndian.Uint32(buf[20:])
	server, ok := cbPara.(*server)
	if !ok {
		logger.ERR("RecvCallBackFunc cbPara type invalid")
		return nil
	}
	id := (int(asyncId)+int(seq))%common.ConfigProcRespRoutineNum
	select {
	case <-server.closeFlag:
		logger.INFO("exit recv routine.")
	case server.respMsgChanList[id] <- buf:
	}
	return nil
}

func (s *server) processRsp(msg *tcaplus_protocol_cs.TCaplusPkg) {
	switch int(msg.Head.Cmd) {
	case cmd.TcaplusApiAppSignUpRes:
		logger.INFO("recv proxy %s response %s", s.proxyUrl, common.CsHeadVisualize(msg.Head))
		if 0 == msg.Head.Result {
			s.signUpFlag = SignUpSuccess
			logger.INFO("zone %d proxy %s signUp success", s.zoneId, s.proxyUrl)
		} else {
			s.signUpFlag = SignUpFail
			logger.ERR("zone %d proxy %s signUp failed, ret %d", s.zoneId, s.proxyUrl, msg.Head.Result)
		}

	case cmd.TcaplusApiNotifyStopReq:
		logger.INFO("recv TcaplusApiNotifyStopReq from %s", s.proxyUrl)
		s.prepareStop = true
		time.AfterFunc(time.Second, func() {
			s.sendStopNotifyRes(msg.Head.AsynID)
		})

	case cmd.TcaplusApiHeartBeatRes:
		curTime := time.Now().UnixNano() / 1000
		cost := curTime - msg.Body.HeartBeatRes.ApiTimeUs
		if cost > 1000*10 {
			sendTm := time.Unix(msg.Body.HeartBeatRes.ApiTimeUs/1000000, 0)
			proxyTm := time.Unix(msg.Body.HeartBeatRes.ProxyTimeUs/1000000, 0)
			logger.WARN("proxy %s Heartbeat delay %d us > 10ms, sendTm:%s:%d, proxyTm:%s:%d",
				s.proxyUrl, cost,
				sendTm.Format("2006-01-02 15:04:05"), msg.Body.HeartBeatRes.ApiTimeUs%1000000,
				proxyTm.Format("2006-01-02 15:04:05"), msg.Body.HeartBeatRes.ProxyTimeUs%1000000)
		} else {
			logger.DEBUG("proxy %s Heartbeat delay %d us", s.proxyUrl, cost)
		}

	case cmd.TcaplusApiInsertRes,
		 cmd.TcaplusApiGetRes,
		 cmd.TcaplusApiUpdateRes,
		 cmd.TcaplusApiReplaceRes,
		 cmd.TcaplusApiDeleteRes,
		 cmd.TcaplusApiBatchGetRes,
		 cmd.TcaplusApiGetByPartkeyRes,
		 cmd.TcaplusApiDeleteByPartkeyRes,
		 cmd.TcaplusApiIncreaseRes,
		 cmd.TcaplusApiListGetAllRes,
		 cmd.TcaplusApiListAddAfterRes,
		 cmd.TcaplusApiListGetRes,
		 cmd.TcaplusApiListDeleteRes,
		 cmd.TcaplusApiListDeleteAllRes,
		 cmd.TcaplusApiListReplaceRes,
		 cmd.TcaplusApiListDeleteBatchRes,
		 cmd.TcaplusApiSqlRes,
		 cmd.TcaplusApiMetadataGetRes,
		 cmd.TcaplusApiPBFieldGetRes,
		 cmd.TcaplusApiPBFieldUpdateRes,
		 cmd.TcaplusApiPBFieldIncreaseRes,
		 cmd.TcaplusApiGetShardListRes,
		 cmd.TcaplusApiTableTraverseRes,
		 cmd.TcaplusApiGetTableRecordCountRes:
		 	if logger.LogConf.LogLevel == "DEBUG" {
				logger.DEBUG("recv proxy %s response %s", s.proxyUrl, common.CsHeadVisualize(msg.Head))
			}
		router := s.router.(*Router)
		if msg.Head.Cmd == cmd.TcaplusApiTableTraverseRes || msg.Head.Cmd == cmd.TcaplusApiGetShardListRes {
			drop := false
			router.TM.OnRecvResponse(s.zoneId, msg, &drop)
			if drop {
				return
			}
		}
		router.processRouterMsg(msg)
		logger.DEBUG("recv proxy response finish")

	default:
		logger.ERR("zone %d proxy %s msg %s invalid", s.zoneId, s.proxyUrl, common.CsHeadVisualize(msg.Head))
	}
}
