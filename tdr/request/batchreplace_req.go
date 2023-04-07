package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type batchReplaceRequest struct {
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

func newBatchReplaceRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*batchReplaceRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.BatchReplaceReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Head.KeyInfo.FieldNum = 0
	pkg.Body.BatchReplaceReq.AllowMultiResponses = 0
	pkg.Body.BatchReplaceReq.CheckVersiontType = 1
	pkg.Body.BatchReplaceReq.Flag = 0
	pkg.Body.BatchReplaceReq.RecordNum = 0
	pkg.Body.BatchReplaceReq.ValueLen = 0
	pkg.Body.BatchReplaceReq.ValueInfo = nil
	pkg.Body.BatchReplaceReq.KeyInfo = nil
	pkg.Body.BatchReplaceReq.SplitTableKeyBuffs = nil
	pkg.Body.BatchReplaceReq.ValueIndex = nil
	req := &batchReplaceRequest{
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

func (req *batchReplaceRequest) AddRecord(index int32) (*record.Record, error) {
	if len(req.record) >= 1024 {
		logger.ERR("record num > 1024")
		return nil, &terror.ErrorCode{Code: terror.RecordNumOverMax}
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
		SplitTableKeyBuff: nil,
	}

	rec.KeySet = new(tcaplus_protocol_cs.TCaplusKeySet)
	rec.ValueSet = tcaplus_protocol_cs.NewTCaplusValueSet_()
	rec.FieldIndex = tcaplus_protocol_cs.NewFieldIndex()

	rec.SplitTableKeyBuff = new(tcaplus_protocol_cs.SplitTableKeyBuff)
	req.pkg.Body.BatchReplaceReq.SplitTableKeyBuffs = append(req.pkg.Body.BatchReplaceReq.SplitTableKeyBuffs,
		rec.SplitTableKeyBuff)

	req.record = append(req.record, rec)

	return rec, nil
}

func (req *batchReplaceRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *batchReplaceRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.BatchReplaceReq.CheckVersiontType = p
	return nil
}

func (req *batchReplaceRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "BatchReplace not Support ResultFlag"}
}

func (req *batchReplaceRequest) Pack() ([]byte, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return nil, &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}
	if len(req.record) == 0 {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	for _, rec := range req.record {
		if err := rec.PackKey(); err != nil {
			logger.ERR("record pack key failed, %s", err.Error())
			return nil, err
		}
		req.pkg.Body.BatchReplaceReq.RecordNum += 1
		req.pkg.Body.BatchReplaceReq.KeyInfo = append(req.pkg.Body.BatchReplaceReq.KeyInfo, rec.KeySet)

		if err := rec.PackValue(nil); err != nil {
			logger.ERR("record pack key failed, %s", err.Error())
			return nil, err
		}
		rec.FieldIndex.Size = rec.ValueSet.CompactValueSet.ValueBufLen
		rec.FieldIndex.Offset = req.pkg.Body.BatchReplaceReq.ValueLen
		req.pkg.Body.BatchReplaceReq.ValueIndex = append(req.pkg.Body.BatchReplaceReq.ValueIndex,
			rec.FieldIndex)
		req.pkg.Body.BatchReplaceReq.ValueInfo = append(req.pkg.Body.BatchReplaceReq.ValueInfo,
			rec.ValueSet.CompactValueSet.ValueBuf...)
		req.pkg.Body.BatchReplaceReq.ValueLen += rec.ValueSet.CompactValueSet.ValueBufLen
	}

	req.pkg.Body.BatchReplaceReq.FieldName.FieldNum = 0
	for key, _ := range req.valueNameMap {
		req.pkg.Body.BatchReplaceReq.FieldName.FieldNum += 1
		req.pkg.Body.BatchReplaceReq.FieldName.FieldName = append(req.pkg.Body.BatchReplaceReq.FieldName.FieldName, key)
	}
	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("batchReplaceRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *batchReplaceRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *batchReplaceRequest) GetKeyHash() (uint32, error) {
	if len(req.record) == 0 {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.record[0].KeySet)
}

func (req *batchReplaceRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *batchReplaceRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *batchReplaceRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *batchReplaceRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *batchReplaceRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *batchReplaceRequest) SetAddableIncreaseFlag(increase_flag byte) int32 {
	return int32(terror.GEN_ERR_SUC)
}

func (req *batchReplaceRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	if 1 == multi_flag {
		req.pkg.Body.BatchReplaceReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.BatchReplaceReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *batchReplaceRequest) SetResultFlagForSuccess(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	flag = flag << 4
	flag |= 1 << 6
	req.pkg.Body.BatchReplaceReq.Flag |= flag
	return terror.GEN_ERR_SUC
}

func (req *batchReplaceRequest) SetResultFlagForFail(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	flag = flag << 2
	flag |= 1 << 6
	req.pkg.Body.BatchReplaceReq.Flag |= flag
	return terror.GEN_ERR_SUC
}

func (req *batchReplaceRequest) SetPerfTest(sendTime uint64) int {
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

func (req *batchReplaceRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *batchReplaceRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *batchReplaceRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
