package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
)

func FieldUpdate() {
	msg := &tcaplusservice.GamePlayers{
		PlayerId:     10805514,
		PlayerName:   "Calvin",
		PlayerEmail:  "calvin@test.com",
		GameServerId: 15,
		Pay: &tcaplusservice.Payment{
			Amount: 1000,
		},
	}

	opt := &option.PBOpt{
		// （非必须）设置记录版本的检查类型，用于乐观锁
		//VersionPolicy: option.CheckDataVersionAutoIncrease,
		//Version: 100,
		FieldNames: []string{"game_server_id", "pay.amount"},
	}
	err := client.DoFieldUpdate(msg, opt)
	if err != nil {
		logger.ERR("DoFieldUpdate error:%s", err)
		return
	}
	fmt.Println(tools.ConvertToJson(msg))
	fmt.Println(opt.Version) //记录版本
	fmt.Println("DoFieldUpdate success")
}
