package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type increaseRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
}

func newIncreaseRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*increaseRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.IncreaseReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	req := &increaseRequest{
		appId:     appId,
		zoneId:    zoneId,
		tableName: tableName,
		cmd:       cmd,
		seq:       seq,
		record:    nil,
		pkg:       pkg,
		valueNameMap: make(map[string]bool),
	}
	return req, nil
}

func (req *increaseRequest) AddRecord(index int32) (*record.Record, error) {
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
	}

	//key value set
	rec.KeySet = req.pkg.Head.KeyInfo
	rec.UpdFieldSet = req.pkg.Body.IncreaseReq.ValueInfo
	req.record = rec
	return rec, nil
}

func (req *increaseRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *increaseRequest) SetVersionPolicy(p uint8) error {
	if p != policy.CheckDataVersionAutoIncrease && p != policy.NoCheckDataVersionAutoIncrease &&
		p != policy.NoCheckDataVersionOverwrite {
		logger.ERR("policy type Invalid %d", p)
		return terror.ErrorCode{Code: terror.InvalidPolicy}
	}
	req.pkg.Body.IncreaseReq.CheckVersiontType = p
	return nil
}

func (req *increaseRequest) SetResultFlag(flag int) error {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "ResultFlag invalid"}
	}
	req.pkg.Body.IncreaseReq.Flag = byte(flag)
	return nil
}

func (req *increaseRequest) Pack() ([]byte, error) {
	if req.record == nil {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	if err := req.record.PackKey(); err != nil {
		logger.ERR("record pack key failed, %s", err.Error())
		return nil, err
	}

	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("increaseRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *increaseRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *increaseRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *increaseRequest) SetFieldNames(valueNameList []string) error {
	for _, name := range valueNameList {
		req.valueNameMap[name] = true
	}
	return nil
}

func (req *increaseRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *increaseRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *increaseRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}
func (req *increaseRequest)SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *increaseRequest)SetAddableIncreaseFlag(increase_flag byte) int32{
	req.pkg.Body.IncreaseReq.AddableIncrease = increase_flag
	return int32(terror.GEN_ERR_SUC)
}

func (req *increaseRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *increaseRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *increaseRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
