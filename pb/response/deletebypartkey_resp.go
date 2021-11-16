package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type deleteByPartKeyResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	offset int32
	idx    int32
}

func newDeleteByPartKeyResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*deleteByPartKeyResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.DeleteByPartkeyRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &deleteByPartKeyResponse{pkg: pkg}, nil
}

func (res *deleteByPartKeyResponse) GetResult() int {
	ret := int(res.pkg.Body.DeleteByPartkeyRes.Result)
	return ret
}

func (res *deleteByPartKeyResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *deleteByPartKeyResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *deleteByPartKeyResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *deleteByPartKeyResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *deleteByPartKeyResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *deleteByPartKeyResponse) GetRecordCount() int {
	return int(res.pkg.Body.DeleteByPartkeyRes.SucNum)
}

func (res *deleteByPartKeyResponse) FetchRecord() (*record.Record, error) {
	s := res.pkg.Body.DeleteByPartkeyRes
	//	logger.DEBUG("Result:%d, TotalNum:%d, OffSet:%d, SucNum:%d, FailNum:%d, IsCompleteFlag:%d",
	//	s.Result, s.TotalNum, s.OffSet, s.SucNum, s.FailNum, s.IsCompleteFlag)
	//	logger.DEBUG("success buff  len : %d： %v", s.SucKeysBuffLen,s.SucKeysBuff)
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
		KeySet:      res.pkg.Body.DeleteByPartkeyRes.FailKeys, //res.pkg.Head.KeyInfo,
		ValueSet:    nil,                                      //res.pkg.Body.DeleteByPartkeyRes.ResultInfo,
		UpdFieldSet: nil,
	}

	result := int32(0)
	readBytes := int32(0)
	if err := unpackSucKeysBuffLen(s.SucKeysBuff[res.offset:s.SucKeysBuffLen], s.SucKeysBuffLen-res.offset,
		&result, rec.KeyMap, &readBytes); err != nil {
		logger.ERR("record unpack succ keys failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}
	res.idx += 1
	res.offset += readBytes
	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec
	return rec, nil
}

func (res *deleteByPartKeyResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *deleteByPartKeyResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *deleteByPartKeyResponse) HaveMoreResPkgs() int {
	if 0 != res.pkg.Body.DeleteByPartkeyRes.Result || 1 == res.pkg.Body.DeleteByPartkeyRes.IsCompleteFlag {
		return 0
	} else {
		return 1
	}
}
func (res *deleteByPartKeyResponse) GetTotalNum() int {
	return int(res.pkg.Body.DeleteByPartkeyRes.TotalNum)
}

func (res *deleteByPartKeyResponse) GetFailedNum() int {
	return int(res.pkg.Body.DeleteByPartkeyRes.FailNum)
}

func (res *deleteByPartKeyResponse) FetchErrorRecord() (*record.Record, error) {
	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Body.DeleteByPartkeyRes.FailKeys, //res.pkg.Head.KeyInfo,
		ValueSet:    nil,                                      //res.pkg.Body.DeleteByPartkeyRes.ResultInfo,
		UpdFieldSet: nil,
	}
	if nil != rec.KeySet {
		rec.UnPackKey()
	}
	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	return rec, nil
}

func (res *deleteByPartKeyResponse) GetRecordMatchCount() int {
	return int(res.pkg.Body.DeleteByPartkeyRes.SucNum)
}

func (res *deleteByPartKeyResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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
