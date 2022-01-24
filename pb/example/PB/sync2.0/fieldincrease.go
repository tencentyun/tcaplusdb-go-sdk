package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
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

	msg := &tcaplusservice.GamePlayers{
		PlayerId:     10805514,
		PlayerName:   "Calvin",
		PlayerEmail:  "calvin@test.com",
		GameServerId: 15,
		Pay: &tcaplusservice.Payment{
			Amount: 10,
		},
	}

	opt := &option.PBOpt{
		FieldNames: []string{"game_server_id", "pay.amount"},
	}

	// 发送请求,接收响应
	err := client.DoFieldIncrease(msg, opt)
	if err != nil {
		logger.ERR("DoFieldIncrease error:%s", err)
		return
	}
	fmt.Println(tools.ConvertToJson(msg))
	fmt.Println(opt.Version) //记录版本
	fmt.Println("DoFieldIncrease success")
}
