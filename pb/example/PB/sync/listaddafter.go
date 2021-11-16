package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"time"
)

func main() {
	// 创建 client，配置日志，连接数据库
	client := tools.InitPBSyncClient()
	defer client.Close()

	// 生成 listaddafter 请求
	req, err := client.NewRequest(tools.ZoneId, "tb_online_list", cmd.TcaplusApiListAddAfterReq)
	if err != nil {
		logger.ERR("NewRequest error:%s", err)
		return
	}

	// 向请求中添加记录，对于 list 表 index 代表list数组下标，如填 1 表示查到第二条记录之后（下标从0开始）。
	// 插入元素位置在最后面
	//const TCAPLUS_LIST_LAST_INDEX int64 = -1
	// 插入元素位置在最前面
	//const TCAPLUS_LIST_PRE_FIRST_INDEX int64 = -2
	record, err := req.AddRecord(-1)
	if err != nil {
		logger.ERR("AddRecord error:%s", err)
		return
	}

	// 向记录中填充数据
	msg := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}
	client.ListDeleteAll(msg)
	// 第一个返回值为记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应，此处无用
	// key 字段必填，通过 proto 文件设置 key
	// 本例中为 option(tcaplusservice.tcaplus_primary_key) = "openid,tconndid,timekey";
	_, err = record.SetPBData(msg)
	if err != nil {
		logger.ERR("SetPBData error:%s", err)
		return
	}

	// （非必须，默认为 0）insert 请求设置 2 时将返回此次插入的记录，0 1 3 不返回记录
	req.SetResultFlagForSuccess(2)

	// （非必须）设置userbuf，在响应中带回。这个是个开放功能，比如某些临时字段不想保存在全局变量中，
	// 可以通过设置userbuf在发送端接收短传递，也可以起异步id的作用
	req.SetUserBuff([]byte("user buffer test"))

	// （非必须）防止此条记录已存在
	client.ListDeleteAll(msg)

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

	// 获取userbuf
	fmt.Println(string(resp.GetUserBuffer()))

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

		fmt.Println(tools.ConvertToJson(newMsg))
	}

	logger.INFO("listaddafter success")
	fmt.Println("listaddafter success")
}
