package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"google.golang.org/protobuf/proto"
)

func SetTTLExample() {
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
	// set ttl
	opt := &option.PBOpt{
		BatchTTL: []option.TTLInfo{{TTL: 5000}},
	}
	msgs := []proto.Message{msg}
	err := client.DoSetTTLBatch(msgs, nil, opt)
	if err != nil {
		fmt.Printf("DoSetTTLBatch fail, %s", err.Error())
		return
	}
	fmt.Println("DoSetTTLBatch success")
}

func GetTTLExample() {
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
	msgs := []proto.Message{msg}
	// get ttl
	opt := &option.PBOpt{}
	err := client.DoGetTTLBatch(msgs, nil, opt)
	if err != nil {
		fmt.Printf("DoGetTTLBatch fail, %s", err.Error())
		return
	}
	fmt.Println(opt.BatchTTL)
	fmt.Println("DoGetTTLBatch success")
}
