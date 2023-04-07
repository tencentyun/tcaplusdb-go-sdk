package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"strconv"
	"sync"
	"time"
)

func listAddafterBatchExample() {
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

			//带回请求的异步ID
			fmt.Printf("resp success, AsyncId:%d\n", resp.GetAsyncId())
			tcapluserr := resp.GetResult()
			if tcapluserr != 0 {
				fmt.Printf("response ret %s\n",
					"errCode: "+strconv.Itoa(tcapluserr)+", errMsg: "+terror.ErrorCodes[tcapluserr])
				return
			}
			//response中带有获取的记录
			fmt.Printf("response success record count %d\n", resp.GetRecordCount())
			for i := 0; i < resp.GetRecordCount(); i++ {
				record, err := resp.FetchRecord()
				if err != nil {
					fmt.Printf("FetchRecord failed %s\n", err.Error())
					return
				}
				//通过GetData获取记录
				data := tcaplus_tb.NewTable_Traverser_List()
				if err := record.GetData(data); err != nil {
					fmt.Printf("record.GetData failed %s\n", err.Error())
					return
				}
				fmt.Printf("response record data %+v\n", data)
			}
			//判断是否有分包
			if 1 == resp.HaveMoreResPkgs() {
				continue
			}
			return
		}
	}()
	//创建请求
	req, err := client.NewRequest(ZoneId, TABLE_TRAVERSER_LIST, cmd.TcaplusApiListAddAfterBatchReq)
	if err != nil {
		fmt.Printf("NewRequest failed %v\n", err.Error())
		return
	}

	//允许分包
	req.SetMultiResponseFlag(1)
	// 设置ResultFlag标志。0标志只需返回成功与否,1标志返回同请求一致的值,2标志返回操作后所有字段的值,3标志返回操作前所有字段的值
	req.SetResultFlagForSuccess(2)

	for i := 0; i < 10; i++ {
		//TCAPLUS_API_LIST_PRE_FIRST_INDEX(-2)：新元素插入在第一个元素之前
		//TCAPLUS_API_LIST_LAST_INDEX(-1)：新元素插入在最后一个元素之后
		rec, err := req.AddRecord(-1)
		if err != nil {
			fmt.Printf("AddRecord failed %v\n", err.Error())
			return
		}

		//申请tdr结构体并赋值，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
		// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
		data := tcaplus_tb.NewTable_Traverser_List()
		data.Key = 1
		data.Name = 255
		data.Level = uint32(i)
		data.Value1 = "value1"
		data.Value2 = "value2"
		//将tdr的数据设置到请求的记录中
		if err := rec.SetData(data); err != nil {
			fmt.Printf("SetData failed %v\n", err.Error())
			return
		}
	}
	//发送请求
	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	wg.Wait()
}