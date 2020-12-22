package api

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"testing"
)

//case1 记录存在时，Get返回成功
func TestPBGetBypartkeySuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetByPartkeyReq, "game_players")

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

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBPartKeys(oldData, []string{"player_id", "player_name"}); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

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
