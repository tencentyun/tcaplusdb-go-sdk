package api

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"testing"
	"time"
)

//记录存在的时候batchget
func TestBatchGetSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiBatchGetReq, "game_players")

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

	//add record
	rec, err := req.AddRecord(0)
	req.SetMultiResponseFlag(1)
	req.SetResultLimit(0, 1)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	oldData.PlayerEmail = "zhang"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "zhang"
		client.Delete(oldData)
	}()

	//add record
	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	oldData.PlayerEmail = "li"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "li"
		client.Delete(oldData)
	}()

	//add record
	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
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

	if 3 != resp.GetRecordCount() {
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

		if newMsg.PlayerEmail != "zhang" &&
			newMsg.PlayerEmail != "wang" &&
			newMsg.PlayerEmail != "li" {
			t.Errorf("resData != reqData")
			return
		}
	}

}
func TestBatchGetMultiRespSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiBatchGetReq, "game_players")

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

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	oldData.PlayerEmail = "zhang"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "zhang"
		client.Delete(oldData)
	}()

	//add record
	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	oldData.PlayerEmail = "li"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "li"
		client.Delete(oldData)
	}()

	//add record
	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	req.SetMultiResponseFlag(1)

	resps, err := client.DoMore(req, 5*time.Second)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	// 即使设置了分包返回，但是记录没有超过10M,还是只返回一个包。
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

//单条记录超过256KB,设置分包返回。
func TestBatchGetMultiRespSuccess_256KB(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiBatchGetReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldData.LoginTimestamp = []string{tools.Bytes256K()}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"

	oldData.PlayerEmail = "wang"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "wang"
		client.Delete(oldData)
	}()

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	oldData.PlayerEmail = "zhang"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "zhang"
		client.Delete(oldData)
	}()

	//add record
	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	oldData.PlayerEmail = "li"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "li"
		client.Delete(oldData)
	}()

	//add record
	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	req.SetMultiResponseFlag(1)

	resps, err := client.DoMore(req, 10*time.Second)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	// 即使设置了分包返回，记录超过了256KB,就进行分包返回。
	if len(resps) != 3 {
		t.Errorf("recvResponse fail, len(resps) %d != 3", len(resps))
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

//单条记录超过256KB,不设置分包返回。默认分包，设置无效
func TestBatchGetMultiRespSuccess_256KB_Multi_0(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiBatchGetReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	oldData.LoginTimestamp = []string{tools.Bytes256K()}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"

	oldData.PlayerEmail = "wang"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "wang"
		client.Delete(oldData)
	}()

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	oldData.PlayerEmail = "zhang"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "zhang"
		client.Delete(oldData)
	}()

	//add record
	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	oldData.PlayerEmail = "li"
	client.Insert(oldData)
	defer func() {
		oldData.PlayerEmail = "li"
		client.Delete(oldData)
	}()

	//add record
	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(oldData); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	// 默认分包设置无效
	req.SetMultiResponseFlag(0)

	resps, err := client.DoMore(req, 5*time.Second)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	// 即使设置了分包返回，记录超过了256KB,就进行分包返回。
	if len(resps) != 3 {
		t.Errorf("recvResponse fail, len(resps) %d != 3", len(resps))
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
