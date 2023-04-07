package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
)

func GetExample(){
	msg := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}
	opt := &option.PBOpt{}
	err := client.DoGet(msg, opt)
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}
	fmt.Println(tools.ConvertToJson(msg))
	fmt.Println(opt.Version) //记录版本
	fmt.Println("DoFieldUpdate success")
}