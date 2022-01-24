package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type setTtlRequest struct {
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

func newSetTtlRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*setTtlRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusSetTTLReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.TCaplusSetTTLReq.RecordNum = 0
	pkg.Body.TCaplusSetTTLReq.KeyInfo = nil
	pkg.Body.TCaplusSetTTLReq.TTL = nil
	pkg.Body.TCaplusSetTTLReq.TTLType = nil
	pkg.Body.TCaplusSetTTLReq.SplitTableKeyBuffs = nil
	pkg.Body.TCaplusSetTTLReq.Index = nil
	req := &setTtlRequest{
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

func (req *setTtlRequest) AddRecord(index int32) (*record.Record, error) {
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

	req.pkg.Body.TCaplusSetTTLReq.Index = append(req.pkg.Body.TCaplusSetTTLReq.Index,
		index)

	rec.SplitTableKeyBuff = new(tcaplus_protocol_cs.SplitTableKeyBuff)
	req.pkg.Body.TCaplusSetTTLReq.SplitTableKeyBuffs = append(req.pkg.Body.TCaplusSetTTLReq.SplitTableKeyBuffs,
		rec.SplitTableKeyBuff)

	req.pkg.Body.TCaplusSetTTLReq.TTL = append(req.pkg.Body.TCaplusSetTTLReq.TTL,
		0)

	req.pkg.Body.TCaplusSetTTLReq.TTLType = append(req.pkg.Body.TCaplusSetTTLReq.TTLType,
		0)

	rec.Ttl = &req.pkg.Body.TCaplusSetTTLReq.TTL[req.idx]
	rec.TtlType = &req.pkg.Body.TCaplusSetTTLReq.TTLType[req.idx]

	req.idx++

	req.record = append(req.record, rec)
	return rec, nil
}

func (req *setTtlRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *setTtlRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support VersionPolicy"}
}

func (req *setTtlRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support ResultFlag"}
}

func (req *setTtlRequest) Pack() ([]byte, error) {
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
		req.pkg.Body.TCaplusSetTTLReq.RecordNum += 1
		req.pkg.Body.TCaplusSetTTLReq.KeyInfo = append(req.pkg.Body.TCaplusSetTTLReq.KeyInfo, rec.KeySet)
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack %d request %s", req.pkg.Body.TCaplusSetTTLReq.RecordNum, common.CsHeadVisualize(req.pkg.Head))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("setTtlRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *setTtlRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *setTtlRequest) GetKeyHash() (uint32, error) {
	if len(req.record) == 0 {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.record[0].KeySet)
}

func (req *setTtlRequest) SetFieldNames(valueNameList []string) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "set ttl not Support SetFieldNames"}
}

func (req *setTtlRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *setTtlRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *setTtlRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *setTtlRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *setTtlRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *setTtlRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *setTtlRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *setTtlRequest) SetPerfTest(sendTime uint64) int {
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

func (req *setTtlRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *setTtlRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *setTtlRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
