package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type listAddAfterRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	isPB      bool
}

func newListAddAfterRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*listAddAfterRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListAddAfterReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	pkg.Body.ListAddAfterReq.ElementValueInfo.EncodeType = 1
	pkg.Body.ListAddAfterReq.ElementValueInfo.Version_ = 0
	pkg.Body.ListAddAfterReq.ElementValueInfo.Fields_ = nil
	pkg.Body.ListAddAfterReq.ElementValueInfo.FieldNum_ = 0
	pkg.Body.ListAddAfterReq.ElementValueInfo.CompactValueSet.FieldIndexNum = 0
	pkg.Body.ListAddAfterReq.ElementValueInfo.CompactValueSet.FieldIndexs = nil
	pkg.Body.ListAddAfterReq.ElementValueInfo.CompactValueSet.ValueBufLen = 0
	pkg.Body.ListAddAfterReq.ElementValueInfo.CompactValueSet.ValueBuf = nil
	pkg.Body.ListAddAfterReq.ShiftFlag = byte(tcaplus_protocol_cs.TCAPLUS_LIST_SHIFT_HEAD)
	pkg.Body.ListAddAfterReq.ElementIndex = -1
	pkg.Body.ListAddAfterReq.CheckVersiontType = policy.CheckDataVersionAutoIncrease
	pkg.Body.ListAddAfterReq.Flag = 0
	req := &listAddAfterRequest{
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

func (req *listAddAfterRequest) AddRecord(index int32) (*record.Record, error) {
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

	//key value set
	rec.ShardingKey = &req.pkg.Head.SplitTableKeyBuff
	rec.ShardingKeyLen = &req.pkg.Head.SplitTableKeyBuffLen
	rec.KeySet = req.pkg.Head.KeyInfo
	rec.ValueSet = req.pkg.Body.ListAddAfterReq.ElementValueInfo
	req.pkg.Body.ListAddAfterReq.ElementIndex = index
	req.record = rec
	return rec, nil
}

func (req *listAddAfterRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listAddAfterRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return &terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.ListAddAfterReq.CheckVersiontType = p
	return nil
}

func (req *listAddAfterRequest) SetResultFlag(flag int) error {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "ResultFlag invalid"}
	}
	req.pkg.Body.ListAddAfterReq.Flag = byte(flag)
	return nil
}

func (req *listAddAfterRequest) Pack() ([]byte, error) {
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

	if err := req.record.PackValue(nil); err != nil {
		logger.ERR("record pack value failed, %s", err.Error())
		return nil, err
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
		logger.DEBUG("%s", common.CovertToJson(req.pkg.Body.ListAddAfterReq))
	}

	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("listAddAfterRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listAddAfterRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listAddAfterRequest) GetKeyHash() (uint32, error) {
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

func (req *listAddAfterRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *listAddAfterRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listAddAfterRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listAddAfterRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *listAddAfterRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listAddAfterRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listAddAfterRequest) SetListShiftFlag(shiftFlag byte) int32 {
	req.pkg.Body.ListAddAfterReq.ShiftFlag = shiftFlag
	return int32(terror.GEN_ERR_SUC)
}

func (req *listAddAfterRequest) SetResultFlagForSuccess(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	flag = flag << 4
	flag |= 1 << 6
	req.pkg.Body.ListAddAfterReq.Flag |= flag
	return terror.GEN_ERR_SUC
}

func (req *listAddAfterRequest) SetResultFlagForFail(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	flag = flag << 2
	flag |= 1 << 6
	req.pkg.Body.ListAddAfterReq.Flag |= flag
	return terror.GEN_ERR_SUC
}

func (req *listAddAfterRequest) SetPerfTest(sendTime uint64) int {
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

func (req *listAddAfterRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *listAddAfterRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *listAddAfterRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
