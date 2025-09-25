package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
)

func deleteByPartKeyExample() {
	// 向记录中填充部分key
	msg := &tcaplusservice.GamePlayers{
		PlayerId:   10805514,
		PlayerName: "Calvin",
	}
	// 本例中使用的是本地索引 option(tcaplusservice.tcaplus_index) = "index_1(player_id, player_name)";
	rspMsgs, err := client.DoDeleteByPartKey(msg, []string{"player_id", "player_name"}, nil)
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}
	for _, msg := range rspMsgs {
		fmt.Println(tools.ConvertToJson(msg))
	}
	fmt.Println("delete by part key success, size ", len(rspMsgs))
}
