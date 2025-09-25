package api_test

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/option"
	"testing"
	"time"
)

// insert update replace 不设置resultFlag也会返回version

//insert返回version
func TestInsertVersion(t *testing.T) {
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

	opt := &option.TDROpt{}
	err = client.DoInsert(TestTableName, data, opt)
	if err != nil {
		t.Errorf("DoInsert fail, %s", err.Error())
		return
	}
	if opt.Version != 1 {
		t.Errorf("DoInsert version %v != 1", opt.Version)
		return
	}
	fmt.Println("insert version", opt.Version)
}

//update返回version
func TestUpdateVersion(t *testing.T) {
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

	opt := &option.TDROpt{}
	err = client.DoInsert(TestTableName, data, opt)
	if err != nil {
		t.Errorf("DoInsert fail, %s", err.Error())
		return
	}
	if opt.Version != 1 {
		t.Errorf("DoInsert version %v != 1", opt.Version)
		return
	}
	fmt.Println("insert version", opt.Version)

	opt = &option.TDROpt{}
	err = client.DoUpdate(TestTableName, data, opt)
	if err != nil {
		t.Errorf("DoUpdate fail, %s", err.Error())
		return
	}
	if opt.Version != 2 {
		t.Errorf("DoUpdate version %v != 1", opt.Version)
		return
	}
	fmt.Println("update version", opt.Version)
}

//Replace返回version
func TestReplaceVersion(t *testing.T) {
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

	opt := &option.TDROpt{}
	err = client.DoReplace(TestTableName, data, opt)
	if err != nil {
		t.Errorf("DoReplace fail, %s", err.Error())
		return
	}
	if opt.Version != 1 {
		t.Errorf("DoReplace version %v != 1", opt.Version)
		return
	}
	fmt.Println("1 DoReplace version", opt.Version)

	opt = &option.TDROpt{}
	err = client.DoReplace(TestTableName, data, opt)
	if err != nil {
		t.Errorf("DoReplace fail, %s", err.Error())
		return
	}
	if opt.Version != 2 {
		t.Errorf("DoReplace version %v != 1", opt.Version)
		return
	}
	fmt.Println("2 DoReplace version", opt.Version)
}
