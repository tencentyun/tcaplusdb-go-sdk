package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"strconv"
	"time"
)

func getPartKeyExample() {
	//创建Get请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiGetByPartkeyReq)
	if err != nil {
		fmt.Printf("getPartKeyExample NewRequest TcaplusApiGetReq failed %v\n", err.Error())
		return
	}
	fmt.Printf("getPartKeyExample NewRequest TcaplusApiGetReq finish\n")

	//为request添加一条记录
	rec, err := req.AddRecord(0)
	if err != nil {
		fmt.Printf("getPartKeyExample AddRecord failed %v\n", err.Error())
		return
	}
	req.SetResultLimit(2000, 5000)
	fmt.Printf("getPartKeyExample AddRecord finish\n")
	//申请tdr结构体并赋值Key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := service_info.NewService_Info()
	data.Gameid = "dev"
	//data.Envdata = "oaasqomk"
	data.Name = "com"
	//将tdr的数据设置到请求的记录中
	//flist := []string {"updatetime"}
	var flist []string = nil
	if err := rec.SetDataWithIndexAndField(data, flist, "Index_Gameid_Name"); err != nil {
		fmt.Printf("SetData failed %v\n", err.Error())
		return
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

		//response中带有获取的记录
		fmt.Printf("getPartKeyExample response success record count %d, have more :%d\n",
			resp.GetRecordCount(), resp.HaveMoreResPkgs())
	}
	fmt.Printf("getPartKeyExample total count %d,\n", totalCnt)

}
