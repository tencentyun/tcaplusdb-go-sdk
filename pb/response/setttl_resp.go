package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type setTtlResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	offset int32
	idx    int32
}

func newSetTtlResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*setTtlResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusSetTTLRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &setTtlResponse{pkg: pkg}, nil
}

func (res *setTtlResponse) GetResult() int {
	ret := int(res.pkg.Body.TCaplusSetTTLRes.Result)
	return ret
}

func (res *setTtlResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *setTtlResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *setTtlResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *setTtlResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *setTtlResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *setTtlResponse) GetRecordCount() int {
	if res.pkg.Body.TCaplusSetTTLRes.Result == 0 {
		return int(res.pkg.Body.TCaplusSetTTLRes.RecordNum)
	}
	return 0
}

func (res *setTtlResponse) FetchRecord() (*record.Record, error) {

	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	data := res.pkg.Body.TCaplusSetTTLRes
	if res.idx >= int32(data.RecordNum) || res.idx >= int32(len(data.KeyInfo)) ||
		res.idx >= int32(len(data.Index)) || res.idx >= int32(len(data.RecordResult)) {
		logger.ERR("resp fetch record over, current idx: %d, ", res.idx)
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}

	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Head.KeyInfo,
		ValueSet:    nil,
		UpdFieldSet: nil,
	}

	rec.Index = data.Index[res.idx]
	rec.KeySet = data.KeyInfo[res.idx]

	//unpack
	if err := rec.UnPackKey(); err != nil {
		logger.ERR("record unpack key failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}

	res.idx += 1

	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec

	if ret := int(data.RecordResult[res.idx-1]); ret != terror.GEN_ERR_SUC {
		return rec, &terror.ErrorCode{Code: ret}
	}
	return rec, nil
}

func (res *setTtlResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *setTtlResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *setTtlResponse) HaveMoreResPkgs() int {
	return 0
}

func (res *setTtlResponse) GetTotalNum() int {
	return 0
}

func (res *setTtlResponse) GetFailedNum() int {
	return 0
}

func (res *setTtlResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}

func (res *setTtlResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *setTtlResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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