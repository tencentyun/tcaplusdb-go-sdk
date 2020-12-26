package response

import (
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type getResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
}

func newGetResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*getResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.GetRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &getResponse{pkg: pkg}, nil
}

func (res *getResponse) GetResult() int {
	ret := int(res.pkg.Body.GetRes.Result)
	return ret
}

func (res *getResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen])
	return tableName
}

func (res *getResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *getResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *getResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *getResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *getResponse) GetRecordCount() int {
	if res.pkg.Body.GetRes.Result == 0 || res.pkg.Body.GetRes.Result == int32(terror.COMMON_INFO_DATA_NOT_MODIFIED) {
		return 1
	}
	return 0
}

func (res *getResponse) FetchRecord() (*record.Record, error) {
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
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Head.KeyInfo,
		ValueSet:    res.pkg.Body.GetRes.ResultInfo,
		UpdFieldSet: nil,
	}

	//unpack
	if err := rec.UnPackKey(); err != nil {
		logger.ERR("record unpack key failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}

	if err := rec.UnPackValue(); err != nil {
		logger.ERR("record unpack value failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}
	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec
	return rec, nil
}

func (res *getResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *getResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *getResponse) HaveMoreResPkgs() int {
	return 0
}
func (res *getResponse) GetTotalNum() int {
	return 0
}
func (res *getResponse) GetFailedNum() int {
	return 0
}

func (res *getResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}
func (res *getResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}
