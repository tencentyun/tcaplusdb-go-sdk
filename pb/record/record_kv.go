package record

import (
	"bytes"
	"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

/**
	@brief  通用的key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，最大长度1024字节，必须明确数据类型，必须和tdr xml表中的类型一致
								支持bool, byte, int8, int16, uint16, int32, uint32, int64, uint64, float32, float64，string, []byte
	@notice		请根据xml表准确填写类型，最好调用SetKeyInt8等接口
*/
func (r *Record) setKey(name string, data interface{}) error {
	if len(name) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME) {
		logger.ERR("key name len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME)
		return &terror.ErrorCode{Code: terror.KeyNameLenOverMax}
	}

	//check type
	var value []byte
	switch t := data.(type) {
	case bool, byte, int8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
			return err
		}
		value = buf.Bytes()
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
			value = []byte(str)
			value = append(value, 0)
		}
		break
	default:
		logger.ERR("key type not support %v", t)
		return &terror.ErrorCode{Code: terror.RecordKeyTypeInvalid}
	}

	if len(value) > int(tcaplus_protocol_cs.TCAPLUS_BIG_RECORD_MAX_VALUE_BUF_LEN) {
		logger.ERR("key len over %d", tcaplus_protocol_cs.TCAPLUS_BIG_RECORD_MAX_VALUE_BUF_LEN)
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
func (r *Record) setKeyInt8(name string, data int8) error {
	return r.setKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int16
*/
func (r *Record) setKeyInt16(name string, data int16) error {
	return r.setKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int32
*/
func (r *Record) setKeyInt32(name string, data int32) error {
	return r.setKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int64
*/
func (r *Record) setKeyInt64(name string, data int64) error {
	return r.setKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，float32
*/
func (r *Record) setKeyFloat32(name string, data float32) error {
	return r.setKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，float64
*/
func (r *Record) setKeyFloat64(name string, data float64) error {
	return r.setKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，string
*/
func (r *Record) setKeyStr(name string, data string) error {
	return r.setKey(name, data)
}

/**
	@brief  key字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，[]byte
*/
func (r *Record) setKeyBlob(name string, data []byte) error {
	return r.setKey(name, data)
}

/**
	@brief  通用的value字段内容设置
 	@param  [in] name         	字段名称，最大长度32
	@param  [in] data         	字段内容，最大长度1024字节，必须明确数据类型，必须和tdr xml表中的类型一致
								支持bool, byte, int8, int16, uint16, int32, uint32, int64, uint64, float32, float64，string, []byte
	@notice		请根据xml表准确填写类型，最好调用SetValueInt8等接口
*/
func (r *Record) setValue(name string, data interface{}) error {
	if len(name) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME) {
		logger.ERR("value name len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME)
		return &terror.ErrorCode{Code: terror.ValueNameLenOverMax}
	}

	//check type
	var value []byte
	switch t := data.(type) {
	case bool, byte, int8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
			return err
		}
		value = buf.Bytes()
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
			value = []byte(str)
			value = append(value, 0)
		}
		break
	default:
		logger.ERR("value type not support %v", t)
		return &terror.ErrorCode{Code: terror.RecordKeyTypeInvalid}
	}

	if len(value) > int(tcaplus_protocol_cs.TCAPLUS_BIG_RECORD_MAX_VALUE_BUF_LEN) {
		logger.ERR("value len over %d", tcaplus_protocol_cs.TCAPLUS_BIG_RECORD_MAX_VALUE_BUF_LEN)
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
func (r *Record) setValueInt8(name string, data int8) error {
	return r.setValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int16
*/
func (r *Record) setValueInt16(name string, data int16) error {
	return r.setValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int32
*/
func (r *Record) setValueInt32(name string, data int32) error {
	return r.setValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，int64
*/
func (r *Record) setValueInt64(name string, data int64) error {
	return r.setValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，float32
*/
func (r *Record) setValueFloat32(name string, data float32) error {
	return r.setValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，float64
*/
func (r *Record) setValueFloat64(name string, data float64) error {
	return r.setValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，string
*/
func (r *Record) setValueStr(name string, data string) error {
	return r.setValue(name, data)
}

/**
	@brief  Value字段内容设置
    @param  [in] name         	字段名称，最大长度32
    @param  [in] data         	字段内容，[]byte
*/
func (r *Record) setValueBlob(name string, data []byte) error {
	return r.setValue(name, data)
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
		case *bool, *byte, *int8, *int16, *uint16, *int32, *uint32, *int64, *uint64, *float32, *float64:
			if err := binary.Read(bytes.NewReader(keyData), binary.LittleEndian, data); err != nil {
				logger.ERR("zone %d table %s key %s binary.Read err %s", r.TableName, r.ZoneId, name, err.Error())
				return err
			}
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
func (r *Record) getKeyInt8(name string) (int8, error) {
	data := int8(0)
	err := r.getKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int16
	@ret  [out] error         	nil success
*/
func (r *Record) getKeyInt16(name string) (int16, error) {
	data := int16(0)
	err := r.getKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int32
	@ret  [out] error         	nil success
*/
func (r *Record) getKeyInt32(name string) (int32, error) {
	data := int32(0)
	err := r.getKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，int64
	@ret  [out] error         	nil success
*/
func (r *Record) getKeyInt64(name string) (int64, error) {
	data := int64(0)
	err := r.getKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，float32
	@ret  [out] error         	nil success
*/
func (r *Record) getKeyFloat32(name string) (float32, error) {
	data := float32(0)
	err := r.getKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，float32
	@ret  [out] error         	nil success
*/
func (r *Record) getKeyFloat64(name string) (float64, error) {
	data := float64(0)
	err := r.getKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，string
	@ret  [out] error         	nil success
*/
func (r *Record) getKeyStr(name string) (string, error) {
	var data string
	err := r.getKey(name, &data)
	return data, err
}

/**
	@brief  key字段内容获取
    @param  [in] name         	字段名称，最大长度32
    @ret  [out] data         	字段内容，[]byte
	@ret  [out] error         	nil success
*/
func (r *Record) getKeyBlob(name string) ([]byte, error) {
	var data []byte
	err := r.getKey(name, &data)
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
	if len(name) >= int(tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME) {
		logger.ERR("value name len over %d", tcaplus_protocol_cs.TCAPLUS_MAX_FIELD_NAME)
		return &terror.ErrorCode{Code: terror.KeyNameLenOverMax}
	}

	if valueData, exist := r.ValueMap[name]; !exist {
		logger.ERR("zone %d table %s value %s not exist", r.ZoneId, r.TableName, name)
		return &terror.ErrorCode{Code: terror.RecordValueNotExist}
	} else {
		switch t := data.(type) {
		case *bool, *byte, *int8, *int16, *uint16, *int32, *uint32, *int64, *uint64, *float32, *float64:
			if err := binary.Read(bytes.NewReader(valueData), binary.LittleEndian, data); err != nil {
				logger.ERR("zone %d table %s key %s binary.Read err %s", r.ZoneId, r.TableName, name, err.Error())
				return err
			}
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
