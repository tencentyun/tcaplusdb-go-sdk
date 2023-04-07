package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"time"
)

// 针对同一个表只能存在一个遍历器
// 同步请求的遍历只针对数据量较小的表，由于会把所有的数据存放在内存中，占用内存
// 数据量较大的表可以使用异步遍历，边遍历边处理数据
func TraverseExample() {
	// 获取遍历器，遍历器最多同时8个工作，如果超过会返回nil
	tra := client.GetTraverser(tools.Zone, "game_players")
	if tra == nil {
		logger.ERR("GetTraverser fail")
		return
	}
	// 调用stop才能释放资源，防止获取遍历器失败
	defer tra.Stop()
	// （非必须）限制本次遍历记录条数，默认不限制
	tra.SetLimit(1000)

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
