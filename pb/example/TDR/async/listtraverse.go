package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/response"
	"time"

	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/TDR/async/service_info"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
)

//数据量大可以使用异步的遍历，可以在遍历的过程中处理数据
//数据量小可以使用同步遍历，否则记录堆积会占较多内存
func ListTraverseExample() {
	// 创建异步协程接收遍历响应
	respChan := make(chan response.TcaplusResponse, 1024)
	go func() {
		for {
			// resp err 均为 nil 说明响应池中没有任何响应
			resp, err := client.RecvResponse()
			if err != nil {
				fmt.Printf("RecvResponse error:%s\n", err)
				continue
			} else if resp == nil {
				time.Sleep(time.Microsecond * 5)
				continue
			}
			// 同步异步 id 找到对应的响应
			if resp.GetAsyncId() == 12345 {
				respChan <- resp
			}
		}
	}()

	tra := client.GetListTraverser(ZoneId, TABLE_TRAVERSER_LIST)
	defer tra.Stop()

	tra.SetFieldNames([]string{"filterdata", "updatetime"})
	// （非必须）限制本次遍历记录条数，默认不限制
	tra.SetLimit(10)
	//设置 异步 id
	tra.SetAsyncId(12345)

	err := tra.Start()
	if err != nil {
		fmt.Printf("SendRequest error:%s\n", err)
		return
	}

	timeOutChan := time.NewTimer(30 * time.Second)
	for {
		timeOutChan.Reset(30 * time.Second)
		select {
		case <-timeOutChan.C:
			// 防止一直卡主
			fmt.Println("30 秒没有响应包了 stat ", tra.State())
			return

		// 等待收取响应
		case resp := <-respChan:
			// 获取响应结果
			errCode := resp.GetResult()
			if errCode != terror.GEN_ERR_SUC {
				fmt.Printf("insert error:%s\n", terror.GetErrMsg(errCode))
				return
			}

			// 如果有返回记录则用以下接口进行获取
			for i := 0; i < resp.GetRecordCount(); i++ {
				record, err := resp.FetchRecord()
				if err != nil {
					fmt.Printf("FetchRecord failed %s\n", err.Error())
					return
				}

				data := service_info.NewService_Info()
				if err := record.GetData(data); err != nil {
					fmt.Printf("record.GetData failed %s\n", err.Error())
					return
				}
				fmt.Println("record service_info : ", data)
			}
			if resp.HaveMoreResPkgs() == 0 {
				//finished
				fmt.Println("traverse finish")
				return
			}
		}
	}
}
