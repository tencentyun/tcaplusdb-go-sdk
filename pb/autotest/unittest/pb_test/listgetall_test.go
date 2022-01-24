package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"testing"
	"time"
)

//case2 记录存在时，Get返回成功
func TestListGetAllSuccess(t *testing.T) {
	data := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: "lol",
	}
	oldJson := tools.StToJson(data)
	fmt.Println(oldJson)

	/*------------------------------------------------------- 清理记录 ------------------------------------------*/
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListDeleteAllReq, "tb_online_list")

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 && err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	/*------------------------------------------------------- 插入记录 ------------------------------------------*/
	for i := 0; i < 5; i++ {
		client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListAddAfterReq, "tb_online_list")

		//add record
		rec, err = req.AddRecord(int32(tcaplus_protocol_cs.TCAPLUS_LIST_LAST_INDEX))
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}

		if _, err := rec.SetPBData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}

		if err := req.SetResultFlag(2); err != nil {
			t.Errorf("SetResultFlag fail, %s", err.Error())
			return
		}

		if err := client.SendRequest(req); err != nil {
			t.Errorf("SendRequest fail, %s", err.Error())
			return
		}

		//recv resp
		resp, err = tools.RecvResponse(client)
		if err != nil {
			t.Errorf("recvResponse fail, %s", err.Error())
			return
		}

		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
			return
		}

		if 1 != resp.GetRecordCount() {
			t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
			return
		}
	}

	/*------------------------------------------------------- 获取记录 ------------------------------------------*/
	client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListGetAllReq, "tb_online_list")

	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	req.SetResultLimit(100, 0)

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err = tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	fmt.Println(resp.GetRecordMatchCount())

	if 5 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := &tcaplusservice.TbOnlineList{
			Openid:    1,
			Tconndid:  2,
			Timekey:   "test",
			Gamesvrid: "lol",
		}
		if err := record.GetPBData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newData)
		fmt.Println(newJson)
		if newJson != oldJson {
			t.Errorf("resData != reqData")
			return
		}
	}
}

//记录不存在的时候，返回261
func TestGetListGetAllFail(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListGetAllReq, "tb_online_list")
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	//if err := req.SetResultFlag(2); err != nil {
	//      t.Errorf("SetResultFlag failed %v", err.Error())
	//      return
	//}
	data := &tcaplusservice.TbOnlineList{
		Openid:    1000,
		Tconndid:  2000,
		Timekey:   "notexist",
		Gamesvrid: "lol",
	}
	rec, err := req.AddRecord(0)

	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	req.SetResultLimit(100, 0)

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	//recv resp
	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	resp, err := tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 261 {
		t.Errorf("resp.GetResult err,%d, %s", err, terror.GetErrMsg(err))
		return
	}
}

func TestListGetAllMultiResp(t *testing.T) {
	data := &tcaplusservice.TbOnlineList{
		Openid:    1,
		Tconndid:  2,
		Timekey:   "test",
		Gamesvrid: tools.Bytes256K(),
	}
	//oldJson := tools.StToJson(data)
	//fmt.Println(oldJson)

	/*------------------------------------------------------- 清理记录 ------------------------------------------*/
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListDeleteAllReq, "tb_online_list")

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 && err != terror.TXHDB_ERR_RECORD_NOT_EXIST {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	/*------------------------------------------------------- 插入记录 ------------------------------------------*/
	for i := 0; i < 30; i++ {
		client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListAddAfterReq, "tb_online_list")

		//add record
		rec, err = req.AddRecord(int32(tcaplus_protocol_cs.TCAPLUS_LIST_LAST_INDEX))
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}

		if i == 29 {
			data.Gamesvrid = ""
		}

		if _, err := rec.SetPBData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}

		if err := req.SetResultFlag(2); err != nil {
			t.Errorf("SetResultFlag fail, %s", err.Error())
			return
		}

		if err := client.SendRequest(req); err != nil {
			t.Errorf("SendRequest fail, %s", err.Error())
			return
		}

		//recv resp
		resp, err = tools.RecvResponse(client)
		if err != nil {
			t.Errorf("recvResponse fail, %s", err.Error())
			return
		}

		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
			return
		}

		if 1 != resp.GetRecordCount() {
			t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
			return
		}
	}

	/*------------------------------------------------------- 获取记录 ------------------------------------------*/ //for {
	client, req = tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListGetAllReq, "tb_online_list")

	rec, err = req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	req.SetMultiResponseFlag(1)

	resps, err := client.DoMore(req, 5*time.Second)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	fmt.Println("len", len(resps))
	count := 0

	for _, resp := range resps {
		count += resp.GetRecordCount()
	}
	fmt.Println("count", count)

	if count != 30 {
		t.Errorf("count %d != 30", count)
		return
	}
}

// 条件符合
func TestListGetAllWithCondition(t *testing.T) {
	data := &tcaplusservice.ListUser{}
	data.Id = 1
	data.Name = "a"
	oldJson := tools.StToJson(data)
	fmt.Println(oldJson)

	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiListGetAllReq, "list_user")

	client.ListDeleteAll(data)
	for i := int32(1); i <= 10; i++ {
		data.Rank = i
		client.ListAddAfter(data, -1)
	}
	defer client.ListDeleteAll(data)

	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	if _, err := rec.SetPBData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	rec.SetCondition("rank > 2 and rank< 7")

	if err := client.SendRequest(req); err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err , %s", terror.GetErrMsg(err))
		return
	}

	if 4 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 4", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newData := &tcaplusservice.ListUser{}
		if err := record.GetPBData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newData)
		fmt.Println(newJson)
		if newData.Rank <= 2 || newData.Rank >= 7 {
			t.Errorf("newData.Rank <= 2 || newData.Rank >= 7")
			return
		}
	}
}
