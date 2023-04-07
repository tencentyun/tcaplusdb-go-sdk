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

func IndexQueryExample() {
	// 生成 index query 请求
	req, err := client.NewRequest(tools.Zone, "game_players", cmd.TcaplusApiSqlReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return
	}

	query := fmt.Sprintf("select * from game_players where player_id=10805514 and player_name=Calvin")
	// 设置 sql ，仅用于二级索引请求
	req.SetSql(query)

	// 发送请求,接收响应
	resp, err := client.Do(req, 5*time.Second)
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}

	// 获取响应结果
	errCode := resp.GetResult()
	if errCode != terror.GEN_ERR_SUC {
		logger.ERR("insert error:%s", terror.GetErrMsg(errCode))
		return
	}

	// 获取查询类型，仅用于 二级索引接口
	fmt.Println(resp.GetSqlType())

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

	query = fmt.Sprintf("select count(*) from game_players where player_id=10805514 and player_name=Calvin")
	// 设置 sql ，仅用于二级索引请求
	req.SetSql(query)
	// 发送请求,接收响应
	resp, err = client.Do(req, 5*time.Second)
	if err != nil {
		logger.ERR("SendRequest error:%s", err)
		return
	}
	// 获取响应结果
	errCode = resp.GetResult()
	if errCode != terror.GEN_ERR_SUC {
		logger.ERR("insert error:%s", terror.GetErrMsg(errCode))
		return
	}
	// 获取查询类型，仅用于 二级索引接口
	fmt.Println(resp.GetSqlType())
	// 获取聚合查询结果
	fmt.Println(resp.ProcAggregationSqlQueryType())

	logger.INFO("index query success")
	fmt.Println("index query success")
}
