package response

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"unsafe"
)

type indexQueryResponse struct {
	record  *record.Record
	pkg     *tcaplus_protocol_cs.TCaplusPkg
	offset  int32
	idx     int32
	listidx int32
}

func newIndexQueryResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*indexQueryResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.TCaplusSqlRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &indexQueryResponse{pkg: pkg}, nil
}

func (res *indexQueryResponse) GetResult() int {
	ret := int(res.pkg.Body.TCaplusSqlRes.Result)
	return ret
}

func (res *indexQueryResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0 : res.pkg.Head.RouterInfo.TableNameLen-1])
	return tableName
}

func (res *indexQueryResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *indexQueryResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *indexQueryResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *indexQueryResponse) GetAsyncId() uint64 {
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
func (res *indexQueryResponse) GetRecordCount() int {
	if res.pkg.Body.TCaplusSqlRes.Result == 0 {
		return int(res.pkg.Body.TCaplusSqlRes.RecordNum)
	}
	return 0
}

func (res *indexQueryResponse) FetchRecord() (*record.Record, error) {
	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	data := res.pkg.Body.TCaplusSqlRes
	if res.idx >= int32(data.RecordNum) || res.offset >= data.ValueLen {
		logger.ERR("resp fetch record over, current idx: %d, ", res.idx)
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}
	logger.DEBUG("read bytes: %d, total bytes: %d", res.offset, data.ValueLen)
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
	readBytes, err := unpackRecordKV(data.Value[res.offset:data.ValueLen],
		data.ValueLen-res.offset, rec.KeyMap, rec.ValueMap, &rec.Version)
	if err != nil {
		logger.ERR("record unpack failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}
	logger.DEBUG("record unpack success, key: %+v, value: %+v", rec.KeyMap, rec.ValueMap)
	res.idx += 1
	res.offset += readBytes

	logger.DEBUG("record unpack success, app %d zone %d table %s", rec.AppId, rec.ZoneId, rec.TableName)
	res.record = rec
	return rec, nil
}

func (res *indexQueryResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *indexQueryResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *indexQueryResponse) HaveMoreResPkgs() int {
	if res.pkg.Body.TCaplusSqlRes.IsCompleteFlag == 0 {
		return 1
	}
	return 0
}

func (res *indexQueryResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}

func (res *indexQueryResponse) GetPerfTest(recvTime uint64) *tcaplus_protocol_cs.PerfTest {
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

func (res *indexQueryResponse) GetSqlType() int {
	sqlType := res.pkg.Body.TCaplusSqlRes.SqlType
	if sqlType == policy.RECORD_SQL_QUERY_TYPE {
		return policy.RECORD_SQL_QUERY_TYPE
	} else if sqlType == policy.AGGREGATIONS_SQL_QUERY_TYPE {
		return policy.AGGREGATIONS_SQL_QUERY_TYPE
	} else {
		return policy.INVALID_SQL_TYPE
	}
}

func (res *indexQueryResponse) ProcAggregationSqlQueryType() ([]string, error) {
	var rs []string
	result, err := res.FetchSqlResult()
	if err != nil {
		return rs, err
	}
	for i := uint32(0); i < result.rowsNum; i++ {
		r, err := result.FetchRow()
		if err != nil {
			return rs, err
		}
		str := ""
		for j := int32(0); j < r.fieldNum; j++ {
			if j > 0 {
				str += ","
			}
			f, err := r.FetchField()
			if err != nil {
				return rs, err
			}
			switch f.fieldType {
			case policy.TYPE_BOOL:
				str += fmt.Sprint(f.GetBool())
			case policy.TYPE_INT8:
				str += fmt.Sprint(f.GetInt8())
			case policy.TYPE_UINT8:
				str += fmt.Sprint(f.GetUInt8())
			case policy.TYPE_INT16:
				str += fmt.Sprint(f.GetInt16())
			case policy.TYPE_UINT16:
				str += fmt.Sprint(f.GetUInt16())
			case policy.TYPE_INT32:
				str += fmt.Sprint(f.GetInt32())
			case policy.TYPE_UINT32:
				str += fmt.Sprint(f.GetUInt32())
			case policy.TYPE_INT64:
				str += fmt.Sprint(f.GetInt64())
			case policy.TYPE_UINT64:
				str += fmt.Sprint(f.GetUInt64())
			case policy.TYPE_FLOAT:
				str += fmt.Sprint(f.GetFloat32())
			case policy.TYPE_DOUBLE:
				str += fmt.Sprint(f.GetFloat64())
			case policy.TYPE_STRING:
				str += fmt.Sprint(f.GetString())
			default:
				desc := fmt.Sprintf("field type %d not support", f.fieldType)
				return nil, &terror.ErrorCode{Code: terror.GEN_ERR_ERR, Message: desc}
			}
		}
		rs = append(rs, str)
	}
	return rs, nil
}

func (res *indexQueryResponse) FetchSqlResult() (*sqlResult, error) {
	if res.pkg.Body.TCaplusSqlRes.SqlType != policy.AGGREGATIONS_SQL_QUERY_TYPE {
		return nil, &terror.ErrorCode{Code: terror.GEN_ERR_ERR, Message: "sql type not AGGREGATIONS_SQL_QUERY_TYPE"}
	}
	r := &sqlResult{}
	s := res.pkg.Body.TCaplusSqlRes
	r.Set(s.Result, s.SqlType, s.Version, s.RecordNum, s.Value, s.ValueLen)
	return r, nil
}

type sqlResult struct {
	result   int32
	sqlType  int32
	version  int32
	rowsNum  uint32
	value    []byte
	valueLen int32
	cursor   int32
}

func (s *sqlResult) Result() int32 {
	return s.result
}

func (s *sqlResult) SqlType() int32 {
	return s.sqlType
}

func (s *sqlResult) Version() int32 {
	return s.version
}

func (s *sqlResult) RowsNum() uint32 {
	return s.rowsNum
}

func (s *sqlResult) Set(result, sqlType, version int32, rowsNum uint32, value []byte, size int32) int {
	if result == 0 && sqlType != policy.RECORD_SQL_QUERY_TYPE && sqlType != policy.AGGREGATIONS_SQL_QUERY_TYPE {
		return terror.GEN_ERR_ERR
	}
	s.result, s.sqlType, s.rowsNum, s.version, s.value, s.valueLen = result, sqlType, rowsNum, version, value, size
	return terror.GEN_ERR_SUC
}

func (s *sqlResult) FetchRow() (*row, error) {
	if len(s.value) == 0 || s.valueLen == 0 || s.rowsNum == 0 {
		return nil, &terror.ErrorCode{Code: terror.GEN_ERR_ERR, Message: "value is nil"}
	}

	if s.cursor >= s.valueLen {
		return nil, &terror.ErrorCode{Code: terror.GEN_ERR_ERR, Message: "no more row"}
	}

	if s.cursor+8 > s.valueLen {
		return nil, &terror.ErrorCode{Code: terror.GEN_ERR_ERR, Message: "value is not format"}
	}

	length := *(*int32)(unsafe.Pointer(&s.value[s.cursor]))
	s.cursor += 4
	fieldNum := *(*int32)(unsafe.Pointer(&s.value[s.cursor]))
	r := &row{}
	r.set(s.value[s.cursor+4:s.cursor+length], length-4, fieldNum)
	s.cursor += length
	return r, nil
}

type row struct {
	value    []byte
	size     int32
	fieldNum int32
	cursor   int32
}

func (r *row) set(value []byte, size int32, fieldNum int32) int {
	r.value, r.size, r.fieldNum = value, size, fieldNum
	return 0
}

func (r *row) FetchField() (*field, error) {
	if len(r.value) == 0 || r.size == 0 {
		return nil, &terror.ErrorCode{Code: terror.GEN_ERR_ERR, Message: "field is nil"}
	}

	if r.fieldNum <= 0 {
		return nil, &terror.ErrorCode{Code: terror.GEN_ERR_ERR, Message: "fieldNum <= 0"}
	}

	if r.cursor >= r.size {
		return nil, &terror.ErrorCode{Code: terror.GEN_ERR_ERR, Message: "no more field"}
	}

	ftype := *(*int32)(unsafe.Pointer(&r.value[r.cursor]))
	r.cursor += 4
	length := *(*int32)(unsafe.Pointer(&r.value[r.cursor]))
	r.cursor += 4
	f := &field{}
	f.set(ftype, r.value[r.cursor:r.cursor+length], length)
	r.cursor += length
	return f, nil
}

func (r *row) FieldsNum() int32 { return r.fieldNum }

type field struct {
	fieldType int32
	value     []byte
	size      int32
}

func (f *field) set(fieldType int32, value []byte, size int32) {
	f.fieldType, f.value, f.size = fieldType, value, size
}

func (f *field) FieldType() int32 { return f.fieldType }

func (f *field) GetBool() bool { return *(*bool)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetInt8() int8 { return *(*int8)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetUInt8() uint8 { return *(*uint8)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetInt16() int16 { return *(*int16)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetUInt16() uint16 { return *(*uint16)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetInt32() int32 { return *(*int32)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetUInt32() uint32 { return *(*uint32)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetInt64() int64 { return *(*int64)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetUInt64() uint64 { return *(*uint64)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetFloat32() float32 { return *(*float32)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetFloat64() float64 { return *(*float64)(unsafe.Pointer(&f.value[0])) }

func (f *field) GetString() string { return string(f.value) }
