package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb"
)

const (
	AppId                = uint64(2)
	ZoneId               = uint32(3)
	DirUrl               = "tcp://x.x.x.x:9999"
	Signature            = "xxxxxx"
	TableName            = "service_info"
	TABLE_TRAVERSER_LIST = "table_traverser_list"
)

var client *tcaplus.Client

func main() {
	client = tcaplus.NewClient()
	//日志配置，不配置则debug打印到控制台
	if err := client.SetLogCfg("./logconf.xml"); err != nil {
		fmt.Println(err.Error())
		return
	}
	//client连接tcaplus
	err := client.Dial(AppId, []uint32{ZoneId}, DirUrl, Signature, 60)
	if err != nil {
		fmt.Printf("init failed %v\n", err.Error())
		return
	}
	fmt.Printf("Dial finish\n")
	getPartKeyExample()
	deleteByPartKeyExample()
	insertExample()
	getExample()
	updateExample()
	replaceExample()
	deleteExample()
}
