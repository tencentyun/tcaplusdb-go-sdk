package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"google.golang.org/protobuf/proto"
)

func BatchGetExample() {
	msgs := []proto.Message{
		&tcaplusservice.GamePlayers{
			PlayerId:    10805514,
			PlayerName:  "Calvin",
			PlayerEmail: "zhang@test.com",
		}, &tcaplusservice.GamePlayers{
			PlayerId:    10805514,
			PlayerName:  "Calvin",
			PlayerEmail: "calvin@test.com",
		}}

	// 发送请求,接收响应
	opt := &option.PBOpt{
		MultiFlag: 1,
	}
	err := client.DoBatchGet(msgs, opt)
	if err != nil {
		logger.ERR("DoBatchGet error:%s", err)
		return
	}

	for i, rspMsg := range msgs {
		fmt.Println(tools.ConvertToJson(rspMsg))
		//单条记录的错误码
		fmt.Println(opt.BatchResult[i])
		//记录version
		fmt.Println(opt.BatchVersion[i])
	}

	logger.INFO("batch get success")
	fmt.Println("batch get success")
}
