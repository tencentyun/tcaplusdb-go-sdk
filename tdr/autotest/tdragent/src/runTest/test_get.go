package runTest

import (
	"errors"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/tdragent/src/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"time"
)

func (r *RunTest) TestGet() error {
	req, err := r.TcaplusClient.NewRequest(uint32(r.TcaplusCase.Head.ZoneID), r.TcaplusCase.Head.TableName, cmd.TcaplusApiGetReq)
	if err != nil {
		logger.ERR("NewRequest err %s", err.Error())
		return errors.New("NewRequest err")
	}
	req.SetAsyncId(r.asyncId)
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
		//从cache中get
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

	fieldName := make([]string, len(r.TcaplusCase.Body.ValueInfoList), len(r.TcaplusCase.Body.ValueInfoList))
	for i, vInfo := range r.TcaplusCase.Body.ValueInfoList {
		fieldName[i] = vInfo.FieldName
	}
	req.SetFieldNames(fieldName)

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
