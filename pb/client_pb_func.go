package tcaplus

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"google.golang.org/protobuf/proto"
)

/*Client接口的简单封装，方便用户编码,测试阶段仅供参考*/

func (c *PBClient) doSimple(msg proto.Message, apiCmd int, opt *option.PBOpt, zoneId uint32) error {
	table := msg.ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}
	if opt != nil {
		err = c.setReqOpt(req, opt)
		if err != nil {
			logger.ERR("setReqOpt error:%s", err)
			return err
		}
	}

	rec, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return err
	}

	if opt != nil && opt.Version > 0 {
		rec.SetVersion(opt.Version)
	}

	_, err = rec.SetPBData(msg)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return err
	}

	var res response.TcaplusResponse
	if opt != nil && opt.Timeout > 0 {
		res, err = c.Do(req, opt.Timeout)
	} else {
		res, err = c.Do(req, c.defTimeout)
	}
	if err != nil {
		logger.ERR("Do request error:%s", err)
		return err
	}

	ret := res.GetResult()
	if ret != 0 {
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

		if opt != nil {
			opt.Version = record.GetVersion()
		}
	}
	return nil
}

//调用前保证opt中field不为空
func (c *PBClient) doField(msg proto.Message, apiCmd int, opt *option.PBOpt, zoneId uint32) error {
	table := msg.ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	err = c.setReqOpt(req, opt)
	if err != nil {
		logger.ERR("setReqOpt error:%s", err)
		return err
	}

	rec, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return err
	}

	if opt.Version > 0 {
		rec.SetVersion(opt.Version)
	}

	_, err = rec.SetPBFieldValues(msg, opt.FieldNames)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return err
	}

	var res response.TcaplusResponse
	if opt.Timeout > 0 {
		res, err = c.Do(req, opt.Timeout)
	} else {
		res, err = c.Do(req, c.defTimeout)
	}
	if err != nil {
		logger.ERR("Do request error:%s", err)
		return err
	}

	ret := res.GetResult()
	if ret != 0 {
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
		opt.Version = record.GetVersion()
	}
	return nil
}
func (c *PBClient) setReqOpt(req request.TcaplusRequest, opt *option.PBOpt) error {
	var err error
	//req opt
	if opt.VersionPolicy != 0 {
		err = req.SetVersionPolicy(opt.VersionPolicy)
		if err != nil {
			logger.ERR("SetVersionPolicy error:%s", err)
			return err
		}
	}

	if opt.ResultFlag != 0 {
		err = req.SetResultFlag(int(opt.ResultFlag))
		if err != nil {
			logger.ERR("SetResultFlag error:%s", err)
			return err
		}
	} else {
		if opt.ResultFlagForFail != 0 {
			ret := req.SetResultFlagForFail(opt.ResultFlagForFail)
			if ret != 0 {
				err = &terror.ErrorCode{Code: ret}
				logger.ERR("SetResultFlagForFail error:%s", err)
				return err
			}
		}

		if opt.ResultFlagForSuccess != 0 {
			ret := req.SetResultFlagForSuccess(opt.ResultFlagForSuccess)
			if ret != 0 {
				err = &terror.ErrorCode{Code: ret}
				logger.ERR("SetResultFlagForSuccess error:%s", err)
				return err
			}
		}
	}

	if opt.ListShiftFlag != 0 {
		ret := req.SetListShiftFlag(opt.ListShiftFlag)
		if ret != 0 {
			err = &terror.ErrorCode{Code: int(ret)}
			logger.ERR("SetListShiftFlag error:%s", err)
			return err
		}
	}

	if opt.Flags != 0 {
		ret := req.SetFlags(opt.Flags)
		if ret != 0 {
			err = &terror.ErrorCode{Code: int(ret)}
			logger.ERR("SetFlags error:%s", err)
			return err
		}
	}
	return nil
}
func (c *PBClient) setBatchReqOpt(req request.TcaplusRequest, opt *option.PBOpt) error {
	var err error
	//req opt
	if opt.VersionPolicy != 0 {
		err = req.SetVersionPolicy(opt.VersionPolicy)
		if err != nil {
			logger.ERR("SetVersionPolicy error:%s", err)
			return err
		}
	}

	if opt.ResultFlag != 0 {
		err = req.SetResultFlag(int(opt.ResultFlag))
		if err != nil {
			logger.ERR("SetResultFlag error:%s", err)
			return err
		}
	} else {
		if opt.ResultFlagForFail != 0 {
			ret := req.SetResultFlagForFail(opt.ResultFlagForFail)
			if ret != 0 {
				err = &terror.ErrorCode{Code: ret}
				logger.ERR("SetResultFlagForFail error:%s", err)
				return err
			}
		}

		if opt.ResultFlagForSuccess != 0 {
			ret := req.SetResultFlagForSuccess(opt.ResultFlagForSuccess)
			if ret != 0 {
				err = &terror.ErrorCode{Code: ret}
				logger.ERR("SetResultFlagForSuccess error:%s", err)
				return err
			}
		}
	}

	if opt.ListShiftFlag != 0 {
		ret := req.SetListShiftFlag(opt.ListShiftFlag)
		if ret != 0 {
			err = &terror.ErrorCode{Code: int(ret)}
			logger.ERR("SetListShiftFlag error:%s", err)
			return err
		}
	}

	if opt.Flags != 0 {
		ret := req.SetFlags(opt.Flags)
		if ret != 0 {
			err = &terror.ErrorCode{Code: int(ret)}
			logger.ERR("SetFlags error:%s", err)
			return err
		}
	}

	if opt.MultiFlag != 0 {
		req.SetMultiResponseFlag(1)
	}

	if opt.Limit != 0 || opt.Offset != 0 {
		req.SetResultLimit(opt.Limit, opt.Offset)
	}
	return nil
}

func (c *PBClient) doBatch(msgs []proto.Message, versions *[]int32,
	apiCmd int, opt *option.PBOpt, zoneId uint32) error {
	if len(msgs) == 0 {
		logger.ERR("messages is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "messages is nil"}
	}

	table := msgs[0].ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	if versions != nil && len(*versions) != len(msgs) {
		logger.ERR("msgs and versions count not equal")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "msgs and opt.VersionBatch count not equal"}
	}

	if opt != nil {
		err = c.setBatchReqOpt(req, opt)
		if err != nil {
			logger.ERR("setReqOpt error:%s", err)
			return err
		}
	}

	msgMap := make(map[string]int, len(msgs))
	for i, msg := range msgs {
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

		if versions != nil && len(*versions) != 0 && (*versions)[i] > 0 {
			rec.SetVersion((*versions)[i])
		}
		msgMap[string(key)] = i
	}
	var resps []response.TcaplusResponse
	if opt != nil && opt.Timeout > 0 {
		resps, err = c.DoMore(req, opt.Timeout)
	} else {
		resps, err = c.DoMore(req, c.defTimeout)
	}
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

			key, err := record.GetPBKey(nil)
			if err != nil {
				globalErr = err
				logger.ERR("GetPBKey error:%s", err)
				continue
			}

			index, exist := msgMap[string(key)]
			if !exist {
				globalErr = &terror.ErrorCode{Code: terror.RespNotMatchReq}
				logger.ERR("response message is diff request")
				continue
			}

			err = record.GetPBData(msgs[index])
			if err != nil {
				globalErr = err
				logger.ERR("GetPBData error:%s", err)
				continue
			}
			if versions != nil && len(*versions) > 0 {
				(*versions)[index] = record.GetVersion()
			}
		}
	}
	return globalErr
}

/**
    @brief 插入记录,记录存在时报错
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoInsert(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(msg, cmd.TcaplusApiInsertReq, opt, zoneId[0])
	}
	return c.doSimple(msg, cmd.TcaplusApiInsertReq, opt, uint32(c.defZone))
}

/**
    @brief 查询记录
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoGet(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(msg, cmd.TcaplusApiGetReq, opt, zoneId[0])
	}
	return c.doSimple(msg, cmd.TcaplusApiGetReq, opt, uint32(c.defZone))
}

/**
    @brief 替换记录，记录不存在时插入，存在时更新
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoReplace(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(msg, cmd.TcaplusApiReplaceReq, opt, zoneId[0])
	}
	return c.doSimple(msg, cmd.TcaplusApiReplaceReq, opt, uint32(c.defZone))
}

/**
    @brief 更新记录，记录不存在时返错
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoUpdate(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(msg, cmd.TcaplusApiUpdateReq, opt, zoneId[0])
	}
	return c.doSimple(msg, cmd.TcaplusApiUpdateReq, opt, uint32(c.defZone))
}

/**
    @brief 删除记录，记录不存在时返错
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoDelete(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(msg, cmd.TcaplusApiDeleteReq, opt, zoneId[0])
	}
	return c.doSimple(msg, cmd.TcaplusApiDeleteReq, opt, uint32(c.defZone))
}

/**
    @brief 获取部分value,字段通过opt.FieldNames设置
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoFieldGet(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	if opt == nil {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt is nil"}
	}

	if len(opt.FieldNames) == 0 {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt.FieldNames is empty"}
	}

	if len(zoneId) == 1 {
		return c.doField(msg, cmd.TcaplusApiPBFieldGetReq, opt, zoneId[0])
	}
	return c.doField(msg, cmd.TcaplusApiPBFieldGetReq, opt, uint32(c.defZone))
}

/**
    @brief 获取部分value更新
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoFieldUpdate(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	if opt == nil {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt is nil"}
	}

	if len(opt.FieldNames) == 0 {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt.FieldNames is empty"}
	}

	if len(zoneId) == 1 {
		return c.doField(msg, cmd.TcaplusApiPBFieldUpdateReq, opt, zoneId[0])
	}
	return c.doField(msg, cmd.TcaplusApiPBFieldUpdateReq, opt, uint32(c.defZone))
}

/**
    @brief 获取部分value自增
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoFieldIncrease(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	if opt == nil {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt is nil"}
	}

	if len(opt.FieldNames) == 0 {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt.FieldNames is empty"}
	}

	if len(zoneId) == 1 {
		return c.doField(msg, cmd.TcaplusApiPBFieldIncreaseReq, opt, zoneId[0])
	}
	return c.doField(msg, cmd.TcaplusApiPBFieldIncreaseReq, opt, uint32(c.defZone))
}

/**
    @brief 分布式索引查询
	@param [IN] query sql 查询语句 详情见 https://iwiki.woa.com/pages/viewpage.action?pageId=419645505
	@retval []proto.Message 非聚合查询结果
	@retval []string 聚合查询结果
    @retval error 错误码
**/
func (c *PBClient) DoIndexQuery(query string, opt *option.PBOpt, zoneId ...uint32) ([]proto.Message, []string, error) {
	if len(zoneId) == 1 {
		return c.indexQuery(nil, query, cmd.TcaplusApiSqlReq, zoneId[0])
	}
	return c.indexQuery(nil, query, cmd.TcaplusApiSqlReq, uint32(c.defZone))
}

/**
    @brief 同一个表的批量查询
	@param [IN/OUT] msgs proto.Message 由proto文件生成的记录结构体数组, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等
	@param [OUT] versions 记录的版本号，不关心则传nil
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoBatchGet(msgs []proto.Message, versions *[]int32, opt *option.PBOpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doBatch(msgs, versions, cmd.TcaplusApiBatchGetReq, opt, zoneId[0])
	}
	return c.doBatch(msgs, versions, cmd.TcaplusApiBatchGetReq, opt, uint32(c.defZone))
}
