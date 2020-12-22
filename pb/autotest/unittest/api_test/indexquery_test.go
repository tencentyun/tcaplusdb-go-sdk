package api

import (
	"fmt"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplusservice"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
	"testing"
	"time"
)

func TestPBIndexQuerySuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "game_players")

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

	query := fmt.Sprintf("select * from game_players where player_id=233 and player_name=jiahua")
	req.SetSql(query)

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}

	if 3 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 3", resp.GetRecordCount())
		return
	}

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

		if newMsg.PlayerEmail != "zhang" &&
			newMsg.PlayerEmail != "wang" &&
			newMsg.PlayerEmail != "li" {
			t.Errorf("resData != reqData")
			return
		}
	}
}

func TestPBIndexQueryFieldSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "game_players")

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

	query := fmt.Sprintf("select pay.amount, pay.pay_id from game_players where player_id=233 and player_name=jiahua")
	req.SetSql(query)

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}

	if 3 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 3", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newMsg := &tcaplusservice.GamePlayers{}
		err = record.GetPBDataWithValues(newMsg, []string{"pay.amount", "pay.pay_id"})
		if err != nil {
			t.Errorf("GetPBData failed %s", err.Error())
			return
		}

		if newMsg.Pay.PayId != oldData.Pay.PayId ||
			newMsg.Pay.Amount != oldData.Pay.Amount {
			t.Errorf("resData != reqData")
			return
		}
	}
}
