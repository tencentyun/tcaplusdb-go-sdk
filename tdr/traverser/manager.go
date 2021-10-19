package traverser

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
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
	lock        sync.Mutex
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
		return t
	}
	if len(m.traverseMap) >= 8 {
		logger.ERR("Traverser map is full")
		return nil
	}
	t = newTraverser(zoneId, table)
	t.client = m.client
	m.traverseMap[zoneTable] = t
	return t
}

func (m *TraverserManager) OnRecvResponse(zoneId uint32, msg *tcaplus_protocol_cs.TCaplusPkg, drop *bool) error {
	if msg == nil || msg.Head == nil {
		logger.ERR("msg invalid")
		return &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID}
	}
	table := string(msg.Head.RouterInfo.TableName[:msg.Head.RouterInfo.TableNameLen-1])
	zoneTable := fmt.Sprintf("%d|%s", zoneId, table)

	t, exist := m.traverseMap[zoneTable]
	if !exist {
		logger.ERR("traverse %s not find", zoneTable)
		return &terror.ErrorCode{Code: terror.API_ERR_TRAVERSER_IS_NOT_EXIST}
	}

	if TraverseStateNormal != t.state {
		logger.ERR("Traverser %s state %d not normal", zoneTable, t.state)
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}

	if cmd.TcaplusApiTableTraverseRes == msg.Head.Cmd {
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

	m.lock.Lock()
	defer m.lock.Unlock()
	for k, v := range m.traverseMap {
		if TraverseStateStop == v.state {
			logger.DEBUG("zoneTable %s traverse stop", k)
			delete(m.traverseMap, k)
			continue
		}
		if v.busy && TraverseStateNormal == v.state {
			err := v.continueTraverse()
			if err != nil {
				logger.ERR("continueTraverse %s error %s", k, err)
			}
		}
	}
}
