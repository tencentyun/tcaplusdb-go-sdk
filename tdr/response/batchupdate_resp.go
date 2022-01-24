package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type batchUpdateResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	offset int32
	idx    int32
}

func newBatchUpdateResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*batchUpdateResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.BatchUpdateRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &batchUpdateResponse{pkg: pkg}, nil
}

func (res *batchUpdateResponse) GetResult() int {
	ret := int(res.pkg.Body.BatchUpdateRes.Result)
	return ret
}

func (res *batchUpdateResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen])
	return tableName
}

func (res *batchUpdateResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *batchUpdateResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *batchUpdateResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *batchUpdateResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *batchUpdateResponse) GetRecordCount() int {
	if res.pkg.Body.BatchUpdateRes.Result == 0 ||
		res.pkg.Body.BatchUpdateRes.Result == int32(terror.COMMON_INFO_DATA_NOT_MODIFIED) {
		return int(res.pkg.Body.BatchUpdateRes.RecordNum)
	}
	return 0
}

func (res *batchUpdateResponse) FetchRecord() (*record.Record, error) {

	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	data := res.pkg.Body.BatchUpdateRes
	if res.idx >= int32(data.RecordNum) || res.offset >= data.BatchValueLen {
		logger.ERR("resp fetch record over, current idx: %d, ", res.idx)
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	logger.DEBUG("read bytes: %d, total bytes: %d", res.offset, data.BatchValueLen)
	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Head.KeyInfo,
		ValueSet:    nil, //res.pkg.Body.BatchUpdateRes,
		UpdFieldSet: nil,
	}

	//unpack
	readBytes, err := unpackRecordKV(data.BatchValueInfo[res.offset:data.BatchValueLen],
		data.BatchValueLen-res.offset, rec.KeyMap, rec.ValueMap, &rec.Version)
	if err != nil {
		res.idx += 1
		logger.ERR("record unpack failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}
	logger.DEBUG("record unpack success, key: %+v, value: %+v", rec.KeyMap, rec.ValueMap)
	res.idx += 1
	res.offset += readBytes

	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec

	if ret := int(data.RecordResult[res.idx-1]); ret != terror.GEN_ERR_SUC {
		return rec, &terror.ErrorCode{Code: ret}
	}
	return rec, nil
}

func (res *batchUpdateResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *batchUpdateResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *batchUpdateResponse) HaveMoreResPkgs() int {
	if 0 == res.pkg.Body.BatchUpdateRes.Result && 0 != res.pkg.Body.BatchUpdateRes.LeftNum {
		return 1
	} else {
		return 0
	}
}

func (res *batchUpdateResponse) GetTotalNum() int {
	return int(res.pkg.Body.BatchUpdateRes.TotalNum)
}

func (res *batchUpdateResponse) GetFailedNum() int {
	return 0
}

func (res *batchUpdateResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}

func (res *batchUpdateResponse) GetRecordMatchCount() int {
	if res.pkg.Body.BatchUpdateRes.Result == 0 {
		return int(res.pkg.Body.BatchUpdateRes.TotalNum)
	}
	return terror.GEN_ERR_ERR
}

func (res *batchUpdateResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
	if res.pkg.Head.PerfTestLen == 0 {
		return nil
	}
	perf := tcaplus_protocol_cs.NewPerfTest()
	err := perf.Unpack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion, res.pkg.Head.PerfTest)
	if err != nil {
		logger.ERR("unpack perf error: %s", err)
		return nil
	}
	perf.ApiRecvTime = recvTime
	return perf
}
