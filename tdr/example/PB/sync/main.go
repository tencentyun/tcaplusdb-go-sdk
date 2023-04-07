package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
)

var client *tcaplus.PBClient

func main() {
	// 创建 client，配置日志，连接数据库
	client = tools.NewPBClient()
	if client == nil {
		fmt.Println("NewPBClient failed")
		return
	}
	GetExample()
}