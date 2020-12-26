package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"testing"
	"time"
)

//////////////////////////////INSERT////////////////////////
//case1 记录不存在时insert success resultFlag = 0
//case2 记录存在时，重复插入，返回记录已存在 resultFlag = 0
func TestDupInsertFlag0(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 插入成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	////2 插入记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)

	oldJson2 := StToJson(data)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_RECORD_EXIST, but nil")
		return
	} else if err != terror.SVR_ERR_FAIL_RECORD_EXIST {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_RECORD_EXIST,but %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp2.GetRecordCount() {
		t.Errorf("resp2.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}
}

//case3 记录不存在时insert success resultFlag = 1
//case4 记录存在时，重复插入，返回记录已存在 resultFlag = 1
func TestDupInsertFlag1(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 插入成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
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

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	////2 插入记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_RECORD_EXIST, but nil")
		return
	} else if err != terror.SVR_ERR_FAIL_RECORD_EXIST {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_RECORD_EXIST,but %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp2.GetRecordCount() %d != 1", resp2.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson2 {
			t.Errorf("resData != reqData")
			return
		}
	}
}

//case5 记录不存在时insert success resultFlag = 3
//case6 记录存在时，重复插入，返回记录已存在 resultFlag = 3
func TestDupInsertFlag3(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 插入成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	////2 插入记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_RECORD_EXIST, but nil")
		return
	} else if err != terror.SVR_ERR_FAIL_RECORD_EXIST {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_RECORD_EXIST,but %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp2.GetRecordCount() {
		t.Errorf("resp2.GetRecordCount() %d != 1", resp2.GetRecordCount())
		return
	}
}

//////////////////////////////REPLACE////////////////////////
//case7 记录不存在时Replace success resultFlag = 0
//case8 记录存在时，重复插入，返回记录已存在 resultFlag = 0
func TestDupReplaceFlag0(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 replace成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	////2 replace记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}
}

//case9 记录不存在时Replace success resultFlag = 1
//case10 记录存在时，重复插入，返回记录已存在 resultFlag = 1
func TestDupReplaceFlag1(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 replace成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
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

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	////2 replace记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson2 {
			t.Errorf("resData2 != reqData2")
			return
		}
	}
}

//case11 记录不存在时Replace success resultFlag = 2
//case12 记录存在时，重复插入，返回记录已存在 resultFlag = 2
func TestDupReplaceFlag2(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 replace成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
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

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	////2 replace记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson2 {
			t.Errorf("resData2 != reqData2")
			return
		}
	}
}

//case13 记录不存在时Replace success resultFlag = 3
//case14 记录存在时，重复插入，返回记录已存在 resultFlag = 3
func TestDupReplaceFlag3(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 replace成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}

	////2 replace记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData2 != reqData2")
			return
		}
	}
}

//////////////////////////////UPDATE////////////////////////
//case13 记录不存在时update fail resultFlag = 0
func TestUpdateFlag0(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 insert成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiUpdateReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case14 记录存在时，update success resultFlag = 0
func TestDupUpdateFlag0(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 insert成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	////2 update记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiUpdateReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case15 记录不存在时update fail resultFlag = 1
func TestUpdateFlag1(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiUpdateReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case16 记录存在时，update success resultFlag = 1
func TestDupUpdateFlag1(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 insert成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
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

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	////2 update记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiUpdateReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson2 {
			t.Errorf("resData2 != reqData2")
			return
		}
	}
}

//case17 记录不存在时update fail resultFlag = 2
func TestUpdateFlag2(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiUpdateReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case18 记录存在时，update success resultFlag = 2
func TestDupUpdateFlag2(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 insert成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
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

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	////2 update记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiUpdateReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson2 {
			t.Errorf("resData2 != reqData2")
			return
		}
	}
}

//case19 记录不存在时update fail resultFlag = 3
func TestUpdateFlag3(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiUpdateReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case20 记录存在时，update success resultFlag = 3
func TestDupUpdateFlag3(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 insert成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}

	////2 update记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiUpdateReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData2 != reqData2")
			return
		}
	}
}

//////////////////////////////DELETE////////////////////////
//case21 记录不存在时delete fail resultFlag = 0
func TestDeleteFlag0(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiDeleteReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case22 记录存在时，delete success resultFlag = 0
func TestDupDeleteFlag0(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 insert成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	////2 delete记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiDeleteReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(0); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case23 记录不存在时delete fail resultFlag = 1
func TestDeleteFlag1(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiDeleteReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case24 记录存在时，delete success resultFlag = 1
func TestDupDeleteFlag1(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 insert成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
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

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	////2 update记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiDeleteReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(1); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}
}

//case25 记录不存在时delete fail resultFlag = 2
func TestDeleteFlag2(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiDeleteReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case26 记录存在时，delete success resultFlag = 2
func TestDupDeleteFlag2(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 insert成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
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

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}

	////2 update记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiDeleteReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData2 != reqData2")
			return
		}
	}
}

//case27 记录不存在时delete fail resultFlag = 3
func TestDeleteFlag3(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiDeleteReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}
}

//case28 记录存在时，delete success resultFlag = 3
func TestDupDeleteFlag3(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	///////1 insert成功
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	fmt.Println(oldJson)
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
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 0 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 0", resp.GetRecordCount())
		return
	}

	////2 update记录已存在
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiDeleteReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req2.SetResultFlag(3); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//相同的key
	data2 := newGenericTableRec()
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	fmt.Println(oldJson2)
	//add record
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetData(data2); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req2); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	if 1 != resp2.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData2 != reqData2")
			return
		}
	}
}
