package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

func BatchReplaceExample() {
	// 生成 batch 请求
	req, err := client.NewRequest(tools.Zone, "game_players", cmd.TcaplusApiBatchReplaceReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return
	}
	// 设置分包
	req.SetMultiResponseFlag(1)
	// 向请求中添加记录，对于 generic 表 index 无意义，填 0 即可
	record1, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return
	}

	// 向记录中填充数据
	msg1 := &tcaplusservice.GamePlayers{
		PlayerId:        10805515,
		PlayerName:      "Calvin1",
		PlayerEmail:     "calvin@test.com",
		GameServerId:    10,
		LoginTimestamp:  []string{"2019-12-12 15:00:00"},
		LogoutTimestamp: []string{"2019-12-12 16:00:00"},
		IsOnline:        false,
		Pay: &tcaplusservice.Payment{
			PayId:  10101,
			Amount: 1000,
			Method: 1,
		},
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
		PlayerId:        10805515,
		PlayerName:      "Calvin2",
		PlayerEmail:     "calvin@test.com",
		GameServerId:    10,
		LoginTimestamp:  []string{"2019-12-12 15:00:00"},
		LogoutTimestamp: []string{"2019-12-12 16:00:00"},
		IsOnline:        false,
		Pay: &tcaplusservice.Payment{
			PayId:  10101,
			Amount: 1000,
			Method: 1,
		},
	}
	// 第一个返回值为记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应，此处无用
	_, err = record2.SetPBData(msg2)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return
	}

	// 发送请求,接收响应
	resps, err := client.DoMore(req, 5*time.Second)
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}

	for _, resp := range resps {
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

	logger.INFO("batch success")
	fmt.Println("batch success")
}
