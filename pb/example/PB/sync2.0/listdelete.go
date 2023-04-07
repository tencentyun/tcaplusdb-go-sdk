package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
)

func ListDeleteExample() {
	// 向记录中填充数据
	msg := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}

	opt := &option.PBOpt{
		ResultFlag: option.TcaplusResultFlagAllOldValue,
	}
	err := client.DoListDelete(msg, 0, opt)
	if err != nil {
		logger.ERR("DoListDelete error:%s", err)
		return
	}

	fmt.Println(tools.ConvertToJson(msg))
	fmt.Println(opt.Version)
	fmt.Println("listdelete success")
}
