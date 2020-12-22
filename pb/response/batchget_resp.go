package response

import (
	"git.code.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type batchGetResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	offset int32
	idx    int32
}

func newBatchGetResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*batchGetResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.BatchGetRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &batchGetResponse{pkg: pkg}, nil
}

func (res *batchGetResponse) GetResult() int {
	ret := int(res.pkg.Body.BatchGetRes.Result)
	return ret
}

func (res *batchGetResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *batchGetResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *batchGetResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *batchGetResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *batchGetResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *batchGetResponse) GetRecordCount() int {
	if res.pkg.Body.BatchGetRes.Result == 0 ||
		res.pkg.Body.BatchGetRes.Result == int32(terror.COMMON_INFO_DATA_NOT_MODIFIED) {
		return int(res.pkg.Body.BatchGetRes.RecordNum)
	}
	return 0
}

func (res *batchGetResponse) FetchRecord() (*record.Record, error) {

	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	data := res.pkg.Body.BatchGetRes
	if res.idx >= int32(data.RecordNum) || res.offset >= data.BatchValueLen{
		logger.ERR("resp fetch record over, current idx: %d, ",res.idx)
		return nil , &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	logger.DEBUG("read bytes: %d, total bytes: %d", res.offset, data.BatchValueLen)
	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen-1]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Head.KeyInfo,
		ValueSet:    nil,//res.pkg.Body.BatchGetRes,
		UpdFieldSet: nil,
	}

	//unpack
	read_bytes, err := unpack_record_k_v(data.BatchValueInfo[res.offset: data.BatchValueLen],
		data.BatchValueLen - res.offset,  rec.KeyMap, rec.ValueMap)
	if err != nil{
		logger.ERR("record unpack failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}
	logger.DEBUG("record unpack success, key: %+v, value: %+v",  rec.KeyMap,  rec.ValueMap)
	res.idx += 1
	res.offset += read_bytes

	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec

	if ret := int(data.RecordResult[res.idx - 1 ]); ret != terror.GEN_ERR_SUC{
		return rec, &terror.ErrorCode{Code:  ret}
	}
	return rec, nil
}

func (res *batchGetResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *batchGetResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *batchGetResponse) HaveMoreResPkgs() int {
	if 0 == res.pkg.Body.BatchGetRes.Result && 0 != res.pkg.Body.BatchGetRes.LeftNum{
		return 1
	}else {
		return 0
	}
}

func (res *batchGetResponse) GetTotalNum() int {
	return int(res.pkg.Body.BatchGetRes.TotalNum)
}

func (res *batchGetResponse) GetFailedNum() int {
	return 0
}

func (res *batchGetResponse) FetchErrorRecord() (*record.Record, error) {
	return nil,nil
}

func (res *batchGetResponse) GetRecordMatchCount() int{
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}