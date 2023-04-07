package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
)

func GetByPartKeyExample() {
	// 向记录中填充部分key
	msg := &tcaplusservice.GamePlayers{
		PlayerId:   10805514,
		PlayerName: "Calvin",
	}
	opt := &option.PBOpt{}
	// 本例中使用的是本地索引 option(tcaplusservice.tcaplus_index) = "index_1(player_id, player_name)";
	rspMsgs, err := client.DoGetByPartKey(msg, []string{"player_id", "player_name"}, opt)
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}
	for i, msg := range rspMsgs {
		fmt.Println(tools.ConvertToJson(msg))
		//记录version
		fmt.Println(opt.BatchVersion[i])
	}
	fmt.Println("get by part key success")
}
