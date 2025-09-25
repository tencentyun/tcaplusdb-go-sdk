package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

type pbBatchFieldGetRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       []*record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
	idx          int
	isPB         bool
}

func newPBBatchFieldGetRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*pbBatchFieldGetRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusPbBatchFieldGetReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.TCaplusPbBatchFieldGetReq.ValueInfo.EncodeType = 1
	pkg.Body.TCaplusPbBatchFieldGetReq.ValueInfo.Version_ = 0
	pkg.Body.TCaplusPbBatchFieldGetReq.ValueInfo.CompactValueSet.ValueBuf = nil
	pkg.Body.TCaplusPbBatchFieldGetReq.ValueInfo.CompactValueSet.ValueBufLen = 0
	pkg.Body.TCaplusPbBatchFieldGetReq.ValueInfo.CompactValueSet.FieldIndexs = nil
	pkg.Body.TCaplusPbBatchFieldGetReq.ValueInfo.CompactValueSet.FieldIndexNum = 0
	pkg.Body.TCaplusPbBatchFieldGetReq.ValueInfo.FieldNum_ = 0
	pkg.Body.TCaplusPbBatchFieldGetReq.SplitTableKeyBuffs = nil
	pkg.Body.TCaplusPbBatchFieldGetReq.KeyInfo = nil
	pkg.Body.TCaplusPbBatchFieldGetReq.RecordNum = 0
	req := &pbBatchFieldGetRequest{
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

func (req *pbBatchFieldGetRequest) AddRecord(index int32) (*record.Record, error) {
	if len(req.record) >= 1024 {
		logger.ERR("record num > 1024")
		return nil, &terror.ErrorCode{Code: terror.RecordNumOverMax}
	}

	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return nil, &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}

	rec := &record.Record{
		AppId:             req.appId,
		ZoneId:            req.zoneId,
		TableName:         req.tableName,
		Cmd:               req.cmd,
		KeyMap:            make(map[string][]byte),
		ValueMap:          make(map[string][]byte),
		Version:           -1,
		KeySet:            nil,
		ValueSet:          nil,
		UpdFieldSet:       nil,
		PBFieldMap:        make(map[string]bool),
		SplitTableKeyBuff: nil,
		IsPB:              req.isPB,
	}
	rec.PBValueSet = req.pkg.Body.TCaplusPbBatchFieldGetReq.ValueInfo
	rec.KeySet = new(tcaplus_protocol_cs.TCaplusKeySet)
	req.pkg.Body.TCaplusPbBatchFieldGetReq.KeyInfo = append(req.pkg.Body.TCaplusPbBatchFieldGetReq.KeyInfo, rec.KeySet)
	rec.SplitTableKeyBuff = new(tcaplus_protocol_cs.SplitTableKeyBuff)
	req.pkg.Body.TCaplusPbBatchFieldGetReq.SplitTableKeyBuffs = append(req.pkg.Body.TCaplusPbBatchFieldGetReq.SplitTableKeyBuffs,
		rec.SplitTableKeyBuff)
	req.record = append(req.record, rec)
	return rec, nil
}

func (req *pbBatchFieldGetRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *pbBatchFieldGetRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "FieldGet not Support VersionPolicy"}
}

func (req *pbBatchFieldGetRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "FieldGet not Support ResultFlag"}
}

func (req *pbBatchFieldGetRequest) Pack() ([]byte, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return nil, &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}

	if len(req.record) == 0 {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	req.pkg.Body.TCaplusPbBatchFieldGetReq.RecordNum = 0
	for _, rec := range req.record {
		if err := rec.PackKey(); err != nil {
			logger.ERR("record pack key failed, %s", err.Error())
			return nil, err
		}
		req.pkg.Body.TCaplusPbBatchFieldGetReq.RecordNum++
	}

	if err := req.record[0].PackPBFieldValue(); err != nil {
		logger.ERR("record pack value failed, %s", err.Error())
		return nil, err
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
		logger.DEBUG("%s", common.CovertToJson(req.pkg.Body.TCaplusPbBatchFieldGetReq))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("pbBatchFieldGetRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *pbBatchFieldGetRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *pbBatchFieldGetRequest) GetKeyHash() (uint32, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return uint32(terror.RequestHasHasNoPkg), &terror.ErrorCode{Code: terror.RequestHasHasNoPkg,
			Message: "Request can not second use"}
	}
	defer func() {
		cs_pool.PutTcaplusCSPkg(req.pkg)
		req.pkg = nil
	}()

	if len(req.record) == 0 {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return uint32(time.Now().UnixNano()), nil
}

func (req *pbBatchFieldGetRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *pbBatchFieldGetRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *pbBatchFieldGetRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *pbBatchFieldGetRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *pbBatchFieldGetRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *pbBatchFieldGetRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *pbBatchFieldGetRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *pbBatchFieldGetRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *pbBatchFieldGetRequest) SetPerfTest(sendTime uint64) int {
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

func (req *pbBatchFieldGetRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *pbBatchFieldGetRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *pbBatchFieldGetRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
