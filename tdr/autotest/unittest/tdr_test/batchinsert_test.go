package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"strings"
	"testing"
	"time"
)

//case1 BatchInsert success
func TestBatchInsertSuccess(t *testing.T) {
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

	err = client.DoBatchInsert(TestTableName, dataSlice, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}

	//2 batchGet success
	var dataSlice2 []record.TdrTableSt
	for i := 0; i < 10; i++ {
		data := tcaplus_tb.NewTable_Generic()
		data.Uin = uin
		data.Name = fmt.Sprintf("%d", i)
		data.Key3 = "key3"
		data.Key4 = "key4"
		dataSlice2 = append(dataSlice2, data)
	}

	err = client.DoBatchGet(TestTableName, dataSlice2, nil)
	if err != nil {
		t.Errorf("DoBatchGet fail, %s", err.Error())
		return
	}

	for i, data := range dataSlice2 {
		if data.(*tcaplus_tb.Table_Generic).Level != int32(i) {
			t.Errorf("DoBatchGet fail, %+v", data)
			return
		}
	}
}
//case1 BatchInsert success
func TestBatchInsertSuccess_1024(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//批量插入1024 条记录成功
	var dataSlice []record.TdrTableSt
	uin := uint64(time.Now().UnixNano())
	for i := 0; i < 1024; i++ {
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

	//2 batchGet success
	var dataSlice2 []record.TdrTableSt
	for i := 0; i < 1024; i++ {
		data := tcaplus_tb.NewTable_Generic()
		data.Uin = uin
		data.Name = fmt.Sprintf("%d", i)
		data.Key3 = "key3"
		data.Key4 = "key4"
		dataSlice2 = append(dataSlice2, data)
	}

	err = client.DoBatchGet(TestTableName, dataSlice2, nil)
	if err != nil {
		t.Errorf("DoBatchGet fail, %s", err.Error())
		return
	}

	for i, data := range dataSlice2 {
		if data.(*tcaplus_tb.Table_Generic).Level != int32(i) {
			t.Errorf("DoBatchGet fail, %+v", data)
			return
		}
	}
}
//case1 BatchInsert success
func TestBatchInsertSuccess_1025(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//批量插入1025 条记录成功
	var dataSlice []record.TdrTableSt
	uin := uint64(time.Now().UnixNano())
	for i := 0; i < 1025; i++ {
		data := newGenericTableRec()
		data.Uin = uin
		data.Name = fmt.Sprintf("%d", i)
		data.Level = int32(i)
		data.Info = fmt.Sprintf("%d", i)
		dataSlice = append(dataSlice, data)
	}

	err = client.DoBatchInsert(TestTableName, dataSlice, nil)
	if !strings.Contains(err.Error(),"-4126") {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}

	//2 batchGet success
	var dataSlice2 []record.TdrTableSt
	for i := 0; i < 1025; i++ {
		data := tcaplus_tb.NewTable_Generic()
		data.Uin = uin
		data.Name = fmt.Sprintf("%d", i)
		data.Key3 = "key3"
		data.Key4 = "key4"
		dataSlice2 = append(dataSlice2, data)
	}

	err = client.DoBatchGet(TestTableName, dataSlice2, nil)
	if !strings.Contains(err.Error(),"-4126") {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}
}