package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"strconv"
	"sync"
	"testing"
	"time"
)

func syncInsert(client *tcaplus.Client, req request.TcaplusRequest, t *testing.T, oldJson string) {
	resp, err := client.Do(req, time.Duration(2*time.Second))
	if err != nil {
		t.Errorf("Do fail, %s", err.Error())
		wg.Done()
		return
	}
	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		wg.Done()
		return
	}
	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			wg.Done()
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			wg.Done()
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			wg.Done()
			return
		}
	}
	wg.Done()
}

func syncGet(client *tcaplus.Client, req request.TcaplusRequest, t *testing.T, oldJson string) {
	resp, err := client.Do(req, time.Duration(2*time.Second))
	if err != nil {
		wg.Done()
		t.Errorf("Do fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		wg.Done()
		t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
		return
	}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			wg.Done()
			return
		}

		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			wg.Done()
			return
		}

		newJson := StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData2 != reqData")
			wg.Done()
			return
		}
	}
	wg.Done()
}

func asyncInsert(client *tcaplus.Client, req request.TcaplusRequest, t *testing.T, oldJson string) {
	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		wg.Done()
		return
	}

	//recv resp
	resp, err := recvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		wg.Done()
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		wg.Done()
		return
	}
	wg.Done()
}

//case1 同步发送 多协程
var wg sync.WaitGroup

func TestSyncConcurrencyInsert(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	var uinKeyList [100]uint64
	var oldJsonList [100]string
	for i := 0; i < 100; i++ {
		wg.Add(1)
		req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
		if err != nil {
			t.Errorf("NewRequest fail, %s", err.Error())
			return
		}

		if err := req.SetResultFlag(2); err != nil {
			t.Errorf("SetResultFlag failed %v", err.Error())
			return
		}

		uinKey := time.Now().UnixNano()
		data := newGenericTableRec()
		data.Name = strconv.Itoa(i)
		data.Uin = uint64(uinKey)
		oldJson := StToJson(data)
		fmt.Println(oldJson)

		uinKeyList[i] = uint64(uinKey)
		oldJsonList[i] = oldJson
		rec, err := req.AddRecord(0)
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}
		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		go syncInsert(client, req, t, oldJson)
	}

	wg.Wait()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiGetReq)
		if err != nil {
			t.Errorf("NewRequest fail, %s", err.Error())
			return
		}

		data2 := tcaplus_tb.NewTable_Generic()
		data2.Name = strconv.Itoa(i)
		data2.Key3 = "key3"
		data2.Key4 = "key4"
		data2.Uin = uinKeyList[i]
		rec, err := req.AddRecord(0)
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}

		if err := rec.SetData(data2); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		go syncGet(client, req, t, oldJsonList[i])
	}
	wg.Wait()
	fmt.Println("sync request over")
	return
}

//case2 同步加异步并发发送
func TestSyncConcurrencyInsert2(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	var uinKeyList [100]uint64
	var oldJsonList [100]string
	for i := 0; i < 100; i++ {
		wg.Add(1)
		req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiInsertReq)
		if err != nil {
			t.Errorf("NewRequest fail, %s", err.Error())
			return
		}

		if err := req.SetResultFlag(2); err != nil {
			t.Errorf("SetResultFlag failed %v", err.Error())
			return
		}

		uinKey := time.Now().UnixNano()
		data := newGenericTableRec()
		data.Name = strconv.Itoa(i)
		data.Uin = uint64(uinKey)
		oldJson := StToJson(data)
		fmt.Println(oldJson)

		uinKeyList[i] = uint64(uinKey)
		oldJsonList[i] = oldJson
		rec, err := req.AddRecord(0)
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}
		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		if i%2 == 0 {
			go asyncInsert(client, req, t, oldJson)
		} else {
			go syncInsert(client, req, t, oldJson)
		}

	}

	wg.Wait()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		req, err := client.NewRequest(cfg.ApiConfig.ZoneId, TestTableName, cmd.TcaplusApiGetReq)
		if err != nil {
			t.Errorf("NewRequest fail, %s", err.Error())
			return
		}

		data2 := tcaplus_tb.NewTable_Generic()
		data2.Name = strconv.Itoa(i)
		data2.Key3 = "key3"
		data2.Key4 = "key4"
		data2.Uin = uinKeyList[i]
		rec, err := req.AddRecord(0)
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}

		if err := rec.SetData(data2); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		go syncGet(client, req, t, oldJsonList[i])
	}
	wg.Wait()
	fmt.Println("sync request over")
	return
}
