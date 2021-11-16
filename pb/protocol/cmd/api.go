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
	//TCAPLUS_API_LIST_TABLE_TRAVERSE_REQ = 0x0057
	//
	///**\brief List table的遍历响应*/
	//TCAPLUS_API_LIST_TABLE_TRAVERSE_RES = 0x0058

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
	//TCAPLUS_API_PB_BATCH_FIELD_GET_REQ           = 0x0075
	//
	///** \brief protobuf部分字段自增响应 */
	//TCAPLUS_API_PB_BATCH_FIELD_GET_RES           = 0x0076

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
