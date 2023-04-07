package router

import (
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/config"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/tnet"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/traverser"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcapdir_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

// 同步请求结构体
type SyncRequest struct {
	//sync response msg chan
	syncMsgPipe chan *tcaplus_protocol_cs.TCaplusPkg
	//request package
	requestPkg request.TcaplusRequest
	delFromMap byte //0 添加，1 删除，2 收到响应后自动删除
	syncId     uint32
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
		go func(id int) {
			r.ctrl.Add(1)
			defer r.ctrl.Done()
			ch := r.requestChanList[id]
			rmap := r.requestChanMap[id]
			proc := func(p interface{}) {
				if req, ok := p.(*SyncRequest); ok {
					if req.delFromMap == 1 {
						delete(rmap, req.syncId)
					} else {
						if req.requestPkg != nil {
							err := r.sendRequest(req.requestPkg)
							if err != nil {
								logger.ERR("Send failed %s", err.Error())
								return
							}
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
		go func(ch chan *tnet.PKG) {
			r.ctrl.Add(1)
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
	r.respMsgQueue = list.New()
	r.ctrl = ctrl

	//放最后
	r.createRequestRoutine()
	r.createResponseRoutine()
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
	syncRequest.delFromMap = 1
	select {
	case <-r.chanClose:
		logger.INFO("processSyncOperate exit")
		return -1
	case r.requestChanList[syncId] <- syncRequest:
	}
	return 0
}

func (r *Router) SetHeartbeatInterval(heartbeatInterval time.Duration) {
	r.heartbeatInterval = heartbeatInterval
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
	for _, proxy := range r.proxyMap {
		proxy.update()
	}

	//发送心跳
	diff := time.Now().Sub(r.lastHeartbeatTime)
	if diff > r.heartbeatInterval*time.Second {
		for _, proxy := range r.proxyMap {
			proxy.sendHeartbeat()
		}
		r.lastHeartbeatTime = time.Now()
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
			prepareServerList: make(map[string]*server),
			removeServerList:  make(map[string]*server),
		}
		p.processTablesAndAccessMsg(msg)
		r.proxyMap[uint32(msg.ZoneID)] = p
	}
}

func (r *Router) sendRequest(req request.TcaplusRequest) error {
	//打包
	data, err := req.Pack()
	if err != nil {
		logger.ERR("req pack failed %s", err.Error())
		return err
	}
	//获取keyHash
	code, err := req.GetKeyHash()
	if err != nil {
		logger.ERR("get key hash failed %s", err.Error())
		return err
	}
	err = r.Send(code, req.GetZoneId(), data)
	if err != nil {
		logger.ERR("Send failed %s", err.Error())
		return err
	}
	return nil
}

func (r *Router) Send(hashCode uint32, zoneId uint32, data []byte) error {
	if p, exist := r.proxyMap[zoneId]; exist {
		return p.send(hashCode, data)
	}
	logger.ERR("zone %d not connect", zoneId)
	return &terror.ErrorCode{Code: terror.SendRequestFail, Message: "zone proxy not connect"}
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
		for _, svr := range v.prepareServerList {
			if svr == nil {
				continue
			}
			svr.disConnect()
		}
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
