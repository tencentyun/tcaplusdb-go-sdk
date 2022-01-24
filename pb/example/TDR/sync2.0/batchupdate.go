package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"time"
)

func batchUpdateExample() {
	//result flag + version success
	opt := &option.TDROpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
		VersionPolicy:        option.CheckDataVersionAutoIncrease,
	}

	var dataSlice []record.TdrTableSt
	for i := 0; i < 10; i++ {
		data := service_info.NewService_Info()
		data.Gameid = "dev"
		data.Envdata = "oa"
		data.Name = fmt.Sprintf("%d", i)
		data.Filterdata = "update"
		data.Updatetime = uint64(time.Now().UnixNano())
		data.Inst_Max_Num = 2
		data.Inst_Min_Num = 3
		route := "test"
		data.Routeinfo_Len = uint32(len(route))
		data.Routeinfo = []byte(route)
		dataSlice = append(dataSlice, data)
		opt.BatchVersion = append(opt.BatchVersion, 2) //校验version
	}

	if err := client.DoBatchUpdate(TableName, dataSlice, opt); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	for _, data := range dataSlice {
		fmt.Printf("%+v", data)
	}

	fmt.Println("Batch Update SUCCESS")
}
