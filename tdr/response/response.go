package response

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

/*
	大多数响应都会用到的函数放到commonInterface接口中
	个别或者极少数响应的特殊方法放到TcaplusResponse
*/
type TcaplusResponse interface {
	commonInterface

	/** @brief    获取受影响的Records的条数.
	*  @retval >= 0               受影响的Records的条数.
	*  @retval < 0                操作失败，具体错误参见 @link ErrorCode @endlink
	*  注意，在当前版本中该函数仅对List类型的DeleteAll操作有效.
	 */
	GetAffectedRecordNum() int32

	/**
	@brief  获取PartkeyUpdate和PartkeyDelete符合条件的记录
	@retval 记录条数
	*/
	GetTotalNum() int

	/**
	@brief  获取PartkeyUpdate和PartkeyDelete失败的记录数
	@retval 记录条数
	*/
	GetFailedNum() int

	/** \brief  从PartkeyUpdate/PartkeyDelete结果中获取失败的记录信息
	 *  \retval(2)   指向返回记录的指针, 获取记录成功
	 */
	FetchErrorRecord() (*record.Record, error)

	/*
	   @brief 该函数仅用于索引查询类型为聚合查询时获取聚合结果
	   @retval(2) 返回聚合查询结果，返回NULL表示获取失败
	   @note 如果索引查询类型为记录查询，该函数将返回NULL
	         如果索引查询类型为聚合查询，则根据SqlResult类提供的函数来获取查询结果
	*/
	FetchSqlResult() (*sqlResult, error)

	/*
		仅聚合操作可以使用，请在使用前判断索引类型
		处理聚合分布式索引响应并返回处理后的结果，切片长度为结果行数，以逗号作为字段分割符。
	*/
	ProcAggregationSqlQueryType() ([]string, error)

	/*
		获取分布式索引类型 2 为聚合操作  1 为非聚合操作  0 为无效操作
	*/
	GetSqlType() int
}

type commonInterface interface {
	/*
		@brief  获取响应结果
		@retval int tcaplus api自定义错误码。 0，表示请求成功；非0,有错误码，可从terror.GetErrMsg(int)得到错误消息
	*/
	GetResult() int

	/*
		@brief  获取响应表名
		@retval string 响应消息对应的表名称
	*/
	GetTableName() string

	/*
		@brief  获取响应appId
		@retval uint64 响应消息对应的appId
	*/
	GetAppId() uint64

	/*
		@brief  获取响应zoneId
		@retval uint32 响应消息对应的zoneId
	*/
	GetZoneId() uint32

	/*
		@brief  获取响应命令
		@retval int 响应消息命令字，cmd包中的响应命令字
	*/
	GetCmd() int

	/*
		@brief  获取响应异步id，和请求对应
		@retval uint64 响应消息对应的异步id和请求对应
	*/
	GetAsyncId() uint64

	/*
		@brief  获取本响应中结果记录条数
		@retval int 响应中结果记录条数
	*/
	GetRecordCount() int

	/*
		@brief  从结果中获取一条记录
		@retval *record.Record 记录指针
		@retval error 错误码
	*/
	FetchRecord() (*record.Record, error)

	/*
		@brief  获取响应消息中的用户缓存信息
		@retval []byte 用户缓存二进制，和请求消息中的buffer内容一致
	*/
	GetUserBuffer() []byte

	/*
		@获取用户seq
	*/
	GetSeq() int32
	/*
	   @判断是否有更多的回包
	   @retval  1 有， 0 没有
	*/
	HaveMoreResPkgs() int

	/**
	    @brief  获取整个结果中的记录条数。既包括本响应返回的记录数，也包括本响应未返回的记录数。
	    @retval 记录条数
	    @note 该函数只能用于以下请求：
		(1）TCAPLUS_API_GET_BY_PARTKEY_REQ（部分key查询）；
		(2）TCAPLUS_API_LIST_GETALL_REQ（list表查询所有匹配的记录）；
		(3）TCAPLUS_API_BATCH_GET_REQ（批量查询）

	    @note 该函数的作用是：
	       (1）当使用了SetResultLimit()来限制返回的记录时，
	               使用该函数可以获取所有匹配的记录的个数，
	               包括本响应返回的记录数和本响应未返回的记录数，
	               而GetRecordCount()函数只能获取本响应返回的记录的个数；
	       (2）当所返回的记录很多时，需要分包的时候，
	               使用该函数可以获取总共的记录数，
	               即多个分包所有记录数的总和，
	               而GetRecordCount()函数只能返回单个分包中的(本响应中的)记录数.
	*/
	GetRecordMatchCount() int
}

func NewResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (TcaplusResponse, error) {
	if nil == pkg {
		logger.ERR("para pkg is nil")
		return nil, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "pkg is nil"}
	}

	resp := &tcapResponse{}
	var err error

	switch pkg.Head.Cmd {
	case cmd.TcaplusApiInsertRes:
		resp.commonInterface, err = newInsertResponse(pkg)
	case cmd.TcaplusApiGetRes:
		resp.commonInterface, err = newGetResponse(pkg)
	case cmd.TcaplusApiIncreaseRes:
		resp.commonInterface, err = newIncreaseResponse(pkg)
	case cmd.TcaplusApiUpdateRes:
		resp.commonInterface, err = newUpdateResponse(pkg)
	case cmd.TcaplusApiReplaceRes:
		resp.commonInterface, err = newReplaceResponse(pkg)
	case cmd.TcaplusApiDeleteRes:
		resp.commonInterface, err = newDeleteResponse(pkg)
	case cmd.TcaplusApiBatchGetRes:
		resp.commonInterface, err = newBatchGetResponse(pkg)
	case cmd.TcaplusApiGetByPartkeyRes:
		resp.commonInterface, err = newGetByPartKeyResponse(pkg)
	case cmd.TcaplusApiDeleteByPartkeyRes:
		resp.commonInterface, err = newDeleteByPartKeyResponse(pkg)
	case cmd.TcaplusApiUpdateByPartkeyRes:
		resp.commonInterface, err = newUpdataByPartKeyResponse(pkg)
	case cmd.TcaplusApiListGetAllRes:
		resp.commonInterface, err = newListGetAllResponse(pkg)
	case cmd.TcaplusApiListAddAfterRes:
		resp.commonInterface, err = newListAddAfterResponse(pkg)
	case cmd.TcaplusApiListGetRes:
		resp.commonInterface, err = newListGetResponse(pkg)
	case cmd.TcaplusApiListDeleteRes:
		resp.commonInterface, err = newListDeleteResponse(pkg)
	case cmd.TcaplusApiListDeleteAllRes:
		resp.commonInterface, err = newListDeleteAllResponse(pkg)
	case cmd.TcaplusApiListReplaceRes:
		resp.commonInterface, err = newListReplaceResponse(pkg)
	case cmd.TcaplusApiListDeleteBatchRes:
		resp.commonInterface, err = newListDeleteBatchResponse(pkg)
	case cmd.TcaplusApiSqlRes:
		resp.commonInterface, err = newIndexQueryResponse(pkg)
	default:
		logger.ERR("invalid cmd %d", pkg.Head.Cmd)
		return nil, &terror.ErrorCode{Code: terror.InvalidCmd}
	}

	return resp, err
}

/*
*
按照用户设置的标识位, 对Flag进行置位规则:

0|1|成功场景下设置的2位|失败场景下设置的2位|00

成功:
0: 00 形如  0|1|00|失败场景下设置的2位|00
1: 01 形如  0|1|01|失败场景下设置的2位|00
2: 10 形如  0|1|10|失败场景下设置的2位|00
3: 11 形如  0|1|11|失败场景下设置的2位|00

失败:
0: 00 形如 0|1|成功场景下设置的2位|00|00
1: 01 形如 0|1|成功场景下设置的2位|01|00
2: 10 形如 0|1|成功场景下设置的2位|10|00
3: 11 形如 0|1|成功场景下设置的2位|11|00

*/
func GetResultFlagByBit(flag byte, success bool) int64 {
	var iFirstIndex uint32
	var iSecondIndex uint32
	if success {
		iFirstIndex = 5
		iSecondIndex = 4
	} else {
		iFirstIndex = 3
		iSecondIndex = 2
	}

	if 0 == (flag&(1<<iFirstIndex)) && 0 == (flag&(1<<iSecondIndex)) {
		// 00
		return tcaplus_protocol_cs.TCaplusValueFlag_NOVALUE
	} else if 0 == (flag&(1<<iFirstIndex)) && 0 != (flag&(1<<iSecondIndex)) {
		// 01
		return tcaplus_protocol_cs.TCaplusValueFlag_SAMEWITHREQUEST
	} else if 0 != (flag&(1<<iFirstIndex)) && 0 == (flag&(1<<iSecondIndex)) {
		// 10
		return tcaplus_protocol_cs.TCaplusValueFlag_ALLVALUE
	} else if 0 != (flag&(1<<iFirstIndex)) && 0 != (flag&(1<<iSecondIndex)) {
		// 11
		return tcaplus_protocol_cs.TCaplusValueFlag_ALLOLDVALUE
	} else {
		return -1
	}
}

type tcapResponse struct {
	commonInterface
}

func (res *tcapResponse) GetAffectedRecordNum() int32 {
	switch res.commonInterface.(type) {
	case *listDeleteAllResponse:
		return res.commonInterface.(*listDeleteAllResponse).GetAffectedRecordNum()
	case *listDeleteBatchResponse:
		return res.commonInterface.(*listDeleteBatchResponse).GetAffectedRecordNum()
	default:
		return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
	}
}

func (res *tcapResponse) GetTotalNum() int {
	switch res.commonInterface.(type) {
	case *updataByPartKeyResponse:
		return res.commonInterface.(*updataByPartKeyResponse).GetTotalNum()
	case *deleteByPartKeyResponse:
		return res.commonInterface.(*deleteByPartKeyResponse).GetTotalNum()
	default:
		return 0
	}
}

func (res *tcapResponse) GetFailedNum() int {
	switch res.commonInterface.(type) {
	case *updataByPartKeyResponse:
		return res.commonInterface.(*updataByPartKeyResponse).GetFailedNum()
	case *deleteByPartKeyResponse:
		return res.commonInterface.(*deleteByPartKeyResponse).GetFailedNum()
	default:
		return 0
	}
}

func (res *tcapResponse) FetchErrorRecord() (*record.Record, error) {
	switch res.commonInterface.(type) {
	case *updataByPartKeyResponse:
		return res.commonInterface.(*updataByPartKeyResponse).FetchErrorRecord()
	case *deleteByPartKeyResponse:
		return res.commonInterface.(*deleteByPartKeyResponse).FetchErrorRecord()
	default:
		return nil, nil
	}
}

func (res *tcapResponse) FetchSqlResult() (*sqlResult, error) {
	switch res.commonInterface.(type) {
	case *indexQueryResponse:
		return res.commonInterface.(*indexQueryResponse).FetchSqlResult()
	default:
		return nil, &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH}
	}
}

func (res *tcapResponse) ProcAggregationSqlQueryType() ([]string, error) {
	switch res.commonInterface.(type) {
	case *indexQueryResponse:
		return res.commonInterface.(*indexQueryResponse).ProcAggregationSqlQueryType()
	default:
		return nil, &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH}
	}
}

func (res *tcapResponse) GetSqlType() int {
	switch res.commonInterface.(type) {
	case *indexQueryResponse:
		return res.commonInterface.(*indexQueryResponse).GetSqlType()
	default:
		return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
	}
}
