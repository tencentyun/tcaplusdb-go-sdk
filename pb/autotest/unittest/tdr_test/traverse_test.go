package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"testing"
	"time"
)

func TestTDRTraverse(t *testing.T) {
	client, _ := tools.InitClientAndReqWithTableName(cmd.TcaplusApiTableTraverseReq, "table_generic")

	tra := client.GetTraverser(cfg.ApiConfig.ZoneId, "table_generic")
	defer tra.Stop()

	tra.SetFieldNames([]string{"level", "count", "info"})
	tra.SetLimit(10)

	resps, err := client.DoTraverse(tra, 60*time.Second)
	if err != nil {
		t.Errorf("RecvResponse fail, %s", err.Error())
		return
	}
	fmt.Println(len(resps))

	for _, resp := range resps {
		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult err %s", terror.GetErrMsg(err))
			return
		}

		for i := 0; i < resp.GetRecordCount(); i++ {
			record, err := resp.FetchRecord()
			if err != nil {
				t.Errorf("FetchRecord failed %s", err.Error())
				return
			}

			key1, err := record.GetKeyInt64("uin")
			if err != nil {
				t.Errorf("record.GetValueInt8 failed %s", err.Error())
				return
			}
			value1, err := record.GetValueInt32("level")
			if err != nil {
				t.Errorf("record.GetValueInt8 failed %s", err.Error())
				return
			}
			value2, err := record.GetValueInt8("count")
			if err != nil {
				t.Errorf("record.GetValueInt8 failed %s", err.Error())
				return
			}
			value3, err := record.GetValueStr("info")
			if err != nil {
				t.Errorf("record.GetValueInt8 failed %s", err.Error())
				return
			}
			fmt.Println(key1, value1, value2, value3)
		}
	}
}
