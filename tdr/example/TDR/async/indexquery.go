package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"strconv"
	"sync"
	"time"
)

func sqlExample() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 在另一协程处理响应消息
	go func() {
		defer wg.Done()
		for {
			// resp err 均为 nil 说明响应池中没有任何响应
			resp, err := client.RecvResponse()
			if err != nil {
				logger.ERR("RecvResponse error:%s", err)
				continue
			} else if resp == nil {
				time.Sleep(time.Microsecond * 5)
				continue
			}

			//带回请求的异步ID
			fmt.Printf("resp success, AsyncId:%d\n", resp.GetAsyncId())
			tcapluserr := resp.GetResult()
			if tcapluserr != 0 {
				fmt.Printf("response ret %s\n",
					"errCode: "+strconv.Itoa(tcapluserr)+", errMsg: "+terror.ErrorCodes[tcapluserr])
				return
			}
			//response中带有获取的记录
			fmt.Printf("response success record count %d\n", resp.GetRecordCount())
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
				fmt.Printf("response record data %+v, route: %s\n",
					data, string(data.Routeinfo[0:data.Routeinfo_Len]))
			}
			//判断是否有分包
			if 1 == resp.HaveMoreResPkgs() {
				continue
			}
			return
		}
	}()

	//创建请求，设置sql语句
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiSqlReq)
	if err != nil {
		fmt.Printf("NewRequest TcaplusApiSqlReq failed %v\n", err.Error())
		return
	}

	query := fmt.Sprintf("select * from service_info where Inst_Max_Num > 100 limit 10 offset 0; ")
	req.SetSql(query)

	//发送请求
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	wg.Wait()
}

//聚合查询
func sqlExample2() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 在另一协程处理响应消息
	go func() {
		defer wg.Done()
		for {
			// resp err 均为 nil 说明响应池中没有任何响应
			resp, err := client.RecvResponse()
			if err != nil {
				logger.ERR("RecvResponse error:%s", err)
				continue
			} else if resp == nil {
				time.Sleep(time.Microsecond * 5)
				continue
			}

			//带回请求的异步ID
			fmt.Printf("resp success, AsyncId:%d\n", resp.GetAsyncId())
			tcapluserr := resp.GetResult()
			if tcapluserr != 0 {
				fmt.Printf("response ret %s\n",
					"errCode: "+strconv.Itoa(tcapluserr)+", errMsg: "+terror.ErrorCodes[tcapluserr])
				return
			}
			//获取聚合查询结果
			r, err := resp.ProcAggregationSqlQueryType()
			if err != nil {
				fmt.Printf("ProcAggregationSqlQueryType failed %s", err.Error())
				return
			}
			fmt.Printf("%s", common.CovertToJson(r))
			return
		}
	}()
	//创建请求，设置sql语句
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiSqlReq)
	if err != nil {
		fmt.Printf("NewRequest TcaplusApiSqlReq failed %v\n", err.Error())
		return
	}

	query := fmt.Sprintf("select count(Inst_Max_Num), min(Inst_Max_Num), from service_info where Inst_Max_Num > 100;")
	req.SetSql(query)
	//发送请求
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	wg.Wait()
}
