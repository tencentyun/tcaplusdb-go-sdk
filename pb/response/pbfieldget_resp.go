package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type pbFieldGetResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
}

func newPbFieldGetResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*pbFieldGetResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusPbFieldGetRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &pbFieldGetResponse{pkg: pkg}, nil
}

func (res *pbFieldGetResponse) GetResult() int {
	ret := int(res.pkg.Body.TCaplusPbFieldGetRes.Result)
	return ret
}

func (res *pbFieldGetResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *pbFieldGetResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *pbFieldGetResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *pbFieldGetResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *pbFieldGetResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *pbFieldGetResponse) GetRecordCount() int {
	if res.pkg.Body.TCaplusPbFieldGetRes.Result == 0 || res.pkg.Body.TCaplusPbFieldGetRes.Result == int32(terror.COMMON_INFO_DATA_NOT_MODIFIED) {
		return 1
	}
	return 0
}

func (res *pbFieldGetResponse) FetchRecord() (*record.Record, error) {
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
		PBValueSet:  res.pkg.Body.TCaplusPbFieldGetRes.ResultInfo,
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

func (res *pbFieldGetResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *pbFieldGetResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *pbFieldGetResponse) HaveMoreResPkgs() int {
	return 0
}
func (res *pbFieldGetResponse) GetTotalNum() int {
	return 0
}
func (res *pbFieldGetResponse) GetFailedNum() int {
	return 0
}

func (res *pbFieldGetResponse) FetchErrorRecord() (*record.Record, error) {
	return nil,nil
}
func (res *pbFieldGetResponse) GetRecordMatchCount() int{
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}