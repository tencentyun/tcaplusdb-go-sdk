package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
)

func listAddAfterBatchExample() {
	//list addafter 10 record
	var dataSlice []record.TdrTableSt
	var indexs []int32
	for i := 0; i < 10; i++ {
		data := tcaplus_tb.NewTable_Traverser_List()
		data.Key = 1
		data.Name = 255
		data.Level = uint32(i)
		data.Value1 = "value1"
		data.Value2 = "value2"
		dataSlice = append(dataSlice, data)
		indexs = append(indexs, -1)
	}

	opt := &option.TDROpt{}
	err := client.DoListAddAfterBatch(TABLE_TRAVERSER_LIST, dataSlice, indexs, opt)
	if err != nil {
		fmt.Printf("DoListAddAfterBatch failed %s", err.Error())
		return
	}

	fmt.Println(indexs)
	fmt.Println(opt.BatchResult)
	fmt.Println(opt.Version)
	fmt.Println("DoListAddAfterBatch success")
}
