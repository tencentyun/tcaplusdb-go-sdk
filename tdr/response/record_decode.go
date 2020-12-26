package response

import (
	"bytes"
	"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	//	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	//	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

/*
| FieldNum 2 Bytes | [{namelen  2 bytes | name  n bytes | fieldlen 2 bytes | filed n bytes} ....{ }]
*/
func unpack_record_key(buff []byte, len int32, KeyMap map[string][]byte, offset_out *int32) error {
	offset := int32(0)
	FieldNum := uint16(0)
	if err := binary.Read(bytes.NewReader(buff[0:2]), binary.LittleEndian, &FieldNum); err != nil {
		logger.ERR("read FieldNum failed %s", err.Error())
		return err
	}
	offset += 2
	len -= 2

	namelen := uint16(0)
	for idx := uint16(0); idx < FieldNum; idx++ {
		if len < 2 {
			logger.ERR("KeyNumOverMax: len:%d", len)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		// name 获取长度再获取值
		namelen = 0
		if err := binary.Read(bytes.NewReader(buff[offset:offset+2]), binary.LittleEndian, &namelen); err != nil {
			logger.ERR("read name len failed %s", err.Error())
			return err
		}

		offset += 2
		len -= 2
		//  name必须不能为空
		if int32(namelen) > len || namelen == 0 {
			logger.ERR("name len %d, name len error", namelen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name := string(buff[offset : offset+int32(namelen)-1])

		logger.DEBUG("namelen: %d, name : %s", namelen, name)
		offset += int32(namelen)
		len -= int32(namelen)

		//field
		FieldLen := uint32(0)
		if err := binary.Read(bytes.NewReader(buff[offset:offset+4]), binary.LittleEndian, &FieldLen); err != nil {
			logger.ERR("read FieldLen failed %s", err.Error())
			return err
		}

		offset += 4
		len -= 4

		if int32(FieldLen) > len {
			logger.ERR("Field len %d, out of data", FieldLen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		KeyMap[name] = buff[offset : offset+int32(FieldLen)]
		//logger.DEBUG("FieldLen: %d, Fields : %s", FieldLen, Fields)
		offset += int32(FieldLen)
		len -= int32(FieldLen)
	}
	*offset_out = offset
	return nil
}

/*
| FieldNum 2 Bytes |version 4 bytes | [{namelen  2 bytes | name  n bytes | fieldlen 4 bytes | filed n bytes} ....{ }]
*/
func unpack_record_value(value_buff []byte, len int32, valueMap map[string][]byte, offset_out *int32) error {
	FieldNum := uint32(0)
	offset := int32(0)

	if err := binary.Read(bytes.NewReader(value_buff[offset:offset+4]), binary.LittleEndian, &FieldNum); err != nil {
		logger.ERR("read FieldNum failed %s", err.Error())
		return err
	}
	offset += 4
	len -= 4

	version := uint32(0)
	if err := binary.Read(bytes.NewReader(value_buff[offset:offset+4]), binary.LittleEndian, &version); err != nil {
		logger.ERR("read version failed %s", err.Error())
		return err
	}

	offset += 4
	len -= 4

	for idx := uint32(0); idx < FieldNum; idx++ {
		namelen := int16(0)
		if err := binary.Read(bytes.NewReader(value_buff[offset:offset+2]), binary.LittleEndian, &namelen); err != nil {
			logger.ERR("read FieldNum failed %s", err.Error())
			return err
		}
		len -= 2
		offset += 2
		if int32(namelen) > len || namelen == 0 {
			logger.ERR("name len %d, out of data", namelen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name := string(value_buff[offset : offset+int32(namelen)-1])

		len -= int32(namelen)
		offset += int32(namelen)

		bufflen := int32(0)
		if err := binary.Read(bytes.NewReader(value_buff[offset:offset+4]), binary.LittleEndian, &bufflen); err != nil {
			logger.ERR("read FieldNum failed %s", err.Error())
			return err
		}
		len -= 4
		offset += 4
		if int32(bufflen) > len {
			logger.ERR("buff len %d, out of data", bufflen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}

		valueMap[name] = value_buff[offset : offset+bufflen]
		len -= bufflen
		offset += bufflen
	}
	*offset_out = offset
	return nil
}

func unpack_record_k_v(buff []byte, len int32, KeyMap map[string][]byte, valueMap map[string][]byte) (int32, error) {

	var offset_key int32 = 0
	if err := unpack_record_key(buff, len, KeyMap, &offset_key); err != nil {
		return 0, err
	}
	var offset_value int32 = 0
	if err := unpack_record_value(buff[offset_key:len], len-offset_key, valueMap, &offset_value); err != nil {
		return 0, err
	}
	return offset_key + offset_value, nil
}

/*
just for TCaplusGetByPartKeyResultSucc.BatchValueInfo  please check if ok for others
*/
func unpack_batchvalue_with_version_1(buff []byte, len int32, KeyMap map[string][]byte,
	valueMap map[string][]byte) (int32, error) {
	OffSet := int32(0)
	nEncodeVersion := uint16(0)
	if err := binary.Read(bytes.NewReader(buff[OffSet:OffSet+2]), binary.LittleEndian, &nEncodeVersion); err != nil {
		logger.ERR("read nEncodeVersion failed %s", err.Error())
		return 0, err
	}
	if nEncodeVersion != 1 {
		logger.ERR("nEncodeVersion %d error", nEncodeVersion)
		return 0, &terror.ErrorCode{Code: -3}
	}
	OffSet += 2

	ullLastAccessTime := uint64(0)
	if err := binary.Read(bytes.NewReader(buff[OffSet:OffSet+8]), binary.LittleEndian, &ullLastAccessTime); err != nil {
		logger.ERR("read nEncodeVersion failed %s", err.Error())
		return 0, err
	}
	OffSet += 8

	read_bytes := int32(0)
	var err error
	if read_bytes, err = unpack_record_k_v(buff[OffSet:len], len-OffSet, KeyMap, valueMap); err != nil {
		return 0, err
	}
	OffSet += read_bytes
	return OffSet, nil
}

/*
just for TCaplusGetByPartKeyResultSucc.BatchValueInfo  please check if ok for others
*/

func unpack_record(buff []byte, len int32, KeyMap map[string][]byte, valueMap map[string][]byte) (int32, error) {
	OffSet := int32(0)

	nEncodeVersion := uint16(0)
	if err := binary.Read(bytes.NewReader(buff[OffSet:OffSet+2]), binary.LittleEndian, &nEncodeVersion); err != nil {
		logger.ERR("read nEncodeVersion failed %s", err.Error())
		return 0, err
	}
	if nEncodeVersion != 1 {
		logger.ERR("nEncodeVersion %d error", nEncodeVersion)
		return 0, &terror.ErrorCode{Code: -3}
	}
	OffSet += 2

	read_bytes := int32(0)
	var err error
	if read_bytes, err = unpack_batchvalue_with_accesstime(buff[OffSet:len], len-OffSet,
		KeyMap, valueMap, nil); err != nil {
		return 0, err
	}
	OffSet += read_bytes
	return OffSet, nil
}

func unpack_batchvalue_with_accesstime(buff []byte, len int32, KeyMap map[string][]byte, valueMap map[string][]byte,
	lastAccessTime *uint64) (int32, error) {
	OffSet := int32(0)

	ullLastAccessTime := uint64(0)
	if err := binary.Read(bytes.NewReader(buff[OffSet:OffSet+8]), binary.LittleEndian, &ullLastAccessTime); err != nil {
		logger.ERR("read nEncodeVersion failed %s", err.Error())
		return 0, err
	}
	if nil != lastAccessTime {
		*lastAccessTime = ullLastAccessTime
	}

	OffSet += 8

	read_bytes := int32(0)
	var err error
	if read_bytes, err = unpack_record_k_v(buff[OffSet:len], len-OffSet, KeyMap, valueMap); err != nil {
		return 0, err
	}
	OffSet += read_bytes
	return OffSet, nil
}

func unpack_suc_keys_buffLen(key_buff []byte, len int32, result *int32,
	keyMap map[string][]byte, offset_out *int32) error {
	if len < 8 {
		logger.ERR("len %d, less then 8", len)
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}

	offset := int32(0)

	*result = 0
	if err := binary.Read(bytes.NewReader(key_buff[offset:offset+4]), binary.LittleEndian, result); err != nil {
		logger.ERR("read result failed %s", err.Error())
		return err
	}
	offset += 4
	len -= 4

	FieldNum := uint32(0)
	if err := binary.Read(bytes.NewReader(key_buff[offset:offset+4]), binary.LittleEndian, &FieldNum); err != nil {
		logger.ERR("read FieldNum failed %s", err.Error())
		return err
	}
	offset += 4
	len -= 4

	for idx := uint32(0); idx < FieldNum; idx++ {
		namelen := int16(0)
		if err := binary.Read(bytes.NewReader(key_buff[offset:offset+2]), binary.LittleEndian, &namelen); err != nil {
			logger.ERR("read FieldNum failed %s", err.Error())
			return err
		}
		len -= 2
		offset += 2
		if int32(namelen) > len || namelen == 0 {
			logger.ERR("name len %d, out of data", namelen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name := string(key_buff[offset : offset+int32(namelen)-1])

		len -= int32(namelen)
		offset += int32(namelen)

		bufflen := int32(0)
		if err := binary.Read(bytes.NewReader(key_buff[offset:offset+4]), binary.LittleEndian, &bufflen); err != nil {
			logger.ERR("read FieldNum failed %s", err.Error())
			return err
		}
		len -= 4
		offset += 4
		if int32(bufflen) > len {
			logger.ERR("buff len %d, out of data", bufflen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		keyMap[name] = key_buff[offset : offset+bufflen]
		logger.DEBUG("upack key: %s,value: %s", name, keyMap[name])

		len -= bufflen
		offset += bufflen
	}
	*offset_out = offset
	return nil
}
func check_less(small uint32, big uint32) error {
	if small > big {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	return nil
}
func unpack_element_buff(start_buff []byte, offset_in uint32, elem_len uint32, index *int32, offset_out *uint32,
	valueMap map[string][]byte) error {
	if elem_len < 4+offset_in {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	curr_buff := start_buff[offset_in:elem_len]
	*offset_out = 0
	offset := uint32(0)
	*index = 0
	if offset_in+offset+4 > elem_len {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	if err := binary.Read(bytes.NewReader(curr_buff[offset:offset+4]), binary.LittleEndian, index); err != nil {
		logger.ERR("read FieldNum failed %s", err.Error())
		return err
	}
	logger.DEBUG("index: %d", *index)
	offset += 4
	element_vserion := int32(0)
	if offset_in+offset+4 > elem_len {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	if err := binary.Read(bytes.NewReader(curr_buff[offset:offset+4]), binary.LittleEndian, &element_vserion); err != nil {
		logger.ERR("read FieldNum failed %s", err.Error())
		return err
	}
	logger.DEBUG("element_vserion: %d", element_vserion)
	offset += 4
	field_num := uint32(0)
	if offset_in+offset+4 > elem_len {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	if err := binary.Read(bytes.NewReader(curr_buff[offset:offset+4]), binary.LittleEndian, &field_num); err != nil {
		logger.ERR("read FieldNum failed %s", err.Error())
		return err
	}
	if offset_in+offset+4 > elem_len {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	logger.DEBUG("field_num: %d", field_num)
	offset += 4
	//FieldsBuffLen, 这里这个值没啥用，跳过去了
	offset += 4
	for idx := uint32(0); idx < field_num; idx++ {
		name_len := bytes.IndexByte(curr_buff[offset:], 0)
		if name_len < 0 {
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name_len += 1
		if offset_in+offset+uint32(name_len) > elem_len {
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		filed_len := uint32(0)
		if err := binary.Read(bytes.NewReader(curr_buff[offset+uint32(name_len):offset+uint32(name_len)+4]),
			binary.LittleEndian, &filed_len); err != nil {
			logger.ERR("read FieldNum failed %s", err.Error())
			return err
		}
		if offset_in+offset+uint32(name_len)+4+filed_len > elem_len {
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name := string(curr_buff[offset : offset+uint32(name_len)-1])
		valueMap[name] = curr_buff[offset+uint32(name_len)+4 : offset+uint32(name_len)+4+uint32(filed_len)]
		offset = offset + uint32(name_len) + 4 + uint32(filed_len)
		//logger.ERR("name:[%s], len:%d",name, len(name))
	}
	*offset_out = offset
	return nil
}
