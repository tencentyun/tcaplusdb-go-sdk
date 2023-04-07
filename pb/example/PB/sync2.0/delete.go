package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
)

func DeleteExample() {
	msg := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}

	//可选参数设置
	opt := &option.PBOpt{
		ResultFlag: option.TcaplusResultFlagAllOldValue,
	}
	err := client.DoDelete(msg, opt)
	if err != nil {
		logger.ERR("DoDelete error:%s", err)
		return
	}

	//设置了resultflag，svr返回的msg会覆盖原msg
	fmt.Println(opt.Version)
	fmt.Println(tools.ConvertToJson(msg))
	fmt.Println("delete success")
}