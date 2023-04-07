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

type listDeleteRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	isPB      bool
}

func newListDeleteRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*listDeleteRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListDeleteReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.ListDeleteReq.Flag = 0
	pkg.Body.ListDeleteReq.ElementIndex = 0
	pkg.Body.ListDeleteReq.CheckVersiontType = 1
	pkg.Body.ListDeleteReq.Condition = ""
	req := &listDeleteRequest{
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

func (req *listDeleteRequest) AddRecord(index int32) (*record.Record, error) {
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
	req.pkg.Body.ListDeleteReq.ElementIndex = index
	rec.Condition = &req.pkg.Body.ListDeleteReq.Condition
	req.record = rec
	return rec, nil
}

func (req *listDeleteRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listDeleteRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return &terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.ListDeleteReq.CheckVersiontType = p
	return nil
}

func (req *listDeleteRequest) SetResultFlag(flag int) error {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "ResultFlag invalid"}
	}
	req.pkg.Body.ListDeleteReq.Flag = byte(flag)
	return nil
}

func (req *listDeleteRequest) Pack() ([]byte, error) {
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
		logger.DEBUG("%s", common.CovertToJson(req.pkg.Body.ListDeleteReq))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("listDeleteRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listDeleteRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listDeleteRequest) GetKeyHash() (uint32, error) {
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

func (req *listDeleteRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *listDeleteRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listDeleteRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listDeleteRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *listDeleteRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listDeleteRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listDeleteRequest) SetResultFlagForSuccess(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	flag = flag << 4
	flag |= 1 << 6
	req.pkg.Body.ListDeleteReq.Flag |= flag
	return terror.GEN_ERR_SUC
}

func (req *listDeleteRequest) SetResultFlagForFail(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	flag = flag << 2
	flag |= 1 << 6
	req.pkg.Body.ListDeleteReq.Flag |= flag
	return terror.GEN_ERR_SUC
}

func (req *listDeleteRequest) SetPerfTest(sendTime uint64) int {
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

func (req *listDeleteRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *listDeleteRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *listDeleteRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
