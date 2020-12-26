package api_test

import (
	"fmt"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/cfg"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"strings"
	"testing"
)

//case1 key字段长度为31正常插入，见insert_test, table_generic的bound_31_byte_test_012345678901字段已经覆盖
//case2 key字段长度为32失败
func TestKeyNameLen32Fail(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
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
	if err := rec.SetKeyInt64("bound_32_byte_test_0123456789010", 1); err == nil {
		t.Errorf("rec.SetKeyInt64 expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -2334") {
		t.Errorf("rec.SetKeyInt64 expect errCode: -2334, errMsg: record中key名称长度超限, but %s", err.Error())
		return
	}
}

//case3 value字段长度为32失败
func TestValueNameLen32Fail(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
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
	if err := rec.SetValueInt64("bound_32_byte_test_0123456789010", 1); err == nil {
		t.Errorf("rec.SetKeyInt64 expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -3102") {
		t.Errorf("errCode: -3102, errMsg: record中value名称长度超限, but %s", err.Error())
		return
	}
}

//case4 key个数超过8失败
func TestKeyNum9Fail(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
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

	for i := 0; i < 8; i++ {
		name := fmt.Sprintf("keyName%d", i)
		if err := rec.SetKeyInt64(name, int64(i)); err != nil {
			t.Errorf("rec.SetKeyInt64 err %s", err.Error())
			return
		}
	}

	//使用SetKey SetValue接口
	if err := rec.SetKeyInt64("keyName9", 1); err == nil {
		t.Errorf("rec.SetKeyInt64 expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -2846") {
		t.Errorf("expect errCode: -2846, errMsg: record中key数量超限, but %s", err.Error())
		return
	}
}

//case4 value个数超过256失败
func TestValueNum257Fail(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
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

	for i := 0; i < 256; i++ {
		name := fmt.Sprintf("valueName%d", i)
		if err := rec.SetValueInt64(name, int64(i)); err != nil {
			t.Errorf("rec.SetKeyInt64 err %s", err.Error())
			return
		}
	}

	//使用SetKey SetValue接口
	if err := rec.SetValueInt64("valueName257", 1); err == nil {
		t.Errorf("rec.SetValueInt64 expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -3614") {
		t.Errorf("expect errCode: -3614, errMsg: record中value数量超限, but %s", err.Error())
		return
	}
}

//case5 key的data超过最大长度1024
func TestKeyData1025Fail(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
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

	//使用SetKey
	if err := rec.SetKeyStr("str1024", string(make([]byte, 1024))); err == nil {
		t.Errorf("rec.SetValueInt64 expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -2590") {
		t.Errorf("expect errCode: -2590, errMsg: record中key值长度超限, but %s", err.Error())
		return
	}

	//使用SetKey
	if err := rec.SetKeyBlob("byte1025", make([]byte, 1025)); err == nil {
		t.Errorf("rec.SetValueInt64 expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -2590") {
		t.Errorf("expect errCode: -2590, errMsg: record中key值长度超限, but %s", err.Error())
		return
	}
}

//key的data超过最大长度1024,并且设置key的长度为1024
func TestKeyData1024Succ(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
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

	//使用SetKey
	if err := rec.SetKeyStr("str1024", string(make([]byte, 1024))); err == nil {
		t.Errorf("rec.SetValueInt64 expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -2590") {
		t.Errorf("expect errCode: -2590, errMsg: record中key值长度超限, but %s", err.Error())
		return
	}

}

// value的data不超过最大长度256kB
func TestValueData256KBSucc(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
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

	//使用SetValue
	if err := rec.SetValueStr("str256kb", string(make([]byte, 256*1024))); err == nil {
		t.Errorf("rec.SetValueStr expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -3358") {
		t.Errorf("expect errCode: -3358, errMsg: record中value值长度超限, but %s", err.Error())
		return
	}

}

//case5 value的data超过最大长度256KB
func TestValueData256KBFail(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
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

	//使用SetValue
	if err := rec.SetValueStr("str256kb", string(make([]byte, 256*1024))); err == nil {
		t.Errorf("rec.SetValueStr expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -3358") {
		t.Errorf("expect errCode: -3358, errMsg: record中value值长度超限, but %s", err.Error())
		return
	}

	//使用SetValue
	if err := rec.SetValueBlob("byte256kb", make([]byte, 256*1024+1)); err == nil {
		t.Errorf("rec.SetValueBlob expect err, but nil")
		return
	} else if !strings.Contains(err.Error(), "errCode: -3358") {
		t.Errorf("expect errCode: -3358, errMsg: record中value值长度超限, but %s", err.Error())
		return
	}
}
