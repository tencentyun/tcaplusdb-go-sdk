package api

import (
	"fmt"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplusservice"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/logger"
	"google.golang.org/protobuf/proto"
	"testing"
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

func TestPBIndexQuery(t *testing.T) {
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
