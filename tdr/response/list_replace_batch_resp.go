package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type listReplaceBatchResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	offset int32
	idx    int32
}

func newListReplaceBatchResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*listReplaceBatchResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListReplaceBatchRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &listReplaceBatchResponse{pkg: pkg}, nil
}

func (res *listReplaceBatchResponse) GetResult() int {
	ret := int(res.pkg.Body.ListReplaceBatchRes.Result)
	return ret
}

func (res *listReplaceBatchResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *listReplaceBatchResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *listReplaceBatchResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *listReplaceBatchResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *listReplaceBatchResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *listReplaceBatchResponse) GetRecordCount() int {
	return int(res.pkg.Body.ListReplaceBatchRes.ElementNum)
}

func (res *listReplaceBatchResponse) FetchRecord() (*record.Record, error) {
	data := res.pkg.Body.ListReplaceBatchRes.ResultInfo
	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("%s", common.CovertToJson(res.pkg.Body.ListReplaceBatchRes))
	}

	if res.idx >= int32(res.pkg.Body.ListReplaceBatchRes.ElementNum) {
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
		Version:     res.pkg.Head.KeyInfo.Version,
		KeySet:      res.pkg.Head.KeyInfo,
		ValueSet:    nil,
		UpdFieldSet: nil,
	}
	//unpack
	if err := rec.UnPackKey(); err != nil {
		logger.ERR("record unpack key failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		res.idx++
		return nil, err
	}

	ret := int(res.pkg.Body.ListReplaceBatchRes.ElementResult[res.idx])
	if res.pkg.Body.ListReplaceBatchRes.HasElementBuff[res.idx] == 0 {
		if ret != 0 {
			res.idx++
			return rec, &terror.ErrorCode{Code: ret}
		}
		res.idx++
		return rec, nil
	}

	readBytes := uint32(0)
	err := unpackElementBuff(data.ElementsBuff, uint32(res.offset), data.ElementsBuffLen, &rec.Index,
		&readBytes, rec.ValueMap)
	rec.Index = res.pkg.Body.ListReplaceBatchRes.ElementIndexArray[res.idx]
	if err != nil {
		res.idx += 1
		logger.ERR("unpackElementBuff failed %s", err.Error())
		return rec, err
	}

	res.offset += int32(readBytes)
	if ret != 0 {
		res.idx += 1
		return rec, &terror.ErrorCode{Code: ret}
	}

	res.idx += 1
	return rec, nil
}

func (res *listReplaceBatchResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *listReplaceBatchResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *listReplaceBatchResponse) HaveMoreResPkgs() int {
	if 0 == res.pkg.Body.ListReplaceBatchRes.Result && 0 != res.pkg.Body.ListReplaceBatchRes.LeftNum {
		return 1
	} else {
		return 0
	}
}

func (res *listReplaceBatchResponse) GetRecordMatchCount() int {
	return int(res.pkg.Body.ListReplaceBatchRes.TotalNum)
}

func (res *listReplaceBatchResponse) GetAffectedRecordNum() int32 {
	return int32(res.pkg.Body.ListReplaceBatchRes.TotalNum)
}

func (res *listReplaceBatchResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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
