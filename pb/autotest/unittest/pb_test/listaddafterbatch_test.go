package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

// 全部索引存在
func TestListAddAfterBatchSuccess(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 addafterbatch 10 条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	var msgs []proto.Message
	var indexs []int32
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
		msgs = append(msgs, data)
		indexs = append(indexs, -1)
	}
	opt := &option.PBOpt{}
	err := client.DoListAddAfterBatch(msgs, indexs, opt)
	if err != nil {
		t.Errorf("DoListAddAfterBatch fail, %s", err.Error())
		return
	}
	fmt.Println("DoListAddAfterBatch success")
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println(indexs)

	//BatchGet 10条记录
	indexs = nil
	for i := 0; i < 10; i++ {
		indexs = append(indexs, int32(i))
	}
	msg := &tcaplusservice.TbOnlineList{
		Openid:   1,
		Tconndid: 2,
		Timekey:  key,
	}
	opt = &option.PBOpt{}
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
	//list get all 10Record
	opt = &option.PBOpt{}
	resMsgs, err = client.DoListGetAll(msg, opt)
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
	//list delete all
	err = client.DoListDeleteAll(msg, nil)
	if err != nil {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}

	//list get all empty
	opt = &option.PBOpt{}
	_, err = client.DoListGetAll(msg, opt)
	fmt.Println(err)
	if err == nil {
		t.Errorf("DoListGetBatch fail must empty")
		return
	}
}

func TestListAddAfterBatchVersionFail(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 addafterbatch 10 条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	var msgs []proto.Message
	var indexs []int32
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
		msgs = append(msgs, data)
		indexs = append(indexs, -1)
	}
	opt := &option.PBOpt{}
	err := client.DoListAddAfterBatch(msgs, indexs, opt)
	if err != nil {
		t.Errorf("DoListAddAfterBatch fail, %s", err.Error())
		return
	}
	fmt.Println("DoListAddAfterBatch success")
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println(indexs)

	//version success
	opt = &option.PBOpt{
		Version:       10,
		VersionPolicy: option.CheckDataVersionAutoIncrease,
	}
	indexs = nil
	for i := 0; i < 10; i++ {
		indexs = append(indexs, -1)
	}
	err = client.DoListAddAfterBatch(msgs, indexs, opt)
	if err != nil {
		t.Errorf("DoListAddAfterBatch fail, %s", err.Error())
		return
	}
	fmt.Println("DoListAddAfterBatch success")
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println(indexs)

	//version fail
	opt = &option.PBOpt{
		Version:       100,
		VersionPolicy: option.CheckDataVersionAutoIncrease,
	}
	indexs = nil
	for i := 0; i < 10; i++ {
		indexs = append(indexs, -1)
	}
	err = client.DoListAddAfterBatch(msgs, indexs, opt)
	if err == nil {
		t.Errorf("DoListAddAfterBatch must version failed")
		return
	}
	fmt.Println("DoListAddAfterBatch success")
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println(indexs)
}
// 全部索引存在
func TestListAddAfterBatchSuccess_1023(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 addafterbatch 10 条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	var msgs []proto.Message
	var indexs []int32
	for i := 0; i < 1023; i++ {
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
		msgs = append(msgs, data)
		indexs = append(indexs, -1)
	}
	opt := &option.PBOpt{}
	err := client.DoListAddAfterBatch(msgs, indexs, opt)
	if err != nil {
		t.Errorf("DoListAddAfterBatch fail, %s", err.Error())
		return
	}
	fmt.Println("DoListAddAfterBatch success")
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println("*******************")
	fmt.Println(indexs)

	//BatchGet 10条记录
	indexs = nil
	for i := 0; i < 1023; i++ {
		indexs = append(indexs, int32(i))
	}
	msg := &tcaplusservice.TbOnlineList{
		Openid:   1,
		Tconndid: 2,
		Timekey:  key,
	}
	opt = &option.PBOpt{}
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
	//list get all 10Record
	opt = &option.PBOpt{}
	resMsgs, err = client.DoListGetAll(msg, opt)
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
	//list delete all
	err = client.DoListDeleteAll(msg, nil)
	if err != nil {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}

	//list get all empty
	opt = &option.PBOpt{}
	_, err = client.DoListGetAll(msg, opt)
	fmt.Println(err)
	if err == nil {
		t.Errorf("DoListGetBatch fail must empty")
		return
	}
}

// 全部索引存在
func TestListAddAfterBatchSuccess_1025(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 addafterbatch 1025条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	var msgs []proto.Message
	var indexs []int32
	for i := 0; i < 1025; i++ {
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
		msgs = append(msgs, data)
		indexs = append(indexs, -1)
	}
	opt := &option.PBOpt{}
	err := client.DoListAddAfterBatch(msgs, indexs, opt)
	if err == nil {
		t.Errorf("DoListAddAfterBatch fail, %s", err.Error())
		return
	}
}

// 全部索引存在，addafterbatch 1 条记录
func TestListAddAfterBatchSuccess_01(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 addafterbatch 1 条记录
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	var msgs []proto.Message
	var indexs []int32
	for i := 0; i < 1; i++ {
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
		msgs = append(msgs, data)
		indexs = append(indexs, -1)
	}
	opt := &option.PBOpt{}
	err := client.DoListAddAfterBatch(msgs, indexs, opt)
	if err != nil {
		t.Errorf("DoListAddAfterBatch fail, %s", err.Error())
		return
	}
	fmt.Println("DoListAddAfterBatch success")
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println("*******************")
	fmt.Println(indexs)

	//BatchGet 1条记录
	indexs = nil
	for i := 0; i < 1; i++ {
		indexs = append(indexs, int32(i))
	}
	msg := &tcaplusservice.TbOnlineList{
		Openid:   1,
		Tconndid: 2,
		Timekey:  key,
	}
	opt = &option.PBOpt{}
	resMsgs, err := client.DoListGetBatch(msg, indexs, opt)
	/*
	if err != nil {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}

	 */
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
	//list get all 10Record
	opt = &option.PBOpt{}
	resMsgs, err = client.DoListGetAll(msg, opt)
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
	//list delete all
	err = client.DoListDeleteAll(msg, nil)
	if err != nil {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}

	//list get all empty
	opt = &option.PBOpt{}
	_, err = client.DoListGetAll(msg, opt)
	fmt.Println(err)
	if err == nil {
		t.Errorf("DoListGetBatch fail must empty")
		return
	}
}
