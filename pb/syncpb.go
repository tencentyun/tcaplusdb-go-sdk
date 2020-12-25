package tcaplus

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/metadata"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
	"strings"
	"time"
)

type PBClient struct {
	*client
	defZone    int32
	defTimeout time.Duration
}

func NewPBClient() *PBClient {
	c := new(PBClient)
	c.client = newClient()
	c.defZone = -1
	c.defTimeout = 5 * time.Second
	return c
}

// 连接 dir proxy 初始化 注意在Dial前初始化 logger 否则会打到控制台
func (c *PBClient) Dial(appId uint64, zoneList []uint32, dirUrl string, signature string, timeout uint32, zoneTables map[uint32][]string) error {
	err := c.client.Dial(appId, zoneList, dirUrl, signature, timeout)
	if err != nil {
		return err
	}
	return c.initTableMeta(zoneTables)
}

// 初始化表元数据
func (c *PBClient) initTableMeta(zoneTables map[uint32][]string) error {
	if len(zoneTables) == 0 {
		zoneTables = make(map[uint32][]string)
		for _, zone := range c.zoneList {
			zoneTables[zone] = c.netServer.router.GetZoneTables(zone)
		}
	}

	for zone, tables := range zoneTables {
		if c.defZone == -1 {
			c.defZone = int32(zone)
			logger.DEBUG("init default zone %d", c.defZone)
		}
		for _, table := range tables {
			req, err := c.NewRequest(zone, table, cmd.TcaplusApiMetadataGetReq)
			if err != nil {
				logger.ERR("NewRequest error:%s")
				return err
			}
			resp, err := c.Do(req, c.defTimeout*time.Second)
			if err != nil {
				logger.ERR("Do request error:%s", err)
				return err
			}
			if r := resp.GetResult(); r != 0 {
				errMsg := fmt.Sprintf("get zone %d table %s metadata error:%s", zone, table,
					terror.GetErrMsg(r))
				logger.ERR(errMsg)
				return &terror.ErrorCode{Code: r, Message: errMsg}
			}
			if resp.GetTcaplusPackagePtr() == nil {
				errMsg := fmt.Sprintf("get zone %d table %s metadata error:response pkg is nil", zone, table)
				logger.ERR(errMsg)
				return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH, Message: errMsg}
			}
			metares := resp.GetTcaplusPackagePtr().Body.MetadataGetRes
			if metares.IdlType != 2 {
				errMsg := fmt.Sprintf("get zone %d table %s metadata error:table type %d not proto",
					zone, table, metares.IdlType)
				logger.ERR(errMsg)
				return &terror.ErrorCode{Code: terror.MetadataNotProtobuf, Message: errMsg}
			}
			err = metadata.GetMetaManager().AddTableDesGrp(c.appId, zone, table,
				metares.IdlContent[:metares.IdlConLen])
			if err != nil {
				errMsg := fmt.Sprintf("add app %d zone %d table %s metadata error:%s",
					c.appId, zone, table, err)
				logger.ERR(errMsg)
				return err
			}
		}
	}
	return nil
}

/**
    @brief 设置默认zoneId
	@param [IN] zoneId zoneID
    @retval error 错误码，如果未dial调用此接口将会返错 ClientNotDial
**/
func (c *PBClient) SetDefaultZoneId(zoneId uint32) error {
	// 等于 -1 说明未进行dial初始化
	if c.defZone == -1 {
		logger.ERR("client not dial init")
		return &terror.ErrorCode{Code: terror.ClientNotDial}
	}
	c.defZone = int32(zoneId)
	return nil
}

/**
    @brief 设置默认超时时间
	@param [IN] t time.Duration
    @retval error 错误码，如果未dial调用此接口将会返错 ClientNotDial
**/
func (c *PBClient) SetDefaultTimeOut(t time.Duration) error {
	if c.defZone == -1 {
		logger.ERR("client not dial init")
		return &terror.ErrorCode{Code: terror.ClientNotDial}
	}
	c.defTimeout = t
	return nil
}

func (c *PBClient) simpleOperate(msg proto.Message, apicmd int, zoneId uint32) error {
	if c.defZone == -1 {
		logger.ERR("client not dial init")
		return &terror.ErrorCode{Code: terror.ClientNotDial}
	}

	table := msg.ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apicmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	rec, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return err
	}

	_, err = rec.SetPBData(msg)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return err
	}

	req.SetResultFlag(3)

	res, err := c.Do(req, c.defTimeout)
	if err != nil {
		logger.ERR("Do request error:%s", err)
		return err
	}

	ret := res.GetResult()
	if ret != 0 {
		//logger.ERR("result is %d, error:%s", ret, terror.GetErrMsg(ret))
		return &terror.ErrorCode{Code: ret}
	}

	if res.GetRecordCount() > 0 {
		record, err := res.FetchRecord()
		if err != nil {
			logger.ERR("FetchRecord error:%s", err)
			return err
		}

		err = record.GetPBData(msg)
		if err != nil {
			logger.ERR("GetPBData error:%s", err)
			return err
		}
	}

	return nil
}

func (c *PBClient) batchOperate(msgs []proto.Message, apicmd int, zoneId uint32) error {
	if c.defZone == -1 {
		logger.ERR("client not dial init")
		return &terror.ErrorCode{Code: terror.ClientNotDial}
	}

	if len(msgs) == 0 {
		logger.ERR("messages is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "messages is nil"}
	}

	table := msgs[0].ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apicmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	msgMap := make(map[string]proto.Message, len(msgs))

	for _, msg := range msgs {
		rec, err := req.AddRecord(0)
		if err != nil {
			logger.ERR("AddRecord error:%s", err)
			return err
		}

		key, err := rec.SetPBData(msg)
		if err != nil {
			logger.ERR("SetPBData error:%s", err)
			return nil
		}

		msgMap[string(key)] = msg
	}

	req.SetResultFlag(2)
	req.SetMultiResponseFlag(1)

	resps, err := c.DoMore(req, c.defTimeout)
	if err != nil {
		logger.ERR("DoMore request error:%s", err)
		return err
	}

	var globalErr error

	for _, res := range resps {
		ret := res.GetResult()
		if ret != 0 {
			globalErr = &terror.ErrorCode{Code: ret}
			logger.ERR("result is %d, error:%s", ret, terror.GetErrMsg(ret))
			continue
		}

		for i := 0; i < res.GetRecordCount(); i++ {
			record, err := res.FetchRecord()
			if err != nil {
				globalErr = err
				logger.ERR("FetchRecord error:%s", err)
				continue
			}

			key, err := record.GetPBKey()
			if err != nil {
				globalErr = err
				logger.ERR("GetPBKey error:%s", err)
				continue
			}

			msg, exist := msgMap[string(key)]
			if !exist {
				globalErr = &terror.ErrorCode{Code: terror.RespNotMatchReq}
				logger.ERR("response message is diff request")
				continue
			}

			err = record.GetPBData(msg)
			if err != nil {
				globalErr = err
				logger.ERR("GetPBData error:%s", err)
				continue
			}
		}
	}

	return globalErr
}

func (c *PBClient) partkeyOperate(msg proto.Message, keys []string, apicmd int, zoneId uint32) ([]proto.Message, error) {
	if c.defZone == -1 {
		logger.ERR("client not dial init")
		return nil, &terror.ErrorCode{Code: terror.ClientNotDial}
	}

	table := msg.ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apicmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return nil, err
	}

	rec, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return nil, err
	}

	_, err = rec.SetPBPartKeys(msg, keys)
	if err != nil {
		logger.ERR("SetPBPartKeys error:%s", err)
		return nil, err
	}

	req.SetMultiResponseFlag(1)

	resps, err := c.DoMore(req, c.defTimeout)
	if err != nil {
		logger.ERR("DoMore request error:%s", err)
		return nil, err
	}

	var msgs []proto.Message
	var globalErr error

	for _, res := range resps {
		ret := res.GetResult()
		if ret != 0 {
			globalErr = &terror.ErrorCode{Code: ret}
			logger.ERR("result is %d, error:%s", ret, terror.GetErrMsg(ret))
			continue
		}

		for i := 0; i < res.GetRecordCount(); i++ {
			record, err := res.FetchRecord()
			if err != nil {
				globalErr = err
				logger.ERR("FetchRecord error:%s", err)
				continue
			}

			err = record.GetPBData(msg)
			if err != nil {
				globalErr = err
				logger.ERR("GetPBData error:%s", err)
				continue
			}

			msgs = append(msgs, proto.Clone(msg))
		}
	}

	return msgs, globalErr
}

func (c *PBClient) fieldOperate(msg proto.Message, values []string, apicmd int, zoneId uint32) error {
	if c.defZone == -1 {
		logger.ERR("client not dial init")
		return &terror.ErrorCode{Code: terror.ClientNotDial}
	}

	table := msg.ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apicmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	rec, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return err
	}

	_, err = rec.SetPBFieldValues(msg, values)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return err
	}

	res, err := c.Do(req, c.defTimeout)
	if err != nil {
		logger.ERR("Do request error:%s", err)
		return err
	}

	ret := res.GetResult()
	if ret != 0 {
		logger.ERR("result is %d, error:%s", ret, terror.GetErrMsg(ret))
		return &terror.ErrorCode{Code: ret}
	}

	if res.GetRecordCount() > 0 {
		record, err := res.FetchRecord()
		if err != nil {
			logger.ERR("FetchRecord error:%s", err)
			return err
		}

		err = record.GetPBFieldValues(msg)
		if err != nil {
			logger.ERR("GetPBData error:%s", err)
			return err
		}
	}

	return nil
}

func (c *PBClient) indexQuery(query string, apicmd int, zoneId uint32) ([]proto.Message, []string, error) {
	if c.defZone == -1 {
		logger.ERR("client not dial init")
		return nil, nil, &terror.ErrorCode{Code: terror.ClientNotDial}
	}

	// select * from table where a>100;
	// 至少分为 6 段
	conditions := strings.Fields(query)
	if len(conditions) < 6 {
		logger.ERR("sql field length %d less 6", len(conditions))
		return nil, nil, &terror.ErrorCode{Code: terror.SqlQueryFormatError}
	}

	// 以select开头，忽略大小写
	if !strings.EqualFold(conditions[0], "select") {
		logger.ERR("sql first word not select")
		return nil, nil, &terror.ErrorCode{Code: terror.SqlQueryFormatError}
	}

	var i = int(2)
	for i < len(conditions) {
		if strings.EqualFold(conditions[i], "from") {
			break
		}
		i++
	}

	// from 不能没有或者位于最后
	if i >= len(conditions)-1 {
		logger.ERR("sql not find from or from at last")
		return nil, nil, &terror.ErrorCode{Code: terror.SqlQueryFormatError}
	}

	fieldstr := ""
	var fields []string
	for j := 1; j < i; j++ {
		fieldstr += conditions[j]
	}
	if fieldstr != "*" {
		fields = strings.Split(fieldstr, ",")
	}
	logger.DEBUG("select fields %+v", fields)

	table := conditions[i+1]
	req, err := c.NewRequest(zoneId, table, apicmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return nil, nil, err
	}

	ret := req.SetSql(query)
	if ret != 0 {
		logger.ERR("SetSql %s", terror.GetErrMsg(ret))
		return nil, nil, &terror.ErrorCode{Code: ret}
	}

	resps, err := c.DoMore(req, c.defTimeout)
	if err != nil {
		logger.ERR("DoMore request error:%s", err)
		return nil, nil, err
	}

	sqlType := resps[0].GetSqlType()
	if sqlType == policy.AGGREGATIONS_SQL_QUERY_TYPE {
		res, err := resps[0].ProcAggregationSqlQueryType()
		if err != nil {
			logger.ERR("ProcAggregationSqlQueryType error:%s", err)
			return nil, nil, err
		}
		return nil, res, nil
	}

	zoneTable := fmt.Sprintf("%d|%d|%s", c.appId, c.defZone, table)
	grp := metadata.GetMetaManager().GetTableDesGrp(zoneTable)
	if grp == nil {
		logger.ERR("not find zoneTable %s", zoneTable)
		return nil, nil, &terror.ErrorCode{Code: terror.ParameterInvalid}
	}

	var msgs []proto.Message
	var globalErr error

	for _, res := range resps {
		ret = res.GetResult()
		if ret != 0 {
			globalErr = &terror.ErrorCode{Code: ret}
			logger.ERR("result is %d, error:%s", ret, terror.GetErrMsg(ret))
			continue
		}

		for i := 0; i < res.GetRecordCount(); i++ {
			record, err := res.FetchRecord()
			if err != nil {
				globalErr = err
				logger.ERR("FetchRecord error:%s", err)
				continue
			}

			msg := dynamicpb.NewMessage(grp.Desc)
			err = record.GetPBDataWithValues(msg, fields)
			if err != nil {
				globalErr = err
				logger.ERR("GetPBData error:%s", err)
				continue
			}

			msgs = append(msgs, msg)
		}
	}

	return msgs, nil, globalErr
}

/**
    @brief 插入记录
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (c *PBClient) Insert(msg proto.Message) error {
	return c.simpleOperate(msg, cmd.TcaplusApiInsertReq, uint32(c.defZone))
}

/**
    @brief 插入记录
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) InsertWithZone(msg proto.Message, zoneId uint32) error {
	return c.simpleOperate(msg, cmd.TcaplusApiInsertReq, zoneId)
}

/**
    @brief 替换记录，记录不存在时插入
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (c *PBClient) Replace(msg proto.Message) error {
	return c.simpleOperate(msg, cmd.TcaplusApiReplaceReq, uint32(c.defZone))
}

/**
    @brief 替换记录，记录不存在时插入。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) ReplaceWithZone(msg proto.Message, zoneId uint32) error {
	return c.simpleOperate(msg, cmd.TcaplusApiReplaceReq, zoneId)
}

/**
    @brief 获取记录
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (c *PBClient) Get(msg proto.Message) error {
	return c.simpleOperate(msg, cmd.TcaplusApiGetReq, uint32(c.defZone))
}

/**
    @brief 获取记录。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) GetWithZone(msg proto.Message, zoneId uint32) error {
	return c.simpleOperate(msg, cmd.TcaplusApiGetReq, zoneId)
}

/**
    @brief 删除记录
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (c *PBClient) Delete(msg proto.Message) error {
	return c.simpleOperate(msg, cmd.TcaplusApiDeleteReq, uint32(c.defZone))
}

/**
    @brief 删除记录。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) DeleteWithZone(msg proto.Message, zoneId uint32) error {
	return c.simpleOperate(msg, cmd.TcaplusApiDeleteReq, zoneId)
}

/**
    @brief 修改记录，记录不存在时返错
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (c *PBClient) Update(msg proto.Message) error {
	return c.simpleOperate(msg, cmd.TcaplusApiUpdateReq, uint32(c.defZone))
}

/**
    @brief 修改记录，记录不存在时返错。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) UpdateWithZone(msg proto.Message, zoneId uint32) error {
	return c.simpleOperate(msg, cmd.TcaplusApiUpdateReq, zoneId)
}

/**
    @brief 批量获取记录
	@param [IN] msgs []proto.Message 需获取的记录列表
    @retval error 错误码
**/
func (c *PBClient) BatchGet(msgs []proto.Message) error {
	return c.batchOperate(msgs, cmd.TcaplusApiBatchGetReq, uint32(c.defZone))
}

/**
    @brief 批量获取记录。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msgs []proto.Message 需获取的记录列表
	@param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) BatchGetWithZone(msgs []proto.Message, zoneId uint32) error {
	return c.batchOperate(msgs, cmd.TcaplusApiBatchGetReq, zoneId)
}

/**
    @brief 部分key获取记录
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] keys []string 部分key，根据 proto 文件中的 index 选择填写
	@retval []proto.Message 返回记录，可能匹配到多条记录
    @retval error 错误码
**/
func (c *PBClient) GetByPartKey(msg proto.Message, keys []string) ([]proto.Message, error) {
	return c.partkeyOperate(msg, keys, cmd.TcaplusApiGetByPartkeyReq, uint32(c.defZone))
}

/**
    @brief 部分key获取记录。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] keys []string 部分key，根据 proto 文件中的 index 选择填写
	@param [IN] zoneId 指定表所在zone
	@retval []proto.Message 返回记录，可能匹配到多条记录
    @retval error 错误码
**/
func (c *PBClient) GetByPartKeyWithZone(msg proto.Message, keys []string, zoneId uint32) ([]proto.Message, error) {
	return c.partkeyOperate(msg, keys, cmd.TcaplusApiGetByPartkeyReq, zoneId)
}

/**
    @brief 获取记录部分字段value
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] values []string 部分字段名，根据需要选择填写
    @retval error 错误码
**/
func (c *PBClient) FieldGet(msg proto.Message, values []string) error {
	return c.fieldOperate(msg, values, cmd.TcaplusApiPBFieldGetReq, uint32(c.defZone))
}

/**
    @brief 获取记录部分字段value。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] values []string 部分字段名，根据需要选择填写
	@param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) FieldGetWithZone(msg proto.Message, values []string, zoneId uint32) error {
	return c.fieldOperate(msg, values, cmd.TcaplusApiPBFieldGetReq, zoneId)
}

/**
    @brief 更新记录部分字段value
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] values []string 部分字段名，根据需要选择填写
    @retval error 错误码
**/
func (c *PBClient) FieldUpdate(msg proto.Message, values []string) error {
	return c.fieldOperate(msg, values, cmd.TcaplusApiPBFieldUpdateReq, uint32(c.defZone))
}

/**
    @brief 更新记录部分字段value。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] values []string 部分字段名，根据需要选择填写
	@param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) FieldUpdateWithZone(msg proto.Message, values []string, zoneId uint32) error {
	return c.fieldOperate(msg, values, cmd.TcaplusApiPBFieldUpdateReq, zoneId)
}

/**
    @brief 自增记录部分字段value
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] values []string 部分字段名，根据需要选择填写
    @retval error 错误码
**/
func (c *PBClient) FieldIncrease(msg proto.Message, values []string) error {
	return c.fieldOperate(msg, values, cmd.TcaplusApiPBFieldIncreaseReq, uint32(c.defZone))
}

/**
    @brief 自增记录部分字段value。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] values []string 部分字段名，根据需要选择填写
	@param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) FieldIncreaseWithZone(msg proto.Message, values []string, zoneId uint32) error {
	return c.fieldOperate(msg, values, cmd.TcaplusApiPBFieldIncreaseReq, zoneId)
}

/**
    @brief 分布式索引查询
	@param [IN] query sql 查询语句 详情见 https://iwiki.woa.com/pages/viewpage.action?pageId=419645505
	@retval []proto.Message 非聚合查询结果
	@retval []string 聚合查询结果
    @retval error 错误码
**/
func (c *PBClient) IndexQuery(query string) ([]proto.Message, []string, error) {
	return c.indexQuery(query, cmd.TcaplusApiSqlReq, uint32(c.defZone))
}

/**
    @brief 自增记录部分字段value
	@param [IN] query sql 查询语句 详情见 https://iwiki.woa.com/pages/viewpage.action?pageId=419645505
	@param [IN] zoneId 指定表所在zone
	@retval []proto.Message 非聚合查询结果
	@retval []string 聚合查询结果
    @retval error 错误码
**/
func (c *PBClient) IndexQueryWithZone(query string, zoneId uint32) ([]proto.Message, []string, error) {
	return c.indexQuery(query, cmd.TcaplusApiSqlReq, zoneId)
}
