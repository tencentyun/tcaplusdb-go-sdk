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

//case1 BatchDelete  success
func TestBatchDeleteSuccess(t *testing.T) {
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

	for i := 0; i < len(dataSlice); i++ {
		opt.BatchVersion = append(opt.BatchVersion, 1)
	}
	err = client.DoBatchDelete(TestTableName, dataSlice, opt)
	if err != nil {
		t.Errorf("DoBatchDelete fail, %s", err.Error())
		return
	}

	fmt.Println(opt.BatchVersion)
	fmt.Println(opt.BatchResult)
	for i, data := range dataSlice {
		if data.(*tcaplus_tb.Table_Generic).Level != int32(i) || data.(*tcaplus_tb.Table_Generic).Info != fmt.Sprintf("%d", i) {
			t.Errorf("DoBatchDelete fail, %+v", data)
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
	fmt.Println(opt.BatchVersion)
	fmt.Println(opt.BatchResult)
	if err == nil {
		t.Errorf("DoBatchGet fail, must err not exist")
		return
	}
}

func TestBatchDeleteVersionFail(t *testing.T) {
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

	for i := 0; i < len(dataSlice); i++ {
		opt.BatchVersion = append(opt.BatchVersion, 10)
	}
	err = client.DoBatchDelete(TestTableName, dataSlice, opt)
	if err == nil {
		t.Errorf("DoBatchDelete version fail, must version error")
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
