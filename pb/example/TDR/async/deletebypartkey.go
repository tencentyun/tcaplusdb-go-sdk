package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

func deleteByPartKeyExample() {
	//创建Delete请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiDeleteByPartkeyReq)
	if err != nil {
		fmt.Printf("NewRequest TcaplusApiDeleteByPartkeyReq failed %v\n", err.Error())
		return
	}
	fmt.Printf("deleteByPartKeyExample NewRequest TcaplusApiDeleteByPartkeyReq finish\n")
	//设置异步请求ID，异步请求通过ID让响应和请求对应起来
	req.SetAsyncId(770)
	//设置结果标记位，删除成功后，返回tcaplus端的旧数据，默认为0
	if err := req.SetVersionPolicy(3); err != nil {
		fmt.Printf("SetResultFlag failed %v\n", err.Error())
		return
	}
	fmt.Printf("deleteByPartKeyExample SetResultFlag finish\n")
	//为request添加一条记录，（index只有在list表中支持，generic不校验）
	rec, err := req.AddRecord(0)
	if err != nil {
		fmt.Printf("AddRecord failed %v\n", err.Error())
		return
	}
	req.SetResultLimit(20000, 0)
	fmt.Printf("deleteByPartKeyExample AddRecord finish\n")
	//申请tdr结构体并赋值key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "dev"
	//data.Envdata = "oaasqomk"
	data.Name = "com"
	//将tdr的数据设置到请求的记录中
	if err := rec.SetDataWithIndexAndField(data, nil, "Index_Gameid_Name"); err != nil {
		fmt.Printf("SetData failed %v\n", err.Error())
		return
	}
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	fmt.Printf("deleteByPartKeyExample send finish\n")
	var total_num int = 0
	//recv response
	for {
		resp, err := recvResponse(client)
		if err != nil {
			fmt.Printf("recv err %s\n", err.Error())
			return
		}
		//带回请求的异步ID
		fmt.Printf("deleteByPartKeyExample resp success, AsyncId:%d\n", resp.GetAsyncId())
		tcapluserr := resp.GetResult()
		if tcapluserr != 0 {
			fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
			return
		}
		has_more := resp.HaveMoreResPkgs()
		total_num += resp.GetRecordCount()
		fmt.Printf("deleteByPartKeyExample resp success, total_num:%d\n", total_num)
		//  response中带有获取的旧记录
		fmt.Printf("deleteByPartKeyExample response success record count %d\n", resp.GetRecordCount())
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
			fmt.Printf("gameid:%s, Envdata %s, name:%s, Expansion:%s\n",
				oldData.Gameid, oldData.Envdata, oldData.Name, oldData.Expansion)
			//fmt.Printf("\ndeleteByPartKeyExample response record data %+v, route: %s",
			// oldData, string(oldData.Routeinfo[0:oldData.Routeinfo_Len]))
			//fmt.Printf("\ndeleteByPartKeyExample request  record data %+v, route: %s",
			// data, string(data.Routeinfo[0:data.Routeinfo_Len]))
		}
		if 0 == has_more {
			break
		}
	}
}
