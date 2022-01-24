package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"testing"
	"time"
	"unsafe"
)

// 全部索引存在
func TestListAddAfterBatchSuccess(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//list addafter 10 record
	var dataSlice []record.TdrTableSt
	var indexs []int32
	uin := time.Now().UnixNano()
	key := *(*uint32)(unsafe.Pointer(&uin))
	for i := 0; i < 10; i++ {
		data := tcaplus_tb.NewTable_Traverser_List()
		data.Key = key
		data.Name = 255
		data.Level = uint32(i)
		data.Value1 = "value1"
		data.Value2 = "value2"
		dataSlice = append(dataSlice, data)
		indexs = append(indexs, -1)
	}

	opt := &option.TDROpt{}
	err = client.DoListAddAfterBatch(TABLE_TRAVERSER_LIST, dataSlice, indexs, opt)
	if err != nil {
		t.Errorf("DoListAddAfterBatch failed %s", err.Error())
		return
	}
	fmt.Println(indexs)
	fmt.Println(opt.BatchResult)
	fmt.Println(opt.Version)

	//list batch get
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = key
	data.Name = 255
	var indexs2 []int32
	for i := 0; i < 10; i++ {
		indexs2 = append(indexs2, int32(i))
	}
	recs, err := client.DoListGetBatch(TABLE_TRAVERSER_LIST, data, indexs2, nil)
	if err != nil {
		t.Errorf("DoListGetBatch failed %s", err.Error())
		return
	}

	for i, rec := range recs {
		ver := rec.GetVersion()
		index := rec.GetIndex()
		fmt.Println("version:", ver)
		fmt.Println("index:", index)
		if ver != 10 {
			t.Errorf("Version invalid %d", ver)
			return
		}
		if index != int32(i) {
			t.Errorf("index invalid %d", index)
			return
		}
		data = tcaplus_tb.NewTable_Traverser_List()
		err := rec.GetData(data)
		if err != nil {
			t.Errorf("GetData failed %s", err.Error())
			return
		}
		fmt.Println(data)
		if data.Level != uint32(i) {
			t.Errorf("data.Level invalid %d", data.Level)
			return
		}
	}

	//list get all
	data = tcaplus_tb.NewTable_Traverser_List()
	data.Key = key
	data.Name = 255
	recs, err = client.DoListGetAll(TABLE_TRAVERSER_LIST, data, nil)
	if err != nil {
		t.Errorf("DoListGetBatch failed %s", err.Error())
		return
	}

	for i, rec := range recs {
		ver := rec.GetVersion()
		index := rec.GetIndex()
		fmt.Println("version:", ver)
		fmt.Println("index:", index)
		if ver != 10 {
			t.Errorf("Version invalid %d", ver)
			return
		}
		if index != int32(i) {
			t.Errorf("index invalid %d", index)
			return
		}
		data = tcaplus_tb.NewTable_Traverser_List()
		err := rec.GetData(data)
		if err != nil {
			t.Errorf("GetData failed %s", err.Error())
			return
		}
		fmt.Println(data)
		if data.Level != uint32(i) {
			t.Errorf("data.Level invalid %d", data.Level)
			return
		}
	}
	//list delete all
	err = client.DoListDeleteAll(TABLE_TRAVERSER_LIST, data, nil)
	if err != nil {
		t.Errorf("DoListDeleteAll failed %s", err.Error())
		return
	}

	//empty
	_, err = client.DoListGetAll(TABLE_TRAVERSER_LIST, data, nil)
	fmt.Println(err)
	if err == nil {
		t.Errorf("DoListGetAll must empty")
		return
	}
}

func TestListAddAfterBatchVersionFail(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//list addafter 10 record
	var dataSlice []record.TdrTableSt
	var indexs []int32
	uin := time.Now().UnixNano()
	key := *(*uint32)(unsafe.Pointer(&uin))
	for i := 0; i < 10; i++ {
		data := tcaplus_tb.NewTable_Traverser_List()
		data.Key = key
		data.Name = 255
		data.Level = uint32(i)
		data.Value1 = "value1"
		data.Value2 = "value2"
		dataSlice = append(dataSlice, data)
		indexs = append(indexs, -1)
	}

	opt := &option.TDROpt{}
	err = client.DoListAddAfterBatch(TABLE_TRAVERSER_LIST, dataSlice, indexs, opt)
	if err != nil {
		t.Errorf("DoListAddAfterBatch failed %s", err.Error())
		return
	}
	fmt.Println(indexs)
	fmt.Println(opt.BatchResult)
	fmt.Println(opt.Version)

	//succ
	opt = &option.TDROpt{
		VersionPolicy: option.CheckDataVersionAutoIncrease,
		Version:       10,
	}
	indexs = nil
	for i := 0; i < 10; i++ {
		indexs = append(indexs, -1)
	}
	err = client.DoListAddAfterBatch(TABLE_TRAVERSER_LIST, dataSlice, indexs, opt)
	if err != nil {
		t.Errorf("DoListAddAfterBatch failed %s", err.Error())
		return
	}
	fmt.Println(indexs)
	fmt.Println(opt.BatchResult)
	fmt.Println(opt.Version)

	//fail
	opt = &option.TDROpt{
		VersionPolicy: option.CheckDataVersionAutoIncrease,
		Version:       100,
	}
	indexs = nil
	for i := 0; i < 10; i++ {
		indexs = append(indexs, -1)
	}
	err = client.DoListAddAfterBatch(TABLE_TRAVERSER_LIST, dataSlice, indexs, opt)
	if err == nil {
		t.Errorf("DoListAddAfterBatch must version failed")
		return
	}
	fmt.Println(indexs)
	fmt.Println(opt.BatchResult)
	fmt.Println(opt.Version)
}
