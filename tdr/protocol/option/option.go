package option

import "time"

const (
	CheckDataVersionAutoIncrease   byte = 1
	NoCheckDataVersionOverwrite    byte = 2
	NoCheckDataVersionAutoIncrease byte = 3

	TcaplusListShiftNone byte = 0
	TcaplusListShiftHead byte = 1
	TcaplusListShiftTail byte = 2

	TcaplusAddableIncreaseTrue  byte = 1
	TcaplusAddableIncreaseFalse byte = 0

	TcaplusResultFlagNoValue         byte = 0
	TcaplusResultFlagSameWithRequest byte = 1
	TcaplusResultFlagAllNewValue     byte = 2
	TcaplusResultFlagAllOldValue     byte = 3

	TcaplusFlagFetchOnlyIfModified              int32 = 1
	TcaplusFlagFetchOnlyIfExpired               int32 = 2
	TcaplusFlagOnlyReadFromSlave                int32 = 4
	TcaplusFlagListReserveIndexHavingNoElements int32 = 8
	TcaplusFlagInsertRecordIfNotExist           int32 = 16 //PB的FieldUpdate使用，数据不存在则插入
)

/** brief  自增自减字段操作
 *  param  [in] field_name         字段名称
 *  param  [in] incData            加减数值，和表中定义的字段类型保持一致
 *  param  [in] operation          操作类型，cmd.TcaplusApiOpPlus 加操作 cmd.TcaplusApiOpMinus 减操作
 *  param [IN] lower_limit         操作结果值下限，如果比这个值小，返回 TcapErrCode::SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE
 *  param [IN] upper_limit         操作结果值上限，如果比这个值大，返回 TcapErrCode::SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE
 *  note                           lower_limit == upper_limit 时，存储端不对操作结果进行范围检测
 */
type IncFieldInfo struct {
	FieldName  string      //字段名称
	IncData    interface{} //加减数值，和表中定义的字段类型保持一致,支持double
	Operation  uint32      //操作类型，cmd.TcaplusApiOpPlus 加操作 cmd.TcaplusApiOpMinus 减操作
	LowerLimit int64       // 操作结果值下限
	UpperLimit int64       // 操作结果值上限
}

type TDROpt struct {
	/**
	@brief  更新操作时设置记录版本的检查类型，用于乐观锁，版本检测类型，取值可以为:
			CheckDataVersionAutoIncrease: 表示检测记录版本号，只有当record.SetVersion函数传入的参数
				version的值>0,并且版本号与服务器端的版本号相同时，请求才会成功同时在服务器端该版本号会自增1；
				如果version <=0，则仍然表示不关心版本号
			NoCheckDataVersionOverwrite: 表示不检测记录版本号。
				当record.SetVersion函数传入的参数version的值>0,覆盖服务端的版本号；
				如果Version的version <=0，则仍然表示不关心版本号
			NoCheckDataVersionAutoIncrease: 表示不检测记录版本号，将服务器端的数据记录版本号自增1，
				若服务器端新写入数据记录则新写入的数据记录的版本号为1
	**/
	VersionPolicy byte
	Version       int32
	BatchVersion  []int32 //批量操作， partKey时读取该version
	/**
	@brief  批量操作，单条记录的操作结果
	**/
	BatchResult []error

	/**
	@brief  设置响应标志。主要用于Generic表的insert、increase、replace、update、delete操作, 请求标志:
			TcaplusResultFlagNoValue表示: 只需返回操作执行成功与否
			TcaplusResultFlagSameWithRequest表示: 操作成功，响应返回与请求字段一致
			TcaplusResultFlagAllNewValue表示: 操作成功，响应返回变更记录的所有字段最新数据
			TcaplusResultFlagAllOldValue表示: 操作成功，响应返回变更记录的所有字段旧数据
		ResultFlag 是本次请求执行成功后返回给前端的数据
		ResultFlagForSuccess是本次请求执行成功后返回给前端的数据
		ResultFlagForFail是本次请求执行失败后返回给前端的数据
	NOTE：ResultFlag有历史包袱，某些场景并不准确，推荐使用ResultFlagForSuccess
	**/
	ResultFlag           byte
	ResultFlagForSuccess byte
	ResultFlagForFail    byte

	/**
	  	@brief  设置LIST满时，插入元素时，删除旧元素的模式
			TcaplusListShiftNone: 不允许删除元素，若LIST满，插入失败；
			TcaplusListShiftHead: 移除最前面的元素；
			TcaplusListShiftTail: 移除最后面的元素
			如果表是排序List,必须要进行淘汰,且淘汰规则是根据字段的排序顺序进行自动制定的,用户调用该接口会失败
	*/
	ListShiftFlag byte

	/**
	@brief  设置请求的通用标志位，可以通过"按位或"操作同时设定多个值
			有效的标志位包括：
	*  TcaplusFlagFetchOnlyIfModified:
	*       "数据变更才取回"标志位。在发起读操作之前，用户代码通过 TcaplusServiceRecord::SetVersion()
	*       带上本地缓存数据的版本号，并将此标志置位，那么存储端检测到当前数据与API本地缓存的数据版本
	*       一致时，表明该记录未发生过修改，API缓存的数据是最新的，因此在响应中将不会携带实际的数据，
	*       只是返回 TcapErrCode::COMMON_INFO_DATA_NOT_MODIFIED 的错误码
	*
	*       只有如下请求支持设置此标志：
	*           TCAPLUS_API_GET_REQ,
	*           TCAPLUS_API_LIST_GET_REQ,
	*           TCAPLUS_API_LIST_GETALL_REQ
	*
	*  TCAPLUS_FLAG_FETCH_ONLY_IF_EXPIRED:
	*       "数据过期才取回"标志位。在发起读操作之前，用户代码通过 SetExpireTime() 设定数据过期时间，
	*       并将此标志置位，那么存储端若检测到记录在指定时间内发生过更新，则将数据返回，
	*       否则不返回实际数据，只是返回 TcapErrCode::COMMON_INFO_DATA_NOT_MODIFIED 的错误码。

	*       只有如下请求支持设置此标志：
	*           TCAPLUS_API_BATCH_GET_REQ
	*
	*  TCAPLUS_FLAG_ONLY_READ_FROM_SLAVE
	*       设置此标志后，读请求将会直接发送给Tcapsvr Slave 节点。
	*       Tcapsvr Slave 通常比较空闲，设置此标志有助于充分利用Tcapsvr Slave 资源。
	*
	*       适用场景:
	*                              对于数据实时性要求不高的读请求，
	*                              包括generic表和list表的所有读请求以及batchget，遍历请求
	*
	*  TCAPLUS_FLAG_LIST_RESERVE_INDEX_HAVING_NO_ELEMENTS
	*       设置此标志后，List表删除最后一个元素时需要保留index和version。
	*       ListDelete ListDeleteBatch ListDeleteAll操作在删除list表最后一个元素时，
	*          设置此标志在写入新的List记录时，版本号依次增长，不会被重置为1。
	*
	*       适用场景:
	*                              业务需要确定某个表在删除最后一个元素时是否需要保留index和version
	*                              主要涉及List表的使用体验
	*
	*/
	Flags int32

	/**
	  @brief  如果此请求会返回多条记录，通过此接口对返回的记录做一些限制
	  @param [IN] limit       需要查询的记录条数, limit若等于-1表示操作或返回所有匹配的数据记录.
	  @param [IN] offset      记录起始编号；若设置为负值(-N, N>0)，则从倒数第N个记录开始返回结果
	  @retval 0               设置成功
	  @retval <0              设置失败，具体错误参见 \link ErrorCode \endlink
	  @note 对于Generic类型的部分Key查询，limit表示所要获取Record的条数，offset表示所要获取Record的开始下标；
		   	对于List类型的GetAll操作，limit表示所要获取Record的条数，offset表示所要获取Record的开始下标，
			在当前版本中这些Record一定属于同一个List.
		   	该函数仅仅对于GET_BY_PARTKEY(Generic类型的部分Key查询), UPDATE_BY_PARTKEY,
			DELETE_BY_PARTKEY, LIST_GETALL(List类型的GetAll操)这4种操作类型有效。
	*/
	Limit  int32
	Offset int32

	/**
	  @brief  设置是否允许一个请求包可以自动响应多个应答包，仅对ListGetAll和BatchGet协议有效。
	  @param [IN] multi_flag   多响应包标示，1表示允许一个请求包可以自动响应多个应答包,
							   0表示不允许一个请求包自动响应多个应答包
	  @retval 0                设置成功
	  @retval <0               设置失败，具体错误参见 \link ErrorCode \endlink
	  @note	分包应答，目前只支持ListGetAll和BatchGet操作；其他操作设置该值是没有意义的，
			函数会返回<0的错误码。
	*/
	MultiFlag byte
	/**
	@brief  tdr表设置需要查询或更新的Value字段名称列表，即部分Value字段查询和更新，可用于get、replace、update操作。
			需要查询或更新的字段名称列表
	**/
	FieldNames []string

	/**
	@brief  自增自减字段设置，IncreaseReq使用
	**/
	IncField []IncFieldInfo

	/**
	@brief  超时时间,不设置默认5s
	**/
	Timeout time.Duration
	/**
	@brief  设置空记录自增允许标志。用于Generic表的increase操作, 空记录自增允许标志。
			TcaplusAddableIncreaseFalse表示不允许,TcaplusAddableIncreaseTrue表示允许，
			当记录不存在时，将按字段默认值创建新记录再自增；若无默认值则返回错误
	*/
	AddableIncreaseFlag byte

	//同步请求意义不大，设置用户缓存，响应消息将携带返回
	UserBuff []byte
	//同步请求意义不大，设置请求的异步事务ID
	AsyncId uint64
}

type PBOpt struct {
	/**
	@brief  更新操作时设置记录版本的检查类型，用于乐观锁，版本检测类型，取值可以为:
			CheckDataVersionAutoIncrease: 表示检测记录版本号，只有当record.SetVersion函数传入的参数
				version的值>0,并且版本号与服务器端的版本号相同时，请求才会成功同时在服务器端该版本号会自增1；
				如果version <=0，则仍然表示不关心版本号
			NoCheckDataVersionOverwrite: 表示不检测记录版本号。
				当record.SetVersion函数传入的参数version的值>0,覆盖服务端的版本号；
				如果Version的version <=0，则仍然表示不关心版本号
			NoCheckDataVersionAutoIncrease: 表示不检测记录版本号，将服务器端的数据记录版本号自增1，
				若服务器端新写入数据记录则新写入的数据记录的版本号为1
	**/
	VersionPolicy byte
	Version       int32
	BatchVersion  []int32 //批量操作， partKey时读取该version
	/**
	@brief  批量操作，单条记录的操作结果
	**/
	BatchResult []error

	/**
	@brief  设置响应标志。主要用于Generic表的insert、increase、replace、update、delete操作, 请求标志:
			TcaplusResultFlagNoValue表示: 只需返回操作执行成功与否
			TcaplusResultFlagSameWithRequest表示: 操作成功，响应返回与请求字段一致
			TcaplusResultFlagAllNewValue表示: 操作成功，响应返回变更记录的所有字段最新数据
			TcaplusResultFlagAllOldValue表示: 操作成功，响应返回变更记录的所有字段旧数据
		ResultFlagForSuccess是本次请求执行成功后返回给前端的数据
		ResultFlagForFail是本次请求执行失败后返回给前端的数据
	**/
	ResultFlag           byte
	ResultFlagForSuccess byte
	ResultFlagForFail    byte

	/**
	  	@brief  设置LIST满时，插入元素时，删除旧元素的模式
			TcaplusListShiftNone: 不允许删除元素，若LIST满，插入失败；
			TcaplusListShiftHead: 移除最前面的元素；
			TcaplusListShiftTail: 移除最后面的元素
			如果表是排序List,必须要进行淘汰,且淘汰规则是根据字段的排序顺序进行自动制定的,用户调用该接口会失败
	*/
	ListShiftFlag byte

	/**
	@brief  设置请求的通用标志位，可以通过"按位或"操作同时设定多个值
			有效的标志位包括：
	*  TcaplusFlagFetchOnlyIfModified:
	*       "数据变更才取回"标志位。在发起读操作之前，用户代码通过 TcaplusServiceRecord::SetVersion()
	*       带上本地缓存数据的版本号，并将此标志置位，那么存储端检测到当前数据与API本地缓存的数据版本
	*       一致时，表明该记录未发生过修改，API缓存的数据是最新的，因此在响应中将不会携带实际的数据，
	*       只是返回 TcapErrCode::COMMON_INFO_DATA_NOT_MODIFIED 的错误码
	*
	*       只有如下请求支持设置此标志：
	*           TCAPLUS_API_GET_REQ,
	*           TCAPLUS_API_LIST_GET_REQ,
	*           TCAPLUS_API_LIST_GETALL_REQ
	*
	*  TCAPLUS_FLAG_FETCH_ONLY_IF_EXPIRED:
	*       "数据过期才取回"标志位。在发起读操作之前，用户代码通过 SetExpireTime() 设定数据过期时间，
	*       并将此标志置位，那么存储端若检测到记录在指定时间内发生过更新，则将数据返回，
	*       否则不返回实际数据，只是返回 TcapErrCode::COMMON_INFO_DATA_NOT_MODIFIED 的错误码。

	*       只有如下请求支持设置此标志：
	*           TCAPLUS_API_BATCH_GET_REQ
	*
	*  TCAPLUS_FLAG_ONLY_READ_FROM_SLAVE
	*       设置此标志后，读请求将会直接发送给Tcapsvr Slave 节点。
	*       Tcapsvr Slave 通常比较空闲，设置此标志有助于充分利用Tcapsvr Slave 资源。
	*
	*       适用场景:
	*                              对于数据实时性要求不高的读请求，
	*                              包括generic表和list表的所有读请求以及batchget，遍历请求
	*
	*  TCAPLUS_FLAG_LIST_RESERVE_INDEX_HAVING_NO_ELEMENTS
	*       设置此标志后，List表删除最后一个元素时需要保留index和version。
	*       ListDelete ListDeleteBatch ListDeleteAll操作在删除list表最后一个元素时，
	*          设置此标志在写入新的List记录时，版本号依次增长，不会被重置为1。
	*
	*       适用场景:
	*                              业务需要确定某个表在删除最后一个元素时是否需要保留index和version
	*                              主要涉及List表的使用体验
	*
	*/
	Flags int32

	/**
	  @brief  如果此请求会返回多条记录，通过此接口对返回的记录做一些限制
	  @param [IN] limit       需要查询的记录条数, limit若等于-1表示操作或返回所有匹配的数据记录.
	  @param [IN] offset      记录起始编号；若设置为负值(-N, N>0)，则从倒数第N个记录开始返回结果
	  @retval 0               设置成功
	  @retval <0              设置失败，具体错误参见 \link ErrorCode \endlink
	  @note 对于Generic类型的部分Key查询，limit表示所要获取Record的条数，offset表示所要获取Record的开始下标；
		   	对于List类型的GetAll操作，limit表示所要获取Record的条数，offset表示所要获取Record的开始下标，
			在当前版本中这些Record一定属于同一个List.
		   	该函数仅仅对于GET_BY_PARTKEY(Generic类型的部分Key查询), UPDATE_BY_PARTKEY,
			DELETE_BY_PARTKEY, LIST_GETALL(List类型的GetAll操)这4种操作类型有效。
	*/
	Limit  int32
	Offset int32

	/**
	  @brief  设置是否允许一个请求包可以自动响应多个应答包，仅对ListGetAll和BatchGet协议有效。
	  @param [IN] multi_flag   多响应包标示，1表示允许一个请求包可以自动响应多个应答包,
							   0表示不允许一个请求包自动响应多个应答包
	  @retval 0                设置成功
	  @retval <0               设置失败，具体错误参见 \link ErrorCode \endlink
	  @note	分包应答，目前只支持ListGetAll和BatchGet操作；其他操作设置该值是没有意义的，
			函数会返回<0的错误码。
	*/
	MultiFlag byte

	/**
	@brief  超时时间,不设置默认5s
	**/
	Timeout time.Duration

	/**
	@brief  PB FieldGet和FieldSet使用，获取或更新部分字段
	**/
	FieldNames []string
	//同步请求意义不大，设置用户缓存，响应消息将携带返回
	UserBuff []byte
	//同步请求意义不大，设置请求的异步事务ID
	AsyncId uint64
}
