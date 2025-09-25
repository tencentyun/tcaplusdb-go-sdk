package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
)

func ListGetAllExample() {
	// 向记录中填充数据
	msg := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}

	opt := &option.PBOpt{
		MultiFlag: 1,
	}
	idx, rspMsgs, err := client.DoListGetAllV2(msg, opt)
	if err != nil {
		logger.ERR("DoListGetAll error:%s", err)
		return
	}
	//记录version
	fmt.Println(opt.Version)
	for i, msg := range rspMsgs {
		fmt.Println(tools.ConvertToJson(msg))
		//list index
		fmt.Println(idx[i])
	}
	fmt.Println("listgetall success")
}
