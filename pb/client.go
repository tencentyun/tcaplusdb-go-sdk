package tcaplus

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/config"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/version"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/router"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/traverser"
	"hash/crc32"
	"sort"
	"sync/atomic"
	"time"
)

const (
	NotInit     = 0
	InitSuccess = 1
	InitFail    = 2
)

/**
	@brief tcaplus api客户端
	@param [IN] appId 业务id
	@param [IN] zoneList 区列表
	@param [IN] dirUrl	dir地址
	@param [IN] initFlag 是否初始化
	@param [IN] netServer 服务管理
**/
type client struct {
	appId      uint64
	zoneList   []uint32
	dirUrl     string
	initFlag   int
	netServer  netServer
	tm         *traverser.TraverserManager
	isPB       bool
	reqSeq     uint32
	ctrl       *config.ClientCtrl
	defZone    int32
	defTimeout time.Duration
}

/**
   @brief 创建一个tcaplus api客户端
   @retval 返回客户端指针
**/
func newClient(isPB bool) *client {
	c := new(client)
	c.tm = traverser.NewTraverserManager(c)
	c.netServer.router.TM = c.tm
	c.initFlag = NotInit
	c.isPB = isPB
	c.reqSeq = 1
	c.ctrl = &config.ClientCtrl{
		Option: config.NewDefaultClientOption(),
	}
	c.defZone = -1
	return c
}

/**
   @brief                   设置客户端可选参数，请在client.Dial之前调用
   @param [IN] opt      	客户端可选参数
**/
func (c *client) SetOpt(opt *config.ClientOption) {
	*c.ctrl.Option = *opt
}

/**
   @brief                   设置API日志配置文件全路径log.conf(json格式，example下有示例)，请在client.Dial之前调用
   @param [IN] cfgPath      日志配置文件全路径log.conf
   @retval 					错误码
   @note                    Api日志默认使用的zap，用户也可自行实现日志接口logger.LogInterface，调用SetLogger进行设置
**/
func (c *client) SetLogCfg(cfgPath string) error {
	return logger.SetLogCfg(cfgPath)
}

/**
   @brief                   自定义API日志接口,用户实现logger.LogInterface日志接口，日志将打印到用户的日志接口中，请在client.Dial之前调用
   @param [IN] handle       logger.LogInterface类型的日志接口
   @retval                  错误码
**/
func (c *client) SetLogger(handle logger.LogInterface) {
	logger.SetLogger(handle)
}

/**
   @brief 连接tcaplue函数
   @param [IN] appId         appId，在网站注册相应服务以后，你可以得到该appId
   @param [IN] zoneList      需要操作表的区服ID列表，操作的表在多个不同的zone，填zoneId列表；操作的表在一个zone，zone列表填一个zoneId
   @param [IN] signature     签名/密码，在网站注册相应服务以后，你可以得到该字符串
   @param [IN] dirUrl        目录服务器的url，形如"tcp://172.25.40.181:10600"
   @param [IN] timeout       second, 连接所有表对应的tcaplus proxy服务器。若所有的proxy连通且鉴权通过，则立即返回成功；
							     若到达超时时间，只要有一个proxy连通且鉴权通过，也会返回成功；否则返回超时错误。
   @retval                   错误码
**/
func (c *client) Dial(appId uint64, zoneList []uint32, dirUrl string, signature string, timeout uint32) error {
	c.appId = appId
	c.zoneList = make([]uint32, len(zoneList))
	copy(c.zoneList, zoneList)
	c.dirUrl = dirUrl
	if len(c.zoneList) == 1 {
		c.defZone = int32(c.zoneList[0])
	}
	//log init
	logger.Init()
	if err := c.netServer.init(appId, zoneList, dirUrl, signature, timeout, c.ctrl); err != nil {
		logger.ERR("net start failed %s", err.Error())
		c.initFlag = InitFail
		return err
	}
	logger.INFO("Tcaplus Go Api Version: %s", version.Version)

	//wait init success
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if ret, err := c.netServer.router.CanStartUp(); ret == 0 {
			c.initFlag = InitSuccess
			logger.INFO("init success.")
			return nil
		} else if ret == 2 {
			c.initFlag = InitSuccess
			logger.INFO("init part success.")
			return nil
		} else if ret == -1 {
			logger.ERR("init failed %v", err.Error())
			c.initFlag = InitFail
			c.netServer.stopNetWork <- true
			return err
		}
		logger.ERR("init timeout %v", timeout)
		c.initFlag = InitFail
		c.netServer.stopNetWork <- true
		return &terror.ErrorCode{Code: terror.ClientInitTimeOut, Message: "init timeout"}
	case ret := <-c.netServer.initResult:
		if ret != nil {
			logger.ERR("init failed. %s", ret.Error())
			c.initFlag = InitFail
			c.netServer.stopNetWork <- true
			return ret
		} else {
			c.initFlag = InitSuccess
			logger.INFO("init success.")
			return nil
		}
	}
}

/**
    @brief 设置默认zoneId
	@param [IN] zoneId zoneID
    @retval error 错误码，如果未dial调用此接口将会返错 ClientNotDial
**/
func (c *client) SetDefaultZoneId(zoneId uint32) error {
	if c.initFlag != InitSuccess {
		return &terror.ErrorCode{Code: terror.ClientNotInit}
	}
	c.defZone = int32(zoneId)
	return nil
}

/**
    @brief 设置请求默认超时时间
**/
func (c *client) SetDefaultReqTimeout(timeout time.Duration) {
	c.defTimeout = timeout
}

/**
@brief 创建指定分区指定表的请求
@param [IN] zoneId              分区ID
@param [IN] tableName           表名
@param [IN] cmd                 命令字(cmd.TcaplusApiGetReq等)
@retval request.TcaplusRequest  tcaplus请求
@retval error                   错误码
*/
func (c *client) NewRequest(zoneId uint32, tableName string, cmd int) (request.TcaplusRequest, error) {
	if c.initFlag != InitSuccess {
		return nil, &terror.ErrorCode{Code: terror.ClientNotInit}
	}

	if len(tableName) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_TABLE_NAME_LEN) {
		return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "table name too long"}
	}

	if err := c.netServer.router.CheckTable(zoneId, tableName); err != nil {
		return nil, err
	}

	for _, z := range c.zoneList {
		if z == zoneId {
			if req, err := request.NewRequest(c.appId, zoneId, tableName, cmd, c.isPB); err != nil {
				logger.ERR("new request failed, %s", err.Error())
				return nil, err
			} else {
				return req, nil
			}
		}
	}
	return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "zoneId is invalid"}
}

/**
@brief 发送tcaplus请求
@param [IN] req       tcaplus请求
@retval error         错误码
*/
func (c *client) SendRequest(req request.TcaplusRequest) error {
	if c.initFlag != InitSuccess {
		return &terror.ErrorCode{Code: terror.ClientNotInit}
	}
	err := c.netServer.sendRequest(req)
	logger.DEBUG("SendRequest Finish")
	return err
}

/**
	@brief 异步接收tcaplus响应
	@retval response.TcaplusResponse tcaplus响应
	@retval error 错误码
	@note   error nil，response nil 成功但当前无响应消息
			error nil, response 非nil，成功获取响应消息
            error 非nil，接收响应出错
*/
func (c *client) RecvResponse() (response.TcaplusResponse, error) {
	if c.initFlag != InitSuccess {
		return nil, &terror.ErrorCode{Code: terror.ClientNotInit}
	}
	return c.netServer.recvResponse()
}

/**
    @brief 发送tcaplus同步请求并接受响应
	@param [IN] req tcaplus请求
	@param [IN] timeout 超时时间
    @retval response.TcaplusResponse tcaplus响应
    @retval error 错误码
            error nil，response nil 成功但当前无响应消息
            error nil, response 非nil，成功获取响应消息
            error 非nil，接收响应出错
**/
func (c *client) Do(req request.TcaplusRequest, timeout time.Duration) (response.TcaplusResponse, error) {
	if c.initFlag != InitSuccess {
		return nil, &terror.ErrorCode{Code: terror.ClientNotInit}
	}

	requestSeq := int32(atomic.AddUint32(&c.reqSeq, 1))
	if requestSeq == 0 {
		requestSeq = int32(atomic.AddUint32(&c.reqSeq, 1))
	}
	req.SetSeq(requestSeq)
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	var synrequestPkg router.SyncRequest
	synrequestPkg.Init(req)
	if c.netServer.router.RequestChanMapAdd(&synrequestPkg) == -1 {
		return nil, &terror.ErrorCode{Code: terror.RouterIsClosed}
	}

	for {
		select {
		case <-timer.C:
			logger.ERR("requestSeq %d :%s, timeout", requestSeq, timeout.String())
			c.netServer.router.RequestChanMapClean(&synrequestPkg)
			return nil, &terror.ErrorCode{Code: terror.TimeOut, Message: timeout.String() + ", timeout"}
		case routerPkg := <-synrequestPkg.GetSyncChan():
			return response.NewResponse(routerPkg)
		}
	}
}

/**
    @brief 发送tcaplus同步请求并接受响应
	@param [IN] req tcaplus请求
	@param [IN] timeout 超时时间
    @retval []response.TcaplusResponse tcaplus响应
    @retval error 错误码
            error nil，response nil 成功但当前无响应消息
            error nil, response 非nil，成功获取响应消息
            error 非nil，response 非nil 接收部分回包正确，但是收到了错误包或者超时退出
**/
func (c *client) DoMore(req request.TcaplusRequest, timeout time.Duration) ([]response.TcaplusResponse, error) {
	requestSeq := int32(atomic.AddUint32(&c.reqSeq, 1))
	if requestSeq == 0 {
		requestSeq = int32(atomic.AddUint32(&c.reqSeq, 1))
	}
	req.SetSeq(requestSeq)
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	var synrequestPkg router.SyncRequest
	synrequestPkg.InitMoreChan(req, 1024)

	if c.netServer.router.RequestChanMapAdd(&synrequestPkg) == -1 {
		return nil, &terror.ErrorCode{Code: terror.RouterIsClosed}
	}
	defer c.netServer.router.RequestChanMapClean(&synrequestPkg)

	var resp_list []response.TcaplusResponse
	var idx int = 0
	for {
		select {
		case <-timer.C:
			logger.ERR("requestSeq %d :%s, timeout, current pkg num %d", requestSeq, timeout.String(), idx)
			return resp_list, &terror.ErrorCode{Code: terror.TimeOut, Message: timeout.String() + ", timeout"}
		case routerPkg := <-synrequestPkg.GetSyncChan():
			resp, err := response.NewResponse(routerPkg)
			idx += 1
			if err == nil {
				resp_list = append(resp_list, resp)
				if 1 == resp.HaveMoreResPkgs() {
					continue
				} else {
					return resp_list, nil
				}
			} else {
				logger.ERR("requestSeq %d, current pkg num: %d,  %s", requestSeq, idx, err.Error())
				return resp_list, err
			}
		}
	}
}

/**
    @brief 发送tcaplus同步请求并接受响应
	@param [IN] tra 遍历器
	@param [IN] timeout 超时时间
    @retval []response.TcaplusResponse tcaplus响应
    @retval error 错误码
            error nil，response nil 成功但当前无响应消息
            error nil, response 非nil，成功获取响应消息
            error 非nil，response 非nil 接收部分回包正确，但是收到了错误包或者超时退出
**/
func (c *client) DoTraverse(tra *traverser.Traverser, timeout time.Duration) ([]response.TcaplusResponse, error) {
	requestSeq := int32(atomic.AddUint32(&c.reqSeq, 1))
	if requestSeq == 0 {
		requestSeq = int32(atomic.AddUint32(&c.reqSeq, 1))
	}
	err := tra.SetSeq(requestSeq)
	if err != nil {
		return nil, err
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	var synrequestPkg router.SyncRequest
	synrequestPkg.InitTraverseChan(requestSeq, 1024)

	if c.netServer.router.RequestChanMapAdd(&synrequestPkg) == -1 {
		return nil, &terror.ErrorCode{Code: terror.RouterIsClosed}
	}
	defer c.netServer.router.RequestChanMapClean(&synrequestPkg)

	err = tra.Start()
	if err != nil {
		return nil, err
	}

	var resp_list []response.TcaplusResponse
	var idx int = 0
	for {
		select {
		case <-timer.C:
			logger.ERR("requestSeq %d :%s, timeout, current pkg num %d", requestSeq, timeout.String(), idx)
			return resp_list, &terror.ErrorCode{Code: terror.TimeOut, Message: timeout.String() + ", timeout"}
		case routerPkg := <-synrequestPkg.GetSyncChan():
			resp, err := response.NewResponse(routerPkg)
			idx += 1
			if err == nil {
				resp_list = append(resp_list, resp)
				if traverser.TraverseStateNormal == tra.State() {
					continue
				} else {
					logger.INFO("traverse state is %d", tra.State())
					return resp_list, nil
				}
			} else {
				logger.ERR("requestSeq %d, current pkg num: %d,  %s", requestSeq, idx, err.Error())
				return resp_list, err
			}
		}
	}
}

/**
    @brief 获取遍历器（存在则直接获取，不存在则新建一个）
	@param [IN] zoneId tcaplus请求
	@param [IN] table 超时时间
    @retval *traverser.Traverser 遍历器，一个client最多分配8个遍历器，超过将会返回 nil
**/
func (c *client) GetTraverser(zoneId uint32, table string) *traverser.Traverser {
	return c.tm.GetTraverser(zoneId, table)
}

/*
	@brief 获取本次连接的appId
	@retval int appId
*/
func (c *client) GetAppId() uint64 {
	return c.appId
}

func (c *client) GetProxyUrl(keySet *tcaplus_protocol_cs.TCaplusKeySet, zoneId uint32) string {
	if keySet.FieldNum <= 0 {
		return ""
	}

	field := keySet.Fields[0:keySet.FieldNum]
	sort.Slice(field, func(i, j int) bool {
		if field[i].FieldName < field[j].FieldName {
			return true
		}
		return false
	})

	var buf []byte
	for _, v := range field {
		buf = append(buf, v.FieldBuff[0:v.FieldLen]...)
	}
	if len(buf) <= 0 {
		return ""
	}
	return c.netServer.router.GetProxyUrl(crc32.ChecksumIEEE(buf), zoneId)
}

/*
	@brief 关闭client，释放资源。
	注：关闭接口，关闭各协程的操作是异步的，并不保证在调用之后的下一刻就关闭了所有打开的协程
		不要复用同一个Client以免触发未知bug
*/
func (c *client) Close() {
	c.netServer.stopNetWork <- true
	c.initFlag = NotInit
	c.ctrl.Wait()
}

/*
	@brief 指定IP，主要用于无法访问docker内部ip的情况
*/
func (c *client) SetPublicIP(publicIP string) {
	common.PublicIP = publicIP
}
