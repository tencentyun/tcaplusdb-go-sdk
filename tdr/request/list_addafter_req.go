package request

import (
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/common"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/policy"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type listAddAfterRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
}

func newListAddAfterRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*listAddAfterRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListAddAfterReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	pkg.Body.ListAddAfterReq.ElementValueInfo.EncodeType = 1
	pkg.Body.ListAddAfterReq.ShiftFlag = byte(tcaplus_protocol_cs.TCAPLUS_LIST_SHIFT_HEAD)
	pkg.Body.ListAddAfterReq.ElementIndex = -1
	pkg.Body.ListAddAfterReq.CheckVersiontType = policy.CheckDataVersionAutoIncrease
	req := &listAddAfterRequest{
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

func (req *listAddAfterRequest) AddRecord(index int32) (*record.Record, error) {
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

	//key value set
	rec.KeySet = req.pkg.Head.KeyInfo
	rec.ValueSet = req.pkg.Body.ListAddAfterReq.ElementValueInfo
	req.pkg.Body.ListAddAfterReq.ElementIndex = index
	req.record = rec
	return rec, nil
}

func (req *listAddAfterRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listAddAfterRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.ListAddAfterReq.CheckVersiontType = p
	return nil
}

func (req *listAddAfterRequest) SetResultFlag(flag int) error {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "ResultFlag invalid"}
	}
	req.pkg.Body.ListAddAfterReq.Flag = byte(flag)
	return nil
}

func (req *listAddAfterRequest) Pack() ([]byte, error) {
	if req.record == nil {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	if err := req.record.PackKey(); err != nil {
		logger.ERR("record pack key failed, %s", err.Error())
		return nil, err
	}

	if err := req.record.PackValue(nil); err != nil {
		logger.ERR("record pack value failed, %s", err.Error())
		return nil, err
	}

	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("listAddAfterRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listAddAfterRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listAddAfterRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *listAddAfterRequest) SetFieldNames(valueNameList []string) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "list insert not Support SetFieldNames"}
}

func (req *listAddAfterRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listAddAfterRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listAddAfterRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *listAddAfterRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listAddAfterRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listAddAfterRequest) SetListShiftFlag(shiftFlag byte) int32 {
	req.pkg.Body.ListAddAfterReq.ShiftFlag = shiftFlag
	return int32(terror.GEN_ERR_SUC)
}

func (req *listAddAfterRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *listAddAfterRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
