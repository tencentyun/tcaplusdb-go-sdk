package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

func getCountExample() {
	// 创建请求
	req, err := client.NewRequest(ZoneId, TableName, cmd.TcaplusApiGetTableRecordCountReq)
	if err != nil {
		fmt.Printf("getExample NewRequest failed %v\n", err.Error())
		return
	}

	// 使用客户端同步发送请求并接收响应
	if resp, err := client.Do(req, time.Duration(2*time.Second)); err != nil {
		fmt.Printf("recv err %s\n", err.Error())
		return
	} else {
		// 获取响应消息的错误码
		tcapluserr := resp.GetResult()
		if tcapluserr != 0 {
			fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
			return
		}

		// 从响应消息中提取记录
		fmt.Printf("response success record count %d\n", resp.GetTableRecordCount())
	}
}
