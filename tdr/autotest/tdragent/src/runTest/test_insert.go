package runTest

import (
	"errors"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/autotest/tdragent/src/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"time"
)

func (r *RunTest) TestInsert() error {
	req, err := r.TcaplusClient.NewRequest(uint32(r.TcaplusCase.Head.ZoneID), r.TcaplusCase.Head.TableName, cmd.TcaplusApiInsertReq)
	if err != nil {
		logger.ERR("NewRequest err %s", err.Error())
		return errors.New("NewRequest err")
	}
	req.SetAsyncId(r.asyncId)

	if err := req.SetResultFlag(r.TcaplusCase.Head.ResultFlag); err != nil {
		logger.ERR("SetResultFlag err %s", err.Error())
		return errors.New("SetResultFlag err")
	}

	rec, err := req.AddRecord(0)
	if err != nil {
		logger.ERR("AddRecord err %s", err.Error())
		return errors.New("AddRecord err")
	}

	recKey, err := r.MakeKey(rec)
	if err != nil {
		logger.ERR("MakeKey err %s", err.Error())
		return err
	}

	if err := r.MakeValue(rec); err != nil {
		logger.ERR("MakeValue err %s", err.Error())
		return err
	}

	//发送请求
	go r.TcaplusClient.SendRequest(req)
	//if err := r.TcaplusClient.SendRequest(req); err != nil {
	//	logger.ERR("SendRequest failed %v", err.Error())
	//	return err
	//}

	recKey.SendTime = time.Now()
	r.sendKeyMap[r.asyncId] = recKey
	return nil
}
