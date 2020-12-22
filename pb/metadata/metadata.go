// @Title  metadata.go
// @Description  解析 tcaplus proxy 返回的 protobuf 文件元数据
// @Author  jiahuazhang  2020-12-07 10:20:00
package metadata

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"sort"
	"strings"
	"sync"
)

// 单例
var ins *metaManager

// 用于实现单例
var once sync.Once

// 元数据管理
type metaManager struct {
	// 读写锁，由于map非线程安全的，操作时上锁防止panic
	lock   sync.RWMutex

	// protobuf表元数据 map key 为  zoneID|tableName
	tables map[string]*msgDescGroup
}

// 元数据结构体
type msgDescGroup struct {
	// table 元数据
	Desc        protoreflect.MessageDescriptor

	// primary key 列表
	Keys        []string

	// shardingKey 用户指定，或者计算索引交集得到
	ShardingKey string

	// 必填字段列表
	Required    []string

	// 全字段  key为字段路径 a.b.c
	FieldMap	map[string][]protowire.Number

	// key  为字段number路径 1.1.1
	NumberMap	map[string][]protowire.Number

	// 检查元数据，仅检查一次，所以不要多个版本同时用
	Checked		bool
}

// 获取元数据管理（单例）
func GetMetaManager() *metaManager {
	once.Do(func() {
		if ins == nil {
			ins = &metaManager{tables: make(map[string]*msgDescGroup)}
		}
	})
	return ins
}

// @title    AddTableDesGrp
// @description  填加表元数据
// @param     appId   uint64  表所在的 app ID
// @param     zoneId   uint32  表所在的 zone ID
// @param     tableName  string  表名
// @param     metaData  []byte   元数据
// @return    error  错误信息
func (m *metaManager) AddTableDesGrp(appId uint64, zoneId uint32, tableName string, metaData []byte) error {
	group := &msgDescGroup{
		FieldMap: make(map[string][]protowire.Number),
		NumberMap: make(map[string][]protowire.Number),
	}

	// 解析元数据
	files, err := m.DescriptorPool(metaData)
	if err != nil {
		logger.ERR("DescriptorPool error:%s", err)
		return &terror.ErrorCode{Code: terror.API_ERR_UNPACK_MESSAGE}
	}

	// 获取表描述符
	files.RangeFiles(func(f protoreflect.FileDescriptor) bool {
		if f.Name() == "google/protobuf/descriptor.proto" ||
			f.Name() == "google/protobuf/any.proto" ||
			f.Name() == "tcaplusservice.optionv1.proto" {
			return true
		}
		group.Desc = f.Messages().ByName(protoreflect.Name(tableName))
		if group.Desc != nil {
			return false
		}
		return true
	})
	if group.Desc == nil {
		errMsg := "not find table message descriptor"
		logger.ERR(errMsg)
		return &terror.ErrorCode{Code: terror.API_ERR_UNPACK_MESSAGE, Message: errMsg}
	}

	// 获取表primarykey
	opts := group.Desc.Options()
	var unknowMap map[protowire.Number]interface{}
	strPrimary := proto.GetExtension(opts, tcaplusservice.E_TcaplusPrimaryKey).(string)
	if len(strPrimary) == 0 {
		// 如果用户未初始化 option 文件可能出现不识别的情况，从unknown中获取
		unknowMap = m.unKnownField(opts.ProtoReflect().GetUnknown())
		if value, exist := unknowMap[60000]; exist {
			strPrimary = value.(string)
		} else {
			errMsg := fmt.Sprintf("table %s not find primarykey", tableName)
			logger.ERR(errMsg)
			return &terror.ErrorCode{Code: terror.API_ERR_MISS_PRIMARY_KEY, Message: errMsg}
		}
	}
	kmap := make(map[string]struct{})
	keys := strings.Split(strPrimary, ",")
	for _, k := range keys {
		key := strings.TrimSpace(k)
		kmap[key] = struct{}{}
		group.Keys = append(group.Keys, key)
	}
	// 如果有重复，去重
	if len(kmap) != len(keys) {
		group.Keys = group.Keys[:0]
		for k := range kmap {
			group.Keys = append(group.Keys, k)
		}
	}

	// 获取ShardingKey
	group.ShardingKey = proto.GetExtension(opts, tcaplusservice.E_TcaplusShardingKey).(string)
	if len(unknowMap) != 0 {
		group.ShardingKey = unknowMap[60005].(string)
	}

	// 获取required字段
	required := make(map[string]struct{})
	m.getAllRequiredFields(group.Desc, "", "", nil, required, group.FieldMap, group.NumberMap)
	group.Required = make([]string, 0, len(required))
	for r := range required {
		group.Required = append(group.Required, r)
	}

	logger.DEBUG("init msgDescGroup:\n%+v", group)

	m.lock.Lock()
	m.tables[fmt.Sprintf("%d|%d|%s", appId, zoneId, tableName)] = group
	m.lock.Unlock()
	return nil
}

func (m *metaManager) DescriptorPool(buf []byte) (*protoregistry.Files, error) {
	files := &descriptorpb.FileDescriptorSet{}
	err := proto.Unmarshal(buf, files)
	if err != nil {
		return nil, err
	}
	pbdesc := protodesc.ToFileDescriptorProto(descriptorpb.File_google_protobuf_descriptor_proto)
	pbany := protodesc.ToFileDescriptorProto(anypb.File_google_protobuf_any_proto)
	files.File = append(files.File, pbdesc, pbany)

	return protodesc.NewFiles(files)
}

func (m *metaManager) GetTableDesGrp(zoneTable string) *msgDescGroup {
	m.lock.RLock()
	group, ok := m.tables[zoneTable]
	m.lock.RUnlock()
	if ok {
		return group
	}
	return nil
}

func (m *metaManager) getAllRequiredFields(desc protoreflect.MessageDescriptor, prefixName, prefixNum string,
		prefixNumArr []protowire.Number, requird map[string]struct{},
		fieldMap map[string][]protowire.Number, NumberMap map[string][]protowire.Number) error {
	if desc == nil {
		return fmt.Errorf("DescriptorProto is nil")
	}
	if requird == nil {
		return fmt.Errorf("required map is nil")
	}

	fields := desc.Fields()
	for i := 0; i < fields.Len(); i++ {
		f := desc.Fields().Get(i)
		fname := prefixName
		fnum := prefixNum
		numberArr := make([]protowire.Number, len(prefixNumArr))
		copy(numberArr, prefixNumArr)
		if prefixName == "" {
			fname = string(f.Name())
		} else {
			fname += "." + string(f.Name())
		}
		if prefixNum == "" {
			fnum = fmt.Sprint(f.Number())
		} else {
			fnum += "." + fmt.Sprint(f.Number())
		}
		numberArr = append(numberArr, f.Number())
		fieldMap[fname] = numberArr
		NumberMap[fnum] = numberArr
		if descriptorpb.FieldDescriptorProto_Label(f.Cardinality()) ==
			descriptorpb.FieldDescriptorProto_LABEL_REQUIRED {
			requird[fname] = struct{}{}
		}
		if f.Message() != nil {
			ret := m.getAllRequiredFields(f.Message(), fname, fnum, numberArr, requird, fieldMap, NumberMap)
			if ret != nil {
				return ret
			}
		}
	}

	return nil
}

func (m *metaManager) ExtractMsgPartKey(message proto.Message, keys []string) ([]byte, error) {
	var keybuf []byte
	desc := message.ProtoReflect().Descriptor()
	var err error

	// 去重
	common := make(map[string]struct{})
	for _, k := range keys {
		common[strings.TrimSpace(k)] = struct{}{}
	}

	fields := make([]protoreflect.FieldDescriptor, 0, len(common))
	for name := range common {
		field := desc.Fields().ByName(protoreflect.Name(name))
		if field == nil {
			errMsg := fmt.Sprintf("message not find primary key field %s", name)
			logger.ERR(errMsg)
			return keybuf, terror.ErrorCode{Code: terror.API_ERR_PACK_MESSAGE, Message: errMsg}
		}
		fields = append(fields, field)
	}

	// 按number排序
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Number() < fields[j].Number()
	})

	// 序列化
	for _, f := range fields {
		keybuf, err = proto.MarshalOptions{}.MarshalField(keybuf, f, message.ProtoReflect().Get(f))
		if err != nil {
			errMsg := fmt.Sprintf("MarshalField field %s error:%s", f.JSONName(), err)
			logger.ERR(errMsg)
			return keybuf, terror.ErrorCode{Code: terror.API_ERR_PACK_MESSAGE, Message: errMsg}
		}
	}
	return keybuf, nil
}

func (m *metaManager) CompareMessageMeta(svr, cli protoreflect.MessageDescriptor) error {
	if svr == nil || cli == nil {
		errMsg := fmt.Sprintf("message descriptor is nil")
		logger.ERR(errMsg)
		return fmt.Errorf(errMsg)
	}

	if svr.Syntax() != cli.Syntax() {
		errMsg := fmt.Sprintf("syntax is diff svr-cli:%s-%s", svr.Syntax(), cli.Syntax())
		logger.ERR(errMsg)
		return fmt.Errorf(errMsg)
	}

	strSvrCipherSuite := proto.GetExtension(svr.Options(), tcaplusservice.E_TcaplusFieldCipherSuite).(string)
	strCliCipherSuite := proto.GetExtension(cli.Options(), tcaplusservice.E_TcaplusFieldCipherSuite).(string)
	if strSvrCipherSuite != strCliCipherSuite {
		errMsg := fmt.Sprintf("tcaplus_field_cipher_suite is diff svr-cli:%s-%s",
			strSvrCipherSuite, strCliCipherSuite)
		logger.ERR(errMsg)
		return fmt.Errorf(errMsg)
	}

	strSvrCipherMd5 := proto.GetExtension(svr.Options(), tcaplusservice.E_TcaplusCipherMd5).(string)
	strCliCipherMd5 := proto.GetExtension(cli.Options(), tcaplusservice.E_TcaplusCipherMd5).(string)
	if strSvrCipherMd5 != strCliCipherMd5 {
		errMsg := fmt.Sprintf("tcaplus_cipher_md5 is diff svr-cli:%s-%s",
			strSvrCipherMd5, strCliCipherMd5)
		logger.ERR(errMsg)
		return fmt.Errorf(errMsg)
	}

	for i := 0; i < svr.Fields().Len(); i++ {
		svrField := svr.Fields().Get(i)
		cliField := cli.Fields().Get(i)
		if svrField == nil || cliField == nil {
			errMsg := fmt.Sprintf("message field is nil")
			logger.ERR(errMsg)
			return fmt.Errorf(errMsg)
		}
		bSvrCrypto := proto.GetExtension(svrField.Options(), tcaplusservice.E_TcaplusCrypto).(bool)
		bCliCrypto := proto.GetExtension(svrField.Options(), tcaplusservice.E_TcaplusCrypto).(bool)
		svrLabel := descriptorpb.FieldDescriptorProto_Label(svrField.Cardinality())
		cliLabel := descriptorpb.FieldDescriptorProto_Label(cliField.Cardinality())
		svrType := descriptorpb.FieldDescriptorProto_Type(svrField.Kind())
		cliType := descriptorpb.FieldDescriptorProto_Type(cliField.Kind())

		if svrLabel != descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL &&
			svrLabel != descriptorpb.FieldDescriptorProto_LABEL_REPEATED &&
			svrLabel != cliLabel {
			errMsg := fmt.Sprintf("field %s lable is diff svr-cli:%d-%d", svrField.Name(), svrLabel, cliLabel)
			logger.ERR(errMsg)
			return fmt.Errorf(errMsg)
		}

		if svrField.Name() != cliField.Name() {
			errMsg := fmt.Sprintf("field name is diff svr-cli:%s-%s", svrField.Name(), cliField.Name())
			logger.ERR(errMsg)
			return fmt.Errorf(errMsg)
		}

		if svrType != cliType {
			errMsg := fmt.Sprintf("field %s lable is diff svr-cli:%d-%d", svrField.Name(), svrType, cliType)
			logger.ERR(errMsg)
			return fmt.Errorf(errMsg)
		}

		if bSvrCrypto != bCliCrypto {
			errMsg := fmt.Sprintf("field %s crypto is diff svr-cli:%v-%v", svrField.Name(), bSvrCrypto, bCliCrypto)
			logger.ERR(errMsg)
			return fmt.Errorf(errMsg)
		}

		if svrType == descriptorpb.FieldDescriptorProto_TYPE_GROUP ||
			svrType == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {

			err := m.CompareMessageMeta(svrField.Message(), cliField.Message())
			if err != nil {
				return err
			}
		}
	}

	iSvrIndexs := proto.GetExtension(svr.Options(), tcaplusservice.E_TcaplusIndex).([]string)
	iCliIndexs := proto.GetExtension(cli.Options(), tcaplusservice.E_TcaplusIndex).([]string)
	if len(iSvrIndexs) != len(iCliIndexs) {
		errMsg := fmt.Sprintf("indexs count is diff svr-cli:%d-%d", len(iSvrIndexs), len(iCliIndexs))
		logger.ERR(errMsg)
		return fmt.Errorf(errMsg)
	}

	if len(iSvrIndexs) > 0 {
		strSvrShardingKey := proto.GetExtension(svr.Options(), tcaplusservice.E_TcaplusShardingKey).(string)
		strCliShardingKey := proto.GetExtension(cli.Options(), tcaplusservice.E_TcaplusShardingKey).(string)

		if strCliShardingKey != "" && strSvrShardingKey != strCliShardingKey {
			errMsg := fmt.Sprintf("shardingkey is diff svr-cli:%s-%s", strSvrShardingKey, strCliShardingKey)
			logger.ERR(errMsg)
			return fmt.Errorf(errMsg)
		}

		svrIndexInfos := m.formatIndexInfo(iSvrIndexs)
		if svrIndexInfos == nil {
			errMsg := fmt.Sprintf("svr indexs format error")
			logger.ERR(errMsg)
			return fmt.Errorf(errMsg)
		}

		cliIndexInfos := m.formatIndexInfo(iCliIndexs)
		if cliIndexInfos == nil {
			errMsg := fmt.Sprintf("svr indexs format error")
			logger.ERR(errMsg)
			return fmt.Errorf(errMsg)
		}

		for k, v := range svrIndexInfos {
			value, exist := cliIndexInfos[k]
			if !exist || !m.isEqual(v, value) {
				errMsg := fmt.Sprintf("index %s is diff svr-cli:%+v-%+v", k, v, value)
				logger.ERR(errMsg)
				return fmt.Errorf(errMsg)
			}
		}
	}

	return nil
}

func (m *metaManager) formatIndexInfo(infos []string) map[string][]string {
	if len(infos) == 0 {
		return nil
	}
	infomap := make(map[string][]string, len(infos))
	for _, info := range infos {
		left := strings.Index(info, "(")
		right := strings.Index(info, ")")
		keys := strings.Split(string([]byte(info)[left+1:right]), ",")
		if left == -1 || right == -1 || len(keys) == 0 {
			logger.ERR("index %s format error", info)
			return nil
		}
		name := strings.TrimSpace(string([]byte(info)[:left]))
		infomap[name] = keys
	}
	return infomap
}

func (m *metaManager) isEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Slice(a, func(i, j int) bool {
		return strings.TrimSpace(a[i]) < strings.TrimSpace(a[j])
	})
	sort.Slice(b, func(i, j int) bool {
		return strings.TrimSpace(b[i]) < strings.TrimSpace(b[j])
	})
	for i := 0; i < len(a); i++ {
		if strings.TrimSpace(a[i]) != strings.TrimSpace(b[i]) {
			return false
		}
	}
	return true
}

func (m *metaManager) unKnownField(b protoreflect.RawFields) map[protowire.Number]interface{} {
	unknownMap := make(map[protowire.Number]interface{})
	for len(b) > 0 {
		num, wtyp, n := protowire.ConsumeTag(b)
		if n < 0 {
			return unknownMap
		}
		if protowire.BytesType != wtyp {
			return unknownMap
		}
		b = b[n:]

		v, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return unknownMap
		}
		b = b[n:]
		if num == 60001 {
			if value, exist := unknownMap[num]; exist {
				unknownMap[num] = append(value.([]string), string(v))
			} else {
				unknownMap[num] = []string{string(v)}
			}
		} else {
			unknownMap[num] = string(v)
		}
	}
	return unknownMap
}
