package response

import (
	"git.code.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type pbFieldUpdateResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
}

func newPbFieldUpdateResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*pbFieldUpdateResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusPbFieldUpdateRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &pbFieldUpdateResponse{pkg: pkg}, nil
}

func (res *pbFieldUpdateResponse) GetResult() int {
	ret := int(res.pkg.Body.TCaplusPbFieldUpdateRes.Result)
	return ret
}

func (res *pbFieldUpdateResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *pbFieldUpdateResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *pbFieldUpdateResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *pbFieldUpdateResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *pbFieldUpdateResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *pbFieldUpdateResponse) GetRecordCount() int {
	if res.pkg.Body.TCaplusPbFieldUpdateRes.Result == 0 || res.pkg.Body.TCaplusPbFieldUpdateRes.Result == int32(terror.COMMON_INFO_DATA_NOT_MODIFIED) {
		return 1
	}
	return 0
}

func (res *pbFieldUpdateResponse) FetchRecord() (*record.Record, error) {
	if res.record != nil {
		logger.ERR("all record fetched , no more")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}

	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}

	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[:res.pkg.Head.RouterInfo.TableNameLen-1]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Head.KeyInfo,
		PBValueSet:  res.pkg.Body.TCaplusPbFieldUpdateRes.ResultInfo,
		UpdFieldSet: nil,
	}

	//unpack
	if err := rec.UnPackKey(); err != nil {
		logger.ERR("record unpack key failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}

	if err := rec.UnPackPBValue(); err != nil {
		logger.ERR("record unpack value failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}
	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec
	return rec, nil
}

func (res *pbFieldUpdateResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *pbFieldUpdateResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *pbFieldUpdateResponse) HaveMoreResPkgs() int {
	return 0
}
func (res *pbFieldUpdateResponse) GetTotalNum() int {
	return 0
}
func (res *pbFieldUpdateResponse) GetFailedNum() int {
	return 0
}

func (res *pbFieldUpdateResponse) FetchErrorRecord() (*record.Record, error) {
	return nil,nil
}
func (res *pbFieldUpdateResponse) GetRecordMatchCount() int{
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}