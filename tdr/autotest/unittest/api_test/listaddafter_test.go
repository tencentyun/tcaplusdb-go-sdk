package api_test

import (
	"fmt"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplus_tb"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
	"testing"
)

// 插入记录到末尾
func TestListAddAfterToTail(t *testing.T) {
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
	for i := int32(0); i < 5; i++ {
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

		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetIndex() != i {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

		if record.GetIndex() != int32(i) {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

// 插入记录到开头
func TestListAddAfterToHead(t *testing.T) {
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
	for i := int32(0); i < 5; i++ {
		client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiListAddAfterReq, "table_traverser_list")

		//add record

		rec, err = req.AddRecord(int32(tcaplus_protocol_cs.TCAPLUS_LIST_PRE_FIRST_INDEX))
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}

		if err := rec.SetData(data); err != nil {
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

		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetIndex() != i {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

		if record.GetIndex() != int32(4-i) {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

// 插入到一条不存在index之后
func TestListAddAfterIndexNotExist(t *testing.T) {
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

	record, err := resp.FetchRecord()
	if err != nil {
		t.Errorf("FetchRecord failed %s", err.Error())
		return
	}

	if record.GetIndex() != 0 {
		t.Errorf("index failed %d-%d", record.GetIndex(), 0)
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

	/*------------------------------------------------------- 插入index不存在 ------------------------------------------*/
	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiListAddAfterReq, "table_traverser_list")

	//此时 index=3 并不存在
	rec, err = req.AddRecord(3)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec.SetData(data); err != nil {
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

	if err := resp.GetResult(); err != terror.SVR_ERR_FAIL_INVALID_INDEX {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}
}

// 插入记录到已存在的index后
func TestListAddAfterIndexExist(t *testing.T) {
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

	/*----------------------------------------------------- 插入2条记录 ------------------------------------------*/
	for i := int32(0); i < 2; i++ {
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

		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetIndex() != i {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

	/*--------------------------------------------------- 插入到 index=1 之后 ------------------------------------------*/
	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiListAddAfterReq, "table_traverser_list")

	//add record
	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec.SetData(data); err != nil {
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

	record, err := resp.FetchRecord()
	if err != nil {
		t.Errorf("FetchRecord failed %s", err.Error())
		return
	}

	if record.GetIndex() != 2 {
		t.Errorf("index failed %d-2", record.GetIndex())
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
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	indexs := []int32{0, 2, 1}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetIndex() != indexs[i] {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

// 记录已满继续插入
func TestListAddAfterOnListFull(t *testing.T) {
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

	/*------------------------------------------------------- 插满记录 ------------------------------------------*/
	for i := int32(0); i < 1023; i++ {
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

		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetIndex() != i {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

	/*------------------------------------------------------- shift=0 ------------------------------------------*/
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

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag fail, %s", err.Error())
		return
	}

	if req.SetListShiftFlag(byte(tcaplus_protocol_cs.TCAPLUS_LIST_SHIFT_NONE)) ==
		int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH) {
		t.Errorf("SetListShiftFlag fail, req not support")
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

	if err := resp.GetResult(); err != terror.SVR_ERR_FAIL_LIST_FULL {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	/*------------------------------------------------------- shift=1 ------------------------------------------*/
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

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag fail, %s", err.Error())
		return
	}

	if req.SetListShiftFlag(byte(tcaplus_protocol_cs.TCAPLUS_LIST_SHIFT_HEAD)) ==
		int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH) {
		t.Errorf("SetListShiftFlag fail, req not support")
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

	if 2 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	indexs := []int32{1023, 0}
	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		fmt.Println(i)
		if record.GetIndex() != indexs[i] {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

	/*----------------------------------------------- 获取记录头部被删除 ------------------------------------------*/
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

	req.SetResultLimit(10, 0)

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

	if 10 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	// 记录从 1开始，表示0被删除
	indexs = []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetIndex() != indexs[i] {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

	/*------------------------------------------------------- shift=2 ------------------------------------------*/
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

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag fail, %s", err.Error())
		return
	}

	if req.SetListShiftFlag(byte(tcaplus_protocol_cs.TCAPLUS_LIST_SHIFT_TAIL)) ==
		int32(terror.API_ERR_OPERATION_TYPE_NOT_MATCH) {
		t.Errorf("SetListShiftFlag fail, req not support")
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

	if 2 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	indexs = []int32{1024, 1023}
	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		fmt.Println(i)
		if record.GetIndex() != indexs[i] {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

	/*----------------------------------------------- 获取记录尾部被删除 ------------------------------------------*/
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

	req.SetResultLimit(10, 1018)

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

	// 记录1023被删除
	indexs = []int32{1019, 1020, 1021, 1022, 1024}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetIndex() != indexs[i] {
			t.Errorf("index failed %d-%d", record.GetIndex(), i)
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

// 插入记录且记录key不存在
func TestListAddAfterKeyNotExist(t *testing.T) {
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

	/*------------------------------------------------- 插入记录 ------------------------------------------*/
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

	if err := rec.SetKeyStr("KeyNotExist", "KeyNotExist"); err != nil {
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

	if err := resp.GetResult(); err != terror.SVR_ERR_FAIL_INVALID_FIELD_NAME {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}
}
