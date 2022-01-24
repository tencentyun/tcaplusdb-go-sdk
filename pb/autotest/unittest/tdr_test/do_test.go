package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"testing"
	"time"
	"unsafe"
)
//case1 BatchInsert success
func TestTdrDo(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//插入记录成功
	uin := uint64(time.Now().UnixNano())
	data := newGenericTableRec()
	data.Uin = uin
	data.Name = fmt.Sprintf("%d", 2)
	data.Level = int32(2)
	data.Float_Score = float32(6.6)
	data.Double_Score = float64(8.8)
	data.Info = fmt.Sprintf("%d", 2)
	err = client.DoInsert(TestTableName, data, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}
	data.Level=10
	err = client.DoUpdate(TestTableName, data, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}
	data.Level=11
	err = client.DoReplace(TestTableName, data, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}

	err = client.DoDelete(TestTableName, data, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}
}

// 全部索引存在
func TestTdrListDo(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//list addafter
	uin := time.Now().UnixNano()
	key := *(*uint32)(unsafe.Pointer(&uin))
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = key
	data.Name = 255
	data.Level = uint32(11)
	data.Value1 = "value1"
	data.Value2 = "value2"
	opt := &option.TDROpt{
		ResultFlagForSuccess:option.TcaplusResultFlagAllNewValue,
	}
	idx, err := client.DoListAddAfter(TABLE_TRAVERSER_LIST, data, -1, opt)
	if err != nil {
		t.Errorf("DoListAddAfter failed %s", err.Error())
		return
	}
	fmt.Println(opt.Version)
	fmt.Println(idx)

	err = client.DoListGet(TABLE_TRAVERSER_LIST, data, idx, nil)
	if err != nil {
		t.Errorf("DoListGet failed %s", err.Error())
		return
	}
	fmt.Println(opt.Version)
	fmt.Println(data.Level)
	if data.Level != 11 {
		t.Errorf("DoListGet failed data.Level != 11")
		return
	}

	data.Level =12
	err = client.DoListReplace(TABLE_TRAVERSER_LIST, data, idx, nil)
	if err != nil {
		t.Errorf("DoListGet failed %s", err.Error())
		return
	}

	err = client.DoListDelete(TABLE_TRAVERSER_LIST, data, idx, nil)
	if err != nil {
		t.Errorf("DoListGet failed %s", err.Error())
		return
	}

	err = client.DoListGet(TABLE_TRAVERSER_LIST, data, idx, nil)
	fmt.Println(err.Error())
	if err == nil {
		t.Errorf("DoListGet must not exist")
		return
	}
}