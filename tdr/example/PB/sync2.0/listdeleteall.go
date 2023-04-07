package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
)

func ListDeleteAllExample() {
	msg := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}

	err := client.DoListDeleteAll(msg, nil)
	if err != nil {
		logger.ERR("DoListDeleteAll error:%s", err)
		return
	}

	fmt.Println("listdeleteall success")
}
