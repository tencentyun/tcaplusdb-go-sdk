package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
)

func getPartKeyExample() {
	//申请tdr结构体并赋值Key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "dev"
	//data.Envdata = "oaasqomk"
	data.Name = "com"

	//设置返回记录数，不设置则全部返回
	opt := &option.TDROpt{Limit: 3}
	res, err := client.DoGetByPartKey(TableName, data, "Index_Gameid_Name", opt)
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
	fmt.Printf("DoGetByPartKey success total count %d,\n", len(res))

}
