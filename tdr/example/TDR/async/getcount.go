package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"strconv"
	"sync"
	"time"
)

func getCountExample() {
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
			fmt.Printf("response success record count %d\n", resp.GetTableRecordCount())
			return
		}
	}()

	//创建Get请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiGetTableRecordCountReq)
	if err != nil {
		fmt.Printf("NewRequest failed %v\n", err.Error())
		return
	}

	//设置异步请求ID，异步请求通过ID让响应和请求对应起来
	req.SetAsyncId(667)

	// 异步发送请求
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	wg.Wait()
}
