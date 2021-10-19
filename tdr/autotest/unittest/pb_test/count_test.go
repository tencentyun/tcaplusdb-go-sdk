package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"testing"
	"time"
)

func TestCountSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetTableRecordCountReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerName = "TestCountSuccess"
	oldMsg.PlayerEmail = "dsf"
	for i := 0; i < 5; i++ {
		oldMsg.PlayerId = int64(i)
		client.Delete(oldMsg)
	}

	resp, err := client.Do(req, 5*time.Second)
	if err != nil {
		t.Errorf("client.Do error:%s", err)
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult %d != 0", resp.GetResult())
		return
	}

	old := resp.GetTableRecordCount()
	fmt.Println(resp.GetTableRecordCount())
	for i := 0; i < 5; i++ {
		oldMsg.PlayerId = int64(i)
		err = client.Insert(oldMsg)
		if err != nil {
			t.Errorf("client.Insert error:%s", err)
			return
		}
		defer func(id int) {
			oldMsg.PlayerId = int64(id)
			client.Delete(oldMsg)
		}(i)
	}

	resp, err = client.Do(req, 5*time.Second)
	if err != nil {
		t.Errorf("client.Do error:%s", err)
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult %d != 0", resp.GetResult())
		return
	}

	fmt.Println(resp.GetTableRecordCount())

	if resp.GetTableRecordCount() != old+5 {
		t.Errorf("resp.GetTableRecordCount() %d != %d", resp.GetTableRecordCount(), old+5)
		return
	}
}

func TestSyncCountSuccess(t *testing.T) {
	client := tools.InitPBSyncClient()

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerName = "TestSyncCountSuccess"
	oldMsg.PlayerEmail = "dsf"
	for i := 0; i < 5; i++ {
		oldMsg.PlayerId = int64(i)
		client.Delete(oldMsg)
	}

	old, err := client.GetTableCount("game_players")
	if err != nil {
		t.Errorf("client.GetTableCount error:%s", err)
		return
	}

	fmt.Println(old)
	for i := 0; i < 5; i++ {
		oldMsg.PlayerId = int64(i)
		err = client.Insert(oldMsg)
		if err != nil {
			t.Errorf("client.Insert error:%s", err)
			return
		}
		defer func(id int) {
			oldMsg.PlayerId = int64(id)
			client.Delete(oldMsg)
		}(i)
	}

	count, err := client.GetTableCount("game_players")
	if err != nil {
		t.Errorf("client.GetTableCount error:%s", err)
		return
	}

	fmt.Println(count)

	if count != old+5 {
		t.Errorf("resp.GetTableRecordCount() %d != %d", count, old+5)
		return
	}
}
