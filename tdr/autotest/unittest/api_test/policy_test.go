package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"testing"
	"time"
)

//////////////////////////////////CHECKDATAVERSION_AUTOINCREASE//////////////////////////////
func TestPolicyCheckAutoInc(t *testing.T) {
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

	//set policy CHECKDATAVERSION_AUTOINCREASE
	req.SetVersionPolicy(policy.CheckDataVersionAutoIncrease)

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

	//version 9
	rec.SetVersion(9)

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

	////2 update记录已存在 version error
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiUpdateReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	//set policy CHECKDATAVERSION_AUTOINCREASE
	req2.SetVersionPolicy(policy.CheckDataVersionAutoIncrease)

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
	//version 12
	rec2.SetVersion(12)

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
		t.Errorf("resp.GetResult expect err, but nil")
		return
	} else if err != terror.SVR_ERR_FAIL_INVALID_VERSION {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_INVALID_VERSION, but %s", terror.GetErrMsg(err))
		return
	}

	////3 replace记录已存在 version error
	req3, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	//set policy CHECKDATAVERSION_AUTOINCREASE
	req3.SetVersionPolicy(policy.CheckDataVersionAutoIncrease)

	//相同的key
	data3 := newGenericTableRec()
	data3.Uin = uint64(uinKey)
	//不同的value
	data3.Level = 222

	oldJson3 := StToJson(data3)
	fmt.Println(oldJson3)
	//add record
	rec3, err := req3.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec3.SetData(data3); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	//version 12
	rec3.SetVersion(12)

	if err := client.SendRequest(req3); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp3, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp3.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err, but nil")
		return
	} else if err != terror.SVR_ERR_FAIL_INVALID_VERSION {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_INVALID_VERSION, but %s", terror.GetErrMsg(err))
		return
	}
}

//////////////////////////////////NOCHECKDATAVERSION_AUTOINCREASE//////////////////////////////
func TestUpdatePolicyNoCheckAutoInc(t *testing.T) {
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

	//set policy
	req.SetVersionPolicy(policy.NoCheckDataVersionAutoIncrease)

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

		if record.GetVersion() != 1 {
			t.Errorf("record version %d != 1", record.GetVersion())
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

	req2.SetVersionPolicy(policy.NoCheckDataVersionAutoIncrease)

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

	rec2.SetVersion(111)

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
		t.Errorf("resp2.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetVersion() != 2 {
			t.Errorf("record version %d != 1", record.GetVersion())
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
func TestReplacePolicyNoCheckAutoInc(t *testing.T) {
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

	//set policy
	req.SetVersionPolicy(policy.NoCheckDataVersionAutoIncrease)

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

		if record.GetVersion() != 1 {
			t.Errorf("record version %d != 1", record.GetVersion())
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
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	req2.SetVersionPolicy(policy.NoCheckDataVersionAutoIncrease)

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

	rec2.SetVersion(111)

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
		t.Errorf("resp2.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetVersion() != 2 {
			t.Errorf("record version %d != 1", record.GetVersion())
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

//////////////////////////////////NOCHECKDATAVERSION_OVERWRITE//////////////////////////////
func TestUpdatePolicyNoCheckOverWrite(t *testing.T) {
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

	//set policy
	req.SetVersionPolicy(policy.NoCheckDataVersionOverwrite)

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

	rec.SetVersion(9)
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

		if record.GetVersion() != 1 {
			t.Errorf("record version %d != 1", record.GetVersion())
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

	req2.SetVersionPolicy(policy.NoCheckDataVersionOverwrite)

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

	rec2.SetVersion(11)

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
		t.Errorf("resp2.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetVersion() != 11 {
			t.Errorf("record version %d != 1", record.GetVersion())
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
func TestReplacePolicyNoCheckOverWrite(t *testing.T) {
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

	//set policy
	req.SetVersionPolicy(policy.NoCheckDataVersionOverwrite)

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

	rec.SetVersion(9)
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

		if record.GetVersion() != 1 {
			t.Errorf("record version %d != 1", record.GetVersion())
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
	req2, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	req2.SetVersionPolicy(policy.NoCheckDataVersionOverwrite)

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

	rec2.SetVersion(12)

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
		t.Errorf("resp2.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		if record.GetVersion() != 12 {
			t.Errorf("record version %d != 1", record.GetVersion())
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
