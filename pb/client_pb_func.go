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

/*Client接口的简单封装，方便用户编码*/
func (c *PBClient) DoListSimple(msg proto.Message, index int32, apiCmd int, opt *option.PBOpt, zoneId uint32) (int32,
	error) {
	table := msg.ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return index, err
	}

	rec, err := req.AddRecord(index)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return index, err
	}

	_, err = rec.SetPBData(msg)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return index, err
	}

	if opt != nil {
		err = c.setReqOpt(req, opt)
		if err != nil {
			logger.ERR("setReqOpt error:%s", err)
			return index, err
		}
		if opt.Version > 0 {
			rec.SetVersion(opt.Version)
		}
	}

	var res response.TcaplusResponse
	if opt != nil && opt.Timeout > 0 {
		res, err = c.Do(req, opt.Timeout)
	} else {
		res, err = c.Do(req, c.defTimeout)
	}
	if err != nil {
		logger.ERR("Do request error:%s", err)
		return index, err
	}

	ret := res.GetResult()
	if ret != 0 {
		return index, &terror.ErrorCode{Code: ret}
	}

	if res.GetRecordCount() > 0 {
		record, err := res.FetchRecord()
		if err != nil {
			logger.ERR("FetchRecord error:%s", err)
			return index, err
		}

		if opt != nil {
			opt.Version = record.GetVersion()
		}
		index = record.GetIndex()

		if apiCmd != cmd.TcaplusApiListGetReq && !c.needGetData(opt) {
			return index, nil
		}

		err = record.GetPBData(msg)
		if err != nil {
			logger.ERR("GetPBData error:%s", err)
			return index, err
		}
	}
	return index, nil
}

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

		if opt != nil {
			opt.Version = record.GetVersion()
		}

		if apiCmd != cmd.TcaplusApiGetReq && !c.needGetData(opt) {
			return nil
		}

		err = record.GetPBData(msg)
		if err != nil {
			logger.ERR("GetPBData error:%s", err)
			return err
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
func (c *PBClient) doListBatch(msg proto.Message, indexs []int32,
	apiCmd int, opt *option.PBOpt, zoneId uint32) (map[int32]proto.Message, error) {

	table := msg.ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return nil, err
	}

	rec, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return nil, err
	}

	_, err = rec.SetPBData(msg)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return nil, err
	}

	if opt != nil {
		if opt.BatchResult == nil && len(indexs) > 0 {
			opt.BatchResult = make([]error, len(indexs), len(indexs))
		} else if len(opt.BatchResult) != len(indexs) && len(indexs) > 0 {
			logger.ERR("indexs and BatchResult count not equal")
			return nil, &terror.ErrorCode{Code: terror.ParameterInvalid,
				Message: "indexs and opt.BatchResult count not equal"}
		}
		err = c.setBatchReqOpt(req, opt)
		if err != nil {
			logger.ERR("setReqOpt error:%s", err)
			return nil, err
		}
	}

	tmpIdxMap := map[int32]struct{}{}
	for _, index := range indexs {
		if _, exist := tmpIdxMap[index]; exist{
			logger.ERR("batch record exist duplicate index")
			return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "batch record exist duplicate index"}
		}
		tmpIdxMap[index] = struct{}{}
		req.AddElementIndex(index)
	}

	var resps []response.TcaplusResponse
	if opt != nil && opt.Timeout > 0 {
		resps, err = c.DoMore(req, opt.Timeout)
	} else {
		resps, err = c.DoMore(req, c.defTimeout)
	}
	if err != nil {
		logger.ERR("Do request error:%s", err)
		return nil, err
	}

	msgs := make(map[int32]proto.Message)
	var globalErr error
	offset := 0
	for _, res := range resps {
		ret := res.GetResult()
		if ret != 0 {
			globalErr = &terror.ErrorCode{Code: ret}
			logger.ERR("result is %d, error:%s", ret, terror.GetErrMsg(ret))
			continue
		}

		for i := 0; i < res.GetRecordCount(); i++ {
			record, err := res.FetchRecord()
			if opt != nil && offset < len(indexs) {
				opt.BatchResult[offset] = err
			}
			if err != nil {
				globalErr = err
				logger.ERR("FetchRecord error:%s", err)
				offset++
				continue
			}

			err = record.GetPBData(msg)
			if err != nil {
				globalErr = err
				logger.ERR("GetPBData error:%s", err)
				offset++
				continue
			}

			msgs[record.GetIndex()] = proto.Clone(msg)
			if opt != nil {
				opt.Version = record.GetVersion()
			}
			offset++
		}
	}

	return msgs, globalErr
}

func (c *PBClient) doListBatchRecord(msgs []proto.Message, indexs []int32,
	apiCmd int, opt *option.PBOpt, zoneId uint32) error {
	if len(msgs) != len(indexs) || len(msgs) == 0 {
		logger.ERR("len(dataSlice) %d != len(indexs) %d or empty", len(msgs), len(indexs))
		return terror.ErrorCode{Code: terror.ParameterInvalid, Message: "len(msgs) != len(indexs) or empty"}
	}

	table := msgs[0].ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	if opt != nil {
		if opt.BatchResult == nil {
			opt.BatchResult = make([]error, len(msgs), len(msgs))
		} else if len(opt.BatchResult) != len(msgs) {
			logger.ERR("dataSlice and BatchResult count not equal")
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dataSlice and opt.BatchResult count not equal"}
		}

		err = c.setBatchReqOpt(req, opt)
		if err != nil {
			logger.ERR("setReqOpt error:%s", err)
			return err
		}
	}

	for i, data := range msgs {
		rec, err := req.AddRecord(indexs[i])
		if err != nil {
			logger.ERR("AddRecord error:%s", err)
			return err
		}

		_, err = rec.SetPBData(data)
		if err != nil {
			logger.ERR("SetData error:%s", err)
			return err
		}
		if opt != nil && opt.Version > 0 {
			rec.SetVersion(opt.Version)
		}
	}

	var resps []response.TcaplusResponse
	if opt != nil && opt.Timeout > 0 {
		resps, err = c.DoMore(req, opt.Timeout)
	} else {
		resps, err = c.DoMore(req, c.defTimeout)
	}
	if err != nil {
		logger.ERR("Do request error:%s", err)
		return err
	}

	var globalErr error
	offset := 0
	for _, res := range resps {
		ret := res.GetResult()
		if ret != 0 {
			globalErr = &terror.ErrorCode{Code: ret}
			logger.ERR("result is %d, error:%s", ret, terror.GetErrMsg(ret))
			continue
		}

		for i := 0; i < res.GetRecordCount(); i++ {
			resRec, err := res.FetchRecord()
			if opt != nil && offset < len(indexs) {
				opt.BatchResult[offset] = err
			}
			if err != nil {
				globalErr = err
				logger.ERR("FetchRecord error:%s", err)
				offset++
				continue
			}

			if resRec == nil {
				offset++
				continue
			}

			indexs[offset] = resRec.GetIndex()
			if opt != nil {
				opt.Version = resRec.GetVersion()
				if c.needGetData(opt) {
					err = resRec.GetPBData(msgs[offset])
					if err != nil {
						globalErr = err
						logger.ERR("FetchRecord error:%s", err)
						offset++
						continue
					}
				}
			}
			offset++
		}
	}
	return globalErr
}

func (c *PBClient) needGetData(opt *option.PBOpt) bool {
	if opt == nil {
		return false
	}

	if opt.ResultFlag == option.TcaplusResultFlagAllNewValue || opt.ResultFlag == option.TcaplusResultFlagAllOldValue ||
		opt.ResultFlagForSuccess == option.TcaplusResultFlagAllNewValue ||
		opt.ResultFlagForSuccess == option.TcaplusResultFlagAllOldValue ||
		opt.ResultFlagForFail == option.TcaplusResultFlagAllNewValue ||
		opt.ResultFlagForFail == option.TcaplusResultFlagAllOldValue {
		return true
	}

	return false
}
func (c *PBClient) doBatch(msgs []proto.Message, apiCmd int, opt *option.PBOpt, zoneId uint32) error {
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

	if opt != nil {
		if opt.BatchVersion == nil {
			opt.BatchVersion = make([]int32, len(msgs), len(msgs))
		} else if len(opt.BatchVersion) != len(msgs) {
			logger.ERR("msgs and BatchVersion count not equal")
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "msgs and opt.BatchVersion count not equal"}
		}

		if opt.BatchResult == nil {
			opt.BatchResult = make([]error, len(msgs), len(msgs))
		} else if len(opt.BatchResult) != len(msgs) {
			logger.ERR("msgs and BatchResult count not equal")
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "msgs and opt.BatchResult count not equal"}
		}

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
			return err
		}

		if opt != nil && (opt.BatchVersion)[i] > 0 {
			rec.SetVersion((opt.BatchVersion)[i])
		}

		if _, exist := msgMap[string(key)]; exist {
			logger.ERR("batch record exist duplicate key")
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "batch record exist duplicate key"}
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
			record, recErr := res.FetchRecord()
			if recErr != nil {
				globalErr = recErr
				logger.DEBUG("FetchRecord error:%s", recErr)
				continue
			}

			if record == nil {
				continue
			}

			key, err := record.GetPBKey(nil)
			if err != nil {
				globalErr = err
				logger.ERR("GetPBKey error:%s", err)
				continue
			}

			keyStr := string(key)
			index, exist := msgMap[keyStr]
			if !exist {
				globalErr = &terror.ErrorCode{Code: terror.RespNotMatchReq}
				logger.ERR("response message is diff request")
				continue
			}
			if opt != nil {
				opt.BatchResult[index] = recErr
			}
			delete(msgMap, keyStr)
			if recErr != nil {
				continue
			}

			if opt != nil {
				opt.BatchVersion[index] = record.GetVersion()
			}

			if apiCmd != cmd.TcaplusApiBatchGetReq && !c.needGetData(opt) {
				continue
			}

			err = record.GetPBData(msgs[index])
			if err != nil {
				globalErr = err
				logger.ERR("GetPBData key %s error:%s", keyStr, err)
				continue
			}
		}
	}

	//msgMap not nil
	if len(msgMap) != 0 && globalErr == nil {
		globalErr = &terror.ErrorCode{Code: terror.NoRspWithTheKeyReq,
			Message: "no rsp with key"}
	}

	for key, index := range msgMap {
		logger.ERR("key %s offset %d not rsp", key, index)
		if opt != nil {
			opt.BatchResult[index] = globalErr
			opt.BatchVersion[index] = -1
		}
	}

	return globalErr
}

func (c *PBClient) doPartKeyGet(msg proto.Message, keys []string, apiCmd int, opt *option.PBOpt,
	zoneId uint32) ([]proto.Message, error) {
	table := msg.ProtoReflect().Descriptor().Name()
	req, err := c.NewRequest(zoneId, string(table), apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return nil, err
	}

	if opt != nil {
		err = c.setBatchReqOpt(req, opt)
		if err != nil {
			logger.ERR("setReqOpt error:%s", err)
			return nil, err
		}
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

	var resps []response.TcaplusResponse
	if opt != nil && opt.Timeout > 0 {
		resps, err = c.DoMore(req, opt.Timeout)
	} else {
		resps, err = c.DoMore(req, c.defTimeout)
	}
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
			if opt != nil {
				opt.BatchVersion = append(opt.BatchVersion, record.GetVersion())
			}
		}
	}

	if len(msgs) == 0 && globalErr == nil {
		return nil, &terror.ErrorCode{Code: terror.TXHDB_ERR_RECORD_NOT_EXIST}
	}

	return msgs, globalErr
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
    @brief 同一个表的批量查询
	@param [IN/OUT] msgs proto.Message 由proto文件生成的记录结构体数组, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数, 分包返回等, 记录的version存放在opt.BatchVersion，单条记录结果存放opt.BatchResult
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoBatchGet(msgs []proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doBatch(msgs, cmd.TcaplusApiBatchGetReq, opt, zoneId[0])
	}
	return c.doBatch(msgs, cmd.TcaplusApiBatchGetReq, opt, uint32(c.defZone))
}

/**
    @brief 根据表的部分key字段查询
	@param [IN] msgs proto.Message 由proto文件生成的记录结构体
	@param [IN] indexKeys 使用的索引部分key字段名称
	@param [IN/OUT] opt 可选参数, 分包返回等, 记录的version存放在opt.BatchVersion
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoGetByPartKey(msgs proto.Message, indexKeys []string, opt *option.PBOpt,
	zoneId ...uint32) ([]proto.Message,
	error) {
	if len(indexKeys) == 0 {
		return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "partKeys is empty"}
	}
	if len(zoneId) == 1 {
		return c.doPartKeyGet(msgs, indexKeys, cmd.TcaplusApiGetByPartkeyReq, opt, zoneId[0])
	}
	return c.doPartKeyGet(msgs, indexKeys, cmd.TcaplusApiGetByPartkeyReq, opt, uint32(c.defZone))
}

/**
    @brief list表插入记录
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] index int32 插入到key中的第index条记录之后
				tcaplus_protocol_cs.TCAPLUS_LIST_LAST_INDEX = -1      插入元素位置在最后面
				tcaplus_protocol_cs.TCAPLUS_LIST_PRE_FIRST_INDEX = -2 插入元素位置在最前面
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval int32 如果有返回记录，返回index 索引
    @retval error 错误码
**/
func (c *PBClient) DoListAddAfter(msg proto.Message, index int32, opt *option.PBOpt, zoneId ...uint32) (int32, error) {
	if len(zoneId) == 1 {
		return c.DoListSimple(msg, index, cmd.TcaplusApiListAddAfterReq, opt, zoneId[0])
	}
	return c.DoListSimple(msg, index, cmd.TcaplusApiListAddAfterReq, opt, uint32(c.defZone))
}

/**
    @brief list表删除记录
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] index int32 操作第index条记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoListDelete(msg proto.Message, index int32, opt *option.PBOpt, zoneId ...uint32) error {
	var err error
	if len(zoneId) == 1 {
		_, err = c.DoListSimple(msg, index, cmd.TcaplusApiListDeleteReq, opt, zoneId[0])
	} else {
		_, err = c.DoListSimple(msg, index, cmd.TcaplusApiListDeleteReq, opt, uint32(c.defZone))
	}
	return err
}

/**
    @brief list表更新记录
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] index int32 操作第index条记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoListReplace(msg proto.Message, index int32, opt *option.PBOpt, zoneId ...uint32) error {
	var err error
	if len(zoneId) == 1 {
		_, err = c.DoListSimple(msg, index, cmd.TcaplusApiListReplaceReq, opt, zoneId[0])
	} else {
		_, err = c.DoListSimple(msg, index, cmd.TcaplusApiListReplaceReq, opt, uint32(c.defZone))
	}
	return err
}

/**
    @brief list表查询记录
	@param [IN/OUT] msg proto.Message 由proto文件生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] index int32 操作第index条记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoListGet(msg proto.Message, index int32, opt *option.PBOpt, zoneId ...uint32) error {
	var err error
	if len(zoneId) == 1 {
		_, err = c.DoListSimple(msg, index, cmd.TcaplusApiListGetReq, opt, zoneId[0])
	} else {
		_, err = c.DoListSimple(msg, index, cmd.TcaplusApiListGetReq, opt, uint32(c.defZone))
	}
	return err
}

/**
    @brief list表删除key下所有记录
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN/OUT] opt 可选参数，乐观锁等
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoListDeleteAll(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) error {
	var err error
	if len(zoneId) == 1 {
		_, err = c.doListBatch(msg, nil, cmd.TcaplusApiListDeleteAllReq, opt, zoneId[0])
	} else {
		_, err = c.doListBatch(msg, nil, cmd.TcaplusApiListDeleteAllReq, opt, uint32(c.defZone))
	}
	return err
}

/**
    @brief 批量删除list表下记录
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] indexs []int32 删除key下多个记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoListDeleteBatch(msg proto.Message, indexs []int32, opt *option.PBOpt,
	zoneId ...uint32) (map[int32]proto.Message, error) {
	if len(zoneId) == 1 {
		return c.doListBatch(msg, indexs, cmd.TcaplusApiListDeleteBatchReq, opt, zoneId[0])
	}
	return c.doListBatch(msg, indexs, cmd.TcaplusApiListDeleteBatchReq, opt, uint32(c.defZone))
}

/**
    @brief list表查询key下所有记录
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *PBClient) DoListGetAll(msg proto.Message, opt *option.PBOpt, zoneId ...uint32) (map[int32]proto.Message,
	error) {
	if len(zoneId) == 1 {
		return c.doListBatch(msg, nil, cmd.TcaplusApiListGetAllReq, opt, zoneId[0])
	}
	return c.doListBatch(msg, nil, cmd.TcaplusApiListGetAllReq, opt, uint32(c.defZone))
}
