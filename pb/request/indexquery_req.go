package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

const (
	VERSION_FOR_SQL = 1
)

type indexQueryRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	record    *record.Record
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	valueNameMap map[string]bool
}

func newIndexQueryRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg) (*indexQueryRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusSqlReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.TCaplusSqlReq.Version = VERSION_FOR_SQL
	req := &indexQueryRequest{
		appId:        appId,
		zoneId:       zoneId,
		tableName:    tableName,
		cmd:          cmd,
		seq:          seq,
		record:       nil,
		pkg:          pkg,
		valueNameMap: make(map[string]bool),
	}
	return req, nil
}

func (req *indexQueryRequest) AddRecord(index int32) (*record.Record, error) {
	return nil, &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH,
		Message: "index query not support AddRecord"}
}

func (req *indexQueryRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *indexQueryRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH,
		Message: "list get not support SetVersionPolicy"}
}

func (req *indexQueryRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH,
		Message: "list get not support SetResultFlag"}
}

func (req *indexQueryRequest) Pack() ([]byte, error) {
	logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("indexQueryRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *indexQueryRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *indexQueryRequest) GetKeyHash() (uint32, error) {
	return 5, nil
}

func (req *indexQueryRequest) SetFieldNames(valueNameList []string) error {
	return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH,
		Message: "index query not support SetFieldNames"}
}

func (req *indexQueryRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *indexQueryRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *indexQueryRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *indexQueryRequest)SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *indexQueryRequest)SetMultiResponseFlag(multi_flag byte) int32{
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *indexQueryRequest)SetResultFlagForSuccess(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *indexQueryRequest)SetResultFlagForFail(result_flag byte) int {
	return terror.GEN_ERR_SUC
}

func (req *indexQueryRequest)SetSql(query string) int {
	req.pkg.Body.TCaplusSqlReq.Version = VERSION_FOR_SQL
	req.pkg.Body.TCaplusSqlReq.Sql = query
	return terror.GEN_ERR_SUC
}
