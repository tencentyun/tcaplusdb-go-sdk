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

func increaseExample() {
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
				data := tcaplus_tb.NewTable_Generic()
				if err := record.GetData(data); err != nil {
					fmt.Printf("record.GetData failed %s\n", err.Error())
					return
				}
				if data.Level != 4 {
					fmt.Printf("data.Level invalid %v", data.Level)
					return
				}

				if data.Float_Score != 12.6 {
					fmt.Printf("data.Float_Score invalid %v", data.Float_Score)
					return
				}

				if data.Double_Score != 16.8 {
					fmt.Printf("data.Double_Score invalid %v", data.Double_Score)
					return
				}
			}
			return
		}
	}()
	req, err := client.NewRequest(ZoneId, "table_generic", cmd.TcaplusApiIncreaseReq)
	if err != nil {
		fmt.Printf("NewRequest TcaplusApiReplaceReq failed %v\n", err.Error())
		return
	}

	if err := req.SetResultFlag(2); err != nil {
		fmt.Printf("SetResultFlag failed %v\n", err.Error())
		return
	}

	//为request添加一条记录，（index只有在list表中支持，generic不校验）
	rec, err := req.AddRecord(0)
	if err != nil {
		fmt.Printf("AddRecord failed %v\n", err.Error())
		return
	}

	//申请tdr结构体并赋值，最好调用tdr pkg的NewXXX函数，会将成员初始化为tdr定义的tdr默认值，
	// 不要自己new，自己new，某些结构体未初始化，存在panic的风险
	data := tcaplus_tb.NewTable_Generic()
	data.Uin = uint64(time.Now().UnixNano())
	data.Name = "GoUnitTest"
	data.Key3 = "key3"
	data.Key4 = "key4"
	data.Info = "info"
	data.Name = fmt.Sprintf("%d", 2)
	data.Level = int32(2)
	data.Float_Score = float32(6.6)
	data.Double_Score = float64(8.8)
	data.Info = fmt.Sprintf("%d", 2)
	//将tdr的数据设置到请求的记录中
	if err := rec.SetData(data); err != nil {
		fmt.Printf("SetData failed %v\n", err.Error())
		return
	}

	rec.SetIncValue("level", int32(2), cmd.TcaplusApiOpPlus, 0, 0)
	rec.SetIncValue("float_score", float32(6), cmd.TcaplusApiOpPlus, 0, 0)
	rec.SetIncValue("double_score", float32(8), cmd.TcaplusApiOpPlus, 0, 0)

	if err := client.SendRequest(req); err != nil {
		fmt.Printf("SendRequest failed %v\n", err.Error())
		return
	}
	wg.Wait()
}