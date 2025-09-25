package subscribe

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/config"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"sync"
	"sync/atomic"
	"time"
)

type CheckProxyFunc func(proxy interface{}) bool

type Manager struct {
	topicSeq      int64
	topicMapMutex sync.RWMutex
	topicMap      map[int64]*Topic
	topicKeyMap   map[string]*Topic
	option        config.SubscribeOption
}

func (m *Manager) Init(clientOption *config.ClientCtrl) {
	m.topicMap = make(map[int64]*Topic)
	m.topicKeyMap = make(map[string]*Topic)
	if clientOption != nil && clientOption.Option != nil {
		m.option = clientOption.Option.SubscribeOption
	}
	if m.option.ReSubscribeIntervalTime <= 0 {
		m.option.ReSubscribeIntervalTime = 50 * time.Millisecond
	}
}

func (m *Manager) GetReSubscribeIntervalTime() time.Duration {
	return m.option.ReSubscribeIntervalTime
}

func (m *Manager) FindTopicById(id int64) *Topic {
	m.topicMapMutex.RLock()
	defer m.topicMapMutex.RUnlock()
	if v, exist := m.topicMap[id]; exist {
		return v
	}
	return nil
}

func (m *Manager) DelSubscribeTopic(topic *Topic) {
	topicKeyStr, _ := topic.GetKeyStr()
	m.topicMapMutex.Lock()
	defer m.topicMapMutex.Unlock()
	if _, exist := m.topicMap[topic.asyncId]; exist {
		close(topic.ResponsePipe)
	}
	delete(m.topicMap, topic.asyncId)
	delete(m.topicKeyMap, topicKeyStr)
}

func (m *Manager) NewSubscribeTopic(zoneId uint32, tableName string, keyMap map[string][]byte,
	index int32, expireSecond int32, rspPipeSize int, svr interface{}, flag int, isPb bool) (*Topic, error) {
	if len(keyMap) == 0 {
		return nil, terror.MakeError(terror.ParameterInvalid, "key not exist")
	}
	if rspPipeSize <= 0 || flag == OnlyIndexFlag {
		rspPipeSize = 1
	}
	intervalMs := expireSecond * 1000 / 3

	// 计算topic key
	topic := &Topic{
		TopicKey: TopicKey{
			zoneId:    zoneId,
			tableName: tableName,
			keyMap:    nil,
		},
		asyncId:         0,
		isPb:            isPb,
		expireIndex:     index,
		lastChangeIdxMs: 0,
		expireSecond:    expireSecond,
		retryFlag:       0,
		intervalMs:      intervalMs,
		rspFlag:         flag,
		firstIndex:      -1,
		lastIndex:       -1,
		ResponsePipe:    make(chan response.TcaplusResponse, rspPipeSize),
	}
	topic.DeepCopy4KeyMap(keyMap)
	topicKeyStr, err := topic.GetKeyStr()
	if err != nil {
		return nil, err
	}

	// check exist
	m.topicMapMutex.Lock()
	defer m.topicMapMutex.Unlock()
	if _, exist := m.topicKeyMap[topicKeyStr]; exist {
		return nil, terror.MakeError(terror.ParameterInvalid, fmt.Sprintf("topic %s already exist", topic.GetKeyDebug()))
	}

	// not exist
	topic.asyncId = atomic.AddInt64(&m.topicSeq, 1)
	m.topicMap[topic.asyncId] = topic
	m.topicKeyMap[topicKeyStr] = topic
	return topic, nil
}

func (m *Manager) GetTopicSvr(t *Topic) interface{} {
	return t.svr.Load()
}

func (m *Manager) SetTopicSvr(t *Topic, svr interface{}) {
	t.svr.Store(svr)
}

func (m *Manager) SetTopicExpireIndex(t *Topic, index int32) {
	atomic.StoreInt32(&t.expireIndex, index)
}

func (m *Manager) SetTopicSendMs(t *Topic, curMs int64) {
	atomic.StoreInt64(&t.lastSendMs, curMs)
}

func (m *Manager) SetTopicRecvMs(t *Topic, curMs int64) {
	atomic.StoreInt64(&t.lastRecvMs, curMs)
}

func (m *Manager) SetTopicLastChangeIdxMs(t *Topic, curMs int64) {
	atomic.StoreInt64(&t.lastChangeIdxMs, curMs)
}

func (m *Manager) SetTopicCheckFlag(t *Topic, flag int32) {
	atomic.StoreInt32(&t.retryFlag, flag)
}

func (m *Manager) SetTopicFirstIndex(t *Topic, index int32) {
	atomic.StoreInt32(&t.firstIndex, index)
}

func (m *Manager) SetTopicLastIndex(t *Topic, index int32) {
	atomic.StoreInt32(&t.lastIndex, index)
}

func (m *Manager) IncTopicAsyncId(t *Topic) {
	m.topicMapMutex.Lock()
	defer m.topicMapMutex.Unlock()
	delete(m.topicMap, t.asyncId)
	newId := atomic.AddInt64(&m.topicSeq, 1)
	atomic.StoreInt64(&t.asyncId, newId)
	m.topicMap[newId] = t
}

func (m *Manager) GetNeedRetryTopic(CheckFunc CheckProxyFunc) []*Topic {
	topicSlice := make([]*Topic, 0)
	curMs := time.Now().UnixMilli()
	noResTopicSlice := make([]*Topic, 0)

	// notice lock
	m.topicMapMutex.RLock()
	for _, topic := range m.topicKeyMap {
		if topic == nil {
			continue
		}

		if topic.RspPipeMaybeFull() {
			// 本地队列满的暂时跳过，直到不满时再次续订
			continue
		}

		// 设置了标记的
		if topic.GetCheckFlag() != RetrySubscribeWithNoIndex {
			topicSlice = append(topicSlice, topic)
			continue
		}

		// proxy 不可用
		//if svr := m.GetTopicSvr(topic); svr != nil {
		//	if CheckFunc != nil && !CheckFunc(svr) {
		//		topicSlice = append(topicSlice, topic)
		//		noResTopicSlice = append(noResTopicSlice, topic)
		//		m.SetTopicLastChangeIdxMs(topic, curMs)
		//		continue
		//	}
		//}

		// 到达定时时间的
		if curMs > topic.GetLastSendMs()+int64(topic.intervalMs) {
			// 长时间没有响应包的，特殊处理
			if curMs > topic.GetLastChangeIdxMs()+int64(topic.expireSecond*1000) && topic.GetExpireIndex() >= 0 {
				//noResTopicSlice = append(noResTopicSlice, topic)
				m.SetTopicLastChangeIdxMs(topic, curMs)
				continue
			}

			// 普通续订
			topicSlice = append(topicSlice, topic)
			continue
		}
	}
	m.topicMapMutex.RUnlock()
	if len(noResTopicSlice) == 0 {
		return topicSlice
	}

	// 长时间无包的重新订阅
	m.topicMapMutex.Lock()
	for _, topic := range noResTopicSlice {
		// 更新asyncid
		delete(m.topicMap, topic.asyncId)
		newId := atomic.AddInt64(&m.topicSeq, 1)
		atomic.StoreInt64(&topic.asyncId, newId)
		// index只check
		m.SetTopicCheckFlag(topic, RetrySubscribeWithIndexOnlyCheck)
		m.topicMap[newId] = topic
	}
	m.topicMapMutex.Unlock()
	topicSlice = append(topicSlice, noResTopicSlice...)
	return topicSlice
}
