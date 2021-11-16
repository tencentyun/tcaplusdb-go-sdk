package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/request"
	"testing"
	"time"
)

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

func BenchmarkPack(b *testing.B) {
	logger.SetLogCfg("../cfg/logconf.xml")
	logger.Init()

	//count :=0
	req, _ := request.NewRequest(1, 3, "test", cmd.TcaplusApiReplaceReq, false)
	rec, _ := req.AddRecord(0)
	data := newGenericTableRec()
	//data.Binary_Count =10240
	//data.Max_Binary =  make([]int8, 10240, 10241)
	rec.SetData(data)
	buf, _ := data.Pack(0)
	for n := 0; n < b.N; n++ {
		req, _ = request.NewRequest(1, 3, "test", cmd.TcaplusApiReplaceReq, false)
		rec, _ = req.AddRecord(0)

		//data.Binary_Count =10240
		//data.Max_Binary =  make([]int8, 10240, 10241)
		rec.SetData(data)
		//req.Pack()
		//buf,_ = data.Pack(0)
	}
	fmt.Println(rec)
	fmt.Println(buf)

}
