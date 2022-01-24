package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type batchUpdateRequest struct {
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

func newBatchUpdateRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*batchUpdateRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.BatchUpdateReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Head.KeyInfo.FieldNum = 0
	pkg.Body.BatchUpdateReq.AllowMultiResponses = 0
	pkg.Body.BatchUpdateReq.CheckVersiontType = 1
	pkg.Body.BatchUpdateReq.Flag = 0
	pkg.Body.BatchUpdateReq.RecordNum = 0
	pkg.Body.BatchUpdateReq.ValueLen = 0
	pkg.Body.BatchUpdateReq.ValueInfo = nil
	pkg.Body.BatchUpdateReq.KeyInfo = nil
	pkg.Body.BatchUpdateReq.SplitTableKeyBuffs = nil
	pkg.Body.BatchUpdateReq.ValueIndex = nil
	req := &batchUpdateRequest{
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

func (req *batchUpdateRequest) AddRecord(index int32) (*record.Record, error) {
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
	req.pkg.Body.BatchUpdateReq.SplitTableKeyBuffs = append(req.pkg.Body.BatchUpdateReq.SplitTableKeyBuffs,
		rec.SplitTableKeyBuff)

	req.record = append(req.record, rec)

	return rec, nil
}

func (req *batchUpdateRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *batchUpdateRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.BatchUpdateReq.CheckVersiontType = p
	return nil
}

func (req *batchUpdateRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "BatchUpdate not Support ResultFlag"}
}

func (req *batchUpdateRequest) Pack() ([]byte, error) {
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
		req.pkg.Body.BatchUpdateReq.RecordNum += 1
		req.pkg.Body.BatchUpdateReq.KeyInfo = append(req.pkg.Body.BatchUpdateReq.KeyInfo,
			rec.KeySet)

		if err := rec.PackValue(nil); err != nil {
			logger.ERR("record pack key failed, %s", err.Error())
			return nil, err
		}
		rec.FieldIndex.Size = rec.ValueSet.CompactValueSet.ValueBufLen
		rec.FieldIndex.Offset = req.pkg.Body.BatchUpdateReq.ValueLen
		req.pkg.Body.BatchUpdateReq.ValueIndex = append(req.pkg.Body.BatchUpdateReq.ValueIndex,
			rec.FieldIndex)
		req.pkg.Body.BatchUpdateReq.ValueInfo = append(req.pkg.Body.BatchUpdateReq.ValueInfo,
			rec.ValueSet.CompactValueSet.ValueBuf...)
		req.pkg.Body.BatchUpdateReq.ValueLen += rec.ValueSet.CompactValueSet.ValueBufLen
	}

	req.pkg.Body.BatchUpdateReq.FieldName.FieldNum = 0
	for key, _ := range req.valueNameMap {
		req.pkg.Body.BatchUpdateReq.FieldName.FieldNum += 1
		req.pkg.Body.BatchUpdateReq.FieldName.FieldName = append(req.pkg.Body.BatchUpdateReq.FieldName.FieldName,
			key)
	}
	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("batchUpdateRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *batchUpdateRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *batchUpdateRequest) GetKeyHash() (uint32, error) {
	if len(req.record) == 0 {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.record[0].KeySet)
}

func (req *batchUpdateRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *batchUpdateRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *batchUpdateRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *batchUpdateRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *batchUpdateRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *batchUpdateRequest) SetAddableIncreaseFlag(increase_flag byte) int32 {
	return int32(terror.GEN_ERR_SUC)
}

func (req *batchUpdateRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	if 1 == multi_flag {
		req.pkg.Body.BatchUpdateReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.BatchUpdateReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *batchUpdateRequest) SetResultFlagForSuccess(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	req.pkg.Body.BatchUpdateReq.Flag = flag << 4
	req.pkg.Body.BatchUpdateReq.Flag |= 1 << 6
	return terror.GEN_ERR_SUC
}

func (req *batchUpdateRequest) SetResultFlagForFail(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	req.pkg.Body.BatchUpdateReq.Flag = flag << 2
	req.pkg.Body.BatchUpdateReq.Flag |= 1 << 6
	return terror.GEN_ERR_SUC
}

func (req *batchUpdateRequest) SetPerfTest(sendTime uint64) int {
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

func (req *batchUpdateRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *batchUpdateRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *batchUpdateRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
