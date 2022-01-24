package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
)

func deleteExample() {
	//申请tdr结构体并赋值key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "dev"
	data.Envdata = "oa"
	data.Name = "com"
	err := client.DoDelete(TableName, data, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("delete success")
}
