package router

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcapdir_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/statistics"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/tnet"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type proxy struct {
	appId     uint64
	zoneId    uint32
	signature string
	router    *Router

	//用户协程和网络协程会同时操作
	tbMutex       sync.RWMutex
	tableNameList map[string]bool

	//锁hashList，用户协程，网络协程会同时操作
	hashMutex sync.RWMutex
	hashList  []*server

	//只有网络协程操作
	prepareServerList map[string]bool
	connectStep       int32 // 连接步长
	deleteStep        int32 // 删除步长

	usingServerList  map[string]*server
	removeServerList map[string]*server

	statMgr         statistics.StatMgr
	perfPercent     int32
	lastIsolateTime time.Time

	// 300ms隔离proxy的比例
	isolateProxyRate int32
}

func (p *proxy) GetErrorStr() string {
	var errStr string
	for _, v := range p.usingServerList {
		if v.error != nil {
			errStr = errStr + v.error.Error() + ","
		}
	}
	return errStr
}

//0 所有认证成功， 1 认证中， 2 部分认证成功， -1 全部认证失败
func (p *proxy) CheckAvailable() (int, error) {
	p.hashMutex.RLock()
	defer p.hashMutex.RUnlock()

	if len(p.hashList) == 0 {
		return 1, nil
	}

	signSucCount := 0
	signFailCount := 0
	var err error
	for _, v := range p.hashList {
		if v.isAvailable() {
			signSucCount++
		} else if v.getSignUpStat() == SignUpFail {
			signFailCount++
			err = v.error
		}
	}

	//全部认证失败
	if signFailCount == len(p.hashList) {
		if err != nil {
			return -1, err
		}
		return -1, &terror.ErrorCode{Code: terror.ProxySignUpFailed}
	}

	//全部认证中
	if signSucCount == 0 {
		return 1, nil
	}

	//全部认证成功
	if signSucCount == len(p.hashList) {
		return 0, nil
	}
	//部分认证成功
	return 2, nil
}

func (p *proxy) updateServerList() {
	for _, v := range p.usingServerList {
		v.update(false)
	}
	for _, v := range p.removeServerList {
		v.update(true)
	}
}

func (p *proxy) DeletePartProxyServer() {
	delTotal := int(0)
	availableTotal := int(0)
	// 不可用且带删除标记的全部挪到remove队列
	for url, server := range p.usingServerList {
		if server.hasDeleteFlag && !server.isAvailable() {
			p.removeServerList[url] = server
			delete(p.usingServerList, url)
			logger.INFO("move not available proxy %s to remove list", url)
			continue
		}
		if server.hasDeleteFlag {
			delTotal++
		}
		if server.isAvailable() {
			availableTotal++
		}
	}

	// 不用删除
	if delTotal <= 0 {
		return
	}

	// 计算可删除的数量
	needDeleteCnt := int(p.deleteStep)
	if needDeleteCnt < 1 {
		needDeleteCnt = 1
	}
	if availableTotal >= delTotal+10 {
		// 有足够多的可用proxy
		needDeleteCnt = delTotal
	}

	// 挪动10%的delete到remove队列
	for url, server := range p.usingServerList {
		if availableTotal <= 1 {
			return
		}
		if needDeleteCnt <= 0 {
			return
		}

		if server.hasDeleteFlag {
			p.removeServerList[url] = server
			delete(p.usingServerList, url)
			needDeleteCnt--
			availableTotal--
			logger.INFO("move available proxy %s to remove list, needDeleteCnt %d", url, needDeleteCnt)
		}
	}
}

func (p *proxy) switchServerList() {
	//remove队列不清空不用切换
	if len(p.removeServerList) != 0 {
		return
	}

	// 挪动10%的delete到remove队列
	p.DeletePartProxyServer()
	hasDelete := len(p.removeServerList) > 0

	// 添加10%的prepare到using队列
	hasAdd := false
	needConnectCnt := int32(len(p.removeServerList))
	if needConnectCnt < p.connectStep {
		needConnectCnt = p.connectStep
	}
	if needConnectCnt < 1 {
		needConnectCnt = 1
	}
	for url, _ := range p.prepareServerList {
		if needConnectCnt <= 0 {
			break
		}
		needConnectCnt--
		hasAdd = true
		// new proxy
		svr := &server{
			appId:     p.appId,
			zoneId:    p.zoneId,
			signature: p.signature,
			proxyUrl:  url, signUpFlag: NotSignUp,
			conn:         nil,
			router:       p.router,
			prepareStop:  0,
			isolated:     false,
			isolatedTime: time.Now(),
			connStat:     p.statMgr.NewConnStat(url),
		}
		svr.connect()
		p.usingServerList[url] = svr
		delete(p.prepareServerList, url)
		logger.INFO("add proxy %s to using, needConnectCnt %d", url, needConnectCnt)
	}

	// 路由无变化
	if !hasDelete && !hasAdd {
		return
	}

	//设置选路hash表
	p.hashMutex.Lock()
	p.hashList = make([]*server, 0, len(p.usingServerList))
	for _, v := range p.usingServerList {
		if v.isAvailable() {
			p.hashList = append(p.hashList, v)
		}
	}
	p.hashMutex.Unlock()
	logger.INFO("hashList %v", p.usingServerList)
}

func (p *proxy) updateHashList() {
	//设置选路hash表
	p.hashMutex.Lock()
	p.hashList = make([]*server, 0, len(p.usingServerList))
	for _, v := range p.usingServerList {
		if v.isAvailable() {
			p.hashList = append(p.hashList, v)
		}
	}
	p.hashMutex.Unlock()
}

func (p *proxy) update(curTime time.Time) {
	p.updateServerList()
	p.switchServerList()
	//remove 队列中一个超时时间没有回包的server进行删除操作
	p.clearRemoveServerList()
	p.IsolateProblemSvr(curTime)
}

func (p *proxy) SetIsolateProxyRate(rate int32) {
	if rate == 0 || rate > 100 {
		// 非法值使用 默认值比例10%
		p.isolateProxyRate = 10
	} else {
		// rate 可能小于0表示关闭
		p.isolateProxyRate = rate
	}
	logger.INFO("config rate %d, real SetIsolateProxyRate %d", rate, p.isolateProxyRate)
}

func (p *proxy) GetCanIsolateProxyCnt() int32 {
	if p.isolateProxyRate <= 0 {
		return 0
	}
	// 总数和可用数
	total := int32(len(p.usingServerList))
	availableNum := int32(0)
	for _, v := range p.usingServerList {
		if v.isAvailable() {
			availableNum++
		}
	}
	if availableNum <= 2 {
		return 0
	}

	// 按比例计算，可隔离上限，至少1个
	thresholdCnt := total * p.isolateProxyRate / 100
	if thresholdCnt < 1 {
		thresholdCnt = 1
	}

	// 至少保留2个
	if thresholdCnt > availableNum-2 {
		thresholdCnt = availableNum - 2
	}

	// 已经隔离的数量超过上限
	unavailableNum := total - availableNum
	if unavailableNum >= thresholdCnt {
		return 0
	}
	// 剩余可隔离的数量
	return thresholdCnt - unavailableNum
}

func (p *proxy) IsolateProblemSvr(curTime time.Time) {
	diff := curTime.Sub(p.lastIsolateTime)
	if diff < 300*time.Millisecond {
		return
	}
	p.lastIsolateTime = curTime

	canIsolatedNum := p.GetCanIsolateProxyCnt()
	for _, v := range p.usingServerList {
		v.checkIsolateStat(curTime, &canIsolatedNum)
	}
}

func (p *proxy) clearRemoveServerList() {
	curMs := time.Now().UnixMilli()
	timeoutMs := p.router.ctrl.Option.ProxyConnOption.ConTimeout.Milliseconds()
	if timeoutMs < 1000 {
		timeoutMs = 1000
	}
	for k, v := range p.removeServerList {
		lastRspTimeMs := atomic.LoadInt64(&v.lastRspTime)
		if !v.isConnected() || curMs > lastRspTimeMs+timeoutMs {
			logger.INFO("remove list remove svr %s lastRspTime %v isConnected %v",
				v.proxyUrl, lastRspTimeMs, v.isConnected())
			v.disConnect()
			delete(p.removeServerList, k)
		}
	}
}

func (p *proxy) sendHeartbeat() {
	for _, v := range p.usingServerList {
		v.sendHeartbeat()
	}
}

func (p *proxy) statReset() {
	p.statMgr.Reset()
	for _, v := range p.usingServerList {
		v.connStat.Reset()
	}
}

func (p *proxy) shuffleProxyList(msg *tcapdir_protocol_cs.ResGetTablesAndAccess) {
	// 洗牌
	rand.Seed(time.Now().UnixNano())
	n := int(msg.AccessCount)
	if n > len(msg.AccessUrlList) {
		n = len(msg.AccessUrlList)
	}
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1) // 从 [0, i] 中随机取一个
		msg.AccessUrlList[i], msg.AccessUrlList[j] = msg.AccessUrlList[j], msg.AccessUrlList[i]
	}
}

func (p *proxy) processTablesAndAccessMsg(msg *tcapdir_protocol_cs.ResGetTablesAndAccess) {
	//设置本区下表名称
	p.tbMutex.Lock()
	p.tableNameList = make(map[string]bool)
	for i := 0; i < int(msg.TableCount); i++ {
		p.tableNameList[msg.TableNameList[i]] = true
	}
	p.tbMutex.Unlock()
	if msg.AccessCount <= 0 {
		return
	}
	p.SetIsolateProxyRate(msg.ConfData.MetalibID)

	p.shuffleProxyList(msg)

	//唯一化,校验proxy地址
	maxProxyNumPerZone := p.router.ctrl.Option.ProxyConnOption.ProxyMaxCount
	if maxProxyNumPerZone < 2 {
		maxProxyNumPerZone = 2
	}
	accessUrlMap := make(map[string]bool)
	for i := 0; i < int(msg.AccessCount) && i < maxProxyNumPerZone; i++ {
		url := msg.AccessUrlList[i]
		urlNet, _, urlPort, err := tnet.ParseUrl(&url)
		if err != nil {
			logger.ERR("proxy url is invalid %s", url)
		}
		// 变更IP
		if common.PublicIP != "" {
			url = fmt.Sprintf("%s://%s:%s", urlNet, common.PublicIP, urlPort)
		}
		accessUrlMap[url] = true
	}
	//初始化
	if len(p.usingServerList) == 0 {
		for url, _ := range accessUrlMap {
			svr := &server{
				appId:        p.appId,
				zoneId:       p.zoneId,
				signature:    p.signature,
				proxyUrl:     url,
				signUpFlag:   NotSignUp,
				conn:         nil,
				router:       p.router,
				prepareStop:  0,
				isolated:     false,
				isolatedTime: time.Now(),
				connStat:     p.statMgr.NewConnStat(url),
			}
			svr.connect()
			p.usingServerList[url] = svr
			logger.INFO("new proxy server %s", url)
		}
		//设置选路hash表
		p.hashMutex.Lock()
		p.hashList = make([]*server, 0, len(p.usingServerList))
		for _, v := range p.usingServerList {
			p.hashList = append(p.hashList, v)
		}
		p.hashMutex.Unlock()
		logger.INFO("hashList %v", p.usingServerList)
		return
	}

	// 清空prepare
	p.prepareServerList = make(map[string]bool)

	// using中的所有打上删除标记
	for _, server := range p.usingServerList {
		server.hasDeleteFlag = true
	}

	deleteTotal := len(p.usingServerList)
	// accessUrlMap中，在using中存在的剔除删除标记，否则加入prepare等待链接
	for url, _ := range accessUrlMap {
		//在using队列
		if server, exist := p.usingServerList[url]; exist {
			server.hasDeleteFlag = false
			deleteTotal--
			logger.INFO("proxy %s in using list", url)
			continue
		}
		p.prepareServerList[url] = true
		logger.INFO("add proxy server %s to prepare", url)
	}
	logger.INFO("delete cnt %d add cnt %d", deleteTotal, len(p.prepareServerList))
	percent := p.router.ctrl.Option.SwitchProxyPercent
	if percent <= 0 {
		percent = 10
	}

	p.deleteStep = int32(deleteTotal * percent / 100)
	if p.deleteStep < 1 {
		p.deleteStep = 1
	}

	p.connectStep = int32(len(p.prepareServerList) * percent / 100)
	if p.connectStep < 1 {
		p.connectStep = 1
	}
}

func (p *proxy) send(cmd uint32, hashCode uint32, data []byte) (error, *server) {
	p.hashMutex.RLock()
	defer p.hashMutex.RUnlock()

	if len(p.hashList) == 0 {
		p.statMgr.IncRouteFail()
		return &terror.ErrorCode{Code: terror.ProxyNotAvailable}, nil
	}
	id := hashCode % uint32(len(p.hashList))
	svr := p.hashList[id]
	if svr.isAvailable() {
		return svr.send(cmd, data), svr
	}

	// 不可用则随机一个
	rand.Seed(time.Now().UnixNano())
	preId := uint32(rand.Int() % len(p.hashList))
	id = preId
	for {
		svr := p.hashList[id]
		if svr.isAvailable() {
			return svr.send(cmd, data), svr
		}
		//选择下个节点
		id++
		id = id % uint32(len(p.hashList))
		//一轮之后
		if id == preId {
			p.statMgr.IncRouteFail()
			return &terror.ErrorCode{Code: terror.ProxyNotAvailable}, nil
		}
	}
}
