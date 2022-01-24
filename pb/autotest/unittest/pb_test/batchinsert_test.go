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

//case1 insert success
func TestBatchInsertSuccess(t *testing.T) {
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
	// 2 batch Get
	opt := &option.PBOpt{}
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
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) {
			t.Errorf("DoBatchGet fail, %+v", msg)
			return
		}
	}
}
//case1 insert success
func TestBatchInsertSuccess_1(t *testing.T) {
	client := tools.InitPBSyncClient()

	//batch insert 1 data
	var msgs []proto.Message
	id := time.Now().UnixNano()
	for i := 0; i < 1; i++ {
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
	// 2 batch Get
	opt := &option.PBOpt{}
	var msgs2 []proto.Message
	for i := 0; i < 1; i++ {
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
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) {
			t.Errorf("DoBatchGet fail, %+v", msg)
			return
		}
	}
}
//case1 insert success
func TestBatchInsertSuccess_1024(t *testing.T) {
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
	// 2 batch Get
	opt := &option.PBOpt{}
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
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) {
			t.Errorf("DoBatchGet fail, %+v", msg)
			return
		}
	}
}

//case1 insert success
func TestBatchInsertSuccess_1025(t *testing.T) {
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
	// 2 batch Get
	opt := &option.PBOpt{}
	var msgs2 []proto.Message
	for i := 0; i < 1025; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = id
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", i)
		msgs2 = append(msgs2, data)
	}

	err = client.DoBatchGet(msgs2, opt)
	if !strings.Contains(err.Error(),"-4126") {
		t.Errorf("DoBatchGet fail, %s", err.Error())
		return
	}
}

//case1 insert success
// batch insert 多条记录相同的记录
func TestBatchInsertSuccess_03(t *testing.T) {
	client := tools.InitPBSyncClient()

	//batch insert 10 data
	var msgs []proto.Message
	//id := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		data := &tcaplusservice.GamePlayers{}
		data.PlayerId = 1
		data.PlayerName = "batchInsert"
		data.PlayerEmail = fmt.Sprintf("%d", 1)
		data.Pay = &tcaplusservice.Payment{Amount: uint64(i), PayId: 2, Method: 3}
		msgs = append(msgs, data)
	}

	err := client.DoBatchInsert(msgs, nil)
	if !strings.Contains(err.Error(),"-30") {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}
}

//case1 insert success
//batch insert 10 条记录，batch get前5条
func TestBatchInsertSuccess_04(t *testing.T) {
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
	// 2 batch Get
	opt := &option.PBOpt{}
	var msgs2 []proto.Message
	for i := 0; i < 5; i++ {
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
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i) {
			t.Errorf("DoBatchGet fail, %+v", msg)
			return
		}
	}
}

//case1 insert success
//batch insert 10 条记录，batch get后5条
func TestBatchInsertSuccess_05(t *testing.T) {
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
	// 2 batch Get
	opt := &option.PBOpt{}
	var msgs2 []proto.Message
	for i := 5; i < 10; i++ {
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
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i+5) {
			t.Errorf("DoBatchGet fail, %+v", msg)
			return
		}
	}
}

//case1 insert success
//batch insert 10 条记录，batch get中间的5条记录
func TestBatchInsertSuccess_06(t *testing.T) {
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
	// 2 batch Get
	opt := &option.PBOpt{}
	var msgs2 []proto.Message
	for i := 2; i < 7; i++ {
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
		if msg.(*tcaplusservice.GamePlayers).Pay.Amount != uint64(i+2) {
			t.Errorf("DoBatchGet fail, %+v", msg)
			return
		}
	}
}