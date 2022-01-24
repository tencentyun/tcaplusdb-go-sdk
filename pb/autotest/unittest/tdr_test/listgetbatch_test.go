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

// 全部索引存在
const TABLE_TRAVERSER_LIST = "table_traverser_list"

func TestListGetBatchSuccess(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//list addafter 10 record
	data := tcaplus_tb.NewTable_Traverser_List()
	uin := time.Now().UnixNano()
	key := *(*uint32)(unsafe.Pointer(&uin))
	for i := 0; i < 10; i++ {
		data.Key = key
		data.Name = 255
		data.Level = uint32(i)
		data.Value1 = "value1"
		data.Value2 = "value2"
		_, err = client.DoListAddAfter(TABLE_TRAVERSER_LIST, data, -1, nil)
		if err != nil {
			t.Errorf("DoListAddAfter failed %s", err.Error())
			return
		}
	}

	//list batch get
	data = tcaplus_tb.NewTable_Traverser_List()
	data.Key = key
	data.Name = 255
	var indexs []int32
	for i := 0; i < 10; i++ {
		indexs = append(indexs, int32(i))
	}
	recs, err := client.DoListGetBatch(TABLE_TRAVERSER_LIST, data, indexs, nil)
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

	//部分字段获取
	opt := &option.TDROpt{
		FieldNames: []string{"level"},
	}
	recs, err = client.DoListGetBatch(TABLE_TRAVERSER_LIST, data, indexs, opt)
	if err != nil {
		t.Errorf("DoListGetBatch failed %s", err.Error())
		return
	}

	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
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
		if len(data.Value1) != 0 {
			t.Errorf("data.Value1 invalid %s", data.Value1)
			return
		}
		if len(data.Value2) != 0 {
			t.Errorf("data.Value2 invalid %s", data.Value2)
			return
		}
	}
}
