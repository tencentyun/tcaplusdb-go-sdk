package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"google.golang.org/protobuf/proto"
	"time"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()
	client.SetDefaultZoneId(tools.ZoneId)

	//batch insert 10 data
	var msgs []proto.Message
	id := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		data.Pay = &tcaplusservice.Payment{Amount: uint64(i), PayId: 2, Method: 3}
		msgs = append(msgs, data)
	}

	err := client.DoBatchInsert(msgs, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//batch DoBatchUpdate
	opt := &option.PBOpt{
		VersionPolicy: option.CheckDataVersionAutoIncrease,
	}

	for i := 0; i < len(msgs); i++ {
		opt.BatchVersion = append(opt.BatchVersion, 1)
	}

	err = client.DoBatchDelete(msgs, opt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println("DoBatchDelete success")
}
