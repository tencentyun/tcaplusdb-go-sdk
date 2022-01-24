package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"testing"
)

//case 1记录不存在时delete fail
func TestPBDeleteFail(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiDeleteReq, "game_players")

	oldMsg := &tcaplusservice.GamePlayers{}
	oldMsg.PlayerId = 233
	oldMsg.PlayerName = "TestPBDeleteFail"
	oldMsg.PlayerEmail = "dsf"
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

	// 万一已经存在残余数据
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

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err ,but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult expect TXHDB_ERR_RECORD_NOT_EXIST ,but %s", terror.GetErrMsg(err))
		return
	}
}

//case2 记录存在时，Delete返回成功,resultflag=2
func TestDupDeleteSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiDeleteReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "TestDupDeleteSuccess"
	oldData.PlayerEmail = "dsf"

	// 万一已经存在残余数据
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	fmt.Println(oldJson)
	// 插入一条数据
	if err := client.Insert(oldData); err != nil {
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

	if _, err := rec.SetPBData(oldData); err != nil {
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

//记录存在时，Delete返回成功,resultflag=0
func TestDupDeleteSuccess_Flag_0(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiDeleteReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "TestDupDeleteSuccess"
	oldData.PlayerEmail = "dsf"

	// 万一已经存在残余数据
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	fmt.Println(oldJson)
	// 插入一条数据
	if err := client.Insert(oldData); err != nil {
		t.Errorf("Insert fail, %s", err.Error())
		return
	}

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
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

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}

}

//记录存在时，Delete返回成功,resultflag=1
func TestDupDeleteSuccess_Flag_1(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiDeleteReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "TestDupDeleteSuccess"
	oldData.PlayerEmail = "dsf"

	// 万一已经存在残余数据
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	// 插入一条数据
	if err := client.Insert(oldData); err != nil {
		t.Errorf("Insert fail, %s", err.Error())
		return
	}

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	oldData.Pay = nil
	oldJson := tools.StToJson(oldData)
	fmt.Println(oldJson)
	if _, err := rec.SetPBData(oldData); err != nil {
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
		_, err = record.GetPBKey(newMsg)
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

//记录存在时，Delete返回成功,resultflag=3
func TestDupDeleteSuccess_Flag_3(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiDeleteReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "TestDupDeleteSuccess"
	oldData.PlayerEmail = "dsf"

	// 万一已经存在残余数据
	client.Delete(oldData)

	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldJson := tools.StToJson(oldData)
	fmt.Println(oldJson)
	// 插入一条数据
	if err := client.Insert(oldData); err != nil {
		t.Errorf("Insert fail, %s", err.Error())
		return
	}

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
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

// case 条件符合
func TestPBDeleteWithCondition(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiDeleteReq, "user")

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

	req.SetResultFlagForSuccess(3)

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
func TestPBDeleteWithCondition_Fail(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiDeleteReq, "user")

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

	req.SetResultFlagForSuccess(3)
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
