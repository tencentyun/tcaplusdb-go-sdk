package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"testing"
	"time"
)

func TestPBTraverse(t *testing.T) {
	client := tools.InitPBSyncClient()
	tra := client.GetTraverser(cfg.ApiConfig.ZoneId, "game_players")
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
