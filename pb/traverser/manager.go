package traverser

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/request"
	"sync"
	"sync/atomic"
)

const (
	TraverseStateIdle          = 1
	TraverseStateReady         = 2
	TraverseStateNormal        = 4
	TraverseStateStop          = 8
	TraverseStateRecoverable   = 16
	TraverseStateUnRecoverable = 32
)

type ClientInf interface {
	NewRequest(zoneId uint32, tableName string, cmd int) (request.TcaplusRequest, error)
	SendRequest(req request.TcaplusRequest) error
}

type TraverserManager struct {
	lock        sync.RWMutex
	traverseMap map[string]*Traverser
	client      ClientInf
}

func NewTraverserManager(client ClientInf) *TraverserManager {
	tm := &TraverserManager{
		traverseMap: make(map[string]*Traverser, 8),
		client:      client,
	}
	return tm
}

func (m *TraverserManager) CheckTraverserFinish(t *Traverser) bool {
	if t != nil && atomic.LoadInt32(&t.isFinished) == 1 {
		atomic.StoreInt32(&t.state, TraverseStateIdle)
		return true
	}
	return false
}

func (m *TraverserManager) GetTraverser(zoneId uint32, table string, isPb bool) *Traverser {
	m.lock.Lock()
	defer m.lock.Unlock()
	zoneTable := fmt.Sprintf("%d|%s", zoneId, table)
	t, exist := m.traverseMap[zoneTable]
	if exist {
		t.tableType = 0
		t.isPb = isPb
		return t
	}
	if len(m.traverseMap) >= 8 {
		logger.ERR("Traverser map is full")
		return nil
	}
	t = newTraverser(zoneId, table)
	t.tableType = 0
	t.client = m.client
	t.tm = m
	m.traverseMap[zoneTable] = t
	t.isPb = isPb
	return t
}

func (m *TraverserManager) GetListTraverser(zoneId uint32, table string, isPb bool) *Traverser {
	t := m.GetTraverser(zoneId, table, isPb)
	if t != nil {
		t.tableType = 1
		t.isPb = isPb
	}
	return t
}

func (m *TraverserManager) OnRecvResponse(zoneId uint32, msg *tcaplus_protocol_cs.TCaplusPkg, drop *bool) *Traverser {
	if msg == nil || msg.Head == nil {
		logger.ERR("msg invalid")
		*drop = true
		return nil
	}
	table := string(msg.Head.RouterInfo.TableName[:msg.Head.RouterInfo.TableNameLen-1])
	zoneTable := fmt.Sprintf("%d|%s", zoneId, table)
	m.lock.RLock()
	defer m.lock.RUnlock()
	t, exist := m.traverseMap[zoneTable]
	if !exist {
		logger.ERR("traverse %s not find", zoneTable)
		*drop = true
		return nil
	}

	if TraverseStateNormal != atomic.LoadInt32(&t.state) {
		*drop = true
		logger.ERR("Traverser %s state %d not normal", zoneTable, atomic.LoadInt32(&t.state))
		return nil
	}

	if cmd.TcaplusApiTableTraverseRes == msg.Head.Cmd || cmd.TcaplusApiListTableTraverseRes == msg.Head.Cmd {
		asyncId := t.asyncId
		if 0 == t.asyncId {
			asyncId = uint64(t.traverseId)<<32 | uint64(t.requestId)
		}

		if asyncId != msg.Head.AsynID {
			*drop = true
			logger.WARN("zone %d, tableName %s traverse recvived expire response cmd:%d.asyncId %d msg.Head.AsynID %d",
				t.zoneId, t.tableName, msg.Head.Cmd, asyncId, msg.Head.AsynID)
			return nil
		}
	}
	t.onRecvResponse(msg, drop)
	return t
}

func (m *TraverserManager) ContinueTraverse() {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if len(m.traverseMap) == 0 {
		return
	}
	for k, v := range m.traverseMap {
		if v.busy.Load().(bool) && TraverseStateNormal == atomic.LoadInt32(&v.state) {
			err := v.continueTraverse()
			if err != nil {
				logger.ERR("continueTraverse %s error %s", k, err)
			}
		}
	}
}
