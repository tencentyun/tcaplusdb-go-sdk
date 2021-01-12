package record

import (
	"bytes"
	"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

func (r *Record) PackKey() error {
	if nil == r.KeySet {
		logger.ERR("record keySet is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "record keySet is nil"}
	}

	r.KeySet.Fields = make([]*tcaplus_protocol_cs.TCaplusKeyField, len(r.KeyMap))
	r.KeySet.FieldNum = 0
	for name, v := range r.KeyMap {
		r.KeySet.Fields[r.KeySet.FieldNum] = tcaplus_protocol_cs.NewTCaplusKeyField()
		//key-name
		r.KeySet.Fields[r.KeySet.FieldNum].FieldName = name
		//key-value
		r.KeySet.Fields[r.KeySet.FieldNum].FieldBuff = v
		//key-value-len
		r.KeySet.Fields[r.KeySet.FieldNum].FieldLen = uint32(len(v))
		//key-num ++
		r.KeySet.FieldNum++
	}
	r.KeySet.Version = r.Version
	return nil
}

/*
 * 使用10M的紧凑模式和batchget的连续编码
 * ----------------------------------------------------------------------------------------
 * | field_num(4B) | version(4B) | [ size of name(2B) |  name ... | buf_len(4B) | buf ... ]
 * ----------------------------------------------------------------------------------------
 */
func (r *Record) PackValue(valueNameMap map[string]bool) error {
	if nil == r.ValueSet {
		logger.ERR("record valueSet is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "record valueSet is nil"}
	}
	r.ValueSet.Version_ = r.Version
	return r.packCompactValueSet(r.ValueSet.CompactValueSet, valueNameMap)
}

func (r *Record) packCompactValueSet(compactValueSet *tcaplus_protocol_cs.CompactValueSet,
					valueNameMap map[string]bool) error {
	if nil == compactValueSet {
		errMsg := "record valueSet compactValueSet is nil"
		logger.ERR(errMsg)
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: errMsg}
	}

	valueBuf := new(bytes.Buffer)
	compactValueSet.ValueBufLen = 8 //field_num(4B) + version(4B)
	//set fieldNum
	if err := binary.Write(valueBuf, binary.LittleEndian, int32(0)); err != nil {
		return err
	}
	//set version

	if err := binary.Write(valueBuf, binary.LittleEndian, r.Version); err != nil {
		return err
	}

	//set name + data + index
	compactValueSet.FieldIndexs = make([]*tcaplus_protocol_cs.FieldIndex, len(r.ValueMap))
	compactValueSet.FieldIndexNum = 0
	for name, v := range r.ValueMap {
		//部分value字段查询和更新
		if valueNameMap != nil && len(valueNameMap) > 0 {
			if _, exist := valueNameMap[name]; !exist {
				continue
			}
		}

		//check total len( namelen(2B) + namedata + "\0" + dataLen(4B) + data)
		totalLen := 2 + len(name) + 1 + 4 + len(v)
		if int(compactValueSet.ValueBufLen)+totalLen >= int(tcaplus_protocol_cs.TCAPLUS_GROSS_MAX_VALUE_BUFFER_LEN) {
			logger.ERR("record value pack too large %d", int(compactValueSet.ValueBufLen)+totalLen)
			return &terror.ErrorCode{Code: terror.ValuePackOverMax}
		}

		//write name len
		nameLen := int16(len(name) + 1)
		if err := binary.Write(valueBuf, binary.LittleEndian, nameLen); err != nil {
			return err
		}

		//write name + "\0"
		valueBuf.Write(common.StringToCByte(name))

		//write data len
		vLen := int32(len(v))
		if err := binary.Write(valueBuf, binary.LittleEndian, vLen); err != nil {
			return err
		}

		//write data
		valueBuf.Write(v)

		//set index
		index := tcaplus_protocol_cs.NewFieldIndex()
		index.Offset = compactValueSet.ValueBufLen
		index.Size = int32(valueBuf.Len()) - compactValueSet.ValueBufLen
		index.Flag = 0
		compactValueSet.FieldIndexs[compactValueSet.FieldIndexNum] = index
		compactValueSet.FieldIndexNum++
		compactValueSet.ValueBufLen = int32(valueBuf.Len())
	}
	compactValueSet.ValueBuf = valueBuf.Bytes()

	//reset fieldNum
	fieldNumBuf := new(bytes.Buffer)
	if err := binary.Write(fieldNumBuf, binary.LittleEndian, compactValueSet.FieldIndexNum); err != nil {
		return err
	}
	copy(compactValueSet.ValueBuf[0:], fieldNumBuf.Bytes())
	return nil
}

func (r *Record) PackPBFieldValue() error {
	if nil == r.PBValueSet {
		logger.ERR("record PBValueSet is nil")
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "record PBValueSet is nil"}
	}
	r.PBValueSet.Version_ = r.Version
	return r.packCompactValueSet(r.PBValueSet.CompactValueSet, r.PBFieldMap)
}
