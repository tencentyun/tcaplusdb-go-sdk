package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"testing"
	"time"
)

func TestPBTraverse(t *testing.T) {
	client := tools.InitPBSyncClient()
	tra := client.GetTraverser(1, "game_players")
	defer tra.Stop()

	msg := &tcaplusservice.GamePlayers{}
	msg.PlayerId = 233
	msg.PlayerName = "TestPBTraverse"
	msg.PlayerEmail = "dsf"
	for i := 0; i < 5; i++ {
		msg.PlayerId = 233 * int64(i)
		client.Insert(msg)
		defer func(id int64) {
			msg.PlayerId = 233 * id
			client.Delete(msg)
		}(int64(i))
	}

	msgs, err := client.DoTraverse(tra, 60*time.Second)
	if err != nil {
		t.Errorf("RecvResponse fail, %s", err.Error())
		return
	}
	fmt.Println(len(msgs))
}

func TestPBSyncTraverse(t *testing.T) {
	client := tools.InitPBSyncClient()

	msg := &tcaplusservice.GamePlayers{}
	msg.PlayerId = 233
	msg.PlayerName = "TestPBSyncTraverse"
	msg.PlayerEmail = "dsf"
	// 插入500条数据
	for i := 0; i < 5; i++ {
		msg.PlayerId = 233 * int64(i)
		client.Insert(msg)
		defer func(id int64) {
			msg.PlayerId = 233 * id
			client.Delete(msg)
		}(int64(i))
	}

	// 遍历，参数为定义的 proto message， 返回  message 列表与错误
	table := &tcaplusservice.GamePlayers{}
	client.SetDefaultTimeOut(30 * time.Second)
	msgs, err := client.Traverse(table)
	if err != nil {
		t.Errorf("start error:%s", err)
	}

	fmt.Println(len(msgs))
}

func TestPBSyncTraverse2(t *testing.T) {
	client := tools.InitPBSyncClient()

	msg := &tcaplusservice.GamePlayers{}
	msg.PlayerId = 233
	msg.PlayerName = "TestPBSyncTraverse2"
	msg.PlayerEmail = "dsf"
	// 插入500条数据
	for i := 0; i < 5; i++ {
		msg.PlayerId = 233 * int64(i)
		client.Insert(msg)
		defer func(id int64) {
			msg.PlayerId = 233 * id
			client.Delete(msg)
		}(int64(i))
	}

	// 遍历，参数为定义的 proto message， 返回  message 列表与错误
	table := &tcaplusservice.GamePlayers{}
	client.SetDefaultTimeOut(30 * time.Second)

	for i := 0; i < 10; i++ {
		start := time.Now()
		msgs, err := client.Traverse(table)
		if err != nil {
			t.Errorf("start error:%s", err)
		}

		fmt.Println(len(msgs))
		fmt.Println(time.Since(start))
		time.Sleep(10 * time.Millisecond)
	}

}

// 条件符合
func TestPBTraverseCondition(t *testing.T) {
	client := tools.InitPBSyncClient()
	tra := client.GetTraverser(1, "user")
	defer tra.Stop()

	msg := &tcaplusservice.User{}
	msg.Id = 1
	msg.Name = "aaa"
	client.Insert(msg)
	defer func() {
		msg.Id = 1
		msg.Name = "aaa"
		client.Delete(msg)
	}()

	msg.Id = 2
	msg.Name = "bbb"
	client.Insert(msg)
	defer func() {
		msg.Id = 2
		msg.Name = "bbb"
		client.Delete(msg)
	}()

	msg.Id = 3
	msg.Name = "ccc"
	client.Insert(msg)
	defer func() {
		msg.Id = 3
		msg.Name = "ccc"
		client.Delete(msg)
	}()

	msg.Id = 4
	msg.Name = "ddd"
	client.Insert(msg)
	defer func() {
		msg.Id = 4
		msg.Name = "ddd"
		client.Delete(msg)
	}()

	msg.Id = 5
	msg.Name = "eee"
	client.Insert(msg)
	defer func() {
		msg.Id = 5
		msg.Name = "eee"
		client.Delete(msg)
	}()

	tra.SetCondition("id > 2 AND name != \"eee\"")

	resps, err := client.DoTraverse(tra, 60*time.Second)
	if err != nil {
		t.Errorf("RecvResponse fail, %s", err.Error())
		return
	}

	for _, resp := range resps {
		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
			return
		}

		for i := 0; i < resp.GetRecordCount(); i++ {
			record, err := resp.FetchRecord()
			if err != nil {
				t.Errorf("FetchRecord failed %s", err.Error())
				return
			}

			newMsg := &tcaplusservice.User{}
			err = record.GetPBData(newMsg)
			if err != nil {
				t.Errorf("GetPBData failed %s", err.Error())
				return
			}

			newJson := tools.StToJson(newMsg)
			fmt.Println(newJson)
		}
	}
}
