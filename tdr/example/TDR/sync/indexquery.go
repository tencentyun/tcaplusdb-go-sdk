package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

func sqlExample() {
	//创建请求，设置sql语句
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiSqlReq)
	if err != nil {
		fmt.Printf("NewRequest TcaplusApiSqlReq failed %v\n", err.Error())
		return
	}

	query := fmt.Sprintf("select * from service_info where Inst_Max_Num > 100 limit 10 offset 0; ")
	req.SetSql(query)

	//发送请求
	if resps, err := client.DoMore(req, time.Duration(2*time.Second)); err != nil {
		fmt.Printf("recv err %s\n", err.Error())
		return
	} else {
		for _, resp := range resps {
			tcapluserr := resp.GetResult()
			if tcapluserr != 0 {
				fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
				return
			}
			//response中带有获取的记录
			fmt.Printf("replaceExample response success record count %d\n", resp.GetRecordCount())
			for i := 0; i < resp.GetRecordCount(); i++ {
				record, err := resp.FetchRecord()
				if err != nil {
					fmt.Printf("FetchRecord failed %s\n", err.Error())
					return
				}
				resData := service_info.NewService_Info()
				if err := record.GetData(resData); err != nil {
					fmt.Printf("record.GetData failed %s\n", err.Error())
					return
				}
				fmt.Printf("sqlExample response record data %+v, route: %s\n",
					resData, string(resData.Routeinfo[0:resData.Routeinfo_Len]))
			}
		}
	}
}

//聚合查询
func sqlExample2() {
	//创建请求，设置sql语句
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiSqlReq)
	if err != nil {
		fmt.Printf("NewRequest TcaplusApiSqlReq failed %v\n", err.Error())
		return
	}

	query := fmt.Sprintf("select count(Inst_Max_Num), min(Inst_Max_Num), from service_info where Inst_Max_Num > 100;")
	req.SetSql(query)

	//发送请求
	if resp, err := client.Do(req, time.Duration(2*time.Second)); err != nil {
		fmt.Printf("recv err %s\n", err.Error())
		return
	} else {
		tcapluserr := resp.GetResult()
		if tcapluserr != 0 {
			fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
			return
		}
		//获取聚合查询结果
		r, err := resp.ProcAggregationSqlQueryType()
		if err != nil {
			fmt.Printf("ProcAggregationSqlQueryType failed %s", err.Error())
			return
		}
		fmt.Printf("%s", common.CovertToJson(r))
	}
}
