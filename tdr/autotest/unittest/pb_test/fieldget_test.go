package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestFieldGetSimple(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldGetReq, "game_players")

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

	oldJson := tools.StToJson(oldData)
	fmt.Println(oldJson)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

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

		if newMsg.Pay.PayId != oldData.Pay.PayId ||
			newMsg.Pay.Amount != oldData.Pay.Amount {
			t.Errorf("resData != reqData")
			return
		}
	}
}

func TestFieldGetSuccess(t *testing.T) {
	client := tools.InitPBSyncClient()

	// 插入一条数据
	oldData := &tcaplusservice.GamePlayers{}
	oldData.PlayerId = 444
	oldData.PlayerName = "jiahua"
	oldData.PlayerEmail = "dsf"
	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	client.Insert(oldData)
	defer client.Delete(oldData)

	// 读取数据
	newData := &tcaplusservice.GamePlayers{}
	newData.PlayerId = 444
	newData.PlayerName = "jiahua"
	newData.PlayerEmail = "dsf"
	fields := []string{"pay", "pay.pay_id"}
	err := client.FieldGet(newData, fields)
	if err != nil {
		t.Errorf("Get failed %s", err)
		return
	}

	logger.DEBUG("%+v", newData)

	if !proto.Equal(oldData.Pay, newData.Pay) {
		t.Errorf("data diff \n%+v \n%+v", oldData.Pay, newData.Pay)
		return
	}

	fields = []string{"pay.pay_id", "pay.amount"}
	err = client.FieldGet(newData, fields)
	if err != nil {
		t.Errorf("Get failed %s", err)
		return
	}

	if oldData.Pay.PayId != newData.Pay.PayId {
		t.Errorf("data diff \n%+v \n%+v", oldData.Pay.PayId, newData.Pay.PayId)
		return
	}

	if oldData.Pay.Amount != newData.Pay.Amount {
		t.Errorf("data diff \n%+v \n%+v", oldData.Pay.Amount, newData.Pay.Amount)
		return
	}
}

// case 条件符合
func TestPBFieldGetWithCondition(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldGetReq, "user")

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

	if _, err := rec.SetPBFieldValues(oldMsg, []string{"rank"}); err != nil {
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
		err = record.GetPBFieldValues(newMsg)
		if err != nil {
			t.Errorf("GetPBData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newMsg)
		fmt.Println(newJson)
		if newMsg.Rank != 1 {
			t.Errorf("newMsg.Rank != 1")
			return
		}
	}
}

// case 条件不符
func TestPBFieldGetWithCondition_Fail(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldGetReq, "user")

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

	if _, err := rec.SetPBFieldValues(oldMsg, []string{"rank"}); err != nil {
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

func TestPBFieldGetByPath(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldGetReq, "tb_map")

	initMsg := initTbMap(3, 4, 5, 4, 3)

	client.Insert(initMsg)
	defer client.Delete(initMsg)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	oldMsg := initTbMap(0, 0, 0, 0, 0)
	paths := []string{
		"i32",
		"str",
		"all_type.u64",
		"i32_array",
		"str_map['1']",
		"str_map['2']",
		"int_map",
	}
	if _, err := rec.SetPBFieldValues(oldMsg, paths); err != nil {
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

		newMsg := &tcaplusservice.TbMap{}
		err = record.GetPBFieldValues(newMsg)
		if err != nil {
			t.Errorf("GetPBData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newMsg)
		fmt.Println(newJson)

		if newMsg.I32 != 1 || newMsg.Str != "a" || newMsg.AllType.U64 != 1 || len(newMsg.I32Array) != 5 ||
			len(newMsg.IntMap) != 3 || len(newMsg.StrMap) != 2 {
			t.Errorf("resData != reqData")
			return
		}
	}
}

func TestPBFieldGetByPathCorrect(t *testing.T) {
	client, _ := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldGetReq, "tb_map")

	initMsg := initTbMap(3, 4, 5, 4, 3)

	client.Insert(initMsg)
	defer client.Delete(initMsg)

	oldMsg := initTbMap(0, 0, 0, 0, 0)
	err := client.FieldGet(oldMsg, []string{"int_map[10]"})
	if err == nil || err.(*terror.ErrorCode).Code != terror.SVR_ERR_FAIL_PROTOBUF_FIELD_GET {
		t.Errorf(err.Error())
		return
	}
	err = client.FieldGet(oldMsg, []string{"i32_array[-1]"})
	if err == nil || err.(*terror.ErrorCode).Code != terror.SVR_ERR_FAIL_PROTOBUF_FIELD_GET {
		t.Errorf(err.Error())
		return
	}
	err = client.FieldGet(oldMsg, []string{"str_map['3'].ri32[0]"})
	if err == nil || err.(*terror.ErrorCode).Code != terror.SVR_ERR_FAIL_PROTOBUF_FIELD_GET {
		t.Errorf(err.Error())
		return
	}
}
