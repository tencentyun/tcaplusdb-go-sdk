package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

//case1 replace result + version success
func TestBatchDeleteSuccess(t *testing.T) {
	client := tools.InitPBSyncClient()

	//batch insert 10 data
	var msgs []proto.Message
	id := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		data.Pay = &tcaplusservice.Payment{Amount: uint64(i), PayId: 2, Method: 3}
		msgs = append(msgs, data)
	}

	err := client.DoBatchInsert(msgs, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}

	//2 batch replace
	opt := &option.PBOpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
		VersionPolicy:        option.CheckDataVersionAutoIncrease,
	}

	for _, msg := range msgs {
		msg.(*tcaplusservice.GamePlayers).Pay.PayId = 3
		opt.BatchVersion = append(opt.BatchVersion, 1)
	}

	err = client.DoBatchDelete(msgs, opt)
	if err != nil {
		t.Errorf("DoBatchDelete fail, %s", err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
	for i, msg := range msgs {
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) ||
			msg.(*tcaplusservice.GamePlayers).Pay.PayId != 2 {
			t.Errorf("DoBatchDelete fail, %+v", msg)
			return
		}
	}

	// 3 batch Get
	opt = &option.PBOpt{}
	var msgs2 []proto.Message
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		msgs2 = append(msgs2, data)
	}

	err = client.DoBatchGet(msgs2, opt)
	if err == nil {
		t.Errorf("DoBatchGet must be 261")
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
}

func TestBatchDeleteVersionFail(t *testing.T) {
	client := tools.InitPBSyncClient()

	//batch insert 10 data
	var msgs []proto.Message
	id := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		data.Pay = &tcaplusservice.Payment{Amount: uint64(i), PayId: 2, Method: 3}
		msgs = append(msgs, data)
	}

	err := client.DoBatchInsert(msgs, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}

	//2 batch replace must error
	opt := &option.PBOpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
		VersionPolicy:        option.CheckDataVersionAutoIncrease,
	}

	for _, msg := range msgs {
		msg.(*tcaplusservice.GamePlayers).Pay.PayId = 3
		opt.BatchVersion = append(opt.BatchVersion, 10)
	}

	err = client.DoBatchDelete(msgs, opt)
	if err == nil {
		t.Errorf("DoBatchDelete version fail, must version error")
		return
	}

	fmt.Println(opt.BatchVersion)
	fmt.Println(opt.BatchResult)
	for i, _ := range msgs {
		if opt.BatchResult[i] == nil {
			t.Errorf("DoBatchDelete fail, must version error  %+v", opt.BatchResult)
			return
		}
	}
}
//case1 replace result + version success
func TestBatchDeleteSuccess_1024(t *testing.T) {
	client := tools.InitPBSyncClient()

	//batch insert 1024 data
	var msgs []proto.Message
	id := time.Now().UnixNano()
	for i := 0; i < 1024; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		data.Pay = &tcaplusservice.Payment{Amount: uint64(i), PayId: 2, Method: 3}
		msgs = append(msgs, data)
	}

	err := client.DoBatchInsert(msgs, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}

	//2 batch delete
	opt := &option.PBOpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
		VersionPolicy:        option.CheckDataVersionAutoIncrease,
	}

	for _, msg := range msgs {
		msg.(*tcaplusservice.GamePlayers).Pay.PayId = 3
		opt.BatchVersion = append(opt.BatchVersion, 1)
	}

	err = client.DoBatchDelete(msgs, opt)
	if err != nil {
		t.Errorf("DoBatchDelete fail, %s", err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
	for i, msg := range msgs {
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) ||
			msg.(*tcaplusservice.GamePlayers).Pay.PayId != 2 {
			t.Errorf("DoBatchDelete fail, %+v", msg)
			return
		}
	}

	// 3 batch Get
	opt = &option.PBOpt{}
	var msgs2 []proto.Message
	for i := 0; i < 1024; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		msgs2 = append(msgs2, data)
	}

	err = client.DoBatchGet(msgs2, opt)
	if err == nil {
		t.Errorf("DoBatchGet must be 261")
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
}

//case1 replace result + version success
func TestBatchDeleteSuccess_Partial_Record(t *testing.T) {
	client := tools.InitPBSyncClient()

	//batch insert 1024 data
	var msgs []proto.Message
	id := time.Now().UnixNano()
	for i := 0; i < 1024; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		data.Pay = &tcaplusservice.Payment{Amount: uint64(i), PayId: 2, Method: 3}
		msgs = append(msgs, data)
	}

	err := client.DoBatchInsert(msgs, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}

	//2 batch delete
	opt := &option.PBOpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
		VersionPolicy:        option.CheckDataVersionAutoIncrease,
	}
	msgs = nil
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		data.Pay = &tcaplusservice.Payment{Amount: uint64(i), PayId: 2, Method: 3}
		msgs = append(msgs, data)
	}
	err = client.DoBatchDelete(msgs, opt)
	if err != nil {
		t.Errorf("DoBatchDelete fail, %s", err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
	for i, msg := range msgs {
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) ||
			msg.(*tcaplusservice.GamePlayers).Pay.PayId != 2 {
			t.Errorf("DoBatchDelete fail, %+v", msg)
			return
		}
	}

	// 3 batch Get
	opt = &option.PBOpt{}
	var msgs2 []proto.Message
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		msgs2 = append(msgs2, data)
	}

	err = client.DoBatchGet(msgs2, opt)
	if err == nil {
		t.Errorf("DoBatchGet must be 261")
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
}