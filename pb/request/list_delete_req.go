package request

import (
	"git.code.com/gcloud_storage_group/tcaplus-go-api/common"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/policy"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type listDeleteRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
}

func newListDeleteRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*listDeleteRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListDeleteReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	req := &listDeleteRequest{
		appId:        appId,
		zoneId:       zoneId,
		tableName:    tableName,
		cmd:          cmd,
		seq:          seq,
		record:       nil,
		pkg:          pkg,
	}
	return req, nil
}

func (req *listDeleteRequest) AddRecord(index int32) (*record.Record, error) {
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
	req.pkg.Body.ListDeleteReq.ElementIndex = index
	req.record = rec
	return rec, nil
}

func (req *listDeleteRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listDeleteRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.ListDeleteReq.CheckVersiontType = p
	return nil
}

func (req *listDeleteRequest) SetResultFlag(flag int) error {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "ResultFlag invalid"}
	}
	req.pkg.Body.ListDeleteReq.Flag = byte(flag)
	return nil
}

func (req *listDeleteRequest) Pack() ([]byte, error) {
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
		logger.ERR("listDeleteRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listDeleteRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listDeleteRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *listDeleteRequest) SetFieldNames(valueNameList []string) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "list delete not Support SetFieldNames"}
}

func (req *listDeleteRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listDeleteRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listDeleteRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *listDeleteRequest)SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listDeleteRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listDeleteRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *listDeleteRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
