package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type listGetBatchRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	isPB      bool
}

func newListGetBatchRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*listGetBatchRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListGetBatchReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.ListGetBatchReq.ElementValueNames.FieldNum = 0
	pkg.Body.ListGetBatchReq.ElementValueNames.FieldName = nil
	pkg.Body.ListGetBatchReq.AllowMultiResponses = 0
	pkg.Body.ListGetBatchReq.ElementNum = 0
	pkg.Body.ListGetBatchReq.ElementIndexArray = nil
	req := &listGetBatchRequest{
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

func (req *listGetBatchRequest) AddRecord(index int32) (*record.Record, error) {
	if req.record != nil {
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
	}

	rec.KeySet = req.pkg.Head.KeyInfo
	rec.ShardingKey = &req.pkg.Head.SplitTableKeyBuff
	rec.ShardingKeyLen = &req.pkg.Head.SplitTableKeyBuffLen
	req.record = rec
	return rec, nil
}

func (req *listGetBatchRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listGetBatchRequest) SetVersionPolicy(p uint8) error {
	return nil
}

func (req *listGetBatchRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH, Message: "ResultFlag not support"}
}

func (req *listGetBatchRequest) Pack() ([]byte, error) {
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

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("getRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listGetBatchRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listGetBatchRequest) GetKeyHash() (uint32, error) {
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

func (req *listGetBatchRequest) SetFieldNames(valueNameList []string) error {
	req.pkg.Body.ListGetBatchReq.ElementValueNames.FieldName = valueNameList
	req.pkg.Body.ListGetBatchReq.ElementValueNames.FieldNum = uint32(len(valueNameList))
	return nil
}

func (req *listGetBatchRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listGetBatchRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listGetBatchRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *listGetBatchRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listGetBatchRequest) SetMultiResponseFlag(multiFlag byte) int32 {
	if 1 == multiFlag {
		req.pkg.Body.ListGetBatchReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.ListGetBatchReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *listGetBatchRequest) SetResultFlagForSuccess(flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *listGetBatchRequest) SetResultFlagForFail(flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *listGetBatchRequest) AddElementIndex(idx int32) int32 {
	bodyReq := req.pkg.Body.ListGetBatchReq
	for i := uint32(0); i < bodyReq.ElementNum; i++ {
		if idx == bodyReq.ElementIndexArray[i] {
			return int32(terror.GEN_ERR_SUC)
		}
	}
	if int64(bodyReq.ElementNum) > tcaplus_protocol_cs.TCAPLUS_MAX_LIST_ELEMENTS_NUM {
		return int32(terror.API_ERR_OVER_MAX_LIST_INDEX_NUM)
	}
	bodyReq.ElementIndexArray = append(bodyReq.ElementIndexArray, idx)
	bodyReq.ElementNum++
	return int32(terror.GEN_ERR_SUC)
}

func (req *listGetBatchRequest) SetPerfTest(sendTime uint64) int {
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

func (req *listGetBatchRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *listGetBatchRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *listGetBatchRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
