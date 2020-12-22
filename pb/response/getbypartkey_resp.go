package response

import (
	//"bytes"
	//"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	//"runtime/debug"
)

type getByPartKeyResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	offset int32
	idx    int32
}

func newGetByPartKeyResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*getByPartKeyResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.GetByPartKeyRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &getByPartKeyResponse{pkg: pkg}, nil
}

func (res *getByPartKeyResponse) GetResult() int {
	ret := int(res.pkg.Body.GetByPartKeyRes.Result)
	return ret
}

func (res *getByPartKeyResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *getByPartKeyResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *getByPartKeyResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *getByPartKeyResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *getByPartKeyResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *getByPartKeyResponse) GetRecordCount() int {
	if res.pkg.Body.GetByPartKeyRes.Result == 0 ||
		res.pkg.Body.GetByPartKeyRes.Result == int32(terror.COMMON_INFO_DATA_NOT_MODIFIED) {
		Result := res.pkg.Body.GetByPartKeyRes.RecordResult.Result
		//logger.DEBUG(" ********result： %d, num: %d, OffSet:%d, RecordNum: %d, BatchValueLen:%d, IsCompleteFlag:%d，",
		//	res.pkg.Body.GetByPartKeyRes.Result, Result.TotalNum, Result.OffSet, Result.RecordNum, Result.BatchValueLen,
		//	Result.IsCompleteFlag/*, debug.Stack()*/)
		return int(Result.RecordNum)
	}
	return 0
}

func (res *getByPartKeyResponse) FetchRecord() (*record.Record, error) {
	data := res.pkg.Body.GetByPartKeyRes.RecordResult.Result
	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	if res.idx >= int32(data.RecordNum) || res.offset >= data.BatchValueLen{
		logger.ERR("resp fetch record over, current idx: %d, ",res.idx)
		return nil , &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}

	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen-1]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Head.KeyInfo,
		ValueSet: nil,
		UpdFieldSet: nil,
	}

	read_bytes, err := unpack_record(data.BatchValueInfo[res.offset: data.BatchValueLen],
		data.BatchValueLen - res.offset,  rec.KeyMap, rec.ValueMap)
	if err !=nil{
		logger.ERR("record unpack failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}

	res.idx += 1
	res.offset += read_bytes
	res.record = rec
	return rec, nil
}

func (res *getByPartKeyResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *getByPartKeyResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *getByPartKeyResponse) HaveMoreResPkgs() int {
	if 0 != res.pkg.Body.GetByPartKeyRes.Result ||
		1 == res.pkg.Body.GetByPartKeyRes.RecordResult.Result.IsCompleteFlag  {
		return 0
	} else{
		return 1
	}
}

func (res *getByPartKeyResponse) GetTotalNum() int {
	return int(res.pkg.Body.GetByPartKeyRes.RecordResult.Result.TotalNum)
}

func (res *getByPartKeyResponse) GetFailedNum() int {
	return 0
}

func (res *getByPartKeyResponse) FetchErrorRecord() (*record.Record, error) {
	return nil,nil
}

func (res *getByPartKeyResponse) GetRecordMatchCount() int{
	if 0 == res.pkg.Body.GetByPartKeyRes.Result{
		return int(res.pkg.Body.GetByPartKeyRes.RecordResult.Result.TotalNum)
	}
	return terror.GEN_ERR_ERR;
}