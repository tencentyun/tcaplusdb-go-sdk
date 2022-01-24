package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
)

func listGetBatchExample() {
	//list batch get 10 条记录
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255
	var indexs []int32
	for i := 0; i < 10; i++ {
		indexs = append(indexs, int32(i))
	}
	recs, err := client.DoListGetBatch(TABLE_TRAVERSER_LIST, data, indexs, nil)
	if err != nil {
		fmt.Println("DoListGetBatch failed,", err.Error())
		return
	}

	for _, rec := range recs {
		ver := rec.GetVersion()
		index := rec.GetIndex()
		fmt.Println("version:", ver)
		fmt.Println("index:", index)

		err := rec.GetData(data)
		if err != nil {
			fmt.Println("GetData failed", err.Error())
			return
		}
		fmt.Println(data)
	}
	fmt.Println("listGetBatchExample success")
}
