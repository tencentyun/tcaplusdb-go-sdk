package api_test

import (
	"fmt"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/table/tcaplus_tb"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/common"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
	"testing"
)

// 用法参考  https://iwiki.woa.com/pages/viewpage.action?pageId=419645505

// 普通查询
func TestIndexQuerySimple(t *testing.T) {
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query := fmt.Sprintf("select * from table_generic where uin > 100 and uin < 88888888 limit 10 offset 0; ")
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

func TestIndexQueryBetween(t *testing.T) {
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query := fmt.Sprintf("select * from table_generic where uin between 1 and 100 and level > 0 and level < 100; ")
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

func TestIndexQueryLike(t *testing.T) {
	// 匹配 %
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query := `select * from table_generic where name like "nam%" ;`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
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

	// 匹配  _
	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query = `select * from table_generic where  key3 like "key_" ;`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
	req.SetSql(query)

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

	// 匹配  *
	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query = `select * from table_generic where key4 like "*ey4";`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
	req.SetSql(query)

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

	// 完全匹配
	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query = `select * from table_generic where key4 like "key4";`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
	req.SetSql(query)

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

func TestIndexQueryNotLike(t *testing.T) {
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query := `select * from table_generic where name not like "key%";`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
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
	query := fmt.Sprintf("select count(uin), count(distinct(name)), sum(uin), max(uin), min(uin), avg(uin) from table_generic where uin > 0 limit 2 offset 0;")
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

// 部分字段
func TestIndexQueryField(t *testing.T) {
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query := `select uin from table_generic where name not like "key%";`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
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

// 部分字段
func TestIndexQueryFix(t *testing.T) {
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query := `select *, uin from table_generic where name not like "key%";`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
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

	if err := resp.GetResult(); err != terror.PROXY_ERR_THIS_SQL_IS_NOT_SUPPORT {
		t.Errorf("resp.GetResult err %d , %s", err, terror.GetErrMsg(err))
		return
	}

	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query = `select sum(uin), uin from table_generic where name not like "key%";`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
	req.SetSql(query)

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

	if err := resp.GetResult(); err != terror.PROXY_ERR_THIS_SQL_IS_NOT_SUPPORT {
		t.Errorf("resp.GetResult err %d , %s", err, terror.GetErrMsg(err))
		return
	}

	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query = `select sum(*), * from table_generic where name not like "key%";`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
	req.SetSql(query)

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

	if err := resp.GetResult(); err != terror.PROXY_ERR_THIS_SQL_IS_NOT_SUPPORT {
		t.Errorf("resp.GetResult err %d , %s", err, terror.GetErrMsg(err))
		return
	}
}

func TestIndexQueryNotSupport(t *testing.T) {
	client, req := tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query := `select * from table_generic where uin > 0 order by uin;`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
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

	if err := resp.GetResult(); err != terror.PROXY_ERR_THIS_SQL_IS_NOT_SUPPORT {
		t.Errorf("resp.GetResult err %d , %s", err, terror.GetErrMsg(err))
		return
	}

	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query = `select * from table_generic where uin > 0 group by uin;`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
	req.SetSql(query)

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

	if err := resp.GetResult(); err != terror.PROXY_ERR_THIS_SQL_IS_NOT_SUPPORT {
		t.Errorf("resp.GetResult err %d , %s", err, terror.GetErrMsg(err))
		return
	}

	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query = `select * from table_generic where uin > 0 having sum(uin) > 100;`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
	req.SetSql(query)

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

	if err := resp.GetResult(); err != terror.PROXY_ERR_QUERY_FOR_CONVERT_TCAPLUS_REQ_TO_INDEX_SERVER_REQ_FAILED {
		t.Errorf("resp.GetResult err %d , %s", err, terror.GetErrMsg(err))
		return
	}

	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query = `select * from table_generic where uin > 0 and level in (select level from table_generic where level > 0);`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
	req.SetSql(query)

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

	if err := resp.GetResult(); err != terror.TCAPLUS_INDEX_ERR_SEND_TO_INDEX_SERVER_FAILED_FOR_OTHER_REASON {
		t.Errorf("resp.GetResult err %d , %s", err, terror.GetErrMsg(err))
		return
	}

	client, req = tools.InitClientAndReqWithTableName(cmd.TcaplusApiSqlReq, "table_generic")

	query = `select uin myuin from table_generic where uin > 0 and level in (select level from table_generic where level > 0);`

	//query := fmt.Sprintf(`select * from table_generic where name like "name"; `)
	req.SetSql(query)

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

	if err := resp.GetResult(); err != terror.TCAPLUS_INDEX_ERR_SEND_TO_INDEX_SERVER_FAILED_FOR_OTHER_REASON {
		t.Errorf("resp.GetResult err %d , %s", err, terror.GetErrMsg(err))
		return
	}

}
