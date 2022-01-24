package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"google.golang.org/protobuf/proto"
	"strings"
	"testing"
	"time"
)

//case1 replace result + version success
func TestBatchUpdateSuccess(t *testing.T) {
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

	err = client.DoBatchUpdate(msgs, opt)
	if err != nil {
		t.Errorf("DoBatchUpdate fail, %s", err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
	for i, msg := range msgs {
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) ||
			msg.(*tcaplusservice.GamePlayers).Pay.PayId != 2 {
			t.Errorf("DoBatchGet fail, %+v", msg)
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
	if err != nil {
		t.Errorf("DoBatchGet fail, %s", err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
	for i, msg := range msgs2 {
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) ||
			msg.(*tcaplusservice.GamePlayers).Pay.PayId != 3 {
			t.Errorf("DoBatchGet fail, %+v", msg)
			return
		}
	}
}

func TestBatchUpdateVersionFail(t *testing.T) {
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

	err = client.DoBatchUpdate(msgs, opt)
	if err == nil {
		t.Errorf("DoBatchUpdate version fail, must version error")
		return
	}

	fmt.Println(opt.BatchVersion)
	fmt.Println(opt.BatchResult)
	for i, _ := range msgs {
		if opt.BatchResult[i] == nil {
			t.Errorf("DoBatchUpdate fail, must version error  %+v", opt.BatchResult)
			return
		}
	}
}
//case1 replace result + version success
func TestBatchUpdateSuccess_1024(t *testing.T) {
	client := tools.InitPBSyncClient()

	//batch insert 10 data
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

	//2 batch replace
	opt := &option.PBOpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllOldValue,
		VersionPolicy:        option.CheckDataVersionAutoIncrease,
	}

	for _, msg := range msgs {
		msg.(*tcaplusservice.GamePlayers).Pay.PayId = 3
		opt.BatchVersion = append(opt.BatchVersion, 1)
	}

	err = client.DoBatchUpdate(msgs, opt)
	if err != nil {
		t.Errorf("DoBatchUpdate fail, %s", err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
	for i, msg := range msgs {
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) ||
			msg.(*tcaplusservice.GamePlayers).Pay.PayId != 2 {
			t.Errorf("DoBatchGet fail, %+v", msg)
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
	if err != nil {
		t.Errorf("DoBatchGet fail, %s", err.Error())
		return
	}

	fmt.Println(opt.BatchResult)
	fmt.Println(opt.BatchVersion)
	for i, msg := range msgs2 {
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) ||
			msg.(*tcaplusservice.GamePlayers).Pay.PayId != 3 {
			t.Errorf("DoBatchGet fail, %+v", msg)
			return
		}
	}
}


//case1 replace result + version success
func TestBatchUpdateSuccess_1025(t *testing.T) {
	client := tools.InitPBSyncClient()

	//batch insert 10 data
	var msgs []proto.Message
	id := time.Now().UnixNano()
	for i := 0; i < 1025; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		data.Pay = &tcaplusservice.Payment{Amount: uint64(i), PayId: 2, Method: 3}
		msgs = append(msgs, data)
	}

	err := client.DoBatchInsert(msgs, nil)
	if !strings.Contains(err.Error(),"-4126") {
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

	err = client.DoBatchUpdate(msgs, opt)
	if !strings.Contains(err.Error(),"-4126") {
		t.Errorf("DoBatchUpdate fail, %s", err.Error())
		return
	}
}
//记录不存在的时候batch update
//case1 update result + version success
func TestBatchUpdateFail(t *testing.T) {
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
		msg.(*tcaplusservice.GamePlayers).PlayerId = 100
		msg.(*tcaplusservice.GamePlayers).Pay.PayId = 3
		opt.BatchVersion = append(opt.BatchVersion, 1)
	}


	err = client.DoBatchUpdate(msgs, opt)
	if err.Error() != "errCode: 261, errMsg: txhdb_record_not_exist"{
		t.Errorf("DoBatchUpdate fail, %s", err.Error())
		return
	}
}