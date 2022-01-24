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
	// （非必须） 防止记录不存在
	client.DoListDeleteAll(msg, nil)
	client.DoListAddAfter(&tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}, -1, nil)
	client.DoListAddAfter(&tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}, -1, nil)

	opt := &option.PBOpt{}
	rspMsgs, err := client.DoListGetAll(msg, opt)
	if err != nil {
		logger.ERR("DoListGetAll error:%s", err)
		return
	}
	//记录version
	fmt.Println(opt.Version)
	for i, msg := range rspMsgs {
		fmt.Println(tools.ConvertToJson(msg))
		//list index
		fmt.Println(i)
	}
	fmt.Println("listgetall success")
}
