package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
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
	isPB         bool
}

func newUpdateByPartKeyRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*updateByPartKeyRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.UpdateByPartkeyReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.UpdateByPartkeyReq.ValueInfo.EncodeType = 1
	pkg.Body.UpdateByPartkeyReq.ValueInfo.Version_ = 0
	pkg.Body.UpdateByPartkeyReq.ValueInfo.CompactValueSet.ValueBuf = nil
	pkg.Body.UpdateByPartkeyReq.ValueInfo.CompactValueSet.ValueBufLen = 0
	pkg.Body.UpdateByPartkeyReq.ValueInfo.CompactValueSet.FieldIndexs = nil
	pkg.Body.UpdateByPartkeyReq.ValueInfo.CompactValueSet.FieldIndexNum = 0
	pkg.Body.UpdateByPartkeyReq.ValueInfo.FieldNum_ = 0
	pkg.Body.UpdateByPartkeyReq.ValueInfo.Fields_ = nil
	pkg.Body.UpdateByPartkeyReq.OffSet = 0
	pkg.Body.UpdateByPartkeyReq.Limit = -1
	pkg.Body.UpdateByPartkeyReq.CheckVersiontType = 1

	req := &updateByPartKeyRequest{
		appId:        appId,
		zoneId:       zoneId,
		tableName:    tableName,
		cmd:          cmd,
		seq:          seq,
		record:       nil,
		pkg:          pkg,
		valueNameMap: make(map[string]bool),
		isPB:         isPB,
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
		IsPB:        req.isPB,
	}
	rec.ShardingKey = &req.pkg.Head.SplitTableKeyBuff
	rec.ShardingKeyLen = &req.pkg.Head.SplitTableKeyBuffLen
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
		return &terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.UpdateByPartkeyReq.CheckVersiontType = p
	return nil
}

func (req *updateByPartKeyRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "not Support ResultFlag"}
}

func (req *updateByPartKeyRequest) Pack() ([]byte, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return nil, &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}

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

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	}
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
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return uint32(terror.RequestHasHasNoPkg), &terror.ErrorCode{Code: terror.RequestHasHasNoPkg,
			Message: "Request can not second use"}
	}
	defer func() {
		cs_pool.PutTcaplusCSPkg(req.pkg)
		req.pkg = nil
	}()
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

func (req *updateByPartKeyRequest) SetResultLimit(limit int32, offset int32) int32 {
	req.pkg.Body.UpdateByPartkeyReq.OffSet = offset
	req.pkg.Body.UpdateByPartkeyReq.Limit = limit
	return int32(terror.GEN_ERR_SUC)
}

func (req *updateByPartKeyRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *updateByPartKeyRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *updateByPartKeyRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *updateByPartKeyRequest) SetPerfTest(sendTime uint64) int {
	perf := tcaplus_protocol_cs.NewPerfTest()
	perf.ApiSendTime = sendTime
	perf.Version = tcaplus_protocol_cs.PerfTestCurrentVersion
	p, err := perf.Pack(tcaplus_protocol_cs.PerfTestCurrentVersion)
	if err != nil {
		logger.ERR("pack perf error: %s", err)
		return terror.API_ERR_PARAMETER_INVALID
	}
	req.pkg.Head.PerfTest = p
	req.pkg.Head.PerfTestLen = uint32(len(p))
	return terror.GEN_ERR_SUC
}

func (req *updateByPartKeyRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *updateByPartKeyRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *updateByPartKeyRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
