package request

import (
	"bytes"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type listAddAfterBatchRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    []*record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	isPB      bool
}

func newListAddAfterBatchRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*listAddAfterBatchRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListAddAfterBatchReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.ListAddAfterBatchReq.AllowMultiResponses = 0
	pkg.Body.ListAddAfterBatchReq.ElementNum = 0
	pkg.Body.ListAddAfterBatchReq.ElementIndexArray = nil
	pkg.Body.ListAddAfterBatchReq.ShiftFlag = 1
	pkg.Body.ListAddAfterBatchReq.ValueIndex = nil
	pkg.Body.ListAddAfterBatchReq.ValueLen = 0
	pkg.Body.ListAddAfterBatchReq.ValueInfo = nil
	pkg.Body.ListAddAfterBatchReq.CheckVersiontType = 1

	req := &listAddAfterBatchRequest{
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

func (req *listAddAfterBatchRequest) AddRecord(index int32) (*record.Record, error) {
	if len(req.record) >= 1024 {
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
		Index:       index,
	}

	rec.KeySet = req.pkg.Head.KeyInfo
	rec.ValueSet = tcaplus_protocol_cs.NewTCaplusValueSet_()
	rec.ShardingKey = &req.pkg.Head.SplitTableKeyBuff
	rec.ShardingKeyLen = &req.pkg.Head.SplitTableKeyBuffLen
	rec.FieldIndex = tcaplus_protocol_cs.NewFieldIndex()
	req.record = append(req.record, rec)
	return rec, nil
}

func (req *listAddAfterBatchRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listAddAfterBatchRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return &terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.ListAddAfterBatchReq.CheckVersiontType = p
	return nil
}

func (req *listAddAfterBatchRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH, Message: "ResultFlag not support"}
}

func (req *listAddAfterBatchRequest) Pack() ([]byte, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return nil, &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}
	if len(req.record) == 0 {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	if err := req.record[0].PackKey(); err != nil {
		logger.ERR("record pack key failed, %s", err.Error())
		return nil, err
	}

	req.pkg.Body.ListAddAfterBatchReq.ValueIndex = make([]*tcaplus_protocol_cs.FieldIndex, len(req.record),
		len(req.record))
	req.pkg.Body.ListAddAfterBatchReq.ElementIndexArray = make([]int32, len(req.record), len(req.record))
	req.pkg.Body.ListAddAfterBatchReq.ElementNum = uint32(len(req.record))
	for i, rec := range req.record {
		if err := rec.PackValue(nil); err != nil {
			logger.ERR("record pack key failed, %s", err.Error())
			return nil, err
		}
		rec.FieldIndex.Size = rec.ValueSet.CompactValueSet.ValueBufLen
		rec.FieldIndex.Offset = req.pkg.Body.ListAddAfterBatchReq.ValueLen
		req.pkg.Body.ListAddAfterBatchReq.ValueIndex[i] = rec.FieldIndex
		req.pkg.Body.ListAddAfterBatchReq.ElementIndexArray[i] = rec.Index
		req.pkg.Body.ListAddAfterBatchReq.ValueLen += rec.ValueSet.CompactValueSet.ValueBufLen
	}
	valueBuf := new(bytes.Buffer)
	valueBuf.Grow(int(req.pkg.Body.ListAddAfterBatchReq.ValueLen))
	valueBuf.Reset()
	for _, rec := range req.record {
		valueBuf.Write(rec.ValueSet.CompactValueSet.ValueBuf)
	}
	req.pkg.Body.ListAddAfterBatchReq.ValueInfo = valueBuf.Bytes()

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
		logger.DEBUG("%s", common.CovertToJson(req.pkg.Body.ListAddAfterBatchReq))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("listAddAfterBatchRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listAddAfterBatchRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listAddAfterBatchRequest) GetKeyHash() (uint32, error) {
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

func (req *listAddAfterBatchRequest) SetFieldNames(valueNameList []string) error {
	return terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH,
		Message: "listAddAfterBatchRequest not support SetFieldNames"}
}

func (req *listAddAfterBatchRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listAddAfterBatchRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listAddAfterBatchRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *listAddAfterBatchRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listAddAfterBatchRequest) SetMultiResponseFlag(multiFlag byte) int32 {
	if 1 == multiFlag {
		req.pkg.Body.ListAddAfterBatchReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.ListAddAfterBatchReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *listAddAfterBatchRequest) SetResultFlagForSuccess(flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *listAddAfterBatchRequest) SetResultFlagForFail(flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *listAddAfterBatchRequest) SetPerfTest(sendTime uint64) int {
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

func (req *listAddAfterBatchRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *listAddAfterBatchRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *listAddAfterBatchRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}

func (req *listAddAfterBatchRequest) SetListShiftFlag(shiftFlag byte) int32 {
	req.pkg.Body.ListAddAfterBatchReq.ShiftFlag = shiftFlag
	return int32(terror.GEN_ERR_SUC)
}
