package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"google.golang.org/protobuf/proto"
)

func BatchDeleteExample() {
	//batch delete 10 data
	var msgs []proto.Message
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = int64(i)
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		data.Pay = &tcaplusservice.Payment{Amount: uint64(i), PayId: 2, Method: 3}
		msgs = append(msgs, data)
	}

	//batch DoBatchUpdate
	opt := &option.PBOpt{
		ResultFlagForSuccess: 3,
	}

	err := client.DoBatchDelete(msgs, opt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println("DoBatchDelete success")
}
