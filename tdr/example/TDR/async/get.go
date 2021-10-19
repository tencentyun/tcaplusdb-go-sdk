package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"strconv"
)

func getExample() {
	//创建Get请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiGetReq)
	if err != nil {
		fmt.Printf("getExample NewRequest TcaplusApiGetReq failed %v\n", err.Error())
		return
	}
	fmt.Printf("getExample NewRequest TcaplusApiGetReq finish\n")
	//设置异步请求ID，异步请求通过ID让响应和请求对应起来
	req.SetAsyncId(667)
	//为request添加一条记录，（index只有在list表中支持，generic不校验）
	rec, err := req.AddRecord(0)
	if err != nil {
		fmt.Printf("getExample AddRecord failed %v\n", err.Error())
		return
	}
	fmt.Printf("getExample AddRecord finish\n")
	//申请tdr结构体并赋值Key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "dev"
	data.Envdata = "oa"
	data.Name = "com"
	//将tdr的数据设置到请求的记录中
	if err := rec.SetData(data); err != nil {
		fmt.Printf("SetData failed %v\n", err.Error())
		return
	}
	fmt.Printf("getExample SetData finish\n")
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	fmt.Printf("getExample send finish\n")
	resp, err := recvResponse(client)
	if err != nil {
		fmt.Printf("recv err %s\n", err.Error())
		return
	}
	//带回请求的异步ID
	fmt.Printf("getExample resp success, AsyncId:%d\n", resp.GetAsyncId())
	tcapluserr := resp.GetResult()
	if tcapluserr != 0 {
		fmt.Printf("response ret %s\n",
			"errCode: "+strconv.Itoa(tcapluserr)+", errMsg: "+terror.ErrorCodes[tcapluserr])
		return
	}
	//response中带有获取的记录
	fmt.Printf("getExample response success record count %d\n", resp.GetRecordCount())
	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			fmt.Printf("FetchRecord failed %s\n", err.Error())
			return
		}
		//通过GetData获取记录
		data := service_info.NewService_Info()
		if err := record.GetData(data); err != nil {
			fmt.Printf("record.GetData failed %s\n", err.Error())
			return
		}
		fmt.Printf("getExample response record data %+v, route: %s\n",
			data, string(data.Routeinfo[0:data.Routeinfo_Len]))
	}
}
