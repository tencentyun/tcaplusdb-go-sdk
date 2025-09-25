package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
)

func deletePartKeyExample() {
	//申请tdr结构体并赋值Key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "dev"
	data.Name = "com"
	data.Inst_Max_Num = 29

	res, err := client.DoDeleteByPartKey(TableName, data, "Index_Gameid_Name", nil)
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
	fmt.Printf("DoDeleteByPartKey success total count %d,\n", len(res))

}
