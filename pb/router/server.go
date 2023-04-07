package router

import (
	"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
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
	lastRspTime time.Time
	router      *Router
	prepareStop bool //proxy准备stop
	error       error
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
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

func (s *server) send(data []byte) error {
	logger.DEBUG("send to proxy %s", s.proxyUrl)
	if s.conn != nil {
		err := s.conn.Send(data)
		return err
	}
	logger.ERR("proxy svr %s conn is empty", s.proxyUrl)
	return &terror.ErrorCode{Code: terror.SendRequestFail, Message: "proxy con is nil"}
}

func (s *server) connect() {
	if s.conn == nil {
		//连接proxy, 3s超时
		conn, err := tnet.NewConn(s.proxyUrl, s.router.ctrl.Option.ProxyConnOption.ConTimeout, ParseProxyPkgLen,
			ProxyCallBackFunc, s,
			s.router.ctrl.Option.ProxyConnOption.BufSizePerCon)
		if err != nil {
			s.error = &terror.ErrorCode{Code: terror.API_ERR_PROXY_CONNECT_FAILED, Message: "connect proxy failed:" + s.proxyUrl}
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
			} else if s.signUpFlag != SignUpSuccess && time.Now().Sub(s.signUpTime).Seconds() > 3 {
				s.error = &terror.ErrorCode{Code: terror.API_ERR_DIR_SIGNUP_FAILED, Message: "sign up proxy timeout(3s):" + s.proxyUrl}
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
			s.error = &terror.ErrorCode{Code: terror.API_ERR_PROXY_CONNECT_FAILED, Message: "connect proxy timeout(3s):" + s.proxyUrl}
			logger.ERR("connect proxy %v failed, conn stat %v, retry connect", s.proxyUrl, s.conn.GetStat())
			s.disConnect()
			conn, err := tnet.NewConn(s.proxyUrl, s.router.ctrl.Option.ProxyConnOption.ConTimeout, ParseProxyPkgLen, ProxyCallBackFunc, s,
				s.router.ctrl.Option.ProxyConnOption.BufSizePerCon)
			if err != nil {
				s.error = &terror.ErrorCode{Code: terror.API_ERR_PROXY_CONNECT_FAILED, Message: "connect proxy failed:" + s.proxyUrl}
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
		s.send(buf)
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
		s.send(buf)
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
		s.send(buf)
	}
}

//判断是否收到完整proxy包
//TCaplusPkgHead = Magic(2) + Version(2) + HeadLen(4) + BodyLen(4) = 12
func ParseProxyPkgLen(buf []byte) int {
	if len(buf) >= 12 {
		headLen := binary.BigEndian.Uint32(buf[4:8])
		bodyLen := binary.BigEndian.Uint32(buf[8:12])
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
func ProxyCallBackFunc(url *string, pkg *tnet.PKG) error {
	server, ok := pkg.GetCbPara().(*server)
	if !ok {
		logger.ERR("url %s RecvCallBackFunc cbPara type invalid", *url)
		return nil
	}
	server.router.ResponseChanAdd(pkg)
	return nil
}

func (s *server) processRsp(msg *tcaplus_protocol_cs.TCaplusPkg) {
	s.lastRspTime = time.Now()
	switch int(msg.Head.Cmd) {
	case cmd.TcaplusApiAppSignUpRes:
		logger.INFO("recv proxy %s response %s", s.proxyUrl, common.CsHeadVisualize(msg.Head))
		if 0 == msg.Head.Result {
			s.signUpFlag = SignUpSuccess
			logger.INFO("zone %d proxy %s signUp success", s.zoneId, s.proxyUrl)
		} else {
			s.signUpFlag = SignUpFail
			s.error = &terror.ErrorCode{Code: int(msg.Head.Result)}
			logger.ERR("zone %d proxy %s signUp failed, ret %d", s.zoneId, s.proxyUrl, msg.Head.Result)
		}

	case cmd.TcaplusApiNotifyStopReq:
		logger.INFO("recv TcaplusApiNotifyStopReq from %s", s.proxyUrl)
		s.prepareStop = true
		s.sendStopNotifyRes(msg.Head.AsynID)

	case cmd.TcaplusApiHeartBeatRes:
		curTime := s.lastRspTime.UnixNano() / 1000
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
		cmd.TcaplusApiListTableTraverseRes,
		cmd.TcaplusApiGetTableRecordCountRes,
		cmd.TcaplusApiSetTtlRes,
		cmd.TcaplusApiGetTtlRes,
		cmd.TcaplusApiBatchDeleteRes,
		cmd.TcaplusApiBatchInsertRes,
		cmd.TcaplusApiBatchReplaceRes,
		cmd.TcaplusApiBatchUpdateRes,
		cmd.TcaplusApiListAddAfterBatchRes,
		cmd.TcaplusApiListGetBatchRes,
		cmd.TcaplusApiListReplaceBatchRes:
		if logger.GetLogLevel() == "DEBUG" {
			s.printRsp(msg)
		}
		router := s.router
		if msg.Head.Cmd == cmd.TcaplusApiTableTraverseRes || msg.Head.Cmd == cmd.TcaplusApiListTableTraverseRes ||
			msg.Head.Cmd == cmd.TcaplusApiGetShardListRes {
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

func (s *server) printRsp(msg *tcaplus_protocol_cs.TCaplusPkg) {
	logger.DEBUG("recv proxy %s response %s", s.proxyUrl, common.CsHeadVisualize(msg.Head))
	switch int(msg.Head.Cmd) {
	case cmd.TcaplusApiAppSignUpRes:
		logger.DEBUG("recv proxy %s TcaplusApiAppSignUpRes %s", s.proxyUrl, common.CovertToJson(msg.Body.AppSignupRes))
	case cmd.TcaplusApiNotifyStopReq:
		logger.DEBUG("recv proxy %s TcaplusApiNotifyStopReq", s.proxyUrl)
	case cmd.TcaplusApiInsertRes:
		logger.DEBUG("recv proxy %s TcaplusApiInsertRes %s", s.proxyUrl, common.CovertToJson(msg.Body.InsertRes))
	case cmd.TcaplusApiGetRes:
		logger.DEBUG("recv proxy %s TcaplusApiGetRes %s", s.proxyUrl, common.CovertToJson(msg.Body.GetRes))
	case cmd.TcaplusApiUpdateRes:
		logger.DEBUG("recv proxy %s TcaplusApiUpdateRes %s", s.proxyUrl, common.CovertToJson(msg.Body.UpdateRes))
	case cmd.TcaplusApiReplaceRes:
		logger.DEBUG("recv proxy %s TcaplusApiReplaceRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ReplaceRes))
	case cmd.TcaplusApiDeleteRes:
		logger.DEBUG("recv proxy %s TcaplusApiDeleteRes %s", s.proxyUrl, common.CovertToJson(msg.Body.DeleteRes))
	case cmd.TcaplusApiBatchGetRes:
		logger.DEBUG("recv proxy %s TcaplusApiBatchGetRes %s", s.proxyUrl, common.CovertToJson(msg.Body.BatchGetRes))
	case cmd.TcaplusApiGetByPartkeyRes:
		logger.DEBUG("recv proxy %s TcaplusApiGetByPartkeyRes %s", s.proxyUrl, common.CovertToJson(msg.Body.GetByPartKeyRes))
	case cmd.TcaplusApiDeleteByPartkeyRes:
		logger.DEBUG("recv proxy %s TcaplusApiDeleteByPartkeyRes %s", s.proxyUrl, common.CovertToJson(msg.Body.DeleteByPartkeyRes))
	case cmd.TcaplusApiIncreaseRes:
		logger.DEBUG("recv proxy %s TcaplusApiIncreaseRes %s", s.proxyUrl, common.CovertToJson(msg.Body.IncreaseRes))
	case cmd.TcaplusApiListGetAllRes:
		logger.DEBUG("recv proxy %s TcaplusApiListGetAllRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListGetAllRes))
	case cmd.TcaplusApiListAddAfterRes:
		logger.DEBUG("recv proxy %s TcaplusApiListAddAfterRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListAddAfterRes))
	case cmd.TcaplusApiListGetRes:
		logger.DEBUG("recv proxy %s TcaplusApiListGetRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListGetRes))
	case cmd.TcaplusApiListDeleteRes:
		logger.DEBUG("recv proxy %s TcaplusApiListDeleteRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListDeleteRes))
	case cmd.TcaplusApiListDeleteAllRes:
		logger.DEBUG("recv proxy %s TcaplusApiListDeleteAllRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListDeleteAllRes))
	case cmd.TcaplusApiListReplaceRes:
		logger.DEBUG("recv proxy %s TcaplusApiListReplaceRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListReplaceRes))
	case cmd.TcaplusApiListDeleteBatchRes:
		logger.DEBUG("recv proxy %s TcaplusApiListDeleteBatchRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListDeleteBatchRes))
	case cmd.TcaplusApiSqlRes:
		logger.DEBUG("recv proxy %s TcaplusApiSqlRes %s", s.proxyUrl, common.CovertToJson(msg.Body.TCaplusSqlRes))
	case cmd.TcaplusApiMetadataGetRes:
		logger.DEBUG("recv proxy %s TcaplusApiMetadataGetRes %s", s.proxyUrl, common.CovertToJson(msg.Body.MetadataGetRes))
	case cmd.TcaplusApiPBFieldGetRes:
		logger.DEBUG("recv proxy %s TcaplusApiPBFieldGetRes %s", s.proxyUrl, common.CovertToJson(msg.Body.TCaplusPbFieldGetRes))
	case cmd.TcaplusApiPBFieldUpdateRes:
		logger.DEBUG("recv proxy %s TcaplusApiPBFieldUpdateRes %s", s.proxyUrl, common.CovertToJson(msg.Body.TCaplusPbFieldUpdateRes))
	case cmd.TcaplusApiPBFieldIncreaseRes:
		logger.DEBUG("recv proxy %s TcaplusApiPBFieldIncreaseRes %s", s.proxyUrl, common.CovertToJson(msg.Body.TCaplusPbFieldIncRes))
	case cmd.TcaplusApiGetShardListRes:
		logger.DEBUG("recv proxy %s TcaplusApiGetShardListRes %s", s.proxyUrl, common.CovertToJson(msg.Body.GetShardListRes))
	case cmd.TcaplusApiTableTraverseRes:
		logger.DEBUG("recv proxy %s TcaplusApiTableTraverseRes %s", s.proxyUrl, common.CovertToJson(msg.Body.TableTraverseRes))
	case cmd.TcaplusApiListTableTraverseRes:
		logger.DEBUG("recv proxy %s TcaplusApiListTableTraverseRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListTableTraverseRes))
	case cmd.TcaplusApiGetTableRecordCountRes:
		logger.DEBUG("recv proxy %s TcaplusApiGetTableRecordCountRes %s", s.proxyUrl, common.CovertToJson(msg.Body.GetTableRecordCountRes))
	case cmd.TcaplusApiSetTtlRes:
		logger.DEBUG("recv proxy %s TcaplusApiSetTtlRes %s", s.proxyUrl, common.CovertToJson(msg.Body.TCaplusSetTTLRes))
	case cmd.TcaplusApiGetTtlRes:
		logger.DEBUG("recv proxy %s TcaplusApiGetTtlRes %s", s.proxyUrl, common.CovertToJson(msg.Body.TCaplusGetTTLRes))
	case cmd.TcaplusApiBatchDeleteRes:
		logger.DEBUG("recv proxy %s TcaplusApiBatchDeleteRes %s", s.proxyUrl, common.CovertToJson(msg.Body.BatchDeleteRes))
	case cmd.TcaplusApiBatchInsertRes:
		logger.DEBUG("recv proxy %s TcaplusApiBatchInsertRes %s", s.proxyUrl, common.CovertToJson(msg.Body.BatchInsertRes))
	case cmd.TcaplusApiBatchReplaceRes:
		logger.DEBUG("recv proxy %s TcaplusApiBatchReplaceRes %s", s.proxyUrl, common.CovertToJson(msg.Body.BatchReplaceRes))
	case cmd.TcaplusApiBatchUpdateRes:
		logger.DEBUG("recv proxy %s TcaplusApiBatchUpdateRes %s", s.proxyUrl, common.CovertToJson(msg.Body.BatchUpdateRes))
	case cmd.TcaplusApiListAddAfterBatchRes:
		logger.DEBUG("recv proxy %s TcaplusApiListAddAfterBatchRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListAddAfterBatchRes))
	case cmd.TcaplusApiListGetBatchRes:
		logger.DEBUG("recv proxy %s TcaplusApiListGetBatchRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListGetBatchRes))
	case cmd.TcaplusApiListReplaceBatchRes:
		logger.DEBUG("recv proxy %s TcaplusApiListReplaceBatchRes %s", s.proxyUrl, common.CovertToJson(msg.Body.ListReplaceBatchRes))
	}
}
