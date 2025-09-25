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

//case insert update replace支持ttl
func TestWriteWithTTLSuccess(t *testing.T) {
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

	///////1 insert with ttl 成功
	data1 := newGenericTableRec()
	data1.Name = "insertWithTtl1"
	opt := &option.TDROpt{
		TTL: &option.TTLInfo{
			TTL: 5000,
		},
	}
	err = client.DoInsert(TestTableName, data1, opt)
	if err != nil {
		t.Errorf("DoInsert fail, %s", err.Error())
		return
	}

	////1 Get ttl
	dataSlice := []record.TdrTableSt{data1}
	opt = &option.TDROpt{}
	err = client.DoGetTTLBatch(TestTableName, dataSlice, nil, opt)
	if err != nil {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}
	if len(opt.BatchTTL) == 0 {
		t.Errorf("opt.BatchTTL empty")
		return
	}
	if opt.BatchTTL[0].TTL < 0 || opt.BatchTTL[0].TTL > 5000 {
		t.Errorf("opt.BatchTTL %d invalid", opt.BatchTTL[0].TTL)
		return
	}
	fmt.Println("after insert ttl", opt.BatchTTL[0].TTL)

	///2 update ttl
	opt = &option.TDROpt{
		TTL: &option.TTLInfo{
			TTL: 10000,
		},
	}
	err = client.DoUpdate(TestTableName, data1, opt)
	if err != nil {
		t.Errorf("DoUpdate fail, %s", err.Error())
		return
	}

	////2 Get ttl
	opt = &option.TDROpt{}
	err = client.DoGetTTLBatch(TestTableName, dataSlice, nil, opt)
	if err != nil {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}
	if len(opt.BatchTTL) == 0 {
		t.Errorf("opt.BatchTTL empty")
		return
	}
	if opt.BatchTTL[0].TTL < 5000 || opt.BatchTTL[0].TTL > 10000 {
		t.Errorf("opt.BatchTTL %d invalid", opt.BatchTTL[0].TTL)
		return
	}
	fmt.Println("after update ttl", opt.BatchTTL[0].TTL)

	///3 replace ttl
	opt = &option.TDROpt{
		TTL: &option.TTLInfo{
			TTL: 15000,
		},
	}
	err = client.DoReplace(TestTableName, data1, opt)
	if err != nil {
		t.Errorf("DoUpdate fail, %s", err.Error())
		return
	}

	////3 Get ttl
	opt = &option.TDROpt{}
	err = client.DoGetTTLBatch(TestTableName, dataSlice, nil, opt)
	if err != nil {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}
	if len(opt.BatchTTL) == 0 {
		t.Errorf("opt.BatchTTL empty")
		return
	}
	if opt.BatchTTL[0].TTL < 10000 || opt.BatchTTL[0].TTL > 15000 {
		t.Errorf("opt.BatchTTL %d invalid", opt.BatchTTL[0].TTL)
		return
	}
	fmt.Println("after replace ttl", opt.BatchTTL[0].TTL)
}
