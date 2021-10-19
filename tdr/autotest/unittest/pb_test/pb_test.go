package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestPBSimple(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	ret = client.Get(newData)
	if ret != nil {
		t.Errorf("Get failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.Delete(newData)
}

//get记录不存在的时候 insert 不存在的记录
func TestPBSimple_NonExist(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 111
	newData.PlayerName = "jiahua1"
	newData.PlayerEmail = "dsf1"
	ret = client.Get(newData)
	if ret == nil {
		t.Errorf("Get failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson == newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.Delete(newData)
}

//get 某个key字段不存在
func TestPBGet_Key_NonExist(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 111
	newData.PlayerName = "jiahua1"
	//不存在的key
	newData.PlayerEmail = "dsf111111111"
	ret = client.Get(newData)
	if ret.(*terror.ErrorCode).Code != 261 {
		t.Errorf("Get failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson == newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.Delete(newData)
}

func TestPBBatchGet(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "zhang"
	client.Delete(oldData)

	oldData2 := &tcaplusservice.GamePlayers{}
	oldData2.PlayerId = 234
	oldData2.PlayerName = "jiahua"
	oldData2.PlayerEmail = "zhang"
	client.Delete(oldData2)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}

	data, _ := proto.Marshal(oldData)
	logger.DEBUG("%+v-%d", data, len(data))

	oldJson := tools.StToJson(oldData)
	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}
	client.Insert(oldData2)

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 233
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "zhang"
	ret = client.BatchGet([]proto.Message{newData, oldData2})
	if ret != nil {
		t.Errorf("Get failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.Delete(newData)
	client.Delete(oldData2)
}

//Batch get记录不存在的时候
func TestPBBatchGet_NonExist(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "zhang"
	client.Delete(oldData)

	oldData2 := &tcaplusservice.GamePlayers{}
	oldData2.PlayerId = 234
	oldData2.PlayerName = "jiahua"
	oldData2.PlayerEmail = "zhang"
	client.Delete(oldData2)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}

	data, _ := proto.Marshal(oldData)
	logger.DEBUG("%+v-%d", data, len(data))

	oldJson := tools.StToJson(oldData)
	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}
	client.Insert(oldData2)

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 23322
	newData.PlayerName = "jiahua222"
	newData.PlayerEmail = "zhang222"
	ret = client.BatchGet([]proto.Message{newData, oldData2})
	if ret.(*terror.ErrorCode).Code != 261 {
		t.Errorf("Get failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson == newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}

	client.Delete(newData)
	client.Delete(oldData2)

}

func TestPBGetByPartKey(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "wang"
	client.Delete(oldData)

	oldData2 := &tcaplusservice.GamePlayers{}
	oldData2.PlayerId = 233
	oldData2.PlayerName = "jiahua"
	oldData2.PlayerEmail = "zhang"
	client.Delete(oldData2)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}

	data, _ := proto.Marshal(oldData)
	logger.DEBUG("%+v-%d", data, len(data))

	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}
	client.Insert(oldData2)

	oldData2.PlayerEmail = "li"
	client.Insert(oldData2)

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 233
	newData.PlayerName = "jiahua"
	keys := []string{"player_id", "player_name"}
	msgs, err := client.GetByPartKey(newData, keys)
	if err != nil {
		t.Errorf("Get failed %s", err)
		return
	}
	if len(msgs) != 3 {
		t.Errorf("data len %d", len(msgs))
		return
	}

	client.Delete(newData)
	client.Delete(oldData2)
}

//getbypartkey key不存在的时候
func TestPBGetByPartKey_NonExist(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "wang"
	client.Delete(oldData)

	oldData2 := &tcaplusservice.GamePlayers{}
	oldData2.PlayerId = 233
	oldData2.PlayerName = "jiahua"
	oldData2.PlayerEmail = "zhang"
	client.Delete(oldData2)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}

	data, _ := proto.Marshal(oldData)
	logger.DEBUG("%+v-%d", data, len(data))

	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}
	client.Insert(oldData2)

	oldData2.PlayerEmail = "li"
	client.Delete(oldData2)
	client.Insert(oldData2)

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 2333311111
	newData.PlayerName = "jiahua333"
	keys := []string{"player_id", "player_name"}
	msgs, err := client.GetByPartKey(newData, keys)
	if err.(*terror.ErrorCode).Code != 261 {
		t.Errorf("getbpartkey failed %s", err)
		return
	}
	if len(msgs) != 0 {
		t.Errorf("data len %d", len(msgs))
		return
	}

	client.Delete(newData)
	client.Delete(oldData2)
}

func TestPBIndexQuery_succ(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"

	oldData.PlayerEmail = "wang"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "wang"
		client.Delete(oldData)
	}()

	oldData.PlayerEmail = "zhang"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "zhang"
		client.Delete(oldData)
	}()

	oldData.PlayerEmail = "li"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "li"
		client.Delete(oldData)
	}()

	time.Sleep(time.Second)

	query := fmt.Sprintf("select pay.pay_id, pay.amount from game_players where player_id=233")
	msgs, res, err := client.IndexQuery(query)
	if err != nil {
		t.Errorf("IndexQuery err:%s", err)
		return
	}

	if len(msgs) != 3 && len(res) != 0 {
		t.Errorf("IndexQuery err:%s", err)
		return
	}

	query = fmt.Sprintf("select count(*) from game_players where player_id=233")
	msgs, res, err = client.IndexQuery(query)
	if err != nil {
		t.Errorf("IndexQuery err:%s", err)
		return
	}

	if len(msgs) != 0 && len(res) != 1 {
		t.Errorf("IndexQuery err:%s", err)
		return
	}

}

func TestPBFieldGet(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "zhang"
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	//oldJson := tools.StToJson(oldData)
	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 233
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "zhang"
	newData.Pay = &tcaplusservice.Payment{Amount: 3, PayId: 2, Method: 1}

	err := client.FieldIncrease(newData, []string{"pay.pay_id"})
	if err != nil {
		t.Errorf("Insert failed %d", err)
		return
	}

	logger.DEBUG("%+v", newData)
}

//insert 存在的记录
func TestPBInsert_Exist(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}
	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	ret = client.Get(newData)
	if ret != nil {
		t.Errorf("Get failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}
	//insert 已经存在的数据
	ret1 := client.Insert(oldData)
	if ret1 == nil {
		t.Errorf("Insert failed %d", ret)
		return
	}

	client.Delete(newData)

}

//update 记录存在的时候
func TestPBUpdate_Exist(t *testing.T) {

	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}
	//update 操作
	oldData1 := &tcaplusservice.GamePlayers{}
	oldData1.PlayerId = 444
	oldData1.PlayerName = "jiahua"
	oldData1.PlayerEmail = "dsf"
	oldData1.GameServerId = 10000
	oldData1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson1 := tools.StToJson(oldData1)
	ret1 := client.Update(oldData1)
	if ret1 != nil {
		t.Errorf("Update failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	ret = client.Get(newData)
	if ret != nil {
		t.Errorf("Get failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson1 != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson1, newJson)
		return
	}
	client.Delete(newData)
}

//update 记录不存在的时候
func TestPBUpdate_NoneExist(t *testing.T) {

	client := tools.InitPBSyncClient()
	//update 操作
	oldData1 := &tcaplusservice.GamePlayers{}
	oldData1.PlayerId = 4444444
	oldData1.PlayerName = "jiahua44444"
	oldData1.PlayerEmail = "dsf4444444"
	oldData1.GameServerId = 10000
	oldData1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret := client.Update(oldData1)
	if ret.(*terror.ErrorCode).Code != 261 {
		t.Errorf("Update failed %d", ret)
		return
	}
}

//replace 记录不存在的时候
func TestPBReplace_NonExist(t *testing.T) {
	client := tools.InitPBSyncClient()
	//update 操作
	oldData1 := &tcaplusservice.GamePlayers{}
	oldData1.PlayerId = 444
	oldData1.PlayerName = "jiahua"
	oldData1.PlayerEmail = "dsf"
	oldData1.GameServerId = 10000
	client.Delete(oldData1)
	oldData1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData1)
	ret := client.Replace(oldData1)
	if ret != nil {
		t.Errorf("Update failed %d", ret)
		return
	}
	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	ret = client.Get(newData)
	if ret != nil {
		t.Errorf("Get failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson, newJson)
		return
	}
	client.Delete(newData)

}

//replace 记录存在的时候
func TestPBReplace_Exist(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret := client.Insert(oldData)
	if ret != nil {
		t.Errorf("Insert failed %d", ret)
		return
	}
	//replace 操作
	oldData1 := &tcaplusservice.GamePlayers{}
	oldData1.PlayerId = 444
	oldData1.PlayerName = "jiahua"
	oldData1.PlayerEmail = "dsf"
	oldData1.GameServerId = 10000
	oldData1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson1 := tools.StToJson(oldData1)
	ret1 := client.Replace(oldData1)
	if ret1 != nil {
		t.Errorf("Update failed %d", ret)
		return
	}

	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	ret = client.Get(newData)
	if ret != nil {
		t.Errorf("Get failed %d", ret)
		return
	}
	newJson := tools.StToJson(newData)
	if oldJson1 != newJson {
		t.Errorf("data diff \n%s \n%s", oldJson1, newJson)
		return
	}
	client.Delete(newData)
}

func TestPBListSimple(t *testing.T) {
	client := tools.InitPBSyncClient()

	msg := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}
	client.ListDeleteAll(msg)

	err := client.ListAddAfter(msg, -1)
	if err != nil {
		t.Errorf("Insert failed %s", err)
		return
	}

	err = client.ListAddAfter(msg, -1)
	if err != nil {
		t.Errorf("Insert failed %s", err)
		return
	}

	msg.Gamesvrid = ""
	err = client.ListGet(msg, 0)
	if err != nil {
		t.Errorf("Insert failed %s", err)
		return
	}
	if msg.Gamesvrid != "lol" {
		t.Errorf("Gamesvrid %s != lol", msg.Gamesvrid)
		return
	}

	msgs, err := client.ListGetAll(msg)
	if err != nil {
		t.Errorf("Insert failed %s", err)
		return
	}
	fmt.Println(msgs)

	msgs, err = client.ListDeleteBatch(msg, []int32{0, 1})
	if err != nil {
		t.Errorf("Insert failed %s", err)
		return
	}
	fmt.Println(msgs)

	err = client.ListAddAfter(msg, -1)
	if err != nil {
		t.Errorf("Insert failed %s", err)
		return
	}

	msg.Gamesvrid = "cf"
	err = client.ListReplace(msg, 0)
	if err != nil {
		t.Errorf("Insert failed %s", err)
		return
	}
	msg.Gamesvrid = ""
	err = client.ListGet(msg, 0)
	if err != nil {
		t.Errorf("Insert failed %s", err)
		return
	}
	if msg.Gamesvrid != "cf" {
		t.Errorf("Gamesvrid %s != lol", msg.Gamesvrid)
		return
	}

	msg.Gamesvrid = ""
	err = client.ListDelete(msg, 0)
	if err != nil {
		t.Errorf("Insert failed %s", err)
		return
	}
	err = client.ListGet(msg, 0)
	if err.(*terror.ErrorCode).Code != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("Insert failed %s", err)
		return
	}

}
