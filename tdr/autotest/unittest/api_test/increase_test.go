package api_test

import (
	"bytes"
	"encoding/binary"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/table/tcaplus_tb"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"testing"
	"time"
)

// 测试没有记录，返回错误码，且不容许创建新记录
func TestIncrease_no_record(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiIncreaseReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	req.SetResultFlag(1)
	nameList := []string{string("level")}
	req.SetFieldNames(nameList)
	data := newGenericTableRec()
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	data.Uin = uint64(time.Now().UnixNano())
	data.Name = "GoUnitTest"
	data.Key3 = "key_3"
	data.Simple_Struct.C_Int64 = 100
	data.Key4 = "key4___12"
	//data.Level = 8

	s1 := make([]byte, 0)
	buf := bytes.NewBuffer(s1)
	var i1 int32 = 8
	binary.Write(buf, binary.LittleEndian, i1)

	rec.AddValueOperation("level", buf.Bytes(), uint32(buf.Len()), cmd.TcaplusApiOpPlus, 0, 400)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	req.SetAddableIncreaseFlag(0)
	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == terror.TXHDB_ERR_RECORD_NOT_EXIST {
		return
	} else {
		t.Errorf("resp.GetResult expect TXHDB_ERR_RECORD_NOT_EXIST ,but %s", terror.GetErrMsg(err))
		return
	}
}

// 设置容许自增没记录的时候先创建记录，但是没有初始值，这样自增会失败
func TestIncrease_no_record_no_defalutvalue_err(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiIncreaseReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	req.SetResultFlag(1)
	nameList := []string{string("c_uint64"), string("level")}
	req.SetFieldNames(nameList)
	data := newGenericTableRec()
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	data.Uin = uint64(time.Now().UnixNano()) + 1
	data.Name = "GoUnitTest"
	data.Key3 = "key_3__"
	data.Key4 = "key4___12"
	data.Simple_Struct.C_Int64 = 100
	s1 := make([]byte, 0)
	buf := bytes.NewBuffer(s1)
	var i1 int32 = 8
	binary.Write(buf, binary.LittleEndian, i1)

	rec.AddValueOperation("c_uint64", buf.Bytes(), uint32(buf.Len()), cmd.TcaplusApiOpPlus, 0, 400)
	//rec.AddValueOperation("level", buf.Bytes(), uint32(buf.Len()), cmd.TcaplusApiOpPlus,0, 400)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	req.SetAddableIncreaseFlag(1)
	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != terror.SVR_ERR_FAIL_INVALID_FIELD_TYPE {
		t.Errorf("resp.GetResult expect error SVR_ERR_FAIL_INVALID_FIELD_TYPE ,but error %s, errid: %d", terror.GetErrMsg(err), err)
		return
	}
	return
}

// 正常自增
func TestIncrease(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiIncreaseReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	req.SetResultFlag(1)
	nameList := []string{string("level")}
	req.SetFieldNames(nameList)
	data := newGenericTableRec()
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	data.Uin = uint64(time.Now().UnixNano())
	data.Name = "GoUnitTest"
	data.Key3 = "key3"

	data.Key4 = "key4___12"
	//data.Level = 8
	data.Simple_Struct.C_Int64 = 100
	s1 := make([]byte, 0)
	buf := bytes.NewBuffer(s1)
	var i1 int32 = 8
	binary.Write(buf, binary.LittleEndian, i1)

	rec.AddValueOperation("level", buf.Bytes(), uint32(buf.Len()), cmd.TcaplusApiOpPlus, 0, 0xffffffff)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	req.SetAddableIncreaseFlag(1)
	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err == 0 {
		for idx := int(0); idx < resp.GetRecordCount(); idx++ {
			rec, err := resp.FetchRecord()
			if err == nil {
				oldData := tcaplus_tb.NewTable_Generic()
				if err := rec.GetData(oldData); err == nil {
					if oldData.Level != i1+1 {
						t.Errorf("resp.GetResult expect %d ,but result %d", i1+1, oldData.Level)
						return
					} else {
						// here is right result
						return
					}
				} else {
					t.Errorf("resp.GetResult expect no error ,but error %s", err.Error())
					return
				}
			} else {
				t.Errorf("resp.GetResult expect no error ,but error %s", err.Error())
				return
			}
		}

	} else {
		t.Errorf("resp.GetResult expect no error ,but error %s, errid: %d", terror.GetErrMsg(err), err)
		return
	}
	t.Errorf("resp.GetResult expect no error ,but error ")
	return
}

// 测试自增越界
func TestIncrease_out_range(t *testing.T) {
	client, req := InitClientAndReq(cmd.TcaplusApiIncreaseReq)
	if nil == client || nil == req {
		t.Errorf("init client and req fail")
		return
	}
	req.SetAsyncId(723)
	req.SetResultFlag(1)
	nameList := []string{string("level")}
	req.SetFieldNames(nameList)
	data := newGenericTableRec()
	rec, err := req.AddRecord(0)
	if err != nil {
		t.Errorf("AddRecord fail, %s", err.Error())
		return
	}
	uin := uint64(time.Now().UnixNano())
	data.Uin = uin
	data.Name = "GoUnitTest"
	data.Key3 = "key3"

	data.Key4 = "key4___12"
	//data.Level = 8

	s1 := make([]byte, 0)
	buf := bytes.NewBuffer(s1)
	var i1 int32 = 8
	binary.Write(buf, binary.LittleEndian, i1)

	rec.AddValueOperation("level", buf.Bytes(), uint32(buf.Len()), cmd.TcaplusApiOpPlus, 0, 0x7)

	if err := rec.SetData(data); err != nil {
		t.Errorf("SetData fail, %s", err.Error())
		return
	}
	req.SetAddableIncreaseFlag(1)
	//recv resp
	resp, err := AsyncSendAndGetRes(client, req)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if err := resp.GetResult(); err != terror.SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE {
		t.Errorf("resp.GetResult expect error SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE,but error %s, errid: %d", terror.GetErrMsg(err), err)
		return
	}
}

// 测试减法,分三步， 正常加法，正常减法，越界减法
func TestIncrease_out_range_sub(t *testing.T) {
	uin := uint64(time.Now().UnixNano())
	data := newGenericTableRec()
	data.Uin = uin
	data.Name = "GoUnitTest"
	data.Key3 = "key3"
	data.Key4 = "key4___12"
	//data.Level = 8

	// 自增插入一条数据
	{
		client, req := InitClientAndReq(cmd.TcaplusApiIncreaseReq)
		if nil == client || nil == req {
			t.Errorf("init client and req fail")
			return
		}
		req.SetAsyncId(723)
		req.SetResultFlag(1)
		nameList := []string{string("level")}
		req.SetFieldNames(nameList)

		rec, err := req.AddRecord(0)
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}

		s1 := make([]byte, 0)
		buf := bytes.NewBuffer(s1)
		var i1 int32 = 8
		binary.Write(buf, binary.LittleEndian, i1)
		// 不检查是否越界
		rec.AddValueOperation("level", buf.Bytes(), uint32(buf.Len()), cmd.TcaplusApiOpPlus, 0, 0)

		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		req.SetAddableIncreaseFlag(1)
		//recv resp
		resp, err := AsyncSendAndGetRes(client, req)
		if err != nil {
			t.Errorf("recvResponse fail, %s", err.Error())
			return
		}
		if err := resp.GetResult(); err != 0 {
			t.Errorf("resp.GetResult expect no error ,but error %s, errid: %d", terror.GetErrMsg(err), err)
			return
		}
	}
	// 减法正常操作
	{
		client, req := InitClientAndReq(cmd.TcaplusApiIncreaseReq)
		if nil == client || nil == req {
			t.Errorf("init client and req fail")
			return
		}
		req.SetAsyncId(723)
		req.SetResultFlag(1)
		nameList := []string{string("level")}
		req.SetFieldNames(nameList)

		rec, err := req.AddRecord(0)
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}

		s1 := make([]byte, 0)
		buf := bytes.NewBuffer(s1)
		var i1 int32 = 2
		binary.Write(buf, binary.LittleEndian, i1)
		// 不检查是否越界
		rec.AddValueOperation("level", buf.Bytes(), uint32(buf.Len()), cmd.TcaplusApiOpMinus, 0, 0)

		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		req.SetAddableIncreaseFlag(1)
		//recv resp
		resp, err := AsyncSendAndGetRes(client, req)
		if err != nil {
			t.Errorf("recvResponse fail, %s", err.Error())
			return
		}
		if err := resp.GetResult(); err == 0 {
			for idx := int(0); idx < resp.GetRecordCount(); idx++ {
				rec, err := resp.FetchRecord()
				if err == nil {
					oldData := tcaplus_tb.NewTable_Generic()
					if err := rec.GetData(oldData); err == nil {
						if 7 != oldData.Level {
							t.Errorf("resp.GetResult expect %d, but %d", 7, oldData.Level)
							return
						}
					} else {
						t.Errorf("resp.GetResult expect no error ,but error %s", err.Error())
						return
					}

				} else {
					t.Errorf("resp.GetResult expect no error ,but error %s", err.Error())
					return
				}
			}
		} else {
			t.Errorf("resp.GetResult expect no error ,but error")
			return
		}

	}
	// 减法越界检查
	{
		client, req := InitClientAndReq(cmd.TcaplusApiIncreaseReq)
		if nil == client || nil == req {
			t.Errorf("init client and req fail")
			return
		}
		req.SetAsyncId(723)
		req.SetResultFlag(1)
		nameList := []string{string("level")}
		req.SetFieldNames(nameList)

		rec, err := req.AddRecord(0)
		if err != nil {
			t.Errorf("AddRecord fail, %s", err.Error())
			return
		}

		s1 := make([]byte, 0)
		buf := bytes.NewBuffer(s1)
		var i1 int32 = 3
		binary.Write(buf, binary.LittleEndian, i1)
		// 不检查是否越界
		rec.AddValueOperation("level", buf.Bytes(), uint32(buf.Len()), cmd.TcaplusApiOpMinus, 5, 100)

		if err := rec.SetData(data); err != nil {
			t.Errorf("SetData fail, %s", err.Error())
			return
		}
		req.SetAddableIncreaseFlag(1)
		//recv resp
		resp, err := AsyncSendAndGetRes(client, req)
		if err != nil {
			t.Errorf("recvResponse fail, %s", err.Error())
			return
		}

		if err := resp.GetResult(); err != terror.SVR_ERR_FAIL_OUT_OF_USER_DEF_RANGE {
			t.Errorf("resp.GetResult expect no error ,but error %s, errid: %d", terror.GetErrMsg(err), err)
			return
		}

	}
}
