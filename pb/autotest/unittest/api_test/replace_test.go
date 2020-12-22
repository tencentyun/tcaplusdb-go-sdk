package api

import (
	"fmt"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplusservice"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
	"testing"
)

//case 1记录不存在时Replace success resultFlag = 2
func TestPBReplaceSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiReplaceReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 233
	oldMsg.PlayerName = "TestPBReplaceSuccess"
	oldMsg.PlayerEmail = "dsf"
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

	client.Delete(oldMsg)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldMsg); err != nil {
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
		err = record.GetPBData(newMsg)
		if err != nil {
			t.Errorf("GetPBData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newMsg)
		fmt.Println(newJson)
		if oldJson != newJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	client.Delete(oldMsg)
}

//case2 记录存在时，重复Replace，返回成功 resultFlag = 2
func TestDupReplace(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiReplaceReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 233
	oldMsg.PlayerName = "TestDupReplace"
	oldMsg.PlayerEmail = "dsf"
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

	client.Insert(oldMsg)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldMsg); err != nil {
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
		err = record.GetPBData(newMsg)
		if err != nil {
			t.Errorf("GetPBData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newMsg)
		fmt.Println(newJson)
		if oldJson != newJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	client.Delete(oldMsg)
}
