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

func ListGetAllExample() {
	// 生成请求
	req, err := client.NewRequest(tools.Zone, "tb_online_list", cmd.TcaplusApiListGetAllReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return
	}

	record, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return
	}
	//允许分包
	req.SetMultiResponseFlag(1)

	// 向记录中填充数据
	msg := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}

	// 第一个返回值为记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应，此处无用
	// key 字段必填，通过 proto 文件设置 key
	// 本例中为 option(tcaplusservice.tcaplus_primary_key) = "openid,tconndid,timekey";
	_, err = record.SetPBData(msg)
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

			newMsg := &tcaplusservice.TbOnlineList{}
			err = record.GetPBData(newMsg)
			if err != nil {
				logger.ERR("GetPBData failed %s", err.Error())
				return
			}
			index := record.GetIndex()
			fmt.Println(tools.ConvertToJson(newMsg), index)
		}
	}

	logger.INFO("listgetall success")
	fmt.Println("listgetall success")
}
