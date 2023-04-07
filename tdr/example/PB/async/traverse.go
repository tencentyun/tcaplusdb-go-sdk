package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/traverser"
	"time"
)

// 针对同一个表只能存在一个遍历器
// 同步请求的遍历只针对数据量较小的表，由于会把所有的数据存放在内存中，占用内存
// 数据量较大的表可以使用异步遍历，边遍历边处理数据
func TraverseExample() {
	// 创建异步协程接收请求
	respChan := make(chan response.TcaplusResponse)
	go func() {
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
			// 同步异步 id 找到对应的响应
			if resp.GetAsyncId() == 12345 {
				respChan <- resp
			}
		}
	}()

	// 获取遍历器，遍历器最多同时8个工作，如果超过会返回nil
	tra := client.GetTraverser(tools.Zone, "game_players")
	if tra == nil {
		logger.ERR("GetTraverser fail")
		return
	}
	// 调用stop才能释放资源，防止获取遍历器失败
	defer tra.Stop()

	// （非必须）设置 异步 id
	tra.SetAsyncId(12345)

	// （非必须）限制本次遍历记录条数，默认不限制
	tra.SetLimit(1000)

	// 发送请求
	err := tra.Start()
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}

	for {
		timeOutChan := time.After(5 * time.Second)
		select {
		case <-timeOutChan:
			if tra.State() == traverser.TraverseStateNormal {
				fmt.Println("continue ......")
			} else if tra.State() == traverser.TraverseStateIdle {
				fmt.Println("traverse finish")
				return
			} else {
				fmt.Println("traverse stat err ", tra.State())
				return
			}
		// 等待收取响应
		case resp := <-respChan:
			// 获取响应结果
			errCode := resp.GetResult()
			if errCode != terror.GEN_ERR_SUC {
				logger.ERR("insert error:%s", terror.GetErrMsg(errCode))
				return
			}

			// 如果有返回记录则用以下接口进行获取
			for i := 0; i < resp.GetRecordCount(); i++ {
				record, err := resp.FetchRecord()
				if err != nil {
					logger.ERR("FetchRecord failed %s", err.Error())
					return
				}

				newMsg := &tcaplusservice.GamePlayers{}
				err = record.GetPBData(newMsg)
				if err != nil {
					logger.ERR("GetPBData failed %s", err.Error())
					return
				}

				fmt.Println(tools.ConvertToJson(newMsg))
			}
		}
	}
}
