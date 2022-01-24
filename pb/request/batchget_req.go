package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type batchGetRequest struct {
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

func newBatchGetRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*batchGetRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.BatchGetReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Head.KeyInfo.FieldNum = 0
	pkg.Body.BatchGetReq.ValueInfo.FieldNum = 0
	pkg.Body.BatchGetReq.KeyInfo = nil
	pkg.Body.BatchGetReq.AllowMultiResponses = 0
	pkg.Body.BatchGetReq.ExpireTime = 0
	pkg.Body.BatchGetReq.RecordNum = 0
	pkg.Body.BatchGetReq.SplitTableKeyBuffs = nil
	req := &batchGetRequest{
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

func (req *batchGetRequest) AddRecord(index int32) (*record.Record, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return nil, &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}

	if len(req.record) >= 1024 {
		logger.ERR("record num > 1024")
		return nil, &terror.ErrorCode{Code: terror.RecordNumOverMax}
	}
	//batchGetReq := req.pkg.Body.BatchGetReq
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
		IsPB:              req.isPB,
	}

	rec.KeySet = new(tcaplus_protocol_cs.TCaplusKeySet)

	rec.SplitTableKeyBuff = new(tcaplus_protocol_cs.SplitTableKeyBuff)
	req.pkg.Body.BatchGetReq.SplitTableKeyBuffs = append(req.pkg.Body.BatchGetReq.SplitTableKeyBuffs,
		rec.SplitTableKeyBuff)

	req.record = append(req.record, rec)
	return rec, nil
}

func (req *batchGetRequest) SetAsyncId(id uint64) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
	}
	req.pkg.Head.AsynID = id
}

func (req *batchGetRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "BatchGet not Support VersionPolicy"}
}

func (req *batchGetRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "BatchGet not Support ResultFlag"}
}

func (req *batchGetRequest) Pack() ([]byte, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return nil, &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}

	if len(req.record) == 0 {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	req.pkg.Body.BatchGetReq.RecordNum = 0
	req.pkg.Body.BatchGetReq.KeyInfo = make([]*tcaplus_protocol_cs.TCaplusKeySet, len(req.record))
	for _, rec := range req.record {
		if err := rec.PackKey(); err != nil {
			logger.ERR("record pack key failed, %s", err.Error())
			return nil, err
		}
		req.pkg.Body.BatchGetReq.KeyInfo[req.pkg.Body.BatchGetReq.RecordNum] = rec.KeySet
		req.pkg.Body.BatchGetReq.RecordNum++
	}

	if req.isPB {
		req.pkg.Body.BatchGetReq.ValueInfo.FieldNum = 3
		req.pkg.Body.BatchGetReq.ValueInfo.FieldName = []string{"klen", "vlen", "value"}
	} else {
		if len(req.valueNameMap) > 0 {
			req.record[0].ValueMap = make(map[string][]byte)
			for name, _ := range req.valueNameMap {
				req.record[0].ValueMap[name] = []byte{}
			}
		}
		nameSet := req.pkg.Body.BatchGetReq.ValueInfo
		nameSet.FieldNum = 0
		nameSet.FieldName = make([]string, len(req.record[0].ValueMap))
		for key, _ := range req.record[0].ValueMap {
			nameSet.FieldName[nameSet.FieldNum] = key
			nameSet.FieldNum++
		}
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("batchGetRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *batchGetRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *batchGetRequest) GetKeyHash() (uint32, error) {
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
	return keyHashCode(req.record[0].KeySet)
}

func (req *batchGetRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *batchGetRequest) SetUserBuff(userBuffer []byte) error {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *batchGetRequest) GetSeq() int32 {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return 0
	}
	return req.pkg.Head.Seq
}

func (req *batchGetRequest) SetSeq(seq int32) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return
	}
	req.pkg.Head.Seq = seq
}

func (req *batchGetRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *batchGetRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return int32(terror.RequestHasHasNoPkg)
	}
	if 0 != multi_flag {
		req.pkg.Body.BatchGetReq.AllowMultiResponses = 1
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *batchGetRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *batchGetRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *batchGetRequest) SetPerfTest(sendTime uint64) int {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return int(terror.RequestHasHasNoPkg)
	}
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

func (req *batchGetRequest) SetFlags(flag int32) int {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return int(terror.RequestHasHasNoPkg)
	}
	return setFlags(req.pkg, flag)
}

func (req *batchGetRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *batchGetRequest) GetFlags() int32 {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return int32(terror.RequestHasHasNoPkg)
	}
	return req.pkg.Head.Flags
}
