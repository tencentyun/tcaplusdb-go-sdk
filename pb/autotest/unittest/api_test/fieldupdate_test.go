package api

import (
	"fmt"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplusservice"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestFieldUpdateSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldUpdateReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"

	oldData.PlayerEmail = "wang"
	client.Insert(oldData)
	newData := &tcaplusservice.GamePlayers{PlayerId: 233, PlayerName: "jiahua", PlayerEmail: "wang"}
	client.Get(newData)
	if !proto.Equal(oldData, newData) {
		t.Errorf("data diff")
		return
	}
	defer func() {
		oldData.PlayerEmail = "wang"
		client.Delete(oldData)
	}()

	oldJson := tools.StToJson(oldData)
	fmt.Println(oldJson)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	oldData.Pay = &tcaplusservice.Payment{Amount: 3, PayId: 1}
	if _, err := rec.SetPBFieldValues(oldData, []string{"pay.pay_id", "pay.amount"}); err != nil {
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

	if 1 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newMsg := &tcaplusservice.GamePlayers{}
		err = record.GetPBFieldValues(newMsg)
		if err != nil {
			t.Errorf("GetPBData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newMsg)
		fmt.Println(newJson)

		if !proto.Equal(oldData, newMsg) {
			t.Errorf("resData != reqData")
			return
		}
	}
}

