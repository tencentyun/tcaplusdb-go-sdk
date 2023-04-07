package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
)

func ListGetBatchExample() {
	var indexs []int32
	for i := 0; i < 10; i++ {
		indexs = append(indexs, int32(i))
	}
	msg := &tcaplusservice.TbOnlineList{
		Openid:   1,
		Tconndid: 2,
		Timekey:  "key",
	}
	opt := &option.PBOpt{
		MultiFlag: 1,
	}
	resMsgs, err := client.DoListGetBatch(msg, indexs, opt)
	if err != nil {
		fmt.Printf("DoListGetBatch fail, %s", err.Error())
		return
	}

	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	for index, msg := range resMsgs {
		fmt.Println("index", index)
		fmt.Println(msg)
	}
	fmt.Println("DoListGetBatch success")
}
