package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb"
)

const (
	AppId     = uint64(2)
	ZoneId    = uint32(3)
	DirUrl    = "tcp://x.x.x.x:xxxx"
	Signature = "xxxxxxxxxxxxx"
	TableName = "service_info"
)

var client *tcaplus.Client

func main() {
	client = tcaplus.NewClient()
	if err := client.SetLogCfg("./logconf.xml"); err != nil {
		fmt.Println(err.Error())
		return
	}

	err := client.Dial(AppId, []uint32{ZoneId}, DirUrl, Signature, 60)
	if err != nil {
		fmt.Printf("init failed %v\n", err.Error())
		return
	}
	fmt.Printf("Dial finish\n")
	getPartKeyExample()
	insertExample()
	getExample()
	updateExample()
	replaceExample()
	deleteExample()
}
