package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type updateByPartKeyRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       *record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
}

func newUpdateByPartKeyRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*updateByPartKeyRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.UpdateByPartkeyReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.UpdateByPartkeyReq.ValueInfo.EncodeType = 1
	req := &updateByPartKeyRequest{
		appId:        appId,
		zoneId:       zoneId,
		tableName:    tableName,
		cmd:          cmd,
		seq:          seq,
		record:       nil,
		pkg:          pkg,
		valueNameMap: make(map[string]bool),
	}
	return req, nil
}

func (req *updateByPartKeyRequest) AddRecord(index int32) (*record.Record, error) {
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
	rec.ValueSet = req.pkg.Body.UpdateByPartkeyReq.ValueInfo
	req.record = rec
	return rec, nil
}

func (req *updateByPartKeyRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *updateByPartKeyRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.UpdateByPartkeyReq.CheckVersiontType = p
	return nil
}

func (req *updateByPartKeyRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "not Support ResultFlag"}
}

func (req *updateByPartKeyRequest) Pack() ([]byte, error) {
	if req.record == nil {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	if err := req.record.PackKey(); err != nil {
		logger.ERR("record pack key failed, %s", err.Error())
		return nil, err
	}

	//req.pkg.Body.UpdateByPartkeyReq.OffSet = 0
	//req.pkg.Body.UpdateByPartkeyReq.Limit = -1
	//req.pkg.Body.UpdateByPartkeyReq.ValueInfo.FieldNum = 0

	//for key, _ := range req.record.ValueMap {
	//	req.pkg.Body.UpdateByPartkeyReq.ValueInfo.FieldNum += 1
	//	req.pkg.Body.UpdateByPartkeyReq.ValueInfo.FieldName =
	//	append(req.pkg.Body.UpdateByPartkeyReq.ValueInfo.FieldName, key)
	//}

	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("getRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *updateByPartKeyRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *updateByPartKeyRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *updateByPartKeyRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *updateByPartKeyRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *updateByPartKeyRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *updateByPartKeyRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *updateByPartKeyRequest)SetResultLimit(limit int32, offset int32) int32 {
	req.pkg.Body.UpdateByPartkeyReq.OffSet = offset
	req.pkg.Body.UpdateByPartkeyReq.Limit = limit
	return int32(terror.GEN_ERR_SUC)
}

func (req *updateByPartKeyRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *updateByPartKeyRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *updateByPartKeyRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
