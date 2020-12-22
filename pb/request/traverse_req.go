package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type traverseRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       *record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
}

func newTraverseRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*traverseRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TableTraverseReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	req := &traverseRequest{
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

func (req *traverseRequest) AddRecord(index int32) (*record.Record, error) {
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
	rec.NameSet = req.pkg.Body.TableTraverseReq.ValueInfo
	req.record = rec
	return rec, nil
}

func (req *traverseRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *traverseRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support VersionPolicy"}
}

func (req *traverseRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support ResultFlag"}
}

func (req *traverseRequest) Pack() ([]byte, error) {
	req.pkg.Body.TableTraverseReq.ValueInfo.FieldNum = 3
	req.pkg.Body.TableTraverseReq.ValueInfo.FieldName = []string{"klen", "vlen", "value"}

	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("traverseRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *traverseRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *traverseRequest) GetKeyHash() (uint32, error) {
	return 5, nil
}

func (req *traverseRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *traverseRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *traverseRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *traverseRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *traverseRequest)SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *traverseRequest)SetAddableIncreaseFlag(increase_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *traverseRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *traverseRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *traverseRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *traverseRequest) GetTcaplusPackagePtr() *tcaplus_protocol_cs.TCaplusPkg {
	return req.pkg
}
