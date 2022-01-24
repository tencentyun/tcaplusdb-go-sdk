package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"testing"
)

//case 1记录不存在时Get fail
func TestPBGetFail(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 55555
	oldMsg.PlayerName = "NotExist"
	oldMsg.PlayerEmail = "dsf"
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

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

//case2 记录存在时，Get返回成功
func TestPBGetSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 555
	oldMsg.PlayerName = "jiahua"
	oldMsg.PlayerEmail = "dsf"
	client.Insert(oldMsg)
	defer client.Delete(oldMsg)
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

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

//case3 某个key字段不存在，报错.
func TestPBGetFail_Key_NonExist(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 555
	oldMsg.PlayerName = "jiahua"
	oldMsg.PlayerEmail = "dsf"
	client.Insert(oldMsg)
	defer client.Delete(oldMsg)
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)
	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	//设置不存在的key
	oldMsg.PlayerId = 55
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

	if err := resp.GetResult(); err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}
}

//case4 某个Value字段不存在，get能够成功.
func TestPBGetSuccess_Value_NonExist(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 555
	oldMsg.PlayerName = "jiahua"
	oldMsg.PlayerEmail = "dsf"
	client.Insert(oldMsg)
	defer client.Delete(oldMsg)
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)
	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	//设置不存在的value
	oldMsg.GameServerId = 34
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

}

// case 条件符合
func TestPBGetWithCondition(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetReq, "user")

	oldMsg := &tcaplusservice.User{}
	oldMsg.Id = 1
	oldMsg.Name = "a"
	oldMsg.Rank = 1
	mail := &tcaplusservice.UserMail{
		Title:   "tcaplus",
		Content: "...",
	}
	oldMsg.Mailbox = append(oldMsg.Mailbox, mail)
	client.Insert(oldMsg)
	defer client.Delete(oldMsg)
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

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

	rec.SetCondition("mailbox CONTAINS(title == \"tcaplus\")")

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

		newMsg := &tcaplusservice.User{}
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

// case 条件不符
func TestPBGetWithCondition_Fail(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetReq, "user")

	oldMsg := &tcaplusservice.User{}
	oldMsg.Id = 1
	oldMsg.Name = "a"
	oldMsg.Rank = 1
	mail := &tcaplusservice.UserMail{
		Title:   "tcaplus",
		Content: "...",
	}
	oldMsg.Mailbox = append(oldMsg.Mailbox, mail)
	client.Insert(oldMsg)
	defer client.Delete(oldMsg)
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

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

	rec.SetCondition("mailbox NOT CONTAINS(title == \"tcaplus\")")

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

	if err := resp.GetResult(); err != 281 {
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}
}
