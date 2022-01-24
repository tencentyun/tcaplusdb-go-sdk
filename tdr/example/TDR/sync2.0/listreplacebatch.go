package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
)

func listReplaceBatchExample() {
	//list addafter 10 record
	var dataSlice []record.TdrTableSt
	var indexs []int32
	for i := 0; i < 10; i++ {
		data := tcaplus_tb.NewTable_Traverser_List()
		data.Key = 1
		data.Name = 255
		data.Level = uint32(i)
		data.Value1 = "replace"
		data.Value2 = "replace"
		dataSlice = append(dataSlice, data)
		indexs = append(indexs, int32(i))
	}

	opt := &option.TDROpt{}
	err := client.DoListReplaceBatch(TABLE_TRAVERSER_LIST, dataSlice, indexs, opt)
	if err != nil {
		fmt.Printf("DoListReplaceBatch failed %s", err.Error())
		return
	}

	fmt.Println(indexs)
	fmt.Println(opt.BatchResult)
	fmt.Println(opt.Version)
	fmt.Println("DoListReplaceBatch success")
}
