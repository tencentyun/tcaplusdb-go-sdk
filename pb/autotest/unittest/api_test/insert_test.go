package api

import (
	"fmt"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplusservice"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/terror"
	"testing"
)

//case1 记录不存在时insert success resultFlag = 2
func TestPBInsertSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiInsertReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 233
	oldMsg.PlayerName = "TestPBInsertSuccess"
	oldMsg.PlayerEmail = "dsf"
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

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

//case2 记录存在时，重复插入，返回记录已存在 resultFlag = 2
func TestDupInsert(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiInsertReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 233
	oldMsg.PlayerName = "TestDupInsert"
	oldMsg.PlayerEmail = "dsf"
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

	// 万一已经存在残余数据
	client.Delete(oldMsg)

	// 插入一条数据
	if err := client.Insert(oldMsg); err != nil {
		t.Errorf("Insert fail, %s", err.Error())
		return
	}

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

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_RECORD_EXIST, but nil")
		return
	} else if err != terror.SVR_ERR_FAIL_RECORD_EXIST {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_RECORD_EXIST,but %s", terror.GetErrMsg(err))
		return
	}

	client.Delete(oldMsg)
}
