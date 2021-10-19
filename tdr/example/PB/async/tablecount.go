package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()

	// 获取下当前记录总数
	record := &tcaplusservice.GamePlayers{
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
	}
	client.Delete(record)
	old, _ := client.GetTableCount("game_players")

	// 创建异步协程接收请求
	respChan := make(chan response.TcaplusResponse)
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

	// 生成 get 请求
	req, err := client.NewRequest(tools.ZoneId, "game_players", cmd.TcaplusApiGetTableRecordCountReq)
	if err != nil {
		fmt.Printf("NewRequest error:%s\n", err)
		return
	}

	// （非必须）设置 异步 id
	req.SetAsyncId(12345)

	// （非必须）设置userbuf，在响应中带回。这个是个开放功能，比如某些临时字段不想保存在全局变量中，
	// 可以通过设置userbuf在发送端接收短传递，也可以起异步id的作用
	req.SetUserBuff([]byte("user buffer test"))

	// （非必须） 防止记录不存在
	client.Insert(record)
	defer client.Delete(record)

	// 发送请求
	err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("SendRequest error:%s\n", err)
		return
	}

	// 等待收取响应
	resp := <-respChan

	// 获取响应结果
	errCode := resp.GetResult()
	if errCode != terror.GEN_ERR_SUC {
		fmt.Printf("insert error:%s\n", terror.GetErrMsg(errCode))
		return
	}

	// 获取userbuf
	fmt.Println(string(resp.GetUserBuffer()))

	if resp.GetTableRecordCount() != old+1 {
		fmt.Printf("resp.GetTableRecordCount() %d != %d\n", resp.GetTableRecordCount(), old+1)
		return
	}

	logger.INFO("count success")
	fmt.Println("count success")
}
