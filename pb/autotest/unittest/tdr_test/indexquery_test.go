package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"testing"
)

// 用法参考  https://iwiki.woa.com/pages/viewpage.action?pageId=419645505

// 普通查询
func TestIndexQuerySimple(t *testing.T) {
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query := fmt.Sprintf("select * from table_generic where uin > 100 limit 10 offset 0; ")
	req.SetSql(query)

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

	for i := 0; i < resp.GetRecordCount(); i++ {
		record, err := resp.FetchRecord()
		if err != nil {
			t.Errorf("FetchRecord failed %s", err.Error())
			return
		}
		newData := tcaplus_tb.NewTable_Generic()
		if err := record.GetData(newData); err != nil {
			t.Errorf("record.GetData failed %s", err.Error())
			return
		}

		newJson := tools.StToJson(newData)
		logger.DEBUG("%s", newJson)
		fmt.Println(newJson)
		//if newJson != oldJson {
		//	t.Errorf("resData != reqData")
		//	return
		//}
	}
}

// 聚合查询
func TestIndexQueryAggregation(t *testing.T) {
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	// count(uin)
	query := fmt.Sprintf("select count(uin), count(distinct(name)), sum(uin), max(uin), min(uin), avg(uin) from table_generic where uin > 0")
	req.SetSql(query)

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

	r, err := resp.ProcAggregationSqlQueryType()
	if err != nil {
		t.Errorf("ProcAggregationSqlQueryType failed %s", err.Error())
		return
	}
	logger.DEBUG("%s", common.CovertToJson(r))
}
