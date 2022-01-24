package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"testing"
	"time"
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

func TestPBGetBypartkeyMultiRespSuccess(t *testing.T) {
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

	req.SetMultiResponseFlag(1)

	resps, err := client.DoMore(req, 5*time.Second)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	// 这里修改为想要的响应数量，设置了分包返回，但是记录没有超过10MB,还是只返回一个包。
	if len(resps) != 1 {
		t.Errorf("recvResponse fail, len(resps) %d != 1", len(resps))
		return
	}

	for _, resp := range resps {
		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
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
}

//构造超过256B的记录。设置反包返回。
func TestPBGetBypartkeyMultiRespSuccess_256B(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetByPartkeyReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldData.LoginTimestamp = []string{tools.Bytes256K()}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"

	oldData.PlayerEmail = "wang"
	//client.Insert(oldData)
	err := client.Insert(oldData)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	defer func() {
		oldData.PlayerEmail = "wang"
		client.Delete(oldData)
	}()

	oldData.PlayerEmail = "wang"
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

	req.SetMultiResponseFlag(1)

	resps, err := client.DoMore(req, 5*time.Second)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	// 这里修改为想要的响应数量，设置了分包返回，超过256K,就进行分包返回。
	if len(resps) != 2 {
		t.Errorf("recvResponse fail, len(resps) %d != 2", len(resps))
		return
	}

	for _, resp := range resps {
		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
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
}

//构造超过256B的记录。设置反包返回。partkey不支持
/*func TestPBGetBypartkeyMultiRespSuccess_256B_Multi_0(t *testing.T){
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetByPartkeyReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldData.LoginTimestamp = []string{tools.Bytes256K()}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"

	oldData.PlayerEmail = "wang"
	//client.Insert(oldData)
	err := client.Insert(oldData)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	defer func() {
		oldData.PlayerEmail = "wang"
		client.Delete(oldData)
	}()

	oldData.PlayerEmail = "wang"
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

	req.SetMultiResponseFlag(0)

	resps, err := client.DoMore(req, 5*time.Second)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	// 这里修改为想要的响应数量，没有设置分包返回，超过256K,也只返回一个包。
	if len(resps) != 1 {
		t.Errorf("recvResponse fail, len(resps) %d != 1", len(resps))
		return
	}

	for _, resp := range resps {
		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
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
}
*/

// case 附带条件
func TestPBGetBypartkeyWithCondition(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetByPartkeyReq, "user")

	oldMsg := &tcaplusservice.User{}
	oldMsg.Id = 1
	oldMsg.Name = "aaa"
	oldMsg.Rank = 1
	client.Insert(oldMsg)
	defer client.Delete(oldMsg)

	oldMsg.Id = 2
	oldMsg.Name = "aaa"
	client.Insert(oldMsg)
	defer func() {
		oldMsg.Id = 2
		oldMsg.Name = "aaa"
		client.Delete(oldMsg)
	}()

	oldMsg.Id = 3
	oldMsg.Name = "bbb"
	client.Insert(oldMsg)
	defer func() {
		oldMsg.Id = 3
		oldMsg.Name = "bbb"
		client.Delete(oldMsg)
	}()

	oldMsg.Id = 4
	oldMsg.Name = "bbb"
	client.Insert(oldMsg)
	defer func() {
		oldMsg.Id = 4
		oldMsg.Name = "bbb"
		client.Delete(oldMsg)
	}()

	oldMsg.Id = 5
	oldMsg.Name = "aaa"
	client.Insert(oldMsg)
	defer func() {
		oldMsg.Id = 5
		oldMsg.Name = "aaa"
		client.Delete(oldMsg)
	}()

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	oldMsg.Id = 1
	oldMsg.Name = "aaa"
	if _, err := rec.SetPBPartKeys(oldMsg, []string{"id", "name"}); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	// 条件
	rec.SetCondition("id > 0")

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

		fmt.Println(newMsg)
	}
}
