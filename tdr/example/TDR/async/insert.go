package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

func insertExample() {
	//创建insert请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		fmt.Printf("NewRequest TcaplusApiInsertReq failed %v\n", err.Error())
		return
	}
	fmt.Printf("insertExample NewRequest TcaplusApiInsertReq finish\n")
	//设置异步请求ID，异步请求通过ID让响应和请求对应起来
	req.SetAsyncId(666)
	//为request添加一条记录，（index只有在list表中支持，generic不校验）
	rec, err := req.AddRecord(0)
	if err != nil {
		fmt.Printf("AddRecord failed %v\n", err.Error())
		return
	}
	fmt.Printf("insertExample AddRecord finish\n")
	//申请tdr结构体并赋值，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "dev"
	data.Envdata = "oa"
	data.Name = "com"
	data.Filterdata = time.Now().Format("2006-01-02T15:04:05.000000Z")
	data.Updatetime = uint64(time.Now().UnixNano())
	data.Inst_Max_Num = 2
	data.Inst_Min_Num = 3
	//数组类型为slice需要准确赋值长度，与refer保持一致
	route := "test"
	data.Routeinfo_Len = uint32(len(route))
	data.Routeinfo = []byte(route)
	//将tdr的数据设置到请求的记录中
	if err := rec.SetData(data); err != nil {
		fmt.Printf("SetData failed %v\n", err.Error())
		return
	}
	//发送请求
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	fmt.Printf("insertExample send finish\n")
	//接收响应
	resp, err := recvResponse(client)
	if err != nil {
		fmt.Printf("recv err %s\n", err.Error())
		return
	}
	//带回请求的异步ID
	fmt.Printf("insertExample resp success, AsyncId:%d\n", resp.GetAsyncId())
	tcapluserr := resp.GetResult()
	if tcapluserr != 0 {
		fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
		return
	}
}
