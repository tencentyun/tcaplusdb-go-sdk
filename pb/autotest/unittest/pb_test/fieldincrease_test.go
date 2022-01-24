package api

import (
	"fmt"
	"testing"

	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

func TestFieldIncreaseSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldIncreaseReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"

	oldData.PlayerEmail = "wang"
	client.Delete(oldData)
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

		if newMsg.Pay.Amount != 4 ||
			newMsg.Pay.PayId != 3 {
			t.Errorf("resData != reqData")
			return
		}
	}
}

//设置flag,req.SetAddableIncreaseFlag(1),浮点类型的字段，自增浮点型。
func TestFieldIncreaseNotExist(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldIncreaseReq, "game_players")

	recordData := &tcaplusservice.GamePlayers{}
	recordData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2}
	recordData.PlayerId = 233
	recordData.PlayerName = "jiahua"
	recordData.PlayerEmail = "jiahua@qq.com"
	// 删除记录，保证记录不存在
	client.Delete(recordData)
	defer func() {
		client.Delete(recordData)
	}()

	recordJson := tools.StToJson(recordData)
	fmt.Println(recordJson)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	// 设置，当记录不存在则创建新的空记录
	req.SetAddableIncreaseFlag(1)
	//
	recordData.Pay = &tcaplusservice.Payment{Amount: 3, PayId: 1, Anker: 1.1111, Anker_01: 3.3333}
	if _, err := rec.SetPBFieldValues(recordData, []string{"pay.pay_id", "pay.amount", "pay.anker", "pay.anker_01"}); err != nil {
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

		if newMsg.Pay.Amount != 3 ||
			newMsg.Pay.PayId != 1 ||
			newMsg.Pay.Anker != 1.1111 ||
			newMsg.Pay.Anker_01 != 3.3333 {
			t.Errorf("resData != reqData")
			return
		}
	}
}

//设置flag,req.SetAddableIncreaseFlag(1),浮点类型的字段，自增整型。
func TestFieldIncreaseNotExist_03(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldIncreaseReq, "game_players")

	recordData := &tcaplusservice.GamePlayers{}
	recordData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2}
	recordData.PlayerId = 233
	recordData.PlayerName = "jiahua"
	recordData.PlayerEmail = "jiahua@qq.com"
	// 删除记录，保证记录不存在
	client.Delete(recordData)
	defer func() {
		client.Delete(recordData)
	}()

	recordJson := tools.StToJson(recordData)
	fmt.Println(recordJson)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	// 设置，当记录不存在则创建新的空记录
	req.SetAddableIncreaseFlag(1)
	//
	recordData.Pay = &tcaplusservice.Payment{Amount: 3, PayId: 1, Anker: 1, Anker_01: 3}
	if _, err := rec.SetPBFieldValues(recordData, []string{"pay.pay_id", "pay.amount", "pay.anker", "pay.anker_01"}); err != nil {
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

		if newMsg.Pay.Amount != 3 ||
			newMsg.Pay.PayId != 1 ||
			newMsg.Pay.Anker != 1 ||
			newMsg.Pay.Anker_01 != 3 {
			t.Errorf("resData != reqData")
			return
		}
	}
}

//不设置flag
func TestFieldIncreaseNotExist_02(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldIncreaseReq, "game_players")

	recordData := &tcaplusservice.GamePlayers{}
	recordData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2}
	recordData.PlayerId = 233
	recordData.PlayerName = "jiahua"
	recordData.PlayerEmail = "jiahua@qq.com"
	// 删除记录，保证记录不存在
	client.Delete(recordData)
	defer func() {
		client.Delete(recordData)
	}()

	recordJson := tools.StToJson(recordData)
	fmt.Println(recordJson)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	// 设置，当记录不存在则创建新的空记录
	//req.SetAddableIncreaseFlag(1)

	recordData.Pay = &tcaplusservice.Payment{Amount: 3, PayId: 1}
	if _, err := rec.SetPBFieldValues(recordData, []string{"pay.pay_id", "pay.amount"}); err != nil {
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

	if err := resp.GetResult(); err != 261 {
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}

}

func TestFieldIncreaseNotExist_01(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldIncreaseReq, "game_players")

	recordData := &tcaplusservice.GamePlayers{}
	recordData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2}
	recordData.PlayerId = 233
	recordData.PlayerName = "jiahua"
	recordData.PlayerEmail = "jiahua@qq.com"
	// 删除记录，保证记录不存在
	client.Delete(recordData)
	defer func() {
		client.Delete(recordData)
	}()

	recordJson := tools.StToJson(recordData)
	fmt.Println(recordJson)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	// 设置，当记录不存在则创建新的空记录
	req.SetAddableIncreaseFlag(0)

	recordData.Pay = &tcaplusservice.Payment{Amount: 3, PayId: 1}
	if _, err := rec.SetPBFieldValues(recordData, []string{"pay.pay_id", "pay.amount"}); err != nil {
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

	if err := resp.GetResult(); err != 261 {
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}
}

// case 更新附带条件与操作
func TestPBFieldIncreaseWithOperateCondition(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldIncreaseReq, "user")

	oldMsg := &tcaplusservice.User{}
	oldMsg.Id = 1
	oldMsg.Name = "a"
	oldMsg.Rank = 1
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

	oldMsg.Rank = -1
	if _, err := rec.SetPBFieldValues(oldMsg, []string{"rank"}); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	// 条件
	rec.SetCondition("rank > 0")
	// 操作
	rec.SetOperation("PUSH gameids #[-1][$=123]", 0)

	req.SetResultFlag(2)

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
	}

	msg := &tcaplusservice.User{Id: 1, Name: "a"}
	client.Get(msg)
	fmt.Println(msg)
	if msg.Rank != 0 || msg.Gameids[0] != 123 {
		t.Errorf("msg.Rank != 0 || msg.Gameids[0] != 123")
		return
	}
}

// case 更新附带条件
func TestPBFieldIncreaseWithCondition(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldIncreaseReq, "user")

	oldMsg := &tcaplusservice.User{}
	oldMsg.Id = 1
	oldMsg.Name = "a"
	oldMsg.Rank = 0
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

	oldMsg.Rank = -1
	if _, err := rec.SetPBFieldValues(oldMsg, []string{"rank"}); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	// 条件
	rec.SetCondition("rank > 0")

	req.SetResultFlag(2)

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
