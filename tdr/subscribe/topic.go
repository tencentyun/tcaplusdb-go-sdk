package subscribe

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

const RetrySubscribeWithNoIndex = 0        // 不携带index续订阅
const RetrySubscribeWithIndex = 1          // 携带index续订阅，index不存在返回错误码
const RetrySubscribeWithIndexOnlyCheck = 2 // 携带index续订阅，index不存在不返回错误

const (
	RecordFlag    = iota // 0 订阅记录, 一个响应中可能有多条记录，最大256Kb
	OnlyIndexFlag        // 只订阅index
	OneRecordFlag        // 控制每个订阅响应中只有一个记录，保证队列可控
)

type TopicKey struct {
	zoneId    uint32
	tableName string            //表名
	keyMap    map[string][]byte //订阅的List Key列表
}

func (tk *TopicKey) GetZoneId() uint32 {
	return tk.zoneId
}

func (tk *TopicKey) GetTableName() string {
	return tk.tableName
}

func (tk *TopicKey) Copy2KeyMap(keyMap map[string][]byte) {
	if keyMap == nil {
		return
	}
	for k, v := range tk.keyMap {
		keyMap[k] = v
	}
}

func (tk *TopicKey) DeepCopy4KeyMap(keyMap map[string][]byte) {
	tk.keyMap = make(map[string][]byte, len(keyMap))
	for k, v := range keyMap {
		newSlice := make([]byte, len(v))
		copy(newSlice, v)
		tk.keyMap[k] = newSlice
	}
}

func (tk *TopicKey) GetKeyDebug() string {
	return fmt.Sprintf("%d|%s|%+v", tk.zoneId, tk.tableName, tk.keyMap)
}

func (tk *TopicKey) GetKeyStr() (string, error) {
	if len(tk.keyMap) == 0 {
		return "", terror.ErrorCode{Code: terror.ParameterInvalid, Message: "TopicKey keyMap empty"}
	}
	keys := make([]string, 0, len(tk.keyMap))
	for k := range tk.keyMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// zone | table | Key
	topicKey := strconv.FormatUint(uint64(tk.zoneId), 10) + "|" + tk.tableName
	for _, k := range keys {
		topicKey = topicKey + "|" + k                    // kName
		topicKey = topicKey + ":" + string(tk.keyMap[k]) //kBuf
	}
	return topicKey, nil
}

type Topic struct {
	TopicKey                     // 固定
	asyncId         int64        // 订阅使用的异步ID, 唯一但是，丢包时会更新
	isPb            bool         // 是否是pb表订阅,固定
	expireIndex     int32        // 期待下个包的index，收包更新
	lastChangeIdxMs int64        // 最近一次收到订阅响应，并修改idx的时间
	expireSecond    int32        // 订阅的超时时间，秒，固定
	svr             atomic.Value // 订阅发送的svr，断开会更新
	lastSendMs      int64        // 最近一次发送订阅的时间，发请求更新
	lastRecvMs      int64        // 最近一次收到订阅响应的时间，收响应更新
	intervalMs      int32        // 多久续订阅一次，固定
	retryFlag       int32        // 重订阅标记：RetrySubscribeWithNoIndex/RetrySubscribeWithIndex/RetrySubscribeWithIndexOnlyCheck
	rspFlag         int          //RecordFlag/OnlyIndexFlag
	firstIndex      int32        //OnlyIndexFlag时收到第一个订阅包的起始index
	lastIndex       int32        //lastIndex时收到最新订阅包的起始index
	// 响应消息channel,可直接监听（直接监听则不要使用PeekRsp方法）
	// 注意检查是否关闭，关闭了表示订阅断开
	// 如果pipe满，则会丢掉订阅包;等到不满时，会拉取后续订阅包，填充队列，会保证收到的订阅包仍然是连续的
	ResponsePipe chan response.TcaplusResponse
	// 以下用于peek和pop方法
	peekMutex sync.RWMutex             // peek锁，如果不想使用peek，则直接监听ResponsePipe，性能更高
	peekRsp   response.TcaplusResponse // 用来临时存放peek的响应, 用于peek和pop方法配合使用
}

func (t *Topic) GetRspFlag() int {
	return t.rspFlag
}

func (t *Topic) GetFirstIndex() int32 {
	return atomic.LoadInt32(&t.firstIndex)
}

func (t *Topic) GetLastIndex() int32 {
	return atomic.LoadInt32(&t.lastIndex)
}

func (t *Topic) RspPipeMaybeFull() bool {
	return len(t.ResponsePipe) >= cap(t.ResponsePipe)
}

// GetRspPipeLen 获取队列的缓存的响应数量
func (t *Topic) GetRspPipeLen() int {
	t.peekMutex.RLock()
	defer t.peekMutex.RUnlock()
	if t.peekRsp != nil {
		return len(t.ResponsePipe) + 1
	}
	return len(t.ResponsePipe)
}

// PeekRsp 获取topic缓存的头部响应，PeekRsp + PopRsp 配合使用
func (t *Topic) PeekRsp() (response.TcaplusResponse, error) {
	// 先尝试读取数据
	t.peekMutex.RLock()
	rsp := t.peekRsp
	if rsp != nil {
		t.peekMutex.RUnlock()
		return response.NewResponse(rsp.GetTcaplusPackagePtr())
	}
	t.peekMutex.RUnlock()

	// 数据不存在时，获取写锁，避免并发写入
	t.peekMutex.Lock()
	defer t.peekMutex.Unlock()

	// 再次检查 rsp 是否已被更新（防止其他 goroutine 更新了数据）
	if rsp = t.peekRsp; rsp != nil {
		return response.NewResponse(rsp.GetTcaplusPackagePtr())
	}

	// 从管道获取数据
	select {
	case rsp, ok := <-t.ResponsePipe:
		if !ok {
			return nil, terror.MakeError(terror.SubscribeChannelClosed, "")
		}
		// 更新peekRsp
		t.peekRsp = rsp
		return response.NewResponse(rsp.GetTcaplusPackagePtr())
	default:
		// 如果没有数据可用，返回 nil
		return nil, nil
	}
}

// PopRsp 从topic缓存中弹出一个响应
func (t *Topic) PopRsp() (response.TcaplusResponse, error) {
	t.peekMutex.Lock()
	defer t.peekMutex.Unlock()

	rsp := t.peekRsp
	if rsp != nil {
		t.peekRsp = nil
		return rsp, nil
	}

	select {
	case rsp, ok := <-t.ResponsePipe:
		if !ok {
			return nil, terror.MakeError(terror.SubscribeChannelClosed, "")
		}
		return rsp, nil
	default:
		// 如果没有数据可用，返回 nil
		return nil, nil
	}
}

func (t *Topic) GetAsyncId() int64 {
	return t.asyncId
}

func (t *Topic) IsPb() bool {
	return t.isPb
}

func (t *Topic) GetCheckFlag() int32 {
	return atomic.LoadInt32(&t.retryFlag)
}

func (t *Topic) GetExpireIndex() int32 {
	return atomic.LoadInt32(&t.expireIndex)
}

func (t *Topic) GetExpireSec() int32 {
	return t.expireSecond
}

func (t *Topic) GetLastRecvMs() int64 {
	return atomic.LoadInt64(&t.lastRecvMs)
}

func (t *Topic) GetLastSendMs() int64 {
	return atomic.LoadInt64(&t.lastSendMs)
}

func (t *Topic) GetLastChangeIdxMs() int64 {
	return atomic.LoadInt64(&t.lastChangeIdxMs)
}
