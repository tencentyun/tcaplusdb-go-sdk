package terror

import "strconv"

var (
	//go api 定义错误码
	ParameterInvalid       = -0x0000001e
	DirSignUpFailed        = -0x0000011e
	ClientInitTimeOut      = -0x0000021e
	ProxySignUpFailed      = -0x0000031e
	ZoneIdNotExist         = -0x0000041e
	TableNotExist          = -0x0000051e
	InvalidCmd             = -0x0000061e
	InvalidPolicy          = -0x0000071e
	RecordToMax            = -0x0000081e
	KeyNameLenOverMax      = -0x0000091e
	KeyLenOverMax          = -0x00000a1e
	KeyNumOverMax          = -0x00000b1e
	ValueNameLenOverMax    = -0x00000c1e
	ValueLenOverMax        = -0x00000d1e
	ValueNumOverMax        = -0x00000e1e
	ValuePackOverMax       = -0x00000f1e
	RecordNumOverMax       = -0x0000101e
	ProxyNotAvailable      = -0x0000111e
	RequestHasNoRecord     = -0x0000121e
	RequestHasNoKeyField   = -0x0000131e
	RecordKeyTypeInvalid   = -0x0000141e
	RecordValueTypeInvalid = -0x0000151e
	OperationNotSupport    = -0x0000161e
	ClientNotInit          = -0x0000171e
	RecordUnpackFailed     = -0x0000181e
	RecordKeyNotExist      = -0x0000191e
	RecordValueNotExist    = -0x00001a1e
	ClientNotDial		   = -0x00001b1e
	RespNotMatchReq		   = -0x00001c1e
	MetadataNotProtobuf	   = -0x00001d1e
	SqlQueryFormatError	   = -0x00001e1e

	/*****************************************************************************************
	 **********************************C版本错误码*********************************************
	 ****************************************************************************************/
	//GENERAL BUSINESS (module id 0x00) Error Code defined below
	GEN_ERR_SUC                             = 0x00000000
	GEN_ERR_ERR                             = -0x00000100 /*-256*/
	GEN_ERR_ECMGR_INVALID_MODULE_ID         = -0x00000200 /*-512*/
	GEN_ERR_ECMGR_INVALID_ERROR_CODE        = -0x00000300 /*-768*/
	GEN_ERR_ECMGR_NULL_ERROR_STRING         = -0x00000400 /*-1024*/
	GEN_ERR_ECMGR_DUPLICATED_ERROR_CODE     = -0x00000500 /*-1280*/
	GEN_ERR_TXLOG_NULL_POINTER_FROM_TSD     = -0x00000600 /*-1536*/
	GEN_ERR_TABLE_READONLY                  = -0x00000700 /*-1792*/
	GEN_ERR_TABLE_READ_DELETE               = -0x00000800 /*-2048*/
	GEN_ERR_ACCESS_DENIED                   = -0x00000900 /*-2304*/
	GEN_ERR_INVALID_ARGUMENTS               = -0x00000A00 /*-2560*/
	GEN_ERR_UNSUPPORT_OPERATION             = -0x00000B00 /*-2816*/
	GEN_ERR_NOT_ENOUGH_MEMORY               = -0x00000C00 /*-3072*/
	GEN_ERR_NOT_SATISFY_INSERT_FOR_SORTLIST = -0x00000D00 /*-3328*/
	GEN_ERR_BASE64_ENCODE_FAILED            = -0x00000E00 /*-3584*/
	GEN_ERR_BASE64_DECODE_FAILED            = -0x00000F00 /*-3840*/

	//LINELOC BUSINESS (module id 0x02) Error Code defined below
	LOC_ERR__0x00000102 = -0x00000102 /*-258*/
	LOC_ERR__0x00000202 = -0x00000202 /*-514*/
	LOC_ERR__0x00000302 = -0x00000302 /*-770*/
	LOC_ERR__0x00000402 = -0x00000402 /*-1026*/
	LOC_ERR__0x00000502 = -0x00000502 /*-1282*/
	LOC_ERR__0x00000602 = -0x00000602 /*-1538*/
	LOC_ERR__0x00000702 = -0x00000702 /*-1794*/
	LOC_ERR__0x00000802 = -0x00000802 /*-2050*/
	LOC_ERR__0x00000902 = -0x00000902 /*-2306*/
	LOC_ERR__0x00000A02 = -0x00000A02 /*-2562*/
	LOC_ERR__0x00000B02 = -0x00000B02 /*-2818*/
	LOC_ERR__0x00000C02 = -0x00000C02 /*-3074*/
	LOC_ERR__0x00000D02 = -0x00000D02 /*-3330*/
	LOC_ERR__0x00000E02 = -0x00000E02 /*-3586*/
	LOC_ERR__0x00000F02 = -0x00000F02 /*-3842*/
	LOC_ERR__0x00001002 = -0x00001002 /*-4098*/
	LOC_ERR__0x00001102 = -0x00001102 /*-4354*/
	LOC_ERR__0x00001202 = -0x00001202 /*-4610*/
	LOC_ERR__0x00001302 = -0x00001302 /*-4866*/
	LOC_ERR__0x00001402 = -0x00001402 /*-5122*/
	LOC_ERR__0x00001502 = -0x00001502 /*-5378*/
	LOC_ERR__0x00001602 = -0x00001602 /*-5634*/
	LOC_ERR__0x00001702 = -0x00001702 /*-5890*/
	LOC_ERR__0x00001802 = -0x00001802 /*-6146*/
	LOC_ERR__0x00001902 = -0x00001902 /*-6402*/
	LOC_ERR__0x00001A02 = -0x00001A02 /*-6658*/
	LOC_ERR__0x00001B02 = -0x00001B02 /*-6914*/
	LOC_ERR__0x00001C02 = -0x00001C02 /*-7170*/
	LOC_ERR__0x00001D02 = -0x00001D02 /*-7426*/
	LOC_ERR__0x00001E02 = -0x00001E02 /*-7682*/
	LOC_ERR__0x00001F02 = -0x00001F02 /*-7938*/
	LOC_ERR__0x00002002 = -0x00002002 /*-8194*/
	LOC_ERR__0x00002802 = -0x00002802 /*-10242*/
	LOC_ERR__0x00003002 = -0x00003002 /*-12290*/
	LOC_ERR__0x00003802 = -0x00003802 /*-14338*/
	LOC_ERR__0x00004002 = -0x00004002 /*-16386*/
	LOC_ERR__0x00004802 = -0x00004802 /*-18434*/
	LOC_ERR__0x00005002 = -0x00005002 /*-20482*/
	LOC_ERR__0x00005802 = -0x00005802 /*-22530*/
	LOC_ERR__0x00006002 = -0x00006002 /*-24578*/
	LOC_ERR__0x00006802 = -0x00006802 /*-26626*/
	LOC_ERR__0x00007002 = -0x00007002 /*-28674*/
	LOC_ERR__0x00007802 = -0x00007802 /*-30722*/
	LOC_ERR__0x00008002 = -0x00008002 /*-32770*/
	LOC_ERR__0x00008802 = -0x00008802 /*-34818*/
	LOC_ERR__0x00009002 = -0x00009002 /*-36866*/
	LOC_ERR__0x00009802 = -0x00009802 /*-38914*/
	LOC_ERR__0x0000A002 = -0x0000A002 /*-40962*/
	LOC_ERR__0x0000A802 = -0x0000A802 /*-43010*/
	LOC_ERR__0x0000B002 = -0x0000B002 /*-45058*/
	LOC_ERR__0x0000B802 = -0x0000B802 /*-47106*/
	LOC_ERR__0x0000C002 = -0x0000C002 /*-49154*/
	LOC_ERR__0x0000C802 = -0x0000C802 /*-51202*/
	LOC_ERR__0x0000FF02 = -0x0000FF02 /*-65282*/

	//TXHDB SYSTEM (module id 0x05) Error Code defined below
	TXHDB_ERR_RECORD_NOT_EXIST                              = 0x00000105 /*261*/
	TXHDB_ERR_ITERATION_NO_MORE_RECORDS                     = 0x00000205 /*517*/
	TXHDB_ERR_MUTEX_TRYLOCK_BUSY                            = 0x00000305 /*773*/
	TXHDB_ERR_MUTEX_TIMEDLOCK_TIMEOUT                       = 0x00000405 /*1029*/
	TXHDB_ERR_RWLOCK_TRYWRLOCK_BUSY                         = 0x00000505 /*1285*/
	TXHDB_ERR_RWLOCK_TRYRDLOCK_BUSY                         = 0x00000605 /*1541*/
	TXHDB_ERR_SPIN_TRYLOCK_BUSY                             = 0x00000705 /*1797*/
	TXHDB_ERR_ITERATION_EXCEED_MAX_ALLOWED_TIME_OF_ONE_ITER = 0x00000805 /*2053*/

	TXHDB_ERR_INVALID_ARGUMENTS                                                    = -0x00000105 /*-261*/
	TXHDB_ERR_INVALID_MEMBER_VARIABLE_VALUE                                        = -0x00000205 /*-517*/
	TXHDB_ERR_ALREADY_OPEN                                                         = -0x00000305 /*-773*/
	TXHDB_ERR_MUTEX_LOCK_FAIL                                                      = -0x00000405 /*-1029*/
	TXHDB_ERR_MUTEX_TRYLOCK_FAIL                                                   = -0x00000505 /*-1285*/
	TXHDB_ERR_MUTEX_TIMEDLOCK_FAIL                                                 = -0x00000605 /*-1541*/
	TXHDB_ERR_MUTEX_UNLOCK_FAIL                                                    = -0x00000705 /*-1797*/
	TXHDB_ERR_RWLOCK_WRLOCK_FAIL                                                   = -0x00000805 /*-2053*/
	TXHDB_ERR_RWLOCK_TRYWRLOCK_FAIL                                                = -0x00000905 /*-2309*/
	TXHDB_ERR_RWLOCK_RDLOCK_FAIL                                                   = -0x00000a05 /*-2565*/
	TXHDB_ERR_RWLOCK_TRYRDLOCK_FAIL                                                = -0x00000b05 /*-2821*/
	TXHDB_ERR_RWLOCK_UNLOCK_FAIL                                                   = -0x00000c05 /*-3077*/
	TXHDB_ERR_SPIN_LOCK_FAIL                                                       = -0x00000d05 /*-3333*/
	TXHDB_ERR_SPIN_UNLOCK_FAIL                                                     = -0x00000e05 /*-3589*/
	TXHDB_ERR_FILE_EXISTS_BUT_STATUS_ERROR                                         = -0x00000f05 /*-3845*/
	TXHDB_ERR_FILE_OPEN_FAIL                                                       = -0x00001005 /*-4101*/
	TXHDB_ERR_FILE_READ_SIZE_INVALID                                               = -0x00001105 /*-4357*/
	TXHDB_ERR_FILE_INVALID_FILE_PATH                                               = -0x00001205 /*-4613*/
	TXHDB_ERR_FILE_LOCK_FILE_FAIL                                                  = -0x00001305 /*-4869*/
	TXHDB_ERR_FILE_NOT_A_REGULAR_FILE                                              = -0x00001405 /*-5125*/
	TXHDB_ERR_FILE_MMAP_FAIL                                                       = -0x00001505 /*-5381*/
	TXHDB_ERR_FILE_MUNMAP_FAIL                                                     = -0x00001605 /*-5637*/
	TXHDB_ERR_FILE_CLOSE_FAIL                                                      = -0x00001705 /*-5893*/
	TXHDB_ERR_FILE_SPACE_NOT_ENOUGH_IN_HEAD                                        = -0x00001805 /*-6149*/
	TXHDB_ERR_FILE_FTRUNCATE_FAIL                                                  = -0x00001905 /*-6405*/
	TXHDB_ERR_FILE_INCONSISTANT_FILE_SIZE                                          = -0x00001a05 /*-6661*/
	TXHDB_ERR_FILE_MSIZ_LESSER_THAN_TXHDB_WHOLE_REC_OFFSET                         = -0x00001b05 /*-6917*/
	TXHDB_ERR_FILE_MSIZ_CHANGE_NOT_PERMIT                                          = -0x00001c05 /*-7173*/
	TXHDB_ERR_FILE_FSTAT_FAIL                                                      = -0x00001d05 /*-7429*/
	TXHDB_ERR_FILE_MSYNC_FAIL                                                      = -0x00001e05 /*-7685*/
	TXHDB_ERR_FILE_FSYNC_FAIL                                                      = -0x00001f05 /*-7941*/
	TXHDB_ERR_FILE_FCNTL_LOCK_FILE_FAIL                                            = -0x00002005 /*-8197*/
	TXHDB_ERR_FILE_FCNTL_UNLOCK_FILE_FAIL                                          = -0x00002105 /*-8453*/
	TXHDB_ERR_FILE_PREAD_FAIL_WITH_SPECIFIED_ERRNO                                 = -0x00002205 /*-8709*/
	TXHDB_ERR_FILE_PREAD_FAIL_WITH_UNSPECIFIED_ERRNO                               = -0x00002305 /*-8965*/
	TXHDB_ERR_FILE_PWRITE_FAIL_WITH_SPECIFIED_ERRNO                                = -0x00002405 /*-9221*/
	TXHDB_ERR_FILE_PWRITE_FAIL_WITH_UNSPECIFIED_ERRNO                              = -0x00002505 /*-9477*/
	TXHDB_ERR_FILE_READ_EXCEED_FILE_BOUNDARY                                       = -0x00002605 /*-9733*/
	TXHDB_ERR_FILE_READ_FAIL_DURING_COPY                                           = -0x00002705 /*-9989*/
	TXHDB_ERR_FILE_WRITE_FAIL_DURING_COPY                                          = -0x00002805 /*-10245*/
	TXHDB_ERR_FILE_INVALID_FREE_BLOCK_POOL_METADATA                                = -0x00002905 /*-10501*/
	TXHDB_ERR_FILE_INVALID_MAGIC                                                   = -0x00002a05 /*-10757*/
	TXHDB_ERR_FILE_INVALID_LIBRARY_VERSION                                         = -0x00002b05 /*-11013*/
	TXHDB_ERR_FILE_INVALID_LIBRARY_REVISION                                        = -0x00002c05 /*-11269*/
	TXHDB_ERR_FILE_INVALID_FORMAT_VERSION                                          = -0x00002d05 /*-11525*/
	TXHDB_ERR_FILE_INVALID_EXTDATA_FORMAT_VERSION                                  = -0x00002e05 /*-11781*/
	TXHDB_ERR_FILE_INVALID_DBTYPE                                                  = -0x00002f05 /*-12037*/
	TXHDB_ERR_FILE_HEAD_CRC_UNMATCH                                                = -0x00003005 /*-12293*/
	TXHDB_ERR_FILE_INVALID_METADATA                                                = -0x00003105 /*-12549*/
	TXHDB_ERR_FILE_INVALID_HEADLEN                                                 = -0x00003205 /*-12805*/
	TXHDB_ERR_FILE_DESERIAL_HEAD_SPACE_NOT_ENOUGH                                  = -0x00003305 /*-13061*/
	TXHDB_ERR_FILE_SERIAL_HEAD_SPACE_NOT_ENOUGH                                    = -0x00003405 /*-13317*/
	TXHDB_ERR_FILE_DESERIAL_STAT_SPACE_NOT_ENOUGH                                  = -0x00003505 /*-13573*/
	TXHDB_ERR_FILE_SERIAL_STAT_SPACE_NOT_ENOUGH                                    = -0x00003605 /*-13829*/
	TXHDB_ERR_FILE_SERIAL_FREE_BLOCK_LIST_INFO_WRONG_BUFFLEN                       = -0x00003705 /*-14085*/
	TXHDB_ERR_FILE_IN_EXCEPTIONAL_STATUS                                           = -0x00003805 /*-14341*/
	TXHDB_ERR_DB_NOT_OPENED                                                        = -0x00003905 /*-14597*/
	TXHDB_ERR_DB_WRITE_NOT_PERMIT                                                  = -0x00003a05 /*-14853*/
	TXHDB_ERR_INVALID_OFFSET_FROM_BUCKET                                           = -0x00003b05 /*-15109*/
	TXHDB_ERR_READ_EXTDATA_EXCEED_BUFF_LENGTH                                      = -0x00003c05 /*-15365*/
	TXHDB_ERR_WRITE_EXTDATA_EXCEED_BUFF_LENGTH                                     = -0x00003d05 /*-15621*/
	TXHDB_ERR_FREE_BLOCK_IS_READ_WHEN_GETTING_RECORD                               = -0x00003e05 /*-15877*/
	TXHDB_ERR_INVALID_KEY_DATABLOCK_NUM                                            = -0x00003f05 /*-16133*/
	TXHDB_ERR_INVALID_VALUE_DATABLOCK_NUM                                          = -0x00004005 /*-16389*/
	TXHDB_ERR_GET_RECORD_EXCEED_BUFF_LENGTH                                        = -0x00004105 /*-16645*/
	TXHDB_ERR_COMPRESSION_FAIL                                                     = -0x00004205 /*-16901*/
	TXHDB_ERR_DECOMPRESSION_FAIL                                                   = -0x00004305 /*-17157*/
	TXHDB_ERR_INVALID_OFFSETINEXTDATA_AND_SIZE_WHEN_UPDATING_EXTDATA               = -0x00004405 /*-17413*/
	TXHDB_ERR_UNEXPECTED_FREEBLOCK                                                 = -0x00004505 /*-17669*/
	TXHDB_ERR_VALUE_APOW_LESSER_THAN_KEY_APOW                                      = -0x00004605 /*-17925*/
	TXHDB_ERR_DUPLICATED_FILE_PATH                                                 = -0x00004705 /*-18181*/
	TXHDB_ERR_INVALID_KEY_HEAD_SIZE_IN_TXHDB_META                                  = -0x00004805 /*-18437*/
	TXHDB_ERR_INVALID_FILE_SIZE                                                    = -0x00004905 /*-18693*/
	TXHDB_ERR_INVALID_FREE_BLOCK_SIZE                                              = -0x00004a05 /*-18949*/
	TXHDB_ERR_MMAP_MEMSIZE_CHANGE_NOT_PERMITTED                                    = -0x00004b05 /*-19205*/
	TXHDB_ERR_NEW_FILE_OBJ_FAIL                                                    = -0x00004c05 /*-19461*/
	TXHDB_ERR_RECORD_KEY_OFFSET_LESSER_THAN_TXHDB_WHOLE_REC_OFFSET                 = -0x00004d05 /*-19717*/
	TXHDB_ERR_RECORD_VALUE_OFFSET_LESSER_THAN_TXHDB_WHOLE_REC_OFFSET               = -0x00004e05 /*-19973*/
	TXHDB_ERR_RECORD_OFFSET_LESSER_THAN_TXHDB_WHOLE_REC_OFFSET                     = -0x00004f05 /*-20229*/
	TXHDB_ERR_KEY_BUFFSIZE_LESSER_THAN_KEY_HEADSIZE                                = -0x00005005 /*-20485*/
	TXHDB_ERR_VALUE_BUFFSIZE_LESSER_THAN_VALUE_HEADSIZE                            = -0x00005105 /*-20741*/
	TXHDB_ERR_RECORD_SIZE_LESSER_THAN_KEY_HEADSIZE                                 = -0x00005205 /*-20997*/
	TXHDB_ERR_INVALID_BLOCK_MAGIC                                                  = -0x00005305 /*-21253*/
	TXHDB_ERR_INVALID_FREE_BLOCK_MAGIC                                             = -0x00005405 /*-21509*/
	TXHDB_ERR_INVALID_KEYMAGIC                                                     = -0x00005505 /*-21765*/
	TXHDB_ERR_INVALID_KEYSPLMAGIC                                                  = -0x00005605 /*-22021*/
	TXHDB_ERR_INVALID_VALMAGIC                                                     = -0x00005705 /*-22277*/
	TXHDB_ERR_INVALID_VALSPLMAGIC                                                  = -0x00005805 /*-22533*/
	TXHDB_ERR_UNSUPPORTED_KEY_FORMAT_VERSION                                       = -0x00005905 /*-22789*/
	TXHDB_ERR_UNSUPPORTED_KEY_SPLBLOCK_FORMAT_VERSION                              = -0x00005a05 /*-23045*/
	TXHDB_ERR_UNSUPPORTED_VALUE_FORMAT_VERSION                                     = -0x00005b05 /*-23301*/
	TXHDB_ERR_UNSUPPORTED_VALUE_SPLBLOCK_FORMAT_VERSION                            = -0x00005c05 /*-23557*/
	TXHDB_ERR_UNSUPPORTED_FREE_BLOCK_FORMAT_VERSION                                = -0x00005d05 /*-23813*/
	TXHDB_ERR_KEY_HEAD_CRC_UNMATCH                                                 = -0x00005e05 /*-24069*/
	TXHDB_ERR_KEY_SPLBLOCK_HEAD_CRC_UNMATCH                                        = -0x00005f05 /*-24325*/
	TXHDB_ERR_VALUE_HEAD_CRC_UNMATCH                                               = -0x00006005 /*-24581*/
	TXHDB_ERR_VALUE_SPLBLOCK_HEAD_CRC_UNMATCH                                      = -0x00006105 /*-24837*/
	TXHDB_ERR_FREE_BLOCK_HEAD_CRC_UNMATCH                                          = -0x00006205 /*-25093*/
	TXHDB_ERR_FREE_BLOCK_LIST_INFO_CRC_UNMATCH                                     = -0x00006305 /*-25349*/
	TXHDB_ERR_GET_KEY_READ_BUFFER_FAIL                                             = -0x00006405 /*-25605*/
	TXHDB_ERR_GET_VALUE_READ_BUFFER_FAIL                                           = -0x00006505 /*-25861*/
	TXHDB_ERR_GET_LRU_VALUE_BUFFER_FAIL                                            = -0x00006605 /*-26117*/
	TXHDB_ERR_GET_EXTDATA_READ_BUFFER_FAIL                                         = -0x00006705 /*-26373*/
	TXHDB_ERR_KEY_BLOCK_BODYSIZE_GREATER_THAN_KEY_BODYSIZE                         = -0x00006805 /*-26629*/
	TXHDB_ERR_VALUE_BLOCK_BODYSIZE_GREATER_THAN_VALUE_BODYSIZE                     = -0x00006905 /*-26885*/
	TXHDB_ERR_NULL_RECORD_POINTER                                                  = -0x00006a05 /*-27141*/
	TXHDB_ERR_NULL_RECORD_WRITE_BUFF                                               = -0x00006b05 /*-27397*/
	TXHDB_ERR_SERIALIZE_RECORD_KEY_HEAD                                            = -0x00006c05 /*-27653*/
	TXHDB_ERR_INVALID_IDX_IN_STAT_NUMS_ARRAY                                       = -0x00006d05 /*-27909*/
	TXHDB_ERR_INVALID_ELEMNUM_OF_STAT_KEYNUMS                                      = -0x00006e05 /*-28165*/
	TXHDB_ERR_INVALID_ELEMNUM_OF_STAT_VALNUMS                                      = -0x00006f05 /*-28421*/
	TXHDB_ERR_PRINT_SPACE_NOT_ENOUGH                                               = -0x00007005 /*-28677*/
	TXHDB_ERR_LRU_SHIFTIN_NOT_ENOUGH_MEMORY                                        = -0x00007105 /*-28933*/
	TXHDB_ERR_LRU_SHIFTIN_NO_MORE_LRU_NODE                                         = -0x00007205 /*-29189*/
	TXHDB_ERR_LRU_ADJUST_NO_MORE_LRU_NODE                                          = -0x00007305 /*-29445*/
	TXHDB_ERR_LRU_SHIFTOUT_RECORD_ALREADY_OUTSIDE_OF_MEMORY                        = -0x00007405 /*-29701*/
	TXHDB_ERR_FILE_EXTDATA_LENGTH_CRC_UNMATCH                                      = -0x00007505 /*-29957*/
	TXHDB_ERR_FILE_EXTDATA_INVALID_LENGTH                                          = -0x00007605 /*-30213*/
	TXHDB_ERR_INVALID_VALUE_HEAD_SIZE_IN_TXHDB_META                                = -0x00007705 /*-30469*/
	TXHDB_ERR_INVALID_SPLITDATABLOCK_HEAD_SIZE_IN_TXHDB_META                       = -0x00007805 /*-30725*/
	TXHDB_ERR_KEY_BUCKETIDX_UNMATCH                                                = -0x00007905 /*-30981*/
	TXHDB_ERR_FILE_WRITE_SIZE_INVALID                                              = -0x00007a05 /*-31237*/
	TXHDB_ERR_MODIFY_STAT_UNSUPPORTED_OPERATION_TYPE                               = -0x00007b05 /*-31493*/
	TXHDB_ERR_INVALID_EXTDATAMAGIC                                                 = -0x00007c05 /*-31749*/
	TXHDB_ERR_INVALID_INTERNAL_LIST_TAIL_DURING_POP_LRU_NODELIST                   = -0x00007d05 /*-32005*/
	TXHDB_ERR_GET_LRUNODE_FAIL                                                     = -0x00007e05 /*-32261*/
	TXHDB_ERR_LRUNODE_INVALID_FLAG                                                 = -0x00007f05 /*-32517*/
	TXHDB_ERR_INVALID_FREE_BLOCK_NUM_TOO_MANY_FREE_BLOCKS                          = -0x00008005 /*-32773*/
	TXHDB_ERR_INVALID_ELEMNUM_OF_STAT_NOPADDING_SIZE_KEYNUMS                       = -0x00008105 /*-33029*/
	TXHDB_ERR_INVALID_ELEMNUM_OF_STAT_NOPADDING_SIZE_VALNUMS                       = -0x00008205 /*-33285*/
	TXHDB_ERR_ADD_LSIZE_EXCEEDS_MAX_TSD_VALUE_BUFF_SIZE                            = -0x00008305 /*-33541*/
	TXHDB_ERR_INTERNAL_CONSTANTS_ILLEGAL                                           = -0x00008405 /*-33797*/
	TXHDB_ERR_TOO_BIG_KEY_BIZ_SIZE                                                 = -0x00008505 /*-34053*/
	TXHDB_ERR_TOO_BIG_VALUE_BIZ_SIZE                                               = -0x00008605 /*-34309*/
	TXHDB_ERR_INDEX_NO_EXIST                                                       = -0x00008705 /*-34565*/
	TXHDB_ERR_INVALID_FREE_BLOCK_BASESIZE                                          = -0x00008805 /*-34821*/
	TXHDB_ERR_CANNOT_CREATE_MMAPSHM_BECAUSE_SHM_ALREADY_EXISTED                    = -0x00008905 /*-35077*/
	TXHDB_ERR_INVALID_GENSHM_KEY                                                   = -0x00008a05 /*-35333*/
	TXHDB_ERR_GENSHM_GET_FAIL                                                      = -0x00008b05 /*-35589*/
	TXHDB_ERR_GENSHM_CREATE_FAIL                                                   = -0x00008c05 /*-35845*/
	TXHDB_ERR_GENSHM_STAT_FAIL                                                     = -0x00008d05 /*-36101*/
	TXHDB_ERR_GENSHM_DOES_NOT_EXIST                                                = -0x00008e05 /*-36357*/
	TXHDB_ERR_GENSHM_ATTACH_FAIL_BECAUSE_IT_IS_ALREADY_ATTACHED_BY_OTHER_PROCESSES = -0x00008f05 /*-36613*/
	TXHDB_ERR_GENSHM_ATTACH_FAIL                                                   = -0x00009005 /*-36869*/
	TXHDB_ERR_FILE_INCONSISTANT_MSIZE                                              = -0x00009105 /*-37125*/
	TXHDB_ERR_INVALID_TCAP_GENSHM_MAGIC                                            = -0x00009205 /*-37381*/
	TXHDB_ERR_GENSHM_FIXED_HEAD_BUFFLEN_UNMATCH                                    = -0x00009305 /*-37637*/
	TXHDB_ERR_GENSHM_INVALID_HEADLEN                                               = -0x00009405 /*-37893*/
	TXHDB_ERR_GENSHM_HEAD_CRC_UNMATCH                                              = -0x00009505 /*-38149*/
	TXHDB_ERR_GENSHM_HEAD_INVALID_VERSION                                          = -0x00009605 /*-38405*/
	TXHDB_ERR_GENSHM_INVALID_FILETYPE                                              = -0x00009705 /*-38661*/
	TXHDB_ERR_GET_IPV4ADDR_FAIL                                                    = -0x00009805 /*-39429*/
	TXHDB_ERR_NO_VALID_IPV4ADDR_EXISTS                                             = -0x00009905 /*-39173*/
	TXHDB_ERR_TRANSFER_IPV4ADDR_FAIL                                               = -0x00009a05 /*-39429*/
	TXHDB_ERR_FILE_EXCEEDS_LSIZE_LIMIT                                             = -0x00009b05 /*-39685*/
	TXHDB_ERR_GENSHM_DETACH_FAIL                                                   = -0x00009c05 /*-39941*/
	TXHDB_ERR_TXHDB_HEAD_PARAMETERS_ERROR                                          = -0x00009d05 /*-40197*/
	TXHDB_ERR_TXHDB_HEAD_OLD_VERSION                                               = -0x00009e05 /*-40453*/
	TXHDB_ERR_TXHDB_SHM_COREINFO_UNMATCH                                           = -0x00009f05 /*-40709*/
	TXHDB_ERR_TXHDB_SHM_EXTDATA_UNMATCH                                            = -0x0000a005 /*-40965*/
	TXHDB_ERR_TXHDB_EXTDATA_CHECK_ERROR                                            = -0x0000a105 /*-41221*/
	TXHDB_ERR_CHUNK_BUFFS_CANNOT_BE_ALLOCED_IF_THEY_ARE_NOT_RELEASED               = -0x0000a205 /*-41477*/
	TXHDB_ERR_ALLOCATE_MEMORY_FAIL                                                 = -0x0000a305 /*-41733*/
	TXHDB_ERR_INVALID_CHUNK_RW_MANNER                                              = -0x0000a405 /*-41989*/
	TXHDB_ERR_FILE_PREAD_NOT_COMPLETE                                              = -0x0000a505 /*-42245*/
	TXHDB_ERR_FILE_PWRITE_NOT_COMPLETE                                             = -0x0000a605 /*-42501*/
	TXHDB_ERR_KEY_ONEBLOCK_BUT_NEXT_NOTNULL                                        = -0x0000a705 /*-42757*/
	TXHDB_ERR_VALUE_ONEBLOCK_BUT_NEXT_NOTNULL                                      = -0x0000a805 /*-43013*/
	TXHDB_ERR_VARINT_FORMAT_ERROR                                                  = -0x0000a905 /*-43269*/
	TXHDB_ERR_TXSTAT_ERROR                                                         = -0x0000aa05 /*-43525*/
	TXHDB_ERR_INVALID_VERSION                                                      = -0x0000ab05 /*-43781*/
	TXHDB_ERR_FREE_BLOCK_NOT_ENOUGH                                                = -0x0000ac05 /*-44037*/

	//Engine SYSTEM (module id 0x07) Error Code defined below
	ENG_ERR_INVALID_ARGUMENTS                        = -0x00000107 /*-263*/
	ENG_ERR_INVALID_MEMBER_VARIABLE_VALUE            = -0x00000207 /*-519*/
	ENG_ERR_NEW_TXHCURSOR_FAILED                     = -0x00000307 /*-775*/
	ENG_ERR_TXHCURSOR_KEY_BUFFER_LEGHTH_NOT_ENOUGH   = -0x00000407 /*-1031*/
	ENG_ERR_TXHCURSOR_VALUE_BUFFER_LEGHTH_NOT_ENOUGH = -0x00000507 /*-1287*/
	ENG_ERR_TXHDB_FILEPATH_NULL                      = -0x00000607 /*-1543*/
	ENG_ERR_TCHDB_RELATED_ERROR                      = -0x00000707 /*-1799*/
	ENG_ERR_NULL_CACHE                               = -0x00000807 /*-2055*/
	ENG_ERR_ITER_FAIL_SYSTEM_RECORD                  = -0x00000907 /*-2311*/
	ENG_ERR_SYSTEM_ERROR                             = -0x00000a07 /*-2567*/
	ENG_ERR_ENGINE_ERROR                             = -0x00000b07 /*-2823*/
	ENG_ERR_DATA_ERROR                               = -0x00000c07 /*-3079*/
	ENG_ERR_VERSION_ERROR                            = -0x00000d07 /*-3335*/
	ENG_ERR_SYSTEM_ERROR_BUFF_OVERFLOW               = -0x00000e07 /*-3591*/
	ENG_ERR_METADATA_ERROR                           = -0x00000f07 /*-3847*/
	ENG_ERR_ADD_KEYMETA_FAILED                       = -0x00001007 /*-4103*/
	ENG_ERR_ADD_VALUEMETA_FAILED                     = -0x00001107 /*-4359*/
	ENG_ERR_RESERVED_FIELDNAME                       = -0x00001207 /*-4615*/
	ENG_ERR_KEYNAME_REPEAT                           = -0x00001307 /*-4871*/
	ENG_ERR_VALUENAME_REPEAT                         = -0x00001407 /*-5127*/
	ENG_ERR_MISS_KEYMETA                             = -0x00001507 /*-5383*/
	ENG_ERR_DELETE_KEYFIELD                          = -0x00001607 /*-5639*/
	ENG_ERR_CHANGE_KEYCOUNT                          = -0x00001707 /*-5895*/
	ENG_ERR_CHANGE_KEYTYPE                           = -0x00001807 /*-6151*/
	ENG_ERR_CHANGE_KEYLENGTH                         = -0x00001907 /*-6407*/
	ENG_ERR_CHANGE_VALUETYPE                         = -0x00001a07 /*-6663*/
	ENG_ERR_CHANGE_VALUELENGTH                       = -0x00001b07 /*-6919*/
	ENG_ERR_CHANGE_DEFAULTVALUE                      = -0x00001c07 /*-7175*/
	ENG_ERR_EMPTY_FIELDNAME                          = -0x00001d07 /*-7431*/
	ENG_ERR_INVALID_TARGET_KEYFIELD                  = -0x00001e07 /*-7687*/
	ENG_ERR_INVALID_TARGET_VALUEFIELD                = -0x00001f07 /*-7943*/
	ENG_ERR_INVALID_TABLE_TYPE                       = -0x00002007 /*-8199*/
	ENG_ERR_CHANGE_TABLE_TYPE                        = -0x00002107 /*-8455*/
	ENG_ERR_MISS_VALUEMETA                           = -0x00002207 /*-8711*/
	ENG_ERR_NOT_ENOUGH_BUFF_FOR_FILEPATH             = -0x00002307 /*-8967*/
	ENG_ERR_ENGINE_FILE_NOT_FOUND                    = -0x00002407 /*-9223*/

	ULOG_ERR_INVALID_PARAMS = -0x00000109 /*-265*/
	//SYNCDB SYSTEM (module id 0x0b) Error Code defined below
	SYNCDB_ERR_INVALID_PARAMS                     = -0x0000010b /*-267*/
	SYNCDB_ERR_PAUSE_TO_SEND_FOR_SWITCH_CONNECTOR = -0x0000020b /*-523*/
	SYNCDB_ERR_CONNECTOR_IS_NOT_CONNECTED         = -0x0000030b /*-779*/

	//TCAPSVR SYSTEM (module id 0x0d) Error Code defined below
	SVR_ERR_FAIL_ROUTE                                    = -0x0000010d /*-269*/
	SVR_ERR_FAIL_TIMEOUT                                  = -0x0000020d /*-525*/
	SVR_ERR_FAIL_SHORT_BUFF                               = -0x0000030d /*-781*/
	SVR_ERR_FAIL_SYSTEM_BUSY                              = -0x0000040d /*-1037*/
	SVR_ERR_FAIL_RECORD_EXIST                             = -0x0000050d /*-1293*/
	SVR_ERR_FAIL_INVALID_FIELD_NAME                       = -0x0000060d /*-1549*/
	SVR_ERR_FAIL_VALUE_OVER_MAX_LEN                       = -0x0000070d /*-1805*/
	SVR_ERR_FAIL_INVALID_FIELD_TYPE                       = -0x0000080d /*-2061*/
	SVR_ERR_FAIL_SYNC_WRITE                               = -0x0000090d /*-2317*/
	SVR_ERR_FAIL_WRITE_RECORD                             = -0x00000a0d /*-2573*/
	SVR_ERR_FAIL_DELETE_RECORD                            = -0x00000b0d /*-2829*/
	SVR_ERR_FAIL_DATA_ENGINE                              = -0x00000c0d /*-3085*/
	SVR_ERR_FAIL_RESULT_OVERFLOW                          = -0x00000d0d /*-3341*/
	SVR_ERR_FAIL_INVALID_OPERATION                        = -0x00000e0d /*-3597*/
	SVR_ERR_FAIL_INVALID_SUBSCRIPT                        = -0x00000f0d /*-3853*/
	SVR_ERR_FAIL_INVALID_INDEX                            = -0x0000100d /*-4109*/
	SVR_ERR_FAIL_OVER_MAXE_FIELD_NUM                      = -0x0000110d /*-4365*/
	SVR_ERR_FAIL_MISS_KEY_FIELD                           = -0x0000120d /*-4621*/
	SVR_ERR_FAIL_NEED_SIGNUP                              = -0x0000130d /*-4877*/
	SVR_ERR_FAIL_CROSS_AUTH                               = -0x0000140d /*-5133*/
	SVR_ERR_FAIL_SIGNUP_FAIL                              = -0x0000150d /*-5389*/
	SVR_ERR_FAIL_SIGNUP_INVALID                           = -0x0000160d /*-5645*/
	SVR_ERR_FAIL_SIGNUP_INIT                              = -0x0000170d /*-5901*/
	SVR_ERR_FAIL_LIST_FULL                                = -0x0000180d /*-6157*/
	SVR_ERR_FAIL_LOW_VERSION                              = -0x0000190d /*-6412*/
	SVR_ERR_FAIL_HIGH_VERSION                             = -0x00001a0d /*-6669*/
	SVR_ERR_FAIL_INVALID_RESULT_FLAG                      = -0x00001b0d /*-6925*/
	SVR_ERR_FAIL_PROXY_STOPPING                           = -0x00001c0d /*-7181*/
	SVR_ERR_FAIL_SVR_READONLY                             = -0x00001d0d /*-7437*/
	SVR_ERR_FAIL_SVR_READONLY_BECAUSE_IN_SLAVE_MODE       = -0x00001e0d /*-7693*/
	SVR_ERR_FAIL_INVALID_VERSION                          = -0x00001f0d /*-7949*/
	SVR_ERR_FAIL_SYSTEM_ERROR                             = -0x0000200d /*-8205*/
	SVR_ERR_FAIL_OVERLOAD                                 = -0x0000210d /*-8461*/
	SVR_ERR_FAIL_NOT_ENOUGH_DADADISK_SPACE                = -0x0000220d /*-8717*/
	SVR_ERR_FAIL_NOT_ENOUGH_ULOGDISK_SPACE                = -0x0000230d /*-8973*/
	SVR_ERR_FAIL_UNSUPPORTED_PROTOCOL_MAGIC               = -0x0000240d /*-9229*/
	SVR_ERR_FAIL_UNSUPPORTED_PROTOCOL_CMD                 = -0x0000250d /*-9485*/
	SVR_ERR_FAIL_HIGH_TABLE_META_VERSION                  = -0x0000260d /*-9741*/
	SVR_ERR_FAIL_MERGE_VALUE_FIELD                        = -0x0000270d /*-9997*/
	SVR_ERR_FAIL_CUT_VALUE_FIELD                          = -0x0000280d /*-10253*/
	SVR_ERR_FAIL_PACK_FIELD                               = -0x0000290d /*-10509*/
	SVR_ERR_FAIL_UNPACK_FIELD                             = -0x00002a0d /*-10765*/
	SVR_ERR_FAIL_LOW_API_VERSION                          = -0x00002b0d /*-11021*/
	SVR_ERR_COMMAND_AND_TABLE_TYPE_IS_MISMATCH            = -0x00002c0d /*-11277*/
	SVR_ERR_FAIL_TO_FIND_CACHE                            = -0x00002d0d /*-11533*/
	SVR_ERR_FAIL_TO_FIND_META                             = -0x00002e0d /*-11789*/
	SVR_ERR_FAIL_TO_GET_CURSOR                            = -0x00002f0d /*-12045*/
	SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE                    = -0x0000300d /*-12301*/
	SVR_ERR_INVALID_ARGUMENTS                             = -0x0000310d /*-12557*/
	SVR_ERR_SLAVE_READ_INVALID                            = -0x0000320d /*-12813*/
	SVR_ERR_NULL_CACHE                                    = -0x0000330d /*-13069*/
	SVR_ERR_NULL_CURSOR                                   = -0x0000340d /*-13325*/
	SVR_ERR_METALIB_VERSION_LESS_THAN_ENTRY_VERSION       = -0x0000350d /*-13581*/
	SVR_ERR_INVALID_SELECT_ID_FOR_UNION                   = -0x0000360d /*-13837*/
	SVR_ERR_CAN_NOT_FIND_SELECT_ENTRY_FOR_UNION           = -0x0000370d /*-14093*/
	SVR_ERR_FAIL_DOCUMENT_PACK_VERSION                    = -0x0000380d /*-14349*/
	SVR_ERR_TCAPSVR_PROCESS_NOT_NORMAL                    = -0x0000390d /*-14605*/
	SVR_ERR_TBUSD_PROCESS_NOT_NORMAL                      = -0x00003a0d /*-14861*/
	SVR_ERR_INVALID_ARRAY_COUNT                           = -0x00003b0d /*-15117*/
	SVR_ERR_REJECT_REQUEST_BECAUSE_ROUTE_IN_REJECT_STATUS = -0x00003c0d /*-15373*/
	SVR_ERR_FAIL_GET_ROUTE_HASH_CODE                      = -0x00003d0d /*-15629*/
	SVR_ERR_FAIL_INVALID_FIELD_VALUE                      = -0x00003e0d /*-15885*/
	SVR_ERR_FAIL_PROTOBUF_FIELD_GET                       = -0x00003f0d /*-16141*/
	SVR_ERR_FAIL_PROTOBUF_VALUE_BUFF_EXCEED               = -0x0000400d /*-16397*/
	SVR_ERR_FAIL_PROTOBUF_FIELD_UPDATE                    = -0x0000410d /*-16653*/
	SVR_ERR_FAIL_PROTOBUF_FIELD_INCREASE                  = -0x0000420d /*-16909*/
	SVR_ERR_FAIL_PROTOBUF_FIELD_TAG_MISMATCH              = -0x0000430d /*-17165*/
	SVR_ERR_FAIL_BINLOG_SEQUENCE_TOO_SMALL                = -0x0000440d
	SVR_ERR_FAIL_SVR_IS_NOT_MASTER                        = -0x0000450d
	SVR_ERR_FAIL_BINLOG_INVALID_FILE_PATH                 = -0x0000460d
	SVR_ERR_FAIL_BINLOG_SOCKET_SEND_BUFF_IS_FULL          = -0x0000470d
	SVR_ERR_FAIL_DOCUMENT_NOT_SUPPORT                     = -0x0000480d /*-18445*/
	SVR_ERR_FAIL_PARTKEY_INSERT_NOT_SUPPORT               = -0x0000490d /*-18701*/
	SVR_ERR_FAIL_SQL_FILTER_FAILED                        = -0x00004a0d /*-18957*/
	SVR_ERR_FAIL_NOT_MATCHED_SQL_QUERY_CONDITION          = -0x00004b0d /*-19213*/

	//TCAPDB SYSTEM (module id 0x0f) Error Code defined below
	TCAPDB_ERR_INVALID_PARAMS               = -0x0000010f /*-271*/
	TCAPDB_ERR_ALLOCATE_MEMORY_FAILED       = -0x0000020f /*-527*/
	TCAPDB_ERR_INDEX_SERVER_RETURN_EXISTED  = -0x0000030f /*-783*/
	TCAPDB_ERR_INDEX_SERVER_RETURN_NOT_FIND = -0x0000040f /*-1039*/
	TCAPDB_ERR_INDEX_SERVER_RETURN_OVERLOAD = -0x0000050f /*-1295*/
	TCAPDB_ERR_PACK_FAILED                  = -0x0000060f /*-1551*/
	TCAPDB_ERR_TIMEOUT                      = -0x0000070f /*-1807*/
	TCAPDB_ERR_REJECT_REQ                   = -0x0000080f /*-2063*/

	//TCAPROXY SYSTEM (module id 0x11) Error String defined below
	PROXY_ERR_INVALID_PARAMS                                                = -0x00000111 /*-273*/
	PROXY_ERR_NO_NEED_ROUTE_BATCHGET_ACTION_MSG_WHEN_NODE_IS_IN_SYNC_STATUS = -0x00000211 /*-529*/
	PROXY_ERR_NO_NEED_ROUTE_WHEN_NODE_IS_IN_REJECT_STATUS                   = -0x00000311 /*-785*/
	PROXY_ERR_PROBE_TIMEOUT                                                 = -0x00000411 /*-1041*/
	PROXY_ERR_SYSTEM_ERROR                                                  = -0x00000511 /*-1297*/
	PROXY_ERR_CONFIG_ERROR                                                  = -0x00000611 /*-1553*/
	PROXY_ERR_OVER_MAX_NODE                                                 = -0x00000711 /*-1809*/
	PROXY_ERR_INVALID_SPLIT_SIZE                                            = -0x00000811 /*-2065*/
	PROXY_ERR_INVALID_ROUTE_INDEX                                           = -0x00000911 /*-2321*/
	PROXY_ERR_CONNECT_SERVER                                                = -0x00000a11 /*-2577*/
	PROXY_ERR_COMPOSE_MSG                                                   = -0x00000b11 /*-2833*/
	PROXY_ERR_ROUTE_MSG                                                     = -0x00000c11 /*-3089*/
	PROXY_ERR_SHORT_BUFFER                                                  = -0x00000d11 /*-3345*/
	PROXY_ERR_OVER_MAX_RECORD                                               = -0x00000e11 /*-3601*/
	PROXY_ERR_INVALID_SERVICE_TABLE                                         = -0x00000f11 /*-3857*/
	PROXY_ERR_REGISTER_FAILED                                               = -0x00001011 /*-4113*/
	PROXY_ERR_CREATE_SESSION_HASH                                           = -0x00001111 /*-4369*/
	PROXY_ERR_WRONG_STATUS                                                  = -0x00001211 /*-4625*/
	PROXY_ERR_UNPACK_MSG                                                    = -0x00001311 /*-4881*/
	PROXY_ERR_PACK_MSG                                                      = -0x00001411 /*-5137*/
	PROXY_ERR_SEND_MSG                                                      = -0x00001511 /*-5393*/
	PROXY_ERR_ALLOCATE_MEMORY                                               = -0x00001611 /*-5649*/
	PROXY_ERR_PARSE_MSG                                                     = -0x00001711 /*-5905*/
	PROXY_ERR_INVALID_MSG                                                   = -0x00001811 /*-6161*/
	PROXY_ERR_FAILED_PROC_REQUEST_BECAUSE_NODE_IS_IN_SYNC_STASUS            = -0x00001911 /*-6417*/
	PROXY_ERR_KEY_FIELD_NUM_IS_ZERO                                         = -0x00001a11 /*-6673*/
	PROXY_ERR_LACK_OF_SOME_KEY_FIELDS                                       = -0x00001b11 /*-6929*/
	PROXY_ERR_FAILED_TO_FIND_NODE                                           = -0x00001c11 /*-7185*/
	PROXY_ERR_INVALID_COMPRESS_TYPE                                         = -0x00001d11 /*-7441*/
	PROXY_ERR_REQUEST_OVERSPEED                                             = -0x00001e11 /*-7697*/
	PROXY_ERR_SWIFT_TIMEOUT                                                 = -0x00001f11 /*-7953*/
	PROXY_ERR_SWIFT_ERROR                                                   = -0x00002011 /*-8209*/
	PROXY_ERR_DIRECT_RESPONSE                                               = -0x00002111 /*-8465*/
	PROXY_ERR_INIT_TLOG                                                     = -0x00002211 /*-8721*/
	PROXY_ERR_ASSISTANT_THREAD_NOT_RUN                                      = -0x00002311 /*-8977*/
	PROXY_ERR_REQUEST_ACCESS_CTRL_REJECT                                    = -0x00002411 /*-9233*/
	PROXY_ERR_NOT_ALL_NODES_ARE_IN_NORMAL_OR_WAIT_STATUS                    = -0x00002511 /*-9489*/
	PROXY_ERR_ALREADY_CACHED_REQUEST_TIMEOUT                                = -0x00002611 /*-9745*/
	PROXY_ERR_FAILED_TO_CACHE_REQUEST                                       = -0x00002711 /*-10001*/
	PROXY_ERR_NOT_EXIST_CACHED_REQUEST                                      = -0x00002811 /*-10257*/
	PROXY_ERR_FAILED_NOT_ENOUGH_CACHE_BUFF                                  = -0x00002911 /*-10513*/
	PROXY_ERR_FAILED_PROCESS_CACHED_REQUEST                                 = -0x00002a11 /*-10769*/
	PROXY_ERR_SYNC_ROUTE_HAS_BEEN_CANCELLED                                 = -0x00002b11 /*-11025*/
	PROXY_ERR_FAILED_LOCK_CACHE                                             = -0x00002c11 /*-11281*/
	PROXY_ERR_SWIFT_SEND_BUFFER_FULL                                        = -0x00002d11 /*-11537*/
	PROXY_ERR_REQUEST_OVERLOAD_CTRL_REJECT                                  = -0x00002e11 /*-11793*/
	PROXY_ERR_SQL_QUERY_MGR_IS_NULL                                         = -0x00002f11 /*-12049*/
	PROXY_ERR_SQL_QUERY_INVALID_SQL_TYPE                                    = -0x00003011 /*-12305*/
	PROXY_ERR_GET_TRANSACTION_FAILED                                        = -0x00003111 /*-12561*/
	PROXY_ERR_ADD_TRANSACTION_FAILED                                        = -0x00003211 /*-12817*/
	PROXY_ERR_QUERY_FROM_INDEX_SERVER_FAILED                                = -0x00003311 /*-13073*/
	PROXY_ERR_QUERY_FROM_INDEX_SERVER_TIMEOUT                               = -0x00003411 /*-13329*/
	PROXY_ERR_QUERY_FOR_CONVERT_TCAPLUS_REQ_TO_INDEX_SERVER_REQ_FAILED      = -0x00003511 /*-13585*/
	PROXY_ERR_QUERY_INDEX_FIELD_NOT_EXIST                                   = -0x00003611 /*-13841*/
	PROXY_ERR_THIS_SQL_IS_NOT_SUPPORT                                       = -0x00003711;/*-14097*/

	//API SYSTEM (module id 0x13) Error Code defined below
	API_ERR_OVER_MAX_KEY_FIELD_NUM                           = -0x00000113 /*-275*/
	API_ERR_OVER_MAX_VALUE_FIELD_NUM                         = -0x00000213 /*-531*/
	API_ERR_OVER_MAX_FIELD_NAME_LEN                          = -0x00000313 /*-787*/
	API_ERR_OVER_MAX_FIELD_VALUE_LEN                         = -0x00000413 /*-1043*/
	API_ERR_FIELD_NOT_EXSIST                                 = -0x00000513 /*-1299*/
	API_ERR_FIELD_TYPE_NOT_MATCH                             = -0x00000613 /*-1555*/
	API_ERR_PARAMETER_INVALID                                = -0x00000713 /*-1811*/
	API_ERR_OPERATION_TYPE_NOT_MATCH                         = -0x00000813 /*-2067*/
	API_ERR_PACK_MESSAGE                                     = -0x00000913 /*-2323*/
	API_ERR_UNPACK_MESSAGE                                   = -0x00000a13 /*-2579*/
	API_ERR_PACKAGE_NOT_UNPACKED                             = -0x00000b13 /*-2835*/
	API_ERR_OVER_MAX_RECORD_NUM                              = -0x00000c13 /*-3091*/
	API_ERR_INVALID_COMMAND                                  = -0x00000d13 /*-3347*/
	API_ERR_NO_MORE_RECORD                                   = -0x00000e13 /*-3603*/
	API_ERR_OVER_KEY_FIELD_NUM                               = -0x00000f13 /*-3859*/
	API_ERR_OVER_VALUE_FIELD_NUM                             = -0x00001013 /*-4115*/
	API_ERR_OBJ_NEED_INIT                                    = -0x00001113 /*-4371*/
	API_ERR_INVALID_DATA_SIZE                                = -0x00001213 /*-4627*/
	API_ERR_INVALID_ARRAY_COUNT                              = -0x00001313 /*-4883*/
	API_ERR_INVALID_UNION_SELECT                             = -0x00001413 /*-5139*/
	API_ERR_MISS_PRIMARY_KEY                                 = -0x00001513 /*-5395*/
	API_ERR_UNSUPPORT_FIELD_TYPE                             = -0x00001613 /*-5651*/
	API_ERR_ARRAY_BUFFER_IS_SMALL                            = -0x00001713 /*-5907*/
	API_ERR_IS_NOT_WHOLE_PACKAGE                             = -0x00001813 /*-6163*/
	API_ERR_MISS_PAIR_FIELD                                  = -0x00001913 /*-6419*/
	API_ERR_GET_META_ENTRY                                   = -0x00001a13 /*-6675*/
	API_ERR_GET_ARRAY_META                                   = -0x00001b13 /*-6931*/
	API_ERR_GET_ENTRY_META                                   = -0x00001c13 /*-7187*/
	API_ERR_INCOMPATIBLE_META                                = -0x00001d13 /*-7443*/
	API_ERR_PACK_ARRAY_DATA                                  = -0x00001e13 /*-7669*/
	API_ERR_PACK_UNION_DATA                                  = -0x00001f13 /*-7955*/
	API_ERR_PACK_STRUCT_DATA                                 = -0x00002013 /*-8211*/
	API_ERR_UNPACK_ARRAY_DATA                                = -0x00002113 /*-8467*/
	API_ERR_UNPACK_UNION_DATA                                = -0x00002213 /*-8723*/
	API_ERR_UNPACK_STRUCT_DATA                               = -0x00002313 /*-8979*/
	API_ERR_INVALID_INDEX_NAME                               = -0x00002413 /*-9235*/
	API_ERR_MISS_PARTKEY_FIELD                               = -0x00002513 /*-9491*/
	API_ERR_ALLOCATE_MEMORY                                  = -0x00002613 /*-9747*/
	API_ERR_GET_META_SIZE                                    = -0x00002713 /*-10003*/
	API_ERR_MISS_BINARY_VERSION                              = -0x00002813 /*-10259*/
	API_ERR_INVALID_INCREASE_FIELD                           = -0x00002913 /*-10515*/
	API_ERR_INVALID_RESULT_FLAG                              = -0x00002a13 /*-10771*/
	API_ERR_OVER_MAX_LIST_INDEX_NUM                          = -0x00002b13 /*-11027*/
	API_ERR_INVALID_OBJ_STATUE                               = -0x00002c13 /*-11283*/
	API_ERR_INVALID_REQUEST                                  = -0x00002d13 /*-11539*/
	API_ERR_INVALID_SHARD_LIST                               = -0x00002e13 /*-11795*/
	API_ERR_TABLE_NAME_MISSING                               = -0x00002f13 /*-12051*/
	API_ERR_SOCKET_SEND_BUFF_IS_FULL                         = -0x00003013 /*-12307*/
	API_ERR_INVALID_MAGIC                                    = -0x00003113 /*-12563*/
	API_ERR_TABLE_IS_NOT_EXIST                               = -0x00003213 /*-12819*/
	API_ERR_SHORT_BUFF                                       = -0x00003313 /*-13075*/
	API_ERR_FLOW_CONTROL                                     = -0x00003413 /*-13331*/
	API_ERR_COMPRESS_SWITCH_NOT_SUPPORTED_REGARDING_THIS_CMD = -0x00003513 /*-13587*/
	API_ERR_FAILED_TO_FIND_ROUTE                             = -0x00003613 /*-13843*/
	API_ERR_OVER_MAX_PKG_SIZE                                = -0x00003713 /*-14099*/
	API_ERR_INVALID_VERSION_FOR_TLV                          = -0x00003813 /*-14355*/
	API_ERR_BSON_SERIALIZE                                   = -0x00003913 /*-14611*/
	API_ERR_BSON_DESERIALIZE                                 = -0x00003a13 /*-14867*/
	API_ERR_ADD_RECORD                                       = -0x00003b13 /*-15123*/
	API_ERR_ZONE_IS_NOT_EXIST                                = -0x00003c13 /*-15379*/
	API_ERR_TRAVERSER_IS_NOT_EXIST                           = -0x00003d13 /*-15635*/
	API_ERR_INSTANCE_ID_FULL                                 = -0x00003e13 /*-15891*/
	API_ERR_INSTANCE_INIT_LOG_FAILURE                        = -0x00003f13 /*-16147*/
	API_ERR_CONNECTOR_IS_ABNORMAL                            = -0x00004013 /*-16403*/
	API_ERR_WAIT_RSP_TIMEOUT                                 = -0x00004113 /*-16659*/

	//TCAPCENTER SYSTEM (module id 0x15) Error String defined below
	CENTER_ERR_INVALID_PARAMS      = -0x00000115 /*-277*/
	CENTER_ERR_TABLE_ALREADY_EXIST = -0x00000215 /*-533*/
	CENTER_ERR_TABLE_NOT_EXIST     = -0x00000315 /*-789*/

	//TCAPDIR SYSTEM (module id 0x17) Error Code defined below
	DIR_ERR_SIGN_FAIL                  = -0x00000117 /*-279*/
	DIR_ERR_LOW_VERSION                = -0x00000217 /*-535*/
	DIR_ERR_HIGH_VERSION               = -0x00000317 /*-791*/
	DIR_ERR_GET_DIR_SERVER_LIST        = -0x00000417 /*-1047*/
	DIR_ERR_APP_IS_NOT_FOUNT           = -0x00000517 /*-1303*/
	DIR_ERR_NOT_CONNECT_TCAPCENTER     = -0x00000617 /*-1559*/
	DIR_ERR_ZONE_IS_NOT_FOUNT          = -0x00000717 /*-1815*/
	DIR_ERR_HASH_TABLE_FAILED          = -0x00000817 /*-2071*/
	DIR_ERR_GET_TABLE_AND_ACCESS_LIST  = -0x00000917 /*-2327*/
	DIR_ERR_IS_NOT_THIS_ZONE_API       = -0x00000A17 /*-2583*/
	DIR_ERR_IS_NOT_IN_ZONES_WHITE_LIST = -0x0000FF03 /*-65283*/

	//BSON ERROR
	BSON_ERR_TYPE_IS_NOT_MATCH                          = -0x00000118 /*-280*/
	BSON_ERR_INVALID_DATA_TYPE                          = -0x00000218 /*-536*/
	BSON_ERR_INVALID_VALUE                              = -0x00000318 /*-792*/
	BSON_ERR_BSON_TYPE_UNMATCH_TDR_TYPE                 = -0x00000418 /*-1048*/
	BSON_ERR_BSON_TYPE_IS_NOT_SUPPORT_BY_TCAPLUS        = -0x00000518 /*-1304*/
	BSON_ERR_BSON_ARRAY_COUNT_IS_INVALID                = -0x00000618 /*-1560*/
	BSON_ERR_FAILED_TO_PARSE                            = -0x00000718 /*-1816*/
	BSON_ERR_INVALID_FIELD_NAME_LENGTH                  = -0x00000818 /*-2072*/
	BSON_ERR_INDEX_FIELD_NAME_NOT_EXIST_WITH_ARRAY_TYPE = -0x00000918 /*-2328*/
	BSON_ERR_INVALID_ARRAY_INDEX                        = -0x00000a18 /*-2584*/
	BSON_ERR_TDR_META_LIB_IS_NULL                       = -0x00000b18 /*-2840*/
	BSON_ERR_MATCHED_COUNT_GREATER_THAN_ONE             = -0x00000c18 /*-3096*/
	BSON_ERR_NO_MATCHED                                 = -0x00000d18 /*-3352*/
	BSON_ERR_GREATER_THAN_ARRAY_MAX_COUNT               = -0x00000f18 /*-3864*/
	BSON_ERR_BSON_EXCEPTION                             = -0x00001018 /*-4120*/
	BSON_ERR_STD_EXCEPTION                              = -0x00001118 /*-4376*/
	BSON_ERR_INVALID_KEY                                = -0x00001218 /*-4632*/
	BSON_ERR_TDR_META_LIB_IS_INVALID                    = -0x00001318 /*-4888*/

	//TCAPTCAPCOMMON SYSTEM (module id 0x19) Error Code defined below
	COMMON_ERR_INVALID_ARGUMENTS                = -0x00000119 /*-281*/
	COMMON_ERR_INVALID_MEMBER_VARIABLE_VALUE    = -0x00000219 /*-537*/
	COMMON_ERR_SPINLOCK_INIT_FAIL               = -0x00000319 /*-793*/
	COMMON_ERR_SPINLOCK_DESTROY_FAIL            = -0x00000419 /*-1049*/
	COMMON_ERR_COMPRESS_BUF_NOT_ENOUGH          = -0x00000519 /*-1305*/
	COMMON_ERR_DECOMPRESS_BUF_NOT_ENOUGH        = -0x00000619 /*-1561*/
	COMMON_ERR_DECOMPRESS_INVALID_INPUT         = -0x00000719 /*-1817*/
	COMMON_ERR_CANNOT_FIND_COMPRESS_ALGORITHM   = -0x00000819 /*-2073*/
	COMMON_ERR_CANNOT_FIND_DECOMPRESS_ALGORITHM = -0x00000919 /*-2329*/
	COMMON_ERR_COMPRESS_FAIL                    = -0x00000a19 /*-2585*/
	COMMON_ERR_DECOMPRESS_FAIL                  = -0x00000b19 /*-2841*/
	COMMON_ERR_INVALID_SWITCH_VALUE             = -0x00000c19 /*-3097*/
	COMMON_ERR_LINUX_SYSTEM_CALL_FAIL           = -0x00000d19 /*-3353*/
	COMMON_ERR_NOT_FIND_STAT_CACHE_VALUE        = -0x00000e19 /*-3609*/
	COMMON_ERR_LZO_CHECK_FAIL                   = -0x00000f19 /*-3865*/

	//tcaplus_index SYSTEM (module id 0x1a) Error Code defined below
	TCAPLUS_INDEX_ERR_INVALID_PARAMS                                         = -0x0000011a /*-282*/
	TCAPLUS_INDEX_ERR_ALLOCATE_MEMORY                                        = -0x0000021a /*-538*/
	TCAPLUS_INDEX_ERR_CREATE_CONNECTOR_TO_INDEX_SERVER_FAILED                = -0x0000031a /*-794*/
	TCAPLUS_INDEX_ERR_SEND_TO_INDEX_SERVER_FAILED_FOR_NO_CONNECTOR           = -0x0000041a /*-1050*/
	TCAPLUS_INDEX_ERR_SEND_TO_INDEX_SERVER_FAILED_FOR_NO_AVAILABLE_CONNECTOR = -0x0000051a /*-1306*/
	TCAPLUS_INDEX_ERR_SEND_TO_INDEX_SERVER_FAILED_FOR_OTHER_REASON           = -0x0000061a /*-1562*/
	TCAPLUS_INDEX_ERR_PAUSE_SEND_FOR_CHANGIN_URL_STATUS                      = -0x0000071a /*-1818*/
	TCAPLUS_INDEX_ERR_QUERY_INDEX_SERVER_OVERLOAD                            = -0x0000081a /*-2074*/

	// Non-error (for information purpose)
	COMMON_INFO_DATA_NOT_MODIFIED = 0x00000120 /*288*/

)

var ErrorCodes = map[int]string{
	DirSignUpFailed:        "dir 认证失败",
	ClientInitTimeOut:      "client 初始化超时",
	ProxySignUpFailed:      "proxy 认证失败",
	ZoneIdNotExist:         "注册zoneId不存在",
	TableNotExist:          "注册的表不存在",
	InvalidCmd:             "request/response命令非法",
	InvalidPolicy:          "request版本策略非法",
	RecordToMax:            "request中record数量超过上限",
	KeyNameLenOverMax:      "record中key名称长度超限",
	KeyLenOverMax:          "record中key值长度超限",
	KeyNumOverMax:          "record中key数量超限",
	ValueNameLenOverMax:    "record中value名称长度超限",
	ValueLenOverMax:        "record中value值长度超限",
	ValueNumOverMax:        "record中value数量超限",
	ValuePackOverMax:       "record中value字段打包超限",
	RecordNumOverMax:       "request中record数量超限",
	ProxyNotAvailable:      "没有可用的proxy",
	RequestHasNoRecord:     "请求中没有添加record",
	RequestHasNoKeyField:   "请求中没有key字段",
	RecordKeyTypeInvalid:   "record Key 类型错误",
	RecordValueTypeInvalid: "record Value 类型错误",
	OperationNotSupport:    "操作不支持",
	ClientNotInit:          "client 未初始化",
	RecordUnpackFailed:     "response record解包失败",
	RecordKeyNotExist:      "record中请求的key不存在",
	RecordValueNotExist:    "record中请求的value不存在",
	ClientNotDial:			"client 未进行dial初始化",
	RespNotMatchReq:		"响应未与请求对应",
	MetadataNotProtobuf:	"元数据类型不是protobuf",
	SqlQueryFormatError:	"sql语句格式错误",

	/*****************************************************************************************
	*****************************************C版本错误码*********************************************
	*******************************************************************************************/
	//TXHDB (module id 0x05) Error String defined below
	TXHDB_ERR_RECORD_NOT_EXIST:                                                     "txhdb_record_not_exist",
	TXHDB_ERR_ITERATION_NO_MORE_RECORDS:                                            "txhdb_iteration_no_more_record",
	TXHDB_ERR_MUTEX_TRYLOCK_BUSY:                                                   "txhdb_mutex_trylock_busy",
	TXHDB_ERR_MUTEX_TIMEDLOCK_TIMEOUT:                                              "txhdb_mutex_timedlock_timeout",
	TXHDB_ERR_RWLOCK_TRYWRLOCK_BUSY:                                                "txhdb_rwlock_trywrlock_busy",
	TXHDB_ERR_RWLOCK_TRYRDLOCK_BUSY:                                                "txhdb_rwlock_tryrdlock_busy",
	TXHDB_ERR_SPIN_TRYLOCK_BUSY:                                                    "txhdb_spin_trylock_busy",
	TXHDB_ERR_ITERATION_EXCEED_MAX_ALLOWED_TIME_OF_ONE_ITER:
		"txhdb_err_iteration_exceed_max_allowed_time_of_one_iter",
	TXHDB_ERR_INVALID_ARGUMENTS:                                                    "txhdb_invalid_arguments",
	TXHDB_ERR_INVALID_MEMBER_VARIABLE_VALUE:
		"txhdb_invalid_member_variable_value",
	TXHDB_ERR_ALREADY_OPEN:                                                         "txhdb_already_opened",
	TXHDB_ERR_MUTEX_LOCK_FAIL:                                                      "txhdb_mutex_lock_fail",
	TXHDB_ERR_MUTEX_TRYLOCK_FAIL:                                                   "txhdb_mutex_trylock_fail",
	TXHDB_ERR_MUTEX_TIMEDLOCK_FAIL:                                                 "txhdb_mutex_timedlock_fail",
	TXHDB_ERR_MUTEX_UNLOCK_FAIL:                                                    "txhdb_mutex_unlock_fail",
	TXHDB_ERR_RWLOCK_WRLOCK_FAIL:                                                   "txhdb_rwlock_wrlock_fail",
	TXHDB_ERR_RWLOCK_TRYWRLOCK_FAIL:                                                "txhdb_rwlock_trywrlock_fail",
	TXHDB_ERR_RWLOCK_RDLOCK_FAIL:                                                   "txhdb_rwlock_rdlock_fail",
	TXHDB_ERR_RWLOCK_TRYRDLOCK_FAIL:                                                "txhdb_rwlock_tryrdlock_fail",
	TXHDB_ERR_RWLOCK_UNLOCK_FAIL:                                                   "txhdb_rwlock_unlock_fail",
	TXHDB_ERR_SPIN_LOCK_FAIL:                                                       "txhdb_spin_lock_fail",
	TXHDB_ERR_SPIN_UNLOCK_FAIL:                                                     "txhdb_spin_unlock_fail",
	TXHDB_ERR_FILE_EXISTS_BUT_STATUS_ERROR:                                         "Txhdb_file_exists_but_status_error",
	TXHDB_ERR_FILE_OPEN_FAIL:                                                       "Txhdb_file_open_fail",
	TXHDB_ERR_FILE_READ_SIZE_INVALID:                                               "Txhdb_file_read_size_INVALID",
	TXHDB_ERR_FILE_INVALID_FILE_PATH:                                               "Txhdb_file_invalid_file_path",
	TXHDB_ERR_FILE_LOCK_FILE_FAIL:                                                  "Txhdb_file_lock_file_fail",
	TXHDB_ERR_FILE_NOT_A_REGULAR_FILE:                                              "Txhdb_file_not_a_regular_file",
	TXHDB_ERR_FILE_MMAP_FAIL:                                                       "Txhdb_file_mmap_fail",
	TXHDB_ERR_FILE_MUNMAP_FAIL:                                                     "Txhdb_file_munmap_fail",
	TXHDB_ERR_FILE_CLOSE_FAIL:                                                      "Txhdb_file_close_fail",
	TXHDB_ERR_FILE_SPACE_NOT_ENOUGH_IN_HEAD:
		"Txhdb_file_space_not_enough_in_head",
	TXHDB_ERR_FILE_FTRUNCATE_FAIL:                                                  "Txhdb_file_ftruncate_fail",
	TXHDB_ERR_FILE_INCONSISTANT_FILE_SIZE:                                          "Txhdb_file_inconsistant_file_size",
	TXHDB_ERR_FILE_MSIZ_LESSER_THAN_TXHDB_WHOLE_REC_OFFSET:
		"Txhdb_file_msiz_lesser_than_txhdb_whole_rec_offset",
	TXHDB_ERR_FILE_MSIZ_CHANGE_NOT_PERMIT:                                          "Txhdb_file_msiz_change_not_permit",
	TXHDB_ERR_FILE_FSTAT_FAIL:                                                      "Txhdb_file_fstat_fail",
	TXHDB_ERR_FILE_MSYNC_FAIL:                                                      "Txhdb_file_msync_fail",
	TXHDB_ERR_FILE_FSYNC_FAIL:                                                      "Txhdb_file_fsync_fail",
	TXHDB_ERR_FILE_FCNTL_LOCK_FILE_FAIL:                                            "Txhdb_file_fcntl_lock_file_fail",
	TXHDB_ERR_FILE_FCNTL_UNLOCK_FILE_FAIL:                                          "Txhdb_file_fcntl_unlock_file_fail",
	TXHDB_ERR_FILE_PREAD_FAIL_WITH_SPECIFIED_ERRNO:
		"txhdb_file_pread_fail_with_specified_errno",
	TXHDB_ERR_FILE_PREAD_FAIL_WITH_UNSPECIFIED_ERRNO:
		"txhdb_file_pread_fail_with_unspecified_errno",
	TXHDB_ERR_FILE_PWRITE_FAIL_WITH_SPECIFIED_ERRNO:
		"txhdb_file_pwrite_fail_with_specified_errno",
	TXHDB_ERR_FILE_PWRITE_FAIL_WITH_UNSPECIFIED_ERRNO:
		"txhdb_file_pwrite_fail_with_unspecified_errno",
	TXHDB_ERR_FILE_READ_EXCEED_FILE_BOUNDARY:                                       "txhdb_read_exceed_file_boundary",
	TXHDB_ERR_FILE_READ_FAIL_DURING_COPY:                                           "txhdb_file_read_fail_during_copy",
	TXHDB_ERR_FILE_WRITE_FAIL_DURING_COPY:                                          "txhdb_file_write_fail_during_copy",
	TXHDB_ERR_FILE_INVALID_FREE_BLOCK_POOL_METADATA:
		"Txhdb_file_invalid_free_block_pool_metadata",
	TXHDB_ERR_FILE_INVALID_MAGIC:                                                   "Txhdb_file_invalid_magic",
	TXHDB_ERR_FILE_INVALID_LIBRARY_VERSION:                                         "Txhdb_file_invalid_library_version",
	TXHDB_ERR_FILE_INVALID_LIBRARY_REVISION:                                        "Txhdb_file_invalid_library_revision",
	TXHDB_ERR_FILE_INVALID_FORMAT_VERSION:                                          "Txhdb_file_invalid_format_version",
	TXHDB_ERR_FILE_INVALID_EXTDATA_FORMAT_VERSION:
		"Txhdb_file_invalid_extdata_format_version",
	TXHDB_ERR_FILE_INVALID_DBTYPE:                                                  "Txhdb_file_invalid_dbtype",
	TXHDB_ERR_FILE_HEAD_CRC_UNMATCH:                                                "Txhdb_file_head_crc_unmatch",
	TXHDB_ERR_FILE_INVALID_METADATA:
		"Txhdb_txhdb_err_file_invalid_metadata",
	TXHDB_ERR_FILE_INVALID_HEADLEN:                                                 "Txhdb_txhdb_err_file_invalid_headlen",
	TXHDB_ERR_FILE_DESERIAL_HEAD_SPACE_NOT_ENOUGH:
		"Txhdb_file_deserialhead_space_not_enough",
	TXHDB_ERR_FILE_SERIAL_HEAD_SPACE_NOT_ENOUGH:
		"Txhdb_file_serialhead_space_not_enough",
	TXHDB_ERR_FILE_DESERIAL_STAT_SPACE_NOT_ENOUGH:
		"Txhdb_file_deserialstat_space_not_enough",
	TXHDB_ERR_FILE_SERIAL_STAT_SPACE_NOT_ENOUGH:
		"Txhdb_file_serialstat_space_not_enough",
	TXHDB_ERR_FILE_SERIAL_FREE_BLOCK_LIST_INFO_WRONG_BUFFLEN:
		"txhdb_file_serial_free_block_list_info_wrong_bufflen",
	TXHDB_ERR_FILE_IN_EXCEPTIONAL_STATUS:                                           "txhdb_file_in_exceptional_status",
	TXHDB_ERR_DB_NOT_OPENED:                                                        "Txhdb_not_opened",
	TXHDB_ERR_DB_WRITE_NOT_PERMIT:                                                  "Txhdb_db_write_not_permit",
	TXHDB_ERR_INVALID_OFFSET_FROM_BUCKET:                                           "Txhdb_invalid_offset_from_bucket",
	TXHDB_ERR_READ_EXTDATA_EXCEED_BUFF_LENGTH:
		"Txhdb_read_extdata_exceed_buff_length",
	TXHDB_ERR_WRITE_EXTDATA_EXCEED_BUFF_LENGTH:
		"Txhdb_write_extdata_exceed_buff_length",
	TXHDB_ERR_FREE_BLOCK_IS_READ_WHEN_GETTING_RECORD:
		"Txhdb_free_block_is_read_when_getting_record",
	TXHDB_ERR_INVALID_KEY_DATABLOCK_NUM:                                            "Txhdb_invalid_key_datablock_num",
	TXHDB_ERR_INVALID_VALUE_DATABLOCK_NUM:                                          "Txhdb_invalid_value_datablock_num",
	TXHDB_ERR_GET_RECORD_EXCEED_BUFF_LENGTH:                                        "Txhdb_get_record_exceed_buff_length",
	TXHDB_ERR_COMPRESSION_FAIL:                                                     "Txhdb_compession_fail",
	TXHDB_ERR_DECOMPRESSION_FAIL:                                                   "Txhdb_decompression_fail",
	TXHDB_ERR_INVALID_OFFSETINEXTDATA_AND_SIZE_WHEN_UPDATING_EXTDATA:
		"Txhdb_invalid_offsetInExtdata_and_size_when_updating_extdata",
	TXHDB_ERR_UNEXPECTED_FREEBLOCK:                                                 "Txhdb_unexpected_freeblock",
	TXHDB_ERR_VALUE_APOW_LESSER_THAN_KEY_APOW:
		"Txhdb_value_apow_lesser_than_key_apow__value_apow_should_be_equal_to_or_greater_than_key_apow",
	TXHDB_ERR_DUPLICATED_FILE_PATH:                                                 "Txhdb_duplicated_file_path",
	TXHDB_ERR_INVALID_KEY_HEAD_SIZE_IN_TXHDB_META:
		"Txhdb_invalid_key_head_size_in_txhdb_meta",
	TXHDB_ERR_INVALID_FILE_SIZE:                                                    "Txhdb_invalid_file_size",
	TXHDB_ERR_INVALID_FREE_BLOCK_SIZE:                                              "Txhdb_invalid_free_block_size",
	TXHDB_ERR_MMAP_MEMSIZE_CHANGE_NOT_PERMITTED:
		"Txhdb_mmap_memsize_change_not_permitted",
	TXHDB_ERR_NEW_FILE_OBJ_FAIL:                                                    "Txhdb_new_file_obj_fail",
	TXHDB_ERR_RECORD_KEY_OFFSET_LESSER_THAN_TXHDB_WHOLE_REC_OFFSET:
		"txhdb_record_key_offset_lesser_than_txhdb_whole_rec_offset",
	TXHDB_ERR_RECORD_VALUE_OFFSET_LESSER_THAN_TXHDB_WHOLE_REC_OFFSET:
		"txhdb_record_value_offset_lesser_than_txhdb_whole_rec_offset",
	TXHDB_ERR_RECORD_OFFSET_LESSER_THAN_TXHDB_WHOLE_REC_OFFSET:
		"txhdb_record_offset_lesser_than_txhdb_whole_rec_offset",
	TXHDB_ERR_KEY_BUFFSIZE_LESSER_THAN_KEY_HEADSIZE:
		"txhdb_key_buffsize_lesser_than_key_headsize",
	TXHDB_ERR_VALUE_BUFFSIZE_LESSER_THAN_VALUE_HEADSIZE:
		"txhdb_value_buffsize_lesser_than_value_headsize",
	TXHDB_ERR_RECORD_SIZE_LESSER_THAN_KEY_HEADSIZE:
		"txhdb_record_size_lesser_than_key_headsize",
	TXHDB_ERR_INVALID_BLOCK_MAGIC:                                                  "txhdb_invalid_block_magic",
	TXHDB_ERR_INVALID_FREE_BLOCK_MAGIC:                                             "txhdb_invalid_free_block_magic",
	TXHDB_ERR_INVALID_KEYMAGIC:
		"txhdb_invalid_KEYMAGIC_it_should_be_KEYMAGIC",
	TXHDB_ERR_INVALID_KEYSPLMAGIC:
		"txhdb_invalid_KEYSPLMAGIC_it_should_be_KEYSPLMAGIC",
	TXHDB_ERR_INVALID_VALMAGIC:
		"txhdb_invalid_VALMAGIC_it_should_be_VALMAGIC",
	TXHDB_ERR_INVALID_VALSPLMAGIC:
		"txhdb_invalid_VALSPLMAGIC_it_should_be_VALSPLMAGIC",
	TXHDB_ERR_UNSUPPORTED_KEY_FORMAT_VERSION:
		"txhdb_unsupported_key_format_version",
	TXHDB_ERR_UNSUPPORTED_KEY_SPLBLOCK_FORMAT_VERSION:
		"txhdb_unsupported_key_splblock_format_version",
	TXHDB_ERR_UNSUPPORTED_VALUE_FORMAT_VERSION:
		"txhdb_unsupported_value_format_version",
	TXHDB_ERR_UNSUPPORTED_VALUE_SPLBLOCK_FORMAT_VERSION:
		"txhdb_unsupported_value_splblock_format_version",
	TXHDB_ERR_UNSUPPORTED_FREE_BLOCK_FORMAT_VERSION:
		"txhdb_unsupported_value_splblock_format_version",
	TXHDB_ERR_KEY_HEAD_CRC_UNMATCH:                                                 "txhdb_key_head_crc_unmatch",
	TXHDB_ERR_KEY_SPLBLOCK_HEAD_CRC_UNMATCH:
		"txhdb_key_splblock_head_crc_unmatch",
	TXHDB_ERR_VALUE_HEAD_CRC_UNMATCH:                                               "txhdb_value_head_crc_unmatch",
	TXHDB_ERR_VALUE_SPLBLOCK_HEAD_CRC_UNMATCH:
		"txhdb_value_splblock_head_crc_unmatch",
	TXHDB_ERR_FREE_BLOCK_HEAD_CRC_UNMATCH:                                          "txhdb_free_block_head_crc_unmatch",
	TXHDB_ERR_FREE_BLOCK_LIST_INFO_CRC_UNMATCH:
		"txhdb_free_block_list_info_crc_unmatch",
	TXHDB_ERR_GET_KEY_READ_BUFFER_FAIL:                                             "txhdb_get_key_read_buffer_fail",
	TXHDB_ERR_GET_VALUE_READ_BUFFER_FAIL:                                           "txhdb_get_value_read_buffer_fail",
	TXHDB_ERR_GET_LRU_VALUE_BUFFER_FAIL:                                            "txhdb_get_lru_value_buffer_fail",
	TXHDB_ERR_GET_EXTDATA_READ_BUFFER_FAIL:                                         "txhdb_get_extdata_read_buffer_fail",
	TXHDB_ERR_KEY_BLOCK_BODYSIZE_GREATER_THAN_KEY_BODYSIZE:
		"txhdb_key_block_bodysize_greater_than_key_bodysize",
	TXHDB_ERR_VALUE_BLOCK_BODYSIZE_GREATER_THAN_VALUE_BODYSIZE:
		"txhdb_value_block_bodysize_greater_than_value_bodysize",
	TXHDB_ERR_NULL_RECORD_POINTER:                                                  "txhdb_null_record_pointer",
	TXHDB_ERR_NULL_RECORD_WRITE_BUFF:                                               "txhdb_null_record_write_buff",
	TXHDB_ERR_SERIALIZE_RECORD_KEY_HEAD:                                            "txhdb_serialize_record_key_head",
	TXHDB_ERR_INVALID_IDX_IN_STAT_NUMS_ARRAY:
		"txhdb_invalid_idx_in_stat_nums_array",
	TXHDB_ERR_INVALID_ELEMNUM_OF_STAT_KEYNUMS:
		"txhdb_invalid_elemnum_of_stat_keynums",
	TXHDB_ERR_INVALID_ELEMNUM_OF_STAT_VALNUMS:
		"txhdb_invalid_elemnum_of_stat_valnums",
	TXHDB_ERR_PRINT_SPACE_NOT_ENOUGH:                                               "txhdb_print_space_not_enough",
	TXHDB_ERR_LRU_SHIFTIN_NOT_ENOUGH_MEMORY:                                        "txhdb_lru_shiftin_not_enough_memory",
	TXHDB_ERR_LRU_SHIFTIN_NO_MORE_LRU_NODE:                                         "txhdb_lru_shiftin_no_more_lru_node",
	TXHDB_ERR_LRU_ADJUST_NO_MORE_LRU_NODE:                                          "txhdb_lru_adjust_no_more_lru_node",
	TXHDB_ERR_LRU_SHIFTOUT_RECORD_ALREADY_OUTSIDE_OF_MEMORY:
		"txhdb_lru_shiftout_record_already_outside_of_memory",
	TXHDB_ERR_FILE_EXTDATA_LENGTH_CRC_UNMATCH:
		"txhdb_file_extdata_length_crc_unmatch",
	TXHDB_ERR_FILE_EXTDATA_INVALID_LENGTH:                                          "txhdb_file_extdata_invalid_length",
	TXHDB_ERR_INVALID_VALUE_HEAD_SIZE_IN_TXHDB_META:
		"Txhdb_invalid_value_head_size_in_txhdb_meta",
	TXHDB_ERR_INVALID_SPLITDATABLOCK_HEAD_SIZE_IN_TXHDB_META:
		"Txhdb_invalid_splitdatablock_head_size_in_txhdb_meta",
	TXHDB_ERR_KEY_BUCKETIDX_UNMATCH:                                                "txhdb_key_bucketidx_unmatch",
	TXHDB_ERR_FILE_WRITE_SIZE_INVALID:                                              "Txhdb_file_write_size_invalid",
	TXHDB_ERR_MODIFY_STAT_UNSUPPORTED_OPERATION_TYPE:
		"Txhdb_modify_stat_unsupported_operation_type",
	TXHDB_ERR_INVALID_EXTDATAMAGIC:
		"txhdb_invalid_EXTDATAMAGIC_it_should_be_EXTDATAMAGIC",
	TXHDB_ERR_INVALID_INTERNAL_LIST_TAIL_DURING_POP_LRU_NODELIST:
		"txhdb_invalid_internal_list_tail_during_pop_lru_nodelist",
	TXHDB_ERR_GET_LRUNODE_FAIL:                                                     "txhdb_get_lrunode_fail",
	TXHDB_ERR_LRUNODE_INVALID_FLAG:                                                 "txhdb_lrunode_invalid_flag",
	TXHDB_ERR_INVALID_FREE_BLOCK_NUM_TOO_MANY_FREE_BLOCKS:
		"txhdb_invalid_free_block_num_too_many_free_blocks",
	TXHDB_ERR_INVALID_ELEMNUM_OF_STAT_NOPADDING_SIZE_KEYNUMS:
		"txhdb_invalid_elemnum_of_stat_nopadding_size_keynums",
	TXHDB_ERR_INVALID_ELEMNUM_OF_STAT_NOPADDING_SIZE_VALNUMS:
		"txhdb_invalid_elemnum_of_stat_nopadding_size_valnums",
	TXHDB_ERR_ADD_LSIZE_EXCEEDS_MAX_TSD_VALUE_BUFF_SIZE:
		"txhdb_add_lsize_exceeds_max_tsd_value_buff_size",
	TXHDB_ERR_INTERNAL_CONSTANTS_ILLEGAL:                                           "txhdb_internal_constants_illegal",
	TXHDB_ERR_TOO_BIG_KEY_BIZ_SIZE:                                                 "txhdb_too_big_key_biz_size",
	TXHDB_ERR_TOO_BIG_VALUE_BIZ_SIZE:                                               "txhdb_too_big_value_biz_size",
	TXHDB_ERR_INDEX_NO_EXIST:                                                       "txhdb_index_no_exist",
	TXHDB_ERR_INVALID_FREE_BLOCK_BASESIZE:                                          "txhdb_invalid_free_block_basesize",
	TXHDB_ERR_CANNOT_CREATE_MMAPSHM_BECAUSE_SHM_ALREADY_EXISTED:
		"txhdb_cannot_create_mmapshm_because_shm_already_existed",
	TXHDB_ERR_INVALID_GENSHM_KEY:                                                   "txhdb_invalid_gen_shm_key",
	TXHDB_ERR_GENSHM_GET_FAIL:                                                      "txhdb_invalid_gen_shm_get_fail",
	TXHDB_ERR_GENSHM_CREATE_FAIL:                                                   "txhdb_invalid_gen_shm_create_fail",
	TXHDB_ERR_GENSHM_STAT_FAIL:                                                     "txhdb_invalid_gen_shm_stat_fail",
	TXHDB_ERR_GENSHM_DOES_NOT_EXIST:                                                "txhdb_genshm_does_not_exist",
	TXHDB_ERR_GENSHM_ATTACH_FAIL_BECAUSE_IT_IS_ALREADY_ATTACHED_BY_OTHER_PROCESSES:
		"txhdb_genshm_attach_fail_because_it_is_already_attached_by_other_processes",
	TXHDB_ERR_GENSHM_ATTACH_FAIL:                                                   "txhdb_genshm_attach_fail",
	TXHDB_ERR_FILE_INCONSISTANT_MSIZE:                                              "txhdb_inconsistant_msize",
	TXHDB_ERR_INVALID_TCAP_GENSHM_MAGIC:                                            "txhdb_invalid_tcap_genshm_magic",
	TXHDB_ERR_GENSHM_FIXED_HEAD_BUFFLEN_UNMATCH:
		"txhdb_genshm_fixed_head_bufflen_unmatch",
	TXHDB_ERR_GENSHM_INVALID_HEADLEN:                                               "txhdb_genshm_invalid_headlen",
	TXHDB_ERR_GENSHM_HEAD_CRC_UNMATCH:                                              "txhdb_genshm_head_crc_unmatch",
	TXHDB_ERR_GENSHM_HEAD_INVALID_VERSION:                                          "txhdb_genshm_head_invalid_version",
	TXHDB_ERR_GENSHM_INVALID_FILETYPE:                                              "txhdb_genshm_invalid_filetype",
	TXHDB_ERR_GET_IPV4ADDR_FAIL:                                                    "txhdb_get_ipv4addr_fail",
	TXHDB_ERR_NO_VALID_IPV4ADDR_EXISTS:                                             "txhdb_no_valid_ipv4addr_exists",
	TXHDB_ERR_TRANSFER_IPV4ADDR_FAIL:                                               "txhdb_transfer_ipv4addr_fail",
	TXHDB_ERR_FILE_EXCEEDS_LSIZE_LIMIT:                                             "txhdb_file_exceeds_lsize_limit",
	TXHDB_ERR_GENSHM_DETACH_FAIL:                                                   "txhdb_genshm_detach_fail",
	TXHDB_ERR_TXHDB_HEAD_PARAMETERS_ERROR:                                          "txhdb_genshm_detach_fail",
	TXHDB_ERR_TXHDB_HEAD_OLD_VERSION:                                               "txhdb_head_old_version",
	TXHDB_ERR_TXHDB_SHM_COREINFO_UNMATCH:                                           "txhdb_shm_coreinfo_unmatch",
	TXHDB_ERR_TXHDB_SHM_EXTDATA_UNMATCH:                                            "txhdb_shm_extdata_unmatch",
	TXHDB_ERR_TXHDB_EXTDATA_CHECK_ERROR:                                            "txhdb_extdata_check_error",
	TXHDB_ERR_CHUNK_BUFFS_CANNOT_BE_ALLOCED_IF_THEY_ARE_NOT_RELEASED:
		"txhdb_chunk_buffs_cannot_be_alloced_if_they_are_not_released",
	TXHDB_ERR_ALLOCATE_MEMORY_FAIL:                                                 "txhdb_allocate_memory_fail",
	TXHDB_ERR_INVALID_CHUNK_RW_MANNER:                                              "txhdb_invalid_chunk_rw_manner",
	TXHDB_ERR_FILE_PREAD_NOT_COMPLETE:                                              "txhdb_file_pread_not_complete",
	TXHDB_ERR_FILE_PWRITE_NOT_COMPLETE:                                             "txhdb_file_pwrite_not_complete",
	TXHDB_ERR_KEY_ONEBLOCK_BUT_NEXT_NOTNULL:
		"TXHDB_ERR_KEY_ONEBLOCK_BUT_NEXT_NOTNULL",
	TXHDB_ERR_VALUE_ONEBLOCK_BUT_NEXT_NOTNULL:
		"TXHDB_ERR_VALUE_ONEBLOCK_BUT_NEXT_NOTNULL",
	TXHDB_ERR_VARINT_FORMAT_ERROR:                                                  "TXHDB_ERR_VARINT_FORMAT_ERROR",
	TXHDB_ERR_TXSTAT_ERROR:                                                         "TXHDB_ERR_TXSTAT_ERROR",
	TXHDB_ERR_INVALID_VERSION:                                                      "invalid txhdb version",
	TXHDB_ERR_FREE_BLOCK_NOT_ENOUGH:                                                "TXHDB_ERR_FREE_BLOCK_NOT_ENOUGH",

	//ENGINE SYSTEM (module id 0x07) Error String defined below
	ENG_ERR_INVALID_ARGUMENTS:                        "engine_invalid_arguments",
	ENG_ERR_INVALID_MEMBER_VARIABLE_VALUE:            "engine_invalid_member_variable_value",
	ENG_ERR_NEW_TXHCURSOR_FAILED:                     "engine_new_txhcursor_failed",
	ENG_ERR_TXHCURSOR_KEY_BUFFER_LEGHTH_NOT_ENOUGH:   "engine_txhcursor_key_buffer_leghth_not_enough",
	ENG_ERR_TXHCURSOR_VALUE_BUFFER_LEGHTH_NOT_ENOUGH: "engine_txhcursor_value_buffer_leghth_not_enough",
	ENG_ERR_TXHDB_FILEPATH_NULL:                      "engine_txhdb_filepath_null",
	ENG_ERR_TCHDB_RELATED_ERROR:                      "engine_tchdb_related_error",
	ENG_ERR_NULL_CACHE:                               "engine_null_cache",
	ENG_ERR_ITER_FAIL_SYSTEM_RECORD:                  "engine_interation_fail_system_record",
	ENG_ERR_SYSTEM_ERROR:                             "engine_system_error",
	ENG_ERR_ENGINE_ERROR:                             "engine_engine_error",
	ENG_ERR_DATA_ERROR:                               "engine_data_error",
	ENG_ERR_VERSION_ERROR:                            "engine_version_error",
	ENG_ERR_SYSTEM_ERROR_BUFF_OVERFLOW:               "engine_system_error_buff_overflow",
	ENG_ERR_METADATA_ERROR:                           "engine_metadata_error",
	ENG_ERR_ADD_KEYMETA_FAILED:                       "engine_add_keymate_failed",

	ENG_ERR_ADD_VALUEMETA_FAILED:         "engine_add_valuemeta_failed",
	ENG_ERR_RESERVED_FIELDNAME:           "engine_reserved_fieldname_error",
	ENG_ERR_KEYNAME_REPEAT:               "engine_keyname_repeat_error",
	ENG_ERR_VALUENAME_REPEAT:             "engine_valuename_repeat_error",
	ENG_ERR_MISS_KEYMETA:                 "engine_misss_keymate_error",
	ENG_ERR_DELETE_KEYFIELD:              "engine_delete_keyfield_error",
	ENG_ERR_CHANGE_KEYCOUNT:              "engine_change_keycount_error",
	ENG_ERR_CHANGE_KEYTYPE:               "engine_change_keytype_error",
	ENG_ERR_CHANGE_KEYLENGTH:             "engine_change_keylength_error",
	ENG_ERR_CHANGE_VALUETYPE:             "engine_change_valuetype_error",
	ENG_ERR_CHANGE_VALUELENGTH:           "engine_change_valuelength_error",
	ENG_ERR_CHANGE_DEFAULTVALUE:          "engine_change_defaultvalue_error",
	ENG_ERR_EMPTY_FIELDNAME:              "engine_empty_fieldname_error",
	ENG_ERR_INVALID_TARGET_KEYFIELD:      "engine_invalid_target_keyfield_error",
	ENG_ERR_INVALID_TARGET_VALUEFIELD:    "engine_invalid_target_valuefield_error",
	ENG_ERR_INVALID_TABLE_TYPE:           "engine_invalid_table_type_error",
	ENG_ERR_CHANGE_TABLE_TYPE:            "engine_change_table_type_error",
	ENG_ERR_MISS_VALUEMETA:               "engine_miss_valuemeta_error",
	ENG_ERR_NOT_ENOUGH_BUFF_FOR_FILEPATH: "engine_not_enough_buff_for_filepath",
	ENG_ERR_ENGINE_FILE_NOT_FOUND:        "engine file or index file not found",

	//ULOG SYSTEM (module id 0x09) Error String defined below
	ULOG_ERR_INVALID_PARAMS: "Ulog_invalid parameters",

	//SYNCDB SYSTEM (module id 0x0b) Error String defined below
	SYNCDB_ERR_INVALID_PARAMS:                     "Syncdb_invalid parameters",
	SYNCDB_ERR_PAUSE_TO_SEND_FOR_SWITCH_CONNECTOR: "pause to send req for switch connector",
	SYNCDB_ERR_CONNECTOR_IS_NOT_CONNECTED:         "connector is not connected",

	//TCAPSVR SYSTEM (module id 0x0d) Error String defined below
	SVR_ERR_FAIL_ROUTE:                                    "tcapsvr_fail_route",
	SVR_ERR_FAIL_TIMEOUT:                                  "tcapsvr_fail_timeout",
	SVR_ERR_FAIL_SHORT_BUFF:                               "tcapsvr_fail_short_buf",
	SVR_ERR_FAIL_SYSTEM_BUSY:                              "tcapsvr_fail_system_busy",
	SVR_ERR_FAIL_RECORD_EXIST:                             "tcapsvr_fail_record_exist",
	SVR_ERR_FAIL_INVALID_FIELD_NAME:                       "tcapsvr_fail_invalid_field_name",
	SVR_ERR_FAIL_VALUE_OVER_MAX_LEN:                       "tcapsvr_fail_value_over_max_len",
	SVR_ERR_FAIL_INVALID_FIELD_TYPE:                       "tcapsvr_fail_invalid_field_type",
	SVR_ERR_FAIL_SYNC_WRITE:                               "tcapsvr_fail_sync_write",
	SVR_ERR_FAIL_WRITE_RECORD:                             "tcapsvr_fail_write_record",
	SVR_ERR_FAIL_DELETE_RECORD:                            "tcapsvr_fail_delete_record",
	SVR_ERR_FAIL_DATA_ENGINE:                              "tcapsvr_fail_data_engine",
	SVR_ERR_FAIL_RESULT_OVERFLOW:                          "tcapsvr_fail_result_overflow",
	SVR_ERR_FAIL_INVALID_OPERATION:                        "tcapsvr_fail_invalid_operation",
	SVR_ERR_FAIL_INVALID_SUBSCRIPT:                        "tcapsvr_fail_invalid_subscript",
	SVR_ERR_FAIL_INVALID_INDEX:                            "tcapsvr_fail_invalid_index",
	SVR_ERR_FAIL_OVER_MAXE_FIELD_NUM:                      "tcapsvr_fail_over_max_field_num",
	SVR_ERR_FAIL_MISS_KEY_FIELD:                           "tcapsvr_fail_miss_key_field",
	SVR_ERR_FAIL_NEED_SIGNUP:                              "tcapsvr_fail_need_signup",
	SVR_ERR_FAIL_CROSS_AUTH:                               "tcapsvr_fail_cross_auth",
	SVR_ERR_FAIL_SIGNUP_FAIL:                              "tcapsvr_fail_signup_fail",
	SVR_ERR_FAIL_SIGNUP_INVALID:                           "tcapsvr_fail_signup_invalid",
	SVR_ERR_FAIL_SIGNUP_INIT:                              "tcapsvr_fail_signup_init",
	SVR_ERR_FAIL_LIST_FULL:                                "tcapsvr_fail_list_full",
	SVR_ERR_FAIL_LOW_VERSION:                              "tcapsvr_fail_low_version",
	SVR_ERR_FAIL_HIGH_VERSION:                             "tcapsvr_fail_high_version",
	SVR_ERR_FAIL_INVALID_RESULT_FLAG:                      "tcapsvr_fail_invalid_result_flag",
	SVR_ERR_FAIL_PROXY_STOPPING:                           "tcapsvr_fail_proxy_stopping",
	SVR_ERR_FAIL_SVR_READONLY:                             "tcapsvr_fail_svr_readonly",
	SVR_ERR_FAIL_SVR_READONLY_BECAUSE_IN_SLAVE_MODE:       "tcapsvr_fail_svr_readonly_because_in_slave_mode",
	SVR_ERR_FAIL_INVALID_VERSION:                          "tcapsvr_fail_invalid_version",
	SVR_ERR_FAIL_SYSTEM_ERROR:                             "tcapsvr_fail_system_error",
	SVR_ERR_FAIL_OVERLOAD:                                 "server is overload",
	SVR_ERR_FAIL_NOT_ENOUGH_DADADISK_SPACE:                "tcapsvr_fail_not_enough_datadisk_space",
	SVR_ERR_FAIL_NOT_ENOUGH_ULOGDISK_SPACE:                "tcapsvr_fail_not_enough_ulogdisk_space",
	SVR_ERR_FAIL_UNSUPPORTED_PROTOCOL_MAGIC:               "tcapsvr_fail_unsupported_protocol_magic",
	SVR_ERR_FAIL_UNSUPPORTED_PROTOCOL_CMD:                 "tcapsvr_fail_unsupported_protocol_cmd",
	SVR_ERR_FAIL_HIGH_TABLE_META_VERSION:                  "tcapsvr_fail_api_table_meta_version_too_high",
	SVR_ERR_FAIL_MERGE_VALUE_FIELD:                        "tcapsvr_fail_merge_value_field",
	SVR_ERR_FAIL_CUT_VALUE_FIELD:                          "tcapsvr_fail_cut_value_field",
	SVR_ERR_FAIL_PACK_FIELD:                               "tcapsvr_fail_pack_value_field",
	SVR_ERR_FAIL_UNPACK_FIELD:                             "tcapsvr_fail_unpack_value_field",
	SVR_ERR_FAIL_LOW_API_VERSION:                          "tcapsvr_fail_api_version_too_low",
	SVR_ERR_COMMAND_AND_TABLE_TYPE_IS_MISMATCH:            "the command in request is mismatch to the table type",
	SVR_ERR_FAIL_TO_FIND_CACHE:                            "tcapsvr_fail_to_find_cache",
	SVR_ERR_FAIL_TO_FIND_META:                             "tcapsvr_fail_to_find_meta",
	SVR_ERR_FAIL_TO_GET_CURSOR:                            "tcapsvr_fail_to_get_cursor",
	SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE:                    "field value gets out of the range specified by user",
	SVR_ERR_INVALID_ARGUMENTS:                             "tcapsvr_invalid_arguments",
	SVR_ERR_SLAVE_READ_INVALID:
		"ProcGetDuringMoveFromSrcReq failed because the svr is slave: can't read",
	SVR_ERR_NULL_CACHE:                                    "null cache object",
	SVR_ERR_NULL_CURSOR:                                   "null cursor object",
	SVR_ERR_METALIB_VERSION_LESS_THAN_ENTRY_VERSION:       "the metalib version in request is less than entry version",
	SVR_ERR_INVALID_SELECT_ID_FOR_UNION:                   "invalid select id for union",
	SVR_ERR_CAN_NOT_FIND_SELECT_ENTRY_FOR_UNION:           "can not find the select entry for union",
	SVR_ERR_FAIL_DOCUMENT_PACK_VERSION:                    "document pack version does not match",
	SVR_ERR_TCAPSVR_PROCESS_NOT_NORMAL:                    "tcapsvr process in abnormal",
	SVR_ERR_TBUSD_PROCESS_NOT_NORMAL:                      "tbusd process in abnormal",
	SVR_ERR_INVALID_ARRAY_COUNT:                           "array count invalid",
	SVR_ERR_REJECT_REQUEST_BECAUSE_ROUTE_IN_REJECT_STATUS:
		"reject request because route in reject status: it appears generally during data move",
	SVR_ERR_FAIL_GET_ROUTE_HASH_CODE:
		"get route hash code failed. it perhaps unpack ProxyHeadForReqSendToSvr failed.",
	SVR_ERR_FAIL_INVALID_FIELD_VALUE:                      "invalid field value",
	SVR_ERR_FAIL_PROTOBUF_FIELD_GET:                       "protobuf fail to get field",
	SVR_ERR_FAIL_PROTOBUF_VALUE_BUFF_EXCEED:
		"protobuf value buff exceed TCAPLUS_BIG_RECORD_MAX_VALUE_BUF_LEN(10M)",
	SVR_ERR_FAIL_PROTOBUF_FIELD_UPDATE:                    "protobuf fail to update field",
	SVR_ERR_FAIL_PROTOBUF_FIELD_INCREASE:                  "protobuf fail to increase field",
	SVR_ERR_FAIL_PROTOBUF_FIELD_TAG_MISMATCH:              "protobuf field tag mismatch",
	SVR_ERR_FAIL_BINLOG_SEQUENCE_TOO_SMALL:
		"binlog sequence too small for lossless move binlog sync: maybe the binlog file has already been deleted.",
	SVR_ERR_FAIL_SVR_IS_NOT_MASTER:
		"failed because svr is not master",
	SVR_ERR_FAIL_BINLOG_INVALID_FILE_PATH:                 "invalid binlog path",
	SVR_ERR_FAIL_BINLOG_SOCKET_SEND_BUFF_IS_FULL:
		"socket send buff is full for lossless mov binlog sync",
	SVR_ERR_FAIL_DOCUMENT_NOT_SUPPORT:                     "not support docment operation.",
	SVR_ERR_FAIL_PARTKEY_INSERT_NOT_SUPPORT:               "not support partkeyinsert operation.",
	SVR_ERR_FAIL_SQL_FILTER_FAILED:                        "sql filter failed.",
	SVR_ERR_FAIL_NOT_MATCHED_SQL_QUERY_CONDITION:          "not matched sql query condition.",

	//TCAPDB SYSTEM (module id 0x0f) Error String defined below
	TCAPDB_ERR_INVALID_PARAMS:               "Tcapdb_invalid parameters",
	TCAPDB_ERR_ALLOCATE_MEMORY_FAILED:       "tcapdb_allocate_memeory_failed",
	TCAPDB_ERR_INDEX_SERVER_RETURN_EXISTED:  "tcapdb_index_server_return_existed.",
	TCAPDB_ERR_INDEX_SERVER_RETURN_NOT_FIND: "tcapdb_index_server_return_not_find.",
	TCAPDB_ERR_INDEX_SERVER_RETURN_OVERLOAD: "tcapdb_index_server_return_overload: index server is overload.",
	TCAPDB_ERR_PACK_FAILED:                  "Tcapdb_pack_failed",
	TCAPDB_ERR_TIMEOUT:                      "Tcapdb_error_timeout",
	TCAPDB_ERR_REJECT_REQ:                   "Tcapdb_error_reject_request",

	//TCAPROXY SYSTEM (module id 0x11) Error String defined below
	PROXY_ERR_INVALID_PARAMS: "tcaproxy_invalid_parameters",
	PROXY_ERR_NO_NEED_ROUTE_BATCHGET_ACTION_MSG_WHEN_NODE_IS_IN_SYNC_STATUS:
		"tcaproxy_error_no_need_routes batchget_action_msg_when_node_is_in_sync_status",
	PROXY_ERR_NO_NEED_ROUTE_WHEN_NODE_IS_IN_REJECT_STATUS:
		"tcaproxy_error_no_need_routes_when_node_is_in_reject_status",
	PROXY_ERR_PROBE_TIMEOUT:         "tcaproxy_error_probe_timeout",
	PROXY_ERR_SYSTEM_ERROR:          "tcaproxy_error_system_error",
	PROXY_ERR_CONFIG_ERROR:          "tcaproxy_error_config_error",
	PROXY_ERR_OVER_MAX_NODE:         "tcaproxy_error_over_max_node",
	PROXY_ERR_INVALID_SPLIT_SIZE:    "tcaproxy_error_invalid_split_size",
	PROXY_ERR_INVALID_ROUTE_INDEX:   "tcaproxy_error_invalid_route_index",
	PROXY_ERR_CONNECT_SERVER:        "tcaproxy_error_connect_server",
	PROXY_ERR_COMPOSE_MSG:           "tcaproxy_error_compose_msg",
	PROXY_ERR_ROUTE_MSG:             "tcaproxy_error_route_msg",
	PROXY_ERR_SHORT_BUFFER:          "tcaproxy_error_short_buffer",
	PROXY_ERR_OVER_MAX_RECORD:       "tcaproxy_error_over_max_record",
	PROXY_ERR_INVALID_SERVICE_TABLE: "tcaproxy_error_invalid_service_table",
	PROXY_ERR_REGISTER_FAILED:       "tcaproxy_error_register_failed",
	PROXY_ERR_CREATE_SESSION_HASH:   "tcaproxy_error_create_session_hash",
	PROXY_ERR_WRONG_STATUS:          "tcaproxy_error_wrong_status",
	PROXY_ERR_UNPACK_MSG:            "tcaproxy_error_unpack_msg",
	PROXY_ERR_PACK_MSG:              "tcaproxy_error_pack_msg",
	PROXY_ERR_SEND_MSG:              "tcaproxy_error_send_msg",
	PROXY_ERR_ALLOCATE_MEMORY:       "tcaproxy_error_allocate_memory",
	PROXY_ERR_PARSE_MSG:             "tcaproxy_error_parse_msg",
	PROXY_ERR_INVALID_MSG:           "tcaproxy_error_invalid_msg",
	PROXY_ERR_FAILED_PROC_REQUEST_BECAUSE_NODE_IS_IN_SYNC_STASUS:
		"tcaproxy_error_failed_proc_request_becuase_node_is_in_sync_status",
	PROXY_ERR_KEY_FIELD_NUM_IS_ZERO:                                    "tcaproxy_error_key_field_num_is_zero",
	PROXY_ERR_LACK_OF_SOME_KEY_FIELDS:                                  "tcaproxy_error_lack_of_some_key_fields",
	PROXY_ERR_FAILED_TO_FIND_NODE:                                      "tcaproxy_error_failed_to_find_node",
	PROXY_ERR_INVALID_COMPRESS_TYPE:                                    "tcaproxy_error_invalid_compress_type",
	PROXY_ERR_REQUEST_OVERSPEED:                                        "tcaproxy_error_request_overspeed",
	PROXY_ERR_SWIFT_TIMEOUT:                                            "tcaproxy_error_swift_trans_timeout",
	PROXY_ERR_SWIFT_ERROR:                                              "tcaproxy_error_swift_other_errors",
	PROXY_ERR_DIRECT_RESPONSE:
		"tcaproxy_error_reponse_direct_not_processed_by_svr ",
	PROXY_ERR_INIT_TLOG:                                                "tcaproxy_error_init_tlog",
	PROXY_ERR_ASSISTANT_THREAD_NOT_RUN:                                 "tcaproxy_error_assistant_thread_not_run",
	PROXY_ERR_REQUEST_ACCESS_CTRL_REJECT:                               "tcaproxy_error_request_access_ctrl_reject",
	PROXY_ERR_NOT_ALL_NODES_ARE_IN_NORMAL_OR_WAIT_STATUS:
		"tcaproxy_error_routes_is_not_all_noraml_or_wait",
	PROXY_ERR_ALREADY_CACHED_REQUEST_TIMEOUT:                           "tcaproxy_error_already_cached_request_timeout",
	PROXY_ERR_FAILED_TO_CACHE_REQUEST:                                  "tcaproxy_error_failed_to_cache_request",
	PROXY_ERR_NOT_EXIST_CACHED_REQUEST:                                 "tcaproxy_error_not_exist_cached_request",
	PROXY_ERR_FAILED_NOT_ENOUGH_CACHE_BUFF:                             "tcaproxy_error_failed_not_enough_cache_buff",
	PROXY_ERR_FAILED_PROCESS_CACHED_REQUEST:                            "tcaproxy_error_failed_process_cached_request",
	PROXY_ERR_SYNC_ROUTE_HAS_BEEN_CANCELLED:                            "tcaproxy_sync_route_has_been_cancelled",
	PROXY_ERR_FAILED_LOCK_CACHE:                                        "tcaproxy_failed_lock_cache",
	PROXY_ERR_SWIFT_SEND_BUFFER_FULL:                                   "tcaproxy_swift_send_buffer_full",
	PROXY_ERR_REQUEST_OVERLOAD_CTRL_REJECT:                             "tcaproxy_error_request_overload_ctrl_reject",
	PROXY_ERR_SQL_QUERY_MGR_IS_NULL:                                    "tcaproxy_err_sql_query_mgr_is_null",
	PROXY_ERR_SQL_QUERY_INVALID_SQL_TYPE:                               "proxy_err_sql_query_invalid_sql_type",
	PROXY_ERR_GET_TRANSACTION_FAILED:                                   "proxy_err_get_transaction_failed",
	PROXY_ERR_ADD_TRANSACTION_FAILED:                                   "proxy_err_add_transaction_failed",
	PROXY_ERR_QUERY_FROM_INDEX_SERVER_FAILED:                           "proxy_err_query_from_index_server_failed",
	PROXY_ERR_QUERY_FROM_INDEX_SERVER_TIMEOUT:                          "proxy_err_query_from_index_server_timeout",
	PROXY_ERR_QUERY_FOR_CONVERT_TCAPLUS_REQ_TO_INDEX_SERVER_REQ_FAILED:
		"proxy_err_query_for_convert_tcaplus_req_to_index_server_req_failed",
	PROXY_ERR_QUERY_INDEX_FIELD_NOT_EXIST:                              "proxy_err_query_index_field_not_exist",
	PROXY_ERR_THIS_SQL_IS_NOT_SUPPORT:									"proxy_err_this_sql_is_not_support",

	//API SYSTEM (module id 0x13) Error Code defined below
	API_ERR_OVER_MAX_KEY_FIELD_NUM:                           "api_over_max_key_field_num_error",
	API_ERR_OVER_MAX_VALUE_FIELD_NUM:                         "api_over_max_value_field_num_error",
	API_ERR_OVER_MAX_FIELD_NAME_LEN:                          "api_over_max_field_name_len_error",
	API_ERR_OVER_MAX_FIELD_VALUE_LEN:                         "api_over_max_field_value_len_error",
	API_ERR_FIELD_NOT_EXSIST:                                 "api_field_not_exist_error",
	API_ERR_FIELD_TYPE_NOT_MATCH:                             "api_field_type_not_match_error",
	API_ERR_PARAMETER_INVALID:                                "api_parameter_invalid_error",
	API_ERR_OPERATION_TYPE_NOT_MATCH:                         "api_operation_type_not_match_error",
	API_ERR_PACK_MESSAGE:                                     "api_pack_message_error",
	API_ERR_UNPACK_MESSAGE:                                   "api_unpack_message_error",
	API_ERR_PACKAGE_NOT_UNPACKED:                             "api_package_not_unpacked_error",
	API_ERR_OVER_MAX_RECORD_NUM:                              "api_over_max_record_num_error",
	API_ERR_INVALID_COMMAND:                                  "api_invalid_command_error",
	API_ERR_NO_MORE_RECORD:                                   "api_no_more_record_error",
	API_ERR_OVER_KEY_FIELD_NUM:                               "api_over_key_field_num_error",
	API_ERR_OVER_VALUE_FIELD_NUM:                             "api_over_value_field_num_error",
	API_ERR_OBJ_NEED_INIT:                                    "api_obj_need_init_error",
	API_ERR_INVALID_DATA_SIZE:                                "api_invalid_data_size_error",
	API_ERR_INVALID_ARRAY_COUNT:                              "api_invalid_array_count_error",
	API_ERR_INVALID_UNION_SELECT:                             "api_invalid_union_select_error",
	API_ERR_MISS_PRIMARY_KEY:                                 "api_miss_primary_key_error",
	API_ERR_UNSUPPORT_FIELD_TYPE:                             "api_unsupport_field_type_error",
	API_ERR_ARRAY_BUFFER_IS_SMALL:                            "api_array_buffer_is_small_error",
	API_ERR_IS_NOT_WHOLE_PACKAGE:                             "api_is_not_whole_package_error",
	API_ERR_MISS_PAIR_FIELD:                                  "api_miss_pair_field_error",
	API_ERR_GET_META_ENTRY:                                   "api_get_meta_entry_error",
	API_ERR_GET_ARRAY_META:                                   "api_get_array_meta_error",
	API_ERR_GET_ENTRY_META:                                   "api_get_entry_meta_error",
	API_ERR_INCOMPATIBLE_META:                                "api_incompatible_meta_error",
	API_ERR_PACK_ARRAY_DATA:                                  "api_pack_array_data_error",
	API_ERR_PACK_UNION_DATA:                                  "api_pack_union_data_error",
	API_ERR_PACK_STRUCT_DATA:                                 "api_pack_struct_data_error",
	API_ERR_UNPACK_ARRAY_DATA:                                "api_unpack_array_data_error",
	API_ERR_UNPACK_UNION_DATA:                                "api_unpack_union_data_error",
	API_ERR_UNPACK_STRUCT_DATA:                               "api_unpack_struct_data_error",
	API_ERR_INVALID_INDEX_NAME:                               "api_invalid_index_name_error",
	API_ERR_MISS_PARTKEY_FIELD:                               "api_miss_partkey_field_error",
	API_ERR_ALLOCATE_MEMORY:                                  "api_allocate_memory_error",
	API_ERR_GET_META_SIZE:                                    "api_get_meta_size_error",
	API_ERR_MISS_BINARY_VERSION:                              "api_miss_binary_version_error",
	API_ERR_INVALID_INCREASE_FIELD:                           "api_invalid_increase_field_error",
	API_ERR_INVALID_RESULT_FLAG:                              "api_invalid_result_flag_error",
	API_ERR_OVER_MAX_LIST_INDEX_NUM:                          "api_over_max_list_index_num_error",
	API_ERR_INVALID_OBJ_STATUE:                               "api_invalid_obj_status_error",
	API_ERR_INVALID_REQUEST:                                  "api_invalid_request_error",
	API_ERR_INVALID_SHARD_LIST:                               "api_invalid_shard_list_error",
	API_ERR_TABLE_NAME_MISSING:                               "api_table_name_missing_error",
	API_ERR_SOCKET_SEND_BUFF_IS_FULL:                         "api_socket_send_buff_is_full_error",
	API_ERR_INVALID_MAGIC:                                    "api_invalid_magic",
	API_ERR_TABLE_IS_NOT_EXIST:                               "api_table_is_not_exist_error",
	API_ERR_SHORT_BUFF:                                       "api_buffer_not_enough",
	API_ERR_FLOW_CONTROL:                                     "api_flow_control",
	API_ERR_COMPRESS_SWITCH_NOT_SUPPORTED_REGARDING_THIS_CMD: "api_compress_switch_not_supported_regarding_this_cmd",
	API_ERR_FAILED_TO_FIND_ROUTE:
		"api_failed_to_find_route: perhaps the table is not register or all proxies are not connected",
	API_ERR_OVER_MAX_PKG_SIZE:                                "api_failed_over_max_pkg_size",
	API_ERR_INVALID_VERSION_FOR_TLV:
		"api_failed_invalid_version_for_tlv: the obj version is not equal to the lib version",
	API_ERR_BSON_SERIALIZE:                                   "cannot serailize the BSON object into a string",
	API_ERR_BSON_DESERIALIZE:                                 "cannot build a BSON object from the string",
	API_ERR_ADD_RECORD:                                       "failed to add a new record into request",
	API_ERR_ZONE_IS_NOT_EXIST:                                "zone_is_not_exist_error",
	API_ERR_TRAVERSER_IS_NOT_EXIST:                           "traverser_does_not_exist",
	API_ERR_INSTANCE_ID_FULL:                                 "instance_id_full",
	API_ERR_INSTANCE_INIT_LOG_FAILURE:                        "instance_fail_to_init_log",
	API_ERR_CONNECTOR_IS_ABNORMAL:                            "connector_is_abnoraml",
	API_ERR_WAIT_RSP_TIMEOUT:                                 "wait_rsp_timeout",

	//TCAPCENTER SYSTEM (module id 0x15) Error String defined below
	CENTER_ERR_INVALID_PARAMS:      "Tcapcenter_invalid parameters",
	CENTER_ERR_TABLE_ALREADY_EXIST: "Tcapcenter_table_already_exist",
	CENTER_ERR_TABLE_NOT_EXIST:     "Tcapcenter_table_not_exist",

	//TCAPDIR SYSTEM (module id 0x17) Error Code defined below
	DIR_ERR_SIGN_FAIL:                  "tcapdir_sign_fail",
	DIR_ERR_LOW_VERSION:                "tcapdir_low_version_error",
	DIR_ERR_HIGH_VERSION:               "tcapdir_high_version_error",
	DIR_ERR_GET_DIR_SERVER_LIST:        "tcapdir_get_dir_server_list_error",
	DIR_ERR_APP_IS_NOT_FOUNT:           "tcapdir_app_is_not_found_error",
	DIR_ERR_NOT_CONNECT_TCAPCENTER:     "tcapdir_is not conncted_tcapcenter_error",
	DIR_ERR_ZONE_IS_NOT_FOUNT:          "tcapdir_zone_is_not_found_error",
	DIR_ERR_HASH_TABLE_FAILED:          "tcapdir_hash_table_created_failed_error",
	DIR_ERR_GET_TABLE_AND_ACCESS_LIST:  "tcapdir_get_table_and_access_list_error",
	DIR_ERR_IS_NOT_THIS_ZONE_API:       "tcapdir_this_ip_is_not_this_zone_error",
	DIR_ERR_IS_NOT_IN_ZONES_WHITE_LIST: "tcapdir_this_ip_is_not_in_zones_white_list",

	//BSON ERROR(module id 0x1b) Error Code defined below
	BSON_ERR_TYPE_IS_NOT_MATCH:                          "the bson element type is not match.",
	BSON_ERR_INVALID_DATA_TYPE:                          "the bson element data type is invalid.",
	BSON_ERR_INVALID_VALUE:                              "the value of the bson element is invalid.",
	BSON_ERR_BSON_TYPE_UNMATCH_TDR_TYPE:                 "the bson data type is not match tdr type.",
	BSON_ERR_BSON_TYPE_IS_NOT_SUPPORT_BY_TCAPLUS:        "the bson data type is not support by tcaplus.",
	BSON_ERR_BSON_ARRAY_COUNT_IS_INVALID:
		"the bson array count is invalid: perhaps it is greater than max count.",
	BSON_ERR_FAILED_TO_PARSE:                            "parse the bson string failed.",
	BSON_ERR_INVALID_FIELD_NAME_LENGTH:
		"the field name length is invalid: perhaps it is greater than the max field name length.",
	BSON_ERR_INDEX_FIELD_NAME_NOT_EXIST_WITH_ARRAY_TYPE:
		"the index field name is not exist: but the array field name and index field name must be a pair.",
	BSON_ERR_INVALID_ARRAY_INDEX:                        "the index of the array is invalid.",
	BSON_ERR_TDR_META_LIB_IS_NULL:                       "the meta lib is null.",
	BSON_ERR_MATCHED_COUNT_GREATER_THAN_ONE:
		"the matched count is greater than one in elemMatch include \"$$uniq\" field name and primary key field name.",
	BSON_ERR_NO_MATCHED:                                 "there is no matched element according to $elemMatch.",
	BSON_ERR_GREATER_THAN_ARRAY_MAX_COUNT:               "the array real count is greater than the array max count.",
	BSON_ERR_BSON_EXCEPTION:                             "An exception occurred in bson lib.",
	BSON_ERR_STD_EXCEPTION:                              "std::exception occured.",
	BSON_ERR_INVALID_KEY:                                "bson_err_invalid_key",
	BSON_ERR_TDR_META_LIB_IS_INVALID:                    "bson_err_tdr_meta_lib_is_invalid",

	//TCAPTCAPCOMMON SYSTEM (module id 0x19) Error Code defined below
	COMMON_ERR_INVALID_ARGUMENTS:                "common_invalid_arguments",
	COMMON_ERR_INVALID_MEMBER_VARIABLE_VALUE:    "common_invalid_member_variable_value",
	COMMON_ERR_SPINLOCK_INIT_FAIL:               "common_spinlock_init_fail",
	COMMON_ERR_SPINLOCK_DESTROY_FAIL:            "common_spinlock_destroy_fail",
	COMMON_ERR_COMPRESS_BUF_NOT_ENOUGH:          "common_compress_buf_is_not_enough",
	COMMON_ERR_DECOMPRESS_BUF_NOT_ENOUGH:        "common_decompress_buf_is_not_enough",
	COMMON_ERR_DECOMPRESS_INVALID_INPUT:         "when decompress the input is invalid",
	COMMON_ERR_CANNOT_FIND_COMPRESS_ALGORITHM:   "can't find the compress algorithm.",
	COMMON_ERR_CANNOT_FIND_DECOMPRESS_ALGORITHM: "can't find the decompress algorithm.",
	COMMON_ERR_COMPRESS_FAIL:                    "compress failed.",
	COMMON_ERR_DECOMPRESS_FAIL:                  "decompress failed.",
	COMMON_ERR_INVALID_SWITCH_VALUE:             "invalid_switch_value.",
	COMMON_ERR_LINUX_SYSTEM_CALL_FAIL:           "linux system call failed: such as fopen: fget: sscanf and so on.",
	COMMON_ERR_NOT_FIND_STAT_CACHE_VALUE:
		"can not find old stat cache value: such as cpu:network old stat cache info.",
	COMMON_ERR_LZO_CHECK_FAIL:
		"when use lzo to compress file: it's header contains magic",

	//tcaplus_index SYSTEM (module id 0x1a) Error Code defined below
	TCAPLUS_INDEX_ERR_INVALID_PARAMS:                                         "tcaplus_index_invalid_parameters",
	TCAPLUS_INDEX_ERR_ALLOCATE_MEMORY:                                        "tcaplus_index_allocate_memory_failed",
	TCAPLUS_INDEX_ERR_CREATE_CONNECTOR_TO_INDEX_SERVER_FAILED:
		"tcaplus_index_create_connector_to_index_server_failed",
	TCAPLUS_INDEX_ERR_SEND_TO_INDEX_SERVER_FAILED_FOR_NO_CONNECTOR:
		"tcaplus_index_send_to_index_server_failed_for_no_connector",
	TCAPLUS_INDEX_ERR_SEND_TO_INDEX_SERVER_FAILED_FOR_NO_AVAILABLE_CONNECTOR:
		"tcaplus_index_send_to_index_server_failed_for_no_available_connector",
	TCAPLUS_INDEX_ERR_SEND_TO_INDEX_SERVER_FAILED_FOR_OTHER_REASON:
		"tcaplus_index_send_to_index_server_failed_for_other_reason",
	TCAPLUS_INDEX_ERR_PAUSE_SEND_FOR_CHANGIN_URL_STATUS:
		"tcaplus_index_pause_send_for_changing_url_status",
	TCAPLUS_INDEX_ERR_QUERY_INDEX_SERVER_OVERLOAD:
		"tcaplus_index_err_query_index_server_overload",

	// Non-error (for information purpose)
	COMMON_INFO_DATA_NOT_MODIFIED:
		"TCAPLUS_FLAG_FETCH_ONLY_IF_MODIFIED flag set and version equals: return early without real data",
}

type ErrorCode struct {
	Code    int
	Message string
}

func (e ErrorCode) Error() string {
	if len(e.Message) != 0 {
		return "errCode: " + strconv.Itoa(e.Code) + ", errMsg: " + e.Message
	}
	return "errCode: " + strconv.Itoa(e.Code) + ", errMsg: " + ErrorCodes[e.Code]
}

func GetErrMsg(tcapluserr int) string {
	return ErrorCodes[tcapluserr]
}
