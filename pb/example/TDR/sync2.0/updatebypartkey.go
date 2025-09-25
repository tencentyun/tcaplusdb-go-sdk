package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
)

func updatePartKeyExample() {
	//申请tdr结构体并赋值Key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "dev"
	data.Name = "com"
	data.Inst_Max_Num = 29

	//设置更新的字段
	opt := &option.TDROpt{
		FieldNames: []string{"inst_max_num"},
	}
	res, err := client.DoUpdateByPartKey(TableName, data, "Index_Gameid_Name", opt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, rec := range res {
		err = rec.GetData(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(data)
	}
	fmt.Printf("DoUpdateByPartKey success total count %d,\n", len(res))

}
