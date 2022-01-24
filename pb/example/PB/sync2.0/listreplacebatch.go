package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"google.golang.org/protobuf/proto"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()
	client.SetDefaultZoneId(tools.ZoneId)

	var indexs []int32
	var msgs []proto.Message
	indexs = nil
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   "key",
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i + 10),
			},
		}
		fmt.Println(data)
		msgs = append(msgs, data)
		indexs = append(indexs, int32(i))
	}
	opt := &option.PBOpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
	}
	err := client.DoListReplaceBatch(msgs, indexs, opt)
	if err != nil {
		fmt.Printf("DoListReplaceBatch fail, %s", err.Error())
		return
	}

	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println(indexs)
	for _, msg := range msgs {
		//old msg
		fmt.Println(msg)
	}
	fmt.Println("DoListReplaceBatch success")
}
