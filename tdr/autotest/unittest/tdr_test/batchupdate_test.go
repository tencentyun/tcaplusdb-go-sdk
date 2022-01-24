package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"testing"
	"time"
)

//case1 BatchUpdate success
func TestBatchUpdateSuccess(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//1 批量插入10 条记录成功
	var dataSlice []record.TdrTableSt
	uin := uint64(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		data := newGenericTableRec()
		data.Uin = uin
		data.Name = fmt.Sprintf("%d", i)
		data.Level = int32(i)
		data.Info = fmt.Sprintf("%d", i)
		dataSlice = append(dataSlice, data)
	}

	err = client.DoBatchInsert(TestTableName, dataSlice, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}

	//2 result flag + version success
	opt := &option.TDROpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
		VersionPolicy:        option.CheckDataVersionAutoIncrease,
	}

	for _, data := range dataSlice {
		data.(*tcaplus_tb.Table_Generic).Info = "update"
		opt.BatchVersion = append(opt.BatchVersion, 1)
	}
	err = client.DoBatchUpdate(TestTableName, dataSlice, opt)
	if err != nil {
		t.Errorf("DoBatchUpdate fail, %s", err.Error())
		return
	}

	fmt.Println(opt.BatchVersion)
	fmt.Println(opt.BatchResult)
	for i, data := range dataSlice {
		if data.(*tcaplus_tb.Table_Generic).Level != int32(i) || data.(*tcaplus_tb.Table_Generic).Info != fmt.Sprintf("%d", i) {
			t.Errorf("DoBatchUpdate fail, %+v", data)
			return
		}
	}

	//3 batchGet success
	var dataSlice2 []record.TdrTableSt
	for i := 0; i < 10; i++ {
		data := tcaplus_tb.NewTable_Generic()
		data.Uin = uin
		data.Name = fmt.Sprintf("%d", i)
		data.Key3 = "key3"
		data.Key4 = "key4"
		dataSlice2 = append(dataSlice2, data)
	}

	opt = &option.TDROpt{}
	err = client.DoBatchGet(TestTableName, dataSlice2, opt)
	if err != nil {
		t.Errorf("DoBatchGet fail, %s", err.Error())
		return
	}
	fmt.Println(opt.BatchVersion)
	fmt.Println(opt.BatchResult)
	for i, data := range dataSlice2 {
		if data.(*tcaplus_tb.Table_Generic).Level != int32(i) || data.(*tcaplus_tb.Table_Generic).Info != "update" {
			t.Errorf("DoBatchGet fail, %+v", data)
			return
		}
	}
}

func TestBatchUpdateVersionFail(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//批量插入10 条记录成功
	var dataSlice []record.TdrTableSt
	uin := uint64(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		data := newGenericTableRec()
		data.Uin = uin
		data.Name = fmt.Sprintf("%d", i)
		data.Level = int32(i)
		data.Info = fmt.Sprintf("%d", i)
		dataSlice = append(dataSlice, data)
	}

	err = client.DoBatchReplace(TestTableName, dataSlice, nil)
	if err != nil {
		t.Errorf("DoBatchReplace fail, %s", err.Error())
		return
	}

	//2 result flag + version
	opt := &option.TDROpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllNewValue,
		VersionPolicy:        option.CheckDataVersionAutoIncrease,
	}

	for _, data := range dataSlice {
		data.(*tcaplus_tb.Table_Generic).Info = "replace"
		opt.BatchVersion = append(opt.BatchVersion, 10)
	}
	err = client.DoBatchUpdate(TestTableName, dataSlice, opt)
	if err == nil {
		t.Errorf("DoBatchUpdate version fail, must version error")
		return
	}

	fmt.Println(opt.BatchVersion)
	fmt.Println(opt.BatchResult)
	for i, _ := range dataSlice {
		if opt.BatchResult[i] == nil {
			t.Errorf("DoBatchUpdate fail, must version error  %+v", opt.BatchResult)
			return
		}
	}
}
