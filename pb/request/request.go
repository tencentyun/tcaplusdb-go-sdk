package request

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cs_pool"
	"hash/crc32"
	"sort"

	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	tcaplusCmd "github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

/*
	大多数响应都会用到的函数放到commonInterface接口中
	个别或者极少数响应的特殊方法放到TcaplusRequest
*/
type TcaplusRequest interface {
	commonInterface

	/**
	    @brief  设置空记录自增允许标志。用于Generic表的increase操作。
	    @param  [IN] increase_flag  空记录自增允许标志。
		0表示不允许。1表示允许，当记录不存在时，将按字段默认值创建新记录再自增；若无默认值则返回错误
	    @retval 0    设置成功
	    @retval <0   失败，返回对应的错误码。通常因为未初始化。
	*/
	SetAddableIncreaseFlag(increase_flag byte) int32

	/**
	  @brief  设置LIST满时，插入元素时，删除旧元素的模式
	  @param  [in] chListShiftFlag
				TCAPLUS_LIST_SHIFT_NONE: 不允许删除元素，若LIST满，插入失败；
				TCAPLUS_LIST_SHIFT_HEAD: 移除最前面的元素；
				TCAPLUS_LIST_SHIFT_TAIL: 移除最后面的元素
	          如果表是排序List,必须要进行淘汰,且淘汰规则是根据字段的排序顺序进行自动制定的,用户调用该接口会失败
	  @retval 0              设置成功
	  @retval 非0            设置失败，具体错误参见 \link ErrorCode \endlink
	*/
	SetListShiftFlag(shiftFlag byte) int32

	/**
		@brief  添加LIST记录的元素索引值。该函数只对于
	                    TCAPLUS_API_LIST_DELETE_BATCH_REQ
	                    TCAPLUS_API_LIST_GET_BATCH_REQ
	                    有效，对于其它Command是无效的。
		@param  [in] idx          LIST元素索引值。不可取值TCAPLUS_API_LIST_PRE_FIRST_INDEX，不可取值TCAPLUS_API_LIST_LAST_INDEX。
		@retval 0                 设置成功
		@retval 非0               设置失败，具体错误参见 \link ErrorCode \endlink
	*/
	AddElementIndex(idx int32) int32

	/*
		@brief  添加LIST记录的元素索引值。该函数只对于
						TcaplusApiSqlReq有效
		@param  query sql语句
		@retval 0                 设置成功
		@retval 非0               设置失败，具体错误参见 \link ErrorCode \endlink
	*/
	SetSql(query string) int

	GetTcaplusPackagePtr() *tcaplus_protocol_cs.TCaplusPkg
}

/*
	大多数响应都会用到的函数放到commonInterface接口中
	个别或者极少数响应的特殊方法放到TcaplusRequest
*/
type commonInterface interface {
	/**
	  @brief  向请求中添加一条记录。
	  @param [IN] index         用于List操作,通常>=0，表示该Record在所属List中的Index；
								对于Generic操作，index无意义，设0即可
	  @retval record.Record     返回记录指针
	  @retval error   			错误码
	**/
	AddRecord(index int32) (*record.Record, error)

	/**
	@brief  设置请求的异步事务ID，api会将其值不变地通过对应的响应消息带回来
	@param  [IN] asyncId  请求对应的异步事务ID
	**/
	SetAsyncId(id uint64)

	/**
	@brief  设置记录版本的检查类型，用于乐观锁
	@param [IN] type   版本检测类型，取值可以为:
						CheckDataVersionAutoIncrease: 表示检测记录版本号，只有当record.SetVersion函数传入的参数
							version的值>0,并且版本号与服务器端的版本号相同时，请求才会成功同时在服务器端该版本号会自增1；
							如果record.SetVersion的version <=0，则仍然表示不关心版本号
						NoCheckDataVersionOverwrite: 表示不检测记录版本号。
							当record.SetVersion函数传入的参数version的值>0,覆盖服务端的版本号；
							如果record.SetVersion的version <=0，则仍然表示不关心版本号
				 		NoCheckDataVersionAutoIncrease: 表示不检测记录版本号，将服务器端的数据记录版本号自增1，
							若服务器端新写入数据记录则新写入的数据记录的版本号为1
	@retval error      错误码
	@note 此函数适合Replace, Update操作
	**/
	SetVersionPolicy(p uint8) error

	/**
	@brief  设置响应标志。主要用于Generic表的insert、increase、replace、update、delete操作。
	@param  [IN] flag  请求标志:
							0表示: 只需返回操作执行成功与否
							1表示: 操作成功，响应返回与请求字段一致
							2表示: 操作成功，响应返回变更记录的所有字段最新数据
							3表示: 操作成功，响应返回变更记录的所有字段旧数据
	@retval error      错误码
	NOTE：SetResultFlag有历史包袱，某些场景并不准确，推荐使用SetResultFlagForSuccess
	**/
	SetResultFlag(flag int) error

	/**
	@brief 设置需要查询或更新的Value字段名称列表，即部分Value字段查询和更新，可用于get、replace、update操作。
	@param [IN] valueNameList   需要查询或更新的字段名称列表
	@retval error      			错误码
	@note  在使用该函数设置字段名时，字段名只能包含value字段名，不能包含key字段名；对于数组类型的字段，
				refer字段和数组字段要同时设置或者同时不设置，否则容易数据错乱
	**/
	SetFieldNames(valueNameList []string) error

	/**
	@brief 设置用户缓存，响应消息将携带返回
	@param [IN] userBuffer   用户缓存
	@retval error      		 错误码
	**/
	SetUserBuff(userBuffer []byte) error

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
	SetResultLimit(limit int32, offset int32) int32

	/**
	  @brief  设置是否允许一个请求包可以自动响应多个应答包，仅对ListGetAll和BatchGet协议有效。
	  @param [IN] multi_flag   多响应包标示，1表示允许一个请求包可以自动响应多个应答包,
	                           0表示不允许一个请求包自动响应多个应答包
	  @retval 0                设置成功
	  @retval <0               设置失败，具体错误参见 \link ErrorCode \endlink
	  @note	分包应答，目前只支持ListGetAll和BatchGet操作；其他操作设置该值是没有意义的，
	        函数会返回<0的错误码。
	*/
	SetMultiResponseFlag(multi_flag byte) int32

	/**
		@brief	设置响应标志。主要是本次请求成功执行后返回给前端的数据

		result_flag 的取值范围如下:

	 TCaplusValueFlag_NOVALUE = 0,			  // 不返回任何返回值
	 TCaplusValueFlag_SAMEWITHREQUEST = 1,	  // 返回同请求一致的值
	 TCaplusValueFlag_ALLVALUE = 2, 		  // 返回tcapsvr端操作后所有字段的值
	 TCaplusValueFlag_ALLOLDVALUE = 3,		  // 返回tcapsvr端操作前所有字段的值


	下面是各个支持的命令字在设置不同的result_flag下执行成功后返回给API端的数据详细情况:

	 1. TCAPLUS_API_INSERT_REQ TCAPLUS_API_BATCH_INSERT_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后返回本次insert操作后的数据
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回空数据

	 2. TCAPLUS_API_REPLACE_REQ TCAPLUS_API_BATCH_REPLACE_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后返回本次replace操作后的数据
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回tcapsvr端操作前的数据, 如果tcapsvr端没有数据,即返回为空

	 3. TCAPLUS_API_UPDATE_REQ TCAPLUS_API_BATCH_UPDATE_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后返回本次update操作后的数据
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回tcapsvr端操作前的数据

	 4. TCAPLUS_API_INCREASE_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后返回本次increase操作后的数据
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回tcapsvr端操作前的数据, 如果tcapsvr端没有数据,即返回为空

	 5. TCAPLUS_API_DELETE_REQ TCAPLUS_API_BATCH_DELETE_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后返回空数据
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回tcapsvr端操作前的数据

	 6. TCAPLUS_API_LIST_DELETE_BATCH_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据, 暂时没有实现
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回tcapsvr端操作前的数据, 凡是本次成功删除的index对应的数据都会返回

	 7. TCAPLUS_API_LIST_ADDAFTER_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据, 暂时没有实现
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后, 返回本次插入的记录和本次淘汰的数据记录
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后不返回数据

	 8. TCAPLUS_API_LIST_DELETE_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据, 暂时没有实现
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后返回空数据
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回tcapsvr端listdelete前的数据

	 9. TCAPLUS_API_LIST_REPLACE_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 则操作成功后不返回数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据, 暂时没有实现
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后返回tcapsvr端listreplace后的数据
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回tcapsvr端listreplace前的数据
	10. TCAPLUS_API_LIST_REPLACE_BATCH_REQ
		 如果设置的是TCaplusValueFlag_NOVALUE, 操作成功后返回和请求一致的数据
		 如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作成功后返回和请求一致的数据
		 如果设置的是TCaplusValueFlag_ALLVALUE, 则操作成功后返回tcapsvr端listreplace后的数据
		 如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回tcapsvr端listreplace前的数据

	 @param  [IN] result_flag  请求标志:
								 0表示: 只需返回操作执行成功与否
								 1表示: 返回与请求字段一致
								 2表示: 须返回变更记录的所有字段最新数据
								 3表示: 须返回变更记录的所有字段旧数据

								 对于batch_get请求，该字段设置为大于0时，某个key查询记录不存在或svr端产生的其它错误时会返回对应的key，
								 从而知道是哪个key对应的记录失败了
	 @retval 0	  设置成功
	 @retval <0   失败，返回对应的错误码。通常因为未初始化。

	*/

	SetResultFlagForSuccess(result_flag byte) int

	/**
		@brief	设置响应标志。主要是本次请求执行失败后返回给前端的数据

		result_flag 的取值范围如下:

		TCaplusValueFlag_NOVALUE = 0,			 // 不返回任何返回值
		TCaplusValueFlag_SAMEWITHREQUEST = 1,	 // 返回同请求一致的值
		TCaplusValueFlag_ALLVALUE = 2,			 // 返回tcapsvr端操作后所有字段的值
		TCaplusValueFlag_ALLOLDVALUE = 3,		 // 返回tcapsvr端操作前所有字段的值


	   下面是各个支持的命令字在设置不同的result_flag下执行失败后返回给API端的数据详细情况:

		1. TCAPLUS_API_INSERT_REQ  TCAPLUS_API_BATCH_INSERT_REQ
			如果设置的是TCaplusValueFlag_NOVALUE, 则操作失败后不返回数据
			如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作失败后返回和请求一致的数据
			如果设置的是TCaplusValueFlag_ALLVALUE, 不合理场景
			如果设置的是TCaplusValueFlag_ALLOLDVALUE, 如果获取到了tcapsvr端的数据则返回tcpasvr端的数据,如果没有获取到tcapsvr端的数据则返回空

		2. TCAPLUS_API_REPLACE_REQ  TCAPLUS_API_BATCH_REPLACE_REQ
			如果设置的是TCaplusValueFlag_NOVALUE, 则操作失败后不返回数据
			如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作失败后返回和请求一致的数据
			如果设置的是TCaplusValueFlag_ALLVALUE, 不合理场景
			如果设置的是TCaplusValueFlag_ALLOLDVALUE, 如果获取到了tcapsvr端的数据则返回tcpasvr端的数据,如果没有获取到tcapsvr端的数据则返回空

		3. TCAPLUS_API_UPDATE_REQ  TCAPLUS_API_BATCH_UPDATE_REQ
			如果设置的是TCaplusValueFlag_NOVALUE, 则操作失败后不返回数据
			如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作失败后返回和请求一致的数据
			如果设置的是TCaplusValueFlag_ALLVALUE, 不合理场景
			如果设置的是TCaplusValueFlag_ALLOLDVALUE, 如果获取到了tcapsvr端的数据则返回tcpasvr端的数据,如果没有获取到tcapsvr端的数据则返回空

		4. TCAPLUS_API_INCREASE_REQ
			如果设置的是TCaplusValueFlag_NOVALUE, 则操作失败后不返回数据
			如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作失败后返回和请求一致的数据
			如果设置的是TCaplusValueFlag_ALLVALUE, 不合理场景
			如果设置的是TCaplusValueFlag_ALLOLDVALUE, 如果获取到了tcapsvr端的数据则返回tcpasvr端的数据,如果没有获取到tcapsvr端的数据则返回空

		5. TCAPLUS_API_DELETE_REQ TCAPLUS_API_BATCH_DELETE_REQ
			如果设置的是TCaplusValueFlag_NOVALUE, 则操作失败后不返回数据
			如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作失败后返回和请求一致的数据
			如果设置的是TCaplusValueFlag_ALLVALUE, 不合理场景
			如果设置的是TCaplusValueFlag_ALLOLDVALUE, 如果获取到了tcapsvr端的数据则返回tcpasvr端的数据,如果没有获取到tcapsvr端的数据则返回空

		6. TCAPLUS_API_LIST_DELETE_BATCH_REQ
			如果设置的是TCaplusValueFlag_NOVALUE, 则操作失败后不返回数据
			如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作失败后返回和请求一致的数据, 暂时没有实现
			如果设置的是TCaplusValueFlag_ALLVALUE, 不合理场景
			如果设置的是TCaplusValueFlag_ALLOLDVALUE, 则操作成功后返回tcapsvr端操作前的数据, 凡是本次成功删除的index对应的数据都会返回

		7. TCAPLUS_API_LIST_ADDAFTER_REQ
			如果设置的是TCaplusValueFlag_NOVALUE, 则操作失败后不返回数据
			如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作失败后返回和请求一致的数据, 暂时没有实现
			如果设置的是TCaplusValueFlag_ALLVALUE, 不合理场景
			如果设置的是TCaplusValueFlag_ALLOLDVALUE, 不返回数据

		8. TCAPLUS_API_LIST_DELETE_REQ
			如果设置的是TCaplusValueFlag_NOVALUE, 则操作失败后不返回数据
			如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作失败后返回和请求一致的数据, 暂时没有实现
			如果设置的是TCaplusValueFlag_ALLVALUE, 不合理场景
			如果设置的是TCaplusValueFlag_ALLOLDVALUE, 如果获取到了tcapsvr端的数据则返回tcpasvr端的数据,如果没有获取到tcapsvr端的数据则返回空

		9. TCAPLUS_API_LIST_REPLACE_REQ
			如果设置的是TCaplusValueFlag_NOVALUE, 则操作失败后不返回数据
			如果设置的是TCaplusValueFlag_SAMEWITHREQUEST, 则操作失败后返回和请求一致的数据, 暂时没有实现
			如果设置的是TCaplusValueFlag_ALLVALUE, 不合理场景
			如果设置的是TCaplusValueFlag_ALLOLDVALUE, 如果获取到了tcapsvr端的数据则返回tcpasvr端的数据,如果没有获取到tcapsvr端的数据则返回空

		@param	[IN] result_flag  请求标志:
									0表示: 只需返回操作执行成功与否
									1表示: 返回与请求字段一致
									2表示: 须返回变更记录的所有字段最新数据
									3表示: 须返回变更记录的所有字段旧数据

									对于batch_get请求，该字段设置为大于0时，某个key查询记录不存在或svr端产生的其它错误时会返回对应的key，
									从而知道是哪个key对应的记录失败了
		@retval 0	 设置成功
		@retval <0	 失败，返回对应的错误码。通常因为未初始化。

	*/

	SetResultFlagForFail(result_flag byte) int

	// 设置perf用于定位各阶段耗时，sendTime 为发送时间戳，单位us
	SetPerfTest(sendTime uint64) int

	/**
	    @brief  设置请求的通用标志位，可以通过"按位或"操作同时设定多个值
	    @param  [IN]  flag. 请求标志位的值
	    @retval 0     设置成功
	    @retval <0    失败，返回对应的错误码。通常因为未初始化。
	    @note   有效的标志位包括：
	    *  TCAPLUS_FLAG_FETCH_ONLY_IF_MODIFIED:
	    *       "数据变更才取回"标志位。在发起读操作之前，用户代码通过 TcaplusServiceRecord::SetVersion()
	    *       带上本地缓存数据的版本号，并将此标志置位，那么存储端检测到当前数据与API本地缓存的数据版本
	    *       一致时，表明该记录未发生过修改，API缓存的数据是最新的，因此在响应中将不会携带实际的数据，
	    *       只是返回 TcapErrCode::COMMON_INFO_DATA_NOT_MODIFIED 的错误码
	    *
	    *       在请求中设置了此标志位之后，收到响应后应首先通过 TcaplusServiceResponse::GetFlags() 来获知
	    *       发送请求时是否设置了TCAPLUS_FLAG_FETCH_ONLY_IF_MODIFIED标志.
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
	    *
	    *       在请求中设置了此标志位之后，收到响应后应首先通过 TcaplusServiceResponse::GetFlags() 来获知
	    *       发送请求时是否设置了 TCAPLUS_FLAG_FETCH_ONLY_IF_EXPIRED 标志.
	    *
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
	SetFlags(flag int32) int

	/**
	  @brief  清理请求的通用标志位，可以通过"按位或"操作同时设定多个值
	  @param  [IN]  flag. 请求标志位的值
	  @retval 0     设置成功
	  @retval <0    失败，返回对应的错误码。通常因为未初始化。
	  @note   有效的标志位列表及详细解释请参考 SetFlags()
	*/
	ClearFlags(flag int32) int

	/**
	  @brief   获取请求的通用标志位
	  @return  返回请求的通用标志位
	  @note   有效的标志位列表及详细解释请参考 SetFlags()
	*/
	GetFlags() int32

	//以下非对外
	GetSeq() int32
	SetSeq(seq int32)
	GetZoneId() uint32
	GetKeyHash() (uint32, error)
	Pack() ([]byte, error)
}

func NewRequest(appId uint64, zoneId uint32, tableName string, cmd int, isPB bool) (TcaplusRequest, error) {
	innerSeq := uint32(0)
	pkg := cs_pool.GetTcaplusCSPkg(uint32(cmd))
	pkg.Head.Magic = uint16(tcaplus_protocol_cs.TCAPLUS_PROTOCOL_MAGIC_CS)
	pkg.Head.Version = uint16(tcaplus_protocol_cs.TCaplusPkgCurrentVersion)
	pkg.Head.Cmd = uint32(cmd)
	pkg.Head.Seq = int32(0)

	pkg.Head.RouterInfo.AppID = int32(appId)
	pkg.Head.RouterInfo.ZoneID = int32(zoneId)

	//string转byte
	pkg.Head.RouterInfo.TableName = common.StringToCByte(tableName)
	pkg.Head.RouterInfo.TableNameLen = uint32(len(pkg.Head.RouterInfo.TableName))

	pkg.Head.KeyInfo.Version = -1
	//pkg.Body.Init(int64(cmd))

	req := &tcapRequest{}
	var err error

	switch cmd {
	case tcaplusCmd.TcaplusApiInsertReq:
		req.commonInterface, err = newInsertRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiGetReq:
		req.commonInterface, err = newGetRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiUpdateReq:
		req.commonInterface, err = newUpdateRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiReplaceReq:
		req.commonInterface, err = newReplaceRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiDeleteReq:
		req.commonInterface, err = newdeleteRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiGetByPartkeyReq:
		req.commonInterface, err = newGetByPartKeyRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiDeleteByPartkeyReq:
		req.commonInterface, err = newDeleteByPartKeyRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiIncreaseReq:
		req.commonInterface, err = newIncreaseRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListGetAllReq:
		req.commonInterface, err = newListGetAllRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiBatchGetReq:
		req.commonInterface, err = newBatchGetRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	//  not support yet
	//case tcaplusCmd.TcaplusApiUpdateByPartkeyReq:
	//	req.commonInterface, err = newUpdateByPartKeyRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListAddAfterReq:
		req.commonInterface, err = newListAddAfterRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListGetReq:
		req.commonInterface, err = newListGetRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListDeleteReq:
		req.commonInterface, err = newListDeleteRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListDeleteAllReq:
		req.commonInterface, err = newListDeleteAllRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListReplaceReq:
		req.commonInterface, err = newListReplaceRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListDeleteBatchReq:
		req.commonInterface, err = newListDeleteBatchRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiSqlReq:
		req.commonInterface, err = newIndexQueryRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiMetadataGetReq:
		req.commonInterface, err = newGetMetaRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiPBFieldGetReq:
		req.commonInterface, err = newPBFieldGetRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiPBFieldUpdateReq:
		req.commonInterface, err = newPBFieldUpdateRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiPBFieldIncreaseReq:
		req.commonInterface, err = newPBFieldIncreaseRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiGetShardListReq:
		req.commonInterface, err = newGetShardListRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiTableTraverseReq:
		req.commonInterface, err = newTraverseRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiGetTableRecordCountReq:
		req.commonInterface, err = newCountRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiBatchInsertReq:
		req.commonInterface, err = newBatchInsertRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiBatchDeleteReq:
		req.commonInterface, err = newBatchDeleteRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiBatchReplaceReq:
		req.commonInterface, err = newBatchReplaceRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiBatchUpdateReq:
		req.commonInterface, err = newBatchUpdateRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiGetTtlReq:
		req.commonInterface, err = newGetTtlRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiSetTtlReq:
		req.commonInterface, err = newSetTtlRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListGetBatchReq:
		req.commonInterface, err = newListGetBatchRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListAddAfterBatchReq:
		req.commonInterface, err = newListAddAfterBatchRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	case tcaplusCmd.TcaplusApiListReplaceBatchReq:
		req.commonInterface, err = newListReplaceBatchRequest(appId, zoneId, tableName, cmd, innerSeq, pkg, isPB)
	default:
		logger.ERR("invalid cmd %d", cmd)
		return nil, &terror.ErrorCode{Code: terror.InvalidCmd}
	}
	return req, err
}

func setUserBuffer(pkg *tcaplus_protocol_cs.TCaplusPkg, userBuffer []byte) error {
	if pkg == nil {
		logger.ERR("Request can not second use")
		return &terror.ErrorCode{Code: terror.RequestHasHasNoPkg, Message: "Request can not second use"}
	}
	bufLen := len(userBuffer)
	if bufLen <= 0 {
		return nil
	}

	if bufLen > int(tcaplus_protocol_cs.TCAPLUS_MAX_USERBUFF_LEN) {
		logger.ERR("userBuffer len %d > %d", bufLen, tcaplus_protocol_cs.TCAPLUS_MAX_USERBUFF_LEN)
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_DATA_SIZE}
	}

	pkg.Head.UserBuff = userBuffer
	pkg.Head.UserBuffLen = uint32(bufLen)
	return nil
}

func keyHashCode(keySet *tcaplus_protocol_cs.TCaplusKeySet) (uint32, error) {
	if keySet.FieldNum <= 0 {
		return 0, &terror.ErrorCode{Code: terror.RequestHasNoKeyField}
	}

	field := keySet.Fields[0:keySet.FieldNum]
	sort.Slice(field, func(i, j int) bool {
		if field[i].FieldName < field[j].FieldName {
			return true
		}
		return false
	})

	var buf []byte
	for _, v := range field {
		buf = append(buf, v.FieldBuff[0:v.FieldLen]...)
	}
	if len(buf) <= 0 {
		return 0, nil
	}
	return crc32.ChecksumIEEE(buf), nil
}

var allowdFlagCmdMap = [][]uint32{
	/* bit1 (0x00000001): TCAPLUS_FLAG_FETCH_ONLY_IF_MODIFIED */
	{tcaplusCmd.TcaplusApiGetReq,
		tcaplusCmd.TcaplusApiListGetReq,
		tcaplusCmd.TcaplusApiListGetAllReq},
	/* bit2 (0x00000002): TCAPLUS_FLAG_FETCH_ONLY_IF_EXPIRED */
	{tcaplusCmd.TcaplusApiBatchGetReq},
	/* bit3 (0x00000004): TCAPLUS_FLAG_ONLY_READ_FROM_SLAVE */
	{tcaplusCmd.TcaplusApiGetReq,
		//TCAPLUS_API_GET_TABLE_RECORD_COUNT_REQ,
		tcaplusCmd.TcaplusApiListGetReq,
		tcaplusCmd.TcaplusApiListGetAllReq,
		tcaplusCmd.TcaplusApiGetByPartkeyReq,
		//TCAPLUS_API_BATCH_GET_BY_PARTKEY_REQ,
		//TCAPLUS_API_METADATA_GET_REQ,
		tcaplusCmd.TcaplusApiTableTraverseReq,
		//TCAPLUS_API_LIST_TABLE_TRAVERSE_REQ,
		tcaplusCmd.TcaplusApiBatchGetReq},
	/* bit4 (0x00000008): TCAPLUS_FLAG_LIST_RESERVE_INDEX_HAVING_NO_ELEMENTS */
	{tcaplusCmd.TcaplusApiListDeleteReq,
		tcaplusCmd.TcaplusApiListDeleteAllReq,
		tcaplusCmd.TcaplusApiListDeleteBatchReq},
	/* bit5 (0x00000010): TcaplusFlagInsertRecordIfNotExist int32 = 16 PB的FieldUpdate使用，数据不存在则插入*/
	{tcaplusCmd.TcaplusApiPBFieldUpdateReq},
}

func manipulateFlags(pkg *tcaplus_protocol_cs.TCaplusPkg, flags int32, clear bool) int {
	if pkg == nil {
		return int(terror.RequestHasHasNoPkg)
	}

	// 针对每个flag检查合法
	for i := 0; i < len(allowdFlagCmdMap); i++ {
		if ((1 << uint32(i)) & flags) != 0 {
			continue
		}

		k := 0
		for ; k < len(allowdFlagCmdMap[i]); k++ {
			if allowdFlagCmdMap[i][k] == pkg.Head.Cmd {
				break
			}
		}

		if k >= len(allowdFlagCmdMap[i]) {
			return int(terror.OperationNotSupport)
		}
	}

	if clear {
		pkg.Head.Flags &= ^flags
	} else {
		pkg.Head.Flags |= flags
	}

	return 0
}

func setFlags(pkg *tcaplus_protocol_cs.TCaplusPkg, flags int32) int {
	if pkg == nil {
		logger.ERR("Request can not second use")
		return int(terror.RequestHasHasNoPkg)
	}
	return manipulateFlags(pkg, flags, false)
}

func clearFlags(pkg *tcaplus_protocol_cs.TCaplusPkg, flags int32) int {
	if pkg == nil {
		logger.ERR("Request can not second use")
		return int(terror.RequestHasHasNoPkg)
	}
	return manipulateFlags(pkg, flags, true)
}

type tcapRequest struct {
	commonInterface
}

func (req *tcapRequest) SetListShiftFlag(shiftFlag byte) int32 {
	switch req.commonInterface.(type) {
	case *listAddAfterRequest:
		return req.commonInterface.(*listAddAfterRequest).SetListShiftFlag(shiftFlag)
	case *listAddAfterBatchRequest:
		return req.commonInterface.(*listAddAfterBatchRequest).SetListShiftFlag(shiftFlag)
	default:
		return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
	}
}

func (req *tcapRequest) SetAddableIncreaseFlag(increase_flag byte) int32 {
	switch req.commonInterface.(type) {
	case *increaseRequest:
		return req.commonInterface.(*increaseRequest).SetAddableIncreaseFlag(increase_flag)
	case *pbFieldIncreaseRequest:
		return req.commonInterface.(*pbFieldIncreaseRequest).SetEnableIncreaseNotExist(increase_flag)
	default:
		return int32(terror.GEN_ERR_SUC)
	}
}

func (req *tcapRequest) AddElementIndex(idx int32) int32 {
	switch req.commonInterface.(type) {
	case *listDeleteBatchRequest:
		return req.commonInterface.(*listDeleteBatchRequest).AddElementIndex(idx)
	case *listGetBatchRequest:
		return req.commonInterface.(*listGetBatchRequest).AddElementIndex(idx)
	default:
		return int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH)
	}
}

func (req *tcapRequest) SetSql(query string) int {
	switch req.commonInterface.(type) {
	case *indexQueryRequest:
		return req.commonInterface.(*indexQueryRequest).SetSql(query)
	default:
		return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
	}
}

func (req *tcapRequest) GetTcaplusPackagePtr() *tcaplus_protocol_cs.TCaplusPkg {
	switch req.commonInterface.(type) {
	case *getShardListRequest:
		return req.commonInterface.(*getShardListRequest).GetTcaplusPackagePtr()
	case *traverseRequest:
		return req.commonInterface.(*traverseRequest).GetTcaplusPackagePtr()
	default:
		return nil
	}
}
