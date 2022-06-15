package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"google.golang.org/protobuf/proto"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()
	client.SetDefaultZoneId(tools.ZoneId)

	//1 addafterbatch 10 条记录
	var msgs []proto.Message
	var indexs []int32
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   "key",
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i),
			},
		}
		fmt.Println(data)
		msgs = append(msgs, data)
		indexs = append(indexs, -1)
	}
	opt := &option.PBOpt{}
	err := client.DoListAddAfterBatch(msgs, indexs, opt)
	if err != nil {
		fmt.Printf("DoListAddAfterBatch fail, %s", err.Error())
		return
	}
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println(indexs)
	fmt.Println("DoListAddAfterBatch success")
}