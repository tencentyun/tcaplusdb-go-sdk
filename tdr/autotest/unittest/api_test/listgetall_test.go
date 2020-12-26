package api_test

import (
	"fmt"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplus_tb"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
	"testing"
	"time"
)

func newTableTraverserList() *tcaplus_tb.Table_Traverser_List {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255
	data.Level = 1
	data.Value1 = "value1"
	data.Value2 = "value2"
	return data
}

/*
//case2 记录存在时，Get返回成功
func TestGetListGetAllSuccess(t *testing.T) {
	client, req := InitClientAndReqWithTableName(cmd.TcaplusApiListGetAllReq, "table_traverser_list")
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	//if err := req.SetResultFlag(2); err != nil {
	//	t.Errorf("SetResultFlag failed %v", err.Error())
	//	return
	//}
	uinKey := time.Now().UnixNano()
	data := newTableTraverserList()
	data.Key = uint32(uinKey)
	data.Key = 1
	data.Name = 4

	//oldJson := StToJson(data)
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	req.SetResultLimit(100,0)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err,%d, %s", err,terror.GetErrMsg(err))
		return
	}

	if 3 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}
	for idx := 0; idx < resp.GetRecordCount(); idx++{
		recRes, err :=resp.FetchRecord()
		if err != nil{
			t.Errorf("resp.FetchRecord() %s",  err.Error())
			return
		}
		dataRes := newTableTraverserList()
		if err := recRes.GetData(dataRes); err == nil {
			fmt.Printf("%d, %s, %s\n", dataRes.Level, dataRes.Value1, dataRes.Value2)
		}
	}



}
*/

func TestListGetAllSuccess(t *testing.T) {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255
	data.Level = 1
	data.Value1 = "value1"
	data.Value2 = "value2"
	oldJson := tools.StToJson(data)
	fmt.Println(oldJson)

	/*------------------------------------------------------- 清理记录 ------------------------------------------*/
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiListDeleteAllReq, "table_traverser_list")

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec.SetData(data); err != nil {
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
		client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiListAddAfterReq, "table_traverser_list")

		//add record
		rec, err = req.AddRecord(int32(tcaplus_protocol_cs.TCAPLUS_LIST_LAST_INDEX))
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}

		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}

		if err := req.SetResultFlag(0); err != nil {
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

	/*------------------------------------------------------- 获取记录 ------------------------------------------*/
	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiListGetAllReq, "table_traverser_list")

	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	req.SetResultLimit(100, 0)

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
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Traverser_List()
		if err := record.GetData(newData); err != nil {
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

//记录不存在的时候，返回261
func TestGetListGetAllFail(t *testing.T) {
	client, req := InitClientAndReqWithTableName(cmd.TcaplusApiListGetAllReq, "table_traverser_list")
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	//if err := req.SetResultFlag(2); err != nil {
	//      t.Errorf("SetResultFlag failed %v", err.Error())
	//      return
	//}
	uinKey := time.Now().UnixNano()
	data := newTableTraverserList()
	data.Key = uint32(uinKey)
	data.Key = 10000  //不存在的key
	data.Name = 10000 //不存在的key
	rec, err := req.AddRecord(0)

	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	req.SetResultLimit(100, 0)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 261 {
		t.Errorf("resp.GetResult err,%d, %s", err, terror.GetErrMsg(err))
		return
	}
}
