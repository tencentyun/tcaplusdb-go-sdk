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
	// 清除key
	client.DoListDeleteAll(msg, nil)
	client.DoListAddAfter(&tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}, -1, nil)

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
