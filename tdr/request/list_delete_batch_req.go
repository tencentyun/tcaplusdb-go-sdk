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

type listDeleteBatchRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	isPB      bool
}

func newListDeleteBatchRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*listDeleteBatchRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListDeleteBatchReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.ListDeleteBatchReq.Flag = 0
	pkg.Body.ListDeleteBatchReq.AllowMultiResponses = 0
	pkg.Body.ListDeleteBatchReq.ElementNum = 0
	pkg.Body.ListDeleteBatchReq.ElementIndexArray = nil
	pkg.Body.ListDeleteBatchReq.CheckVersiontType = policy.CheckDataVersionAutoIncrease
	req := &listDeleteBatchRequest{
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

func (req *listDeleteBatchRequest) AddRecord(index int32) (*record.Record, error) {
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

func (req *listDeleteBatchRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listDeleteBatchRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return &terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.ListDeleteBatchReq.CheckVersiontType = p
	return nil
}

func (req *listDeleteBatchRequest) SetResultFlag(flag int) error {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "ResultFlag invalid"}
	}
	req.pkg.Body.ListDeleteBatchReq.Flag = byte(flag)
	return nil
}

func (req *listDeleteBatchRequest) Pack() ([]byte, error) {
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

	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("getRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listDeleteBatchRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listDeleteBatchRequest) GetKeyHash() (uint32, error) {
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

func (req *listDeleteBatchRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *listDeleteBatchRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listDeleteBatchRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listDeleteBatchRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *listDeleteBatchRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listDeleteBatchRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	if 1 == multi_flag {
		req.pkg.Body.ListDeleteBatchReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.ListDeleteBatchReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *listDeleteBatchRequest) SetResultFlagForSuccess(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	req.pkg.Body.ListDeleteBatchReq.Flag = flag << 4
	req.pkg.Body.ListDeleteBatchReq.Flag |= 1 << 6
	return terror.GEN_ERR_SUC
}

func (req *listDeleteBatchRequest) SetResultFlagForFail(flag byte) int {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return terror.ParameterInvalid
	}
	// 0(1个bit位) | 本版本开始该位设置为1(1个bit位) | 成功时的标识(2个bit位) | 失败时的标识(2个bit位) | 本版本以前的标识(2个bit位)
	req.pkg.Body.ListDeleteBatchReq.Flag = flag << 2
	req.pkg.Body.ListDeleteBatchReq.Flag |= 1 << 6
	return terror.GEN_ERR_SUC
}

func (req *listDeleteBatchRequest) AddElementIndex(idx int32) int32 {
	bodyreq := req.pkg.Body.ListDeleteBatchReq
	for i := int32(0); i < bodyreq.ElementNum; i++ {
		if idx == bodyreq.ElementIndexArray[i] {
			return int32(terror.GEN_ERR_SUC)
		}
	}
	if int64(bodyreq.ElementNum) > tcaplus_protocol_cs.TCAPLUS_MAX_LIST_ELEMENTS_NUM {
		return int32(terror.API_ERR_OVER_MAX_LIST_INDEX_NUM)
	}
	bodyreq.ElementIndexArray = append(bodyreq.ElementIndexArray, idx)
	bodyreq.ElementNum++
	return int32(terror.GEN_ERR_SUC)
}

func (req *listDeleteBatchRequest) SetPerfTest(sendTime uint64) int {
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

func (req *listDeleteBatchRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *listDeleteBatchRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *listDeleteBatchRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
