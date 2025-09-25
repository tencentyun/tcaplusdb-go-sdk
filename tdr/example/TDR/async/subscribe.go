package main

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/subscribe"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"strconv"
	"time"
)

// NewSubscribeTopic 创建订阅topic; 只支持list表; 以list的key为订阅topic，当往这个list中插入数据时会收到订阅响应
// 注意: list要尾部插入，淘汰头部，保持list的index是连续的，这样基于index, sdk会做丢包重试
//  index: 订阅的起始index
//         < 0 表示只订阅不拉取数据，
//         >=0 订阅并返回该index之后的所有元素（包括该index）; 如果index不存在，则返回错误，并取消订阅
//         用户不确定从哪个index开始订阅；可以先通过listgetall拉取，再选择index开始订阅
//  expireSecond: 订阅过期时间，秒
//          sdk会自动取该值的三分之一为续订阅时间，sdk定时续订，维持订阅链路；这样sdk异常退出时，svr可以删除超时链路
//  rspPipeSize: topic响应队列channel大小，最小为1
//          如果channel满，则会丢掉订阅包;等到不满时，会拉取后续订阅包，填充队列，sdk会保证收到的订阅包仍然是连续的
//  异常场景：
//    长时间无订阅包：
//          如果一个expireSecond时间，没有收到过订阅响应，sdk也会尝试拉取当前index之后的包，检查是否有遗漏
//    丢包：
//          sdk根据list返回的index是否连续自增判断是否丢包；识别丢包后会自动携带该index去svr拉取；
//          故对于订阅的list，请不要删除中间的元素导致index不连续，不连续会导致订阅中断
//    订阅中断：
//          订阅的index不存在会中断，或者index不连续
//          中断时，topic的ResponsePipe channel会先收到一个带错误码的响应包，然后会收到channel close

// subscribeExample1 直接监听topic的channel快速收订阅包
func subscribeExample1() {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255

	topic, err := client.NewSubscribeTopic(TABLE_TRAVERSER_LIST, data, -1, 30, 10, subscribe.RecordFlag, ZoneId)
	if err != nil {
		fmt.Printf("NewSubscribeTopic failed %v\n", err.Error())
		return
	}
	defer client.DelSubscribeTopic(topic)

	for {
		select {
		case resp, ok := <-topic.ResponsePipe:
			if !ok {
				fmt.Printf("topic %s closed \n", topic.GetKeyDebug())
				return
			}
			Print(resp)
		}
	}

	// 大量topic如下收包
	/*
		topicSlice := make([]*subscribe.Topic, 0, 10)
		topicSlice = append(topicSlice, topic)
		for {
			noRsp := true
			// 可以遍历多个topic的channel收包，channel无包时，select default continue到下一个topic
			for _, tp := range topicSlice {
				select {
				case resp, ok := <-tp.ResponsePipe:
					noRsp = false
					if !ok {
						fmt.Printf("topic %s closed \n", tp.GetKeyDebug())
						return
					}
					Print(resp)
				default:
					continue
				}
			}
			// 无包sleep
			if noRsp {
				time.Sleep(10 * time.Microsecond)
			}
		}
	*/
}

// subscribeExample2 通过peek pop按需收订阅包，收的包会缓存在topic队列
// 满时丢掉订阅包，等到不满时，会拉取后续订阅包，填充队列，sdk会保证收到的订阅包仍然是连续的
func subscribeExample2() {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255

	topic, err := client.NewSubscribeTopic(TABLE_TRAVERSER_LIST, data, -1, 300, 1024, subscribe.RecordFlag, ZoneId)
	if err != nil {
		fmt.Printf("NewSubscribeTopic failed %v\n", err.Error())
		return
	}
	defer client.DelSubscribeTopic(topic)

	for {
		rsp, err := topic.PeekRsp()
		if err != nil {
			fmt.Printf("topic %s PeekRsp error %s\n", topic.GetKeyDebug(), err.Error())
			return
		}

		if rsp == nil {
			time.Sleep(10 * time.Microsecond)
			continue
		}

		Print(rsp)
		_, err = topic.PopRsp()
		if err != nil {
			fmt.Printf("topic %s PopRsp error %s\n", topic.GetKeyDebug(), err.Error())
			return
		}
	}
}

// 只更新topic index模式，通过index拉取数据
func subscribeExample3() {
	data := tcaplus_tb.NewTable_Traverser_List()
	data.Key = 1
	data.Name = 255

	topic, err := client.NewSubscribeTopic(TABLE_TRAVERSER_LIST, data, -1, 30, 10, subscribe.OnlyIndexFlag, ZoneId)
	if err != nil {
		fmt.Printf("NewSubscribeTopic failed %v\n", err.Error())
		return
	}
	defer client.DelSubscribeTopic(topic)

	curIdx := int32(-1)
	for {
		select {
		//此处只有异常退出时，ResponsePipe才会有包
		case resp, ok := <-topic.ResponsePipe:
			if !ok {
				fmt.Printf("topic %s closed \n", topic.GetKeyDebug())
				return
			}
			Print(resp)
		default:
			firstIdx := topic.GetFirstIndex()
			lastIdx := topic.GetLastIndex()
			if firstIdx == -1 || lastIdx == -1 || curIdx == lastIdx {
				// 无订阅包
				time.Sleep(10 * time.Microsecond)
				continue
			}
			fmt.Printf("first idx %d - idx %d ,cur idx %d\n", firstIdx, lastIdx, curIdx)
			for {
				if curIdx == -1 {
					curIdx = firstIdx
				} else {
					curIdx = common.GetSubscribeInt32Next(curIdx)
				}

				// 拉取数据
				err := client.DoListGet(TABLE_TRAVERSER_LIST, data, curIdx, nil)
				if err != nil {
					fmt.Printf("err %s\n", err.Error())
				} else {
					fmt.Printf("response record data %+v index %v\n", data, curIdx)
				}

				// 拉到最后一条
				if curIdx == lastIdx {
					break
				}
			}
		}
	}
}

func Print(resp response.TcaplusResponse) {
	if resp == nil {
		return
	}

	errCode := resp.GetResult()
	if errCode != 0 {
		fmt.Printf("response ret %s\n",
			"errCode: "+strconv.Itoa(errCode)+", errMsg: "+terror.ErrorCodes[errCode])
		return
	}

	//response中带有获取的记录
	fmt.Printf("response success record count %d\n", resp.GetRecordCount())
	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			fmt.Printf("FetchRecord failed %s\n", err.Error())
			continue
		}
		//通过GetData获取记录
		data := tcaplus_tb.NewTable_Traverser_List()
		if err := record.GetData(data); err != nil {
			fmt.Printf("record.GetData failed %s\n", err.Error())
			continue
		}
		fmt.Printf("response record data %+v index %v\n", data, record.GetIndex())

		// 如果需要删除，则调用如下接口删除已消费的元素
		//if err := client.DoListDelete(TABLE_TRAVERSER_LIST, data, record.GetIndex(), nil); err != nil{
		//	fmt.Printf("response record data %+v index %v delete fail %s\n", data, record.GetIndex(), err.Error())
		//}
	}
}
