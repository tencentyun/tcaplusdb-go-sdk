package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

type listDeleteAllResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
	idx    int32
}

func newListDeleteAllResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*listDeleteAllResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListDeleteAllRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &listDeleteAllResponse{pkg: pkg}, nil
}

func (res *listDeleteAllResponse) GetResult() int {
	ret := int(res.pkg.Body.ListDeleteAllRes.Result)
	return ret
}

func (res *listDeleteAllResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen])
	return tableName
}

func (res *listDeleteAllResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *listDeleteAllResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *listDeleteAllResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *listDeleteAllResponse) GetAsyncId() uint64 {
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
func (res *listDeleteAllResponse) GetRecordCount() int {
	return 0
}

func (res *listDeleteAllResponse) FetchRecord() (*record.Record, error) {
	if int(res.idx) > res.GetRecordCount() {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}

	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		Index:       -1,
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

	res.record = rec
	return rec, nil
}

func (res *listDeleteAllResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *listDeleteAllResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *listDeleteAllResponse) HaveMoreResPkgs() int {
	return 0
}

func (res *listDeleteAllResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *listDeleteAllResponse) GetAffectedRecordNum() int32 {
	return res.pkg.Body.ListDeleteAllRes.AffectedElementNum
}
