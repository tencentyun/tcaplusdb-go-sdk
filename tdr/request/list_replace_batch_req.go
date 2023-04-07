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

type listReplaceBatchRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       []*record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
	isPB         bool
}

func newListReplaceBatchRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*listReplaceBatchRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListReplaceBatchReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.ListReplaceBatchReq.AllowMultiResponses = 0
	pkg.Body.ListReplaceBatchReq.ElementNum = 0
	pkg.Body.ListReplaceBatchReq.ElementIndexArray = nil
	pkg.Body.ListReplaceBatchReq.Flag = 0
	pkg.Body.ListReplaceBatchReq.ValueIndex = nil
	pkg.Body.ListReplaceBatchReq.ValueLen = 0
	pkg.Body.ListReplaceBatchReq.ValueInfo = nil
	pkg.Body.ListReplaceBatchReq.CheckVersiontType = 1

	req := &listReplaceBatchRequest{
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

func (req *listReplaceBatchRequest) AddRecord(index int32) (*record.Record, error) {
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

func (req *listReplaceBatchRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listReplaceBatchRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return &terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.ListReplaceBatchReq.CheckVersiontType = p
	return nil
}

func (req *listReplaceBatchRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH, Message: "ResultFlag not support"}
}

func (req *listReplaceBatchRequest) Pack() ([]byte, error) {
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

	req.pkg.Body.ListReplaceBatchReq.ValueIndex = make([]*tcaplus_protocol_cs.FieldIndex, len(req.record),
		len(req.record))
	req.pkg.Body.ListReplaceBatchReq.ElementIndexArray = make([]int32, len(req.record), len(req.record))
	req.pkg.Body.ListReplaceBatchReq.ElementNum = uint32(len(req.record))
	for i, rec := range req.record {
		if err := rec.PackValue(req.valueNameMap); err != nil {
			logger.ERR("record pack key failed, %s", err.Error())
			return nil, err
		}
		rec.FieldIndex.Size = rec.ValueSet.CompactValueSet.ValueBufLen
		rec.FieldIndex.Offset = req.pkg.Body.ListReplaceBatchReq.ValueLen
		req.pkg.Body.ListReplaceBatchReq.ValueIndex[i] = rec.FieldIndex
		req.pkg.Body.ListReplaceBatchReq.ElementIndexArray[i] = rec.Index
		req.pkg.Body.ListReplaceBatchReq.ValueLen += rec.ValueSet.CompactValueSet.ValueBufLen
	}
	valueBuf := new(bytes.Buffer)
	valueBuf.Grow(int(req.pkg.Body.ListReplaceBatchReq.ValueLen))
	valueBuf.Reset()
	for _, rec := range req.record {
		valueBuf.Write(rec.ValueSet.CompactValueSet.ValueBuf)
	}
	req.pkg.Body.ListReplaceBatchReq.ValueInfo = valueBuf.Bytes()

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("listReplaceBatchRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listReplaceBatchRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listReplaceBatchRequest) GetKeyHash() (uint32, error) {
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

func (req *listReplaceBatchRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *listReplaceBatchRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listReplaceBatchRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listReplaceBatchRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *listReplaceBatchRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listReplaceBatchRequest) SetMultiResponseFlag(multiFlag byte) int32 {
	if 1 == multiFlag {
		req.pkg.Body.ListReplaceBatchReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.ListReplaceBatchReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *listReplaceBatchRequest) SetResultFlagForSuccess(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	flag = flag << 4
	flag |= 1 << 6
	req.pkg.Body.ListReplaceBatchReq.Flag |= flag
	return terror.GEN_ERR_SUC
}

func (req *listReplaceBatchRequest) SetResultFlagForFail(flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *listReplaceBatchRequest) SetPerfTest(sendTime uint64) int {
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

func (req *listReplaceBatchRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *listReplaceBatchRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *listReplaceBatchRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
