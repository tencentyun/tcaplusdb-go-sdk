package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"time"
)

func batchInsertExample() {
	var dataSlice []record.TdrTableSt
	for i := 0; i < 10; i++ {
		data := service_info.NewService_Info()
		data.Gameid = "dev"
		data.Envdata = "oa"
		data.Name = fmt.Sprintf("%d", i)
		data.Filterdata = time.Now().Format("2006-01-02T15:04:05.000000Z")
		data.Updatetime = uint64(time.Now().UnixNano())
		data.Inst_Max_Num = 2
		data.Inst_Min_Num = 3
		route := "test"
		data.Routeinfo_Len = uint32(len(route))
		data.Routeinfo = []byte(route)
		dataSlice = append(dataSlice, data)
	}

	if err := client.DoBatchInsert(TableName, dataSlice, nil); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Batch Insert SUCCESS")
}
