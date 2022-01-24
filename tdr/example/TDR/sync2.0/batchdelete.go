package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
)

func batchDeleteExample() {
	//result flag + version
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
		dataSlice = append(dataSlice, data)
		opt.BatchVersion = append(opt.BatchVersion, 3) //校验version
	}

	if err := client.DoBatchDelete(TableName, dataSlice, opt); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	for _, data := range dataSlice {
		fmt.Printf("%+v", data)
	}

	fmt.Println("Batch Delete SUCCESS")
}
