package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestPBSimpleDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.DoDelete(oldData, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	opt := &option.PBOpt{}
	ret = client.DoGet(newData, opt)
	if ret != nil {
		t.Errorf("DoGet failed %d", ret)
		return
	}

	if opt.Version <= 0 {
		t.Errorf("DoGet opt.Version %d", opt.Version)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.DoDelete(newData, nil)
}

//DoGet记录不存在的时候 DoInsert 不存在的记录
func TestPBSimple_NonExistDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.DoDelete(oldData, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 111
	newData.PlayerName = "jiahua1"
	newData.PlayerEmail = "dsf1"
	ret = client.DoGet(newData, nil)
	if ret == nil {
		t.Errorf("DoGet failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson == newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.DoDelete(newData, nil)
}

//DoGet 某个key字段不存在
func TestPBDoGet_Key_NonExistDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.DoDelete(oldData, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 111
	newData.PlayerName = "jiahua1"
	//不存在的key
	newData.PlayerEmail = "dsf111111111"
	ret = client.DoGet(newData, nil)
	if ret.(*terror.ErrorCode).Code != 261 {
		t.Errorf("DoGet failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson == newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.DoDelete(newData, nil)
}

func TestPBBatchGetDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "zhang"
	client.DoDelete(oldData, nil)

	oldData2 := &tcaplusservice.GamePlayers{}
	oldData2.PlayerId = 234
	oldData2.PlayerName = "jiahua"
	oldData2.PlayerEmail = "zhang"
	client.DoDelete(oldData2, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}

	data, _ := proto.Marshal(oldData)
	logger.DEBUG("%+v-%d", data, len(data))

	oldJson := tools.StToJson(oldData)
	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}
	client.DoInsert(oldData2, nil)

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 233
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "zhang"
	opt := &option.PBOpt{}
	ret = client.DoBatchGet([]proto.Message{newData, oldData2}, opt)
	if ret != nil {
		t.Errorf("DoGet failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.DoDelete(newData, nil)
	client.DoDelete(oldData2, nil)
}

//Batch DoGet记录不存在的时候
func TestPBBatchGet_NonExistDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "zhang"
	client.DoDelete(oldData, nil)

	oldData2 := &tcaplusservice.GamePlayers{}
	oldData2.PlayerId = 234
	oldData2.PlayerName = "jiahua"
	oldData2.PlayerEmail = "zhang"
	client.DoDelete(oldData2, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}

	data, _ := proto.Marshal(oldData)
	logger.DEBUG("%+v-%d", data, len(data))

	oldJson := tools.StToJson(oldData)
	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}
	client.DoInsert(oldData2, nil)

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 23322
	newData.PlayerName = "jiahua222"
	newData.PlayerEmail = "zhang222"
	opt := &option.PBOpt{}
	ret = client.DoBatchGet([]proto.Message{newData, oldData2}, opt)
	if ret.(*terror.ErrorCode).Code != 261 {
		t.Errorf("DoBatchGet failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson == newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.DoDelete(newData, nil)
	client.DoDelete(oldData2, nil)

}

func TestPBGetByPartKeyDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "wang"
	client.DoDelete(oldData, nil)

	oldData2 := &tcaplusservice.GamePlayers{}
	oldData2.PlayerId = 233
	oldData2.PlayerName = "jiahua"
	oldData2.PlayerEmail = "zhang"
	client.DoDelete(oldData2, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}

	data, _ := proto.Marshal(oldData)
	logger.DEBUG("%+v-%d", data, len(data))

	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}
	client.DoInsert(oldData2, nil)

	oldData2.PlayerEmail = "li"
	client.DoInsert(oldData2, nil)

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 233
	newData.PlayerName = "jiahua"
	keys := []string{"player_id", "player_name"}
	msgs, err := client.DoGetByPartKey(newData, keys, nil)
	if err != nil {
		t.Errorf("DoGet failed %s", err)
		return
	}
	if len(msgs) != 3 {
		t.Errorf("data len %d", len(msgs))
		return
	}

	client.DoDelete(newData, nil)
	client.DoDelete(oldData2, nil)
}

//DoGetbypartkey key不存在的时候
func TestPBDoGetByPartKey_NonExistDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "wang"
	client.DoDelete(oldData, nil)

	oldData2 := &tcaplusservice.GamePlayers{}
	oldData2.PlayerId = 233
	oldData2.PlayerName = "jiahua"
	oldData2.PlayerEmail = "zhang"
	client.DoDelete(oldData2, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}

	data, _ := proto.Marshal(oldData)
	logger.DEBUG("%+v-%d", data, len(data))

	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}
	client.DoInsert(oldData2, nil)

	oldData2.PlayerEmail = "li"
	client.DoDelete(oldData2, nil)
	client.DoInsert(oldData2, nil)

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 2333311111
	newData.PlayerName = "jiahua333"
	keys := []string{"player_id", "player_name"}
	msgs, err := client.DoGetByPartKey(newData, keys, nil)
	if err.(*terror.ErrorCode).Code != 261 {
		t.Errorf("DoGetbpartkey failed %s", err)
		return
	}
	if len(msgs) != 0 {
		t.Errorf("data len %d", len(msgs))
		return
	}

	client.DoDelete(newData, nil)
	client.DoDelete(oldData2, nil)
}

func TestPBFieldDoGet(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "zhang"
	client.DoDelete(oldData, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	//oldJson := tools.StToJson(oldData)
	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 233
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "zhang"
	newData.Pay = &tcaplusservice.Payment{Amount: 3, PayId: 2, Method: 1}

	opt := &option.PBOpt{
		FieldNames: []string{"pay.pay_id"},
	}
	err := client.DoFieldIncrease(newData, opt)
	if err != nil {
		t.Errorf("DoInsert failed %d", err)
		return
	}

	logger.DEBUG("%+v", newData)
}

//DoInsert 存在的记录
func TestPBDoInsert_ExistDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.DoDelete(oldData, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}
	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	ret = client.DoGet(newData, nil)
	if ret != nil {
		t.Errorf("DoGet failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}
	//DoInsert 已经存在的数据
	ret1 := client.DoInsert(oldData, nil)
	if ret1 == nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	client.DoDelete(newData, nil)

}

//DoUpdate 记录存在的时候
func TestPBUpdate_ExistDo(t *testing.T) {

	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.DoDelete(oldData, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}
	//DoUpdate 操作
	oldData1 := &tcaplusservice.GamePlayers{}
	oldData1.PlayerId = 444
	oldData1.PlayerName = "jiahua"
	oldData1.PlayerEmail = "dsf"
	oldData1.GameServerId = 10000
	oldData1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson1 := tools.StToJson(oldData1)
	ret1 := client.DoUpdate(oldData1, nil)
	if ret1 != nil {
		t.Errorf("DoUpdate failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	ret = client.DoGet(newData, nil)
	if ret != nil {
		t.Errorf("DoGet failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson1 != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson1, newJson)
		return
	}
	client.DoDelete(newData, nil)
}

//DoUpdate 记录不存在的时候
func TestPBUpdate_NoneExistDo(t *testing.T) {

	client := tools.InitPBSyncClient()
	//DoUpdate 操作
	oldData1 := &tcaplusservice.GamePlayers{}
	oldData1.PlayerId = 4444444
	oldData1.PlayerName = "jiahua44444"
	oldData1.PlayerEmail = "dsf4444444"
	oldData1.GameServerId = 10000
	oldData1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret := client.DoUpdate(oldData1, nil)
	if ret.(*terror.ErrorCode).Code != 261 {
		t.Errorf("DoUpdate failed %d", ret)
		return
	}
}

//DoReplace 记录不存在的时候
func TestPBReplace_NonExistDo(t *testing.T) {
	client := tools.InitPBSyncClient()
	//DoUpdate 操作
	oldData1 := &tcaplusservice.GamePlayers{}
	oldData1.PlayerId = 444
	oldData1.PlayerName = "jiahua"
	oldData1.PlayerEmail = "dsf"
	oldData1.GameServerId = 10000
	client.DoDelete(oldData1, nil)
	oldData1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData1)
	ret := client.DoReplace(oldData1, nil)
	if ret != nil {
		t.Errorf("DoUpdate failed %d", ret)
		return
	}
	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	ret = client.DoGet(newData, nil)
	if ret != nil {
		t.Errorf("DoGet failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}
	client.DoDelete(newData, nil)

}

//DoReplace 记录存在的时候
func TestPBReplace_ExistDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.DoDelete(oldData, nil)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret := client.DoInsert(oldData, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}
	//DoReplace 操作
	oldData1 := &tcaplusservice.GamePlayers{}
	oldData1.PlayerId = 444
	oldData1.PlayerName = "jiahua"
	oldData1.PlayerEmail = "dsf"
	oldData1.GameServerId = 10000
	oldData1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson1 := tools.StToJson(oldData1)
	ret1 := client.DoReplace(oldData1, nil)
	if ret1 != nil {
		t.Errorf("DoUpdate failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	ret = client.DoGet(newData, nil)
	if ret != nil {
		t.Errorf("DoGet failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson1 != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson1, newJson)
		return
	}
	client.DoDelete(newData, nil)
}

func TestPBListSimpleDo(t *testing.T) {
	client := tools.InitPBSyncClient()

	msg := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}
	client.DoListDeleteAll(msg, nil)

	_, err := client.DoListAddAfter(msg, -1, nil)
	if err != nil {
		t.Errorf("DoInsert failed %s", err)
		return
	}

	_, err = client.DoListAddAfter(msg, -1, nil)
	if err != nil {
		t.Errorf("DoInsert failed %s", err)
		return
	}

	msg.Gamesvrid = ""
	err = client.DoListGet(msg, 0, nil)
	if err != nil {
		t.Errorf("DoInsert failed %s", err)
		return
	}
	if msg.Gamesvrid != "lol" {
		t.Errorf("Gamesvrid %s != lol", msg.Gamesvrid)
		return
	}

	msgs, err := client.DoListGetAll(msg, nil)
	if err != nil {
		t.Errorf("DoInsert failed %s", err)
		return
	}
	fmt.Println(msgs)

	msgs, err = client.DoListDeleteBatch(msg, []int32{0, 1}, nil)
	if err != nil {
		t.Errorf("DoListDeleteBatch failed %s", err)
		return
	}
	fmt.Println(msgs)

	_, err = client.DoListAddAfter(msg, -1, nil)
	if err != nil {
		t.Errorf("DoInsert failed %s", err)
		return
	}

	msg.Gamesvrid = "cf"
	err = client.DoListReplace(msg, 0, nil)
	if err != nil {
		t.Errorf("DoInsert failed %s", err)
		return
	}
	msg.Gamesvrid = ""
	err = client.DoListGet(msg, 0, nil)
	if err != nil {
		t.Errorf("DoInsert failed %s", err)
		return
	}
	if msg.Gamesvrid != "cf" {
		t.Errorf("Gamesvrid %s != lol", msg.Gamesvrid)
		return
	}

	msg.Gamesvrid = ""
	err = client.DoListDelete(msg, 0, nil)
	if err != nil {
		t.Errorf("DoInsert failed %s", err)
		return
	}
	err = client.DoListGet(msg, 0, nil)
	if err.(*terror.ErrorCode).Code != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("DoInsert failed %s", err)
		return
	}

}
