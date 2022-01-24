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
		PlayerId:     10805514,
		PlayerName:   "Calvin",
		PlayerEmail:  "calvin@test.com",
		GameServerId: 10,
	}, nil)

	msg := &tcaplusservice.GamePlayers{
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
			Method: 3,
		},
	}

	opt := &option.PBOpt{
		ResultFlag: option.TcaplusResultFlagAllNewValue,
	}
	err := client.DoUpdate(msg, opt)
	if err != nil {
		logger.ERR("DoUpdate error:%s", err)
		return
	}

	fmt.Println(tools.ConvertToJson(msg))
	fmt.Println(opt.Version)
	fmt.Println("DoUpdate success")
}
