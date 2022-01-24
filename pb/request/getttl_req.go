package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type getTtlRequest struct {
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

func newGetTtlRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*getTtlRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusGetTTLReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.TCaplusGetTTLReq.RecordNum = 0
	pkg.Body.TCaplusGetTTLReq.KeyInfo = nil
	pkg.Body.TCaplusGetTTLReq.SplitTableKeyBuffs = nil
	pkg.Body.TCaplusGetTTLReq.Index = nil
	req := &getTtlRequest{
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

func (req *getTtlRequest) AddRecord(index int32) (*record.Record, error) {
	if req.idx >= 1024 {
		logger.ERR("record num > 1024")
		return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "record num > 1024"}
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

	rec.KeySet = new(tcaplus_protocol_cs.TCaplusKeySet)

	req.pkg.Body.TCaplusGetTTLReq.Index = append(req.pkg.Body.TCaplusGetTTLReq.Index,
		index)

	rec.SplitTableKeyBuff = new(tcaplus_protocol_cs.SplitTableKeyBuff)
	req.pkg.Body.TCaplusGetTTLReq.SplitTableKeyBuffs = append(req.pkg.Body.TCaplusGetTTLReq.SplitTableKeyBuffs,
		rec.SplitTableKeyBuff)

	req.record = append(req.record, rec)
	return rec, nil
}

func (req *getTtlRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *getTtlRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support VersionPolicy"}
}

func (req *getTtlRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support ResultFlag"}
}

func (req *getTtlRequest) Pack() ([]byte, error) {
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
		req.pkg.Body.TCaplusGetTTLReq.RecordNum += 1
		req.pkg.Body.TCaplusGetTTLReq.KeyInfo = append(req.pkg.Body.TCaplusGetTTLReq.KeyInfo, rec.KeySet)
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("getTtlRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *getTtlRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *getTtlRequest) GetKeyHash() (uint32, error) {
	if len(req.record) == 0 {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.record[0].KeySet)
}

func (req *getTtlRequest) SetFieldNames(valueNameList []string) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "get ttl not Support SetFieldNames"}
}

func (req *getTtlRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *getTtlRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *getTtlRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *getTtlRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *getTtlRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *getTtlRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *getTtlRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *getTtlRequest) SetPerfTest(sendTime uint64) int {
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

func (req *getTtlRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *getTtlRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *getTtlRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
