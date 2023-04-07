package record

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/metadata"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/idl"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/tcapdbproto"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"github.com/tencentyun/tsf4g/tdrcom"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

const (
	TdrFieldTag      = "tdr_field"
	TdrReferTag      = "tdr_refer"
	TdrSelectTag     = "tdr_select"
	TdrSliceMaxCount = "tdr_count"
)

// 设置字段值函数定义
type setFieldFunc func(name string, data interface{}) error

// 获取字段值函数定义
type getFieldFunc func(name string, data interface{}) error

// tdr结构体接口，用于打解包
type TdrTableSt interface {
	GetTDRDBFeilds() *tdrcom.TDRDBFeilds
	Init()
	Pack(cutVer uint32) ([]byte, error)
	Unpack(cutVer uint32, data []byte) error
	GetBaseVersion() uint32
	GetCurrentVersion() uint32
}

// tdr结构体接口，用于打解包
type TdrSt interface {
	Init()
	Pack(cutVer uint32) ([]byte, error)
	PackTo(cutVer uint32, w *tdrcom.Writer) error
	Unpack(cutVer uint32, data []byte) error
	UnpackFrom(cutVer uint32, r *tdrcom.Reader) error
}

// tdr结构体接口，用于打解包
type TdrUnion interface {
	Init(selector int64)
	Pack(cutVer uint32, selector int64) ([]byte, error)
	PackTo(cutVer uint32, w *tdrcom.Writer, selector int64) error
	Unpack(cutVer uint32, data []byte, selector int64) error
	UnpackFrom(cutVer uint32, r *tdrcom.Reader, selector int64) error
}

/**
	@brief  基于TDR描述设置record数据
	@param [IN] data  基于TDR描述record接口数据，tdr的xml通过工具生成的go结构体，包含的TdrTableSt接口的一系列方法
	@retval error     错误码
**/
func (r *Record) SetDataWithIndexAndField(data TdrTableSt, FieldNameList []string, IndexName string) error {
	var keyName string
	if "" == IndexName {
		keyName = data.GetTDRDBFeilds().PrimaryKey
		// 去掉空格
		keyName = strings.Replace(keyName, " ", "", -1)
	} else {
		var flag bool
		keyName, flag = data.GetTDRDBFeilds().Index2Column[IndexName]
		// 去掉空格
		keyName = strings.Replace(keyName, " ", "", -1)
		if false == flag {
			logger.ERR("IndexName PrimaryKey is invalid %s", IndexName)
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "IndexName is invalid"}
		}
	}
	keyList := strings.Split(keyName, ",")
	keyMap := make(map[string]bool)
	for _, v := range keyList {
		if len(v) > 0 {
			keyMap[v] = true
		}
	}
	fullkeyMap := make(map[string]bool)
	if "" == IndexName {
		fullkeyMap = keyMap
	} else {
		primaryKey := data.GetTDRDBFeilds().PrimaryKey
		// 去掉空格
		primaryKey = strings.Replace(primaryKey, " ", "", -1)
		fullKeyList := strings.Split(primaryKey, ",")
		for _, v := range fullKeyList {
			if len(v) > 0 {
				fullkeyMap[v] = true
			}
		}
	}
	var FieldNameMap map[string]bool = nil
	if 0 != len(FieldNameList) {
		FieldNameMap = make(map[string]bool)
		for _, v := range FieldNameList {
			if len(v) > 0 {
				FieldNameMap[v] = true
			}
		}
	}

	if len(keyMap) <= 0 {
		logger.ERR("GetTDRDBFeilds PrimaryKey is invalid %s", keyName)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "GetTDRDBFeilds PrimaryKey is invalid"}
	}

	stValue := reflect.ValueOf(data)
	stType := reflect.TypeOf(data)
	stKind := stType.Kind()
	if stKind != reflect.Ptr && stValue.Elem().Kind() != reflect.Struct {
		logger.ERR("data type invalid %s %v", keyName, stKind)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "data type must be *struct"}
	}
	logger.DEBUG("st type %v", stType)

	//遍历字段
	numField := stType.Elem().NumField()
	fieldFunc := stType.Elem().Field
	valFieldFunc := stValue.Elem().Field

	for i := 0; i < numField; i++ {
		field := fieldFunc(i)
		fieldTag := field.Tag.Get(TdrFieldTag)
		fieldName := field.Name
		fieldType := field.Type
		fieldValue := valFieldFunc(i)
		if len(fieldTag) <= 0 {
			logger.ERR("data name %s has no tag", fieldName)
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "data struct has no tag"}
		}

		//key or value func
		var setFieldFunc setFieldFunc
		if _, exist := keyMap[fieldName]; exist {
			//设置key字段
			setFieldFunc = r.setKey
		} else if r.ValueSet != nil {
			if r.Cmd == cmd.TcaplusApiGetReq {
				//Get请求不关注value的内容
				r.ValueMap[fieldTag] = []byte{}
				continue
			}
			//设置value字段
			setFieldFunc = r.setValue
		} else {
			if r.Cmd == cmd.TcaplusApiGetByPartkeyReq {
				if _, exist := fullkeyMap[fieldName]; false == exist {
					if nil == FieldNameMap {
						logger.DEBUG("nil field list, set value fieldTag :%s", fieldTag)
						r.ValueMap[fieldTag] = []byte{}
					} else {
						if _, exist := FieldNameMap[fieldTag]; true == exist {
							logger.DEBUG("set value fieldTag :%s", fieldTag)
							r.ValueMap[fieldTag] = []byte{}
						}
					}
				}
			}
			// 如果后面未设置feiledname，在这里先设置获取所有字段。
			if r.Cmd == cmd.TcaplusApiListGetReq || r.Cmd == cmd.TcaplusApiBatchGetReq {
				logger.DEBUG("set value fieldTag :%s", fieldTag)
				r.ValueMap[fieldTag] = []byte{}
			}
			//ValueSet nil不需要打包value
			continue
		}

		//设置字段
		switch fieldType.Kind() {
		case reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8,
			reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:

			if err := setFieldFunc(fieldTag, fieldValue.Interface()); err != nil {

				return err
			}

		//struct 二级字段
		case reflect.Ptr:
			if fieldType.Elem().Kind() != reflect.Struct {
				logger.ERR("field type invalid name %s tag %s value %v", fieldName, fieldTag, fieldValue)

				return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}
			}

			st := fieldValue.Interface()
			if st == nil {
				logger.DEBUG("field is nil name %s tag %s value %v", fieldName, fieldTag, fieldValue)
				continue
			}

			var stBuf []byte
			unionTag := field.Tag.Get(TdrSelectTag)
			//判断是否是union类型
			if len(unionTag) > 0 {
				//打包union
				tdrSelector, err := r.getUnionSelectForSetData(data, unionTag)
				if err != nil {
					logger.ERR("getUnionSelectForSetData err %s", err.Error())

					return err
				}

				tdrUnion, ok := st.(TdrUnion)
				if !ok {
					logger.ERR("field type invalid name %s tag %s value %v trans to tdrUnion failed",
						fieldName, fieldTag, fieldValue)

					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}
				}
				packBuf, err := tdrUnion.Pack(0, tdrSelector)
				if err != nil {
					logger.ERR("field type invalid name %s tag %s value %v tdrUnion pack failed",
						fieldName, fieldTag, fieldValue)

					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}
				}
				stBuf = packBuf
			} else {
				//打包struct
				tdrSt, ok := st.(TdrSt)
				if !ok {
					logger.ERR("field type invalid name %s tag %s value %v trans to tdrSt failed",
						fieldName, fieldTag, fieldValue)

					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}
				}
				packBuf, err := tdrSt.Pack(0)
				if err != nil {
					logger.ERR("field type invalid name %s tag %s value %v tdrSt pack failed %s",
						fieldName, fieldTag, fieldValue, err.Error())
					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}
				}
				stBuf = packBuf
			}

			//set data会埋入版本号
			version := int16(data.GetCurrentVersion())
			value := make([]byte, 2, 2)
			binary.LittleEndian.PutUint16(value, uint16(version))
			value = append(value, stBuf...)
			if err := setFieldFunc(fieldTag, value); err != nil {
				return err
			}

		//slice字段
		case reflect.Slice:
			//获取数组的大小
			referTag := field.Tag.Get(TdrReferTag)
			arrayCount := 0
			if len(referTag) > 0 {
				//exist refer
				count, err := r.getSliceCountForSetData(data, referTag)
				if err != nil {
					logger.ERR("field array get refer %s failed, name %s tag %s ", referTag, fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct array field refer invalid"}
				}
				if count > fieldValue.Len() {
					logger.ERR("field array get refer %s too large %d > %d, name %s tag %s ",
						referTag, count, fieldValue.Len(), fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid,
						Message: "struct array field refer too large"}
				}
				arrayCount = count
			} else {
				//没有refer，则认为取最大长度
				countTag := field.Tag.Get(TdrSliceMaxCount)
				if len(countTag) <= 0 {
					logger.ERR("field array tdr_count failed, name %s tag %s ", fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid,
						Message: "struct array field has no tdr_count"}
				}
				count, err := strconv.Atoi(countTag)
				if err != nil {
					logger.ERR("field array tdr_count %s invalid %s, name %s tag %s ",
						countTag, err.Error(), fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct array tdr_count invalid"}
				}

				if count > fieldValue.Len() {
					logger.ERR("field array get tdr_count %s too large %d, name %s tag %s ",
						referTag, count, fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid,
						Message: "struct array field tdr_count too large"}
				}
				arrayCount = count
			}

			//set data会埋入版本号
			version := int16(data.GetCurrentVersion())
			var vFieldBytes []byte
			//打包数组
			if arrayCount > 0 {
				//数组普通元素类型
				switch fieldType.Elem().Kind() {
				case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
					vFieldBytes = make([]byte, arrayCount*int(fieldType.Elem().Size())+2)
					binary.LittleEndian.PutUint16(vFieldBytes, uint16(version))
					p := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: fieldValue.Pointer(),
						Len: arrayCount * int(fieldType.Elem().Size()),
						Cap: arrayCount * int(fieldType.Elem().Size())}))
					copy(vFieldBytes[2:], p)
				//todo 暂时不支持string数组
				//数组类型是结构体
				case reflect.Ptr:
					vBuf := new(bytes.Buffer)
					if err := binary.Write(vBuf, binary.LittleEndian, version); err != nil {
						return err
					}
					unionTag := field.Tag.Get(TdrSelectTag)
					for j := 0; j < arrayCount; j++ {
						eleValue := fieldValue.Index(j)
						if eleValue.Elem().Kind() != reflect.Struct {
							logger.ERR("field array type %v invalid name %s tag %s ",
								fieldValue.Elem().Kind(), fieldName, fieldTag)
							return &terror.ErrorCode{Code: terror.ParameterInvalid,
								Message: "struct array field type invalid"}
						}

						st := eleValue.Interface()
						var stBuf []byte
						//判断是否是union类型
						if len(unionTag) > 0 {
							//打包union
							tdrSelector, err := r.getUnionSelectForSetData(data, unionTag)
							if err != nil {
								logger.ERR("getUnionSelectForSetData err %s", err.Error())
								return err
							}

							tdrUnion, ok := st.(TdrUnion)
							if !ok {
								logger.ERR("field type invalid name %s tag %s value %v trans to tdrUnion failed",
									fieldName, fieldTag, fieldValue)
								return &terror.ErrorCode{Code: terror.ParameterInvalid,
									Message: "struct field type invalid"}
							}
							packBuf, err := tdrUnion.Pack(0, tdrSelector)
							if err != nil {
								logger.ERR("field type invalid name %s tag %s value %v tdrUnion pack failed",
									fieldName, fieldTag, fieldValue)
								return &terror.ErrorCode{Code: terror.ParameterInvalid,
									Message: "struct field type invalid"}
							}
							stBuf = packBuf
						} else {
							//打包struct
							tdrSt, ok := st.(TdrSt)
							if !ok {
								logger.ERR("field type invalid name %s tag %s value %v trans to tdrSt failed",
									fieldName, fieldTag, fieldValue)
								return &terror.ErrorCode{Code: terror.ParameterInvalid,
									Message: "struct field type invalid"}
							}
							packBuf, err := tdrSt.Pack(0)
							if err != nil {
								logger.ERR("field type invalid name %s tag %s value %v tdrSt pack failed %s",
									fieldName, fieldTag, fieldValue, err.Error())
								return &terror.ErrorCode{Code: terror.ParameterInvalid,
									Message: "struct field type invalid"}
							}
							stBuf = packBuf
						}
						vBuf.Write(stBuf)
					}
					vFieldBytes = vBuf.Bytes()
				default:
					logger.ERR("field array type %v invalid name %s tag %s ",
						fieldValue.Elem().Kind(), fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct array field type invalid"}
				}
			} else {
				vFieldBytes = make([]byte, 2)
				binary.LittleEndian.PutUint16(vFieldBytes, uint16(version))
			}
			if err := setFieldFunc(fieldTag, vFieldBytes); err != nil {
				return err
			}

		default:
			logger.ERR("field type invalid name %s tag %s value %v", fieldName, fieldTag, fieldValue)
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}

		}
	}
	logger.DEBUG("SetData success")
	return nil
}
func (r *Record) SetData(data TdrTableSt) error {
	return r.SetDataWithIndexAndField(data, nil, "")
}

/**
	@brief  设置过滤条件
	@param [IN] query     过滤条件例如：fieldValue > 4
	@retval int     错误码
**/
func (r *Record) SetCondition(query string) int {
	if r.Condition == nil {
		logger.ERR("cmd(%d) not support condition", r.Cmd)
		return terror.GEN_ERR_INVALID_ARGUMENTS
	}
	if int64(len(query)) > tcaplus_protocol_cs.TCAPLUS_MAX_EXPR_TEXT_LEN {
		logger.ERR("SetCondition query length %d larger max %d", len(query),
			tcaplus_protocol_cs.TCAPLUS_MAX_EXPR_TEXT_LEN)
		return terror.GEN_ERR_INVALID_ARGUMENTS
	}
	*r.Condition = query
	return 0
}

/**
	@brief  设置附加操作
	@param [IN] operation     附加操作：PUSH gameids #[-1][$=123]
	@param [IN] operateOption 附加操作类型 0|1
	@retval int     错误码
**/
func (r *Record) SetOperation(operation string, operateOption int32) int {
	if r.Operation == nil {
		logger.ERR("cmd(%d) not support operation", r.Cmd)
		return terror.GEN_ERR_INVALID_ARGUMENTS
	}
	if int64(len(operation)) > tcaplus_protocol_cs.TCAPLUS_MAX_EXPR_TEXT_LEN {
		logger.ERR("invalid condition.size=%d", len(operation))
		return terror.GEN_ERR_INVALID_ARGUMENTS
	}

	*r.Operation = operation
	if r.OperateOption != nil {
		*r.OperateOption = operateOption
	}

	return 0
}

/**
	@brief  基于TDR描述读取record数据
	@param [IN] data     基于TDR描述record接口数据，tdr的xml通过工具生成的go结构体，包含的TdrTableSt接口的一系列方法
	@retval error     错误码
**/
func (r *Record) GetData(data TdrTableSt) error {
	logger.DEBUG("GetData start")
	data.Init()
	keyName := data.GetTDRDBFeilds().PrimaryKey
	// 去掉空格
	keyName = strings.Replace(keyName, " ", "", -1)
	keyList := strings.Split(keyName, ",")
	keyMap := make(map[string]bool)
	for _, v := range keyList {
		if len(v) > 0 {
			keyMap[v] = true
		}
	}

	if len(keyMap) <= 0 {
		logger.ERR("GetTDRDBFeilds PrimaryKey is invalid %s", keyName)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "GetTDRDBFeilds PrimaryKey is invalid"}
	}

	stValue := reflect.ValueOf(data)
	stType := reflect.TypeOf(data)
	stKind := stType.Kind()
	if stKind != reflect.Ptr && stValue.Elem().Kind() != reflect.Struct {
		logger.ERR("data type invalid %s %v", keyName, stKind)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "data type must be *struct"}
	}
	logger.DEBUG("st type %v", stType)

	//遍历字段
	numField := stType.Elem().NumField()
	fieldFunc := stType.Elem().Field
	valFieldFunc := stValue.Elem().Field

	for i := 0; i < numField; i++ {
		field := fieldFunc(i)
		fieldTag := field.Tag.Get(TdrFieldTag)
		fieldName := field.Name
		fieldType := field.Type
		fieldValue := valFieldFunc(i)
		if len(fieldTag) <= 0 {
			logger.ERR("data name %s has no tag", fieldName)
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "data struct has no tag"}
		}
		//key or value get func
		var getFieldFunc getFieldFunc
		var findMap map[string][]byte
		if _, exist := keyMap[fieldName]; exist {
			getFieldFunc = r.getKey
			findMap = r.KeyMap
		} else {
			getFieldFunc = r.getValue
			findMap = r.ValueMap
		}

		//field 不存在则不解析
		if _, exist := findMap[fieldTag]; !exist {
			logger.DEBUG("st field name %s tag %s not exist", fieldName, fieldTag)
			continue
		}

		//设置字段
		switch fieldType.Kind() {
		case reflect.Bool:
			var vData bool
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetBool(vData)
		case reflect.Int8:
			var vData int8
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetInt(int64(vData))
		case reflect.Int16:
			var vData int16
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetInt(int64(vData))
		case reflect.Int32:
			var vData int32
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetInt(int64(vData))
		case reflect.Int64:
			var vData int64
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetInt(vData)
		case reflect.Uint8:
			var vData uint8
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetUint(uint64(vData))
		case reflect.Uint16:
			var vData uint16
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetUint(uint64(vData))
		case reflect.Uint32:
			var vData uint32
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetUint(uint64(vData))
		case reflect.Uint64:
			var vData uint64
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetUint(uint64(vData))
		case reflect.Float32:
			var vData float32
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetFloat(float64(vData))
		case reflect.Float64:
			var vData float64
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				return err
			}
			fieldValue.SetFloat(vData)
		case reflect.String:
			vData := reflect.New(fieldType)
			if err := getFieldFunc(fieldTag, vData.Interface()); err != nil {
				logger.ERR("getFieldFunc %s failed err %s", fieldTag, err.Error())
				return err
			}
			fieldValue.Set(vData.Elem())

		//struct 二级字段
		case reflect.Ptr:
			if fieldType.Elem().Kind() != reflect.Struct {
				logger.ERR("field type invalid name %s tag %s value %v", fieldName, fieldTag, fieldValue)
				return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}
			}

			//get value
			var vData []byte
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				logger.ERR("getFieldFunc %s failed err %s", fieldTag, err.Error())
				return err
			}

			if len(vData) <= 2 {
				logger.DEBUG("getFieldFunc %s struct data len < 2", fieldTag)
				continue
			}

			//版本号2B
			version := *(*int16)(unsafe.Pointer(&vData[0]))

			unionTag := field.Tag.Get(TdrSelectTag)
			//判断是union类型
			if len(unionTag) > 0 {
				//解包union
				tdrSelector, err := r.getUnionSelectForGetData(data, unionTag, keyMap)
				if err != nil {
					logger.ERR("getUnionSelectForSetData err %s", err.Error())
					return err
				}

				tdrUnion, ok := fieldValue.Interface().(TdrUnion)
				if !ok {
					logger.ERR("field type invalid name %s tag %s value %v trans to tdrUnion failed",
						fieldName, fieldTag, fieldValue)
					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}
				}
				tdrUnion.Init(tdrSelector)
				if err := tdrUnion.Unpack(uint32(version), vData[2:], tdrSelector); err != nil {
					logger.ERR("zone %d table %s key %s Unpack tdrUnion err %s",
						r.ZoneId, r.TableName, fieldTag, err.Error())
					return err
				}
			} else {
				//解包struct
				tdrSt, ok := fieldValue.Interface().(TdrSt)
				if !ok {
					logger.ERR("field type invalid name %s tag %s value %v trans to tdrSt failed",
						fieldName, fieldTag, fieldValue)
					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}
				}
				tdrSt.Init()
				//st unpack
				if err := tdrSt.Unpack(uint32(version), vData[2:]); err != nil {
					logger.ERR("zone %d table %s key %s Unpack tdrSt err %s",
						r.ZoneId, r.TableName, fieldTag, err.Error())
					return err
				}
			}

		//slice 数组字段
		case reflect.Slice:
			//get value
			var vData []byte
			if err := getFieldFunc(fieldTag, &vData); err != nil {
				logger.ERR("getFieldFunc %s failed err %s", fieldTag, err.Error())
				return err
			}

			if len(vData) <= 2 {
				logger.DEBUG("getFieldFunc %s array data len < 2", fieldTag)
				continue
			}

			//版本号2B
			version := *(*int16)(unsafe.Pointer(&vData[0]))

			//获取数组的大小
			referTag := field.Tag.Get(TdrReferTag)
			arrayCount := 0
			if len(referTag) > 0 {
				//exist refer
				count, err := r.getSliceCountForGetData(data, referTag, keyMap)
				if err != nil {
					logger.ERR("field array get refer %s failed, name %s tag %s ", referTag, fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct array field refer invalid"}
				}
				arrayCount = count
			} else {
				//没有refer，则认为取最大长度
				countTag := stType.Elem().Field(i).Tag.Get(TdrSliceMaxCount)
				if len(countTag) <= 0 {
					logger.ERR("field array tdr_count failed, name %s tag %s ", fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid,
						Message: "struct array field has no tdr_count"}
				}
				count, err := strconv.Atoi(countTag)
				if err != nil {
					logger.ERR("field array tdr_count %s invalid %s, name %s tag %s ",
						countTag, err.Error(), fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct array tdr_count invalid"}
				}
				arrayCount = count
			}

			//len > 0
			if arrayCount > 0 {
				//普通元素类型
				switch fieldType.Elem().Kind() {
				case reflect.Bool:
					valueSlice := reflect.ValueOf(*(*[]bool)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Int8:
					valueSlice := reflect.ValueOf(*(*[]int8)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Uint8:
					valueSlice := reflect.ValueOf(*(*[]uint8)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Int16:
					valueSlice := reflect.ValueOf(*(*[]int16)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Uint16:
					valueSlice := reflect.ValueOf(*(*[]uint16)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Int32:
					valueSlice := reflect.ValueOf(*(*[]int32)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Uint32:
					valueSlice := reflect.ValueOf(*(*[]uint32)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Int64:
					valueSlice := reflect.ValueOf(*(*[]int64)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Uint64:
					valueSlice := reflect.ValueOf(*(*[]uint64)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Float32:
					valueSlice := reflect.ValueOf(*(*[]float32)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				case reflect.Float64:
					valueSlice := reflect.ValueOf(*(*[]float64)(unsafe.Pointer(&reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(&vData[2])), Len: arrayCount, Cap: arrayCount})))
					fieldValue.Set(valueSlice)
				//todo 暂时不支持string数组
				//类型是结构体
				case reflect.Ptr:
					if fieldType.Elem().Elem().Kind() != reflect.Struct {
						logger.ERR("field array type %v invalid name %s tag %s ",
							fieldType.Elem().Elem().Kind(), fieldName, fieldTag)
						return &terror.ErrorCode{Code: terror.ParameterInvalid,
							Message: "struct array field type invalid"}
					}

					vBuf := tdrcom.NewReader(vData[2:])
					valueSlice := reflect.MakeSlice(fieldType, arrayCount, arrayCount)
					for j := 0; j < arrayCount && vBuf.Len() > 0; j++ {
						//解包
						newData := reflect.New(valueSlice.Index(j).Type().Elem())
						st := newData.Interface()
						unionTag := stType.Elem().Field(i).Tag.Get(TdrSelectTag)
						//判断是union类型
						if len(unionTag) > 0 {
							//解包union
							tdrSelector, err := r.getUnionSelectForGetData(data, unionTag, keyMap)
							if err != nil {
								logger.ERR("getUnionSelectForSetData err %s", err.Error())
								return err
							}

							tdrUnion, ok := st.(TdrUnion)
							if !ok {
								logger.ERR("field type invalid name %s tag %s value %v trans to tdrUnion failed",
									fieldName, fieldTag, fieldValue)
								return &terror.ErrorCode{Code: terror.ParameterInvalid,
									Message: "struct field type invalid"}
							}
							tdrUnion.Init(tdrSelector)
							if err := tdrUnion.UnpackFrom(uint32(version), vBuf, tdrSelector); err != nil {
								logger.ERR("zone %d table %s key %s Unpack tdrUnion err %s",
									r.ZoneId, r.TableName, fieldTag, err.Error())
								return err
							}
						} else {
							//解包struct
							tdrSt, ok := st.(TdrSt)
							if !ok {
								logger.ERR("field type invalid name %s tag %s value %v trans to tdrSt failed",
									fieldName, fieldTag, fieldValue)
								return &terror.ErrorCode{Code: terror.ParameterInvalid,
									Message: "struct field type invalid"}
							}
							tdrSt.Init()
							//st unpack
							if err := tdrSt.UnpackFrom(uint32(version), vBuf); err != nil {
								logger.ERR("zone %d table %s key %s Unpack tdrSt err %s",
									r.ZoneId, r.TableName, fieldTag, err.Error())
								return err
							}
						}
						valueSlice.Index(j).Set(newData)
					}
					fieldValue.Set(valueSlice)
				default:
					logger.ERR("field array type %v invalid name %s tag %s ",
						fieldValue.Elem().Kind(), fieldName, fieldTag)
					return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct array field type invalid"}
				}
			}

		default:
			logger.ERR("field type invalid name %s tag %s value %v", fieldName, fieldTag, fieldValue)
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "struct field type invalid"}

		}
	}
	logger.DEBUG("GetData finish")
	return nil
}

func (r *Record) getUnionSelectForSetData(data TdrTableSt, selector string) (int64, error) {
	stValue := reflect.ValueOf(data)
	//遍历字段
	for i := 0; i < stValue.Elem().NumField(); i++ {
		fieldName := stValue.Elem().Type().Field(i).Name
		fieldValue := stValue.Elem().Field(i)
		if selector == fieldName {
			switch fieldValue.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return int64(fieldValue.Int()), nil
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return int64(fieldValue.Uint()), nil
			default:
				logger.ERR("select type invalid %v", fieldValue.Kind())
				return 0, &terror.ErrorCode{Code: terror.ParameterInvalid,
					Message: "getUnionSelectForSetData selector type invalid"}
			}
		}
	}
	return 0, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "getUnionSelectForSetData selector not find"}
}

func (r *Record) getSliceCountForSetData(data TdrTableSt, refer string) (int, error) {
	stValue := reflect.ValueOf(data)
	//遍历字段
	for i := 0; i < stValue.Elem().NumField(); i++ {
		fieldName := stValue.Elem().Type().Field(i).Name
		fieldValue := stValue.Elem().Field(i)
		if refer == fieldName {
			switch fieldValue.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return int(fieldValue.Int()), nil
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return int(fieldValue.Uint()), nil
			default:
				logger.ERR("refer type invalid %v", fieldValue.Kind())
				return 0, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "getArrayCount refer type invalid"}
			}
		}
	}
	return 0, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "getArrayCount refer not find"}
}

func (r *Record) getUnionSelectForGetData(data TdrTableSt, selector string, keyMap map[string]bool) (int64, error) {
	//key or value get func
	var getFieldFunc getFieldFunc
	if _, exist := keyMap[selector]; exist {
		getFieldFunc = r.getKey
	} else {
		getFieldFunc = r.getValue
	}
	stValue := reflect.ValueOf(data)
	stType := reflect.TypeOf(data)

	//遍历字段
	for i := 0; i < stValue.Elem().NumField(); i++ {
		fieldName := stValue.Elem().Type().Field(i).Name
		fieldTag := stType.Elem().Field(i).Tag.Get(TdrFieldTag)
		fieldType := stType.Elem().Field(i).Type
		fieldValue := stValue.Elem().Field(i)
		if selector == fieldName {
			switch fieldValue.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				vData := reflect.New(fieldType)
				if err := getFieldFunc(fieldTag, vData.Interface()); err != nil {
					logger.ERR("getField %s failed err %s", selector, err.Error())
					return 0, err
				}
				return int64(vData.Elem().Int()), nil
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				vData := reflect.New(fieldType)
				if err := getFieldFunc(fieldTag, vData.Interface()); err != nil {
					logger.ERR("getField %s failed err %s", selector, err.Error())
					return 0, err
				}
				return int64(vData.Elem().Uint()), nil
			default:
				return 0, &terror.ErrorCode{Code: terror.ParameterInvalid,
					Message: "getUnionSelectForGetData selector type invalid"}
			}
		}
	}
	return 0, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "getUnionSelectForGetData selector not find"}
}

func (r *Record) getSliceCountForGetData(data TdrTableSt, refer string, keyMap map[string]bool) (int, error) {
	//key or value get func
	var getFieldFunc getFieldFunc
	if _, exist := keyMap[refer]; exist {
		getFieldFunc = r.getKey
	} else {
		getFieldFunc = r.getValue
	}
	stValue := reflect.ValueOf(data)
	stType := reflect.TypeOf(data)

	//遍历字段
	for i := 0; i < stValue.Elem().NumField(); i++ {
		fieldName := stValue.Elem().Type().Field(i).Name
		fieldTag := stType.Elem().Field(i).Tag.Get(TdrFieldTag)
		fieldType := stType.Elem().Field(i).Type
		fieldValue := stValue.Elem().Field(i)
		if refer == fieldName {
			switch fieldValue.Kind() {
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				vData := reflect.New(fieldType)
				if err := getFieldFunc(fieldTag, vData.Interface()); err != nil {
					logger.ERR("getField %s failed err %s", refer, err.Error())
					return 0, err
				}
				return int(vData.Elem().Int()), nil
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				vData := reflect.New(fieldType)
				if err := getFieldFunc(fieldTag, vData.Interface()); err != nil {
					logger.ERR("getField %s failed err %s", refer, err.Error())
					return 0, err
				}
				return int(vData.Elem().Uint()), nil
			default:
				return 0, &terror.ErrorCode{Code: terror.ParameterInvalid,
					Message: "getSliceCountForGetData refer type invalid"}
			}
		}
	}
	return 0, &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "getSliceCountForGetData refer not find"}
}

func (r *Record) AddValueOperation(FieldName string, FieldBuff []byte, FieldLen uint32,
	operation uint32, lower_limit int64, upper_limit int64) error {
	if operation < cmd.TcaplusApiOpPlus || operation > cmd.TcaplusApiOpMinus {
		return &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "AddValueOperation invalid param"}
	}
	r.UpdFieldSet.FieldNum += 1
	//var Fields tcaplus_protocol_cs.TCaplusUpdField
	Fields := &tcaplus_protocol_cs.TCaplusUpdField{
		FieldName:      FieldName,
		FieldLen:       FieldLen,
		FieldBuff:      FieldBuff,
		FieldOperation: operation,
		LowerLimit:     lower_limit,
		UpperLimit:     upper_limit,
	}
	r.UpdFieldSet.Fields = append(r.UpdFieldSet.Fields, Fields)
	return nil
}

/** brief  自增自减字段操作
 *  param  [in] field_name         字段名称
 *  param  [in] incData            加减数值，和表中定义的字段类型保持一致
 *  param  [in] operation          操作类型，cmd.TcaplusApiOpPlus 加操作 cmd.TcaplusApiOpMinus 减操作
 *  param [IN] lower_limit         操作结果值下限，如果比这个值小，返回 TcapErrCode::SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE
									支持double类型，int64(math.Float64bits(3.1415))
 *  param [IN] upper_limit         操作结果值上限，如果比这个值大，返回 TcapErrCode::SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE
									支持double类型，int64(math.Float64bits(3.1415))
 *  note                           lower_limit == upper_limit 时，存储端不对操作结果进行范围检测
*/
func (r *Record) SetIncValue(fieldName string, incData interface{},
	operation uint32, lowerLimit int64, upperLimit int64) error {
	if operation < cmd.TcaplusApiOpPlus || operation > cmd.TcaplusApiOpMinus {
		return &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "AddValueOperation invalid param"}
	}

	//check type
	var value []byte
	LowerUpperLimitType := byte(7) //CAPLUS_RECORD_TYPE_DOUBLE=10 CAPLUS_RECORD_TYPE_INT64=7
	if incData != nil {
		switch t := incData.(type) {
		case int8:
			value = make([]byte, 1, 1)
			value[0] = byte(t)
			break
		case int16:
			value = make([]byte, 2, 2)
			binary.LittleEndian.PutUint16(value, uint16(t))
			break
		case int32:
			value = make([]byte, 4, 4)
			binary.LittleEndian.PutUint32(value, uint32(t))
			break
		case int64:
			value = make([]byte, 8, 8)
			binary.LittleEndian.PutUint64(value, uint64(t))
			break
		case uint8:
			value = make([]byte, 1, 1)
			value[0] = t
			break
		case uint16:
			value = make([]byte, 2, 2)
			binary.LittleEndian.PutUint16(value, t)
			break
		case uint32:
			value = make([]byte, 4, 4)
			binary.LittleEndian.PutUint32(value, t)
			break
		case uint64:
			value = make([]byte, 8, 8)
			binary.LittleEndian.PutUint64(value, t)
			break
		case float32:
			value = make([]byte, 4, 4)
			binary.LittleEndian.PutUint32(value, math.Float32bits(t))
			LowerUpperLimitType = 10
			break
		case float64:
			value = make([]byte, 8, 8)
			binary.LittleEndian.PutUint64(value, math.Float64bits(t))
			LowerUpperLimitType = 10
			break
		default:
			logger.ERR("value type not support %v", t)
			return &terror.ErrorCode{Code: terror.RecordValueTypeInvalid, Message: "value type not support increase"}
		}
	}

	r.UpdFieldSet.FieldNum += 1
	Fields := &tcaplus_protocol_cs.TCaplusUpdField{
		FieldName:           fieldName,
		FieldLen:            uint32(len(value)),
		FieldBuff:           value,
		FieldOperation:      operation,
		LowerLimit:          lowerLimit,
		UpperLimit:          upperLimit,
		LowerUpperLimitType: LowerUpperLimitType,
	}
	r.UpdFieldSet.Fields = append(r.UpdFieldSet.Fields, Fields)
	return nil
}

/**
  @brief  基于 PB Message 设置record数据
  @param [IN] data  PB Message
  @retval []byte 记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应
  @retval error     错误码
*/
func (r *Record) SetPBData(message proto.Message) ([]byte, error) {
	return r.setPBDataCommon(message, nil, nil)
}

/**
    @brief 设置部分value字段，专用于field操作，TcaplusApiPBFieldGetReq TcaplusApiPBFieldUpdateReq TcaplusApiPBFieldIncreaseReq
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] values []string 指定本次设置的 value 字段
    @retval []byte 由记录key字段编码生成，由于多条记录的响应记录是无序的，可以用这个值来匹配记录
    @retval error 错误码
**/
func (r *Record) SetPBFieldValues(message proto.Message, values []string) ([]byte, error) {
	return r.setPBDataCommon(message, nil, values)
}

/**
    @brief 设置部分key字段，专用于partkey操作，TcaplusApiGetByPartkeyReq
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] keys []string 指定本次设置的 key 字段
    @retval []byte 由记录key字段编码生成，由于多条记录的响应记录是无序的，可以用这个值来匹配记录
    @retval error 错误码
**/
func (r *Record) SetPBPartKeys(message proto.Message, keys []string) ([]byte, error) {
	return r.setPBDataCommon(message, keys, nil)
}

// 设置 protobuf 值。
// 当 keys 不为空时，说明为 partkey 操作
// 当 values 不为空时，说明为 field 操作
func (r *Record) setPBDataCommon(message proto.Message, keys, values []string) ([]byte, error) {
	// 检查 message 名与表名是否相符
	table := message.ProtoReflect().Descriptor().Name()
	if string(table) != r.TableName {
		errMsg := fmt.Sprintf("message name (%s) not table name (%s)", table, r.TableName)
		logger.ERR(errMsg)
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: errMsg}
	}

	// 检查是否从服务端获取到了这张表的元数据
	zoneTable := fmt.Sprintf("%d|%d|%s", r.AppId, r.ZoneId, r.TableName)
	msgDesGrp := metadata.GetMetaManager().GetTableDesGrp(zoneTable)
	if msgDesGrp == nil {
		errMsg := fmt.Sprintf("zone %d table %s is not in table map", r.ZoneId, r.TableName)
		logger.ERR(errMsg)
		return nil, &terror.ErrorCode{Code: terror.TableNotExist, Message: errMsg}
	}

	// 对比元数据，防止字段属性有修改但未更新db
	if !msgDesGrp.Checked {
		err := metadata.GetMetaManager().CompareMessageMeta(msgDesGrp.Desc, message.ProtoReflect().Descriptor())
		if err != nil {
			errMsg := fmt.Sprintf("CompareMessageMeta error:%s", err)
			logger.ERR(errMsg)
			return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: errMsg}
		}
		msgDesGrp.Checked = true
	}

	// 如果有shardingkey则设置
	if msgDesGrp.ShardingKey != "" {
		keys := strings.Split(msgDesGrp.ShardingKey, ",")
		shardingKey, _ := metadata.GetMetaManager().ExtractMsgPartKey(message, keys)
		if r.ShardingKey != nil {
			*r.ShardingKey = shardingKey
			*r.ShardingKeyLen = uint32(len(shardingKey))
		}
		if r.SplitTableKeyBuff != nil {
			r.SplitTableKeyBuff.SplitTableKeyBuff = shardingKey
			r.SplitTableKeyBuff.SplitTableKeyBuffLen = uint32(len(shardingKey))
		}
	}

	// 计算 key 值，区分 partkey 操作
	if len(keys) == 0 {
		keys = msgDesGrp.Keys
	}
	keybuf, err := metadata.GetMetaManager().ExtractMsgPartKey(message, keys)
	if err != nil {
		return nil, err
	}

	// 计算 value 值，区分 field 操作
	var buf []byte
	if len(values) == 0 {
		buf, err = proto.Marshal(message)
		if err != nil {
			errMsg := fmt.Sprintf("Marshal message %s error:%s", table, err)
			logger.ERR(errMsg)
			return nil, &terror.ErrorCode{Code: terror.API_ERR_PACK_MESSAGE, Message: errMsg}
		}
	} else {
		if r.PBFieldMap == nil {
			errMsg := fmt.Sprintf("request not field operation")
			logger.ERR(errMsg)
			return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: errMsg}
		}
		if r.Cmd == cmd.TcaplusApiPBFieldGetReq || r.Cmd == cmd.TcaplusApiPBFieldUpdateReq {
			for _, v := range values {
				r.ValueMap[v] = nil
				r.PBFieldMap[v] = true
			}
			r.ValueMap["$"], _ = proto.Marshal(message)
			r.PBFieldMap["$"] = true
		} else {
			fieldMap, err := r.CheckValues(values)
			if err != nil {
				errMsg := fmt.Sprintf("CheckValues error:%s", err)
				logger.ERR(errMsg)
				return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: errMsg}
			}
			r.setPBValues(message, fieldMap)
		}
	}

	data := &idl.Tbl_Idl{Key: keybuf, Klen: int32(len(keybuf)), Value: buf, Vlen: int32(len(buf))}
	return keybuf, r.SetData(data)
}

// 获取 shardingkey
func (r *Record) GetTableShardingKey(message proto.Message) []byte {
	zoneTable := fmt.Sprintf("%d|%d|%s", r.AppId, r.ZoneId, r.TableName)
	msgGrp := metadata.GetMetaManager().GetTableDesGrp(zoneTable)
	if msgGrp == nil {
		return nil
	}
	keys := strings.Split(msgGrp.ShardingKey, ",")
	shardingKey, _ := metadata.GetMetaManager().ExtractMsgPartKey(message, keys)
	return shardingKey
}

/**
    @brief  基于 PB Message 读取record数据
    @param [IN] data   PB Message
    @retval []byte 记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应
    @retval error      错误码
**/
func (r *Record) GetPBData(message proto.Message) error {
	if r.Cmd == cmd.TcaplusApiGetTtlRes || r.Cmd == cmd.TcaplusApiSetTtlRes {
		buf, err := r.getKeyBlob("key")
		if err != nil {
			return err
		}
		err = proto.Unmarshal(buf[2:], message)
		if err != nil {
			logger.ERR(err.Error())
			return err
		}
		return nil
	}
	return r.GetPBDataWithValues(message, nil)
}

// 专用于 field 方法
/**
    @brief 专用于 field 方法，获取响应
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (r *Record) GetPBFieldValues(message proto.Message) error {
	if r.PBValueSet == nil {
		errMsg := fmt.Sprintf("PBValueSet is nil")
		logger.ERR(errMsg)
		return &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: errMsg}
	}

	zoneTable := fmt.Sprintf("%d|%d|%s", r.AppId, r.ZoneId, r.TableName)
	msgDesGrp := metadata.GetMetaManager().GetTableDesGrp(zoneTable)
	if msgDesGrp == nil {
		errMsg := fmt.Sprintf("zone %d table %s is not in table map", r.ZoneId, r.TableName)
		logger.ERR(errMsg)
		return &terror.ErrorCode{Code: terror.TableNotExist, Message: errMsg}
	}

	buf, err := r.GetKeyBlob("key")
	if err != nil {
		return err
	}
	err = proto.Unmarshal(buf[2:], message)
	if err != nil {
		logger.ERR(err.Error())
		return &terror.ErrorCode{Code: terror.API_ERR_UNPACK_MESSAGE}
	}

	if r.Cmd == cmd.TcaplusApiPBFieldGetRes || r.Cmd == cmd.TcaplusApiPBFieldUpdateRes {
		err = proto.UnmarshalOptions{Merge: true}.Unmarshal(r.ValueMap["$"], message)
		if err != nil {
			logger.ERR(err.Error())
			return &terror.ErrorCode{Code: terror.API_ERR_UNPACK_MESSAGE}
		}
		return nil
	}

	for numPath, v := range r.ValueMap {
		tmp := message
		nums := msgDesGrp.NumberMap[numPath]
		for i, num := range nums {
			f := tmp.ProtoReflect().Descriptor().Fields().ByNumber(num)
			if i == len(nums)-1 {
				err := proto.UnmarshalOptions{Merge: true}.Unmarshal(v, tmp)
				if err != nil {
					logger.ERR(err.Error())
					return &terror.ErrorCode{Code: terror.API_ERR_UNPACK_MESSAGE}
				}
			} else {
				tmp = tmp.ProtoReflect().Mutable(f).Message().Interface()
			}
		}
	}

	return nil
}

// 获取指定value， values不为空时，将 values 以外的字段置空
func (r *Record) GetPBDataWithValues(message proto.Message, values []string) error {
	data := &idl.Tbl_Idl{}
	value, exist := r.ValueMap["value"]
	if !exist {
		return &terror.ErrorCode{Code: terror.API_ERR_UNPACK_MESSAGE, Message: "rsp not has pb value"}
	}
	_, exist = r.ValueMap["vlen"]
	if !exist {
		data.Value = value
	} else {
		data.Value = value[2:]
	}

	err := proto.Unmarshal(data.Value, message)
	if err != nil {
		logger.ERR(err.Error())
		return &terror.ErrorCode{Code: terror.API_ERR_UNPACK_MESSAGE}
	}

	if len(values) == 0 {
		return nil
	}

	fmap, err := r.CheckValues(values)
	if err != nil {
		errMsg := fmt.Sprintf("CheckValues error:%s", err)
		logger.ERR(errMsg)
		return &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: errMsg}
	}

	r.cleanField(message.ProtoReflect(), fmap, "")

	return nil
}

// 检查values是否都在message中
func (r *Record) CheckValues(values []string) (map[string][]protowire.Number, error) {
	fmap := make(map[string][]protowire.Number)

	zoneTable := fmt.Sprintf("%d|%d|%s", r.AppId, r.ZoneId, r.TableName)
	msgDesGrp := metadata.GetMetaManager().GetTableDesGrp(zoneTable)
	if msgDesGrp == nil {
		logger.ERR("zone %d table %s is not in table map", r.ZoneId, r.TableName)
		return nil, fmt.Errorf("zone %d table %s is not in table map", r.ZoneId, r.TableName)
	}

	for _, v := range values {
		nums, exist := msgDesGrp.FieldMap[v]
		if !exist {
			logger.ERR("field %s not in table %s", v, r.TableName)
			return nil, fmt.Errorf("field %s not in table %s", v, r.TableName)
		}
		fmap[v] = nums
	}
	return fmap, nil
}

// 设置需要的字段
func (r *Record) setPBValues(message proto.Message, fieldMap map[string][]protowire.Number) {
	for _, nums := range fieldMap {
		tmp := message.ProtoReflect()
		var keybuf []byte
		numPath := ""
		for i, num := range nums {
			f := tmp.Descriptor().Fields().ByNumber(num)
			if numPath == "" {
				numPath = fmt.Sprint(num)
			} else {
				numPath += "." + fmt.Sprint(num)
			}
			if i == len(nums)-1 {
				keybuf, _ = tcapdbproto.MarshalOptions{}.MarshalField(keybuf, f, tmp.Get(f))
			} else {
				tmp = tmp.Get(f).Message()
			}
		}
		r.ValueMap[numPath] = keybuf
		r.PBFieldMap[numPath] = true
	}
}

// 清除不需要的字段
func (r *Record) cleanField(pro protoreflect.Message, fmap map[string][]protowire.Number, prefixName string) {
	fields := pro.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		fname := prefixName
		if prefixName == "" {
			fname = string(f.Name())
		} else {
			fname += "." + string(f.Name())
		}
		if _, exist := fmap[fname]; exist {
			continue
		}
		if f.Message() != nil {
			r.cleanField(pro.Get(f).Message(), fmap, fname)
		} else {
			pro.Clear(f)
		}
	}
}

/**
    @brief 获取记录key编码值
    @retval []byte 由记录key字段编码生成，由于多条记录的响应记录是无序的，可以用这个值来匹配记录
    @retval error 错误码
**/
func (r *Record) GetPBKey(msg proto.Message) ([]byte, error) {
	buf, err := r.GetKeyBlob("key")
	if err != nil {
		logger.ERR(err.Error())
		return nil, err
	}
	if msg != nil {
		err = proto.Unmarshal(buf[2:], msg)
		if err != nil {
			logger.ERR(err.Error())
			return nil, &terror.ErrorCode{Code: terror.API_ERR_UNPACK_MESSAGE}
		}
	}
	return buf[2:], nil
}

/**
@brief  设置记录的生存时间，或者说过期时间，即记录多久之后过期，过期的记录将不会被访问到
@param [IN] ttl 生存时间（过期时间），时间单位为毫秒，如果是相对时间，比如该参数值为10，则表示记录写入10ms之后过期，该参数值为0，则表示记录永不过期
												   如果是绝对时间，比如该参数值为1599105600000, 则表示记录到"20200903 12:00:00"之后过期，该参数值为0，则表示记录永不过期
@param [IN] is_absolute_time 时间类型是否为绝对时间，true表示绝对时间，false表示相对时间，默认是false，即相对时间
@retval 0                       设置成功
@retval 非0                     设置失败，具体错误参见 \link ErrorCode \endlink
@note   该函数当前只支持 TCAPLUS_API_SET_TTL_REQ 响应
@note   设置的ttl值最大不能超过uint64_t最大值的一半，即ttl最大值为 ULONG_MAX/2，超过该值接口会强制设置为该值
@note   设置ttl的请求，在服务端不会增加对应记录的版本号
@note   对于list表，当某个key下面所有记录因为过期删除后，会直接将索引记录也删除
@note   对于设置了ttl的记录，如果是getbypartkey查询，并且只需要返回key字段（即不需要返回value字段）时，此时不会检查该记录是否过期
@note   对于删除操作(generic表和list表的删除)，均不会检验记录是否过期
*/
func (r *Record) SetTTL(ttl uint64, isAbsoluteTime bool) int {
	if r.Cmd != cmd.TcaplusApiSetTtlReq {
		logger.ERR("expect cmd is TCAPLUS_API_SET_TTL_REQ(%d)", cmd.TcaplusApiSetTtlReq)
		return terror.GEN_ERR_ERR
	}

	if r.Ttl == nil {
		logger.ERR("Record ttl is nil")
		return terror.GEN_ERR_ERR
	}

	if r.TtlType == nil {
		logger.ERR("Record ttl type is nil")
		return terror.GEN_ERR_ERR
	}

	//如果ttl的值大于最大值的一半时，将强制设置为最大值的一半
	if ttl > math.MaxUint64/2 {
		ttl = math.MaxUint64 / 2
	}

	*r.Ttl = ttl

	// 如果是绝对时间，则设置ttl_type值为1
	if isAbsoluteTime {
		*r.TtlType = uint8(tcaplus_protocol_cs.TYPE_ABSOLUTE_TIME)
	} else { //如果是相对时间，则设置ttl_type值为0
		*r.TtlType = uint8(tcaplus_protocol_cs.TYPE_RELATIVE_TIME)
	}

	return 0
}

/**
    @brief 专用于 getttl 方法，设置过期
    @param [IN] ttl 单位 ms
    @retval error 错误码
**/
func (r *Record) GetTTL(ttl *uint64) int {
	if r.Cmd != cmd.TcaplusApiGetTtlRes {
		logger.ERR("expect cmd is TCAPLUS_API_GET_TTL_RES(%d)", cmd.TcaplusApiGetTtlRes)
		return terror.GEN_ERR_ERR
	}

	if r.Ttl == nil {
		logger.ERR("Record ttl is nil")
		return terror.GEN_ERR_ERR
	}

	*ttl = *r.Ttl

	return 0
}
