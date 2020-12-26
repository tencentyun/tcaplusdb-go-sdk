package router

import (
	"container/list"
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

type SyncRequest struct {
	//sync response msg chan
	syncMsgPipe chan *tcaplus_protocol_cs.TCaplusPkg

	//request package
	requestPkg request.TcaplusRequest
}

func (S *SyncRequest) InitMoreChan(reqpkg request.TcaplusRequest, num int32) {
	S.syncMsgPipe = make(chan *tcaplus_protocol_cs.TCaplusPkg, num)
	S.requestPkg = reqpkg
}

func (S *SyncRequest) Init(reqpkg request.TcaplusRequest) {
	S.syncMsgPipe = make(chan *tcaplus_protocol_cs.TCaplusPkg, 1)
	S.requestPkg = reqpkg
}

func (S *SyncRequest) GetSyncChan() chan *tcaplus_protocol_cs.TCaplusPkg {
	return S.syncMsgPipe
}

func (S *SyncRequest) SyncChanClose() {
	close(S.syncMsgPipe)
}

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
	reqChanMutex   sync.RWMutex
	requestChanMap map[int32]SyncRequest
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
	r.reqChanMutex.RLock()
	if v, exist := r.requestChanMap[msg.Head.Seq]; exist {
		v.syncMsgPipe <- msg
	} else {
		logger.ERR("Can not find request chan %d", msg.Head.Seq)
	}
	r.reqChanMutex.RUnlock()
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

func (r *Router) Init(appId uint64, zoneList []uint32, signature string) error {
	r.appId = appId
	r.zoneList = make([]uint32, len(zoneList))
	copy(r.zoneList, zoneList)
	r.signature = signature
	r.MsgPipe = make(chan *tcaplus_protocol_cs.TCaplusPkg, 1024)
	r.proxyMap = make(map[uint32]*proxy)
	r.heartbeatInterval = 1
	r.lastHeartbeatTime = time.Now()
	r.respMsgQueue = list.New()
	r.requestChanMap = make(map[int32]SyncRequest)
	return nil
}

func (r *Router) RequestChanMapAdd(reqseq int32, syncrequest SyncRequest) {
	r.reqChanMutex.Lock()
	defer r.reqChanMutex.Unlock()
	r.requestChanMap[reqseq] = syncrequest
}

func (r *Router) RequestChanMapClean(reqsq int32) {
	r.reqChanMutex.Lock()
	defer r.reqChanMutex.Unlock()
	delete(r.requestChanMap, reqsq)
}

func (r *Router) SetHeartbeatInterval(heartbeatInterval time.Duration) {
	r.heartbeatInterval = heartbeatInterval
}

func (r *Router) CheckTable(zoneId uint32, tableName string) error {
	if proxy, exist := r.proxyMap[zoneId]; !exist {
		return &terror.ErrorCode{Code: terror.ZoneIdNotExist}
	} else {
		proxy.tbMutex.RLock()
		defer proxy.tbMutex.RUnlock()
		if _, exist := proxy.tableNameList[tableName]; !exist {
			return &terror.ErrorCode{Code: terror.TableNotExist}
		}
	}

	return nil
}

//0 所有认证成功， 1 有proxy全部认证中， 2 所有proxy部分认证成功，可以启动， -1 有认证失败,启动失败
func (r *Router) CanStartUp() (int, error) {
	if 0 == len(r.proxyMap) {
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

func (r *Router) Send(hashCode uint32, zoneId uint32, data []byte) error {
	return r.proxyMap[zoneId].send(hashCode, data)
}
