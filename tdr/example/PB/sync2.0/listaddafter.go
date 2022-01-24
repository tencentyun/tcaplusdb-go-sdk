package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()
	client.SetDefaultZoneId(tools.ZoneId)

	// 向记录中填充数据
	msg := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}
	client.DoListDeleteAll(msg, nil)

	opt := &option.PBOpt{
		ResultFlag: option.TcaplusResultFlagAllNewValue,
	}
	_, err := client.DoListAddAfter(msg, -1, opt)
	if err != nil {
		logger.ERR("DoListAddAfter error:%s", err)
		return
	}
	fmt.Println(tools.ConvertToJson(msg))
	fmt.Println(opt.Version)
	fmt.Println("listaddafter success")
}
