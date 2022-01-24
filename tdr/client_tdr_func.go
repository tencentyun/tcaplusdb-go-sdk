package tcaplus

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

/*Client接口的简单封装，方便用户编码*/
func (c *Client) DoListSimple(table string, data record.TdrTableSt, index int32, apiCmd int, opt *option.TDROpt,
	zoneId uint32) (int32, error) {
	req, err := c.NewRequest(zoneId, table, apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return index, err
	}

	rec, err := req.AddRecord(index)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return index, err
	}

	err = rec.SetData(data)
	if err != nil {
		logger.ERR("SetData error:%s", err)
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
		resRec, err := res.FetchRecord()
		if err != nil {
			logger.ERR("FetchRecord error:%s", err)
			return index, err
		}

		if opt != nil {
			opt.Version = resRec.GetVersion()
		}
		index = resRec.GetIndex()

		if apiCmd != cmd.TcaplusApiListGetReq && !c.needGetData(opt) {
			return index, nil
		}

		err = resRec.GetData(data)
		if err != nil {
			logger.ERR("GetData error:%s", err)
			return index, err
		}
	}
	return index, nil
}

func (c *Client) needGetData(opt *option.TDROpt) bool {
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

func (c *Client) setRecOpt(rec *record.Record, opt *option.TDROpt) error {
	if opt.Version > 0 {
		rec.SetVersion(opt.Version)
	}

	if len(opt.IncField) > 0 {
		for _, incField := range opt.IncField {
			if err := rec.SetIncValue(incField.FieldName, incField.IncData, incField.Operation, incField.LowerLimit,
				incField.UpperLimit); err != nil {
				logger.ERR("IncValue error:%s", err.Error())
				return err
			}
		}
	}
	return nil
}

func (c *Client) doSimple(table string, data record.TdrTableSt, apiCmd int, opt *option.TDROpt, zoneId uint32) error {
	req, err := c.NewRequest(zoneId, table, apiCmd)
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

	err = rec.SetData(data)
	if err != nil {
		logger.ERR("SetData error:%s", err)
		return err
	}

	if opt != nil {
		err = c.setRecOpt(rec, opt)
		if err != nil {
			logger.ERR("setReqOpt error:%s", err)
			return err
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
		return err
	}

	ret := res.GetResult()
	if ret != 0 {
		return &terror.ErrorCode{Code: ret}
	}

	if res.GetRecordCount() > 0 {
		resRec, err := res.FetchRecord()
		if err != nil {
			logger.ERR("FetchRecord error:%s", err)
			return err
		}

		if opt != nil {
			opt.Version = resRec.GetVersion()
		}

		if apiCmd != cmd.TcaplusApiGetReq && !c.needGetData(opt) {
			return nil
		}

		err = resRec.GetData(data)
		if err != nil {
			logger.ERR("GetData error:%s", err)
			return err
		}
	}
	return nil
}

func (c *Client) setReqOpt(req request.TcaplusRequest, opt *option.TDROpt) error {
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

	if opt.AddableIncreaseFlag != 0 {
		req.SetAddableIncreaseFlag(opt.AddableIncreaseFlag)
	}

	if len(opt.FieldNames) > 0 {
		req.SetFieldNames(opt.FieldNames)
	}
	return nil
}
func (c *Client) setBatchReqOpt(req request.TcaplusRequest, opt *option.TDROpt) error {
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

	if len(opt.FieldNames) > 0 {
		req.SetFieldNames(opt.FieldNames)
	}

	return nil
}
func (c *Client) doListBatch(table string, data record.TdrTableSt, indexs []int32,
	apiCmd int, opt *option.TDROpt, zoneId uint32) ([]*record.Record, error) {
	req, err := c.NewRequest(zoneId, table, apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return nil, err
	}

	rec, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return nil, err
	}

	err = rec.SetData(data)
	if err != nil {
		logger.ERR("SetData error:%s", err)
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

	var recordSlice []*record.Record
	offset := 0
	var globalErr error
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

			recordSlice = append(recordSlice, resRec)
			if opt != nil {
				opt.Version = resRec.GetVersion()
			}
			offset++
		}
	}

	return recordSlice, globalErr
}

func (c *Client) doListBatchRecord(table string, dataSlice []record.TdrTableSt, indexs []int32,
	apiCmd int, opt *option.TDROpt, zoneId uint32) error {
	if len(dataSlice) != len(indexs) {
		logger.ERR("len(dataSlice) %d != len(indexs) %d", len(dataSlice), len(indexs))
		return terror.ErrorCode{Code: terror.ParameterInvalid, Message: "len(dataSlice) != len(indexs)"}
	}

	req, err := c.NewRequest(zoneId, table, apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	if opt != nil {
		if opt.BatchResult == nil {
			opt.BatchResult = make([]error, len(dataSlice), len(dataSlice))
		} else if len(opt.BatchResult) != len(dataSlice) {
			logger.ERR("dataSlice and BatchResult count not equal")
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dataSlice and opt.BatchResult count not equal"}
		}

		err = c.setBatchReqOpt(req, opt)
		if err != nil {
			logger.ERR("setReqOpt error:%s", err)
			return err
		}
	}

	for i, data := range dataSlice {
		rec, err := req.AddRecord(indexs[i])
		if err != nil {
			logger.ERR("AddRecord error:%s", err)
			return err
		}

		err = rec.SetData(data)
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
					err = resRec.GetData(dataSlice[offset])
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

func (c *Client) doBatch(table string, dataSlice []record.TdrTableSt, apiCmd int, opt *option.TDROpt,
	zoneId uint32) error {
	if len(dataSlice) == 0 {
		logger.ERR("dataSlice is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dataSlice is nil"}
	}

	req, err := c.NewRequest(zoneId, table, apiCmd)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	if opt != nil {
		if opt.BatchVersion == nil {
			opt.BatchVersion = make([]int32, len(dataSlice), len(dataSlice))
		} else if len(opt.BatchVersion) != len(dataSlice) {
			logger.ERR("dataSlice and BatchVersion count not equal")
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dataSlice and opt.BatchVersion count not equal"}
		}

		if opt.BatchResult == nil {
			opt.BatchResult = make([]error, len(dataSlice), len(dataSlice))
		} else if len(opt.BatchResult) != len(dataSlice) {
			logger.ERR("dataSlice and BatchResult count not equal")
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dataSlice and opt.BatchResult count not equal"}
		}

		err = c.setBatchReqOpt(req, opt)
		if err != nil {
			logger.ERR("setReqOpt error:%s", err)
			return err
		}
	}

	msgMap := make(map[string]int, len(dataSlice))
	for i, data := range dataSlice {
		rec, err := req.AddRecord(0)
		if err != nil {
			logger.ERR("AddRecord error:%s", err)
			return err
		}

		err = rec.SetData(data)
		if err != nil {
			logger.ERR("SetData error:%s", err)
			return err
		}

		key, err := rec.GetAllKeyBlob()
		if err != nil {
			logger.ERR("GetAllKeyBlob error:%s", err)
			return err
		}

		if opt != nil && (opt.BatchVersion)[i] > 0 {
			rec.SetVersion((opt.BatchVersion)[i])
		}

		if _, exist := msgMap[key]; exist {
			logger.ERR("batch record exist duplicate key")
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "batch record exist duplicate key"}
		}

		msgMap[key] = i
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
			resRec, recErr := res.FetchRecord()
			if recErr != nil {
				globalErr = recErr
				logger.DEBUG("FetchRecord error:%s", recErr)
				continue
			}

			if resRec == nil {
				continue
			}

			key, err := resRec.GetAllKeyBlob()
			if err != nil {
				globalErr = err
				logger.ERR("GetAllKeyBlob error:%s", err)
				continue
			}

			index, exist := msgMap[key]
			if !exist {
				globalErr = &terror.ErrorCode{Code: terror.RespNotMatchReq}
				logger.ERR("response message is diff with request")
				continue
			}
			if opt != nil {
				opt.BatchResult[index] = recErr
			}
			delete(msgMap, key)
			if recErr != nil {
				continue
			}

			if opt != nil {
				opt.BatchVersion[index] = resRec.GetVersion()
			}

			if apiCmd != cmd.TcaplusApiBatchGetReq && !c.needGetData(opt) {
				continue
			}

			err = resRec.GetData(dataSlice[index])
			if err != nil {
				globalErr = err
				logger.ERR("GetData key %s error:%s", key, err)
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

func (c *Client) doPartKeyGet(table string, data record.TdrTableSt, indexName string, apiCmd int, opt *option.TDROpt,
	zoneId uint32) ([]*record.Record, error) {
	req, err := c.NewRequest(zoneId, table, apiCmd)
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

	if opt != nil && opt.FieldNames != nil {
		err = rec.SetDataWithIndexAndField(data, opt.FieldNames, indexName)
	} else {
		err = rec.SetDataWithIndexAndField(data, nil, indexName)
	}
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

	var recordSlice []*record.Record
	var globalErr error
	for _, res := range resps {
		ret := res.GetResult()
		if ret != 0 {
			globalErr = &terror.ErrorCode{Code: ret}
			logger.ERR("result is %d, error:%s", ret, terror.GetErrMsg(ret))
			continue
		}

		for i := 0; i < res.GetRecordCount(); i++ {
			resRec, err := res.FetchRecord()
			if err != nil {
				globalErr = err
				logger.ERR("FetchRecord error:%s", err)
				continue
			}

			recordSlice = append(recordSlice, resRec)
			if opt != nil {
				opt.BatchVersion = append(opt.BatchVersion, resRec.GetVersion())
			}
		}
	}

	if len(recordSlice) == 0 && globalErr == nil {
		return nil, &terror.ErrorCode{Code: terror.TXHDB_ERR_RECORD_NOT_EXIST}
	}

	return recordSlice, globalErr
}

/**
    @brief 插入记录,记录存在时报错
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoInsert(table string, data record.TdrTableSt, opt *option.TDROpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(table, data, cmd.TcaplusApiInsertReq, opt, zoneId[0])
	}
	return c.doSimple(table, data, cmd.TcaplusApiInsertReq, opt, uint32(c.defZone))
}

/**
    @brief 查询记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoGet(table string, data record.TdrTableSt, opt *option.TDROpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(table, data, cmd.TcaplusApiGetReq, opt, zoneId[0])
	}
	return c.doSimple(table, data, cmd.TcaplusApiGetReq, opt, uint32(c.defZone))
}

/**
    @brief 替换记录，记录不存在时插入，存在时更新
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoReplace(table string, data record.TdrTableSt, opt *option.TDROpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(table, data, cmd.TcaplusApiReplaceReq, opt, zoneId[0])
	}
	return c.doSimple(table, data, cmd.TcaplusApiReplaceReq, opt, uint32(c.defZone))
}

/**
    @brief 更新记录，记录不存在时返错
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoUpdate(table string, data record.TdrTableSt, opt *option.TDROpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(table, data, cmd.TcaplusApiUpdateReq, opt, zoneId[0])
	}
	return c.doSimple(table, data, cmd.TcaplusApiUpdateReq, opt, uint32(c.defZone))
}

/**
    @brief 删除记录，记录不存在时返错
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoDelete(table string, data record.TdrTableSt, opt *option.TDROpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doSimple(table, data, cmd.TcaplusApiDeleteReq, opt, zoneId[0])
	}
	return c.doSimple(table, data, cmd.TcaplusApiDeleteReq, opt, uint32(c.defZone))
}

/**
    @brief 部分value自增自减
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoIncrease(table string, data record.TdrTableSt, opt *option.TDROpt, zoneId ...uint32) error {
	if opt == nil {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt is nil"}
	}

	if len(opt.IncField) == 0 {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt.FieldNames is empty"}
	}

	if len(zoneId) == 1 {
		return c.doSimple(table, data, cmd.TcaplusApiIncreaseReq, opt, zoneId[0])
	}
	return c.doSimple(table, data, cmd.TcaplusApiIncreaseReq, opt, uint32(c.defZone))
}

/**
    @brief 同一个表的批量查询
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数, 分包返回等, 记录的version存放在opt.BatchVersion，单条记录结果存放opt.BatchResult
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoBatchGet(table string, dataSlice []record.TdrTableSt, opt *option.TDROpt,
	zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchGetReq, opt, zoneId[0])
	}
	return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchGetReq, opt, uint32(c.defZone))
}

/**
    @brief 根据表的部分key字段查询,
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数, 分包返回等, 记录的version存放在opt.BatchVersion
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoGetByPartKey(table string, data record.TdrTableSt, indexName string, opt *option.TDROpt,
	zoneId ...uint32) ([]*record.Record, error) {
	if len(indexName) == 0 {
		return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "indexName is empty"}
	}
	if len(zoneId) == 1 {
		return c.doPartKeyGet(table, data, indexName, cmd.TcaplusApiGetByPartkeyReq, opt, zoneId[0])
	}
	return c.doPartKeyGet(table, data, indexName, cmd.TcaplusApiGetByPartkeyReq, opt, uint32(c.defZone))
}

/**
    @brief list表插入记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] index int32 插入到key中的第index条记录之后
				tcaplus_protocol_cs.TCAPLUS_LIST_LAST_INDEX = -1      插入元素位置在最后面
				tcaplus_protocol_cs.TCAPLUS_LIST_PRE_FIRST_INDEX = -2 插入元素位置在最前面
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval int32 如果有返回记录，返回index索引
    @retval error 错误码
**/
func (c *Client) DoListAddAfter(table string, data record.TdrTableSt, index int32, opt *option.TDROpt,
	zoneId ...uint32) (int32, error) {
	if len(zoneId) == 1 {
		return c.DoListSimple(table, data, index, cmd.TcaplusApiListAddAfterReq, opt, zoneId[0])
	}
	return c.DoListSimple(table, data, index, cmd.TcaplusApiListAddAfterReq, opt, uint32(c.defZone))
}

/**
    @brief list表删除记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] index int32 操作第index条记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoListDelete(table string, data record.TdrTableSt, index int32, opt *option.TDROpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		_, err := c.DoListSimple(table, data, index, cmd.TcaplusApiListDeleteReq, opt, zoneId[0])
		return err
	}
	_, err := c.DoListSimple(table, data, index, cmd.TcaplusApiListDeleteReq, opt, uint32(c.defZone))
	return err
}

/**
    @brief list表更新记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] index int32 操作第index条记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoListReplace(table string, data record.TdrTableSt, index int32, opt *option.TDROpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		_, err := c.DoListSimple(table, data, index, cmd.TcaplusApiListReplaceReq, opt, zoneId[0])
		return err
	}
	_, err := c.DoListSimple(table, data, index, cmd.TcaplusApiListReplaceReq, opt, uint32(c.defZone))
	return err
}

/**
    @brief list表查询记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] index int32 操作第index条记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoListGet(table string, data record.TdrTableSt, index int32, opt *option.TDROpt, zoneId ...uint32) error {
	if len(zoneId) == 1 {
		_, err := c.DoListSimple(table, data, index, cmd.TcaplusApiListGetReq, opt, zoneId[0])
		return err
	}
	_, err := c.DoListSimple(table, data, index, cmd.TcaplusApiListGetReq, opt, uint32(c.defZone))
	return err
}

/**
    @brief list表删除key下所有记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁等
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoListDeleteAll(table string, data record.TdrTableSt, opt *option.TDROpt, zoneId ...uint32) error {
	var err error
	if len(zoneId) == 1 {
		_, err = c.doListBatch(table, data, nil, cmd.TcaplusApiListDeleteAllReq, opt, zoneId[0])
	} else {
		_, err = c.doListBatch(table, data, nil, cmd.TcaplusApiListDeleteAllReq, opt, uint32(c.defZone))
	}
	return err
}

/**
    @brief 删除list表下同一key的批量记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] indexs []int32 删除key下多个记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoListDeleteBatch(table string, data record.TdrTableSt, indexs []int32, opt *option.TDROpt,
	zoneId ...uint32) ([]*record.Record, error) {
	if len(indexs) == 0 {
		return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "indexs is empty"}
	}
	if len(zoneId) == 1 {
		return c.doListBatch(table, data, indexs, cmd.TcaplusApiListDeleteBatchReq, opt, zoneId[0])
	}
	return c.doListBatch(table, data, indexs, cmd.TcaplusApiListDeleteBatchReq, opt, uint32(c.defZone))
}

/**
    @brief list表查询key下所有记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoListGetAll(table string, data record.TdrTableSt, opt *option.TDROpt, zoneId ...uint32) ([]*record.Record,
	error) {
	if len(zoneId) == 1 {
		return c.doListBatch(table, data, nil, cmd.TcaplusApiListGetAllReq, opt, zoneId[0])
	}
	return c.doListBatch(table, data, nil, cmd.TcaplusApiListGetAllReq, opt, uint32(c.defZone))
}
