package policy

const (
	// 检测记录版本号,只有当该版本号与服务器端的版本号相同时，该版本号才会自增
	CheckDataVersionAutoIncrease uint8 = 1

	//不检测记录版本号，强制把客户端的记录版本号写入到服务器中
	NoCheckDataVersionOverwrite uint8 = 2

	// 不检测记录版本号，将服务器端的版本号自增
	NoCheckDataVersionAutoIncrease uint8 = 3
)

/** \brief 索引查询类型 */
const (
	INVALID_SQL_TYPE = 0 //非法查询
	RECORD_SQL_QUERY_TYPE = 1 //记录查询, select * from test where XXX
	AGGREGATIONS_SQL_QUERY_TYPE = 2 //聚合查询, select sum(level) from test where XXX
)

/** \brief 字段类型，目前主要用于索引查询 */
const (
	TYPE_INVALID = 0
	TYPE_BOOL = 1
	TYPE_INT8 = 2
	TYPE_UINT8 = 3
	TYPE_INT16 = 4
	TYPE_UINT16 = 5
	TYPE_INT32 = 6
	TYPE_UINT32 = 7
	TYPE_INT64 = 8
	TYPE_UINT64 = 9
	TYPE_FLOAT = 10
	TYPE_DOUBLE = 11
	TYPE_STRING = 12
	TYPE_END = 13
)
