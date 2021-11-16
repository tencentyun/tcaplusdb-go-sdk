package response

import (
	"bytes"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"unsafe"
)

/*
| FieldNum 2 Bytes | [{namelen  2 bytes | name  n bytes | fieldlen 2 bytes | filed n bytes} ....{ }]
*/
func unpackRecordKey(buff []byte, len int32, KeyMap map[string][]byte, offsetOut *int32) error {
	offset := int32(0)
	FieldNum := *(*uint16)(unsafe.Pointer(&buff[0]))
	offset += 2
	len -= 2

	nameLen := uint16(0)
	for idx := uint16(0); idx < FieldNum; idx++ {
		if len < 2 {
			logger.ERR("KeyNumOverMax: len:%d", len)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		// name 获取长度再获取值
		nameLen = *(*uint16)(unsafe.Pointer(&buff[offset]))
		offset += 2
		len -= 2
		//  name必须不能为空
		if int32(nameLen) > len || nameLen == 0 {
			logger.ERR("name len %d, name len error", nameLen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name := string(buff[offset : offset+int32(nameLen)-1])

		logger.DEBUG("namelen: %d, name : %s", nameLen, name)
		offset += int32(nameLen)
		len -= int32(nameLen)

		//field
		FieldLen := *(*uint32)(unsafe.Pointer(&buff[offset]))
		offset += 4
		len -= 4

		if int32(FieldLen) > len {
			logger.ERR("Field len %d, out of data", FieldLen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		KeyMap[name] = buff[offset : offset+int32(FieldLen)]
		//logger.DEBUG("FieldLen: %d, Fields : %s", FieldLen, string(KeyMap[name]))
		offset += int32(FieldLen)
		len -= int32(FieldLen)
	}
	*offsetOut = offset
	return nil
}

/*
| FieldNum 2 Bytes |version 4 bytes | [{namelen  2 bytes | name  n bytes | fieldlen 4 bytes | filed n bytes} ....{ }]
*/
func unpackRecordValue(valueBuff []byte, len int32, valueMap map[string][]byte, offsetOut *int32, version *int32) error {
	offset := int32(0)
	FieldNum := *(*uint32)(unsafe.Pointer(&valueBuff[offset]))
	offset += 4
	len -= 4

	*version = *(*int32)(unsafe.Pointer(&valueBuff[offset]))
	logger.DEBUG("unpack_record_value version %d", *version)
	offset += 4
	len -= 4

	for idx := uint32(0); idx < FieldNum; idx++ {
		namelen := *(*uint16)(unsafe.Pointer(&valueBuff[offset]))
		len -= 2
		offset += 2
		if int32(namelen) > len || namelen == 0 {
			logger.ERR("name len %d, out of data", namelen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name := string(valueBuff[offset : offset+int32(namelen)-1])

		len -= int32(namelen)
		offset += int32(namelen)

		bufflen := *(*int32)(unsafe.Pointer(&valueBuff[offset]))
		len -= 4
		offset += 4
		if bufflen > len {
			logger.ERR("buff len %d, out of data", bufflen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}

		valueMap[name] = valueBuff[offset : offset+bufflen]
		len -= bufflen
		offset += bufflen
	}
	*offsetOut = offset
	return nil
}

func unpackRecordKV(buff []byte, len int32, KeyMap map[string][]byte, valueMap map[string][]byte,
	version *int32) (int32, error) {

	var offsetKey int32 = 0
	if err := unpackRecordKey(buff, len, KeyMap, &offsetKey); err != nil {
		return 0, err
	}
	var offsetValue int32 = 0
	if err := unpackRecordValue(buff[offsetKey:len], len-offsetKey, valueMap, &offsetValue, version); err != nil {
		return 0, err
	}
	return offsetKey + offsetValue, nil
}

/*
just for TCaplusGetByPartKeyResultSucc.BatchValueInfo  please check if ok for others
*/
func unpackBatchValueWithVersion1(buff []byte, len int32, KeyMap map[string][]byte,
	valueMap map[string][]byte) (int32, error) {
	OffSet := int32(0)
	nEncodeVersion := *(*uint16)(unsafe.Pointer(&buff[OffSet]))
	if nEncodeVersion != 1 {
		logger.ERR("nEncodeVersion %d error", nEncodeVersion)
		return 0, &terror.ErrorCode{Code: -3}
	}
	OffSet += 2

	ullLastAccessTime := *(*uint64)(unsafe.Pointer(&buff[OffSet]))
	logger.DEBUG("unpack_batchvalue_with_version_1 ullLastAccessTime %d", ullLastAccessTime)
	OffSet += 8

	readBytes := int32(0)
	var err error
	var version int32
	if readBytes, err = unpackRecordKV(buff[OffSet:len], len-OffSet, KeyMap, valueMap, &version); err != nil {
		return 0, err
	}
	OffSet += readBytes
	return OffSet, nil
}

/*
just for TCaplusGetByPartKeyResultSucc.BatchValueInfo  please check if ok for others
*/

func unpackRecord(buff []byte, len int32, KeyMap map[string][]byte, valueMap map[string][]byte, version *int32) (int32,
	error) {
	OffSet := int32(0)

	nEncodeVersion := *(*uint16)(unsafe.Pointer(&buff[OffSet]))
	if nEncodeVersion != 1 {
		logger.ERR("nEncodeVersion %d error", nEncodeVersion)
		return 0, &terror.ErrorCode{Code: -3}
	}
	OffSet += 2

	readBytes := int32(0)
	var err error
	if readBytes, err = unpackBatchValueWithAccessTime(buff[OffSet:len], len-OffSet,
		KeyMap, valueMap, nil, version); err != nil {
		return 0, err
	}
	OffSet += readBytes
	return OffSet, nil
}

func unpackBatchValueWithAccessTime(buff []byte, len int32, KeyMap map[string][]byte, valueMap map[string][]byte,
	lastAccessTime *uint64, version *int32) (int32, error) {
	OffSet := int32(0)

	ullLastAccessTime := *(*uint64)(unsafe.Pointer(&buff[OffSet]))
	if nil != lastAccessTime {
		*lastAccessTime = ullLastAccessTime
	}

	OffSet += 8

	readBytes := int32(0)
	var err error
	if readBytes, err = unpackRecordKV(buff[OffSet:len], len-OffSet, KeyMap, valueMap, version); err != nil {
		return 0, err
	}
	OffSet += readBytes
	return OffSet, nil
}

func unpackSucKeysBuffLen(keyBuff []byte, len int32, result *int32,
	keyMap map[string][]byte, offsetOut *int32) error {
	if len < 8 {
		logger.ERR("len %d, less then 8", len)
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}

	offset := int32(0)

	*result = *(*int32)(unsafe.Pointer(&keyBuff[offset]))
	offset += 4
	len -= 4

	FieldNum := *(*uint32)(unsafe.Pointer(&keyBuff[offset]))
	offset += 4
	len -= 4

	for idx := uint32(0); idx < FieldNum; idx++ {
		nameLen := *(*uint16)(unsafe.Pointer(&keyBuff[offset]))
		len -= 2
		offset += 2
		if int32(nameLen) > len || nameLen == 0 {
			logger.ERR("name len %d, out of data", nameLen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name := string(keyBuff[offset : offset+int32(nameLen)-1])

		len -= int32(nameLen)
		offset += int32(nameLen)

		buffLen := *(*int32)(unsafe.Pointer(&keyBuff[offset]))
		len -= 4
		offset += 4
		if buffLen > len {
			logger.ERR("buff len %d, out of data", buffLen)
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		keyMap[name] = keyBuff[offset : offset+buffLen]
		logger.DEBUG("upack key: %s,value: %s", name, keyMap[name])

		len -= buffLen
		offset += buffLen
	}
	*offsetOut = offset
	return nil
}
func checkLess(small uint32, big uint32) error {
	if small > big {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	return nil
}
func unpackElementBuff(startBuff []byte, offsetIn uint32, elemLen uint32, index *int32, offsetOut *uint32,
	valueMap map[string][]byte) error {
	if elemLen < 4+offsetIn {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	currBuff := startBuff[offsetIn:elemLen]
	*offsetOut = 0
	offset := uint32(0)
	if offsetIn+offset+4 > elemLen {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	*index = *(*int32)(unsafe.Pointer(&currBuff[offset]))
	logger.DEBUG("index: %d", *index)
	offset += 4

	if offsetIn+offset+4 > elemLen {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	elementVersion := *(*int32)(unsafe.Pointer(&currBuff[offset]))
	logger.DEBUG("elementVersion: %d", elementVersion)
	offset += 4

	if offsetIn+offset+4 > elemLen {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	fieldNum := *(*uint32)(unsafe.Pointer(&currBuff[offset]))
	logger.DEBUG("fieldNum: %d", fieldNum)
	offset += 4

	//FieldsBuffLen, 这里这个值没啥用，跳过去了
	if offsetIn+offset+4 > elemLen {
		return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
	}
	offset += 4

	for idx := uint32(0); idx < fieldNum; idx++ {
		nameLen := bytes.IndexByte(currBuff[offset:], 0)
		if nameLen < 0 {
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		nameLen += 1
		if offsetIn+offset+uint32(nameLen) > elemLen {
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		filedLen := *(*uint32)(unsafe.Pointer(&currBuff[offset+uint32(nameLen)]))

		if offsetIn+offset+uint32(nameLen)+4+filedLen > elemLen {
			return &terror.ErrorCode{Code: terror.RecordUnpackFailed}
		}
		name := string(currBuff[offset : offset+uint32(nameLen)-1])
		valueMap[name] = currBuff[offset+uint32(nameLen)+4 : offset+uint32(nameLen)+4+uint32(filedLen)]
		offset = offset + uint32(nameLen) + 4 + uint32(filedLen)
		//logger.ERR("name:[%s], len:%d",name, len(name))
	}
	*offsetOut = offset
	return nil
}
