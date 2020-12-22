package api

import (
	"fmt"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplusservice"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/traverser"
	"testing"
)

func TestPBTraverse(t *testing.T) {
	client := tools.InitPBSyncClient()
	tra := client.GetTraverser(tools.ZoneId, "game_players")
	err := tra.Start()
	if err != nil {
		t.Errorf("start error:%s", err)
	}

	msg := &tcaplusservice.GamePlayers{}
	msg.PlayerId = 233
	msg.PlayerName = "TestPBTraverse"
	msg.PlayerEmail = "dsf"
	for i:=0;i<500;i++{
		msg.PlayerId = 233 * int64(i)
		client.Insert(msg)
		defer func(id int64) {
			msg.PlayerId = 233 * id
			client.Delete(msg)
		}(int64(i))
	}

	count := 0

	for {
		//recv resp
		resp, err := tools.RecvResponse(client)
		if err != nil {
			if tra.State() == traverser.TraverseStateIdle {
				break
			}
			t.Errorf("recvResponse fail, %s", err.Error())
			return
		}

		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
			return
		}

		count += resp.GetRecordCount()

		for i := 0; i < resp.GetRecordCount(); i++ {
			record, err := resp.FetchRecord()
			if err != nil {
				t.Errorf("FetchRecord failed %s", err.Error())
				return
			}

			newMsg := &tcaplusservice.GamePlayers{}
			err = record.GetPBData(newMsg)
			if err != nil {
				t.Errorf("GetPBData failed %s", err.Error())
				return
			}

			newJson := tools.StToJson(newMsg)
			fmt.Println(newJson)
		}
	}

	fmt.Println(count)
}
