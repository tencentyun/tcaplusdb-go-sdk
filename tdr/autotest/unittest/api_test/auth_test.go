package api_test

import (
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/cfg"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"strings"
	"testing"
)

//appId不存在
func TestDirAppNotExist(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	//日志配置，不配置则debug打印到控制台
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(100, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err == nil {
		t.Errorf("excepted dial fail")
		return
	}

	if !strings.Contains(err.Error(), "errCode: -279") {
		t.Errorf("excepted dir auth fail, but real:%s", err.Error())
		return
	}
}

//ZoneId不存在,dir没有做判断处理，跳过改case
/*
func TestDirZoneNotExist(t *testing.T){
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{100}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature,30)
	if err == nil {
		t.Errorf("excepted dial fail")
		return
	}

	if !strings.Contains(err.Error(), "errCode: -279") {
		t.Errorf("excepted dir auth fail, but real:%s", err.Error())
		return
	}
}*/

//signature not correct
func TestDirSignatureNotCorrect(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature+"notCorrect", 30)
	if err == nil {
		t.Errorf("excepted dial fail")
		return
	}

	if !strings.Contains(err.Error(), "errCode: -279") {
		t.Errorf("excepted dir auth fail, but real:%s", err.Error())
		return
	}
}

//signature correct
func TestDirSignatureCorrect(t *testing.T) {
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
}

//NewRequest zone not exist
func TestNewRequestZoneNotExist(t *testing.T) {
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

	req, err := client.NewRequest(100, TestTableName, cmd.TcaplusApiInsertReq)
	if err == nil {
		t.Errorf("excepted NewRequest fail")
		return
	}

	if !strings.Contains(err.Error(), "errCode: -1054") {
		t.Errorf("excepted NewRequest fail, but real:%s", err.Error())
		return
	}

	if req != nil {
		t.Errorf("excepted NewRequest nil")
		return
	}
}

//NewRequest zone not exist
func TestNewRequestTableNotExist(t *testing.T) {
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

	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, "NotExistTable", cmd.TcaplusApiInsertReq)
	if err == nil {
		t.Errorf("excepted NewRequest fail")
		return
	}

	if !strings.Contains(err.Error(), "errCode: -1310") {
		t.Errorf("excepted NewRequest fail, but real:%s", err.Error())
		return
	}

	if req != nil {
		t.Errorf("excepted NewRequest nil")
		return
	}
}
