package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"strconv"
	"time"
)

func batchGetExample() {
	//创建请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiBatchGetReq)
	if err != nil {
		fmt.Printf("NewRequest  failed %v\n", err.Error())
		return
	}
	//允许分包
	req.SetMultiResponseFlag(1)

	for i := 0; i < 10; i++ {
		//为request添加一条记录,最多1024
		rec, err := req.AddRecord(0)
		if err != nil {
			fmt.Printf("AddRecord failed %v\n", err.Error())
			return
		}
		//申请tdr结构体并赋值Key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
		// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
		data := service_info.NewService_Info()
		data.Gameid = "dev"
		data.Envdata = fmt.Sprintf("%d", i)
		data.Name = "com"
		if err := rec.SetData(data); err != nil {
			fmt.Printf("SetData failed %v\n", err.Error())
			return
		}
	}

	respList, err := client.DoMore(req, time.Duration(10*time.Second))
	if err != nil {
		fmt.Printf("recv err %s\n", err.Error())

	}
	var totalCnt int = 0
	for _, resp := range respList {
		tcapluserr := resp.GetResult()
		if tcapluserr != 0 {
			fmt.Printf("response ret %s\n",
				"errCode: "+strconv.Itoa(tcapluserr)+", errMsg: "+terror.ErrorCodes[tcapluserr])
			break
		}
		totalCnt += resp.GetRecordCount()
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
		fmt.Printf("response record data %+v, route: %s\n",
			data, string(data.Routeinfo[0:data.Routeinfo_Len]))
	}
	fmt.Printf("total count %d,\n", totalCnt)
}
