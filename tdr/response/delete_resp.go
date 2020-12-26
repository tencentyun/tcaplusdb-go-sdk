package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type deleteResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
}

func newDeleteResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*deleteResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.DeleteRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &deleteResponse{pkg: pkg}, nil
}

func (res *deleteResponse) GetResult() int {
	ret := int(res.pkg.Body.DeleteRes.Result)
	return ret
}

func (res *deleteResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen])
	return tableName
}

func (res *deleteResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *deleteResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *deleteResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *deleteResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *deleteResponse) GetRecordCount() int {
	if 0 != (res.pkg.Body.DeleteRes.Flag & (1 << 6)) {
		//新版本的result flag 通过ResultFlag判断
		if res.pkg.Body.DeleteRes.Result == 0 {
			ret := GetResultFlagByBit(res.pkg.Body.DeleteRes.Flag, true)
			if tcaplus_protocol_cs.TCaplusValueFlag_ALLOLDVALUE == ret {
				return 1
			}
		} else {
			ret := GetResultFlagByBit(res.pkg.Body.DeleteRes.Flag, false)
			if tcaplus_protocol_cs.TCaplusValueFlag_ALLOLDVALUE == ret &&
				res.pkg.Body.DeleteRes.ResultInfo.CompactValueSet.FieldIndexNum > 0 {
				return 1
			}
		}
	} else {
		//老版本的result flag 通过ResultFlag判断
		if (res.pkg.Body.DeleteRes.Flag == 1 || res.pkg.Body.DeleteRes.Flag == 2 ||
			(res.pkg.Body.DeleteRes.Flag == 3 && res.pkg.Body.DeleteRes.ResultInfo.CompactValueSet.FieldIndexNum > 0)) &&
			(res.pkg.Body.DeleteRes.Result == 0) {
			return 1
		}
	}
	return 0
}

func (res *deleteResponse) FetchRecord() (*record.Record, error) {
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
		ValueSet:    res.pkg.Body.DeleteRes.ResultInfo,
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

func (res *deleteResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *deleteResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *deleteResponse) HaveMoreResPkgs() int {
	return 0
}
func (res *deleteResponse) GetTotalNum() int {
	return 0
}
func (res *deleteResponse) GetFailedNum() int {
	return 0
}
func (res *deleteResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}
func (res *deleteResponse) GetRecordMatchCount() int {
	return int(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}
