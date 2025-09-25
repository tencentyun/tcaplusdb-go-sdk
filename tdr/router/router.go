package router

import (
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/config"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/dir"
	tcaplusCmd "github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/statistics"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/subscribe"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/tnet"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/traverser"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcapdir_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

// 同步请求结构体
type SyncRequest struct {
	//sync response msg chan
	syncMsgPipe chan *tcaplus_protocol_cs.TCaplusPkg
	//request package
	requestPkg  request.TcaplusRequest
	delFromMap  byte //0 添加，1 删除，2 收到响应后自动删除
	syncId      uint32
	proxyServer *server      // 发送到了哪个个proxyServer
	err         atomic.Value // 错误
}

func (S *SyncRequest) SetErr(err error) {
	S.err.Store(err)
}

func (S *SyncRequest) GetErr() error {
	return S.err.Load().(error)
}

func (S *SyncRequest) InitTraverseChan(seq, num int32) {
	S.syncMsgPipe = make(chan *tcaplus_protocol_cs.TCaplusPkg, num)
	S.syncId = uint32(seq)
	S.delFromMap = 0
}

func (S *SyncRequest) InitMoreChan(reqpkg request.TcaplusRequest, num int32) {
	S.syncMsgPipe = make(chan *tcaplus_protocol_cs.TCaplusPkg, num)
	S.requestPkg = reqpkg
	S.syncId = uint32(reqpkg.GetSeq())
	S.delFromMap = 0
}

func (S *SyncRequest) Init(reqpkg request.TcaplusRequest) {
	S.syncMsgPipe = make(chan *tcaplus_protocol_cs.TCaplusPkg, 1)
	S.requestPkg = reqpkg
	S.syncId = uint32(reqpkg.GetSeq())
	S.delFromMap = 2
}

func (S *SyncRequest) GetSyncChan() chan *tcaplus_protocol_cs.TCaplusPkg {
	return S.syncMsgPipe
}

func (S *SyncRequest) SyncChanClose() {
	close(S.syncMsgPipe)
}

func (S *SyncRequest) IncTimeoutStatistic() {
	if S.proxyServer != nil {
		S.proxyServer.IncTimeoutStatistic()
	}
}

func (S *SyncRequest) GetProxyUrl() string {
	if S.proxyServer != nil {
		return S.proxyServer.GetUrl()
	}
	return ""
}

// 用户请求路由结构体
type Router struct {
	appId     uint64
	zoneList  []uint32
	signature string
	//zone-->proxy
	proxyMap map[uint32]*proxy
	MsgPipe  chan *tcaplus_protocol_cs.TCaplusPkg //用户proxy消息通道

	//心跳时间间隔s
	heartbeatInterval time.Duration
	lastHeartbeatTime time.Time
	lastStatPrintTime time.Time

	//res msg queue
	respCount    int64
	respMsgMutex sync.Mutex
	respMsgQueue *list.List

	//req chan map

	TM   *traverser.TraverserManager
	ctrl *config.ClientCtrl

	chanClose chan struct{}
	//打包协程
	requestRoutineNum int
	requestChanMap    []map[uint32]*SyncRequest
	requestChanList   []chan interface{}

	//解包协程
	responseRoutineNum int
	responseChanList   []chan *tnet.PKG

	// user + passwd
	userName string
	passwd   string

	statOutPutInterface statistics.StatInterface // 输出统计到自定义接口
	dirReportStat       *dir.ReportStat          // dir上报统计

	// 订阅管理器
	subscribeMgr subscribe.Manager
}

func (r *Router) SetDirReportStat(dirReportStat *dir.ReportStat) {
	r.dirReportStat = dirReportStat
}

func (r *Router) SetStatisticsOutInterface(statInterface statistics.StatInterface) {
	r.statOutPutInterface = statInterface
}

func (r *Router) SetUserNameAndPassword(userName string, password string) {
	r.userName = userName
	r.passwd = password
}

func (r *Router) createRequestRoutine() {
	if r.ctrl.Option.PackRoutineCount > 0 {
		r.requestRoutineNum = r.ctrl.Option.PackRoutineCount
	} else {
		r.requestRoutineNum = common.ConfigProcReqRoutineNum
	}
	r.requestChanMap = make([]map[uint32]*SyncRequest, r.requestRoutineNum)
	r.requestChanList = make([]chan interface{}, r.requestRoutineNum)
	r.chanClose = make(chan struct{})
	for i := 0; i < r.requestRoutineNum; i++ {
		r.requestChanMap[i] = make(map[uint32]*SyncRequest)
		r.requestChanList[i] = make(chan interface{}, common.ConfigProcReqDepth)
		r.ctrl.Add(1)
		go func(id int) {
			defer r.ctrl.Done()
			ch := r.requestChanList[id]
			rmap := r.requestChanMap[id]
			proc := func(p interface{}) {
				if req, ok := p.(*SyncRequest); ok {
					if req.delFromMap == 1 {
						delete(rmap, req.syncId)
					} else {
						if req.requestPkg != nil {
							err, svr := r.SendRequest(req.requestPkg)
							if err != nil {
								logger.ERR("Send failed %s", err.Error())
								req.SetErr(err)
								req.SyncChanClose()
								return
							}
							req.proxyServer = svr
						}
						rmap[req.syncId] = req
					}
				} else {
					msg := p.(*tcaplus_protocol_cs.TCaplusPkg)
					if v, exist := rmap[uint32(msg.Head.Seq)]; exist {
						v.syncMsgPipe <- msg
						if v.delFromMap == 2 {
							delete(rmap, uint32(msg.Head.Seq))
						}
					} else {
						logger.ERR("Can not find request chan %d", msg.Head.Seq)
					}
				}
			}

			for {
				select {
				case <-r.chanClose:
					// 退出前处理完管道中的包
					for len(ch) > 0 {
						for i := 0; i < len(ch); i++ {
							proc(<-ch)
						}
					}
					logger.INFO("processSyncOperate exit")
					return
				case p := <-ch:
					proc(p)
				}
			}
		}(i)
	}
}

func (r *Router) createResponseRoutine() {
	if r.ctrl.Option.UnPackRoutineCount > 0 {
		r.responseRoutineNum = r.ctrl.Option.UnPackRoutineCount
	} else {
		r.responseRoutineNum = common.ConfigProcRespRoutineNum
	}
	r.responseChanList = make([]chan *tnet.PKG, r.responseRoutineNum)
	for i := 0; i < r.responseRoutineNum; i++ {
		r.responseChanList[i] = make(chan *tnet.PKG, common.ConfigProcRespDepth)
		r.ctrl.Add(1)
		go func(ch chan *tnet.PKG) {
			defer r.ctrl.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.DEBUG("Recovered %v.", r)
				}
			}()
			var start time.Time
			proc := func(pkg *tnet.PKG) {
				buf := pkg.GetData()
				server, ok := pkg.GetCbPara().(*server)
				if !ok {
					logger.ERR("Recv pkg cbPara type invalid")
					pkg.Done()
					return
				}
				if logger.GetLogLevel() == "DEBUG" {
					start = time.Now()
				}
				resp := tcaplus_protocol_cs.NewTCaplusPkg()
				err := resp.Unpack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion, buf)
				// 之后这个pkg不会再被用到，回收到对象池中
				pkg.Done()
				if err != nil {
					logger.ERR("Unpack proxy msg failed, url %s err %v", server.proxyUrl, err.Error())
					return
				}

				if logger.GetLogLevel() == "DEBUG" {
					if time.Now().Sub(start) > 10*time.Millisecond {
						logger.WARN("unpack > 10ms data %v.", buf)
					}
				}
				server.processRsp(resp)
			}

			for {
				select {
				case <-r.chanClose:
					for len(ch) > 0 {
						for i := 0; i < len(ch); i++ {
							proc(<-ch)
						}
					}
					logger.INFO("exit recv routine.")
					return
				case buf := <-ch:
					proc(buf)
				}
			}
		}(r.responseChanList[i])
	}
}

//此处的是用户pkg，非用户pkg已在func (s *server)processRsp回调中处理
func (r *Router) processRouterMsg(msg *tcaplus_protocol_cs.TCaplusPkg) {
	if common.HasFlag(msg.Head.Flags, option.TcaplusFlagSubscribe) {
		r.ProcSubscribeRes(msg)
		return
	}

	if msg.Head.Seq == 0 {
		r.respMsgMutex.Lock()
		defer r.respMsgMutex.Unlock()
		r.respMsgQueue.PushBack(msg)
		atomic.AddInt64(&r.respCount, 1)
		logger.DEBUG("add one msg to queue, %d", r.respCount)
		return
	}

	select {
	case <-r.chanClose:
		logger.INFO("processSyncOperate exit")
	case r.requestChanList[uint32(msg.Head.Seq)%uint32(r.requestRoutineNum)] <- msg:
	}
}

func (r *Router) RecvResponse() (response.TcaplusResponse, error) {
	if atomic.LoadInt64(&r.respCount) <= 0 {
		return nil, nil
	}
	logger.DEBUG("pop one msg from queue, %d", r.respCount)

	r.respMsgMutex.Lock()
	defer r.respMsgMutex.Unlock()
	ele := r.respMsgQueue.Front()
	if ele != nil {
		pkg := ele.Value.(*tcaplus_protocol_cs.TCaplusPkg)
		r.respMsgQueue.Remove(ele)
		atomic.AddInt64(&r.respCount, -1)
		return response.NewResponse(pkg)
	}
	return nil, nil
}

func (r *Router) Init(appId uint64, zoneList []uint32, signature string, ctrl *config.ClientCtrl) error {
	r.appId = appId
	r.zoneList = make([]uint32, len(zoneList))
	copy(r.zoneList, zoneList)
	r.signature = signature
	r.MsgPipe = make(chan *tcaplus_protocol_cs.TCaplusPkg, 1024)
	r.proxyMap = make(map[uint32]*proxy)
	r.heartbeatInterval = 1
	r.lastHeartbeatTime = time.Now()
	r.lastStatPrintTime = time.Now()
	r.respMsgQueue = list.New()
	r.ctrl = ctrl
	r.subscribeMgr.Init(ctrl)
	//放最后
	r.createRequestRoutine()
	r.createResponseRoutine()
	r.startSubscribeRoutine()
	return nil
}

func (r *Router) RequestChanMapAdd(syncRequest *SyncRequest) int {
	syncId := syncRequest.syncId % uint32(r.requestRoutineNum)
	select {
	case <-r.chanClose:
		logger.INFO("processSyncOperate exit")
		return -1
	case r.requestChanList[syncId] <- syncRequest:
	}
	return 0
}

func (r *Router) ResponseChanAdd(pkg *tnet.PKG) {
	buf := pkg.GetData()
	asyncId := binary.BigEndian.Uint64(buf[12:])
	seq := binary.BigEndian.Uint32(buf[20:])
	id := (asyncId + uint64(seq)) % uint64(r.responseRoutineNum)
	select {
	case <-r.chanClose:
		logger.INFO("processSyncOperate exit")
		return
	case r.responseChanList[id] <- pkg:
	}
}

func (r *Router) RequestChanMapClean(syncRequest *SyncRequest) int {
	syncId := syncRequest.syncId % uint32(r.requestRoutineNum)
	var synReq SyncRequest
	synReq.syncId = syncRequest.syncId
	synReq.delFromMap = 1
	select {
	case <-r.chanClose:
		logger.INFO("processSyncOperate exit")
		return -1
	case r.requestChanList[syncId] <- &synReq:
	}
	return 0
}

func (r *Router) SetHeartbeatInterval(heartbeatInterval time.Duration) {
	r.heartbeatInterval = heartbeatInterval
}

func (r *Router) SetPerfPercent(zone int32, percent int32) {
	if p, exist := r.proxyMap[uint32(zone)]; exist {
		atomic.StoreInt32(&p.perfPercent, percent)
		logger.INFO("zone %d perfPercent is %d", zone, percent)
	}
}

func (r *Router) CheckTable(zoneId uint32, tableName string) error {
	if proxy, exist := r.proxyMap[zoneId]; !exist {
		return &terror.ErrorCode{Code: terror.ZoneIdNotExist,
			Message: fmt.Sprintf("zone %d not exit", zoneId)}
	} else {
		proxy.tbMutex.RLock()
		defer proxy.tbMutex.RUnlock()
		if _, exist := proxy.tableNameList[tableName]; !exist {
			return &terror.ErrorCode{Code: terror.TableNotExist,
				Message: fmt.Sprintf("zone %d table %s not exit", zoneId, tableName)}
		}
	}

	return nil
}

func (r *Router) GetZoneTables(zoneId uint32) []string {
	var tables []string
	proxy, exist := r.proxyMap[zoneId]
	if exist {
		tables = make([]string, 0, len(proxy.tableNameList))
		proxy.tbMutex.RLock()
		for table := range proxy.tableNameList {
			tables = append(tables, table)
		}
		proxy.tbMutex.RUnlock()
	}
	return tables
}

//初始化超时，获取dir的error信息
func (r *Router) GetError() error {
	var errStr string
	for _, proxy := range r.proxyMap {
		proxyErr := proxy.GetErrorStr()
		if len(proxyErr) > 0 {
			errStr = errStr + proxyErr + ","
		}
	}

	if len(errStr) > 0 {
		return errors.New(errStr)
	}

	if len(r.proxyMap) == 0 {
		return &terror.ErrorCode{Code: terror.API_ERR_DIR_GET_PROXYLIST_TIMEOUT}
	}
	return nil
}

//0 所有认证成功， 1 有proxy全部认证中， 2 所有proxy部分认证成功，可以启动， -1 有认证失败,启动失败
func (r *Router) CanStartUp() (int, error) {
	if len(r.zoneList) != len(r.proxyMap) {
		return 1, nil
	}

	sucCount := 0
	partSucCount := 0
	for _, proxy := range r.proxyMap {
		ret, err := proxy.CheckAvailable()
		if ret == 0 {
			sucCount++
			partSucCount++
		} else if ret == 2 {
			partSucCount++
		} else if ret == -1 {
			return -1, err
		}
	}

	if sucCount == len(r.proxyMap) {
		return 0, nil
	}

	if partSucCount == len(r.proxyMap) {
		return 2, nil
	}

	return 1, nil
}

func (r *Router) Update() {
	curTime := time.Now()
	for _, proxy := range r.proxyMap {
		proxy.update(curTime)
	}
	//发送心跳
	diff := curTime.Sub(r.lastHeartbeatTime)
	if diff > r.heartbeatInterval*time.Second {
		for _, proxy := range r.proxyMap {
			proxy.sendHeartbeat()
		}
		r.lastHeartbeatTime = curTime
	}

	diff = curTime.Sub(r.lastStatPrintTime)
	if diff > 60*time.Second {
		for _, proxy := range r.proxyMap {
			proxy.statReset()
		}
		r.lastStatPrintTime = curTime
	}
}

func (r *Router) ProcessTablesAndAccessMsg(msg *tcapdir_protocol_cs.ResGetTablesAndAccess) {
	if nil == msg {
		return
	}

	if p, exist := r.proxyMap[uint32(msg.ZoneID)]; exist {
		p.processTablesAndAccessMsg(msg)
	} else {
		p := &proxy{
			zoneId:            uint32(msg.ZoneID),
			appId:             r.appId,
			signature:         r.signature,
			router:            r,
			tableNameList:     make(map[string]bool),
			hashList:          make([]*server, 0, 10),
			usingServerList:   make(map[string]*server),
			prepareServerList: make(map[string]bool),
			removeServerList:  make(map[string]*server),
			lastIsolateTime:   time.Now(),
			isolateProxyRate:  10,
		}
		p.statMgr.Init(uint32(p.appId), p.zoneId, r.statOutPutInterface, r.dirReportStat)
		p.processTablesAndAccessMsg(msg)
		r.proxyMap[uint32(msg.ZoneID)] = p
	}
}

func (r *Router) SendRequest(req request.TcaplusRequest) (error, *server) {
	if p, exist := r.proxyMap[req.GetZoneId()]; exist {
		cmd := req.GetCmd()
		if req.GetSeq()%100 < atomic.LoadInt32(&p.perfPercent) && tcaplusCmd.IsSimpleReq(cmd) {
			req.SetPerfTest(uint64(time.Now().UnixMicro()))
		}
		data, err := req.Pack()
		if err != nil {
			logger.ERR("req pack failed %s", err.Error())
			return err, nil
		}
		//获取keyHash
		code, err := req.GetKeyHash()
		if err != nil {
			logger.ERR("get key hash failed %s", err.Error())
			return err, nil
		}
		return p.send(cmd, code, data)
	}
	logger.ERR("zone %d not connect", req.GetZoneId())
	return &terror.ErrorCode{Code: terror.SendRequestFail, Message: "zone proxy not connect"}, nil
}

func (r *Router) GetProxyUrl(hashCode, zoneId uint32) string {
	p, exist := r.proxyMap[zoneId]
	if !exist {
		return ""
	}
	p.hashMutex.RLock()
	defer p.hashMutex.RUnlock()
	id := hashCode % uint32(len(p.hashList))
	preId := id
	for {
		svr := p.hashList[id]
		if svr.isAvailable() {
			return svr.proxyUrl
		}

		//选择下个节点
		hashCode++
		id = hashCode % uint32(len(p.hashList))
		//一轮之后
		if id == preId {
			return ""
		}
	}
}

func (r *Router) Close() {
	for _, v := range r.proxyMap {
		if v == nil {
			continue
		}
		v.prepareServerList = nil
		for _, svr := range v.usingServerList {
			if svr == nil {
				continue
			}
			svr.disConnect()
		}
		for _, svr := range v.removeServerList {
			if svr == nil {
				continue
			}
			svr.disConnect()
		}
	}
	close(r.chanClose)
}

func (r *Router) UpdateHashList() {
	for _, p := range r.proxyMap {
		p.updateHashList()
	}
}

func (r *Router) NewSubscribeTopic(zoneId uint32, table string, keyMap map[string][]byte,
	index int32, expireSecond int32, rspPipeSize int, flag int, isPb bool) (*subscribe.Topic, error) {

	if err := r.CheckTable(zoneId, table); err != nil {
		return nil, err
	}

	topic, err := r.subscribeMgr.NewSubscribeTopic(zoneId, table, keyMap, index, expireSecond, rspPipeSize, nil, flag, isPb)
	if err != nil {
		return nil, err
	}

	// start subscribe
	if err := r.StartSubscribe(topic); err != nil {
		r.subscribeMgr.DelSubscribeTopic(topic)
		return nil, err
	}
	return topic, nil
}

func (r *Router) DelSubscribe(topic *subscribe.Topic) {
	r.subscribeMgr.DelSubscribeTopic(topic)
	svr := r.subscribeMgr.GetTopicSvr(topic)
	if svr == nil {
		return
	}
	proxySvr := svr.(*server)
	if !proxySvr.isAvailable() {
		return
	}

	req, err := request.NewRequest(r.appId, topic.GetZoneId(), topic.GetTableName(), tcaplusCmd.TcaplusApiListGetAllReq, topic.IsPb())
	if err != nil {
		return
	}
	// 设置为取消订阅模式
	if ret := req.SetSubscribe(0, topic.GetExpireIndex(), false, uint32(topic.GetRspFlag())); ret != 0 {
		return
	}
	req.SetAsyncId(uint64(topic.GetAsyncId()))
	rec, err := req.AddRecord(0)
	if err != nil {
		return
	}
	topic.Copy2KeyMap(rec.KeyMap)

	data, err := req.Pack()
	if err != nil {
		return
	}
	err = proxySvr.send(req.GetCmd(), data)
	if err != nil {
		return
	}
}

func (r *Router) StartSubscribe(topic *subscribe.Topic) error {
	req, err := request.NewRequest(r.appId, topic.GetZoneId(), topic.GetTableName(), tcaplusCmd.TcaplusApiListGetAllReq, topic.IsPb())
	if err != nil {
		return err
	}
	// 设置为订阅模式
	if ret := req.SetSubscribe(topic.GetExpireSec(), topic.GetExpireIndex(), false, uint32(topic.GetRspFlag())); ret != 0 {
		return terror.MakeError(int(ret), fmt.Sprintf("req.SetSubscribe failed"))
	}

	req.SetAsyncId(uint64(topic.GetAsyncId()))
	rec, err := req.AddRecord(0)
	if err != nil {
		return err
	}

	topic.Copy2KeyMap(rec.KeyMap)
	err, svr := r.SendRequest(req)
	if err != nil {
		return err
	}
	r.subscribeMgr.SetTopicSvr(topic, svr)
	r.subscribeMgr.SetTopicSendMs(topic, time.Now().UnixMilli())
	return nil
}

func (r *Router) ProcSubscribeRes(msg *tcaplus_protocol_cs.TCaplusPkg) {
	asyncId := msg.Head.AsynID
	if msg.Head.Cmd != tcaplusCmd.TcaplusApiListGetAllRes {
		logger.DEBUG("Subscribe rsp cmd %d != TcaplusApiListGetAllRes, drop it", msg.Head.Cmd)
		return
	}
	// find topic
	topic := r.subscribeMgr.FindTopicById(int64(asyncId))
	if topic == nil {
		logger.DEBUG("topic asyncId %d not found", asyncId)
		return
	}
	curMs := time.Now().UnixMilli()
	r.subscribeMgr.SetTopicRecvMs(topic, curMs)

	// 创建响应
	res, err := response.NewResponse(msg)
	if err != nil {
		logger.ERR("topic %s response.NewResponse err %s", topic.GetKeyDebug(), err.Error())
		return
	}

	expireIdx := topic.GetExpireIndex()
	minIdx := msg.Body.ListGetAllRes.EmptyIndexFlag
	maxIdx := int32(msg.Body.ListGetAllRes.BiggestIdx)

	if msg.Head.Result != 0 {
		// 返回错误 取消订阅
		logger.ERR("topic %s expire index %d res idx %d - %d return err %d",
			topic.GetKeyDebug(), expireIdx, minIdx, maxIdx, msg.Head.Result)
		select {
		case topic.ResponsePipe <- res:
			r.subscribeMgr.DelSubscribeTopic(topic)
			return
		default:
			logger.ERR("topic %s channel full, len %d cap %d, drop last pkg, del topic", topic.GetKeyDebug(), len(topic.ResponsePipe), cap(topic.ResponsePipe))
			r.subscribeMgr.DelSubscribeTopic(topic)
			return
		}
	}

	if minIdx < 0 {
		return
	}

	// check idx
	if expireIdx >= 0 && minIdx != expireIdx {
		logger.ERR("topic %s expireIdx %d res idx %d - %d, maybe lost pkg, retry subscribe", topic.GetKeyDebug(), expireIdx, minIdx, maxIdx)
		// 包编号不对说明丢包, 修改topic的id，历史的消息不会再收到，直接丢弃
		// 下个订阅周期，重新从下个index开始订阅，
		r.subscribeMgr.IncTopicAsyncId(topic)
		r.subscribeMgr.SetTopicCheckFlag(topic, subscribe.RetrySubscribeWithIndex)
		return
	}

	if topic.GetRspFlag() == subscribe.OnlyIndexFlag {
		// 不需要记录只用更新index
		if topic.GetFirstIndex() == -1 {
			r.subscribeMgr.SetTopicFirstIndex(topic, minIdx)
		}
		r.subscribeMgr.SetTopicLastIndex(topic, maxIdx)
		r.subscribeMgr.SetTopicExpireIndex(topic, common.GetSubscribeInt32Next(maxIdx))
		r.subscribeMgr.SetTopicLastChangeIdxMs(topic, curMs)
		return
	}

	// 空记录，drop
	if res.GetRecordCount() <= 0 {
		return
	}

	// success
	select {
	case topic.ResponsePipe <- res:
		r.subscribeMgr.SetTopicExpireIndex(topic, common.GetSubscribeInt32Next(maxIdx))
		r.subscribeMgr.SetTopicLastChangeIdxMs(topic, curMs)
		return
	default:
		// 队列满导致丢包
		// 下个订阅周期，重新从下个index开始订阅，
		r.subscribeMgr.IncTopicAsyncId(topic)
		r.subscribeMgr.SetTopicCheckFlag(topic, subscribe.RetrySubscribeWithIndexOnlyCheck)
		logger.INFO("topic %s channel full, len %d cap %d, drop pkg", topic.GetKeyDebug(), len(topic.ResponsePipe), cap(topic.ResponsePipe))
	}
}

func (r *Router) retrySubscribeTopic(topic *subscribe.Topic) {
	flag := topic.GetCheckFlag()
	expireIndex := topic.GetExpireIndex()
	OnlyCheck := true
	// 不论成功失败，都置位为发送
	r.subscribeMgr.SetTopicSendMs(topic, time.Now().UnixMilli())
	r.subscribeMgr.SetTopicCheckFlag(topic, subscribe.RetrySubscribeWithNoIndex)

	req, err := request.NewRequest(r.appId, topic.GetZoneId(), topic.GetTableName(), tcaplusCmd.TcaplusApiListGetAllReq, topic.IsPb())
	if err != nil {
		logger.ERR("topic %s NewRequest err %s", topic.GetKeyDebug(), err.Error())
		return
	}

	if flag == subscribe.RetrySubscribeWithNoIndex || expireIndex < 0 {
		expireIndex = -1
	} else if flag == subscribe.RetrySubscribeWithIndex {
		OnlyCheck = false
	}

	logger.DEBUG("retry SubscribeTopic %d", topic.GetAsyncId())
	// 设置为订阅模式
	if ret := req.SetSubscribe(topic.GetExpireSec(), expireIndex, OnlyCheck, uint32(topic.GetRspFlag())); ret != 0 {
		logger.ERR("topic %s req.SetSubscribe failed", topic.GetKeyDebug())
		return
	}

	req.SetAsyncId(uint64(topic.GetAsyncId()))
	rec, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("topic %s req.AddRecord failed", topic.GetKeyDebug())
		return
	}
	topic.Copy2KeyMap(rec.KeyMap)

	// 如果svr可用
	if svr := r.subscribeMgr.GetTopicSvr(topic); svr != nil {
		proxySvr := svr.(*server)
		if proxySvr.isAvailable() {
			data, err := req.Pack()
			if err != nil {
				logger.ERR("topic %s req pack failed %s", topic.GetKeyDebug(), err.Error())
				return
			}

			err = proxySvr.send(req.GetCmd(), data)
			if err != nil {
				logger.ERR("topic %s send failed %s", topic.GetKeyDebug(), err.Error())
				return
			}

			// success
			return
		}
	}

	// 使用新svr
	err, newSvr := r.SendRequest(req)
	if err != nil {
		logger.ERR("%s", err.Error())
		return
	}
	r.subscribeMgr.SetTopicSvr(topic, newSvr)
	return
}

func CheckProxyFunc(proxy interface{}) bool {
	if proxy == nil {
		return false
	}
	proxyServer, ok := proxy.(*server)
	if !ok {
		return false
	}
	if !proxyServer.isAvailable() {
		return false
	}
	return true
}

// SubscribeRoutine 定时续订阅
func (r *Router) startSubscribeRoutine() {
	r.ctrl.Add(1)
	go func() {
		defer r.ctrl.Done()
		// 50ms检查一次需要续订的idx
		reSubscribeTimer := time.NewTimer(r.subscribeMgr.GetReSubscribeIntervalTime())
		for {
			select {
			case <-reSubscribeTimer.C:
				reSubscribeTimer.Reset(r.subscribeMgr.GetReSubscribeIntervalTime())
				// 获取要重订阅的topic
				topicSlice := r.subscribeMgr.GetNeedRetryTopic(CheckProxyFunc)
				for _, topic := range topicSlice {
					if topic == nil {
						continue
					}
					r.retrySubscribeTopic(topic)
				}
			case <-r.chanClose:
				return
			}
		}
	}()
}
