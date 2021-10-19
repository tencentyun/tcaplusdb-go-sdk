package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
)

func deleteExample() {
	//创建Delete请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiDeleteReq)
	if err != nil {
		fmt.Printf("NewRequest TcaplusApiDeleteReq failed %v\n", err.Error())
		return
	}
	fmt.Printf("deleteExample NewRequest TcaplusApiDeleteReq finish\n")
	//设置异步请求ID，异步请求通过ID让响应和请求对应起来
	req.SetAsyncId(670)
	//设置结果标记位，删除成功后，返回tcaplus端的旧数据，默认为0
	if err := req.SetResultFlag(3); err != nil {
		fmt.Printf("SetResultFlag failed %v\n", err.Error())
		return
	}
	fmt.Printf("deleteExample SetResultFlag finish\n")
	//为request添加一条记录，（index只有在list表中支持，generic不校验）
	rec, err := req.AddRecord(0)
	if err != nil {
		fmt.Printf("AddRecord failed %v\n", err.Error())
		return
	}
	fmt.Printf("deleteExample AddRecord finish\n")
	//申请tdr结构体并赋值key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "gid_1"
	data.Envdata = "oagfgadsf"
	data.Name = "com"
	data.Expansion = "fds"
	//将tdr的数据设置到请求的记录中
	if err := rec.SetData(data); err != nil {
		fmt.Printf("SetData failed %v\n", err.Error())
		return
	}
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	fmt.Printf("deleteExample send finish\n")
	//recv response
	resp, err := recvResponse(client)
	if err != nil {
		fmt.Printf("recv err %s\n", err.Error())
		return
	}
	//带回请求的异步ID
	fmt.Printf("deleteExample resp success, AsyncId:%d\n", resp.GetAsyncId())
	tcapluserr := resp.GetResult()
	if tcapluserr != 0 {
		fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
		return
	}
	//response中带有获取的旧记录
	fmt.Printf("deleteExample response success record count %d\n", resp.GetRecordCount())
	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			fmt.Printf("FetchRecord failed %s\n", err.Error())
			return
		}
		oldData := service_info.NewService_Info()
		if err := record.GetData(oldData); err != nil {
			fmt.Printf("record.GetData failed %s\n", err.Error())
			return
		}
		fmt.Printf("\ndeleteExample response record data %+v, route: %s",
			oldData, string(oldData.Routeinfo[0:oldData.Routeinfo_Len]))
		fmt.Printf("\ndeleteExample request  record data %+v, route: %s",
			data, string(data.Routeinfo[0:data.Routeinfo_Len]))
	}
}
