package response

import (
	//"bytes"
	//"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	//"runtime/debug"
)

type listGetAllResponse struct {
	record  *record.Record
	pkg     *tcaplus_protocol_cs.TCaplusPkg
	offset  int32
	idx     int32
	listidx int32
}

func newListGetAllResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*listGetAllResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListGetAllRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &listGetAllResponse{pkg: pkg}, nil
}

func (res *listGetAllResponse) GetResult() int {
	ret := int(res.pkg.Body.ListGetAllRes.Result)
	return ret
}

func (res *listGetAllResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *listGetAllResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *listGetAllResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *listGetAllResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *listGetAllResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *listGetAllResponse) GetRecordCount() int {
	if res.pkg.Body.ListGetAllRes.Result == 0 {
		Result := res.pkg.Body.ListGetAllRes.ResultInfo
		return int(Result.ElementNum)
	}
	return 0
}

func (res *listGetAllResponse) FetchRecord() (*record.Record, error) {
	data := res.pkg.Body.ListGetAllRes.ResultInfo
	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	if res.idx >= int32(data.ElementNum) {
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

	//unpack
	if err := rec.UnPackKey(); err != nil {
		logger.ERR("record unpack key failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}

	readBytes := uint32(0)
	err := unpackElementBuff(data.ElementsBuff, uint32(res.offset), data.ElementsBuffLen, &rec.Index,
		&readBytes, rec.ValueMap)
	res.idx += 1
	res.offset += int32(readBytes)
	res.record = rec
	return rec, err
}

func (res *listGetAllResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *listGetAllResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *listGetAllResponse) HaveMoreResPkgs() int {
	if 1 == res.pkg.Body.ListGetAllRes.IsCompleteFlag {
		return 0
	} else {
		return 1
	}
}

func (res *listGetAllResponse) GetTotalNum() int {
	return int(res.pkg.Body.ListGetAllRes.TotalElementNumOnServer)
}

func (res *listGetAllResponse) GetFailedNum() int {
	return 0
}

func (res *listGetAllResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}

func (res *listGetAllResponse) GetRecordMatchCount() int {
	if 0 == res.pkg.Body.ListGetAllRes.Result {
		return int(res.pkg.Body.ListGetAllRes.TotalElementNumOnServer)
	}
	return terror.GEN_ERR_ERR
}

func (res *listGetAllResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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
