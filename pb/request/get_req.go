package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type getRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       *record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
}

func newGetRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*getRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.GetReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.GetReq.ValueInfo.EncodeType = 1
	req := &getRequest{
		appId:        appId,
		zoneId:       zoneId,
		tableName:    tableName,
		cmd:          cmd,
		seq:          seq,
		record:       nil,
		pkg:          pkg,
		valueNameMap: make(map[string]bool),
	}
	return req, nil
}

func (req *getRequest) AddRecord(index int32) (*record.Record, error) {
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
	rec.ValueSet = req.pkg.Body.GetReq.ValueInfo
	req.record = rec
	return rec, nil
}

func (req *getRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *getRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support VersionPolicy"}
}

func (req *getRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support ResultFlag"}
}

func (req *getRequest) Pack() ([]byte, error) {
	if req.record == nil {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	if err := req.record.PackKey(); err != nil {
		logger.ERR("record pack key failed, %s", err.Error())
		return nil, err
	}

	if len(req.valueNameMap) > 0 {
		for name, _ := range req.valueNameMap {
			req.record.ValueMap[name] = []byte{}
		}
	}

	if err := req.record.PackValue(req.valueNameMap); err != nil {
		logger.ERR("record pack value failed, %s", err.Error())
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

func (req *getRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *getRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *getRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *getRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *getRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *getRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *getRequest)SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *getRequest)SetAddableIncreaseFlag(increase_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *getRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *getRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *getRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
