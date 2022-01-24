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

	// 向记录中填充数据
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
