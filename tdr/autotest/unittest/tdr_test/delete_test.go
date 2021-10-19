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

//case 1记录不存在时delete fail
func TestDeleteFail(t *testing.T) {
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
		t.Errorf("resp.GetResult expect err ,but nil")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult expect TXHDB_ERR_RECORD_NOT_EXIST ,but %s", terror.GetErrMsg(err))
		return
	}
}

//case2 记录存在时，Delete返回成功
func TestDupDeleteSuccess(t *testing.T) {
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
	data.Simple_Struct.C_Int64 = 100
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

	////2 Get记录已存在
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
			t.Errorf("resData2 != reqData")
			return
		}
	}
}

//case3 Delete错误的key字段，失败
func TestDeleteErrKeyFail(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiGetReq)
	if err != nil {
		t.Errorf("NewRequest fail, %s", err.Error())
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
		t.Errorf("resp.GetResult expect error, but nil")
		return
	} else if err != terror.SVR_ERR_FAIL_MISS_KEY_FIELD {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_MISS_KEY_FIELD, but %s", terror.GetErrMsg(err))
		return
	}
}
