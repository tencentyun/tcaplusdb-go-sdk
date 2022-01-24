package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/option"
	"testing"
	"time"
)

//case1 BatchInsert success
func TestDoIncreaseSuccess(t *testing.T) {
	client, err := tools.InitClient()
	if err != nil {
		t.Errorf("InitClient failed %s", err.Error())
		return
	}

	//插入记录成功
	uin := uint64(time.Now().UnixNano())
	data := newGenericTableRec()
	data.Uin = uin
	data.Name = fmt.Sprintf("%d", 2)
	data.Level = int32(2)
	data.Float_Score = float32(6.6)
	data.Double_Score = float64(8.8)
	data.Info = fmt.Sprintf("%d", 2)
	err = client.DoInsert(TestTableName, data, nil)
	if err != nil {
		t.Errorf("DoBatchInsert fail, %s", err.Error())
		return
	}

	//自增记录
	opt := &option.TDROpt{
		ResultFlagForSuccess: option.TcaplusResultFlagAllNewValue,
		IncField: []option.IncFieldInfo{
			option.IncFieldInfo{
				FieldName: "level",
				IncData:   int32(2),
				Operation: cmd.TcaplusApiOpPlus,
			},
		},
	}
	err = client.DoIncrease(TestTableName, data, opt)
	if err != nil {
		t.Errorf("DoIncrease fail, %s", err.Error())
		return
	}

	if data.Level != 4 {
		t.Errorf("data.Level invalid %v", data.Level)
		return
	}
}
