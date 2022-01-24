package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type getByPartKeyRequest struct {
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

func newGetByPartKeyRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*getByPartKeyRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.GetByPartKeyReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.GetByPartKeyReq.OffSet = 0
	pkg.Body.GetByPartKeyReq.Limit = -1
	pkg.Body.GetByPartKeyReq.ValueInfo.FieldNum = 0
	pkg.Body.GetByPartKeyReq.ValueInfo.FieldName = nil
	req := &getByPartKeyRequest{
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

func (req *getByPartKeyRequest) AddRecord(index int32) (*record.Record, error) {
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
	rec.Condition = &req.pkg.Body.GetByPartKeyReq.Condition
	req.record = rec
	return rec, nil
}

func (req *getByPartKeyRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *getByPartKeyRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "GetByPartkey not Support VersionPolicy"}
}

func (req *getByPartKeyRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "GetByPartkey not Support ResultFlag"}
}

func (req *getByPartKeyRequest) Pack() ([]byte, error) {
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
		req.pkg.Body.GetByPartKeyReq.ValueInfo.FieldNum = 3
		req.pkg.Body.GetByPartKeyReq.ValueInfo.FieldName = []string{"klen", "vlen", "value"}
	} else {
		if len(req.valueNameMap) > 0 {
			req.record.ValueMap = make(map[string][]byte)
			for name, _ := range req.valueNameMap {
				req.record.ValueMap[name] = []byte{}
			}
		}
		nameSet := req.pkg.Body.GetByPartKeyReq.ValueInfo
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

func (req *getByPartKeyRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *getByPartKeyRequest) GetKeyHash() (uint32, error) {
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

func (req *getByPartKeyRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *getByPartKeyRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *getByPartKeyRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *getByPartKeyRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}
func (req *getByPartKeyRequest) SetResultLimit(limit int32, offset int32) int32 {
	req.pkg.Body.GetByPartKeyReq.OffSet = offset
	req.pkg.Body.GetByPartKeyReq.Limit = limit
	return int32(terror.GEN_ERR_SUC)
}

func (req *getByPartKeyRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *getByPartKeyRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *getByPartKeyRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *getByPartKeyRequest) SetPerfTest(sendTime uint64) int {
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

func (req *getByPartKeyRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *getByPartKeyRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *getByPartKeyRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
