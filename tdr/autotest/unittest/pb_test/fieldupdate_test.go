package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestFieldUpdateSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldUpdateReq, "game_players")

	oldData := &tcaplusservice.GamePlayers{}
	oldData.Pay = &tcaplusservice.Payment{Amount: 1, PayId: 2}
	oldData.PlayerId = 233
	oldData.PlayerName = "jiahua"

	oldData.PlayerEmail = "wang"
	client.Insert(oldData)
	newData := &tcaplusservice.GamePlayers{PlayerId: 233, PlayerName: "jiahua", PlayerEmail: "wang"}
	client.Get(newData)
	if !proto.Equal(oldData, newData) {
		t.Errorf("data diff")
		return
	}
	defer func() {
		oldData.PlayerEmail = "wang"
		client.Delete(oldData)
	}()

	oldJson := tools.StToJson(oldData)
	fmt.Println(oldJson)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	oldData.Pay = &tcaplusservice.Payment{Amount: 3, PayId: 1}
	if _, err := rec.SetPBFieldValues(oldData, []string{"pay.pay_id", "pay.amount"}); err != nil {
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

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}

	if 1 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newMsg := &tcaplusservice.GamePlayers{}
		err = record.GetPBFieldValues(newMsg)
		if err != nil {
			t.Errorf("GetPBData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newMsg)
		fmt.Println(newJson)

		if !proto.Equal(oldData, newMsg) {
			t.Errorf("resData != reqData")
			return
		}
	}
}

// case 更新附带条件与操作
func TestPBFieldUpdateWithOperateCondition(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldUpdateReq, "user")

	oldMsg := &tcaplusservice.User{}
	oldMsg.Id = 1
	oldMsg.Name = "a"
	oldMsg.Rank = 0
	client.Insert(oldMsg)
	defer client.Delete(oldMsg)
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	oldMsg.Rank = 1
	if _, err := rec.SetPBFieldValues(oldMsg, []string{"rank"}); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	// 条件
	rec.SetCondition("rank == 0")
	// 操作
	rec.SetOperation("PUSH gameids #[-1][$=123]", 0)

	req.SetResultFlag(2)

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
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}

	if 1 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newMsg := &tcaplusservice.User{}
		err = record.GetPBFieldValues(newMsg)
		if err != nil {
			t.Errorf("GetPBData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newMsg)
		fmt.Println(newJson)
	}

	msg := &tcaplusservice.User{Id: 1, Name: "a"}
	client.Get(msg)
	fmt.Println(msg)
	if msg.Rank != 1 || msg.Gameids[0] != 123 {
		t.Errorf("msg.Rank != 0 || msg.Gameids[0] != 123")
		return
	}
}

// case 更新附带条件
func TestPBFieldUpdateWithCondition(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldUpdateReq, "user")

	oldMsg := &tcaplusservice.User{}
	oldMsg.Id = 1
	oldMsg.Name = "a"
	oldMsg.Rank = 1
	client.Insert(oldMsg)
	defer client.Delete(oldMsg)
	oldJson := tools.StToJson(oldMsg)
	fmt.Println(oldJson)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	oldMsg.Rank = 10
	if _, err := rec.SetPBFieldValues(oldMsg, []string{"rank"}); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}

	// 条件
	rec.SetCondition("rank == 0")

	req.SetResultFlag(2)

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

	if err := resp.GetResult(); err != 281 {
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}
}

func initTbMap(intMapSize, strMapSize, intArraySize, strArraySize, allTypeArraySize int32) *tcaplusservice.TbMap {
	t := &tcaplusservice.TbMap{
		Id:   1,
		Name: "a",
		I32:  1,
		Str:  "a",
		AllType: &tcaplusservice.AllTypeT{
			I32: 1,
			U64: 1,
		},
		IntMap: make(map[int32]*tcaplusservice.AllTypeT),
		StrMap: make(map[string]*tcaplusservice.AllTypeT),
	}

	for i := int32(1); i <= intMapSize; i++ {
		t.IntMap[i] = &tcaplusservice.AllTypeT{
			I32: i,
			U64: uint64(i),
		}
	}

	for i := int32(1); i <= strMapSize; i++ {
		t.StrMap[fmt.Sprint(i)] = &tcaplusservice.AllTypeT{
			I32: i,
			U64: uint64(i),
		}
	}

	for i := int32(1); i <= intArraySize; i++ {
		t.I32Array = append(t.I32Array, i)
	}

	for i := int32(1); i <= strArraySize; i++ {
		t.StrArray = append(t.StrArray, fmt.Sprint(i))
	}

	for i := int32(1); i <= allTypeArraySize; i++ {
		t.AllTypeArray = append(t.AllTypeArray, &tcaplusservice.AllTypeT{I32: i, U64: uint64(i)})
	}

	return t
}

func TestPBFieldUpdateByPath(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldUpdateReq, "tb_map")

	initMsg := initTbMap(0, 0, 0, 0, 0)

	client.Insert(initMsg)
	defer client.Delete(initMsg)

	//add record
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}

	oldMsg := initTbMap(3, 4, 5, 4, 3)
	oldMsg.I32 = 11
	oldMsg.Str = "aaa"
	oldMsg.AllType.U64 = 111

	paths := []string{
		"i32",
		"PUSH str",
		"SET all_type.u64",
		"i32_array",
		"str_map['1']",
		"str_map['2']",
		"str_map['3']",
		"str_map['4']",
	}
	if _, err := rec.SetPBFieldValues(oldMsg, paths); err != nil {
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

	if err := resp.GetResult(); err != 0 {
		t.Errorf("resp.GetResult err %d, %s", err, terror.GetErrMsg(err))
		return
	}

	if 1 != resp.GetRecordCount() {
		t.Errorf("resp.GetRecordCount() %d != 1", resp.GetRecordCount())
		return
	}

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}

		newMsg := &tcaplusservice.TbMap{}
		err = record.GetPBFieldValues(newMsg)
		if err != nil {
			t.Errorf("GetPBData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newMsg)
		fmt.Println(newJson)

		if !proto.Equal(oldMsg, newMsg) {
			t.Errorf("resData != reqData")
			return
		}
	}

	initMsg = &tcaplusservice.TbMap{Id: 1, Name: "a"}
	client.Get(initMsg)
	newJson := tools.StToJson(initMsg)
	fmt.Println(newJson)
	if initMsg.I32 != 11 || initMsg.Str != "aaa" || initMsg.AllType.U64 != 111 || len(initMsg.I32Array) != 5 ||
		len(initMsg.StrArray) != 0 || len(initMsg.StrMap) != 4 {
		t.Errorf("resData != reqData")
		return
	}

	oldMsg = initTbMap(3, 5, 5, 4, 3)
	oldMsg.StrMap["3"].I32 = 33
	oldMsg.StrMap["4"].U64 = 444
	paths = []string{
		"POP str",
		"POP all_type.u64",
		"POP i32_array",
		"POP str_map['1']",
		"POP str_map['2']",
		"POP str_map['22']",
		"str_map['3'].i32",
		"str_map['4'].u64",
		"str_map['5']",
		"all_type_array",
		"int_map",
	}
	err = client.FieldUpdate(oldMsg, paths)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	initMsg = &tcaplusservice.TbMap{Id: 1, Name: "a"}
	err = client.Get(initMsg)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	newJson = tools.StToJson(initMsg)
	fmt.Println(newJson)
	if initMsg.Str != "" || initMsg.AllType.U64 != 0 || len(initMsg.I32Array) != 0 ||
		len(initMsg.StrArray) != 0 || len(initMsg.StrMap) != 3 || initMsg.StrMap["3"].I32 != 33 ||
		initMsg.StrMap["4"].U64 != 444 || initMsg.StrMap["5"].I32 != 5 || len(initMsg.AllTypeArray) != 3 ||
		len(initMsg.IntMap) != 3 {
		t.Errorf("resData != reqData")
		return
	}
}

func TestPBFieldUpdateByPathCorrect(t *testing.T) {
	client, _ := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiPBFieldUpdateReq, "tb_map")

	initMsg := initTbMap(1, 0, 0, 0, 0)

	client.Insert(initMsg)
	defer client.Delete(initMsg)

	oldMsg := initTbMap(3, 5, 5, 4, 3)
	err := client.FieldUpdate(oldMsg, []string{"PUSH int_map[1]"})
	if err == nil || err.(*terror.ErrorCode).Code != terror.COMMON_ERR_INVALID_FIELD_NAME {
		fmt.Println(err)
		t.Error("err.(*terror.ErrorCode).Code != terror.COMMON_ERR_INVALID_FIELD_NAME")
		return
	}
	err = client.FieldUpdate(oldMsg, []string{"PUSH i32_array[-1]"})
	if err == nil || err.(*terror.ErrorCode).Code != terror.SVR_ERR_FAIL_PROTOBUF_FIELD_UPDATE {
		fmt.Println(err)
		t.Errorf("err.(*terror.ErrorCode).Code != terror.GEN_ERR_ERR")
		return
	}
	err = client.FieldUpdate(oldMsg, []string{"POP str_map['3'].ri32[0]"})
	if err == nil || err.(*terror.ErrorCode).Code != terror.COMMON_ERR_INVALID_FIELD_NAME {
		fmt.Println(err)
		t.Errorf("err.(*terror.ErrorCode).Code != terror.COMMON_ERR_INVALID_EXPR_TYPE")
		return
	}
	err = client.FieldUpdate(oldMsg, []string{"SET str_map['10']"})
	if err == nil || err.(*terror.ErrorCode).Code != terror.COMMON_ERR_INVALID_FIELD_NAME {
		fmt.Println(err)
		t.Errorf("err.(*terror.ErrorCode).Code != terror.COMMON_ERR_INVALID_FIELD_NAME")
		return
	}
}
