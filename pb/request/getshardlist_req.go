package request

import (
	"git.code.com/gcloud_storage_group/tcaplus-go-api/common"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type getShardListRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       *record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
}

func newGetShardListRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*getShardListRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.GetShardListReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	req := &getShardListRequest{
		appId:        appId,
		zoneId:       zoneId,
		tableName:    tableName,
		cmd:          cmd,
		seq:          seq,
		record:       nil,
		pkg:          pkg,
	}
	return req, nil
}

func (req *getShardListRequest) AddRecord(index int32) (*record.Record, error) {
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
	req.record = rec
	return rec, nil
}

func (req *getShardListRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *getShardListRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support VersionPolicy"}
}

func (req *getShardListRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support ResultFlag"}
}

func (req *getShardListRequest) Pack() ([]byte, error) {
	//if req.record == nil {
	//	return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	//}

	//if err := req.record.PackKey(); err != nil {
	//	logger.ERR("record pack key failed, %s", err.Error())
	//	return nil, err
	//}
	//
	//if len(req.valueNameMap) > 0 {
	//	for name, _ := range req.valueNameMap {
	//		req.record.ValueMap[name] = []byte{}
	//	}
	//}
	//
	//if err := req.record.PackValue(req.valueNameMap); err != nil {
	//	logger.ERR("record pack value failed, %s", err.Error())
	//	return nil, err
	//}

	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("getShardListRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *getShardListRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *getShardListRequest) GetKeyHash() (uint32, error) {
	//if req.record == nil {
	//	return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	//}
	//return keyHashCode(req.pkg.Head.KeyInfo)
	return 5, nil
}

func (req *getShardListRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *getShardListRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *getShardListRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *getShardListRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *getShardListRequest)SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *getShardListRequest)SetAddableIncreaseFlag(increase_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *getShardListRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *getShardListRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *getShardListRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *getShardListRequest) GetTcaplusPackagePtr() *tcaplus_protocol_cs.TCaplusPkg {
	return req.pkg
}
