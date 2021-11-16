package router

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcapdir_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/tnet"
	"sync"
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
	usingServerList   map[string]*server
	prepareServerList map[string]*server
	removeServerList  map[string]*server
}

//0 所有认证成功， 1 认证中， 2 部分认证成功， -1 有认证失败的
func (p *proxy) CheckAvailable() (int, error) {
	p.hashMutex.RLock()
	defer p.hashMutex.RUnlock()

	if len(p.hashList) == 0 {
		return 1, nil
	}

	signSucCount := 0
	for _, v := range p.hashList {
		if v.isAvailable() {
			signSucCount++
		} else if v.getSignUpStat() == SignUpFail {
			return -1, &terror.ErrorCode{Code: terror.ProxySignUpFailed}
		}
	}

	if signSucCount == 0 {
		return 1, nil
	}

	if signSucCount == len(p.usingServerList) {
		return 0, nil
	}
	return 2, nil
}

func (p *proxy) updateServerList() {
	for _, v := range p.usingServerList {
		v.update(false)
	}
	for _, v := range p.prepareServerList {
		v.update(false)
	}
	for _, v := range p.removeServerList {
		v.update(true)
	}
}

func (p *proxy) switchServerList() {
	//为空不用切换
	if len(p.prepareServerList) == 0 {
		return
	}

	//prepare中的server必须有鉴权通过的
	isAvailable := false
	for _, v := range p.prepareServerList {
		if v.isAvailable() {
			isAvailable = true
			break
		}
	}
	if !isAvailable {
		return
	}

	//和usingList相同，不用切换
	if len(p.prepareServerList) == len(p.usingServerList) {
		needSwitch := false
		for k, _ := range p.prepareServerList {
			if _, exist := p.usingServerList[k]; !exist {
				logger.INFO("prepare proxy %s not in using list", k)
				needSwitch = true
				break
			}
		}
		if !needSwitch {
			//清空prepare
			p.prepareServerList = make(map[string]*server)
			logger.INFO("proxy list not changed")
			return
		}
	}

	//切换
	logger.INFO("start switch proxy list!!!")
	//using中存在，prepare中不存在的需要挪到remove队列
	for k, v := range p.usingServerList {
		if _, exist := p.prepareServerList[k]; !exist {
			p.removeServerList[k] = v
			logger.INFO("move proxy %s to remove list", k)
		}
	}

	//using=prepare
	p.usingServerList = make(map[string]*server)
	for k, v := range p.prepareServerList {
		p.usingServerList[k] = v
	}

	//清空prepare
	p.prepareServerList = make(map[string]*server)

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

func (p *proxy) update() {
	p.updateServerList()
	p.switchServerList()

	//TODO remove 队列中1min没有回包的server进行删除操作
}

func (p *proxy) sendHeartbeat() {
	for _, v := range p.usingServerList {
		if v.isAvailable() {
			v.sendHeartbeat()
		}
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

	//唯一化,校验proxy地址
	accessUrlMap := make(map[string]bool)
	for i := 0; i < int(msg.AccessCount) && i < common.MaxProxyNumPerZone; i++ {
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
				appId: p.appId, zoneId: p.zoneId, signature: p.signature, proxyUrl: url, signUpFlag: NotSignUp,
				conn: nil, router: p.router, prepareStop: false,
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
	//prepare不为空，则先将AccessUrlList中不存在但prepare存在的移动到remove队列
	if len(p.prepareServerList) > 0 {
		for url, svr := range p.prepareServerList {
			if _, exist := accessUrlMap[url]; !exist {
				p.removeServerList[url] = svr
				delete(p.prepareServerList, url)
				logger.INFO("proxy %s move from prepare to remove list", url)
			}
		}
	}
	for url, _ := range accessUrlMap {
		//在prepare队列
		if _, exist := p.prepareServerList[url]; exist {
			logger.INFO("proxy %s in prepare list", url)
			continue
		}
		if server, exist := p.usingServerList[url]; exist {
			p.prepareServerList[url] = server
			logger.INFO("proxy %s in using list", url)
			continue
		}
		if server, exist := p.removeServerList[url]; exist {
			p.prepareServerList[url] = server
			delete(p.removeServerList, url)
			logger.INFO("proxy %s in remove list", url)
			continue
		}
		//新的节点
		svr := &server{
			appId: p.appId, zoneId: p.zoneId, signature: p.signature, proxyUrl: url, signUpFlag: NotSignUp,
			conn: nil, router: p.router, prepareStop: false,
		}
		svr.connect()
		p.prepareServerList[url] = svr
		logger.INFO("new proxy server %s to prepare", url)
	}
}

func (p *proxy) send(hashCode uint32, data []byte) error {
	p.hashMutex.RLock()
	defer p.hashMutex.RUnlock()

	if len(p.hashList) == 0 {
		return &terror.ErrorCode{Code: terror.ProxyNotAvailable}
	}
	preId := hashCode % uint32(len(p.hashList))
	id := preId

	for {
		svr := p.hashList[id]
		if svr.isAvailable() {
			return svr.send(data)
		}

		//选择下个节点
		hashCode++
		id = hashCode % uint32(len(p.hashList))
		//一轮之后
		if id == preId {
			return &terror.ErrorCode{Code: terror.ProxyNotAvailable}
		}
	}
}
