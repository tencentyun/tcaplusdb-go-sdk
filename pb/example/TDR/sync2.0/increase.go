package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
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
	data.Float_Score = float32(6.6)
	data.Double_Score = float64(8.8)
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
			option.IncFieldInfo{
				FieldName: "float_score",
				IncData:   float32(6), //+2
				Operation: cmd.TcaplusApiOpPlus,
			},
			option.IncFieldInfo{
				FieldName: "double_score",
				IncData:   float64(8), //+8
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

	if data.Float_Score != 12.6 {
		fmt.Printf("data.Float_Score invalid %v", data.Float_Score)
		return
	}

	if data.Double_Score != 16.8 {
		fmt.Printf("data.Double_Score invalid %v", data.Double_Score)
		return
	}
}
