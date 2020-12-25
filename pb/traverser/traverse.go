package traverser

import (
	log "github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"sync/atomic"
)

var id uint32

type Traverser struct {
	state     int
	zoneId    uint32
	tableName string
	tableType int

	busy bool
	next bool

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


	client ClientInf
}

func newTraverser(zoneId uint32, table string) *Traverser {
	t := &Traverser{
		state:        TraverseStateReady,
		zoneId:       zoneId,
		tableName:    table,
		beginIndex:   -1,
		endIndex:     -1,
		resNumPerReq: 1,
		limit:        -1,
		traverseId:   atomic.AddUint32(&id, 1),
	}
	return t
}

func (t *Traverser) Start() error {
	if t.state != TraverseStateReady {
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.state = TraverseStateNormal
	return t.sendGetShardListRequest()
}

func (t *Traverser) Stop() error {
	t.state = TraverseStateStop
	t.zoneId = 0
	t.tableName = ""
	return nil
}

func (t *Traverser) Resume() error {
	if t.state != TraverseStateRecoverable {
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.state = TraverseStateNormal
	return nil
}

func (t *Traverser) State() int {
	return t.state
}

func (t *Traverser) SetAsyncId(id uint64) error {
	if TraverseStateReady != t.state {
		log.ERR("Traverser state %d not ready", t.state)
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.asyncId = id
	return nil
}

func (t *Traverser) SetOnlyReadFromSlave(flag bool) error {
	if TraverseStateReady != t.state {
		log.ERR("Traverser state %d not ready", t.state)
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.readFromSlave = flag
	return nil
}

func (t *Traverser) SetUserBuff(buf []byte) error {
	if TraverseStateReady != t.state {
		log.ERR("Traverser state %d not ready", t.state)
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.userBuff = buf
	return nil
}

func (t *Traverser) SetLimit(limit int64) error {
	if TraverseStateReady != t.state {
		log.ERR("Traverser state %d not ready", t.state)
		return &terror.ErrorCode{Code: terror.API_ERR_INVALID_OBJ_STATUE}
	}
	t.limit = limit
	if limit < 0 {
		t.limit = -1
	}

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
	if TraverseStateNormal != t.state {
		log.ERR("Traverser state %d not normal", t.state)
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
		t.shardCompleted = 0
		t.offset = 0
		t.shardCurId++
	}
	req.GetTcaplusPackagePtr().Head.RouterInfo.ShardID = t.shardList[t.shardCurId]

	if 0 == t.asyncId {
		req.GetTcaplusPackagePtr().Head.AsynID = uint64(t.traverseId) << 32 | uint64(t.requestId)
	} else {
		req.GetTcaplusPackagePtr().Head.AsynID = t.asyncId
	}

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

	err = t.client.SendRequest(req)
	if err != nil {
		log.ERR("%s", err)
		return err
	}

	t.seq = p.Sequence
	t.expectReceiveSeq = t.seq + 1

	log.DEBUG("send a traverse request successfully, seq=%d, expectReceiveSeq=%d, resNumPerReq=%d",
		t.seq, t.expectReceiveSeq, t.resNumPerReq)

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
			if !t.busy {
				next = true
			}
		} else {
			shardNum := msg.Body.GetShardListRes.ShardNum
			if shardNum <= 0 || int64(shardNum) >= tcaplus_protocol_cs.TCAPLUS_MAX_SHARD_ID_PER_TABLE {
				log.ERR("zone:%d tableName:%s invalid shard list returned: cnt %d",
					t.zoneId, t.tableName, shardNum)

				t.state = TraverseStateUnRecoverable
				return &terror.ErrorCode{Code: terror.API_ERR_INVALID_SHARD_LIST}
			}

			t.shardCnt = shardNum
			t.shardList = msg.Body.GetShardListRes.ShardList
			t.routeKeySet = msg.Body.GetShardListRes.RouteKeySet
			next = true
			t.busy = true
		}
	} else if cmd.TcaplusApiTableTraverseRes == msg.Head.Cmd {
		result := int(msg.Body.TableTraverseRes.Result)
		if 0 != result {
			log.ERR("TcaplusApiTableTraverse error %d, %s", result, terror.GetErrMsg(result))
			t.state = TraverseStateRecoverable
			return &terror.ErrorCode{Code: result}
		}

		if 0 != t.checkIfSwitchMS(msg.Body.TableTraverseRes.CurSrvID) {
			log.ERR("M and S has switch, set state ST_UNRECOVERABLE")
			t.state = TraverseStateUnRecoverable
			return &terror.ErrorCode{Code: terror.GEN_ERR_ERR}
		}

		receivedSeq := msg.Body.TableTraverseRes.Sequence
		if t.expectReceiveSeq < receivedSeq {
			*drop = true
			log.ERR("zone:%d table:%s receive unexpected pkg, received_seq:%d, m_expect_receive_seq:%d",
				t.zoneId, t.tableName, receivedSeq, t.expectReceiveSeq)
		} else if t.expectReceiveSeq > receivedSeq {
			*drop = true
			log.ERR("zone:%d table:%s receive timeout pkg, received_seq:%d, m_expect_receive_seq:%d",
				t.zoneId, t.tableName, receivedSeq, t.expectReceiveSeq)
		} else {
			t.expectReceiveSeq++
			t.offset = msg.Body.TableTraverseRes.Offset
			t.shardCompleted = msg.Body.TableTraverseRes.Completed
			t.traversedCnt = msg.Body.TableTraverseRes.TraversedCnt

			if t.shardCompleted > 0 {
				if t.shardCurId < t.shardCnt {
					next = true
				}
			} else {
				if msg.Body.TableTraverseRes.Sequence == t.seq + uint64(t.resNumPerReq) {
					next = true
				}
			}

			if 0 == msg.Body.TableTraverseRes.RecordNum {
				*drop = true
				log.DEBUG("zone %d, table %s traverse finished with dwRecordNum = 0 " +
					"on shard(%d/%d) m_shard_completed %d, so this response will be dropped",
					t.zoneId, t.tableName, t.shardCurId+1, t.shardCnt, t.shardCompleted)
			} else {
				log.DEBUG("saved traverse response state: offset %d, " +
					"shard %d, completed %d, recnum %d, total %d, zone %d, table_name %s",
					t.offset, t.shardList[t.shardCurId], t.shardCompleted,
					msg.Body.TableTraverseRes.RecordNum, t.traversedCnt, t.zoneId, t.tableName)
			}
		}
	} else {
		*drop = true
		log.ERR("unexpected command %d", msg.Head.Cmd)
		return nil
	}

	finish := false
	if t.shardCompleted > 0 && t.shardCurId >= t.shardCnt - 1 {
		finish = true
	} else {
		// generic 表
		if 0 == t.tableType {
			if t.limit > 0 && t.traversedCnt >= t.limit {
				finish = true
			}
		}
	}

	if finish {
		log.DEBUG("zone %d, table %s traverse all completed", t.zoneId, t.tableName)
		t.state = TraverseStateIdle
		return nil
	}

	t.next = next
	if next {
		t.busy = true
	}

	return nil
}

func (t *Traverser) continueTraverse() error {
	if !t.busy {
		return nil
	}

	if !t.next {
		return nil
	}

	t.state = TraverseStateNormal

	var err error
	if 0 == t.tableType {
		err = t.sendTraverseRequest()
	}

	if err != nil {
		return err
	}

	t.next = false
	return nil
}
