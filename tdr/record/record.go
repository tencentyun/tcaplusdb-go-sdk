package record

import (
	"bytes"
	"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"sort"
	"sync"
)

/*
	注意：用户请通过函数接口进行record的操作
	设置record有两套接口，切记不可混用：
	1. setKey setValue接口设置的数据，只能通过getKey，getValue接口读取
	2. setData接口设置的数据，只能通过getData读取
*/
type Record struct {
	AppId             uint64
	ZoneId            uint32
	TableName         string
	Cmd               int
	KeyMap            map[string][]byte
	ValueMap          map[string][]byte
	Version           int32
	Index             int32
	KeySet            *tcaplus_protocol_cs.TCaplusKeySet
	ValueSet          *tcaplus_protocol_cs.TCaplusValueSet_
	NameSet           *tcaplus_protocol_cs.TCaplusNameSet
	UpdFieldSet       *tcaplus_protocol_cs.TCaplusUpdFieldSet
	SplitTableKeyBuff *tcaplus_protocol_cs.SplitTableKeyBuff
	PBValueSet        *tcaplus_protocol_cs.ProtobufValueSet_
	PBFieldMap        map[string]bool
	ShardingKey       *[]byte
	ShardingKeyLen    *uint32
	IsPB              bool
}

//record缓存池
var recordPool = sync.Pool{
	New: func() interface{} {
		return new(Record)
	},
}

func GetPoolRecord() *Record {
	r := recordPool.Get().(*Record)
	r.AppId = 0
	r.ZoneId = 0
	r.TableName = ""
	r.Cmd = 0
	r.KeyMap = nil
	r.ValueMap = nil
	r.Version = -1
	r.KeySet = nil
	r.ValueSet = nil
	r.UpdFieldSet = nil
	r.SplitTableKeyBuff = nil
	r.PBValueSet = nil
	r.PBFieldMap = nil
	r.ShardingKey = nil
	r.ShardingKeyLen = nil
	r.IsPB = false
	return r
}

func PutPoolRecord(record *Record) {
	if nil == record {
		return
	}
	recordPool.Put(record)
}

/**
    @brief  设置记录版本号
    @param [IN] v     数据记录的版本号:  <=0 表示不关注版本号不关心版本号。具体含义如下。
                当CHECKDATAVERSION_AUTOINCREASE时: 表示检测记录版本号。
					如果Version的值<=0,则仍然表示不关心版本号不关注版本号；
					如果Version的值>0，那么只有当该版本号与服务器端的版本号相同时，
					Replace, Update, Increase, ListAddAfter, ListDelete, ListReplace, ListDeleteBatch操作才会成功同时在服务器端该版本号会自增1。
                当NOCHECKDATAVERSION_OVERWRITE时: 表示不检测记录版本号。
					如果Version的值<=0,则会把版本号1写入服务端的数据记录版本号(服务器端成功写入的数据记录的版本号最少为1)；
					如果Version的值>0，那么会把该版本号写入服务端的数据记录版本号。
                当NOCHECKDATAVERSION_AUTOINCREASE时: 表示不检测记录版本号，将服务器端的数据记录版本号自增1，若服务器端新写入数据记录则新写入的数据记录的版本号为1。
**/
func (r *Record) SetVersion(v int32) {
	if v <= 0 {
		r.Version = -1
		return
	}
	r.Version = v
}

/**
	@brief  获取记录版本号
	@retval 记录版本号
	@note 对于Generic操作表示获取Record的版本；对于List操作表示获取Record所在List的版本。
**/
func (r *Record) GetVersion() int32 {
	return r.Version
}

func (r *Record) GetIndex() int32 {
	return r.Index
}

func (r *Record) GetAllKeyBlob() (string, error) {
	if len(r.KeyMap) <= 0 {
		return "", &terror.ErrorCode{Code: terror.RequestHasNoKeyField}
	}

	keys := make([]string, len(r.KeyMap))
	i := 0
	buffLen := 0
	for k, v := range r.KeyMap {
		keys[i] = k
		i++
		buffLen += len(v)
	}
	sort.Strings(keys)

	buffLen += 4
	buf := new(bytes.Buffer)
	buf.Grow(buffLen)
	buf.Reset()
	for _, k := range keys {
		v := r.KeyMap[k]
		if v != nil {
			buf.Write(v)
		}
	}
	binary.Write(buf, binary.LittleEndian, r.Index)
	return common.Bytes2str(buf.Bytes()), nil
}
