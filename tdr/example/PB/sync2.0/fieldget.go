package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
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

	// 向记录中填充数据
	msg := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}

	// 设置获取字段 game_server_id 和 二级字段 pay.amount
	opt := &option.PBOpt{
		FieldNames: []string{"game_server_id", "pay.amount"},
	}
	// 发送请求,接收响应
	err := client.DoFieldGet(msg, opt)
	if err != nil {
		logger.ERR("DoFieldGet error:%s", err)
		return
	}
	fmt.Println(tools.ConvertToJson(msg))
	fmt.Println("field get success")
}
