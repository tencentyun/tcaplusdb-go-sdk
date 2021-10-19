package tcaplus

import (
	"container/list"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/dir"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcapdir_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/router"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type netServer struct {
	//dir
	dirServer dir.DirServer
	//zone-->proxy
	router      router.Router
	initFlag    int
	initResult  chan error //client读取判断网络协程是否初始化成功
	stopNetWork chan bool  //client停止网络协程

	//res msg queue
	respCount    int64
	respMsgMutex sync.Mutex
	respMsgQueue *list.List

	//定时任务，时间间隔s
	dirListDuration   time.Duration
	proxyListDuration time.Duration
}

func (n *netServer) init(appId uint64, zoneList []uint32, dirUrl string, signature string, timeout uint32) error {
	n.initFlag = NotInit
	n.initResult = make(chan error, 1)
	n.stopNetWork = make(chan bool, 1)
	n.respMsgQueue = list.New()
	n.respCount = 0
	n.dirListDuration = 30
	n.proxyListDuration = 300

	//dir init
	if err := n.dirServer.Init(appId, zoneList, dirUrl, signature); err != nil {
		logger.ERR("dir init failed %s", err.Error())
		return err
	}

	//router init
	if err := n.router.Init(appId, zoneList, signature); err != nil {
		logger.ERR("router init failed %s", err.Error())
		return err
	}

	go n.netPkgProcess()
	return nil
}

//网络协程
func (n *netServer) netPkgProcess() {
	n.dirServer.Update()
	//30s 获取一次dir列表
	dirListTimer := time.NewTimer(n.dirListDuration * time.Second)
	//300s 获取一次proxy列表
	proxyListTimer := time.NewTimer(n.proxyListDuration * time.Second)
	//100ms 一次update
	updateTimer := time.NewTimer(100 * time.Millisecond)

	updateTraverse := time.NewTimer(time.Millisecond)
	updateHashListTimer := time.NewTimer(10 * time.Second)

	// 更新下当前时间，用于不需要精确时间的接口，防止一直去获取时间
	updateTimeNow := time.NewTimer(time.Second)

	for {
		select {
		case <-n.stopNetWork:
			logger.ERR("client net routine exit, close client")
			n.dirServer.DisConnect()
			n.router.Close()
			return
		//dir响应消息
		case dirPkg := <-n.dirServer.MsgPipe:
			n.processDirMsg(dirPkg)
		//proxy控制面响应消息
		case routerPkg := <-n.router.MsgPipe:
			n.processRouterMsg(routerPkg)
		//proxy列表更新定时器300s
		case <-proxyListTimer.C:
			if err := n.dirServer.GetAccessProxy(); err != nil {
				logger.ERR("GetAccessProxy failed, err %v", err.Error())
			} else {
				logger.INFO("GetAccessProxy send")
			}
			proxyListTimer.Reset(n.proxyListDuration * time.Second)
		//dir列表更新定时器30s
		case <-dirListTimer.C:
			if err := n.dirServer.GetDirList(); err != nil {
				logger.ERR("GetDirList failed, err %v", err.Error())
			} else {
				logger.INFO("GetDirList send")
			}
			dirListTimer.Reset(n.dirListDuration * time.Second)

		//update定时器 100ms - 1s
		case <-updateTimer.C:
			n.dirServer.Update()
			n.router.Update()
			if n.initFlag == NotInit {
				//检测所有proxy是否都认证通过
				if ret, err := n.router.CanStartUp(); ret == 0 {
					logger.INFO("router start finish")
					n.initFlag = InitSuccess
					n.initResult <- nil
				} else if ret == -1 {
					logger.INFO("router start failed")
					n.initFlag = InitFail
					n.initResult <- err
				}
				updateTimer.Reset(100 * time.Millisecond)
			} else {
				updateTimer.Reset(1 * time.Second)
			}
		case <-updateTraverse.C:
			n.router.TM.ContinueTraverse()
			updateTraverse.Reset(time.Millisecond)
		case <-updateTimeNow.C:
			common.TimeNow = time.Now()
			updateTimeNow.Reset(time.Second)
		case <-updateHashListTimer.C:
			n.router.UpdateHashList()
			updateHashListTimer.Reset(10 * time.Second)
		}
	}
}

func (n *netServer) processDirMsg(msg *tcapdir_protocol_cs.TCapdirCSPkg) {
	switch int64(msg.Head.Cmd) {
	case tcapdir_protocol_cs.TCAPDIR_CS_CMD_SIGNUP_RES:
		n.dirServer.ProcessSignUpRes(int(msg.Body.ResSignUpApp.Result))
		if n.initFlag == NotInit {
			//第一次启动,拉取proxy列表
			if msg.Body.ResSignUpApp.Result == 0 {
				n.dirServer.GetAccessProxy()
			} else {
				//dir认证失败
				n.initFlag = InitFail
				logger.ERR("dir signUp failed %d", msg.Body.ResSignUpApp.Result)
				n.initResult <- &terror.ErrorCode{Code: int(msg.Body.ResSignUpApp.Result)}
			}
		}
		return
	case tcapdir_protocol_cs.TCAPDIR_CS_CMD_HEARTBEAT_RES:
		logger.DEBUG("recv dir heartbeat res")
		return
	case tcapdir_protocol_cs.TCAPDIR_CS_CMD_GET_DIR_SERVER_LIST_RES:
		logger.INFO("GET_DIR_SERVER_LIST_RES DirServerCount:%d DirServer:%v",
			msg.Body.ResGetDirServerList.DirServerCount,
			msg.Body.ResGetDirServerList.DirServer[0:msg.Body.ResGetDirServerList.DirServerCount])
		n.dirServer.ProcessDirListRes(msg.Body.ResGetDirServerList)
		return

	case tcapdir_protocol_cs.TCAPDIR_CS_CMD_GET_TABLES_AND_ACCESS_RES:
		res := msg.Body.ResGetTablesAndAccess
		logger.INFO("GET_TABLES_AND_ACCESS_RES SetID:%d AppID:%d ZoneID:%d "+
			"TableCount:%d TableNameList:%v AccessCount:%d AccessUrlList:%v AccessIdList:%v"+
			"DirAvailableCheckPeriod:%d DirUpdateListInterval:%d DirUpdateTablesAndAcessInterval:%d"+
			" ApiFromProxyHeartBeatTime:%d ApiFromDirHeartBeatTime:%d",
			res.SetID, res.AppID, res.ZoneID,
			res.TableCount, res.TableNameList[0:res.TableCount],
			res.AccessCount, res.AccessUrlList[0:res.AccessCount],
			res.AccessIdList[0:res.AccessCount],
			res.ConfData.DirAvailableCheckPeriod, res.ConfData.DirUpdateListInterval,
			res.ConfData.DirUpdateTablesAndAcessInterval,
			res.ConfData.ApiFromProxyHeartBeatTime, res.ConfData.ApiFromDirHeartBeatTime)

		//更新周期
		n.dirServer.SetHeartbeatInterval(time.Duration(res.ConfData.ApiFromDirHeartBeatTime))
		if res.ConfData.DirAvailableCheckPeriod <= res.ConfData.DirUpdateListInterval &&
			res.ConfData.DirUpdateListInterval > 0 && res.ConfData.DirUpdateListInterval <= 300 {
			n.dirListDuration = time.Duration(res.ConfData.DirUpdateListInterval)
		}

		if res.ConfData.DirAvailableCheckPeriod <= res.ConfData.DirUpdateTablesAndAcessInterval &&
			res.ConfData.DirUpdateTablesAndAcessInterval > 0 && res.ConfData.DirUpdateTablesAndAcessInterval <= 300 {
			n.proxyListDuration = time.Duration(res.ConfData.DirUpdateTablesAndAcessInterval)
		}

		if res.ConfData.ApiFromProxyHeartBeatTime > 0 && res.ConfData.ApiFromProxyHeartBeatTime <= 5 {
			n.router.SetHeartbeatInterval(time.Duration(res.ConfData.ApiFromProxyHeartBeatTime))
		}

		n.router.ProcessTablesAndAccessMsg(res)
		return
	}
}

//此处的是用户pkg，非用户pkg已在func (s *server)processRsp回调中处理
func (n *netServer) processRouterMsg(msg *tcaplus_protocol_cs.TCaplusPkg) {
	n.respMsgMutex.Lock()
	defer n.respMsgMutex.Unlock()
	n.respMsgQueue.PushBack(msg)
	atomic.AddInt64(&n.respCount, 1)
	logger.DEBUG("add one msg to queue, %d", n.respCount)
}

func (n *netServer) recvResponse() (response.TcaplusResponse, error) {
	return n.router.RecvResponse()
}

func (n *netServer) sendRequest(req request.TcaplusRequest) error {
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

	return n.router.Send(code, req.GetZoneId(), data)
}
