package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type traverseResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	offset int32
	idx    int32
}

func newTraverseResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*traverseResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TableTraverseRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &traverseResponse{pkg: pkg}, nil
}

func (res *traverseResponse) GetResult() int {
	ret := int(res.pkg.Body.TableTraverseRes.Result)
	return ret
}

func (res *traverseResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *traverseResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *traverseResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *traverseResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *traverseResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *traverseResponse) GetRecordCount() int {
	if res.pkg.Body.TableTraverseRes.Result == 0 ||
		res.pkg.Body.TableTraverseRes.Result == int32(terror.COMMON_INFO_DATA_NOT_MODIFIED) {
		return int(res.pkg.Body.TableTraverseRes.RecordNum)
	}
	return 0
}

func (res *traverseResponse) FetchRecord() (*record.Record, error) {

	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	data := res.pkg.Body.TableTraverseRes
	if res.idx >= int32(data.RecordNum) || res.offset >= data.BatchValueLen {
		logger.ERR("resp fetch record over, current idx: %d, ", res.idx)
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	logger.DEBUG("read bytes: %d, total bytes: %d", res.offset, data.BatchValueLen)
	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Head.KeyInfo,
		ValueSet:    nil, //res.pkg.Body.TableTraverseRes,
		UpdFieldSet: nil,
	}

	//unpack
	readBytes, err := unpackRecord(data.BatchValueInfo[res.offset:data.BatchValueLen],
		data.BatchValueLen-res.offset, rec.KeyMap, rec.ValueMap, &rec.Version)
	if err != nil {
		logger.ERR("record unpack failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}
	logger.DEBUG("record unpack success, key: %+v, value: %+v", rec.KeyMap, rec.ValueMap)
	res.idx += 1
	res.offset += readBytes

	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec

	return rec, nil
}

func (res *traverseResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *traverseResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *traverseResponse) HaveMoreResPkgs() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *traverseResponse) GetTotalNum() int {
	return 0
}

func (res *traverseResponse) GetFailedNum() int {
	return 0
}

func (res *traverseResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}

func (res *traverseResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *traverseResponse) GetTcaplusPackagePtr() *tcaplus_protocol_cs.TCaplusPkg {
	return res.pkg
}

func (res *traverseResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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
