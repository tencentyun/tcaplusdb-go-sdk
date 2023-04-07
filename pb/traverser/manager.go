package traverser

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"sync"
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

func (m *TraverserManager) GetTraverser(zoneId uint32, table string) *Traverser {
	m.lock.Lock()
	defer m.lock.Unlock()
	zoneTable := fmt.Sprintf("%d|%s", zoneId, table)
	t, exist := m.traverseMap[zoneTable]
	if exist {
		t.tableType = 0
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
	return t
}

func (m *TraverserManager) GetListTraverser(zoneId uint32, table string) *Traverser {
	t := m.GetTraverser(zoneId, table)
	if t != nil {
		t.tableType = 1
	}
	return t
}

func (m *TraverserManager) OnRecvResponse(zoneId uint32, msg *tcaplus_protocol_cs.TCaplusPkg, drop *bool) error {
	if msg == nil || msg.Head == nil {
		logger.ERR("msg invalid")
		return &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID}
	}
	table := string(msg.Head.RouterInfo.TableName[:msg.Head.RouterInfo.TableNameLen-1])
	zoneTable := fmt.Sprintf("%d|%s", zoneId, table)
	m.lock.RLock()
	defer m.lock.RUnlock()
	t, exist := m.traverseMap[zoneTable]
	if !exist {
		logger.ERR("traverse %s not find", zoneTable)
		return &terror.ErrorCode{Code: terror.API_ERR_TRAVERSER_IS_NOT_EXIST}
	}

	if TraverseStateNormal != t.state {
		logger.ERR("Traverser %s state %d not normal", zoneTable, t.state)
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}

	if cmd.TcaplusApiTableTraverseRes == msg.Head.Cmd || cmd.TcaplusApiListTableTraverseRes == msg.Head.Cmd {
		asyncId := t.asyncId
		if 0 == t.asyncId {
			asyncId = uint64(t.traverseId)<<32 | uint64(t.requestId)
		}

		if asyncId != msg.Head.AsynID {
			logger.WARN("zone %d, tableName %s traverse recvived expire response cmd:%d.asyncId %d msg.Head.AsynID %d",
				t.zoneId, t.tableName, msg.Head.Cmd, asyncId, msg.Head.AsynID)
			return nil
		}
	}

	return t.onRecvResponse(msg, drop)
}

func (m *TraverserManager) ContinueTraverse() {
	if len(m.traverseMap) == 0 {
		return
	}

	m.lock.RLock()
	defer m.lock.RUnlock()
	for k, v := range m.traverseMap {
		if v.busy && TraverseStateNormal == v.state {
			err := v.continueTraverse()
			if err != nil {
				logger.ERR("continueTraverse %s error %s", k, err)
			}
		}
	}
}
