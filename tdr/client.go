package tcaplus

import (
	"errors"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/version"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/request"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/response"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/router"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
	"sync/atomic"
	"time"
)

const (
	NotInit     = 0
	InitSuccess = 1
	InitFail    = 2
)

var reqSeq = uint32(1)

type Client struct {
	appId     uint64
	zoneList  []uint32
	dirUrl    string
	initFlag  int
	netServer netServer
}

/**
   @brief 创建一个tcaplus api客户端
   @retval 返回客户端指针
**/
func NewClient() *Client {
	client := new(Client)
	client.initFlag = NotInit
	return client
}

/**
   @brief                   设置API日志配置文件全路径log.conf(json格式，example下有示例)，请在client.Dial之前调用
   @param [IN] cfgPath      日志配置文件全路径log.conf
   @retval 					错误码
   @note                    Api日志默认使用的zap，用户也可自行实现日志接口logger.LogInterface，调用SetLogger进行设置
**/
func (c *Client) SetLogCfg(cfgPath string) error {
	return logger.SetLogCfg(cfgPath)
}

/**
   @brief                   自定义API日志接口,用户实现logger.LogInterface日志接口，日志将打印到用户的日志接口中，请在client.Dial之前调用
   @param [IN] handle       logger.LogInterface类型的日志接口
   @retval                  错误码
**/
func (c *Client) SetLogger(handle logger.LogInterface) {
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
func (c *Client) Dial(appId uint64, zoneList []uint32, dirUrl string, signature string, timeout uint32) error {
	c.appId = appId
	c.zoneList = make([]uint32, len(zoneList))
	copy(c.zoneList, zoneList)
	c.dirUrl = dirUrl
	//log init
	logger.Init()
	if err := c.netServer.init(appId, zoneList, dirUrl, signature, timeout); err != nil {
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
@brief 创建指定分区指定表的请求
@param [IN] zoneId              分区ID
@param [IN] tableName           表名
@param [IN] cmd                 命令字(cmd.TcaplusApiGetReq等)
@retval request.TcaplusRequest  tcaplus请求
@retval error                   错误码
*/
func (c *Client) NewRequest(zoneId uint32, tableName string, cmd int) (request.TcaplusRequest, error) {
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
			if req, err := request.NewRequest(c.appId, zoneId, tableName, cmd); err != nil {
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
func (c *Client) SendRequest(req request.TcaplusRequest) error {
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
func (c *Client) RecvResponse() (response.TcaplusResponse, error) {
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
func (c *Client) Do(req request.TcaplusRequest, timeout time.Duration) (response.TcaplusResponse, error) {
	requestSeq := int32(atomic.AddUint32(&reqSeq, 1))
	if requestSeq == 0 {
		requestSeq = int32(atomic.AddUint32(&reqSeq, 1))
	}
	req.SetSeq(requestSeq)
	timeOutChan := time.After(timeout)

	var synrequestPkg router.SyncRequest
	synrequestPkg.Init(req)
	defer synrequestPkg.SyncChanClose()

	c.netServer.router.RequestChanMapAdd(requestSeq, synrequestPkg)
	defer c.netServer.router.RequestChanMapClean(requestSeq)

	if err := c.SendRequest(req); err != nil {
		logger.ERR("requestSeq %d :SendRequest failed %v\n", requestSeq, err.Error())
		return nil, err
	}

	for {
		select {
		case <-timeOutChan:
			logger.ERR("requestSeq %d :%s, timeout", requestSeq, timeout.String())
			return nil, errors.New(timeout.String() + ", timeout")
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
func (c *Client) DoMore(req request.TcaplusRequest, timeout time.Duration) ([]response.TcaplusResponse, error) {
	requestSeq := int32(atomic.AddUint32(&reqSeq, 1))
	if requestSeq == 0 {
		requestSeq = int32(atomic.AddUint32(&reqSeq, 1))
	}
	req.SetSeq(requestSeq)
	timeOutChan := time.After(timeout)

	var synrequestPkg router.SyncRequest
	synrequestPkg.InitMoreChan(req, 1024)
	defer synrequestPkg.SyncChanClose()

	c.netServer.router.RequestChanMapAdd(requestSeq, synrequestPkg)
	defer c.netServer.router.RequestChanMapClean(requestSeq)

	if err := c.SendRequest(req); err != nil {
		logger.ERR("requestSeq %d :SendRequest failed %v\n", requestSeq, err.Error())
		return nil, err
	}
	var resp_list []response.TcaplusResponse
	var idx int = 0
	for {
		select {
		case <-timeOutChan:
			logger.ERR("requestSeq %d :%s, timeout, current pkg num %d", requestSeq, timeout.String(), idx)
			return resp_list, errors.New(timeout.String() + ", timeout")
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
