package record

import (
	"bytes"
	"encoding/binary"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
)

func (r *Record) UnPackKey() error {
	if nil == r.KeySet {
		logger.ERR("record keySet is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "record keySet is nil"}
	}

	if r.KeySet.FieldNum <= 0 {
		logger.ERR("record keySet is empty")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "record keySet is empty"}
	}

	for i := 0; i < int(r.KeySet.FieldNum); i++ {
		//key-name
		name := r.KeySet.Fields[i].FieldName
		//key-data
		data := r.KeySet.Fields[i].FieldBuff[0:r.KeySet.Fields[i].FieldLen]
		r.KeyMap[name] = data
	}
	r.Version = r.KeySet.Version
	return nil
}

/*
 * 使用10M的紧凑模式和batchget的连续编码
 * ----------------------------------------------------------------------------------------
 * | field_num(4B) | version(4B) | [ size of name(2B) |  name ... | buf_len(4B) | buf ... ]
 * ----------------------------------------------------------------------------------------
 */
func (r *Record) UnPackValue() error {
	if nil == r.ValueSet {
		logger.ERR("record valueSet is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "record valueSet is nil"}
	}

	return r.unPackCompactValueSet(r.ValueSet.CompactValueSet)
}

func (r *Record) unPackCompactValueSet(compactValueSet *tcaplus_protocol_cs.CompactValueSet) error {
	if nil == compactValueSet {
		errMsg := "record valueSet compactValueSet is nil"
		logger.ERR(errMsg)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: errMsg}
	}
	//get fieldNum(4B)
	fieldNum := int32(0)
	if err := binary.Read(bytes.NewReader(compactValueSet.ValueBuf[0:4]), binary.LittleEndian, &fieldNum); err != nil {
		logger.ERR("read fieldNum failed %s", err.Error())
		return err
	}

	if fieldNum != compactValueSet.FieldIndexNum {
		logger.ERR("compactValueSet fieldNum %d not equal compactValueSet.FieldIndexNum %d",
			compactValueSet.FieldIndexNum)
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}

	//get version(4B)
	version := int32(0)
	if err := binary.Read(bytes.NewReader(compactValueSet.ValueBuf[4:8]), binary.LittleEndian, &version); err != nil {
		logger.ERR("read version failed %s", err.Error())
		return err
	}
	if r.Version != version {
		logger.WARN("value version %d not equal key version %d", version, r.Version)
	}

	for i := 0; i < int(fieldNum); i++ {
		offset := compactValueSet.FieldIndexs[i].Offset
		end := offset + compactValueSet.FieldIndexs[i].Size

		//name-len(2B)
		if offset+2 > end {
			logger.ERR("read offset is invalid %d", offset)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		nameLen := int16(0)
		if err := binary.Read(bytes.NewReader(compactValueSet.ValueBuf[offset:offset+2]),
			binary.LittleEndian, &nameLen); err != nil {
			logger.ERR("read nameLen failed %s", err.Error())
			return err
		}
		if nameLen <= 1 {
			//name 以0结尾必然大于1
			logger.ERR("read nameLen is invalid %d", nameLen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		offset += 2

		//value-name
		if offset+int32(nameLen) > end {
			logger.ERR("read offset is invalid %d", offset)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name := string(compactValueSet.ValueBuf[offset : offset+int32(nameLen)-1])
		offset += int32(nameLen)

		//data-len(4B)
		if offset+4 > end {
			logger.ERR("read offset is invalid %d", offset)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		dataLen := int32(0)
		if err := binary.Read(bytes.NewReader(compactValueSet.ValueBuf[offset:offset+4]),
			binary.LittleEndian, &dataLen); err != nil {
			logger.ERR("read value dataLen failed %s", err.Error())
			return err
		}
		if dataLen < 0 {
			logger.ERR("read dataLen is invalid %d", dataLen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		offset += 4

		//value-data
		if offset+dataLen > end {
			logger.ERR("read offset is invalid %d", offset)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		data := compactValueSet.ValueBuf[offset : offset+dataLen]

		//set map
		r.ValueMap[name] = data
	}
	return nil
}

func (r *Record) UnPackPBValue() error {
	if nil == r.PBValueSet {
		logger.ERR("record valueSet is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "record PBValueSet is nil"}
	}

	return r.unPackCompactValueSet(r.PBValueSet.CompactValueSet)
}
