package api

import (
	"fmt"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplusservice"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
	"testing"
)

//case 1记录不存在时update fail
func TestUpdateFail(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiUpdateReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 233
	oldMsg.PlayerName = "TestUpdateFail"
	oldMsg.PlayerEmail = "dsf"
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

	// 万一存在
	client.Delete(oldMsg)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
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

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err ,but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult expect TXHDB_ERR_RECORD_NOT_EXIST ,but %s", terror.GetErrMsg(err))
		return
	}
}

//case2 记录存在时，update返回成功 resultFlag = 2
func TestDupUpdateSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiUpdateReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 233
	oldMsg.PlayerName = "TestDupUpdateSuccess"
	oldMsg.PlayerEmail = "dsf"
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

	client.Insert(oldMsg)
	defer client.Delete(oldMsg)

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
}
