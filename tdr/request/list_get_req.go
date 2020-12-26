package request

import (
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/common"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type listGetRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	record       *record.Record
	pkg          *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
}

func newListGetRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*listGetRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListGetReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	pkg.Body.ListGetReq.ElementValueNames.FieldNum = 0
	req := &listGetRequest{
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

func (req *listGetRequest) AddRecord(index int32) (*record.Record, error) {
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
	req.pkg.Body.ListGetReq.ElementIndex = index
	req.record = rec
	return rec, nil
}

func (req *listGetRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *listGetRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH,
		Message: "list get not support SetVersionPolicy"}
}

func (req *listGetRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH,
		Message: "list get not support SetResultFlag"}
}

func (req *listGetRequest) Pack() ([]byte, error) {
	if req.record == nil {
		return nil, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}

	if err := req.record.PackKey(); err != nil {
		logger.ERR("record pack key failed, %s", err.Error())
		return nil, err
	}

	if len(req.valueNameMap) > 0 {
		req.record.ValueMap = make(map[string][]byte)
		for name, _ := range req.valueNameMap {
			req.record.ValueMap[name] = []byte{}
		}
	}

	for key, _ := range req.record.ValueMap {
		req.pkg.Body.ListGetReq.ElementValueNames.FieldNum += 1
		req.pkg.Body.ListGetReq.ElementValueNames.FieldName =
			append(req.pkg.Body.ListGetReq.ElementValueNames.FieldName, key)
	}

	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("listGetRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *listGetRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *listGetRequest) GetKeyHash() (uint32, error) {
	if req.record == nil {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoRecord}
	}
	return keyHashCode(req.pkg.Head.KeyInfo)
}

func (req *listGetRequest) SetFieldNames(valueNameList []string) error {
	for _, v := range valueNameList {
		req.valueNameMap[v] = true
	}
	return nil
}

func (req *listGetRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *listGetRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *listGetRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *listGetRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listGetRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *listGetRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *listGetRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}
