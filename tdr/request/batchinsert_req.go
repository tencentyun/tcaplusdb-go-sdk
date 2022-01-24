package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type batchInsertRequest struct {
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

func newBatchInsertRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*batchInsertRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.BatchInsertReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Head.KeyInfo.FieldNum = 0
	pkg.Body.BatchInsertReq.AllowMultiResponses = 0
	pkg.Body.BatchInsertReq.Flag = 0
	pkg.Body.BatchInsertReq.RecordNum = 0
	pkg.Body.BatchInsertReq.ValueLen = 0
	pkg.Body.BatchInsertReq.ValueInfo = nil
	pkg.Body.BatchInsertReq.KeyInfo = nil
	pkg.Body.BatchInsertReq.SplitTableKeyBuffs = nil
	pkg.Body.BatchInsertReq.ValueIndex = nil

	req := &batchInsertRequest{
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

func (req *batchInsertRequest) AddRecord(index int32) (*record.Record, error) {
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
	req.pkg.Body.BatchInsertReq.SplitTableKeyBuffs = append(req.pkg.Body.BatchInsertReq.SplitTableKeyBuffs,
		rec.SplitTableKeyBuff)

	req.record = append(req.record, rec)

	return rec, nil
}

func (req *batchInsertRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *batchInsertRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "BatchInsert not Support VersionPolicy"}
}

func (req *batchInsertRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "BatchInsert not Support ResultFlag"}
}

func (req *batchInsertRequest) Pack() ([]byte, error) {
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
		req.pkg.Body.BatchInsertReq.RecordNum += 1
		req.pkg.Body.BatchInsertReq.KeyInfo = append(req.pkg.Body.BatchInsertReq.KeyInfo, rec.KeySet)

		if err := rec.PackValue(nil); err != nil {
			logger.ERR("record pack key failed, %s", err.Error())
			return nil, err
		}
		rec.FieldIndex.Size = rec.ValueSet.CompactValueSet.ValueBufLen
		rec.FieldIndex.Offset = req.pkg.Body.BatchInsertReq.ValueLen
		req.pkg.Body.BatchInsertReq.ValueIndex = append(req.pkg.Body.BatchInsertReq.ValueIndex,
			rec.FieldIndex)
		req.pkg.Body.BatchInsertReq.ValueInfo = append(req.pkg.Body.BatchInsertReq.ValueInfo,
			rec.ValueSet.CompactValueSet.ValueBuf...)
		req.pkg.Body.BatchInsertReq.ValueLen += rec.ValueSet.CompactValueSet.ValueBufLen
	}
	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("batchInsertRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *batchInsertRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *batchInsertRequest) GetKeyHash() (uint32, error) {
	if len(req.record) == 0 {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.record[0].KeySet)
}

func (req *batchInsertRequest) SetFieldNames(valueNameList []string) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "batch insert not Support SetFieldNames"}
}

func (req *batchInsertRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *batchInsertRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *batchInsertRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *batchInsertRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *batchInsertRequest) SetAddableIncreaseFlag(increaseFlag byte) int32 {
	return int32(terror.GEN_ERR_SUC)
}

func (req *batchInsertRequest) SetMultiResponseFlag(multiFlag byte) int32 {
	if 1 == multiFlag {
		req.pkg.Body.BatchInsertReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.BatchInsertReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *batchInsertRequest) SetResultFlagForSuccess(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	req.pkg.Body.BatchInsertReq.Flag = flag << 4
	req.pkg.Body.BatchInsertReq.Flag |= 1 << 6
	return terror.GEN_ERR_SUC
}

func (req *batchInsertRequest) SetResultFlagForFail(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	req.pkg.Body.BatchInsertReq.Flag = flag << 2
	req.pkg.Body.BatchInsertReq.Flag |= 1 << 6
	return terror.GEN_ERR_SUC
}

func (req *batchInsertRequest) SetPerfTest(sendTime uint64) int {
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

func (req *batchInsertRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *batchInsertRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *batchInsertRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
