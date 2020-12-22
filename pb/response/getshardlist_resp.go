package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type getShardListResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
}

func newGetShardListResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*getShardListResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.GetShardListRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &getShardListResponse{pkg: pkg}, nil
}

func (res *getShardListResponse) GetResult() int {
	ret := int(res.pkg.Body.GetShardListRes.Result)
	return ret
}

func (res *getShardListResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *getShardListResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *getShardListResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *getShardListResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *getShardListResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *getShardListResponse) GetRecordCount() int {
	return 0
}

func (res *getShardListResponse) FetchRecord() (*record.Record, error) {
	return nil, nil
}

func (res *getShardListResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *getShardListResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *getShardListResponse) HaveMoreResPkgs() int {
	return 0
}
func (res *getShardListResponse) GetTotalNum() int {
	return 0
}
func (res *getShardListResponse) GetFailedNum() int {
	return 0
}

func (res *getShardListResponse) FetchErrorRecord() (*record.Record, error) {
	return nil,nil
}

func (res *getShardListResponse) GetRecordMatchCount() int{
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *getShardListResponse) GetTcaplusPackagePtr() *tcaplus_protocol_cs.TCaplusPkg {
	return res.pkg
}
