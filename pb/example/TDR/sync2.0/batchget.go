package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
)

func batchGetExample() {
	var dataSlice []record.TdrTableSt
	for i := 0; i < 10; i++ {
		data := service_info.NewService_Info()
		data.Gameid = "dev"
		data.Envdata = "oa"
		data.Name = fmt.Sprintf("%d", i)
		dataSlice = append(dataSlice, data)
	}
	opt := &option.TDROpt{}
	if err := client.DoBatchGet(TableName, dataSlice, opt); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(opt.BatchVersion)
	fmt.Println(opt.BatchResult)
	for _, data := range dataSlice {
		fmt.Printf("%+v", data)
	}
	fmt.Println("Batch Get SUCCESS")
}
