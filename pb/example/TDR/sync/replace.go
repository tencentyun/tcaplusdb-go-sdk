package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"time"
)

func replaceExample() {
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiReplaceReq)
	if err != nil {
		fmt.Printf("NewRequest TcaplusApiReplaceReq failed %v\n", err.Error())
		return
	}
	fmt.Printf("replaceExample NewRequest TcaplusApiReplaceReq finish\n")
	if err := req.SetResultFlag(3); err != nil {
		fmt.Printf("SetResultFlag failed %v\n", err.Error())
		return
	}
	fmt.Printf("replaceExample SetResultFlag finish\n")

	//为request添加一条记录，（index只有在list表中支持，generic不校验）
	rec, err := req.AddRecord(0)
	if err != nil {
		fmt.Printf("AddRecord failed %v\n", err.Error())
		return
	}
	fmt.Printf("replaceExample AddRecord finish\n")
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
	route := "test"
	data.Routeinfo_Len = uint32(len(route))
	data.Routeinfo = []byte(route)
	//将tdr的数据设置到请求的记录中
	if err := rec.SetData(data); err != nil {
		fmt.Printf("SetData failed %v\n", err.Error())
		return
	}
	if resp, err := client.Do(req, time.Duration(2*time.Second)); err != nil {
		fmt.Printf("recv err %s\n", err.Error())
		return
	} else {
		tcapluserr := resp.GetResult()
		if tcapluserr != 0 {
			fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
			return
		}
		//获取同步请求Seq
		fmt.Printf("request Seq %d\n", req.GetSeq())
		//获取回应消息的序列号
		fmt.Printf("respond seq: %d \n", resp.GetSeq())
		//response中带有获取的旧记录
		fmt.Printf("replaceExample response success record count %d\n", resp.GetRecordCount())
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
			fmt.Printf("replaceExample response record data %+v, route: %s\n",
				oldData, string(oldData.Routeinfo[0:oldData.Routeinfo_Len]))
			fmt.Printf("replaceExample request  record data %+v, route: %s\n",
				data, string(data.Routeinfo[0:data.Routeinfo_Len]))
		}
	}
}
