package api_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/request"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"time"
)

var TestTableName = "table_generic"

//同步接收
func recvResponse(client *tcaplus.Client) (response.TcaplusResponse, error) {
	//recv response
	timeOutChan := time.After(5 * time.Second)
	for {
		select {
		case <-timeOutChan:
			return nil, errors.New("5s timeout")
		default:
			resp, err := client.RecvResponse()
			if err != nil {
				return nil, err
			} else if resp == nil {
				time.Sleep(time.Microsecond * 1)
			} else {
				return resp, nil
			}
		}
	}
}

//将结构体转成json格式
func StToJson(args interface{}) string {
	js, err := json.Marshal(args)
	if err != nil {
		return fmt.Sprintf("%v", args)
	}
	return string(js)
}

//new GenericTable表的一条记录，并赋值
func newGenericTableRec() *tcaplus_tb.Table_Generic {
	//data
	data := tcaplus_tb.NewTable_Generic()
	data.Uin = uint64(time.Now().UnixNano())
	data.Name = "GoUnitTest"
	data.Key3 = "key3"
	data.Key4 = "key4"
	data.Level = 100
	data.Info = "info"
	data.Float_Score = 100.32
	data.Double_Score = 100.64

	data.Count = 2
	data.Items = []uint64{101, 102}

	data.Big_Record_1 = "Big_Record_1"
	data.Big_Record_2 = "Big_Record_2"
	data.Big_Record_3 = "Big_Record_3"
	data.Big_Record_4 = "Big_Record_4"
	data.Big_Record_5 = "Big_Record_5"
	data.C_Int8 = 103
	data.C_Uint8 = 104
	data.C_Int16 = 105
	data.C_Uint16 = 106
	data.C_Int32 = 107
	data.C_Int64 = 109
	data.C_Uint64 = 110
	data.C_Float = 111.32
	data.C_Double = 112.64
	data.C_String = "C_String"

	data.C_Uint32 = 2
	data.C_Binary = []int8{113, 114}

	data.Max_String = "Max_String"

	data.Binary_Count = 3
	data.Max_Binary = []int8{115, 116, 117}

	data.Single_Struct.X = 119
	data.Single_Struct.Y = 120
	data.Single_Struct.Score = 121.64
	data.Single_Struct.Rank = 122
	data.Single_Struct.Title = "Title"
	data.Simple_Struct.C_Int8 = 123
	data.Simple_Struct.C_Uint8 = 124
	data.Simple_Struct.C_Int16 = 125
	data.Simple_Struct.C_Uint16 = 126
	data.Simple_Struct.C_Int32 = 127
	data.Simple_Struct.C_Uint32 = 128
	data.Simple_Struct.C_Int64 = 129
	data.Simple_Struct.C_Uint64 = 130

	data.Single_Union_Selector = 1
	data.Single_Union.Name = "data.Single_Union.Name"

	data.Array_Count = 1
	tableInfo := tcaplus_tb.NewTableInfo()
	tableInfo.C.D = 131
	tableInfo.C.Test = 132
	tableInfo.Test = 133
	tableInfo.String_Array = "String_Array"
	tableInfo.Count = 2
	tableInfo.Binary = []int8{34, 35}
	tableInfo.Bound_31_Byte_Test_012345678901 = 136
	data.Array = []*tcaplus_tb.TableInfo{tableInfo}

	data.Selector = 1
	data.C_Union.Name = "data.C_Union.Name"

	unionSt := tcaplus_tb.NewUnion_Type(1)
	unionSt.Name = "data.Union_Array[0].Name"
	data.Union_Array = []*tcaplus_tb.Union_Type{unionSt}

	data.C_Struct.X = 137
	data.C_Struct.Y = 138
	data.C_Struct.Score = 139.64
	data.C_Struct.Rank = 140
	data.C_Struct.Title = "data.C_Struct.Title"
	data.C_Struct.Level2_Struct.Uin = 141
	data.C_Struct.Level2_Struct.Name = "data.C_Struct.Level2_Struct.Name"

	tmpSt := tcaplus_tb.NewStruct_Type()
	tmpSt.X = 142
	tmpSt.Y = 143
	tmpSt.Score = 144.64
	tmpSt.Rank = 145
	tmpSt.Title = "data.Struct_Array[0].Title"
	tmpSt.Level2_Struct.Uin = 146
	tmpSt.Level2_Struct.Name = "data.Struct_Array[0].Level2_Struct.Name"
	data.Struct_Array = []*tcaplus_tb.Struct_Type{tmpSt}

	data.Bound_31_Byte_Test_012345678901 = 147

	data.Simple_Array_Count = 3
	data.Int_Array = []int32{148, 149, 150}
	data.Double_Array = []float64{151.64, 152.64, 153.64}
	return data
}

func randomStrChar(in string, extern int) string {
	s := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var idx int = 0
	for ; idx < extern; idx++ {
		in += string(s[int64(time.Now().UnixNano()%int64(len(s)))])
	}
	return in
}

var client *tcaplus.Client = nil

func InitClientAndReq(cmd int) (*tcaplus.Client, request.TcaplusRequest) {
	return InitClientAndReqWithTableName(cmd, TestTableName)
}

func InitClientAndReqWithTableName(cmd int, tableName string) (*tcaplus.Client, request.TcaplusRequest) {
	if nil == client {
		if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
			fmt.Printf("ReadApiCfg fail %s", err.Error())
			return nil, nil

		}

		client = tcaplus.NewClient()
		if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
			fmt.Printf("excepted SetLogCfg success")
			return nil, nil
		}

		err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
		if err != nil {
			fmt.Printf("excepted dial success, %s", err.Error())
			return nil, nil
		}
	}
	req, err := client.NewRequest(cfg.ApiConfig.ZoneId, tableName, cmd)
	if err != nil {
		fmt.Printf("NewRequest fail, %s", err.Error())
		return nil, nil
	}
	return client, req
}

func AsyncSendAndGetRes(client *tcaplus.Client, req request.TcaplusRequest) (response.TcaplusResponse, error) {
	if err := client.SendRequest(req); err != nil {
		return nil, err
	}
	return recvResponse(client)
}

func InsertKV(uin uint64, key4 string) {
	client, req := InitClientAndReq(cmd.TcaplusApiInsertReq)
	if nil == client || nil == req {
		fmt.Printf("NewRequest fail")
		return
	}

	//data
	data := newGenericTableRec()
	data.Uin = uin
	data.Name = randomStrChar("nm", 10)
	data.Key3 = randomStrChar("3", 10)
	data.Key4 = key4

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
