package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type listGetAllRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       *record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
	isPB         bool
}

func newListGetAllRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*listGetAllRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListGetAllReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Head.KeyInfo.FieldNum = 0
	pkg.Body.ListGetAllReq.AllowMultiResponses = 0
	pkg.Body.ListGetAllReq.ElementNum = -1
	pkg.Body.ListGetAllReq.StartSubscript = 0
	pkg.Body.ListGetAllReq.ElementValueNames.FieldNum = 0
	pkg.Body.ListGetAllReq.ElementValueNames.FieldName = nil

	req := &listGetAllRequest{
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

func (req *listGetAllRequest) AddRecord(index int32) (*record.Record, error) {
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

	rec.ShardingKey = &req.pkg.Head.SplitTableKeyBuff
	rec.ShardingKeyLen = &req.pkg.Head.SplitTableKeyBuffLen
	rec.KeySet = req.pkg.Head.KeyInfo
	rec.Condition = &req.pkg.Body.ListGetAllReq.Condition
	//rec.ValueSet = req.pkg.Body.ListGetAllReq.ElementValueNames
	req.record = rec
	return rec, nil
}

func (req *listGetAllRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listGetAllRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "listGetAll not Support VersionPolicy"}
}

func (req *listGetAllRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "listGetAll not Support ResultFlag"}
}

func (req *listGetAllRequest) Pack() ([]byte, error) {
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

	if req.isPB {
		//req.pkg.Body.ListGetAllReq.ElementValueNames.FieldNum = 3
		//req.pkg.Body.ListGetAllReq.ElementValueNames.FieldName = []string{"klen", "vlen", "value"}
	} else {
		if len(req.valueNameMap) > 0 {
			for name, _ := range req.valueNameMap {
				req.record.ValueMap[name] = []byte{}
			}
		}
		nameSet := req.pkg.Body.ListGetAllReq.ElementValueNames
		nameSet.FieldNum = 0
		nameSet.FieldName = make([]string, len(req.record.ValueMap))
		for key, _ := range req.record.ValueMap {
			nameSet.FieldName[nameSet.FieldNum] = key
			nameSet.FieldNum++
		}
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

func (req *listGetAllRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listGetAllRequest) GetKeyHash() (uint32, error) {
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

func (req *listGetAllRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *listGetAllRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listGetAllRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listGetAllRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}
func (req *listGetAllRequest) SetResultLimit(limit int32, offset int32) int32 {
	req.pkg.Body.ListGetAllReq.StartSubscript = offset
	req.pkg.Body.ListGetAllReq.ElementNum = limit
	return int32(terror.GEN_ERR_SUC)
}

func (req *listGetAllRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	if 1 == multi_flag {
		req.pkg.Body.ListGetAllReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.ListGetAllReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *listGetAllRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *listGetAllRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *listGetAllRequest) SetPerfTest(sendTime uint64) int {
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

func (req *listGetAllRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *listGetAllRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *listGetAllRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
