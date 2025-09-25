package traverser

import (
	"fmt"
	log "github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"sync/atomic"
)

var id uint32

type Traverser struct {
	state     int32
	zoneId    uint32
	tableName string
	tableType int //0 generic 1 list
	isPb      bool

	busy atomic.Value // bool
	next atomic.Value // bool

	traverseId uint32
	requestId  uint32
	asyncId    uint64

	shardCnt   int32
	shardCurId int32
	shardList  []int32

	expectReceiveSeq uint64
	shardCompleted   int32

	shardCurSvrId string

	userBuff []byte

	readFromSlave bool

	// for traverse request
	offset       uint64
	nameSet      *tcaplus_protocol_cs.TCaplusNameSet
	resNumPerReq uint32
	routeKeySet  *tcaplus_protocol_cs.RouteKeySet
	beginIndex   int32
	endIndex     int32
	seq          uint64
	traversedCnt int64
	limit        int64

	//for list
	keyTraversedCnt int64 //已遍历的key数量
	totalKeyLimit   int64 //遍历的key上限

	seqForSync int32

	client     ClientInf
	tm         *TraverserManager
	condition  string
	isFinished int32
}

func newTraverser(zoneId uint32, table string) *Traverser {
	t := &Traverser{
		state:         TraverseStateReady,
		zoneId:        zoneId,
		tableName:     table,
		beginIndex:    -1,
		endIndex:      -1,
		resNumPerReq:  1,
		limit:         -1,
		totalKeyLimit: -1,
		traverseId:    atomic.AddUint32(&id, 1),
		isFinished:    0,
	}
	t.busy.Store(false)
	t.next.Store(false)
	return t
}

func (t *Traverser) Start() error {
	if atomic.LoadInt32(&t.state) != TraverseStateReady {
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	atomic.StoreInt32(&t.state, TraverseStateNormal)
	atomic.StoreInt32(&t.isFinished, 0)
	return t.sendGetShardListRequest()
}

func (t *Traverser) Stop() error {
	zoneTable := fmt.Sprintf("%d|%s", t.zoneId, t.tableName)
	t.tm.lock.Lock()
	delete(t.tm.traverseMap, zoneTable)
	t.tm.lock.Unlock()
	atomic.StoreInt32(&t.state, TraverseStateStop)
	t.zoneId = 0
	t.tableName = ""
	return nil
}

func (t *Traverser) Resume() error {
	if atomic.LoadInt32(&t.state) != TraverseStateRecoverable {
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	atomic.StoreInt32(&t.state, TraverseStateNormal)
	t.next.Store(true)
	return nil
}

func (t *Traverser) State() int {
	return int(atomic.LoadInt32(&t.state))
}

func (t *Traverser) SetAsyncId(id uint64) error {
	if TraverseStateReady != atomic.LoadInt32(&t.state) {
		log.ERR("Traverser state %d not ready", atomic.LoadInt32(&t.state))
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.asyncId = id
	return nil
}

func (t *Traverser) SetSeq(seq int32) error {
	if TraverseStateReady != atomic.LoadInt32(&t.state) {
		log.ERR("Traverser state %d not ready", atomic.LoadInt32(&t.state))
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.seqForSync = seq
	return nil
}

func (t *Traverser) SetOnlyReadFromSlave(flag bool) error {
	if TraverseStateReady != atomic.LoadInt32(&t.state) {
		log.ERR("Traverser state %d not ready", atomic.LoadInt32(&t.state))
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.readFromSlave = flag
	return nil
}

func (t *Traverser) SetUserBuff(buf []byte) error {
	if TraverseStateReady != atomic.LoadInt32(&t.state) {
		log.ERR("Traverser state %d not ready", atomic.LoadInt32(&t.state))
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.userBuff = buf
	return nil
}

func (t *Traverser) SetLimit(limit int64) error {
	if TraverseStateReady != atomic.LoadInt32(&t.state) {
		log.ERR("Traverser state %d not ready", atomic.LoadInt32(&t.state))
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.limit = limit
	t.totalKeyLimit = limit
	if limit < 0 {
		t.limit = -1
		t.totalKeyLimit = -1
	}

	return nil
}

func (t *Traverser) SetFieldNames(valueNameList []string) error {
	if TraverseStateReady != atomic.LoadInt32(&t.state) {
		log.ERR("Traverser state %d not ready", atomic.LoadInt32(&t.state))
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	if t.nameSet == nil {
		t.nameSet = &tcaplus_protocol_cs.TCaplusNameSet{}
	}
	t.nameSet.FieldNum = 0
	t.nameSet.FieldName = nil
	for _, name := range valueNameList {
		t.nameSet.FieldNum++
		t.nameSet.FieldName = append(t.nameSet.FieldName, name)
	}
	if t.nameSet.FieldNum > 0 && t.isPb {
		t.nameSet.FieldNum++
		t.nameSet.FieldName = append(t.nameSet.FieldName, "$")
	}
	return nil
}

func (t *Traverser) SetCondition(query string) error {
	if TraverseStateReady != atomic.LoadInt32(&t.state) {
		log.ERR("Traverser state %d not ready", atomic.LoadInt32(&t.state))
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	if int64(len(query)) > tcaplus_protocol_cs.TCAPLUS_MAX_EXPR_TEXT_LEN {
		log.ERR("invalid condition.size=%d", len(query))
		return &terror.ErrorCode{Code: terror.GEN_ERR_INVALID_ARGUMENTS,
			Message: fmt.Sprintf("condition len %d over max %d",
				len(query), tcaplus_protocol_cs.TCAPLUS_MAX_EXPR_TEXT_LEN)}
	}
	t.condition = query

	return nil
}

func (t *Traverser) sendGetShardListRequest() error {
	req, err := t.client.NewRequest(t.zoneId, t.tableName, cmd.TcaplusApiGetShardListReq)
	if err != nil {
		log.ERR("zone %d table %s cmd %d NewRequest error:%s",
			t.zoneId, t.tableName, cmd.TcaplusApiGetShardListReq, err)
		return err
	}

	p := req.GetTcaplusPackagePtr().Body.GetShardListReq
	p.BeginIndex = t.beginIndex
	p.EndIndex = t.endIndex

	return t.client.SendRequest(req)
}

func (t *Traverser) sendTraverseRequest() error {
	if TraverseStateNormal != atomic.LoadInt32(&t.state) {
		log.ERR("Traverser state %d not normal", atomic.LoadInt32(&t.state))
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}

	req, err := t.client.NewRequest(t.zoneId, t.tableName, cmd.TcaplusApiTableTraverseReq)
	if err != nil {
		log.ERR("zone %d table %s cmd %d NewRequest error:%s",
			t.zoneId, t.tableName, cmd.TcaplusApiTableTraverseReq, err)
		return err
	}

	if t.readFromSlave {
		req.GetTcaplusPackagePtr().Head.Flags |= int32(tcaplus_protocol_cs.TCAPLUS_FLAG_ONLY_READ_FROM_SLAVE)
	}

	if t.shardCompleted > 0 {
		t.shardCurId++
		if t.shardCurId >= t.shardCnt {
			//shard都遍历完了
			return nil
		}
		t.shardCompleted = 0
		t.offset = 0
	}
	req.GetTcaplusPackagePtr().Head.RouterInfo.ShardID = t.shardList[t.shardCurId]

	if 0 == t.asyncId {
		t.requestId++
		req.GetTcaplusPackagePtr().Head.AsynID = uint64(t.traverseId)<<32 | uint64(t.requestId)
	} else {
		req.GetTcaplusPackagePtr().Head.AsynID = t.asyncId
	}

	req.GetTcaplusPackagePtr().Head.Seq = t.seqForSync

	if len(t.userBuff) != 0 {
		req.GetTcaplusPackagePtr().Head.UserBuff = t.userBuff
		req.GetTcaplusPackagePtr().Head.UserBuffLen = uint32(len(t.userBuff))
	}

	p := req.GetTcaplusPackagePtr().Body.TableTraverseReq
	p.BatchLimit = -1
	p.Offset = t.offset
	p.ResNumPerReq = t.resNumPerReq
	p.RouteKeySet = t.routeKeySet
	p.BeginIndex = t.beginIndex
	p.EndIndex = t.endIndex
	p.Sequence = t.seq + uint64(t.resNumPerReq)
	p.ToTalLimit = t.limit
	p.TraversedCnt = t.traversedCnt
	if t.nameSet != nil {
		p.ValueInfo = t.nameSet
	}
	p.Condition = t.condition

	err = t.client.SendRequest(req)
	if err != nil {
		log.ERR("%s", err)
		return err
	}

	t.seq = p.Sequence
	atomic.StoreUint64(&t.expectReceiveSeq, t.seq+1)

	log.DEBUG("send a traverse request successfully, seq=%d, expectReceiveSeq=%d, resNumPerReq=%d",
		t.seq, atomic.LoadUint64(&t.expectReceiveSeq), t.resNumPerReq)

	return nil
}

func (t *Traverser) sendListTraverseRequest() error {
	if TraverseStateNormal != atomic.LoadInt32(&t.state) {
		log.ERR("Traverser state %d not normal", atomic.LoadInt32(&t.state))
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}

	req, err := t.client.NewRequest(t.zoneId, t.tableName, cmd.TcaplusApiListTableTraverseReq)
	if err != nil {
		log.ERR("zone %d table %s cmd %d NewRequest error:%s",
			t.zoneId, t.tableName, cmd.TcaplusApiListTableTraverseReq, err)
		return err
	}

	if t.readFromSlave {
		req.GetTcaplusPackagePtr().Head.Flags |= int32(tcaplus_protocol_cs.TCAPLUS_FLAG_ONLY_READ_FROM_SLAVE)
	}

	if t.shardCompleted > 0 {
		t.shardCurId++
		if t.shardCurId >= t.shardCnt {
			//shard都遍历完了
			return nil
		}
		t.shardCompleted = 0
		t.offset = 0
	}

	req.GetTcaplusPackagePtr().Head.RouterInfo.ShardID = t.shardList[t.shardCurId]

	if 0 == t.asyncId {
		t.requestId++
		req.GetTcaplusPackagePtr().Head.AsynID = uint64(t.traverseId)<<32 | uint64(t.requestId)
	} else {
		req.GetTcaplusPackagePtr().Head.AsynID = t.asyncId
	}

	req.GetTcaplusPackagePtr().Head.Seq = t.seqForSync

	if len(t.userBuff) != 0 {
		req.GetTcaplusPackagePtr().Head.UserBuff = t.userBuff
		req.GetTcaplusPackagePtr().Head.UserBuffLen = uint32(len(t.userBuff))
	}

	p := req.GetTcaplusPackagePtr().Body.ListTableTraverseReq
	keyLimit := int64(1)
	t.seq = t.seq + uint64(tcaplus_protocol_cs.TCAPLUS_MAX_LIST_ELEMENTS_NUM)

	p.Offset = t.offset
	p.KeyLimit = int32(keyLimit)
	p.RouteKeySet = t.routeKeySet
	p.BeginIndex = t.beginIndex
	p.EndIndex = t.endIndex
	p.Sequence = t.seq
	p.ToTalLimit = t.limit
	p.TraversedCnt = t.traversedCnt
	if t.nameSet != nil {
		p.ValueInfo = t.nameSet
	}
	p.Condition = t.condition

	err = t.client.SendRequest(req)
	if err != nil {
		log.ERR("%s", err)
		return err
	}
	atomic.StoreUint64(&t.expectReceiveSeq, t.seq+1)

	log.DEBUG("send a traverse request successfully, seq=%d, expectReceiveSeq=%d, resNumPerReq=%d",
		t.seq, atomic.LoadUint64(&t.expectReceiveSeq), t.resNumPerReq)

	return nil
}

func (t *Traverser) checkIfSwitchMS(resCurSrvID string) int {
	if t.offset == 0 {
		t.shardCurSvrId = resCurSrvID
	} else {
		if t.shardCurSvrId != resCurSrvID {
			log.ERR("last rsp svrId (%s) not equal current rsp svrId (%s), maybe has Switch Master and Slave",
				t.shardCurSvrId, resCurSrvID)
			return terror.GEN_ERR_ERR
		}
	}
	return 0
}

func (t *Traverser) onRecvResponse(msg *tcaplus_protocol_cs.TCaplusPkg, drop *bool) error {
	// 是否需要发送下一个请求
	next := false

	if cmd.TcaplusApiGetShardListRes == msg.Head.Cmd {
		*drop = true
		if t.shardCnt != 0 || t.shardCurId != 0 {
			//收到重复的GET_SHARD_LIST回包, 忽略
			log.WARN("unexpected GetShardListRes or invalid local shard state, cnt %d, idx %d, zone %d, tableName %s",
				t.shardCnt, t.shardCurId)

			//如果总在Resume,busy又恢复了,给一次发包机会
			if !t.busy.Load().(bool) {
				next = true
			}
		} else {
			shardNum := msg.Body.GetShardListRes.ShardNum
			if shardNum <= 0 || int64(shardNum) >= tcaplus_protocol_cs.TCAPLUS_MAX_SHARD_ID_PER_TABLE ||
				int(shardNum) != len(msg.Body.GetShardListRes.ShardList) {
				log.ERR("zone:%d tableName:%s invalid shard list returned: cnt %d",
					t.zoneId, t.tableName, shardNum)

				atomic.StoreInt32(&t.state, TraverseStateUnRecoverable)
				return &terror.ErrorCode{Code: terror.API_ERR_INVALID_SHARD_LIST}
			}

			t.shardCnt = shardNum
			t.shardList = msg.Body.GetShardListRes.ShardList
			t.routeKeySet = msg.Body.GetShardListRes.RouteKeySet
			next = true
			t.busy.Store(true)
		}
	} else if cmd.TcaplusApiTableTraverseRes == msg.Head.Cmd {
		result := int(msg.Body.TableTraverseRes.Result)
		if 0 != result {
			log.ERR("TcaplusApiTableTraverse error %d, %s", result, terror.GetErrMsg(result))
			atomic.StoreInt32(&t.state, TraverseStateRecoverable)
			return &terror.ErrorCode{Code: result}
		}

		if 0 != t.checkIfSwitchMS(msg.Body.TableTraverseRes.CurSrvID) {
			log.ERR("M and S has switch, set state ST_UNRECOVERABLE")
			atomic.StoreInt32(&t.state, TraverseStateUnRecoverable)
			return &terror.ErrorCode{Code: terror.GEN_ERR_ERR}
		}

		receivedSeq := msg.Body.TableTraverseRes.Sequence
		if atomic.LoadUint64(&t.expectReceiveSeq) < receivedSeq {
			*drop = true
			next = true
			log.ERR("zone:%d table:%s receive unexpected pkg, received_seq:%d, m_expect_receive_seq:%d",
				t.zoneId, t.tableName, receivedSeq, atomic.LoadUint64(&t.expectReceiveSeq))
		} else if atomic.LoadUint64(&t.expectReceiveSeq) > receivedSeq {
			*drop = true
			log.ERR("zone:%d table:%s receive timeout pkg, received_seq:%d, m_expect_receive_seq:%d",
				t.zoneId, t.tableName, receivedSeq, atomic.LoadUint64(&t.expectReceiveSeq))
		} else {
			atomic.AddUint64(&t.expectReceiveSeq, 1)
			t.offset = msg.Body.TableTraverseRes.Offset
			t.shardCompleted = msg.Body.TableTraverseRes.Completed
			t.traversedCnt = msg.Body.TableTraverseRes.TraversedCnt

			if t.shardCompleted > 0 {
				if t.shardCurId < t.shardCnt-1 {
					next = true
				}
			} else {
				if msg.Body.TableTraverseRes.Sequence == t.seq+uint64(t.resNumPerReq) {
					next = true
				}
			}

			if 0 == msg.Body.TableTraverseRes.RecordNum {
				*drop = true
				log.INFO("zone %d, table %s traverse finished with dwRecordNum = 0 "+
					"on shard(%d/%d) m_shard_completed %d, so this response will be dropped",
					t.zoneId, t.tableName, t.shardCurId+1, t.shardCnt, t.shardCompleted)
			} else {
				log.DEBUG("saved traverse response state: offset %d, "+
					"shard %d, completed %d, recnum %d, total %d, zone %d, table_name %s",
					t.offset, t.shardList[t.shardCurId], t.shardCompleted,
					msg.Body.TableTraverseRes.RecordNum, t.traversedCnt, t.zoneId, t.tableName)
			}
		}
	} else if cmd.TcaplusApiListTableTraverseRes == msg.Head.Cmd {
		result := int(msg.Body.ListTableTraverseRes.Result)
		if 0 != result {
			log.ERR("TcaplusApiTableTraverse error %d, %s", result, terror.GetErrMsg(result))
			atomic.StoreInt32(&t.state, TraverseStateRecoverable)
			return &terror.ErrorCode{Code: result}
		}

		if 0 != t.checkIfSwitchMS(msg.Body.ListTableTraverseRes.CurSrvID) {
			log.ERR("M and S has switch, set state ST_UNRECOVERABLE")
			atomic.StoreInt32(&t.state, TraverseStateUnRecoverable)
			return &terror.ErrorCode{Code: terror.GEN_ERR_ERR}
		}

		receivedSeq := msg.Body.ListTableTraverseRes.Sequence
		if atomic.LoadUint64(&t.expectReceiveSeq) < receivedSeq {
			*drop = true
			next = true
			log.ERR("zone:%d table:%s receive unexpected pkg, received_seq:%d, m_expect_receive_seq:%d",
				t.zoneId, t.tableName, receivedSeq, atomic.LoadUint64(&t.expectReceiveSeq))
		} else if atomic.LoadUint64(&t.expectReceiveSeq) > receivedSeq {
			*drop = true
			log.ERR("zone:%d table:%s receive timeout pkg, received_seq:%d, m_expect_receive_seq:%d",
				t.zoneId, t.tableName, receivedSeq, atomic.LoadUint64(&t.expectReceiveSeq))
		} else {
			atomic.AddUint64(&t.expectReceiveSeq, 1)
			t.offset = msg.Body.ListTableTraverseRes.Offset
			t.shardCompleted = msg.Body.ListTableTraverseRes.Completed
			if msg.Body.ListTableTraverseRes.KeyCompleted > 0 && msg.Body.ListTableTraverseRes.RecordNum > 0 {
				t.keyTraversedCnt++
				next = true
			}

			if t.shardCompleted > 0 {
				if t.shardCurId < t.shardCnt-1 {
					next = true
				}
			}

			if 0 == msg.Body.ListTableTraverseRes.RecordNum {
				*drop = true
				log.INFO("zone %d, table %s traverse finished with dwRecordNum = 0 "+
					"on shard(%d/%d) m_shard_completed %d, so this response will be dropped",
					t.zoneId, t.tableName, t.shardCurId+1, t.shardCnt, t.shardCompleted)
				if t.shardCompleted <= 0 {
					next = true
				}
			} else {
				log.DEBUG("saved traverse response state: offset %d, "+
					"shard %d, completed %d, recnum %d, total %d, zone %d, table_name %s",
					t.offset, t.shardList[t.shardCurId], t.shardCompleted,
					msg.Body.ListTableTraverseRes.RecordNum, t.traversedCnt, t.zoneId, t.tableName)
			}
		}
	} else {
		*drop = true
		log.ERR("unexpected command %d", msg.Head.Cmd)
		return nil
	}

	finish := false
	msg.Head.Flags = 1 //复用为finish标记位
	if t.shardCompleted > 0 && t.shardCurId >= t.shardCnt-1 {
		finish = true
	} else {
		if 0 == t.tableType { // generic 表
			if t.limit > 0 && t.traversedCnt >= t.limit {
				finish = true
			}
		} else { //list
			if t.keyTraversedCnt >= t.totalKeyLimit && t.totalKeyLimit > 0 {
				finish = true
			}
		}
	}

	if finish {
		log.INFO("zone %d, table %s traverse all completed", t.zoneId, t.tableName)
		atomic.StoreInt32(&t.isFinished, 1)
		msg.Head.Flags = 0 //复用为finish标记位
		*drop = false
		return nil
	}

	t.next.Store(next)
	if next {
		t.busy.Store(true)
	}

	return nil
}

func (t *Traverser) continueTraverse() error {
	if !t.busy.Load().(bool) {
		return nil
	}

	if !t.next.Load().(bool) {
		return nil
	}

	atomic.StoreInt32(&t.state, TraverseStateNormal)

	var err error
	if 0 == t.tableType {
		err = t.sendTraverseRequest()
	} else {
		err = t.sendListTraverseRequest()
	}

	if err != nil {
		return err
	}

	t.next.Store(false)
	return nil
}
