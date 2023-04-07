package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"time"
)

func getExample() {
	// 创建请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiGetReq)
	if err != nil {
		fmt.Printf("getExample NewRequest TcaplusApiGetReq failed %v\n", err.Error())
		return
	}

	//为request添加一条记录，（index只有在list表中支持，generic不校验）
	rec, err := req.AddRecord(0)
	if err != nil {
		fmt.Printf("getExample AddRecord failed %v\n", err.Error())
		return
	}

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

	// 使用客户端同步发送请求并接收响应
	if resp, err := client.Do(req, time.Duration(2*time.Second)); err != nil {
		fmt.Printf("recv err %s\n", err.Error())
		return
	} else {
		// 获取响应消息的错误码
		tcapluserr := resp.GetResult()
		if tcapluserr != 0 {
			fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
			return
		}

		// 从响应消息中提取记录
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
}
