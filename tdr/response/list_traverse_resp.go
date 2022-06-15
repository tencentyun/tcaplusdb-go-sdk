package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type listTraverseResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	offset int32
	idx    int32
}

func newListTraverseResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*listTraverseResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListTableTraverseRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &listTraverseResponse{pkg: pkg}, nil
}

func (res *listTraverseResponse) GetResult() int {
	ret := int(res.pkg.Body.ListTableTraverseRes.Result)
	return ret
}

func (res *listTraverseResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *listTraverseResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *listTraverseResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *listTraverseResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *listTraverseResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *listTraverseResponse) GetRecordCount() int {
	if res.pkg.Body.ListTableTraverseRes.Result == 0 ||
		res.pkg.Body.ListTableTraverseRes.Result == int32(terror.COMMON_INFO_DATA_NOT_MODIFIED) {
		return int(res.pkg.Body.ListTableTraverseRes.RecordNum)
	}
	return 0
}

func (res *listTraverseResponse) FetchRecord() (*record.Record, error) {
	data := res.pkg.Body.ListTableTraverseRes.ResultValueInfo
	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("%s", common.CovertToJson(res.pkg.Body.ListTableTraverseRes))
	}

	if res.idx >= int32(res.pkg.Body.ListTableTraverseRes.RecordNum) {
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
		KeySet:      res.pkg.Body.ListTableTraverseRes.ResultKeyInfo,
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

	readBytes := uint32(0)
	err := unpackElementBuff(data.ElementsBuff, uint32(res.offset), data.ElementsBuffLen, &rec.Index,
		&readBytes, rec.ValueMap)
	if err != nil {
		res.idx += 1
		logger.ERR("unpackElementBuff failed %s", err.Error())
		return rec, err
	}

	res.offset += int32(readBytes)
	res.idx += 1
	return rec, nil
}

func (res *listTraverseResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *listTraverseResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *listTraverseResponse) HaveMoreResPkgs() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *listTraverseResponse) GetTotalNum() int {
	return 0
}

func (res *listTraverseResponse) GetFailedNum() int {
	return 0
}

func (res *listTraverseResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}

func (res *listTraverseResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *listTraverseResponse) GetTcaplusPackagePtr() *tcaplus_protocol_cs.TCaplusPkg {
	return res.pkg
}

func (res *listTraverseResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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
