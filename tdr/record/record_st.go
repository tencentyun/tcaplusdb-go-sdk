package record

import (
	"bytes"
	"encoding/binary"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
	"git.code.oa.com/tsf4g/tdrcom"
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

type setFieldFunc func(name string, data interface{}) error
type getFieldFunc func(name string, data interface{}) error

type TdrTableSt interface {
	GetTDRDBFeilds() *tdrcom.TDRDBFeilds
	Init()
	Pack(cutVer uint32) ([]byte, error)
	Unpack(cutVer uint32, data []byte) error
	GetBaseVersion() uint32
	GetCurrentVersion() uint32
}

type TdrSt interface {
	Init()
	Pack(cutVer uint32) ([]byte, error)
	PackTo(cutVer uint32, w *tdrcom.Writer) error
	Unpack(cutVer uint32, data []byte) error
	UnpackFrom(cutVer uint32, r *tdrcom.Reader) error
}

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
		fullKeyList := strings.Split(data.GetTDRDBFeilds().PrimaryKey, ",")
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
	for i := 0; i < stType.Elem().NumField(); i++ {
		fieldTag := stType.Elem().Field(i).Tag.Get(TdrFieldTag)
		fieldName := stType.Elem().Field(i).Name
		fieldType := stType.Elem().Field(i).Type
		fieldValue := stValue.Elem().Field(i)
		if len(fieldTag) <= 0 {
			logger.ERR("data name %s has no tag", fieldName)
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "data struct has no tag"}
		}

		//key or value func
		var setFieldFunc setFieldFunc
		if _, exist := keyMap[fieldName]; exist {
			//设置key字段
			setFieldFunc = r.SetKey
		} else if r.ValueSet != nil {
			if r.Cmd == cmd.TcaplusApiGetReq {
				//Get请求不关注value的内容
				r.ValueMap[fieldTag] = []byte{}
				continue
			}
			//设置value字段
			setFieldFunc = r.SetValue
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
			if r.Cmd == cmd.TcaplusApiListGetReq {
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
			unionTag := stType.Elem().Field(i).Tag.Get(TdrSelectTag)
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
			vBuf := new(bytes.Buffer)
			if err := binary.Write(vBuf, binary.LittleEndian, version); err != nil {
				return err
			}
			value := vBuf.Bytes()
			value = append(value, stBuf...)
			if err := setFieldFunc(fieldTag, value); err != nil {
				return err
			}

		//slice字段
		case reflect.Slice:
			//获取数组的大小
			referTag := stType.Elem().Field(i).Tag.Get(TdrReferTag)
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
					unionTag := stType.Elem().Field(i).Tag.Get(TdrSelectTag)
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
	for i := 0; i < stType.Elem().NumField(); i++ {
		fieldTag := stType.Elem().Field(i).Tag.Get(TdrFieldTag)
		fieldName := stType.Elem().Field(i).Name
		fieldType := stType.Elem().Field(i).Type
		fieldValue := stValue.Elem().Field(i)
		if len(fieldTag) <= 0 {
			logger.ERR("data name %s has no tag", fieldName)
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "data struct has no tag"}
		}
		//key or value get func
		var getFieldFunc getFieldFunc
		var findMap map[string][]byte
		if _, exist := keyMap[fieldName]; exist {
			getFieldFunc = r.GetKey
			findMap = r.KeyMap
		} else {
			getFieldFunc = r.GetValue
			findMap = r.ValueMap
		}

		//field 不存在则不解析
		if _, exist := findMap[fieldTag]; !exist {
			logger.INFO("st field name %s tag %s not exist", fieldName, fieldTag)
			continue
		}

		//设置字段
		switch fieldType.Kind() {
		case reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8,
			reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
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
			version := int16(0)
			if err := binary.Read(bytes.NewReader(vData), binary.LittleEndian, &version); err != nil {
				logger.ERR("zone %d table %s key %s binary.Read err %s",
					r.ZoneId, r.TableName, fieldTag, err.Error())
				return err
			}

			newData := reflect.New(fieldType.Elem())
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
				tdrSt, ok := st.(TdrSt)
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
			fieldValue.Set(newData)

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
			version := int16(0)
			if err := binary.Read(bytes.NewReader(vData), binary.LittleEndian, &version); err != nil {
				logger.ERR("zone %d table %s key %s binary.Read err %s",
					r.ZoneId, r.TableName, fieldTag, err.Error())
				return err
			}

			//获取数组的大小
			referTag := stType.Elem().Field(i).Tag.Get(TdrReferTag)
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
		getFieldFunc = r.GetKey
	} else {
		getFieldFunc = r.GetValue
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
		getFieldFunc = r.GetKey
	} else {
		getFieldFunc = r.GetValue
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
