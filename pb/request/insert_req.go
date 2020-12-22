package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type insertRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
}

func newInsertRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*insertRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.InsertReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.InsertReq.ValueInfo.EncodeType = 1
	req := &insertRequest{
		appId:     appId,
		zoneId:    zoneId,
		tableName: tableName,
		cmd:       cmd,
		seq:       seq,
		record:    nil,
		pkg:       pkg,
	}
	return req, nil
}

func (req *insertRequest) AddRecord(index int32) (*record.Record, error) {
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
	rec.ShardingKey = &req.pkg.Head.SplitTableKeyBuff
	rec.ShardingKeyLen = &req.pkg.Head.SplitTableKeyBuffLen
	rec.KeySet = req.pkg.Head.KeyInfo
	rec.ValueSet = req.pkg.Body.InsertReq.ValueInfo
	req.record = rec
	return rec, nil
}

func (req *insertRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *insertRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Insert not Support VersionPolicy"}
}

func (req *insertRequest) SetResultFlag(flag int) error {
	if flag != 0 && flag != 1 && flag != 2 && flag != 3 {
		logger.ERR("result flag invalid %d.", flag)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "ResultFlag invalid"}
	}
	req.pkg.Body.InsertReq.Flag = byte(flag)
	return nil
}

func (req *insertRequest) Pack() ([]byte, error) {
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

	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("insertRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *insertRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *insertRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *insertRequest) SetFieldNames(valueNameList []string) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "insert not Support SetFieldNames"}
}

func (req *insertRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *insertRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *insertRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *insertRequest)SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *insertRequest)SetAddableIncreaseFlag(increase_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *insertRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *insertRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *insertRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
