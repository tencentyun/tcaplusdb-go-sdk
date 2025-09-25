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

//case1 set ttl，Get返回成功
func TestGetTTLSuccess(t *testing.T) {
	client := tools.InitPBSyncClient()

	// insert 3 data
	data1 := &tcaplusservice.GamePlayers{}
	data1.PlayerId = 1
	data1.PlayerName = "jiahua"
	data1.PlayerEmail = "dsf"
	client.DoDelete(data1, nil)

	data2 := &tcaplusservice.GamePlayers{}
	data2.PlayerId = 2
	data2.PlayerName = "jiahua"
	data2.PlayerEmail = "dsf"
	client.DoDelete(data2, nil)

	data3 := &tcaplusservice.GamePlayers{}
	data3.PlayerId = 3
	data3.PlayerName = "jiahua"
	data3.PlayerEmail = "dsf"
	client.DoDelete(data3, nil)

	data1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret := client.DoInsert(data1, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	data2.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret = client.DoInsert(data2, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	data3.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret = client.DoInsert(data3, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	////2 set ttl
	opt := &option.PBOpt{
		BatchTTL: []option.TTLInfo{{TTL: 5000}, {TTL: 5000}, {TTL: 5000}},
	}
	msgs := []proto.Message{data1, data2, data3}
	err := client.DoSetTTLBatch(msgs, nil, opt)
	if err != nil {
		t.Errorf("DoSetTTLBatch fail, %s", err.Error())
		return
	}

	//3 get ttl
	opt = &option.PBOpt{}
	err = client.DoGetTTLBatch(msgs, nil, opt)
	if err != nil {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}
	fmt.Println(opt.BatchTTL)
	for _, ttl := range opt.BatchTTL {
		if ttl.TTL <= 0 || ttl.TTL > 5000 {
			t.Errorf("ttl invalid %d", ttl.TTL)
			return
		}
	}

	//4 5s after
	time.Sleep(5 * time.Second)
	opt = &option.PBOpt{}
	err = client.DoGetTTLBatch(msgs, nil, opt)
	fmt.Println(opt.BatchResult)
	if err == nil {
		t.Errorf("DoGetTTLMust timeout, %v", opt.BatchResult)
		return
	}
}

//case1 set ttl，Get返回成功
func TestGetTTLSuccess_ttl_0(t *testing.T) {
	client := tools.InitPBSyncClient()

	// insert 3 data
	data1 := &tcaplusservice.GamePlayers{}
	data1.PlayerId = 1
	data1.PlayerName = "jiahua"
	data1.PlayerEmail = "dsf"
	client.DoDelete(data1, nil)

	data2 := &tcaplusservice.GamePlayers{}
	data2.PlayerId = 2
	data2.PlayerName = "jiahua"
	data2.PlayerEmail = "dsf"
	client.DoDelete(data2, nil)

	data3 := &tcaplusservice.GamePlayers{}
	data3.PlayerId = 3
	data3.PlayerName = "jiahua"
	data3.PlayerEmail = "dsf"
	client.DoDelete(data3, nil)

	data1.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret := client.DoInsert(data1, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	data2.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret = client.DoInsert(data2, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	data3.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	ret = client.DoInsert(data3, nil)
	if ret != nil {
		t.Errorf("DoInsert failed %d", ret)
		return
	}

	////2 set ttl
	opt := &option.PBOpt{
		//BatchTTL: []option.TTLInfo{{TTL: 5000}, {TTL: 5000, {TTL: 5000}},
		BatchTTL: []option.TTLInfo{{TTL: 0}, {TTL: 0}, {TTL: 0}},
	}
	msgs := []proto.Message{data1, data2, data3}
	err := client.DoSetTTLBatch(msgs, nil, opt)
	if err != nil {
		t.Errorf("DoSetTTLBatch fail, %s", err.Error())
		return
	}

	//3 get ttl
	opt = &option.PBOpt{}
	err = client.DoGetTTLBatch(msgs, nil, opt)
	if err.Error() != "errCode: 2309, errMsg: " {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}

}

// case insert update replace同时支持ttl
func TestWriteWithTTLSuccess(t *testing.T) {
	client := tools.InitPBSyncClient()

	//1 insert with ttl
	data := &tcaplusservice.GamePlayers{}
	data.PlayerId = time.Now().UnixNano()
	data.PlayerName = "test"
	data.PlayerEmail = "dsf"
	data.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2, Method: 3}
	opt := &option.PBOpt{
		TTL: &option.TTLInfo{
			TTL: 5000,
		},
	}
	err := client.DoInsert(data, opt)
	if err != nil {
		t.Errorf("DoInsert failed %v", err)
		return
	}

	//1 get ttl
	msgs := []proto.Message{data}
	opt = &option.PBOpt{}
	err = client.DoGetTTLBatch(msgs, nil, opt)
	if err != nil {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}
	if len(opt.BatchTTL) == 0 {
		t.Errorf("opt.BatchTTL empty")
		return
	}
	if opt.BatchTTL[0].TTL < 0 || opt.BatchTTL[0].TTL > 5000 {
		t.Errorf("opt.BatchTTL %d invalid", opt.BatchTTL[0].TTL)
		return
	}
	fmt.Println("after insert ttl", opt.BatchTTL[0].TTL)

	//2 DoUpdate with ttl
	opt = &option.PBOpt{
		TTL: &option.TTLInfo{
			TTL: 10000,
		},
	}
	err = client.DoUpdate(data, opt)
	if err != nil {
		t.Errorf("DoUpdate failed %v", err)
		return
	}

	//2 get ttl
	opt = &option.PBOpt{}
	err = client.DoGetTTLBatch(msgs, nil, opt)
	if err != nil {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}
	if len(opt.BatchTTL) == 0 {
		t.Errorf("opt.BatchTTL empty")
		return
	}
	if opt.BatchTTL[0].TTL < 5 || opt.BatchTTL[0].TTL > 10000 {
		t.Errorf("opt.BatchTTL %d invalid", opt.BatchTTL[0].TTL)
		return
	}
	fmt.Println("after update ttl", opt.BatchTTL[0].TTL)

	//2 DoReplace with ttl
	opt = &option.PBOpt{
		TTL: &option.TTLInfo{
			TTL: 15000,
		},
	}
	err = client.DoReplace(data, opt)
	if err != nil {
		t.Errorf("DoReplace failed %v", err)
		return
	}

	//2 get ttl
	opt = &option.PBOpt{}
	err = client.DoGetTTLBatch(msgs, nil, opt)
	if err != nil {
		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
		return
	}
	if len(opt.BatchTTL) == 0 {
		t.Errorf("opt.BatchTTL empty")
		return
	}
	if opt.BatchTTL[0].TTL < 10 || opt.BatchTTL[0].TTL > 15000 {
		t.Errorf("opt.BatchTTL %d invalid", opt.BatchTTL[0].TTL)
		return
	}
	fmt.Println("after replace ttl", opt.BatchTTL[0].TTL)
}
