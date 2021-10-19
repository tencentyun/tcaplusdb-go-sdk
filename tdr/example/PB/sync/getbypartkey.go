package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()

	// 生成 get by part key 请求
	req, err := client.NewRequest(tools.ZoneId, "game_players", cmd.TcaplusApiGetByPartkeyReq)
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
	msg := &tcaplusservice.GamePlayers{
		PlayerId:   10805514,
		PlayerName: "Calvin",
	}
	// 第一个返回值为记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应，此处无用
	// 通过部分 key 字段 player_id player_name
	// proto 文件中通过设置 option index 来设置部分key。
	// 本例中使用的是 option(tcaplusservice.tcaplus_index) = "index_1(player_id, player_name)";
	_, err = record.SetPBPartKeys(msg, []string{"player_id", "player_name"})
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return
	}

	// （非必须）对返回记录做限制，此处为从第0条开始返回1条记录
	// req.SetResultLimit(1, 0)

	// （非必须）设置userbuf，在响应中带回。这个是个开放功能，比如某些临时字段不想保存在全局变量中，
	// 可以通过设置userbuf在发送端接收短传递，也可以起异步id的作用
	req.SetUserBuff([]byte("user buffer test"))

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
	client.Insert(&tcaplusservice.GamePlayers{
		PlayerId:        10805514,
		PlayerName:      "Calvin",
		PlayerEmail:     "zhang@test.com",
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

		// 获取userbuf
		fmt.Println(string(resp.GetUserBuffer()))

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

	logger.INFO("get by part key success")
	fmt.Println("get by part key success")
}
