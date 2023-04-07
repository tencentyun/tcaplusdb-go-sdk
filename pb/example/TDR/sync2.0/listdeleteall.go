package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
)

func listDeleteAllExample() {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255
	err := client.DoListDeleteAll(TABLE_TRAVERSER_LIST, data, nil)
	if err != nil {
		fmt.Println("DoListDeleteAll failed,", err.Error())
		return
	}

	fmt.Println("DoListDeleteAll success")
}
