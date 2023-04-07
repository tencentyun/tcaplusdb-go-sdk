package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type updateRequest struct {
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

func newUpdateRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*updateRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.UpdateReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.UpdateReq.ValueInfo.EncodeType = 1
	pkg.Body.UpdateReq.ValueInfo.Version_ = 0
	pkg.Body.UpdateReq.ValueInfo.CompactValueSet.ValueBuf = nil
	pkg.Body.UpdateReq.ValueInfo.CompactValueSet.ValueBufLen = 0
	pkg.Body.UpdateReq.ValueInfo.CompactValueSet.FieldIndexs = nil
	pkg.Body.UpdateReq.ValueInfo.CompactValueSet.FieldIndexNum = 0
	pkg.Body.UpdateReq.ValueInfo.FieldNum_ = 0
	pkg.Body.UpdateReq.ValueInfo.Fields_ = nil
	pkg.Body.UpdateReq.Flag = 0
	pkg.Body.UpdateReq.CheckVersiontType = 1
	pkg.Body.UpdateReq.IncreaseValueInfo.FieldNum = 0
	pkg.Body.UpdateReq.IncreaseValueInfo.Fields = nil
	pkg.Body.UpdateReq.IncreaseValueInfo.Version = 0
	pkg.Body.UpdateReq.Condition = ""
	req := &updateRequest{
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

func (req *updateRequest) AddRecord(index int32) (*record.Record, error) {
	if req.record != nil {
		return nil, &terror.ErrorCode{Code: terror.RecordNumOverMax}
	}
	rec := record.GetPoolRecord()
	rec.AppId = req.appId
	rec.ZoneId = req.zoneId
	rec.TableName = req.tableName
	rec.Cmd = req.cmd
	rec.KeyMap = make(map[string][]byte)
	rec.ValueMap = make(map[string][]byte)
	rec.IsPB = req.isPB

	//key value set
	rec.ShardingKey = &req.pkg.Head.SplitTableKeyBuff
	rec.ShardingKeyLen = &req.pkg.Head.SplitTableKeyBuffLen
	rec.KeySet = req.pkg.Head.KeyInfo
	rec.ValueSet = req.pkg.Body.UpdateReq.ValueInfo
	rec.Condition = &req.pkg.Body.UpdateReq.Condition
	rec.Operation = &req.pkg.Body.UpdateReq.Operation
	rec.UpdFieldSet = req.pkg.Body.UpdateReq.IncreaseValueInfo
	req.record = rec
	return rec, nil
}

func (req *updateRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *updateRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return &terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.UpdateReq.CheckVersiontType = p
	return nil
}

func (req *updateRequest) SetResultFlag(flag int) error {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "ResultFlag invalid"}
	}
	req.pkg.Body.UpdateReq.Flag = byte(flag)
	return nil
}

func (req *updateRequest) Pack() ([]byte, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return nil, &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}

	if req.record == nil {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	defer record.PutPoolRecord(req.record)
	if err := req.record.PackKey(); err != nil {
		logger.ERR("record pack key failed, %s", err.Error())
		return nil, err
	}

	if err := req.record.PackValue(req.valueNameMap); err != nil {
		logger.ERR("record pack value failed, %s", err.Error())
		return nil, err
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
		logger.DEBUG("%s", common.CovertToJson(req.pkg.Body.UpdateReq))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("updateRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *updateRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *updateRequest) GetKeyHash() (uint32, error) {
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

func (req *updateRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *updateRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *updateRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *updateRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}
func (req *updateRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *updateRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *updateRequest) SetResultFlagForSuccess(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	flag = flag << 4
	flag |= 1 << 6
	req.pkg.Body.UpdateReq.Flag |= flag
	return terror.GEN_ERR_SUC
}

func (req *updateRequest) SetResultFlagForFail(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	flag = flag << 2
	flag |= 1 << 6
	req.pkg.Body.UpdateReq.Flag |= flag
	return terror.GEN_ERR_SUC
}

func (req *updateRequest) SetPerfTest(sendTime uint64) int {
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

func (req *updateRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *updateRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *updateRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
