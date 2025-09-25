package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"testing"
	"time"
)

// insert / update /replace 不设置resultflag返回version

// insert返回version
func TestPBInsertVersion(t *testing.T) {
	client := tools.InitPBSyncClient()
	data := &tcaplusservice.GamePlayers{}
	data.PlayerId = time.Now().UnixNano()
	data.PlayerName = "test"
	data.PlayerEmail = "dsf"
	data.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	opt := &option.PBOpt{}
	ret := client.DoInsert(data, opt)
	if ret != nil {
		t.Errorf("DoInsert failed %v", ret)
		return
	}
	if opt.Version != 1 {
		t.Errorf("DoInsert version %v != 1", opt.Version)
		return
	}
	fmt.Println("insert version", opt.Version)
}

// Update返回version
func TestPBUpdateVersion(t *testing.T) {
	client := tools.InitPBSyncClient()
	data := &tcaplusservice.GamePlayers{}
	data.PlayerId = time.Now().UnixNano()
	data.PlayerName = "test"
	data.PlayerEmail = "dsf"
	data.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	opt := &option.PBOpt{}
	ret := client.DoInsert(data, opt)
	if ret != nil {
		t.Errorf("DoInsert failed %v", ret)
		return
	}
	if opt.Version != 1 {
		t.Errorf("DoInsert version %v != 1", opt.Version)
		return
	}
	fmt.Println("insert version", opt.Version)

	ret = client.DoUpdate(data, opt)
	if ret != nil {
		t.Errorf("DoInsert failed %v", ret)
		return
	}
	if opt.Version != 2 {
		t.Errorf("DoInsert version %v != 2", opt.Version)
		return
	}
	fmt.Println("update version", opt.Version)
}

// Update返回version
func TestPBReplaceVersion(t *testing.T) {
	client := tools.InitPBSyncClient()
	data := &tcaplusservice.GamePlayers{}
	data.PlayerId = time.Now().UnixNano()
	data.PlayerName = "test"
	data.PlayerEmail = "dsf"
	data.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	opt := &option.PBOpt{}
	ret := client.DoReplace(data, opt)
	if ret != nil {
		t.Errorf("DoReplace failed %v", ret)
		return
	}
	if opt.Version != 1 {
		t.Errorf("DoReplace version %v != 1", opt.Version)
		return
	}
	fmt.Println("1 DoReplace version", opt.Version)

	ret = client.DoReplace(data, opt)
	if ret != nil {
		t.Errorf("DoReplace failed %v", ret)
		return
	}
	if opt.Version != 2 {
		t.Errorf("DoReplace version %v != 2", opt.Version)
		return
	}
	fmt.Println("2 DoReplace version", opt.Version)
}

// FieldIncrease返回version
func TestPBFieldIncreaseVersion(t *testing.T) {
	client := tools.InitPBSyncClient()
	data := &tcaplusservice.GamePlayers{}
	data.PlayerId = time.Now().UnixNano()
	data.PlayerName = "test"
	data.PlayerEmail = "dsf"
	data.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	opt := &option.PBOpt{}
	ret := client.DoReplace(data, opt)
	if ret != nil {
		t.Errorf("DoReplace failed %v", ret)
		return
	}
	if opt.Version != 1 {
		t.Errorf("DoReplace version %v != 1", opt.Version)
		return
	}
	fmt.Println("1 DoReplace version", opt.Version)

	opt = &option.PBOpt{
		FieldNames: []string{"pay.pay_id"},
	}
	ret = client.DoFieldIncrease(data, opt)
	if ret != nil {
		t.Errorf("DoFieldIncrease failed %v", ret)
		return
	}
	if opt.Version != 2 {
		t.Errorf("DoFieldIncrease version %v != 2", opt.Version)
		return
	}
	fmt.Println("2 DoFieldIncrease version", opt.Version)
}

// FieldUpdate返回version
func TestPBFieldUpdateVersion(t *testing.T) {
	client := tools.InitPBSyncClient()
	data := &tcaplusservice.GamePlayers{}
	data.PlayerId = time.Now().UnixNano()
	data.PlayerName = "test"
	data.PlayerEmail = "dsf"
	data.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	opt := &option.PBOpt{}
	ret := client.DoReplace(data, opt)
	if ret != nil {
		t.Errorf("DoReplace failed %v", ret)
		return
	}
	if opt.Version != 1 {
		t.Errorf("DoReplace version %v != 1", opt.Version)
		return
	}
	fmt.Println("1 DoReplace version", opt.Version)

	opt = &option.PBOpt{
		FieldNames: []string{"pay.pay_id"},
	}
	ret = client.DoFieldUpdate(data, opt)
	if ret != nil {
		t.Errorf("DoFieldUpdate failed %v", ret)
		return
	}
	if opt.Version != 2 {
		t.Errorf("DoFieldUpdate version %v != 2", opt.Version)
		return
	}
	fmt.Println("2 DoFieldUpdate version", opt.Version)
}
