package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"testing"
)

// 全部索引存在
func TestListDeleteBatchSuccess(t *testing.T) {
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
	for i := 0; i < 5; i++ {
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
	}

	/*------------------------------------------------------- 删除记录 ------------------------------------------*/
	client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListDeleteBatchReq, "tb_online_list")

	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	for i := int32(0); i < 5; i++ {
		req.AddElementIndex(i)
	}

	// 顺带测试flag=2时返回被删除数据
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

	if 5 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 5", resp.GetRecordCount())
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

// 部分索引不存在
func TestListDeleteBatchPartNotExist(t *testing.T) {
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
	for i := 0; i < 5; i++ {
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
	}

	/*------------------------------------------------------- 删除记录 ------------------------------------------*/
	client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListDeleteBatchReq, "tb_online_list")

	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	indexs := []int32{0, 5, 1, 6, 2}
	for _, i := range indexs {
		req.AddElementIndex(i)
	}

	// 顺带测试flag=2时返回被删除数据
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

// 全部索引不存在
func TestListDeleteBatchAllNotExist(t *testing.T) {
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
	for i := 0; i < 5; i++ {
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
	}

	/*------------------------------------------------------- 删除记录 ------------------------------------------*/
	client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListDeleteBatchReq, "tb_online_list")

	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	indexs := []int32{54, 55, 56, 57, 58}
	for _, i := range indexs {
		req.AddElementIndex(i)
	}

	// 顺带测试flag=2时返回被删除数据
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

	if err := resp.GetResult(); err != terror.SVR_ERR_FAIL_INVALID_INDEX {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}
}

// 全部索引存在
func TestListDeleteBatchCondition(t *testing.T) {
	data := &tcaplusservice.ListUser{}
	data.Id = 1
	data.Name = "a"
	oldJson := tools.StToJson(data)
	fmt.Println(oldJson)

	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListDeleteBatchReq, "list_user")

	client.ListDeleteAll(data)
	for i := int32(1); i <= 5; i++ {
		data.Rank = i
		client.ListAddAfter(data, -1)
	}
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

	rec.SetCondition("rank != 2")
	req.AddElementIndex(1) // rank = 2
	req.AddElementIndex(2) // rank = 3
	req.AddElementIndex(3) // rank = 4
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
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 2 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 2", resp.GetRecordCount())
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
		if newData.Rank != 3 && newData.Rank != 4 {
			t.Errorf("newData.Rank <= 2 || newData.Rank >= 7")
			return
		}
	}
}
