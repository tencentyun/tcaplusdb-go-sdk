package request

import (
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/common"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/policy"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type listDeleteBatchRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
}

func newListDeleteBatchRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*listDeleteBatchRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListDeleteBatchReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.ListDeleteBatchReq.ElementNum = 0
	pkg.Body.ListDeleteBatchReq.CheckVersiontType = policy.CheckDataVersionAutoIncrease
	req := &listDeleteBatchRequest{
		appId:     appId,
		zoneId:    zoneId,
		tableName: tableName,
		cmd:       cmd,
		seq:       seq,
		record:    nil,
		pkg:       pkg,
	}
	return req, nil
}

func (req *listDeleteBatchRequest) AddRecord(index int32) (*record.Record, error) {
	if req.record != nil {
		return nil, &terror.ErrorCode{Code: terror.RecordNumOverMax}
	}

	rec := &record.Record{
		AppId:       req.appId,
		ZoneId:      req.zoneId,
		TableName:   req.tableName,
		Cmd:         req.cmd,
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      nil,
		ValueSet:    nil,
		UpdFieldSet: nil,
	}

	rec.KeySet = req.pkg.Head.KeyInfo

	req.record = rec
	return rec, nil
}

func (req *listDeleteBatchRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listDeleteBatchRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.ListDeleteBatchReq.CheckVersiontType = p
	return nil
}

func (req *listDeleteBatchRequest) SetResultFlag(flag int) error {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "ResultFlag invalid"}
	}
	req.pkg.Body.ListDeleteBatchReq.Flag = byte(flag)
	return nil
}

func (req *listDeleteBatchRequest) Pack() ([]byte, error) {
	if req.record == nil {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	if err := req.record.PackKey(); err != nil {
		logger.ERR("record pack key failed, %s", err.Error())
		return nil, err
	}

	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("getRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listDeleteBatchRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listDeleteBatchRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *listDeleteBatchRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *listDeleteBatchRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listDeleteBatchRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listDeleteBatchRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}
func (req *listDeleteBatchRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listDeleteBatchRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	if 1 == multi_flag {
		req.pkg.Body.ListDeleteBatchReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.ListDeleteBatchReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *listDeleteBatchRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *listDeleteBatchRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *listDeleteBatchRequest) AddElementIndex(idx int32) int32 {
	bodyreq := req.pkg.Body.ListDeleteBatchReq
	for i := int32(0); i < bodyreq.ElementNum; i++ {
		if idx == bodyreq.ElementIndexArray[i] {
			return int32(terror.GEN_ERR_SUC)
		}
	}
	if int64(bodyreq.ElementNum) > tcaplus_protocol_cs.TCAPLUS_MAX_LIST_ELEMENTS_NUM {
		return int32(terror.API_ERR_OVER_MAX_LIST_INDEX_NUM)
	}
	bodyreq.ElementIndexArray = append(bodyreq.ElementIndexArray, idx)
	bodyreq.ElementNum++
	return int32(terror.GEN_ERR_SUC)
}
