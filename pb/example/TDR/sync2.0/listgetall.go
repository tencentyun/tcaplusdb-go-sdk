package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
)

func listGetAllExample() {
	//list batch get 10 条记录
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255
	opt := &option.TDROpt{
		MultiFlag: 1,
	}
	recs, err := client.DoListGetAll(TABLE_TRAVERSER_LIST, data, opt)
	if err != nil {
		fmt.Println("DoListGetAll failed,", err.Error())
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
	fmt.Println("DoListGetAll success")
}
