package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"google.golang.org/protobuf/proto"
)

func ListAddafterBatchExample() {
	//1 addafterbatch 10 条记录
	var msgs []proto.Message
	var indexs []int32
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   "key",
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i),
			},
		}
		fmt.Println(data)
		msgs = append(msgs, data)
		indexs = append(indexs, -1)
	}
	opt := &option.PBOpt{}
	err := client.DoListAddAfterBatch(msgs, indexs, opt)
	if err != nil {
		fmt.Printf("DoListAddAfterBatch fail, %s", err.Error())
		return
	}
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println(indexs)
	fmt.Println("DoListAddAfterBatch success")
}
