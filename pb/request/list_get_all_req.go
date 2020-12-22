package request

import (
	"git.code.com/gcloud_storage_group/tcaplus-go-api/common"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type listGetAllRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       *record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
}

func newListGetAllRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*listGetAllRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListGetAllReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Head.KeyInfo.FieldNum = 0
	req := &listGetAllRequest{
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

func (req *listGetAllRequest) AddRecord(index int32) (*record.Record, error) {
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

	rec.KeySet = req.pkg.Head.KeyInfo
	//rec.ValueSet = req.pkg.Body.ListGetAllReq.ElementValueNames
	req.record = rec
	return rec, nil
}

func (req *listGetAllRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listGetAllRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "listGetAll not Support VersionPolicy"}
}

func (req *listGetAllRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "listGetAll not Support ResultFlag"}
}

func (req *listGetAllRequest) Pack() ([]byte, error) {
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

	for key, _ := range req.record.ValueMap {
		req.pkg.Body.ListGetAllReq.ElementValueNames.FieldNum += 1
		req.pkg.Body.ListGetAllReq.ElementValueNames.FieldName =
			append(req.pkg.Body.ListGetAllReq.ElementValueNames.FieldName, key)
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

func (req *listGetAllRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listGetAllRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *listGetAllRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *listGetAllRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listGetAllRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listGetAllRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}
func (req *listGetAllRequest)SetResultLimit(limit int32, offset int32) int32 {
	req.pkg.Body.ListGetAllReq.StartSubscript = offset
	req.pkg.Body.ListGetAllReq.ElementNum = limit
	return int32(terror.GEN_ERR_SUC)
}

func (req *listGetAllRequest)SetAddableIncreaseFlag(increase_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *listGetAllRequest)SetMultiResponseFlag(multi_flag byte) int32{
	if 1 == multi_flag {
		req.pkg.Body.ListGetAllReq.AllowMultiResponses = 1
	} else {
		req.pkg.Body.ListGetAllReq.AllowMultiResponses = 0
	}
	return int32(terror.GEN_ERR_SUC)
}

func (req *listGetAllRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *listGetAllRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
