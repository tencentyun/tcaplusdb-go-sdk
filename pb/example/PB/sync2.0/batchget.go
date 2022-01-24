package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"google.golang.org/protobuf/proto"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()
	client.SetDefaultZoneId(tools.ZoneId)

	// （非必须） 防止记录不存在
	client.DoInsert(&tcaplusservice.GamePlayers{
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
	}, nil)
	client.DoInsert(&tcaplusservice.GamePlayers{
		PlayerId:        10805514,
		PlayerName:      "Calvin",
		PlayerEmail:     "zhang@test.com",
		GameServerId:    10,
		LoginTimestamp:  []string{"2019-12-12 15:00:00"},
		LogoutTimestamp: []string{"2019-12-12 16:00:00"},
		IsOnline:        false,
		Pay: &tcaplusservice.Payment{
			PayId:  10101,
			Amount: 1000,
			Method: 2,
		},
	}, nil)

	msgs := []proto.Message{
		&tcaplusservice.GamePlayers{
			PlayerId:    10805514,
			PlayerName:  "Calvin",
			PlayerEmail: "zhang@test.com",
		}, &tcaplusservice.GamePlayers{
			PlayerId:    10805514,
			PlayerName:  "Calvin",
			PlayerEmail: "calvin@test.com",
		}}

	// 发送请求,接收响应
	opt := &option.PBOpt{}
	err := client.DoBatchGet(msgs, opt)
	if err != nil {
		logger.ERR("DoBatchGet error:%s", err)
		return
	}

	for i, rspMsg := range msgs {
		fmt.Println(tools.ConvertToJson(rspMsg))
		//单条记录的错误码
		fmt.Println(opt.BatchResult[i])
		//记录version
		fmt.Println(opt.BatchVersion[i])
	}

	logger.INFO("batch get success")
	fmt.Println("batch get success")
}
