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

func BatchDeleteExample() {
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

	// 生成 batch get 请求
	req, err := client.NewRequest(tools.ZoneId, "game_players", cmd.TcaplusApiBatchDeleteReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return
	}

	// 向请求中添加记录，对于 generic 表 index 无意义，填 0 即可
	record1, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return
	}

	// 向记录中填充数据
	msg1 := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "zhang@test.com",
	}
	// 第一个返回值为记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应，此处无用
	_, err = record1.SetPBData(msg1)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return
	}

	// 向请求中添加记录，对于 generic 表 index 无意义，填 0 即可
	record2, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return
	}

	// 向记录中填充数据
	msg2 := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}
	// 第一个返回值为记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应，此处无用
	_, err = record2.SetPBData(msg2)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return
	}

	// （非必须）设置 异步 id
	req.SetAsyncId(12345)

	// 设置分包
	req.SetMultiResponseFlag(1)

	// 发送请求
	err = client.SendRequest(req)
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}

	wg.Wait()
	logger.INFO("batch get success")
	fmt.Println("batch get success")
}
