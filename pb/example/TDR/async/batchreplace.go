package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"strconv"
	"sync"
	"time"
)

func batchReplaceExample() {
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
	//创建请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiBatchReplaceReq)
	if err != nil {
		fmt.Printf(" NewRequest failed %v\n", err.Error())
		return
	}
	//允许分包
	req.SetMultiResponseFlag(1)
	// 请求标志。0标志只需返回成功与否,1标志返回同请求一致的值,2标志返回操作后所有字段的值,3标志返回操作前所有字段的值
	req.SetResultFlagForSuccess(2)

	for i := 0; i < 10; i++ {
		//为request添加一条记录,最多1024
		rec, err := req.AddRecord(0)
		if err != nil {
			fmt.Printf(" AddRecord failed %v\n", err.Error())
			return
		}
		//申请tdr结构体并赋值Key，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
		// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
		data := service_info.NewService_Info()
		data.Gameid = "dev"
		data.Envdata = fmt.Sprintf("%d", i)
		data.Name = "com"
		data.Filterdata = time.Now().Format("2006-01-02T15:04:05.000000Z")
		data.Updatetime = uint64(time.Now().UnixNano())
		data.Inst_Max_Num = 2
		data.Inst_Min_Num = 3
		//数组类型为slice需要准确赋值长度，与refer保持一致
		route := "test"
		data.Routeinfo_Len = uint32(len(route))
		data.Routeinfo = []byte(route)
		if err := rec.SetData(data); err != nil {
			fmt.Printf("SetData failed %v\n", err.Error())
			return
		}
	}

	//发送请求
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	wg.Wait()
}
