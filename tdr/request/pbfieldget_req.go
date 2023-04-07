package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type pbFieldGetRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	isPB      bool
}

func newPBFieldGetRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*pbFieldGetRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusPbFieldGetReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.TCaplusPbFieldGetReq.ValueInfo.EncodeType = 1
	pkg.Body.TCaplusPbFieldGetReq.ValueInfo.Version_ = 0
	pkg.Body.TCaplusPbFieldGetReq.ValueInfo.CompactValueSet.ValueBuf = nil
	pkg.Body.TCaplusPbFieldGetReq.ValueInfo.CompactValueSet.ValueBufLen = 0
	pkg.Body.TCaplusPbFieldGetReq.ValueInfo.CompactValueSet.FieldIndexs = nil
	pkg.Body.TCaplusPbFieldGetReq.ValueInfo.CompactValueSet.FieldIndexNum = 0
	pkg.Body.TCaplusPbFieldGetReq.ValueInfo.FieldNum_ = 0
	pkg.Body.TCaplusPbFieldGetReq.ValueInfo.Fields_ = nil
	pkg.Body.TCaplusPbFieldGetReq.Condition = ""
	req := &pbFieldGetRequest{
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

func (req *pbFieldGetRequest) AddRecord(index int32) (*record.Record, error) {
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
	rec.PBValueSet = req.pkg.Body.TCaplusPbFieldGetReq.ValueInfo
	rec.Condition = &req.pkg.Body.TCaplusPbFieldGetReq.Condition
	req.record = rec
	return rec, nil
}

func (req *pbFieldGetRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *pbFieldGetRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "FieldGet not Support VersionPolicy"}
}

func (req *pbFieldGetRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "FieldGet not Support ResultFlag"}
}

func (req *pbFieldGetRequest) Pack() ([]byte, error) {
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
		logger.DEBUG("%s", common.CovertToJson(req.pkg.Body.TCaplusPbFieldGetReq))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("pbFieldGetRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *pbFieldGetRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *pbFieldGetRequest) GetKeyHash() (uint32, error) {
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

func (req *pbFieldGetRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *pbFieldGetRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *pbFieldGetRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *pbFieldGetRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *pbFieldGetRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *pbFieldGetRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *pbFieldGetRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *pbFieldGetRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *pbFieldGetRequest) SetPerfTest(sendTime uint64) int {
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

func (req *pbFieldGetRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *pbFieldGetRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *pbFieldGetRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
