package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
)

func listDeleteExample() {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255

	//删除指定index位置的元素
	err := client.DoListDelete(TABLE_TRAVERSER_LIST, data, 1, nil)
	if err != nil {
		fmt.Printf("DoListDelete failed %s", err.Error())
		return
	}
	fmt.Println("DoListDelete success")
}
