package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/policy"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"testing"
	"time"
)

//case 1记录不存在时Get fail
func TestGetBypartKeyFail(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiGetByPartkeyReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	data := newGenericTableRec()
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	if err := rec.SetDataWithIndexAndField(data, nil, "Index3"); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		t.Errorf("resp.GetResult expect err ,but nil, so test may be affected.")
		return
	} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult expect TXHDB_ERR_RECORD_NOT_EXIST ,but %s", terror.GetErrMsg(err))
		return
	} else {
		fmt.Printf("TestGetBypartKeyFail test pass\n")
		//fmt.Printf("return error info 0x%x, not exist flag: 0x%x\n", resp.GetResult(), terror.TXHDB_ERR_RECORD_NOT_EXIST)
	}
}

//case2 记录存在时，Get返回成功
func TestGetBypartkeySuccess(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiInsertReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}
	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	req.SetResultLimit(100, 0)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err,%d, %s", err, terror.GetErrMsg(err))
		return
	}

	if 1 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	////2 Get记录已存在
	client2, req2 := InitClientAndReq(cmd.TcaplusApiGetByPartkeyReq)
	if nil == client || nil == req {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}
	//相同的key
	data2 := tcaplus_tb.NewTable_Generic()
	//data2.Name = "GoUnitTest"
	//data2.Key3 = "key3"
	//data2.Key4 = "key4"
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if err := rec2.SetDataWithIndexAndField(data2, nil, "Index1"); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	//recv resp
	resp2, err := AsyncSendAndGetRes(client2, req2)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		if newJson != oldJson {
			t.Errorf("resData2 != reqData")
			fmt.Println(newJson)
			fmt.Println(oldJson)
			fmt.Println(oldJson2)
			return
		}
		if record.GetVersion() <= 0 {
			t.Errorf("record.GetVersion %d <=0 ", record.GetVersion())
		}
	}
}

//case3 index的key不存在，key不存在的时候get失败
func TestGetByPartKey_Key_NonExist(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiGetByPartkeyReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	data := newGenericTableRec()
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	if err := rec.SetDataWithIndexAndField(data, nil, "Index5"); err == nil {
		return
	}

}

//case4 partkey delete的版本号设置
func TestGetByPartKey_CHECKDATAVERSION_AUTOINCREASE(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiInsertReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}
	//set policy CHECKDATAVERSION_AUTOINCREASE
	req.SetVersionPolicy(policy.CheckDataVersionAutoIncrease)

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	oldJson := StToJson(data)
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	rec.SetVersion(9)
	req.SetResultLimit(100, 0)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err,%d, %s", err, terror.GetErrMsg(err))
		return
	}

	if 1 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}
	//Get记录已存在
	client2, req2 := InitClientAndReq(cmd.TcaplusApiDeleteByPartkeyReq)
	if nil == client2 || nil == req2 {
		t.Errorf("NewRequest fail")
		return
	}
	//相同的key
	data2 := tcaplus_tb.NewTable_Generic()
	//data2.Name = "GoUnitTest"
	//data2.Key3 = "key3"
	//data2.Key4 = "key4"
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	oldJson2 := StToJson(data2)
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	rec2.SetVersion(10)
	if err := rec2.SetDataWithIndexAndField(data2, nil, "Index1"); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	//recv resp
	resp2, err := AsyncSendAndGetRes(client2, req2)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err == 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	} else if err != terror.SVR_ERR_FAIL_INVALID_VERSION {
		t.Errorf("resp.GetResult expect SVR_ERR_FAIL_INVALID_VERSION, but %s", terror.GetErrMsg(err))
		return
	}

	for i := 0; i < resp2.GetRecordCount(); i++ {
		record, err := resp2.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := StToJson(newData)
		if newJson != oldJson {
			t.Errorf("resData2 != reqData")
			fmt.Println(newJson)
			fmt.Println(oldJson)
			fmt.Println(oldJson2)
			return
		}
	}
}

//case6 NOCHECKDATAVERSION_AUTOINCREASE
func TestGetByPartKey_NOCHECKDATAVERSION_AUTOINCREASE(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiInsertReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}
	//set policy NOCHECKDATAVERSION_AUTOINCREASE
	req.SetVersionPolicy(policy.NoCheckDataVersionAutoIncrease)

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	//oldJson := StToJson(data)
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	req.SetResultLimit(100, 0)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err,%d, %s", err, terror.GetErrMsg(err))
		return
	}

	if 1 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}
	//Get记录已存在
	client2, req2 := InitClientAndReq(cmd.TcaplusApiDeleteByPartkeyReq)
	if nil == client2 || nil == req2 {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}
	req2.SetVersionPolicy(policy.NoCheckDataVersionAutoIncrease)
	//相同的key
	data2 := tcaplus_tb.NewTable_Generic()
	//data2.Name = "GoUnitTest"
	//data2.Key3 = "key3"
	//data2.Key4 = "key4"
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	//oldJson2 := StToJson(data2)
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	rec2.SetVersion(111)
	if err := rec2.SetDataWithIndexAndField(data2, nil, "Index1"); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	//recv resp
	resp2, err := AsyncSendAndGetRes(client2, req2)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}
	// get到的数据和删除的数据如果不一致，则认为失败
	if resp2.GetRecordCount() != resp.GetRecordCount() {
		t.Errorf("get and delete by partkey return different recornum %s", err.Error())
		return
	}

}

//case5 NOCHECKDATAVERSION_OVERWRITE
func TestGetByPartKey_NOCHECKDATAVERSION_OVERWRITE(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiInsertReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	if err := req.SetResultFlag(2); err != nil {
		t.Errorf("SetResultFlag failed %v", err.Error())
		return
	}
	//set policy NOCHECKDATAVERSION_OVERWRITE
	req.SetVersionPolicy(policy.NoCheckDataVersionOverwrite)

	uinKey := time.Now().UnixNano()
	data := newGenericTableRec()
	data.Uin = uint64(uinKey)

	//oldJson := StToJson(data)
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	req.SetResultLimit(100, 0)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err,%d, %s", err, terror.GetErrMsg(err))
		return
	}

	if 1 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}
	//Get记录已存在
	client2, req2 := InitClientAndReq(cmd.TcaplusApiDeleteByPartkeyReq)
	if nil == client2 || nil == req2 {
		t.Errorf("NewRequest fail, %s", err.Error())
		return
	}
	req2.SetVersionPolicy(policy.NoCheckDataVersionAutoIncrease)
	//相同的key
	data2 := tcaplus_tb.NewTable_Generic()
	//data2.Name = "GoUnitTest"
	//data2.Key3 = "key3"
	//data2.Key4 = "key4"
	data2.Uin = uint64(uinKey)
	//不同的value
	data2.Level = 222

	//oldJson2 := StToJson(data2)
	rec2, err := req2.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	rec2.SetVersion(111)
	if err := rec2.SetDataWithIndexAndField(data2, nil, "Index1"); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	//recv resp
	resp2, err := AsyncSendAndGetRes(client2, req2)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp2.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}
	// get到的数据和删除的数据如果不一致，则认为失败
	if resp2.GetRecordCount() != resp.GetRecordCount() {
		t.Errorf("get and delete by partkey return different recornum %s", err.Error())
		return
	}

}

//case6 测试insert, getBypartkey, deletebyPartkey, and then getbypartkey again
func TestGetBypartKeyMany(t *testing.T) {
	// 插入数据
	for idx := int(0); idx < 20; idx++ {
		InsertKV(0x55aa7788, "key4")
	}

	// 测试get接口
	client, req := InitClientAndReq(cmd.TcaplusApiGetByPartkeyReq)
	if nil == client || nil == req {
		t.Errorf("InitClientAndReq fail")
		return
	}
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	data := tcaplus_tb.NewTable_Generic()
	data.Key3 = "key3"
	data.Key4 = "key4"
	data.Uin = 0x55aa7788

	if err := rec.SetDataWithIndexAndField(data, nil, "Index4"); err != nil {
		return
	}
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult expect error, but nil")
		return
	}
	// 测试删除接口
	client2, req2 := InitClientAndReq(cmd.TcaplusApiDeleteByPartkeyReq)
	if nil == client2 || nil == req2 {
		t.Errorf("InitClientAndReq fail")
		return
	}
	rec2, err2 := req2.AddRecord(0)
	if err2 != nil {
		t.Errorf("AddRecord fail, %s", err2.Error())
		return
	}
	data2 := tcaplus_tb.NewTable_Generic()
	data2.Key3 = "key3"
	data2.Key4 = "key4"
	data2.Uin = 0x55aa7788

	if err := rec2.SetDataWithIndexAndField(data, nil, "Index4"); err != nil {
		return
	}
	res2, err := AsyncSendAndGetRes(client2, req2)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}
	// get到的数据和删除的数据如果不一致，则认为失败
	if res2.GetRecordCount() != resp.GetRecordCount() {
		t.Errorf("get and delete by partkey return different recornum %d-%d", res2.GetRecordCount(), resp.GetRecordCount())
		return
	}
	{
		client3, req3 := InitClientAndReq(cmd.TcaplusApiGetByPartkeyReq)
		if nil == client3 || nil == req3 {
			t.Errorf("InitClientAndReq fail")
			return
		}
		rec3, err := req3.AddRecord(0)
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}
		data3 := tcaplus_tb.NewTable_Generic()
		data3.Key3 = "key3"
		data3.Key4 = "key4"
		data3.Uin = 0x55aa7788

		if err := rec3.SetDataWithIndexAndField(data3, nil, "Index4"); err != nil {
			return
		}
		resp3, err := AsyncSendAndGetRes(client3, req3)
		if err != nil {
			t.Errorf("recvResponse fail, %s", err.Error())
			return
		}
		if err := resp3.GetResult(); err == 0 {
			t.Errorf("resp.GetResult expect err ,but nil, so test may be affected.")
			return
		} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
			t.Errorf("resp.GetResult expect TXHDB_ERR_RECORD_NOT_EXIST ,but %s", terror.GetErrMsg(err))
			return
		} else {
			//fmt.Printf("TestGetBypartKeyMany test pass")
			fmt.Printf("TestGetBypartKeyMany result 0x%x, expected 0x%x, result right\n", resp3.GetResult(), terror.TXHDB_ERR_RECORD_NOT_EXIST)
		}

	}
}

func TestGetBypartKeyManySync(t *testing.T) {
	// 插入数据
	for idx := int(0); idx < 3800; idx++ {
		//fmt.Printf("current idx:%d\n", idx)
		InsertKV(0x55aa7788, "key4")
	}

	// 测试get接口
	client, req := InitClientAndReq(cmd.TcaplusApiGetByPartkeyReq)
	if nil == client || nil == req {
		t.Errorf("InitClientAndReq fail")
		return
	}
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	data := tcaplus_tb.NewTable_Generic()
	data.Key3 = "key3"
	data.Key4 = "key4"
	data.Uin = 0x55aa7788

	if err := rec.SetDataWithIndexAndField(data, nil, "Index4"); err != nil {
		return
	}
	respMore, err := client.DoMore(req, time.Duration(4*time.Second)) //(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}
	fmt.Printf("has more, pkg num:%d\n", len(respMore))
	var totalCnt int = 0
	for _, resp := range respMore {
		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult expect error, but nil")
			return
		}
		fmt.Printf("has more, current pkg rec num:%d\n", resp.GetRecordCount())
		totalCnt += resp.GetRecordCount()
	}

	// 测试删除接口
	client2, req2 := InitClientAndReq(cmd.TcaplusApiDeleteByPartkeyReq)
	if nil == client2 || nil == req2 {
		t.Errorf("InitClientAndReq fail")
		return
	}
	rec2, err2 := req2.AddRecord(0)
	if err2 != nil {
		t.Errorf("AddRecord fail, %s", err2.Error())
		return
	}
	data2 := tcaplus_tb.NewTable_Generic()
	data2.Key3 = "key3"
	data2.Key4 = "key4"
	data2.Uin = 0x55aa7788

	if err := rec2.SetDataWithIndexAndField(data, nil, "Index4"); err != nil {
		return
	}
	resMore2, err := client2.DoMore(req2, time.Duration(10*time.Second))
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}
	fmt.Printf("has more, pkg num:%d\n", len(resMore2))
	// get到的数据和删除的数据如果不一致，则认为失败
	var totalCnt2 int = 0
	for _, res2 := range resMore2 {
		if err := res2.GetResult(); err != 0 {
			t.Errorf("error : %d", err)
			return
		}
		totalCnt2 += res2.GetRecordCount()
		fmt.Printf("has more, current pkg rec num:%d\n", res2.GetRecordCount())
	}
	if totalCnt != totalCnt2 {
		t.Errorf("get record num: %d, delete record num:%d", totalCnt, totalCnt2)
		return
	} else {
		fmt.Printf("get record num: %d, delete record num:%d", totalCnt, totalCnt2)
	}
	{
		client3, req3 := InitClientAndReq(cmd.TcaplusApiGetByPartkeyReq)
		if nil == client3 || nil == req3 {
			t.Errorf("InitClientAndReq fail")
			return
		}
		rec3, err := req3.AddRecord(0)
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}
		data3 := tcaplus_tb.NewTable_Generic()
		data3.Key3 = "key3"
		data3.Key4 = "key4"
		data3.Uin = 0x55aa7788

		if err := rec3.SetDataWithIndexAndField(data3, nil, "Index4"); err != nil {
			return
		}
		respMore3, err := client3.DoMore(req3, time.Duration(2*time.Second))
		if err != nil {
			t.Errorf("recvResponse fail, %s", err.Error())
			return
		}
		for _, res3 := range respMore3 {
			if err := res3.GetResult(); err == 0 {
				t.Errorf("resp.GetResult expect err ,but nil, so test may be affected.")
				return
			} else if err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
				t.Errorf("resp.GetResult expect TXHDB_ERR_RECORD_NOT_EXIST ,but %s", terror.GetErrMsg(err))
				return
			} else {
				fmt.Printf("TestGetBypartKeyManySync result 0x%x, expected 0x%x, result right\n", res3.GetResult(), terror.TXHDB_ERR_RECORD_NOT_EXIST)
			}
		}
	}
}
