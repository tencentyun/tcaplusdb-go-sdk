package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type batchGetRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       []*record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
	idx          int
}

func newBatchGetRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*batchGetRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.BatchGetReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Head.KeyInfo.FieldNum = 0
	req := &batchGetRequest{
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

func (req *batchGetRequest) AddRecord(index int32) (*record.Record, error) {
	//batchGetReq := req.pkg.Body.BatchGetReq
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
		SplitTableKeyBuff: nil,
	}

	rec.KeySet = new(tcaplus_protocol_cs.TCaplusKeySet)

	rec.SplitTableKeyBuff = new(tcaplus_protocol_cs.SplitTableKeyBuff)
	req.pkg.Body.BatchGetReq.SplitTableKeyBuffs = append(req.pkg.Body.BatchGetReq.SplitTableKeyBuffs,
		rec.SplitTableKeyBuff)

	req.record = append(req.record, rec)
	return rec, nil
}

func (req *batchGetRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *batchGetRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "BatchGet not Support VersionPolicy"}
}

func (req *batchGetRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "BatchGet not Support ResultFlag"}
}

func (req *batchGetRequest) Pack() ([]byte, error) {
	if len(req.record) == 0 {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	for _, rec := range req.record {
		if err := rec.PackKey(); err != nil {
			logger.ERR("record pack key failed, %s", err.Error())
			return nil, err
		}
		req.pkg.Body.BatchGetReq.RecordNum += 1
		req.pkg.Body.BatchGetReq.KeyInfo = append(req.pkg.Body.BatchGetReq.KeyInfo, rec.KeySet)
	}

	req.pkg.Body.BatchGetReq.ValueInfo.FieldNum = 3
	req.pkg.Body.BatchGetReq.ValueInfo.FieldName = []string{"klen", "vlen", "value"}

	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("batchGetRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *batchGetRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *batchGetRequest) GetKeyHash() (uint32, error) {
	if len(req.record) == 0 {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.record[0].KeySet)
}

func (req *batchGetRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *batchGetRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *batchGetRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *batchGetRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *batchGetRequest)SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *batchGetRequest)SetAddableIncreaseFlag(increase_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *batchGetRequest)SetMultiResponseFlag(multi_flag byte) int32{
	if 0 != multi_flag  {
		req.pkg.Body.BatchGetReq.AllowMultiResponses = 1
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *batchGetRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *batchGetRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
