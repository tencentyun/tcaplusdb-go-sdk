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

	if opt != nil && len(opt.Condition) > 0 {
		if ret := rec.SetCondition(opt.Condition); ret != 0 {
			logger.ERR("SetCondition error:%d", ret)
			return index, &terror.ErrorCode{Code: ret, Message: "SetCondition maybe not support or len too long"}
		}
	}
	if opt != nil && len(opt.Operation) > 0 {
		if ret := rec.SetOperation(opt.Operation, 0); ret != 0 {
			logger.ERR("SetOperation error:%d", ret)
			return index, &terror.ErrorCode{Code: ret, Message: "SetOperation failed, maybe not support or len too long"}
		}
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
		err = &terror.ErrorCode{Code: ret}
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
	return index, err
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

	if len(opt.Condition) > 0 {
		if ret := rec.SetCondition(opt.Condition); ret != 0 {
			logger.ERR("SetCondition error:%d", ret)
			return &terror.ErrorCode{Code: ret, Message: "SetCondition failed, maybe not support or len too long"}
		}
	}
	if len(opt.Operation) > 0 {
		if ret := rec.SetOperation(opt.Operation, 0); ret != 0 {
			logger.ERR("SetOperation error:%d", ret)
			return &terror.ErrorCode{Code: ret, Message: "SetOperation failed, maybe not support or len too long"}
		}
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
		err = &terror.ErrorCode{Code: ret}
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
	return err
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

	if opt.ExpireTime != 0 {
		req.SetExpireTime(opt.ExpireTime)
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

	if opt != nil && len(opt.Condition) > 0 {
		if ret := rec.SetCondition(opt.Condition); ret != 0 {
			logger.ERR("SetCondition error:%d", ret)
			return nil, &terror.ErrorCode{Code: ret, Message: "SetCondition failed, maybe not support or len too long"}
		}
	}
	if opt != nil && len(opt.Operation) > 0 {
		if ret := rec.SetOperation(opt.Operation, 0); ret != 0 {
			logger.ERR("SetOperation error:%d", ret)
			return nil, &terror.ErrorCode{Code: ret, Message: "SetOperation failed, maybe not support or len too long"}
		}
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
		if _, exist := tmpIdxMap[index]; exist {
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

		if opt != nil && len(opt.Condition) > 0 {
			if ret := rec.SetCondition(opt.Condition); ret != 0 {
				logger.ERR("SetCondition error:%d", ret)
				return &terror.ErrorCode{Code: ret, Message: "SetCondition failed, maybe not support or len too long"}
			}
		}
		if opt != nil && len(opt.Operation) > 0 {
			if ret := rec.SetOperation(opt.Operation, 0); ret != 0 {
				logger.ERR("SetOperation error:%d", ret)
				return &terror.ErrorCode{Code: ret, Message: "SetOperation failed, maybe not support or len too long"}
			}
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

		if opt != nil && len(opt.Condition) > 0 {
			if ret := rec.SetCondition(opt.Condition); ret != 0 {
				logger.ERR("SetCondition error:%d", ret)
				return &terror.ErrorCode{Code: ret, Message: "SetCondition failed, maybe not support or len too long"}
			}
		}
		if opt != nil && len(opt.Operation) > 0 {
			if ret := rec.SetOperation(opt.Operation, 0); ret != 0 {
				logger.ERR("SetOperation error:%d", ret)
				return &terror.ErrorCode{Code: ret, Message: "SetOperation failed, maybe not support or len too long"}
			}
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

func (c *Client) doPartKey(table string, data record.TdrTableSt, indexName string, apiCmd int, opt *option.TDROpt,
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

	if opt != nil && len(opt.Condition) > 0 {
		if ret := rec.SetCondition(opt.Condition); ret != 0 {
			logger.ERR("SetCondition error:%d", ret)
			return nil, &terror.ErrorCode{Code: ret, Message: "SetCondition failed, maybe not support or len too long"}
		}
	}
	if opt != nil && len(opt.Operation) > 0 {
		if ret := rec.SetOperation(opt.Operation, 0); ret != 0 {
			logger.ERR("SetOperation error:%d", ret)
			return nil, &terror.ErrorCode{Code: ret, Message: "SetOperation failed, maybe not support or len too long"}
		}
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

	if len(recordSlice) == 0 && globalErr == nil && apiCmd == cmd.TcaplusApiGetByPartkeyReq {
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
		return c.doPartKey(table, data, indexName, cmd.TcaplusApiGetByPartkeyReq, opt, zoneId[0])
	}
	return c.doPartKey(table, data, indexName, cmd.TcaplusApiGetByPartkeyReq, opt, uint32(c.defZone))
}

/**
    @brief 根据表的部分key字段删除记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoDeleteByPartKey(table string, data record.TdrTableSt, indexName string, opt *option.TDROpt,
	zoneId ...uint32) ([]*record.Record, error) {
	if len(indexName) == 0 {
		return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "indexName is empty"}
	}
	if len(zoneId) == 1 {
		return c.doPartKey(table, data, indexName, cmd.TcaplusApiDeleteByPartkeyReq, opt, zoneId[0])
	}
	return c.doPartKey(table, data, indexName, cmd.TcaplusApiDeleteByPartkeyReq, opt, uint32(c.defZone))
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
    @retval int32 如果有返回记录，返回index 索引
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

/**
@brief  设置记录的生存时间，或者说过期时间，即记录多久之后过期，过期的记录将不会被访问到
@param [IN]  table 表名
@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
@param [IN/OUT] opt 必填，请设置option.BatchTTL
@param [IN] option.BatchTTL.ttl 生存时间（过期时间），时间单位为毫秒，如果是相对时间，比如该参数值为10，则表示记录写入10ms之后过期，该参数值为0，则表示记录永不过期
												   如果是绝对时间，比如该参数值为1599105600000, 则表示记录到"20200903 12:00:00"之后过期，该参数值为0，则表示记录永不过期
@param [IN] option.BatchTTL.IsAbsolute 时间类型是否为绝对时间，true表示绝对时间，false表示相对时间，默认是false，即相对时间
@note   设置的ttl值最大不能超过uint64_t最大值的一半，即ttl最大值为 ULONG_MAX/2，超过该值接口会强制设置为该值
@note   设置ttl的请求，在服务端不会增加对应记录的版本号
@note   对于list表，当某个key下面所有记录因为过期删除后，会直接将索引记录也删除
@note   对于设置了ttl的记录，如果是getbypartkey查询，并且只需要返回key字段（即不需要返回value字段）时，此时不会检查该记录是否过期
@note   对于删除操作(generic表和list表的删除)，均不会检验记录是否过期
@notice @param [IN]   indexs请设置为NULL，list表目前不支持ttl
*/
func (c *Client) DoSetTTLBatch(table string, dataSlice []record.TdrTableSt, indexs []int32, opt *option.TDROpt,
	zoneId ...uint32) error {
	if opt == nil {
		logger.ERR("opt is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt is nil"}
	}

	if len(dataSlice) == 0 {
		logger.ERR("dataSlice len is 0")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dataSlice len is 0"}
	}

	if len(opt.BatchTTL) == 0 {
		logger.ERR("opt.BatchTTL is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt.BatchTTL is nil"}
	}

	if len(opt.BatchTTL) != len(dataSlice) {
		logger.ERR("len dataSlice != opt.BatchTTL")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "len dataSlice != opt.BatchTTL"}
	}
	if indexs != nil && len(opt.BatchTTL) != len(indexs) {
		logger.ERR("len indexs != opt.BatchTTL")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "len indexs != opt.BatchTTL"}
	}

	zone := uint32(0)
	if len(zoneId) == 1 {
		zone = zoneId[0]
	} else {
		zone = uint32(c.defZone)
	}

	req, err := c.NewRequest(zone, table, cmd.TcaplusApiSetTtlReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	if opt.BatchResult == nil {
		opt.BatchResult = make([]error, len(opt.BatchTTL), len(opt.BatchTTL))
	} else if len(opt.BatchResult) != len(opt.BatchTTL) {
		logger.ERR("dataSlice and BatchResult count not equal")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dataSlice and opt.BatchResult count not equal"}
	}

	msgMap := make(map[string]int, len(opt.BatchTTL))
	for i, ttlRec := range opt.BatchTTL {
		index := int32(0)
		if indexs != nil {
			index = indexs[i]
		}
		rec, err := req.AddRecord(index)
		if err != nil {
			logger.ERR("AddRecord error:%s", err)
			return err
		}

		err = rec.SetData(dataSlice[i])
		if err != nil {
			logger.ERR("SetData error:%s", err)
			return err
		}
		ret := rec.SetTTL(ttlRec.TTL, ttlRec.IsAbsolute)
		if ret != 0 {
			logger.ERR("SetTTL error:%d", ret)
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "SetTTL error"}
		}

		key, err := rec.GetAllKeyBlob()
		if err != nil {
			logger.ERR("GetAllKeyBlob error:%s", err)
			return err
		}

		msgMap[key] = i
	}
	var resps []response.TcaplusResponse
	if opt.Timeout > 0 {
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

			opt.BatchResult[index] = recErr
			delete(msgMap, key)
			if recErr != nil {
				continue
			}
		}

		//msgMap not nil
		for key, index := range msgMap {
			logger.ERR("key %s not rsp", key)
			opt.BatchResult[index] = &terror.ErrorCode{Code: terror.API_ERR_WAIT_RSP_TIMEOUT,
				Message: "no rsp with this key"}
		}
	}
	return globalErr
}

/**
@brief  获取记录的剩余生存时间，或者说剩余过期时间，即记录多久后过期
@param [IN]  table 表名
@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
@param [IN/OUT] opt 必填，返回的ttl在option.BatchTTL中
@note   该函数当前只支持 TCAPLUS_API_GET_TTL_RES 响应
@notice @param [IN]   indexs请设置为NULL，list表目前不支持ttl
*/
func (c *Client) DoGetTTLBatch(table string, dataSlice []record.TdrTableSt, indexs []int32, opt *option.TDROpt,
	zoneId ...uint32) error {
	if opt == nil {
		logger.ERR("opt is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "opt is nil"}
	}

	if len(dataSlice) == 0 {
		logger.ERR("dataSlice len is 0")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dataSlice len is 0"}
	}
	if opt.BatchTTL == nil {
		opt.BatchTTL = make([]option.TTLInfo, len(dataSlice), len(dataSlice))
	} else {
		if len(dataSlice) != len(opt.BatchTTL) {
			logger.ERR("len dataSlice != opt.BatchTTL")
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "len dataSlice != opt.BatchTTL"}
		}
	}

	if indexs != nil && len(dataSlice) != len(indexs) {
		logger.ERR("len dataSlice != opt.BatchTTL")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "len dataSlice != opt.BatchTTL"}
	}

	zone := uint32(0)
	if len(zoneId) == 1 {
		zone = zoneId[0]
	} else {
		zone = uint32(c.defZone)
	}

	req, err := c.NewRequest(zone, table, cmd.TcaplusApiGetTtlReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return err
	}

	if opt.BatchResult == nil {
		opt.BatchResult = make([]error, len(opt.BatchTTL), len(opt.BatchTTL))
	} else if len(opt.BatchResult) != len(opt.BatchTTL) {
		logger.ERR("dataSlice and BatchResult count not equal")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dataSlice and opt.BatchResult count not equal"}
	}

	msgMap := make(map[string]int, len(opt.BatchTTL))
	for i, _ := range opt.BatchTTL {
		index := int32(0)
		if indexs != nil {
			index = indexs[i]
		}
		rec, err := req.AddRecord(index)
		if err != nil {
			logger.ERR("AddRecord error:%s", err)
			return err
		}

		err = rec.SetData(dataSlice[i])
		if err != nil {
			logger.ERR("SetData error:%s", err)
			return err
		}

		key, err := rec.GetAllKeyBlob()
		if err != nil {
			logger.ERR("GetAllKeyBlob error:%s", err)
			return err
		}

		msgMap[key] = i
	}
	var resps []response.TcaplusResponse
	if opt.Timeout > 0 {
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
			delete(msgMap, key)
			opt.BatchResult[index] = recErr
			if recErr != nil {
				continue
			}

			if ret := resRec.GetTTL(&opt.BatchTTL[index].TTL); ret != 0 {
				globalErr = &terror.ErrorCode{Code: ret, Message: "rec GetTTL failed"}
				logger.ERR("resRec.GetTTL failed")
				opt.BatchResult[index] = globalErr
			}
		}

		//msgMap not nil
		for key, index := range msgMap {
			logger.ERR("key %s not rsp", key)
			opt.BatchResult[index] = &terror.ErrorCode{Code: terror.API_ERR_WAIT_RSP_TIMEOUT,
				Message: "no rsp with this key"}
		}
	}
	return globalErr
}

/**
    @brief 同一个表的批量插入
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数, 分包返回等, 记录的version存放在opt.BatchVersion，单条记录结果存放opt.BatchResult
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoBatchInsert(table string, dataSlice []record.TdrTableSt, opt *option.TDROpt,
	zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchInsertReq, opt, zoneId[0])
	}
	return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchInsertReq, opt, uint32(c.defZone))
}

/**
    @brief 同一个表的批量更新,存在更新，不存在则插入
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数, 分包返回等, 记录的version存放在opt.BatchVersion，单条记录结果存放opt.BatchResult
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoBatchReplace(table string, dataSlice []record.TdrTableSt, opt *option.TDROpt,
	zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchReplaceReq, opt, zoneId[0])
	}
	return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchReplaceReq, opt, uint32(c.defZone))
}

/**
    @brief 同一个表的批量更新，不存在则报错
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数, 分包返回等, 记录的version存放在opt.BatchVersion，单条记录结果存放opt.BatchResult
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoBatchUpdate(table string, dataSlice []record.TdrTableSt, opt *option.TDROpt,
	zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchUpdateReq, opt, zoneId[0])
	}
	return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchUpdateReq, opt, uint32(c.defZone))
}

/**
    @brief 同一个表的批量删除
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN/OUT] opt 可选参数, 分包返回等, 记录的version存放在opt.BatchVersion，单条记录结果存放opt.BatchResult
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoBatchDelete(table string, dataSlice []record.TdrTableSt, opt *option.TDROpt,
	zoneId ...uint32) error {
	if len(zoneId) == 1 {
		return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchDeleteReq, opt, zoneId[0])
	}
	return c.doBatch(table, dataSlice, cmd.TcaplusApiBatchDeleteReq, opt, uint32(c.defZone))
}

/**
    @brief list表下同一key获取批量记录
	@param [IN]  table 表名
	@param [IN/OUT] data  tdr结构体由TdrCodeGen生成的记录结构体, 若有记录返回会更新为返回的记录
	@param [IN] indexs []int32 key下多个索引
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoListGetBatch(table string, data record.TdrTableSt, indexs []int32, opt *option.TDROpt,
	zoneId ...uint32) ([]*record.Record, error) {
	if len(indexs) == 0 {
		return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "indexs is empty"}
	}
	if len(zoneId) == 1 {
		return c.doListBatch(table, data, indexs, cmd.TcaplusApiListGetBatchReq, opt, zoneId[0])
	}
	return c.doListBatch(table, data, indexs, cmd.TcaplusApiListGetBatchReq, opt, uint32(c.defZone))
}

/**
    @brief list表下同一key插入批量记录，必须是相同key的批量操作
	@param [IN]  table 表名
	@param [IN/OUT] dataSlice  tdr结构体由TdrCodeGen生成的记录结构体, 该接口必须保证key相同
	@param [IN/OUT] indexs dataSlice对应记录的index，传入可设置为以下参数
				tcaplus_protocol_cs.TCAPLUS_LIST_LAST_INDEX = -1      插入元素位置在最后面
				tcaplus_protocol_cs.TCAPLUS_LIST_PRE_FIRST_INDEX = -2 插入元素位置在最前面
				输出会返回对应的索引
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoListAddAfterBatch(table string, dataSlice []record.TdrTableSt, indexs []int32, opt *option.TDROpt,
	zoneId ...uint32) error {
	if len(indexs) == 0 || len(dataSlice) == 0 {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "indexs or dataSlice is empty"}
	}
	if len(zoneId) == 1 {
		return c.doListBatchRecord(table, dataSlice, indexs, cmd.TcaplusApiListAddAfterBatchReq, opt, zoneId[0])
	}
	return c.doListBatchRecord(table, dataSlice, indexs, cmd.TcaplusApiListAddAfterBatchReq, opt, uint32(c.defZone))
}

/**
    @brief list表下同一key替换批量记录，必须是相同key的批量操作
	@param [IN]  table 表名
	@param [IN/OUT] dataSlice  tdr结构体由TdrCodeGen生成的记录结构体, 该接口必须保证key相同
	@param [IN/OUT] indexs dataSlice对应记录的index
	@param [IN/OUT] opt 可选参数，乐观锁，flag等，若有记录返回，会更新opt中的version为记录的version
	@param [IN] zoneId 可选参数，不设置则取默认zone，默认zone可通过client.SetDefaultZoneId设置
    @retval error 错误码
**/
func (c *Client) DoListReplaceBatch(table string, dataSlice []record.TdrTableSt, indexs []int32, opt *option.TDROpt,
	zoneId ...uint32) error {
	if len(indexs) == 0 || len(dataSlice) == 0 {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "indexs or dataSlice is empty"}
	}
	if len(zoneId) == 1 {
		return c.doListBatchRecord(table, dataSlice, indexs, cmd.TcaplusApiListReplaceBatchReq, opt, zoneId[0])
	}
	return c.doListBatchRecord(table, dataSlice, indexs, cmd.TcaplusApiListReplaceBatchReq, opt, uint32(c.defZone))
}
