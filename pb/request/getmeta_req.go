package request

import (
"git.code.com/gcloud_storage_group/tcaplus-go-api/common"
"git.code.com/gcloud_storage_group/tcaplus-go-api/logger"
"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
"git.code.com/gcloud_storage_group/tcaplus-go-api/record"
"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type getMetaRequest struct {
	appId        uint64
	zoneId       uint32
	tableName    string
	cmd          int
	seq          uint32
	pkg          *tcaplus_protocol_cs.TCaplusPkg
}

func newGetMetaRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*getMetaRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.MetadataGetReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.MetadataGetReq.MetadataVersion = 0
	req := &getMetaRequest{
		appId:        appId,
		zoneId:       zoneId,
		tableName:    tableName,
		cmd:          cmd,
		seq:          seq,
		pkg:          pkg,
	}
	return req, nil
}

func (req *getMetaRequest) AddRecord(index int32) (*record.Record, error) {
	return nil, nil
}

func (req *getMetaRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *getMetaRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support VersionPolicy"}
}

func (req *getMetaRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support ResultFlag"}
}

func (req *getMetaRequest) Pack() ([]byte, error) {
	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("getMetaRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *getMetaRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *getMetaRequest) GetKeyHash() (uint32, error) {
	return 1, nil
}

func (req *getMetaRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *getMetaRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *getMetaRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *getMetaRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *getMetaRequest)SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *getMetaRequest)SetAddableIncreaseFlag(increase_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *getMetaRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.GEN_ERR_SUC)
}

func (req *getMetaRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *getMetaRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *getMetaRequest)SetSplitTableKeyBuff(splitkey []byte) int {
	req.pkg.Head.SplitTableKeyBuff = splitkey
	req.pkg.Head.SplitTableKeyBuffLen = uint32(len(splitkey))
	return terror.GEN_ERR_SUC
}

