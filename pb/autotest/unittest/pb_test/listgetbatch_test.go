package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"testing"
	"time"
)

const TB_ONLINE_LIST = "tb_online_list"

// 全部索引存在
func TestListGetBatchSuccess(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 add 10 条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   key,
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i),
			},
		}
		fmt.Println(data)
		_, err := client.DoListAddAfter(data, -1, nil)
		if err != nil {
			t.Errorf("DoListAddAfter fail, %s", err.Error())
			return
		}
	}

	//BatchGet 10条记录
	var indexs []int32
	for i := 0; i < 10; i++ {
		indexs = append(indexs, int32(i))
	}
	msg := &tcaplusservice.TbOnlineList{
		Openid:   1,
		Tconndid: 2,
		Timekey:  key,
	}
	opt := &option.PBOpt{}
	resMsgs, err := client.DoListGetBatch(msg, indexs, opt)
	if err != nil {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}

	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	for index, msg := range resMsgs {
		fmt.Println("index", index)
		fmt.Println(msg)
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(index) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(index) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}
}

func TestListGetBatchFailed(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 add 10 条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   key,
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i),
			},
		}
		fmt.Println(data)
		_, err := client.DoListAddAfter(data, -1, nil)
		if err != nil {
			t.Errorf("DoListAddAfter fail, %s", err.Error())
			return
		}
	}

	//BatchGet 10条记录,索引不存在
	var indexs []int32
	for i := 10; i < 20; i++ {
		indexs = append(indexs, int32(i))
	}
	msg := &tcaplusservice.TbOnlineList{
		Openid:   1,
		Tconndid: 2,
		Timekey:  key,
	}
	opt := &option.PBOpt{}
	resMsgs, err := client.DoListGetBatch(msg, indexs, opt)
	if err == nil {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}

	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	for index, msg := range resMsgs {
		fmt.Println("index", index)
		fmt.Println(msg)
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(index) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(index) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}

}
// 全部索引存在, batch get部分记录成功, 前几条记录
func TestListGetBatchPartSuccess(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 add 10 条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   key,
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i),
			},
		}
		fmt.Println(data)
		_, err := client.DoListAddAfter(data, -1, nil)
		if err != nil {
			t.Errorf("DoListAddAfter fail, %s", err.Error())
			return
		}
	}

	//BatchGet 5条记录
	var indexs []int32
	for i := 0; i < 5; i++ {
		indexs = append(indexs, int32(i))
	}
	msg := &tcaplusservice.TbOnlineList{
		Openid:   1,
		Tconndid: 2,
		Timekey:  key,
	}
	opt := &option.PBOpt{}
	resMsgs, err := client.DoListGetBatch(msg, indexs, opt)
	if err != nil {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}

	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	for index, msg := range resMsgs {
		fmt.Println("index", index)
		fmt.Println(msg)
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(index) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(index) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}
}
// 全部索引存在, batch get部分记录成功, 后几条记录
func TestListGetBatchPartSuccess_01(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 add 10 条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   key,
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i),
			},
		}
		fmt.Println(data)
		_, err := client.DoListAddAfter(data, -1, nil)
		if err != nil {
			t.Errorf("DoListAddAfter fail, %s", err.Error())
			return
		}
	}

	//BatchGet 5条记录
	var indexs []int32
	for i := 5; i < 10; i++ {
		indexs = append(indexs, int32(i))
	}
	msg := &tcaplusservice.TbOnlineList{
		Openid:   1,
		Tconndid: 2,
		Timekey:  key,
	}
	opt := &option.PBOpt{}
	resMsgs, err := client.DoListGetBatch(msg, indexs, opt)
	if err != nil {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}

	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	for index, msg := range resMsgs {
		fmt.Println("index", index)
		fmt.Println(msg)
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(index) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(index) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}
}

// 全部索引存在, batch get部分记录成功, 取中间几条记录
func TestListGetBatchPartSuccess_02(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 add 10 条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   key,
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i),
			},
		}
		fmt.Println(data)
		_, err := client.DoListAddAfter(data, -1, nil)
		if err != nil {
			t.Errorf("DoListAddAfter fail, %s", err.Error())
			return
		}
	}

	//BatchGet 5条记录
	var indexs []int32
	for i := 3; i < 8; i++ {
		indexs = append(indexs, int32(i))
	}
	msg := &tcaplusservice.TbOnlineList{
		Openid:   1,
		Tconndid: 2,
		Timekey:  key,
	}
	opt := &option.PBOpt{}
	resMsgs, err := client.DoListGetBatch(msg, indexs, opt)
	if err != nil {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}

	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	for index, msg := range resMsgs {
		fmt.Println("index", index)
		fmt.Println(msg)
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(index) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(index) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}
}