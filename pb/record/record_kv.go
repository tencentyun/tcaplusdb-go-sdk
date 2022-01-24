package record

import (
	"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"math"
	"unsafe"
)

/**
	@brief  通用的key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，最大长度1024字节，必须明确数据类型，必须和tdr xml表中的类型一致
								支持bool, byte, int8, int16, uint16, int32, uint32, int64, uint64, float32, float64，string, []byte
	@notice		请根据xml表准确填写类型，最好调用SetKeyInt8等接口
*/
func (r *Record) SetKey(name string, data interface{}) error {
	if r.IsPB {
		return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH}
	}
	return r.setKey(name, data)
}
func (r *Record) setKey(name string, data interface{}) error {
	if len(name) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME) {
		logger.ERR("key name len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME)
		return &terror.ErrorCode{Code: terror.KeyNameLenOverMax}
	}

	//check type
	var value []byte
	switch t := data.(type) {
	case bool:
		value = make([]byte, 1, 1)
		if t {
			value[0] = 1
		} else {
			value[0] = 0
		}
		break
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
		break
	case float64:
		value = make([]byte, 8, 8)
		binary.LittleEndian.PutUint64(value, math.Float64bits(t))
		break
	case []byte:
		if b, ok := data.([]byte); !ok {
			logger.ERR("key type not []byte")
			return &terror.ErrorCode{Code: terror.RecordKeyTypeInvalid}
		} else {
			value = b
		}
		break
	case string:
		//+ "\0"
		if str, ok := data.(string); !ok {
			logger.ERR("key type not string")
			return &terror.ErrorCode{Code: terror.RecordKeyTypeInvalid}
		} else {
			value = common.StringToCByte(str)
		}
		break
	default:
		logger.ERR("key type not support %v", t)
		return &terror.ErrorCode{Code: terror.RecordKeyTypeInvalid}
	}

	if len(value) > int(tcaplus_protocol_cs.TCAPLUS_MAX_KEY_FIELD_LEN) {
		logger.ERR("key len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_KEY_FIELD_LEN)
		return &terror.ErrorCode{Code: terror.KeyLenOverMax}
	}

	if len(r.KeyMap) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_KEY_FIELD_NUM) {
		logger.ERR("key num over %d", tcaplus_protocol_cs.TCAPLUS_MAX_KEY_FIELD_NUM)
		return &terror.ErrorCode{Code: terror.KeyNumOverMax}
	}

	r.KeyMap[name] = value
	return nil
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int8
*/
func (r *Record) SetKeyInt8(name string, data int8) error {
	return r.SetKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int16
*/
func (r *Record) SetKeyInt16(name string, data int16) error {
	return r.SetKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int32
*/
func (r *Record) SetKeyInt32(name string, data int32) error {
	return r.SetKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int64
*/
func (r *Record) SetKeyInt64(name string, data int64) error {
	return r.SetKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，float32
*/
func (r *Record) SetKeyFloat32(name string, data float32) error {
	return r.SetKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，float64
*/
func (r *Record) SetKeyFloat64(name string, data float64) error {
	return r.SetKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，string
*/
func (r *Record) SetKeyStr(name string, data string) error {
	return r.SetKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，[]byte
*/
func (r *Record) SetKeyBlob(name string, data []byte) error {
	return r.SetKey(name, data)
}

/**
	@brief  通用的value字段内容设置
 	@param  [in] name         	字段名称，最大长度32
	@param  [in] data         	字段内容，最大长度1024字节，必须明确数据类型，必须和tdr xml表中的类型一致
								支持bool, byte, int8, int16, uint16, int32, uint32, int64, uint64, float32, float64，string, []byte
	@notice		请根据xml表准确填写类型，最好调用SetValueInt8等接口
*/
func (r *Record) SetValue(name string, data interface{}) error {
	if r.IsPB {
		return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH}
	}
	return r.setValue(name, data)
}
func (r *Record) setValue(name string, data interface{}) error {
	if len(name) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME) {
		logger.ERR("value name len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME)
		return &terror.ErrorCode{Code: terror.ValueNameLenOverMax}
	}

	//check type
	var value []byte
	switch t := data.(type) {
	case bool:
		value = make([]byte, 1, 1)
		if t {
			value[0] = 1
		} else {
			value[0] = 0
		}
		break
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
		break
	case float64:
		value = make([]byte, 8, 8)
		binary.LittleEndian.PutUint64(value, math.Float64bits(t))
		break
	case []byte:
		if b, ok := data.([]byte); !ok {
			logger.ERR("value type not []byte")
			return &terror.ErrorCode{Code: terror.RecordKeyTypeInvalid}
		} else {
			value = b
		}
		break
	case string:
		//+ "\0"
		if str, ok := data.(string); !ok {
			logger.ERR("value type not string")
			return &terror.ErrorCode{Code: terror.RecordKeyTypeInvalid}
		} else {
			value = common.StringToCByte(str)
		}
		break
	default:
		logger.ERR("value type not support %v", t)
		return &terror.ErrorCode{Code: terror.RecordKeyTypeInvalid}
	}

	if len(value) > int(tcaplus_protocol_cs.TCAPLUS_MAX_VALUE_ALL_FIELDS_LEN) {
		logger.ERR("value len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_VALUE_ALL_FIELDS_LEN)
		return &terror.ErrorCode{Code: terror.ValueLenOverMax}
	}

	if len(r.ValueMap) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_VALUE_FIELD_NUM) {
		logger.ERR("value num over %d", tcaplus_protocol_cs.TCAPLUS_MAX_VALUE_FIELD_NUM)
		return &terror.ErrorCode{Code: terror.ValueNumOverMax}
	}

	r.ValueMap[name] = value
	return nil
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int8
*/
func (r *Record) SetValueInt8(name string, data int8) error {
	return r.SetValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int16
*/
func (r *Record) SetValueInt16(name string, data int16) error {
	return r.SetValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int32
*/
func (r *Record) SetValueInt32(name string, data int32) error {
	return r.SetValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int64
*/
func (r *Record) SetValueInt64(name string, data int64) error {
	return r.SetValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，float32
*/
func (r *Record) SetValueFloat32(name string, data float32) error {
	return r.SetValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，float64
*/
func (r *Record) SetValueFloat64(name string, data float64) error {
	return r.SetValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，string
*/
func (r *Record) SetValueStr(name string, data string) error {
	return r.SetValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，[]byte
*/
func (r *Record) SetValueBlob(name string, data []byte) error {
	return r.SetValue(name, data)
}

/**
@brief  加入要操作的字段名称及操作类型，若对应字段名之前已存在，则覆盖之。
@param  [in] field_name         字段名称，最大长度32字节
@param  [in] operation          操作类型，cmd.TcaplusApiOpPlus自增，Tcmd.TcaplusApiOpMinus自减
@param [IN] lower_limit         操作结果值下限，如果比这个值小，返回 SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE
@param [IN] upper_limit         操作结果值上限，如果比这个值大，返回 SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE
@note                           lower_limit == upper_limit 时，存储端不对操作结果进行范围检测
*/
/*
func (r *Record) AddValueOperation(name string, op uint32, lowerLimit int64, upperLimit int64) error {
	if len(name) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME) {
		logger.ERR("value name len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME)
		return &terror.ErrorCode{Code: terror.ValueNameLenOverMax}
	}

	if r.UpdFieldSet == nil {
		logger.ERR("cmd %d not support this op", r.Cmd)
		return &terror.ErrorCode{Code: terror.OperationNotSupport, Message: "cmd not support this func"}
	}

	if op != cmd.TcaplusApiOpPlus && op != cmd.TcaplusApiOpMinus {
		logger.ERR("op %d is invalid", op)
		return &terror.ErrorCode{Code: terror.OperationNotSupport, Message: "op is invalid"}
	}

	if lowerLimit >= upperLimit {
		logger.ERR("lowerLimit %d upperLimit %d", lowerLimit, upperLimit)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "lowerLimit must lower upperLimit"}
	}

	//find + modify
	for i := uint32(0); i < r.UpdFieldSet.FieldNum; i++ {
		if r.UpdFieldSet.Fields[i].FieldName == name {
			r.UpdFieldSet.Fields[i].FieldOperation = op
			r.UpdFieldSet.Fields[i].LowerLimit = lowerLimit
			r.UpdFieldSet.Fields[i].UpperLimit = upperLimit
			return nil
		}
	}

	if r.UpdFieldSet.FieldNum >= uint32(tcaplus_protocol_cs.TCAPLUS_MAX_VALUE_FIELD_NUM) {
		logger.ERR("value num over %d", tcaplus_protocol_cs.TCAPLUS_MAX_VALUE_FIELD_NUM)
		return &terror.ErrorCode{Code: terror.ValueNumOverMax}
	}

	//add
	r.UpdFieldSet.Fields[r.UpdFieldSet.FieldNum].FieldOperation = op
	r.UpdFieldSet.Fields[r.UpdFieldSet.FieldNum].LowerLimit = lowerLimit
	r.UpdFieldSet.Fields[r.UpdFieldSet.FieldNum].UpperLimit = upperLimit
	r.UpdFieldSet.FieldNum++
	return nil
}
*/

/**
	@brief  通用的key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容指针，最大长度1024字节，必须明确数据类型，必须和tdr xml表中的类型一致
								支持*bool, *byte, *int8, *int16, *uint16, *int32, *uint32, *int64, *uint64, *float32, *float64，*string, *[]byte
	@notice		请根据xml表准确填写类型，最好调用GetKeyInt8等接口
*/
func (r *Record) GetKey(name string, data interface{}) error {
	if r.IsPB {
		return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH}
	}
	return r.getKey(name, data)
}
func (r *Record) getKey(name string, data interface{}) error {
	if len(name) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME) {
		logger.ERR("key name len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME)
		return &terror.ErrorCode{Code: terror.KeyNameLenOverMax}
	}

	if keyData, exist := r.KeyMap[name]; !exist {
		logger.ERR("zone %d table %s key %s not exist", r.TableName, r.ZoneId, name)
		return &terror.ErrorCode{Code: terror.RecordKeyNotExist}
	} else {
		switch t := data.(type) {
		case *bool:
			*t = keyData[0] != 0
		case *int8:
			*t = int8(keyData[0])
		case *uint8:
			*t = keyData[0]
		case *int16:
			*t = *(*int16)(unsafe.Pointer(&keyData[0]))
		case *uint16:
			*t = *(*uint16)(unsafe.Pointer(&keyData[0]))
		case *int32:
			*t = *(*int32)(unsafe.Pointer(&keyData[0]))
		case *uint32:
			*t = *(*uint32)(unsafe.Pointer(&keyData[0]))
		case *int64:
			*t = *(*int64)(unsafe.Pointer(&keyData[0]))
		case *uint64:
			*t = *(*uint64)(unsafe.Pointer(&keyData[0]))
		case *float32:
			*t = *(*float32)(unsafe.Pointer(&keyData[0]))
		case *float64:
			*t = *(*float64)(unsafe.Pointer(&keyData[0]))
		case *[]byte:
			*t = keyData
		case *string:
			if len(keyData) > 1 {
				*t = string(keyData[0 : len(keyData)-1])
			} else {
				*t = ""
			}

		default:
			logger.ERR("key %s type not support %v", name, t)
			return &terror.ErrorCode{Code: terror.RecordKeyTypeInvalid}
		}
	}
	return nil
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int8
	@ret  [out] error         	nil success
*/
func (r *Record) GetKeyInt8(name string) (int8, error) {
	data := int8(0)
	err := r.GetKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int16
	@ret  [out] error         	nil success
*/
func (r *Record) GetKeyInt16(name string) (int16, error) {
	data := int16(0)
	err := r.GetKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int32
	@ret  [out] error         	nil success
*/
func (r *Record) GetKeyInt32(name string) (int32, error) {
	data := int32(0)
	err := r.GetKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int64
	@ret  [out] error         	nil success
*/
func (r *Record) GetKeyInt64(name string) (int64, error) {
	data := int64(0)
	err := r.GetKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，float32
	@ret  [out] error         	nil success
*/
func (r *Record) GetKeyFloat32(name string) (float32, error) {
	data := float32(0)
	err := r.GetKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，float32
	@ret  [out] error         	nil success
*/
func (r *Record) GetKeyFloat64(name string) (float64, error) {
	data := float64(0)
	err := r.GetKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，string
	@ret  [out] error         	nil success
*/
func (r *Record) GetKeyStr(name string) (string, error) {
	var data string
	err := r.GetKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，[]byte
	@ret  [out] error         	nil success
*/
func (r *Record) GetKeyBlob(name string) ([]byte, error) {
	if r.IsPB {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH}
	}
	return r.getKeyBlob(name)
}
func (r *Record) getKeyBlob(name string) ([]byte, error) {
	var data []byte
	err := r.GetKey(name, &data)
	return data, err
}

/**
	@brief  通用的value字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容指针，最大长度1024字节，必须明确数据类型，必须和tdr xml表中的类型一致
								支持*bool, *byte, *int8, *int16, *uint16, *int32, *uint32, *int64, *uint64, *float32, *float64，*string, *[]byte
	@notice		请根据xml表准确填写类型，最好调用GetValueInt8等接口
*/
func (r *Record) GetValue(name string, data interface{}) error {
	if r.IsPB {
		return &terror.ErrorCode{Code: terror.API_ERR_OPERATION_TYPE_NOT_MATCH}
	}
	return r.getValue(name, data)
}
func (r *Record) getValue(name string, data interface{}) error {
	if len(name) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME) {
		logger.ERR("value name len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME)
		return &terror.ErrorCode{Code: terror.KeyNameLenOverMax}
	}

	if valueData, exist := r.ValueMap[name]; !exist {
		logger.ERR("zone %d table %s value %s not exist", r.ZoneId, r.TableName, name)
		return &terror.ErrorCode{Code: terror.RecordValueNotExist}
	} else {
		switch t := data.(type) {
		case *bool:
			*t = valueData[0] != 0
		case *int8:
			*t = int8(valueData[0])
		case *uint8:
			*t = valueData[0]
		case *int16:
			*t = *(*int16)(unsafe.Pointer(&valueData[0]))
		case *uint16:
			*t = *(*uint16)(unsafe.Pointer(&valueData[0]))
		case *int32:
			*t = *(*int32)(unsafe.Pointer(&valueData[0]))
		case *uint32:
			*t = *(*uint32)(unsafe.Pointer(&valueData[0]))
		case *int64:
			*t = *(*int64)(unsafe.Pointer(&valueData[0]))
		case *uint64:
			*t = *(*uint64)(unsafe.Pointer(&valueData[0]))
		case *float32:
			*t = *(*float32)(unsafe.Pointer(&valueData[0]))
		case *float64:
			*t = *(*float64)(unsafe.Pointer(&valueData[0]))
		case *[]byte:
			*t = valueData
		case *string:
			if len(valueData) > 1 {
				*t = string(valueData[0 : len(valueData)-1])
			} else {
				*t = ""
			}

		default:
			logger.ERR("value %s type not support %v", name, t)
			return &terror.ErrorCode{Code: terror.RecordValueTypeInvalid}
		}
	}
	return nil
}

/**
	@brief  value字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int8
	@ret  [out] error         	nil success
*/
func (r *Record) GetValueInt8(name string) (int8, error) {
	data := int8(0)
	err := r.GetValue(name, &data)
	return data, err
}

/**
	@brief  value字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int16
	@ret  [out] error         	nil success
*/
func (r *Record) GetValueInt16(name string) (int16, error) {
	data := int16(0)
	err := r.GetValue(name, &data)
	return data, err
}

/**
	@brief  value字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int32
	@ret  [out] error         	nil success
*/
func (r *Record) GetValueInt32(name string) (int32, error) {
	data := int32(0)
	err := r.GetValue(name, &data)
	return data, err
}

/**
	@brief  value字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int64
	@ret  [out] error         	nil success
*/
func (r *Record) GetValueInt64(name string) (int64, error) {
	data := int64(0)
	err := r.GetValue(name, &data)
	return data, err
}

/**
	@brief  value字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，float32
	@ret  [out] error         	nil success
*/
func (r *Record) GetValueFloat32(name string) (float32, error) {
	data := float32(0)
	err := r.GetValue(name, &data)
	return data, err
}

/**
	@brief  value字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，float32
	@ret  [out] error         	nil success
*/
func (r *Record) GetValueFloat64(name string) (float64, error) {
	data := float64(0)
	err := r.GetValue(name, &data)
	return data, err
}

/**
	@brief  value字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，string
	@ret  [out] error         	nil success
*/
func (r *Record) GetValueStr(name string) (string, error) {
	var data string
	err := r.GetValue(name, &data)
	return data, err
}

/**
	@brief  value字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，[]byte
	@ret  [out] error         	nil success
*/
func (r *Record) GetValueBlob(name string) ([]byte, error) {
	var data []byte
	err := r.GetValue(name, &data)
	return data, err
}
