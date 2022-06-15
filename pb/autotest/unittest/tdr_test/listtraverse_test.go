package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"testing"
	"time"
	"unsafe"
)

// 全部索引存在
func TestListTraverse(t *testing.T) {
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

	//list travser
	tra := client.GetListTraverser(cfg.ApiConfig.ZoneId, TABLE_TRAVERSER_LIST)
	defer tra.Stop()

	//tra.SetFieldNames([]string{"level", "count", "info"})
	tra.SetLimit(10)

	resps, err := client.DoTraverse(tra, 60*time.Second)
	if err != nil {
		t.Errorf("RecvResponse fail, %s", err.Error())
		return
	}
	fmt.Println(len(resps))
	for _, resp := range resps {
		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
			return
		}

		for i := 0; i < resp.GetRecordCount(); i++ {
			record, err := resp.FetchRecord()
			if err != nil {
				t.Errorf("FetchRecord failed %s", err.Error())
				return
			}
			data := tcaplus_tb.NewTable_Traverser_List()
			err = record.GetData(data)
			if err != nil {
				t.Errorf("GetData failed %s", err.Error())
				return
			}
			fmt.Println(data)
		}
	}
}

