package api_test

import (
	"fmt"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/cfg"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplus_tb"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
	"testing"
	"time"
)

//case 1记录不存在时Replace success resultFlag = 2
func TestSyncReplaceSuccess(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//data
	data := newGenericTableRec()
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

	resp, err := client.Do(req, time.Duration(2*time.Second))
	if err != nil {
		t.Errorf("Do fail, %s", err.Error())
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
}

//case2 记录存在时，重复Replace，返回成功 resultFlag = 2
func TestSyncDupReplace(t *testing.T) {
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

	resp, err := client.Do(req, time.Duration(2*time.Second))
	if err != nil {
		t.Errorf("Do fail, %s", err.Error())
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

	resp2, err := client.Do(req2, time.Duration(2*time.Second))
	if err != nil {
		t.Errorf("Do fail, %s", err.Error())
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
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

//case3 replace错误的key字段，插入失败
func TestSyncReplaceErrKeyFail(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	//使用SetKey SetValue接口
	rec.SetKeyInt64("uin", 1)
	rec.SetKeyStr("name", "name")
	rec.SetKeyStr("key3", "key3")
	rec.SetKeyStr("notExistKey", "key4")

	rec.SetValueInt32("level", 2)

	resp, err := client.Do(req, time.Duration(2*time.Second))
	if err != nil {
		t.Errorf("Do fail, %s", err.Error())
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect error, but nil")
		return
	} else if err != terror.SVR_ERR_FAIL_MISS_KEY_FIELD {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_MISS_KEY_FIELD, but %s", terror.GetErrMsg(err))
		return
	}
}

//case4 replace错误value字段，插入失败
func TestSyncReplaceErrValueFail(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	//使用SetKey SetValue接口
	rec.SetKeyInt64("uin", 1)
	rec.SetKeyStr("name", "name")
	rec.SetKeyStr("key3", "key3")
	rec.SetKeyStr("key4", "key4")

	rec.SetValueStr("NotExistValue", "value4")

	resp, err := client.Do(req, time.Duration(2*time.Second))
	if err != nil {
		t.Errorf("Do fail, %s", err.Error())
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect error, but nil")
		return
	} else if err != terror.SVR_ERR_FAIL_INVALID_FIELD_NAME {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_INVALID_FIELD_NAME, but %s", terror.GetErrMsg(err))
		return
	}
}
