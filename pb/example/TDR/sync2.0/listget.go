package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
)

func listgetExample() {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255

	//查询指定index位置的元素
	err := client.DoListGet(TABLE_TRAVERSER_LIST, data, 1, nil)
	if err != nil {
		fmt.Printf("DoListGet failed %s", err.Error())
		return
	}
	fmt.Println("DoListAddAfter success")
}
