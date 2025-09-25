package cmd

//brief 操作类型定义
const (
	TcaplusApiInvalidReq = 0x0000

	/** \brief 无效的应答 */
	TcaplusApiInvalidRes = -0x0001

	/** \brief 插入请求 */
	TcaplusApiInsertReq = 0x0001

	/** \brief 插入应答 */
	TcaplusApiInsertRes = 0x0002

	/** \brief 替换/插入请求 */
	TcaplusApiReplaceReq = 0x0003

	/** \brief 替换/插入应答 */
	TcaplusApiReplaceRes = 0x0004

	///** \brief 增量更新请求 */
	TcaplusApiIncreaseReq = 0x0005

	///** \brief 增量更新应答 */
	TcaplusApiIncreaseRes = 0x0006

	/** \brief 单条查询请求 */
	TcaplusApiGetReq = 0x0007

	/** \brief 单条查询应答 */
	TcaplusApiGetRes = 0x0008

	/** \brief 删除请求 */
	TcaplusApiDeleteReq = 0x0009

	/** \brief 删除应答 */
	TcaplusApiDeleteRes = 0x000a

	///** \brief 查询List所有元素请求 */
	TcaplusApiListGetAllReq = 0x000b
	//
	///** \brief 查询List所有元素应答 */
	TcaplusApiListGetAllRes = 0x000c

	/** \brief 删除List所有元素请求 */
	TcaplusApiListDeleteAllReq = 0x000d

	/** \brief 删除List所有元素应答 */
	TcaplusApiListDeleteAllRes = 0x000e

	/** \brief 删除List多个元素请求 */
	TcaplusApiListDeleteBatchReq = 0x0041

	/** \brief 删除List多个元素应答 */
	TcaplusApiListDeleteBatchRes = 0x0042

	/** \brief 查询List单个元素请求 */
	TcaplusApiListGetReq = 0x000f

	/** \brief 查询List单个元素应答 */
	TcaplusApiListGetRes = 0x0010

	/** \brief 插入List元素请求 */
	TcaplusApiListAddAfterReq = 0x0011

	/** \brief 插入List元素应答 */
	TcaplusApiListAddAfterRes = 0x0012

	/** \brief 删除List单个元素请求 */
	TcaplusApiListDeleteReq = 0x0013

	/** \brief 删除List单个元素应答 */
	TcaplusApiListDeleteRes = 0x0014

	/** \brief 替换List单个元素请求 */
	TcaplusApiListReplaceReq = 0x0015

	/** \brief 替换List单个元素应答 */
	TcaplusApiListReplaceRes = 0x0016

	/** \brief 批量查询请求 */
	TcaplusApiBatchGetReq = 0x0017

	/** \brief 批量查询应答 */
	TcaplusApiBatchGetRes = 0x0018

	// Generic批量插入请求
	TcaplusApiBatchInsertReq = 0x0091

	// Generic批量插入响应
	TcaplusApiBatchInsertRes = 0x0092

	// Generic批量替换请求
	TcaplusApiBatchReplaceReq = 0x0093

	// Generic批量替换响应
	TcaplusApiBatchReplaceRes = 0x0094

	// Generic批量更新请求
	TcaplusApiBatchUpdateReq = 0x0095

	// Generic批量更新响应
	TcaplusApiBatchUpdateRes = 0x0096

	// Generic批量删除请求
	TcaplusApiBatchDeleteReq = 0x0097

	// Generic批量删除响应
	TcaplusApiBatchDeleteRes = 0x0098

	/** \brief List批量查询 */
	TcaplusApiListGetBatchReq = 0x0099
	TcaplusApiListGetBatchRes = 0x009a

	/** \brief List批量插入 */
	TcaplusApiListAddAfterBatchReq = 0x009b
	TcaplusApiListAddAfterBatchRes = 0x009c

	/** \brief List批量替换 */
	TcaplusApiListReplaceBatchReq = 0x009d
	TcaplusApiListReplaceBatchRes = 0x009e

	/** \brief 部分Key查询请求 */
	TcaplusApiGetByPartkeyReq = 0x0019

	/** \brief 部分Key查询应答 */
	TcaplusApiGetByPartkeyRes = 0x001a

	/** \brief 更新请求 */
	TcaplusApiUpdateReq = 0x001d

	/** \brief 更新应答 */
	TcaplusApiUpdateRes = 0x001e

	TcaplusApiMetadataGetReq = 0x001b

	TcaplusApiMetadataGetRes = 0x001c

	// 服务化应用身份认证请求
	TcaplusApiAppSignUpReq = 51

	// 服务化应用身份认证响应
	TcaplusApiAppSignUpRes = 52

	// 心跳检查请求
	TcaplusApiHeartBeatReq = 53

	// 心跳检查响应
	TcaplusApiHeartBeatRes = 54

	// tcaproxy通知客户端即将停止运行
	TcaplusApiNotifyStopReq = 67

	// 客户端响应tcaproxy，表示暂时不再发送请求
	TcaplusApiNotifyStopRes = 68

	/** \brief 表遍历请求 */
	TcaplusApiTableTraverseReq = 0x0045

	/** \brief 表遍历响应 */
	TcaplusApiTableTraverseRes = 0x0046

	/** \brief 表遍历前获取shard list请求 */
	TcaplusApiGetShardListReq = 0x0047

	/** \brief 表遍历前获取shard list响应 */
	TcaplusApiGetShardListRes = 0x0048

	///** \brief 批量Partkey查询请求 */
	//TCAPLUS_API_BATCH_GET_BY_PARTKEY_REQ           = 0x0049
	//
	///** \brief 批量Partkey查询响应 */
	//TCAPLUS_API_BATCH_GET_BY_PARTKEY_RES          = 0x004a
	//
	///** \brief Document 操作请求 */
	//TCAPLUS_API_DOCUMENT_OPERATION_REQ            = 0x004b
	//
	///** \brief Document 操作响应 */
	//TCAPLUS_API_DOCUMENT_OPERATION_RES            = 0x004c
	//
	///** \brief Partkey update请求 */
	TcaplusApiUpdateByPartkeyReq = 0x004d
	//
	///** \brief Partkey update响应 */
	TcaplusApiUpdateByPartkeyRes = 0x004e
	//
	///** \brief Partkey delete请求 */
	TcaplusApiDeleteByPartkeyReq = 0x004f
	//
	///** \brief Partkey delete响应 */
	TcaplusApiDeleteByPartkeyRes = 0x0050
	//
	///** \brief 带有相同Partkey的批量insert请求*/
	//TCAPLUS_API_INSERT_BY_PARTKEY_REQ          = 0x0051
	//
	///** \brief 带有相同Partkey的批量insert响应 */
	//TCAPLUS_API_INSERT_BY_PARTKEY_RES          = 0x0052
	//
	/** \brief table的记录总数请求 */
	TcaplusApiGetTableRecordCountReq = 0x0053

	/** \brief table的记录总数响应 */
	TcaplusApiGetTableRecordCountRes = 0x0054
	//
	///**\brief List table的遍历请求*/
	TcaplusApiListTableTraverseReq = 0x0057
	//
	///**\brief List table的遍历响应*/
	TcaplusApiListTableTraverseRes = 0x0058

	/** \brief 设置记录的ttl请求 */
	TcaplusApiSetTtlReq = 0x0059

	/** \brief 设置记录的ttl响应 */
	TcaplusApiSetTtlRes = 0x005a

	/** \brief 获取记录的ttl请求 */
	TcaplusApiGetTtlReq = 0x005b

	/** \brief 获取记录的ttl响应 */
	TcaplusApiGetTtlRes = 0x005c

	/** \brief protobuf部分字段获取请求 */
	TcaplusApiPBFieldGetReq = 0x0067

	/** \brief protobuf部分字段获取响应 */
	TcaplusApiPBFieldGetRes = 0x0068

	/** \brief protobuf部分字段更新请求 */
	TcaplusApiPBFieldUpdateReq = 0x0069

	/** \brief protobuf部分字段更新响应 */
	TcaplusApiPBFieldUpdateRes = 0x006a

	/** \brief protobuf部分字段自增请求 */
	TcaplusApiPBFieldIncreaseReq = 0x006b

	/** \brief protobuf部分字段自增响应 */
	TcaplusApiPBFieldIncreaseRes = 0x006c

	///** \brief protobuf部分字段自增请求 */
	TcaplusApiPBBatchFieldGetReq = 0x0075
	//
	///** \brief protobuf部分字段自增响应 */
	TcaplusApiPBBatchFieldGetRes = 0x0076

	/** \brief 索引查询请求 */
	TcaplusApiSqlReq = 0x0081

	/** \brief 索引查询响应 */
	TcaplusApiSqlRes = 0x0082

	/**\brief API的最大值，为了能够匹配系统内部请求*/
	TcaplusApiMaxNum = 0xffff
)

const (
	//自增
	TcaplusApiOpPlus = 1
	//自减
	TcaplusApiOpMinus = 2
)

var TypeMap = [TcaplusApiMaxNum]int8{
	//1表示简单请求
	TcaplusApiInsertReq:          1,
	TcaplusApiReplaceReq:         1,
	TcaplusApiIncreaseReq:        1,
	TcaplusApiGetReq:             1,
	TcaplusApiDeleteReq:          1,
	TcaplusApiListGetReq:         1,
	TcaplusApiListAddAfterReq:    1,
	TcaplusApiListDeleteReq:      1,
	TcaplusApiListReplaceReq:     1,
	TcaplusApiUpdateReq:          1,
	TcaplusApiPBFieldGetReq:      1,
	TcaplusApiPBFieldUpdateReq:   1,
	TcaplusApiPBFieldIncreaseReq: 1,

	//2是简单响应
	TcaplusApiInsertRes:          2,
	TcaplusApiReplaceRes:         2,
	TcaplusApiIncreaseRes:        2,
	TcaplusApiGetRes:             2,
	TcaplusApiDeleteRes:          2,
	TcaplusApiListGetRes:         2,
	TcaplusApiListAddAfterRes:    2,
	TcaplusApiListDeleteRes:      2,
	TcaplusApiListReplaceRes:     2,
	TcaplusApiUpdateRes:          2,
	TcaplusApiPBFieldGetRes:      2,
	TcaplusApiPBFieldUpdateRes:   2,
	TcaplusApiPBFieldIncreaseRes: 2,

	// 3是复杂请求
	TcaplusApiListGetAllReq:        3,
	TcaplusApiListDeleteAllReq:     3,
	TcaplusApiListDeleteBatchReq:   3,
	TcaplusApiBatchGetReq:          3,
	TcaplusApiBatchInsertReq:       3,
	TcaplusApiBatchReplaceReq:      3,
	TcaplusApiBatchUpdateReq:       3,
	TcaplusApiBatchDeleteReq:       3,
	TcaplusApiListGetBatchReq:      3,
	TcaplusApiListAddAfterBatchReq: 3,
	TcaplusApiListReplaceBatchReq:  3,
	TcaplusApiGetByPartkeyReq:      3,
	TcaplusApiUpdateByPartkeyReq:   3,
	TcaplusApiDeleteByPartkeyReq:   3,
	TcaplusApiSetTtlReq:            3,
	TcaplusApiGetTtlReq:            3,
	TcaplusApiPBBatchFieldGetReq:   3,
	TcaplusApiSqlReq:               3,

	//4是复杂响应
	TcaplusApiListGetAllRes:        4,
	TcaplusApiListDeleteAllRes:     4,
	TcaplusApiListDeleteBatchRes:   4,
	TcaplusApiBatchGetRes:          4,
	TcaplusApiBatchInsertRes:       4,
	TcaplusApiBatchReplaceRes:      4,
	TcaplusApiBatchUpdateRes:       4,
	TcaplusApiBatchDeleteRes:       4,
	TcaplusApiListGetBatchRes:      4,
	TcaplusApiListAddAfterBatchRes: 4,
	TcaplusApiListReplaceBatchRes:  4,
	TcaplusApiGetByPartkeyRes:      4,
	TcaplusApiUpdateByPartkeyRes:   4,
	TcaplusApiDeleteByPartkeyRes:   4,
	TcaplusApiSetTtlRes:            4,
	TcaplusApiGetTtlRes:            4,
	TcaplusApiPBBatchFieldGetRes:   4,
	TcaplusApiSqlRes:               4,
} // 0表示默认不用关注

func IsSimpleReq(cmd uint32) bool {
	if cmd >= TcaplusApiMaxNum {
		return false
	}
	return TypeMap[cmd] == 1
}

func IsSimpleRes(cmd uint32) bool {
	if cmd >= TcaplusApiMaxNum {
		return false
	}
	return TypeMap[cmd] == 2
}

func IsComplexReq(cmd uint32) bool {
	if cmd >= TcaplusApiMaxNum {
		return false
	}
	return TypeMap[cmd] == 3
}

func IsComplexRes(cmd uint32) bool {
	if cmd >= TcaplusApiMaxNum {
		return false
	}
	return TypeMap[cmd] == 4
}
