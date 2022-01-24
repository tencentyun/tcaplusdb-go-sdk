package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"testing"
)

func TestListGetSuccess(t *testing.T) {
	data := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}
	oldJson := tools.StToJson(data)
	fmt.Println(oldJson)

	/*------------------------------------------------------- 清理记录 ------------------------------------------*/
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListDeleteAllReq, "tb_online_list")

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
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

	if err := resp.GetResult(); err != 0 && err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	/*------------------------------------------------------- 插入记录 ------------------------------------------*/
	client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListAddAfterReq, "tb_online_list")

	//add record
	rec, err = req.AddRecord(int32(tcaplus_protocol_cs.TCAPLUS_LIST_LAST_INDEX))
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err = tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
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

		newData := &tcaplusservice.TbOnlineList{}
		if err := record.GetPBData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	/*------------------------------------------------------- 获取记录 ------------------------------------------*/
	client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListGetReq, "tb_online_list")

	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err = tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
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

		newData := &tcaplusservice.TbOnlineList{}
		if err := record.GetPBData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}
}

// 获取不存在的index，失败
func TestListGetIndexNotExist(t *testing.T) {
	data := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}
	oldJson := tools.StToJson(data)
	fmt.Println(oldJson)

	/*------------------------------------------------------- 清理记录 ------------------------------------------*/
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListDeleteAllReq, "tb_online_list")

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
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

	if err := resp.GetResult(); err != 0 && err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	/*------------------------------------------------------- 插入记录 ------------------------------------------*/
	client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListAddAfterReq, "tb_online_list")

	//add record
	rec, err = req.AddRecord(int32(tcaplus_protocol_cs.TCAPLUS_LIST_LAST_INDEX))
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err = tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
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

		newData := &tcaplusservice.TbOnlineList{}
		if err := record.GetPBData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	/*------------------------------------------------------- 获取记录 ------------------------------------------*/
	client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListGetReq, "tb_online_list")

	rec, err = req.AddRecord(233)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	data.Gamesvrid = "newData"
	oldJson = tools.StToJson(data)
	fmt.Println(oldJson)

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err = tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}
}

// 条件符合
func TestListGetWithCondition(t *testing.T) {
	data := &tcaplusservice.ListUser{}
	data.Id = 1
	data.Name = "a"
	data.Rank = 100
	oldJson := tools.StToJson(data)
	fmt.Println(oldJson)

	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListGetReq, "list_user")

	client.ListDeleteAll(data)
	client.ListAddAfter(data, -1)
	defer client.ListDeleteAll(data)

	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	rec.SetCondition("rank == 0 OR $.LastAccessTime > \"2021-01-01 00:00:00\"")

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
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
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

		newData := &tcaplusservice.ListUser{}
		if err := record.GetPBData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}
}

// 条件不符
func TestListGetWithCondition_Fail(t *testing.T) {
	data := &tcaplusservice.ListUser{}
	data.Id = 1
	data.Name = "a"
	data.Rank = 100
	oldJson := tools.StToJson(data)
	fmt.Println(oldJson)

	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListGetReq, "list_user")

	client.ListDeleteAll(data)
	client.ListAddAfter(data, -1)
	defer client.ListDeleteAll(data)

	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	rec.SetCondition("rank == 0")

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
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}
}
