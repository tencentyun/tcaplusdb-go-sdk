package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type pbFieldUpdateRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	isPB      bool
}

func newPBFieldUpdateRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*pbFieldUpdateRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusPbFieldUpdateReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.TCaplusPbFieldUpdateReq.ValueInfo.EncodeType = 1
	pkg.Body.TCaplusPbFieldUpdateReq.ValueInfo.Version_ = 0
	pkg.Body.TCaplusPbFieldUpdateReq.ValueInfo.CompactValueSet.ValueBuf = nil
	pkg.Body.TCaplusPbFieldUpdateReq.ValueInfo.CompactValueSet.ValueBufLen = 0
	pkg.Body.TCaplusPbFieldUpdateReq.ValueInfo.CompactValueSet.FieldIndexs = nil
	pkg.Body.TCaplusPbFieldUpdateReq.ValueInfo.CompactValueSet.FieldIndexNum = 0
	pkg.Body.TCaplusPbFieldUpdateReq.ValueInfo.FieldNum_ = 0
	pkg.Body.TCaplusPbFieldUpdateReq.ValueInfo.Fields_ = nil
	pkg.Body.TCaplusPbFieldUpdateReq.CheckVersionType = 1
	pkg.Body.TCaplusPbFieldUpdateReq.Condition = ""
	req := &pbFieldUpdateRequest{
		appId:     appId,
		zoneId:    zoneId,
		tableName: tableName,
		cmd:       cmd,
		seq:       seq,
		record:    nil,
		pkg:       pkg,
		isPB:      isPB,
	}
	return req, nil
}

func (req *pbFieldUpdateRequest) AddRecord(index int32) (*record.Record, error) {
	if req.record != nil {
		return nil, &terror.ErrorCode{Code: terror.RecordNumOverMax}
	}

	rec := &record.Record{
		AppId:      req.appId,
		ZoneId:     req.zoneId,
		TableName:  req.tableName,
		Cmd:        req.cmd,
		KeyMap:     make(map[string][]byte),
		ValueMap:   make(map[string][]byte),
		Version:    -1,
		PBFieldMap: make(map[string]bool),
		IsPB:       req.isPB,
	}

	//key value set
	rec.ShardingKey = &req.pkg.Head.SplitTableKeyBuff
	rec.ShardingKeyLen = &req.pkg.Head.SplitTableKeyBuffLen
	rec.KeySet = req.pkg.Head.KeyInfo
	rec.PBValueSet = req.pkg.Body.TCaplusPbFieldUpdateReq.ValueInfo
	rec.Condition = &req.pkg.Body.TCaplusPbFieldUpdateReq.Condition
	rec.Operation = &req.pkg.Body.TCaplusPbFieldUpdateReq.Operation
	req.record = rec
	return rec, nil
}

func (req *pbFieldUpdateRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *pbFieldUpdateRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return &terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.TCaplusPbFieldUpdateReq.CheckVersionType = p
	return nil
}

func (req *pbFieldUpdateRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "FieldUpdate not Support ResultFlag"}
}

func (req *pbFieldUpdateRequest) Pack() ([]byte, error) {
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

	if err := req.record.PackPBFieldValue(); err != nil {
		logger.ERR("record pack value failed, %s", err.Error())
		return nil, err
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
		logger.DEBUG("%s", common.CovertToJson(req.pkg.Body.TCaplusPbFieldUpdateReq))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("pbFieldUpdateRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *pbFieldUpdateRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *pbFieldUpdateRequest) GetKeyHash() (uint32, error) {
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

func (req *pbFieldUpdateRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *pbFieldUpdateRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *pbFieldUpdateRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *pbFieldUpdateRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *pbFieldUpdateRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *pbFieldUpdateRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *pbFieldUpdateRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *pbFieldUpdateRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *pbFieldUpdateRequest) SetPerfTest(sendTime uint64) int {
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

func (req *pbFieldUpdateRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *pbFieldUpdateRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *pbFieldUpdateRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
