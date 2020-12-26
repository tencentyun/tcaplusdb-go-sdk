package response

import (
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type replaceResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
}

func newReplaceResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*replaceResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ReplaceRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &replaceResponse{pkg: pkg}, nil
}

func (res *replaceResponse) GetResult() int {
	ret := int(res.pkg.Body.ReplaceRes.Result)
	return ret
}

func (res *replaceResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen])
	return tableName
}

func (res *replaceResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *replaceResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *replaceResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *replaceResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *replaceResponse) GetRecordCount() int {
	if 0 != (res.pkg.Body.ReplaceRes.Flag & (1 << 6)) {
		//新版本的result flag 通过ResultFlag判断
		if res.pkg.Body.ReplaceRes.Result == 0 {
			ret := GetResultFlagByBit(res.pkg.Body.ReplaceRes.Flag, true)
			if tcaplus_protocol_cs.TCaplusValueFlag_SAMEWITHREQUEST == ret ||
				tcaplus_protocol_cs.TCaplusValueFlag_ALLVALUE == ret ||
				(tcaplus_protocol_cs.TCaplusValueFlag_ALLOLDVALUE == ret &&
					res.pkg.Body.ReplaceRes.ResultInfo.CompactValueSet.FieldIndexNum > 0) {
				return 1
			}
		} else {
			ret := GetResultFlagByBit(res.pkg.Body.ReplaceRes.Flag, false)
			if tcaplus_protocol_cs.TCaplusValueFlag_SAMEWITHREQUEST == ret ||
				(tcaplus_protocol_cs.TCaplusValueFlag_ALLOLDVALUE == ret &&
					res.pkg.Body.ReplaceRes.ResultInfo.CompactValueSet.FieldIndexNum > 0) {
				return 1
			}
		}
	} else {
		//老版本的result flag 通过ResultFlag判断
		if (res.pkg.Body.ReplaceRes.Flag == 1 || res.pkg.Body.ReplaceRes.Flag == 2 ||
			(res.pkg.Body.ReplaceRes.Flag == 3 &&
				res.pkg.Body.ReplaceRes.ResultInfo.CompactValueSet.FieldIndexNum > 0)) &&
			(res.pkg.Body.ReplaceRes.Result == 0 ||
				res.pkg.Body.ReplaceRes.Result == int32(terror.SVR_ERR_FAIL_INVALID_VERSION)) {
			return 1
		}
	}
	return 0
}

func (res *replaceResponse) FetchRecord() (*record.Record, error) {
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
		ValueSet:    res.pkg.Body.ReplaceRes.ResultInfo,
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

func (res *replaceResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *replaceResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *replaceResponse) HaveMoreResPkgs() int {
	return 0
}

func (res *replaceResponse) GetTotalNum() int {
	return 0
}

func (res *replaceResponse) GetFailedNum() int {
	return 0
}

func (res *replaceResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}

func (res *replaceResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}
