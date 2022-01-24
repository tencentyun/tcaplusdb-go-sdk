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
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
		VersionPolicy:        option.CheckDataVersionAutoIncrease,
	}

	for _, msg := range msgs {
		msg.(*tcaplusservice.GamePlayers).Pay.PayId = 3
		opt.BatchVersion = append(opt.BatchVersion, 1)
	}

	err = client.DoBatchUpdate(msgs, opt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
	for i, msg := range msgs {
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) ||
			msg.(*tcaplusservice.GamePlayers).Pay.PayId != 2 {
			fmt.Println("error must old value")
			return
		}
	}

	fmt.Println("DoBatchUpdate success")
}
