package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()

	// 获取遍历器，遍历器最多同时8个工作，如果超过会返回nil
	tra := client.GetTraverser(tools.ZoneId, "game_players")
	if tra == nil {
		logger.ERR("GetTraverser fail")
		return
	}
	// 调用stop才能释放资源，防止获取遍历器失败
	defer tra.Stop()

	// （非必须）限制本次遍历记录条数，默认不限制
	tra.SetLimit(1000)

	// （非必须）设置userbuf，在响应中带回。这个是个开放功能，比如某些临时字段不想保存在全局变量中，
	// 可以通过设置userbuf在发送端接收短传递，也可以起异步id的作用
	tra.SetUserBuff([]byte("user buffer test"))

	// （非必须） 防止记录不存在
	client.Insert(&tcaplusservice.GamePlayers{
		PlayerId:        10805514,
		PlayerName:      "Calvin",
		PlayerEmail:     "calvin@test.com",
		GameServerId:    10,
		LoginTimestamp:  []string{"2019-12-12 15:00:00"},
		LogoutTimestamp: []string{"2019-12-12 16:00:00"},
		IsOnline:        false,
		Pay: &tcaplusservice.Payment{
			PayId:  10101,
			Amount: 1000,
			Method: 2,
		},
	})

	// 发送请求
	err := tra.Start()
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}

	resps, err := client.DoTraverse(tra, 30*time.Second)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, res := range resps {
		ret := res.GetResult()
		if ret != 0 {
			logger.ERR("result is %d, error:%s", ret, terror.GetErrMsg(ret))
			continue
		}

		for i := 0; i < res.GetRecordCount(); i++ {
			record, err := res.FetchRecord()
			if err != nil {
				logger.ERR("FetchRecord error:%s", err)
				continue
			}

			newMsg := &tcaplusservice.GamePlayers{}
			err = record.GetPBData(newMsg)
			if err != nil {
				logger.ERR("GetPBData failed %s", err.Error())
				return
			}

			fmt.Println(tools.ConvertToJson(newMsg))
		}
	}

}
