package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"google.golang.org/protobuf/proto"
)

func PbBatchFieldGetExample() {
	msg1 := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}

	msg2 := &tcaplusservice.GamePlayers{
		PlayerId:    10805515,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}

	msg3 := &tcaplusservice.GamePlayers{
		PlayerId:    108055150000,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}
	var msgs []proto.Message
	msgs = append(msgs, msg1)
	msgs = append(msgs, msg2)
	msgs = append(msgs, msg3)

	// 发送请求,接收响应
	opt := &option.PBOpt{
		MultiFlag:  1,
		FieldNames: []string{"game_server_id", "pay.amount"},
	}
	err := client.DoBatchFieldGet(msgs, opt)
	if err != nil {
		logger.ERR("DoBatchGet error:%s", err)
	}

	for i, rspMsg := range msgs {
		fmt.Println(tools.ConvertToJson(rspMsg))
		fmt.Println(rspMsg.(*tcaplusservice.GamePlayers))
		//单条记录的错误码
		fmt.Println("result", opt.BatchResult[i])
		//记录version
		fmt.Println("version", opt.BatchVersion[i])
	}
	logger.INFO("batch field get success")
	fmt.Println("batch field get success")
}
