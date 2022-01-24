package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"testing"
	"time"
)

//case1 set ttl，Get返回成功
func TestGetTTLSuccess(t *testing.T) {
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
	client.SetDefaultZoneId(cfg.ApiConfig.ZoneId)

	///////1 insert 3 条成功
	data1 := newGenericTableRec()
	data1.Name = "ttl1"
	err = client.DoInsert(TestTableName, data1, nil)
	if err != nil {
		t.Errorf("DoInsert fail, %s", err.Error())
		return
	}

	data2 := newGenericTableRec()
	data2.Name = "ttl2"
	err = client.DoInsert(TestTableName, data2, nil)
	if err != nil {
		t.Errorf("DoInsert fail, %s", err.Error())
		return
	}

	data3 := newGenericTableRec()
	data3.Name = "ttl3"
	err = client.DoInsert(TestTableName, data3, nil)
	if err != nil {
		t.Errorf("DoInsert fail, %s", err.Error())
		return
	}

	////2 set ttl
	opt := &option.TDROpt{
		BatchTTL: []option.TTLInfo{{TTL: 5000}, {TTL: 5000}, {TTL: 5000}},
	}
	dataSlice := []record.TdrTableSt{data1, data2, data3}
	err = client.DoSetTTLBatch(TestTableName, dataSlice, nil, opt)
	if err != nil {
		t.Errorf("DoSetTTLBatch fail, %s", err.Error())
		return
	}

	//3 get ttl
	opt = &option.TDROpt{}
	err = client.DoGetTTLBatch(TestTableName, dataSlice, nil, opt)
	if err != nil {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}
	fmt.Println(opt.BatchTTL)
	for _, ttl := range opt.BatchTTL {
		if ttl.TTL <= 0 || ttl.TTL > 5000 {
			t.Errorf("ttl invalid %d", ttl.TTL)
			return
		}
	}

	//4 5s after
	time.Sleep(5 * time.Second)
	opt = &option.TDROpt{}
	err = client.DoGetTTLBatch(TestTableName, dataSlice, nil, opt)
	fmt.Println(opt.BatchResult)
	if err == nil {
		t.Errorf("DoGetTTLMust timeout, %v", opt.BatchResult)
		return
	}
}
//case1 set ttl，Get返回成功
func TestGetTTLSuccess_ttl_0(t *testing.T) {
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
	client.SetDefaultZoneId(cfg.ApiConfig.ZoneId)

	///////1 insert 3 条成功
	data1 := newGenericTableRec()
	data1.Name = "ttl1"
	err = client.DoInsert(TestTableName, data1, nil)
	if err != nil {
		t.Errorf("DoInsert fail, %s", err.Error())
		return
	}

	data2 := newGenericTableRec()
	data2.Name = "ttl2"
	err = client.DoInsert(TestTableName, data2, nil)
	if err != nil {
		t.Errorf("DoInsert fail, %s", err.Error())
		return
	}

	data3 := newGenericTableRec()
	data3.Name = "ttl3"
	err = client.DoInsert(TestTableName, data3, nil)
	if err != nil {
		t.Errorf("DoInsert fail, %s", err.Error())
		return
	}

	////2 set ttl
	opt := &option.TDROpt{
		BatchTTL: []option.TTLInfo{{TTL: 0}, {TTL: 0}, {TTL: 0}},
	}
	dataSlice := []record.TdrTableSt{data1, data2, data3}
	err = client.DoSetTTLBatch(TestTableName, dataSlice, nil, opt)
	if err != nil {
		t.Errorf("DoSetTTLBatch fail, %s", err.Error())
		return
	}

	//3 get ttl
	opt = &option.TDROpt{}
	err = client.DoGetTTLBatch(TestTableName, dataSlice, nil, opt)
	if err.Error() != "errCode: 2309, errMsg: " {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}

}
// case2 list ttl set,list table not support ttl
//func TestListTTLSuccess(t *testing.T) {
//	client,err := tools.InitClient()
//	if err != nil {
//		t.Errorf("InitClient failed %s", err.Error())
//		return
//	}
//	tableName := "table_traverser_list"
//	//add 3 record
//	data := tcaplus_tb.NewTable_Traverser_List()
//	data.Key = 1
//	data.Name = 255
//	data.Level = 1
//	data.Value1 = "value1"
//	data.Value2 = "value2"
//	err = client.DoListAddAfter(tableName, data, -1, nil)
//	if err != nil {
//		t.Errorf("DoListAddAfter failed %s", err.Error())
//		return
//	}
//	err = client.DoListAddAfter(tableName, data, -1, nil)
//	if err != nil {
//		t.Errorf("DoListAddAfter failed %s", err.Error())
//		return
//	}
//	err = client.DoListAddAfter(tableName, data, -1, nil)
//	if err != nil {
//		t.Errorf("DoListAddAfter failed %s", err.Error())
//		return
//	}
//
//	// set ttl
//	opt := &option.TDROpt{
//		BatchTTL: []option.TTLInfo{{TTL:5000},{TTL:4000},{TTL:3000}},
//	}
//	dataSlice := []record.TdrTableSt{data,data,data}
//	indexs := []int32{0,1,2}
//	err = client.DoSetTTLBatch(tableName, dataSlice,indexs, opt)
//	if err != nil {
//		t.Errorf("DoSetTTLBatch fail, %s", err.Error())
//		return
//	}
//
//	//3 get ttl
//	opt = &option.TDROpt{}
//	err = client.DoGetTTLBatch(tableName, dataSlice, indexs, opt)
//	if err != nil {
//		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
//		return
//	}
//	fmt.Println(opt.BatchTTL)
//	for _, ttl := range opt.BatchTTL {
//		if ttl.TTL <=0 || ttl.TTL > 5000{
//			t.Errorf("ttl invalid %d", ttl.TTL)
//			return
//		}
//	}
//
//	//4 5s after
//	time.Sleep(5*time.Second)
//	opt = &option.TDROpt{}
//	err = client.DoGetTTLBatch(tableName, dataSlice,indexs, opt)
//	fmt.Println(opt.BatchResult)
//	if err == nil {
//		t.Errorf("DoGetTTLMust timeout, %v", opt.BatchResult)
//		return
//	}
//}
