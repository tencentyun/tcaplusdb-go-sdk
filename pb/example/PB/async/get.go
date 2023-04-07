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

func GetExample() {
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
		}
	}()

	// 生成 get 请求
	req, err := client.NewRequest(tools.ZoneId, "game_players", cmd.TcaplusApiGetReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return
	}

	// 向请求中添加记录，对于 generic 表 index 无意义，填 0 即可
	record, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return
	}

	// 向记录中填充数据
	// key 字段必填，通过 proto 文件设置 key
	// 本例中为 option(tcaplusservice.tcaplus_primary_key) = "player_id, player_name, player_email";
	msg := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}

	_, err = record.SetPBData(msg)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return
	}

	// 异步发送请求
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	wg.Wait()
}