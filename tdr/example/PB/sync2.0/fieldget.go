package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
)

func FieldGetExample() {
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
