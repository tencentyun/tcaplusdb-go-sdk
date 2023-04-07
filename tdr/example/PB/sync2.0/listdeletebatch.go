package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
)

func ListDeleteBatchExample() {
	msg := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}

	opt := &option.PBOpt{
		ResultFlag: option.TcaplusResultFlagAllOldValue,
	}
	rspMsgs, err := client.DoListDeleteBatch(msg, []int32{0, 1}, nil)
	if err != nil {
		logger.ERR("DoListDeleteBatch error:%s", err)
		return
	}
	//记录version
	fmt.Println(opt.Version)
	for i, msg := range rspMsgs {
		fmt.Println(tools.ConvertToJson(msg))
		//list index
		fmt.Println(i)
	}
	fmt.Println("listdeletebatch success")
}
