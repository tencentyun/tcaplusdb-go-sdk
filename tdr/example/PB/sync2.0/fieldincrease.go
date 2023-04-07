package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
)

func FieldIncrease() {
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
