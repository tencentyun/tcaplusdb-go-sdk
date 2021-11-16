package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/traverser"
	"time"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()

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
	tra := client.GetTraverser(tools.ZoneId, "game_players")
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

	// （非必须）设置userbuf，在响应中带回。这个是个开放功能，比如某些临时字段不想保存在全局变量中，
	// 可以通过设置userbuf在发送端接收短传递，也可以起异步id的作用
	tra.SetUserBuff([]byte("user buffer test"))

	// （非必须） 防止记录不存在
	client.Insert(&tcaplusservice.GamePlayers{
		PlayerId:        10805514,
		PlayerName:      "Calvin",
		PlayerEmail:     "calvin@test.com",
		GameServerId:    10,
		LoginTimestamp:  []string{"2019-12-12 15:00:00"},
		LogoutTimestamp: []string{"2019-12-12 16:00:00"},
		IsOnline:        false,
		Pay: &tcaplusservice.Payment{
			PayId:  10101,
			Amount: 1000,
			Method: 2,
		},
	})

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
			if tra.State() != traverser.TraverseStateNormal {
				logger.INFO("traverse finish")
				fmt.Println("traverse finish")
				return
			}
			fmt.Println("......")
		// 等待收取响应
		case resp := <-respChan:
			// 获取响应结果
			errCode := resp.GetResult()
			if errCode != terror.GEN_ERR_SUC {
				logger.ERR("insert error:%s", terror.GetErrMsg(errCode))
				return
			}

			// 获取userbuf
			fmt.Println(string(resp.GetUserBuffer()))

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
