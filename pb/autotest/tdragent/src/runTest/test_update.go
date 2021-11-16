package runTest

import (
	"errors"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/tdragent/src/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"time"
)

func (r *RunTest) TestUpdate() error {
	req, err := r.TcaplusClient.NewRequest(uint32(r.TcaplusCase.Head.ZoneID), r.TcaplusCase.Head.TableName, cmd.TcaplusApiUpdateReq)
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

	var recKey *RecKeyInfo
	if len(r.keyCache) <= 0 {
		logger.DEBUG("keyCache zero")
		if tRecKey, err := r.MakeKey(rec); err != nil {
			logger.ERR("MakeKey err %s", err.Error())
			return err
		} else {
			recKey = tRecKey
		}
	} else {
		//从cache中update
		for _, v := range r.keyCache {
			if tRecKey, err := r.MakeKeyFromCacheRec(rec, v); err != nil {
				logger.ERR("MakeKeyFromCacheRec err %s", err.Error())
				return err
			} else {
				recKey = tRecKey
				break
			}
		}
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
