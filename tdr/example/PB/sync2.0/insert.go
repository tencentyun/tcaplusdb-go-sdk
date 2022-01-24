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

	// 向记录中填充数据
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
			Method: 1,
		},
	}
	// （非必须）防止此条记录已存在
	client.DoDelete(msg, nil)

	opt := &option.PBOpt{
		ResultFlag: option.TcaplusResultFlagAllNewValue,
	}
	err := client.DoInsert(msg, opt)
	if err != nil {
		logger.ERR("DoInsert error:%s", err)
		return
	}

	fmt.Println(tools.ConvertToJson(msg))
	fmt.Println(opt.Version)
	fmt.Println("insert success")
}
