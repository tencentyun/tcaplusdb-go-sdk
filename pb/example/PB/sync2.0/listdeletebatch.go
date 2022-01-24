package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()
	client.SetDefaultZoneId(tools.ZoneId)

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
