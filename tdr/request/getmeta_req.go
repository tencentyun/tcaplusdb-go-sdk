package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cs_pool"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

type getMetaRequest struct {
	appId     uint64
	zoneId    uint32
	tableName string
	cmd       int
	seq       uint32
	pkg       *tcaplus_protocol_cs.TCaplusPkg
	isPB      bool
}

func newGetMetaRequest(appId uint64, zoneId uint32, tableName string, cmd int,
	seq uint32, pkg *tcaplus_protocol_cs.TCaplusPkg, isPB bool) (*getMetaRequest, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.MetadataGetReq == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}

	pkg.Body.MetadataGetReq.MetadataVersion = 0
	req := &getMetaRequest{
		appId:     appId,
		zoneId:    zoneId,
		tableName: tableName,
		cmd:       cmd,
		seq:       seq,
		pkg:       pkg,
		isPB:      isPB,
	}
	return req, nil
}

func (req *getMetaRequest) AddRecord(index int32) (*record.Record, error) {
	return nil, nil
}

func (req *getMetaRequest) SetAsyncId(id uint64) {
	req.pkg.Head.AsynID = id
}

func (req *getMetaRequest) SetVersionPolicy(p uint8) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support VersionPolicy"}
}

func (req *getMetaRequest) SetResultFlag(flag int) error {
	return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "Get not Support ResultFlag"}
}

func (req *getMetaRequest) Pack() ([]byte, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return nil, &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}

	if logger.GetLogLevel() == "DEBUG" {
		logger.DEBUG("pack request %s", common.CsHeadVisualize(req.pkg.Head))
	}
	data, err := req.pkg.Pack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	if err != nil {
		logger.ERR("getMetaRequest pack failed, %s", err.Error())
		return nil, err
	}
	logger.DEBUG("record pack success, app %d zone %d table %s", req.appId, req.zoneId, req.tableName)
	return data, nil
}

func (req *getMetaRequest) GetZoneId() uint32 {
	return req.zoneId
}

func (req *getMetaRequest) GetKeyHash() (uint32, error) {
	if req.pkg == nil {
		logger.ERR("Request can not second use")
		return uint32(terror.RequestHasHasNoPkg), &terror.ErrorCode{Code: terror.RequestHasHasNoPkg,
			Message: "Request can not second use"}
	}
	defer func() {
		cs_pool.PutTcaplusCSPkg(req.pkg)
		req.pkg = nil
	}()
	return uint32(time.Now().UnixNano()), nil
}

func (req *getMetaRequest) SetFieldNames(valueNameList []string) error {
	return nil
}

func (req *getMetaRequest) SetUserBuff(userBuffer []byte) error {
	return setUserBuffer(req.pkg, userBuffer)
}

func (req *getMetaRequest) GetSeq() int32 {
	return req.pkg.Head.Seq
}

func (req *getMetaRequest) SetSeq(seq int32) {
	req.pkg.Head.Seq = seq
}

func (req *getMetaRequest) SetResultLimit(limit int32, offset int32) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *getMetaRequest) SetMultiResponseFlag(multi_flag byte) int32 {
	return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
}

func (req *getMetaRequest) SetResultFlagForSuccess(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *getMetaRequest) SetResultFlagForFail(result_flag byte) int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (req *getMetaRequest) SetSplitTableKeyBuff(splitkey []byte) int {
	req.pkg.Head.SplitTableKeyBuff = splitkey
	req.pkg.Head.SplitTableKeyBuffLen = uint32(len(splitkey))
	return terror.GEN_ERR_SUC
}

func (req *getMetaRequest) SetPerfTest(sendTime uint64) int {
	perf := tcaplus_protocol_cs.NewPerfTest()
	perf.ApiSendTime = sendTime
	perf.Version = tcaplus_protocol_cs.PerfTestCurrentVersion
	p, err := perf.Pack(tcaplus_protocol_cs.PerfTestCurrentVersion)
	if err != nil {
		logger.ERR("pack perf error: %s", err)
		return terror.API_ERR_PARAMETER_INVALID
	}
	req.pkg.Head.PerfTest = p
	req.pkg.Head.PerfTestLen = uint32(len(p))
	return terror.GEN_ERR_SUC
}

func (req *getMetaRequest) SetFlags(flag int32) int {
	return setFlags(req.pkg, flag)
}

func (req *getMetaRequest) ClearFlags(flag int32) int {
	return clearFlags(req.pkg, flag)
}

func (req *getMetaRequest) GetFlags() int32 {
	return req.pkg.Head.Flags
}
