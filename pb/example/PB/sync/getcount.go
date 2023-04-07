package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"time"
)

func getCountExample() {
	// 生成 get 请求
	req, err := client.NewRequest(tools.Zone, "game_players", cmd.TcaplusApiGetTableRecordCountReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return
	}

	// 发送请求,接收响应
	resp, err := client.Do(req, 5*time.Second)
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}

	// 获取响应结果
	errCode := resp.GetResult()
	if errCode != terror.GEN_ERR_SUC {
		logger.ERR("insert error:%s", terror.GetErrMsg(errCode))
		return
	}
	fmt.Println("table rec count", resp.GetTableRecordCount())
	logger.INFO("get success")
	fmt.Println("get success")
}
