package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr"
)

const (
	AppId                = uint64(2)
	ZoneId               = uint32(3)
	DirUrl               = "tcp://x.x.x.x:xxxx"
	Signature            = "xxxxxxxxxxxxx"
	TableName            = "service_info"
	TABLE_TRAVERSER_LIST = "table_traverser_list"
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
	client.SetDefaultZoneId(ZoneId)
	fmt.Printf("Dial finish\n")
	insertExample()
	getPartKeyExample()
	getExample()
	ttlExample()
	updateExample()
	replaceExample()
	deleteExample()

	batchInsertExample()
	batchGetExample()
	batchReplaceExample()
	batchUpdateExample()
	batchDeleteExample()
}
