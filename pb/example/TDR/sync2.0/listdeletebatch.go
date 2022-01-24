package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
)

func listDeleteBatchExample() {
	//list batch get 10 条记录
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255
	var indexs []int32
	for i := 0; i < 10; i++ {
		indexs = append(indexs, int32(i))
	}
	_, err := client.DoListDeleteBatch(TABLE_TRAVERSER_LIST, data, indexs, nil)
	if err != nil {
		fmt.Println("DoListGetBatch failed,", err.Error())
		return
	}

	fmt.Println("DoListDeleteBatch success")
}
