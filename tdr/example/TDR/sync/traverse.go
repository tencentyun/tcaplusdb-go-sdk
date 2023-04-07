package main

import (
	"fmt"
	"time"

	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

func traverseExample() {
	tra := client.GetTraverser(ZoneId, TableName)
	defer tra.Stop()

	tra.SetFieldNames([]string{"filterdata", "updatetime"})
	tra.SetLimit(10)

	resps, err := client.DoTraverse(tra, 60*time.Second)
	if err != nil {
		fmt.Println("DoTraverse:", err.Error())
		return
	}
	fmt.Println("rsp count ", len(resps))

	for _, resp := range resps {
		if err := resp.GetResult(); err != 0 {
			fmt.Printf("resp.GetResult err %s\n", terror.GetErrMsg(err))
			return
		}

		fmt.Println("one rsp record count ", resp.GetRecordCount())
		for i := 0; i < resp.GetRecordCount(); i++ {
			record, err := resp.FetchRecord()
			if err != nil {
				fmt.Printf("FetchRecord failed %s", err.Error())
				return
			}

			data := service_info.NewService_Info()
			if err := record.GetData(data); err != nil {
				fmt.Printf("record.GetData failed %s\n", err.Error())
				return
			}
			fmt.Println("record service_info : ", data)
		}
	}
}
