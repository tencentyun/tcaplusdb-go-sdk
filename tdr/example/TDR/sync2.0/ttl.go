package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/record"
	"time"
)

func ttlExample() {
	//申请tdr结构体并赋值Key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "dev"
	data.Envdata = "oa"
	data.Name = "com"

	////2 set ttl
	opt := &option.TDROpt{
		BatchTTL: []option.TTLInfo{{TTL: 5000}},
	}
	dataSlice := []record.TdrTableSt{data}
	err := client.DoSetTTLBatch(TableName, dataSlice, nil, opt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//3 get ttl
	opt = &option.TDROpt{}
	err = client.DoGetTTLBatch(TableName, dataSlice, nil, opt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(opt.BatchTTL)
	fmt.Println("ttl success")
}

// 单命令insert update replace也可以设置ttl
func ttlWithInsertExample() {
	//申请tdr结构体并赋值，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "ndev"
	data.Envdata = "noa"
	data.Name = "com"
	data.Filterdata = time.Now().Format("2006-01-02T15:04:05.000000Z")
	data.Updatetime = uint64(time.Now().UnixNano())
	data.Inst_Max_Num = 2
	data.Inst_Min_Num = 3
	//数组类型为slice需要准确赋值长度，与refer保持一致
	route := "test"
	data.Routeinfo_Len = uint32(len(route))
	data.Routeinfo = []byte(route)

	opt := &option.TDROpt{
		TTL: &option.TTLInfo{
			TTL: 5000,
		},
	}
	if err := client.DoInsert(TableName, data, opt); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Insert SUCCESS")
}
