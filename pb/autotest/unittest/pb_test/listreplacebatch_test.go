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
func TestListReplaceBatchSuccess(t *testing.T) {
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

	//2replace + old flag
	msgs = nil
	indexs = nil
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   key,
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i + 10),
			},
		}
		fmt.Println(data)
		msgs = append(msgs, data)
		indexs = append(indexs, int32(i))
	}
	opt = &option.PBOpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
	}
	err = client.DoListReplaceBatch(msgs, indexs, opt)
	if err != nil {
		t.Errorf("DoListReplaceBatch fail, %s", err.Error())
		return
	}
	fmt.Println("DoListReplaceBatch success")
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println(indexs)
	for i, msg := range msgs {
		fmt.Println(msg)
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(i) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(i) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}

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
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(index+10) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(index) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}
}

func TestListReplaceBatchVersionFail(t *testing.T) {
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
		indexs = append(indexs, int32(i))
	}
	err = client.DoListReplaceBatch(msgs, indexs, opt)
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
		indexs = append(indexs, int32(i))
	}
	err = client.DoListReplaceBatch(msgs, indexs, opt)
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
func TestListReplaceBatchSuccess_1023(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 addafterbatch 1023条记录
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
	fmt.Println(indexs)

	//2replace + old flag
	msgs = nil
	indexs = nil
	for i := 0; i < 1023; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   key,
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i + 10),
			},
		}
		fmt.Println(data)
		msgs = append(msgs, data)
		indexs = append(indexs, int32(i))
	}
	opt = &option.PBOpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
	}
	err = client.DoListReplaceBatch(msgs, indexs, opt)
	if err != nil {
		t.Errorf("DoListReplaceBatch fail, %s", err.Error())
		return
	}
	fmt.Println("DoListReplaceBatch success")
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	fmt.Println(indexs)
	for i, msg := range msgs {
		fmt.Println(msg)
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(i) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(i) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}

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
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(index+10) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(index) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}
}

// 全部索引存在,replace 记录不存在检查是否是插入记录
//list replce 就是update generic replace如果记录不存在就会插入新的数据。
func TestListReplaceBatchSuccess_10(t *testing.T) {
	client := tools.InitPBSyncClient()
	uin := time.Now().UnixNano()
	key := fmt.Sprintf("%d", uin)
	var msgs []proto.Message
	var indexs []int32
	opt := &option.PBOpt{}
	msgs = nil
	indexs = nil
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   key,
			Gamesvrid: "lol",
			Pay: &tcaplusservice.TbOnlineListPayInfo{
				TotalMoney: uint64(i),
				PayTimes:   uint64(i + 10),
			},
		}
		fmt.Println(data)
		msgs = append(msgs, data)
		indexs = append(indexs, int32(i))
	}
	opt = &option.PBOpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
	}
	err := client.DoListReplaceBatch(msgs, indexs, opt)
	if err == nil {
		t.Errorf("DoListReplaceBatch fail, %s", err.Error())
		return
	}


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
	if err.Error() != "errCode: 261, errMsg: txhdb_record_not_exist"  {
		t.Errorf("DoListGetBatch fail, %s", err.Error())
		return
	}
	fmt.Println(opt.Version)
	fmt.Println(opt.BatchResult)
	for index, msg := range resMsgs {
		fmt.Println("index", index)
		fmt.Println(msg)
		if msg.(*tcaplusservice.TbOnlineList).Pay.PayTimes != uint64(index+10) {
			t.Errorf("Pay.PayTimes invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
		if msg.(*tcaplusservice.TbOnlineList).Pay.TotalMoney != uint64(index) {
			t.Errorf("Pay.TotalMoney invalid, %v", msg.(*tcaplusservice.TbOnlineList))
			return
		}
	}
}