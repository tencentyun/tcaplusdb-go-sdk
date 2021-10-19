package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type countResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
}

func newCountResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*countResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.GetTableRecordCountRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &countResponse{pkg: pkg}, nil
}

func (res *countResponse) GetResult() int {
	ret := int(res.pkg.Body.GetTableRecordCountRes.Result)
	return ret
}

func (res *countResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *countResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *countResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *countResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *countResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *countResponse) GetRecordCount() int {
	return 0
}

func (res *countResponse) FetchRecord() (*record.Record, error) {
	return nil, nil
}

func (res *countResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *countResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *countResponse) HaveMoreResPkgs() int {
	return 0
}
func (res *countResponse) GetTotalNum() int {
	return 0
}
func (res *countResponse) GetFailedNum() int {
	return 0
}

func (res *countResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}

func (res *countResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *countResponse) GetTableRecordCount() int {
	if res.GetResult() != terror.GEN_ERR_SUC {
		return terror.GEN_ERR_ERR
	}
	return int(res.pkg.Body.GetTableRecordCountRes.Count)
}

func (res *countResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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
