package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"strings"
	"testing"
	"time"
)

func InsertBatchTestKV(uin uint64, name string, key3 string, key4 string) {
	client, req := InitClientAndReq(cmd.TcaplusApiInsertReq)
	if nil == client || nil == req {
		fmt.Printf("NewRequest fail")
		return
	}

	//data
	data := newGenericTableRec()
	data.Uin = uin
	data.Name = name
	data.Key3 = key3
	data.Key4 = key4
	data.Simple_Struct.C_Int64 = 100
	rec, err := req.AddRecord(0)
	if err != nil {
		fmt.Printf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec.SetData(data); err != nil {
		fmt.Printf("SetData fail, %s", err.Error())
		return
	}

	if _, err := AsyncSendAndGetRes(client, req); err != nil {
		fmt.Printf("recvResponse fail, %s", err.Error())
		return
	}
}

func TestBatchGetNoRecord(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiBatchGetReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	req.SetResultFlag(2)
	nameList := []string{string("level"), "count"}
	req.SetFieldNames(nameList)
	start_uin := uint64(time.Now().UnixNano())
	key4 := randomStrChar("key4", int(start_uin)%8)
	for idx := 0; idx < 10; idx++ {
		rec, err := req.AddRecord(0)
		if rec == nil || err != nil {
			fmt.Printf("sfsdf")
		}
		//req.SetFieldNames([]string{"level", "count"})
		data := newGenericTableRec()
		data.Uin = start_uin + uint64(idx)
		data.Name = "TestBatchGet"
		data.Key3 = "key3"
		data.Key4 = key4
		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
	}
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail")
		return
	}

	for idx := int(0); idx < resp.GetRecordCount(); idx++ {
		rec, err := resp.FetchRecord()
		if err != nil {
			errNotExist := terror.ErrorCode{Code: terror.TXHDB_ERR_RECORD_NOT_EXIST}
			if 0 != strings.Compare(errNotExist.Error(), err.Error()) {
				t.Errorf("expect: %s but return: %s", errNotExist.Error(), err.Error())
				return
			}
		} else {
			t.Errorf("expect: no record, but return %+v", rec)
			return
		}
	}
}

func TestBatchGetRecord(t *testing.T) {
	start_uin := uint64(time.Now().UnixNano())
	key4 := randomStrChar("key4", int(start_uin)%8)
	for idx := 0; idx < 10; idx++ {

		data := newGenericTableRec()
		data.Uin = start_uin + uint64(idx)
		data.Name = "TestBatchGet"
		data.Key3 = "key3"
		data.Key4 = key4
		InsertBatchTestKV(data.Uin, data.Name, data.Key3, data.Key4)
	}

	client, req := InitClientAndReq(cmd.TcaplusApiBatchGetReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	req.SetResultFlag(2)
	nameList := []string{string("level"), "count"}
	req.SetFieldNames(nameList)

	for idx := 0; idx < 10; idx++ {
		rec, err := req.AddRecord(0)
		if rec == nil || err != nil {
			fmt.Printf("sfsdf")
		}
		//req.SetFieldNames([]string{"level", "count"})
		data := newGenericTableRec()
		data.Uin = start_uin + uint64(idx)
		data.Name = "TestBatchGet"
		data.Key3 = "key3"
		data.Key4 = key4
		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		InsertBatchTestKV(data.Uin, data.Name, data.Key3, data.Key4)

	}
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail")
		return
	}

	for idx := int(0); idx < resp.GetRecordCount(); idx++ {
		_, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("expect ok but return: %s", err.Error())
		}
		return
	}
}

func TestBatchGetRecordDoMore(t *testing.T) {
	start_uin := uint64(time.Now().UnixNano())
	key4 := randomStrChar("key4", int(start_uin)%8)
	for idx := 0; idx < 10; idx++ {

		data := newGenericTableRec()
		data.Uin = start_uin + uint64(idx)
		data.Name = "TestBatchGet"
		data.Key3 = "key3"
		data.Key4 = key4
		InsertBatchTestKV(data.Uin, data.Name, data.Key3, data.Key4)
	}

	client, req := InitClientAndReq(cmd.TcaplusApiBatchGetReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	req.SetResultFlag(2)
	nameList := []string{string("level"), "count"}
	req.SetFieldNames(nameList)

	for idx := 0; idx < 10; idx++ {
		rec, err := req.AddRecord(0)
		if rec == nil || err != nil {
			fmt.Printf("sfsdf")
		}
		//req.SetFieldNames([]string{"level", "count"})
		data := newGenericTableRec()
		data.Uin = start_uin + uint64(idx)
		data.Name = "TestBatchGet"
		data.Key3 = "key3"
		data.Key4 = key4
		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		InsertBatchTestKV(data.Uin, data.Name, data.Key3, data.Key4)

	}

	resp_arr, err := client.DoMore(req, time.Duration(10*time.Second))
	if err != nil {
		t.Errorf("recvResponse fail")
		return
	}
	fmt.Printf("has more, pkg num:%d\n", len(resp_arr))
	for _, res2 := range resp_arr {
		if err := res2.GetResult(); err != 0 {
			t.Errorf("error : %d", err)
			return
		}
		fmt.Printf("has more, current pkg rec num:%d\n", res2.GetRecordCount())
		for idx := int(0); idx < res2.GetRecordCount(); idx++ {
			_, err := res2.FetchRecord()
			if err != nil {
				t.Errorf("expect ok but return: %s", err.Error())
				return
			}
		}
	}

}

//����1024����¼Ȼ��batch get
func TestBatchGetRecord_1024(t *testing.T) {
	start_uin := uint64(time.Now().UnixNano())
	key4 := randomStrChar("key4", int(start_uin)%8)
	for idx := 0; idx < 1024; idx++ {

		data := newGenericTableRec()
		data.Uin = start_uin + uint64(idx)
		data.Name = "TestBatchGet"
		data.Key3 = "key3"
		data.Key4 = key4
		InsertBatchTestKV(data.Uin, data.Name, data.Key3, data.Key4)
	}

	client, req := InitClientAndReq(cmd.TcaplusApiBatchGetReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	req.SetResultFlag(2)
	nameList := []string{string("level"), "count"}
	req.SetFieldNames(nameList)

	for idx := 0; idx < 1024; idx++ {
		rec, err := req.AddRecord(0)
		if rec == nil || err != nil {
			fmt.Printf("sfsdf")
		}
		//req.SetFieldNames([]string{"level", "count"})
		data := newGenericTableRec()
		data.Uin = start_uin + uint64(idx)
		data.Name = "TestBatchGet"
		data.Key3 = "key3"
		data.Key4 = key4
		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		InsertBatchTestKV(data.Uin, data.Name, data.Key3, data.Key4)

	}
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail")
		return
	}

	for idx := int(0); idx < resp.GetRecordCount(); idx++ {
		_, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("expect ok but return: %s", err.Error())
		}
		return
	}
}
