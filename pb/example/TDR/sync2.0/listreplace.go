package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
)

func listReplaceExample() {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255
	data.Level = 1
	data.Value1 = "value1"
	data.Value2 = "value2"

	//更新index为1的元素
	err := client.DoListReplace(TABLE_TRAVERSER_LIST, data, 1, nil)
	if err != nil {
		fmt.Printf("DoListReplace failed %s", err.Error())
		return
	}

	fmt.Println("DoListReplace success")
}
