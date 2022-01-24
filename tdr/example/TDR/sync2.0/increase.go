package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"time"
)

func increaseExample() {
	//插入记录成功
	data := tcaplus_tb.NewTable_Generic()
	data.Uin = uint64(time.Now().UnixNano())
	data.Name = "GoUnitTest"
	data.Key3 = "key3"
	data.Key4 = "key4"
	data.Info = "info"
	data.Name = fmt.Sprintf("%d", 2)
	data.Level = int32(2)
	data.Info = fmt.Sprintf("%d", 2)
	err := client.DoInsert("table_generic", data, nil)
	if err != nil {
		fmt.Printf("DoBatchInsert fail, %s", err.Error())
		return
	}

	//自增记录
	opt := &option.TDROpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllNewValue,
		IncField: []option.IncFieldInfo{
			option.IncFieldInfo{
				FieldName: "level",
				IncData:   int32(2), //+2
				Operation: cmd.TcaplusApiOpPlus,
			},
		},
	}
	err = client.DoIncrease("table_generic", data, opt)
	if err != nil {
		fmt.Printf("DoIncrease fail, %s", err.Error())
		return
	}

	if data.Level != 4 {
		fmt.Printf("data.Level invalid %v", data.Level)
		return
	}
}
