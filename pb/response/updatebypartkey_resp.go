package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type updataByPartKeyResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	offset int32
	idx    int32
}

func newUpdataByPartKeyResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*updataByPartKeyResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.UpdateByPartkeyRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &updataByPartKeyResponse{pkg: pkg}, nil
}

func (res *updataByPartKeyResponse) GetResult() int {
	ret := int(res.pkg.Body.UpdateByPartkeyRes.Result)
	return ret
}

func (res *updataByPartKeyResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *updataByPartKeyResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *updataByPartKeyResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *updataByPartKeyResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *updataByPartKeyResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *updataByPartKeyResponse) GetRecordCount() int {
	return int(res.pkg.Body.UpdateByPartkeyRes.SucNum)
}

func (res *updataByPartKeyResponse) FetchRecord() (*record.Record, error) {
	s := res.pkg.Body.UpdateByPartkeyRes
	if res.idx >= int32(s.SucNum) || res.offset >= s.SucKeysBuffLen {
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
	result := int32(0)
	readBytes := int32(0)
	if err := unpackSucKeysBuffLen(s.SucKeysBuff[res.offset:s.SucKeysBuffLen],
		s.SucKeysBuffLen-res.offset, &result, rec.KeyMap, &readBytes); err != nil {
		logger.ERR("record unpack succ keys failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}
	res.offset += readBytes
	res.idx += 1
	//logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec
	return rec, nil
}

func (res *updataByPartKeyResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *updataByPartKeyResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *updataByPartKeyResponse) HaveMoreResPkgs() int {
	if 0 != res.pkg.Body.UpdateByPartkeyRes.Result || 1 == res.pkg.Body.UpdateByPartkeyRes.IsCompleteFlag {
		return 0
	} else {
		return 1
	}
}
func (res *updataByPartKeyResponse) GetTotalNum() int {
	return int(res.pkg.Body.UpdateByPartkeyRes.TotalNum)
}

func (res *updataByPartKeyResponse) GetFailedNum() int {
	return int(res.pkg.Body.UpdateByPartkeyRes.FailNum)
}

func (res *updataByPartKeyResponse) FetchErrorRecord() (*record.Record, error) {
	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Body.UpdateByPartkeyRes.FailKeys, //res.pkg.Head.KeyInfo,
		ValueSet:    nil,                                      //res.pkg.Body.DeleteByPartkeyRes.ResultInfo,
		UpdFieldSet: nil,
	}
	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	return rec, nil
}

func (res *updataByPartKeyResponse) GetRecordMatchCount() int {
	return int(res.pkg.Body.UpdateByPartkeyRes.SucNum)
}

func (res *updataByPartKeyResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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
