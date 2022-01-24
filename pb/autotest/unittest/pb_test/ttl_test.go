package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
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
// case2 list ttl set,list table not support ttl
//func TestListTTLSuccess(t *testing.T) {
//	client,err := tools.InitClient()
//	if err != nil {
//		t.Errorf("InitClient failed %s", err.Error())
//		return
//	}
//	tableName := "table_traverser_list"
//	//add 3 record
//	data := tcaplus_tb.NewTable_Traverser_List()
//	data.Key = 1
//	data.Name = 255
//	data.Level = 1
//	data.Value1 = "value1"
//	data.Value2 = "value2"
//	err = client.DoListAddAfter(tableName, data, -1, nil)
//	if err != nil {
//		t.Errorf("DoListAddAfter failed %s", err.Error())
//		return
//	}
//	err = client.DoListAddAfter(tableName, data, -1, nil)
//	if err != nil {
//		t.Errorf("DoListAddAfter failed %s", err.Error())
//		return
//	}
//	err = client.DoListAddAfter(tableName, data, -1, nil)
//	if err != nil {
//		t.Errorf("DoListAddAfter failed %s", err.Error())
//		return
//	}
//
//	// set ttl
//	opt := &option.TDROpt{
//		BatchTTL: []option.TTLInfo{{TTL:5000},{TTL:4000},{TTL:3000}},
//	}
//	dataSlice := []record.TdrTableSt{data,data,data}
//	indexs := []int32{0,1,2}
//	err = client.DoSetTTLBatch(tableName, dataSlice,indexs, opt)
//	if err != nil {
//		t.Errorf("DoSetTTLBatch fail, %s", err.Error())
//		return
//	}
//
//	//3 get ttl
//	opt = &option.TDROpt{}
//	err = client.DoGetTTLBatch(tableName, dataSlice, indexs, opt)
//	if err != nil {
//		t.Errorf("DoGetTTLBatch fail, %s", err.Error())
//		return
//	}
//	fmt.Println(opt.BatchTTL)
//	for _, ttl := range opt.BatchTTL {
//		if ttl.TTL <=0 || ttl.TTL > 5000{
//			t.Errorf("ttl invalid %d", ttl.TTL)
//			return
//		}
//	}
//
//	//4 5s after
//	time.Sleep(5*time.Second)
//	opt = &option.TDROpt{}
//	err = client.DoGetTTLBatch(tableName, dataSlice,indexs, opt)
//	fmt.Println(opt.BatchResult)
//	if err == nil {
//		t.Errorf("DoGetTTLMust timeout, %v", opt.BatchResult)
//		return
//	}
//}
