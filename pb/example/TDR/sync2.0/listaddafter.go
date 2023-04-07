package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
)

func listAddAfterExample() {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255
	data.Level = 1
	data.Value1 = "value1"
	data.Value2 = "value2"

	//TCAPLUS_API_LIST_PRE_FIRST_INDEX(-2)：新元素插入在第一个元素之前
	//TCAPLUS_API_LIST_LAST_INDEX(-1)：新元素插入在最后一个元素之后
	_, err := client.DoListAddAfter(TABLE_TRAVERSER_LIST, data, -1, nil)
	if err != nil {
		fmt.Printf("DoListAddAfter failed %s", err.Error())
		return
	}

	fmt.Println("DoListAddAfter success")
}
