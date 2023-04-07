package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"sync"
	"time"
)

func IndexQueryExample() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	// 在另一协程处理响应消息
	go func() {
		defer wg.Done()
		for {
			// resp err 均为 nil 说明响应池中没有任何响应
			resp, err := client.RecvResponse()
			if err != nil {
				logger.ERR("RecvResponse error:%s", err)
				continue
			} else if resp == nil {
				time.Sleep(time.Microsecond * 5)
				continue
			}

			// 获取响应结果
			errCode := resp.GetResult()
			if errCode != terror.GEN_ERR_SUC {
				logger.ERR("insert error:%s", terror.GetErrMsg(errCode))
				return
			}

			// 如果有返回记录则用以下接口进行获取
			for i := 0; i < resp.GetRecordCount(); i++ {
				record, err := resp.FetchRecord()
				if err != nil {
					logger.ERR("FetchRecord failed %s", err.Error())
					return
				}

				newMsg := &tcaplusservice.GamePlayers{}
				err = record.GetPBData(newMsg)
				if err != nil {
					logger.ERR("GetPBData failed %s", err.Error())
					return
				}
				fmt.Println(tools.ConvertToJson(newMsg))
			}
			//判断是否有分包
			if 1 == resp.HaveMoreResPkgs() {
				continue
			}
			return
		}
	}()

	// 生成 index query 请求
	req, err := client.NewRequest(tools.ZoneId, "game_players", cmd.TcaplusApiSqlReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return
	}

	query := fmt.Sprintf("select * from game_players where player_id=10805514 and player_name=Calvin")
	// 设置 sql ，仅用于二级索引请求
	req.SetSql(query)

	// （非必须）设置 异步 id
	req.SetAsyncId(12345)

	// 发送请求
	err = client.SendRequest(req)
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}

	wg.Wait()
	logger.INFO("index query success")
	fmt.Println("index query success")
}
