package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type deleteByPartKeyRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
}

func newDeleteByPartKeyRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*deleteByPartKeyRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.DeleteByPartkeyReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	req := &deleteByPartKeyRequest{
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

func (req *deleteByPartKeyRequest) AddRecord(index int32) (*record.Record, error) {
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

	//key set
	rec.KeySet = req.pkg.Head.KeyInfo
	req.record = rec
	return rec, nil
}

func (req *deleteByPartKeyRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *deleteByPartKeyRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.DeleteByPartkeyReq.CheckVersiontType = p
	return nil
}

func (req *deleteByPartKeyRequest) SetResultFlag(flag int) error {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "DeleteByPartkey not Support ResultFlag not support"}
//	req.pkg.Body.DeleteByPartkeyReq.Flag = byte(flag)
}

func (req *deleteByPartKeyRequest) Pack() ([]byte, error) {
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
		logger.ERR("deleteByPartKeyRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *deleteByPartKeyRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *deleteByPartKeyRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *deleteByPartKeyRequest) SetFieldNames(valueNameList []string) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "DeleteByPartkey not Support SetFieldNames"}
}

func (req *deleteByPartKeyRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *deleteByPartKeyRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *deleteByPartKeyRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *deleteByPartKeyRequest)SetResultLimit(limit int32, offset int32) int32 {
	req.pkg.Body.DeleteByPartkeyReq.OffSet = offset
	req.pkg.Body.DeleteByPartkeyReq.Limit = limit
	return int32(terror.GEN_ERR_SUC)
}

func (req *deleteByPartKeyRequest)SetAddableIncreaseFlag(increase_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *deleteByPartKeyRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *deleteByPartKeyRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *deleteByPartKeyRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
