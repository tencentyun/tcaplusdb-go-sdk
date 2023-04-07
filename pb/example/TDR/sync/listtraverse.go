package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"time"
)

func ListTraverseExample() {
	tra := client.GetListTraverser(ZoneId, TABLE_TRAVERSER_LIST)
	defer tra.Stop()

	tra.SetFieldNames([]string{"level", "value1"})
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

			data := tcaplus_tb.NewTable_Traverser_List()
			if err := record.GetData(data); err != nil {
				fmt.Printf("record.GetData failed %s\n", err.Error())
				return
			}
			fmt.Println("record : ", data)
		}
	}
}
