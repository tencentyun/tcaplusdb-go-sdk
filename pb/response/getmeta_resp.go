package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

type getMetaResponse struct {
	record *record.Record
	pkg    *tcaplus_protocol_cs.TCaplusPkg
}

func newGetMetaResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*getMetaResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.MetadataGetRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &getMetaResponse{pkg: pkg}, nil
}

func (res *getMetaResponse) GetResult() int {
	ret := int(res.pkg.Body.MetadataGetRes.Result)
	return ret
}

func (res *getMetaResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *getMetaResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *getMetaResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *getMetaResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *getMetaResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

func (res *getMetaResponse) GetRecordCount() int {
	return 0
}

func (res *getMetaResponse) FetchRecord() (*record.Record, error) {
	return nil, nil
}

func (res *getMetaResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *getMetaResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}
func (res *getMetaResponse) HaveMoreResPkgs() int {
	return 0
}
func (res *getMetaResponse) GetTotalNum() int {
	return 0
}
func (res *getMetaResponse) GetFailedNum() int {
	return 0
}

func (res *getMetaResponse) FetchErrorRecord() (*record.Record, error) {
	return nil,nil
}

func (res *getMetaResponse) GetRecordMatchCount() int{
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *getMetaResponse) GetTcaplusPackagePtr() *tcaplus_protocol_cs.TCaplusPkg {
	return res.pkg
}
