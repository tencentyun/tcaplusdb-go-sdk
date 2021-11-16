package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type insertResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
}

func newInsertResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*insertResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.InsertRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &insertResponse{pkg: pkg}, nil
}

func (res *insertResponse) GetResult() int {
	ret := int(res.pkg.Body.InsertRes.Result)
	return ret
}

func (res *insertResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *insertResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *insertResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *insertResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *insertResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

/*
 各个命令请求需要具体问题具体分析
 新版的ResultFlag 的GetRecordCount的规则
 在成功的场景下, 设置了如下的ResultFlag返回数据规则
 1. TCaplusValueFlag_NOVALUE 不返还
 2. TCaplusValueFlag_SAMEWITHREQUEST 返回
 3. TCaplusValueFlag_ALLVALUE 返回
 4. TCaplusValueFlag_ALLOLDVALUE 返回

 在失败的场景下, 设置了如下的ResultFlag会返回数据
 1. TCaplusValueFlag_NOVALUE 不返还
 2. TCaplusValueFlag_SAMEWITHREQUEST 返回
 3. TCaplusValueFlag_ALLVALUE 返回
 4. TCaplusValueFlag_ALLOLDVALUE 返回

*/
func (res *insertResponse) GetRecordCount() int {
	if 0 != (res.pkg.Body.InsertRes.Flag & (1 << 6)) {
		//新版本的result flag 通过ResultFlag判断
		if res.pkg.Body.InsertRes.Result == 0 {
			ret := GetResultFlagByBit(res.pkg.Body.InsertRes.Flag, true)
			if tcaplus_protocol_cs.TCaplusValueFlag_SAMEWITHREQUEST == ret ||
				tcaplus_protocol_cs.TCaplusValueFlag_ALLVALUE == ret {
				return 1
			}
		} else {
			ret := GetResultFlagByBit(res.pkg.Body.InsertRes.Flag, false)
			if (tcaplus_protocol_cs.TCaplusValueFlag_SAMEWITHREQUEST == ret ||
				tcaplus_protocol_cs.TCaplusValueFlag_ALLVALUE == ret) &&
				res.pkg.Body.InsertRes.ResultInfo.CompactValueSet.FieldIndexNum > 0 {
				return 1
			}
		}
	} else {
		//老版本的result flag 通过ResultFlag判断
		if (res.pkg.Body.InsertRes.Flag == 1 || res.pkg.Body.InsertRes.Flag == 2) &&
			(res.pkg.Body.InsertRes.Result == 0 ||
				res.pkg.Body.InsertRes.Result == int32(terror.SVR_ERR_FAIL_RECORD_EXIST)) {
			return 1
		}
	}
	return 0
}

func (res *insertResponse) FetchRecord() (*record.Record, error) {
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
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Head.KeyInfo,
		ValueSet:    res.pkg.Body.InsertRes.ResultInfo,
		UpdFieldSet: nil,
	}

	//unpack
	if err := rec.UnPackKey(); err != nil {
		logger.ERR("record unpack key failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}

	if err := rec.UnPackValue(); err != nil {
		logger.ERR("record unpack value failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}
	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec
	return rec, nil
}

func (res *insertResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *insertResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *insertResponse) HaveMoreResPkgs() int {
	return 0
}

func (res *insertResponse) GetTotalNum() int {
	return 0
}

func (res *insertResponse) GetFailedNum() int {
	return 0
}

func (res *insertResponse) FetchErrorRecord() (*record.Record, error) {
	return nil, nil
}

func (res *insertResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *insertResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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
