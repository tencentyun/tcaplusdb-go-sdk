package response

import (
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/tcaplus_protocol_cs"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
)

type listReplaceResponse struct {
	record  *record.Record
	pkg     *tcaplus_protocol_cs.TCaplusPkg
	offset  int32
	idx     int32
	listidx int32
}

func newListReplaceResponse(pkg *tcaplus_protocol_cs.TCaplusPkg) (*listReplaceResponse, error) {
	if pkg == nil || pkg.Body == nil || pkg.Body.ListReplaceRes == nil {
		return nil, &terror.ErrorCode{Code: terror.API_ERR_PARAMETER_INVALID, Message: "pkg init fail"}
	}
	return &listReplaceResponse{pkg: pkg}, nil
}

func (res *listReplaceResponse) GetResult() int {
	ret := int(res.pkg.Body.ListReplaceRes.Result)
	return ret
}

func (res *listReplaceResponse) GetTableName() string {
	tableName := string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen])
	return tableName
}

func (res *listReplaceResponse) GetAppId() uint64 {
	return uint64(res.pkg.Head.RouterInfo.AppID)
}

func (res *listReplaceResponse) GetZoneId() uint32 {
	return uint32(res.pkg.Head.RouterInfo.ZoneID)
}

func (res *listReplaceResponse) GetCmd() int {
	return int(res.pkg.Head.Cmd)
}

func (res *listReplaceResponse) GetAsyncId() uint64 {
	return res.pkg.Head.AsynID
}

/*
 各个命令请求需要具体问题具体分析
 新版的ResultFlag 的GetRecordCount的规则
 在成功的场景下, 设置了如下的ResultFlag返回数据规则
 1. TCaplusValueFlag_NOVALUE 不返还
 2. TCaplusValueFlag_SAMEWITHREQUEST 返回
 3. TCaplusValueFlag_ALLVALUE 返回
 4. TCaplusValueFlag_ALLOLDVALUE 返回

 在失败的场景下, 设置了如下的ResultFlag会返回数据
 1. TCaplusValueFlag_NOVALUE 不返还
 2. TCaplusValueFlag_SAMEWITHREQUEST 返回
 3. TCaplusValueFlag_ALLVALUE 返回
 4. TCaplusValueFlag_ALLOLDVALUE 返回

*/

func (res *listReplaceResponse) GetRecordCount() int {
	if 0 != (res.pkg.Body.ListReplaceRes.Flag & (1 << 6)) {
		//新版本的result flag 通过ResultFlag判断
		iResultFlagForSuccess := tcaplus_protocol_cs.TCaplusValueFlag_NOVALUE
		iResultFlagForFail := tcaplus_protocol_cs.TCaplusValueFlag_NOVALUE
		if res.pkg.Body.ListReplaceRes.Result == 0 {
			iResultFlagForSuccess = GetResultFlagByBit(res.pkg.Body.ListReplaceRes.Flag, true)
			if tcaplus_protocol_cs.TCaplusValueFlag_SAMEWITHREQUEST == iResultFlagForSuccess ||
				tcaplus_protocol_cs.TCaplusValueFlag_ALLVALUE == iResultFlagForSuccess ||
				tcaplus_protocol_cs.TCaplusValueFlag_ALLOLDVALUE == iResultFlagForSuccess {
				return 1
			}
		} else {
			iResultFlagForFail = GetResultFlagByBit(res.pkg.Body.ListReplaceRes.Flag, false)
			if tcaplus_protocol_cs.TCaplusValueFlag_SAMEWITHREQUEST == iResultFlagForFail ||
				(tcaplus_protocol_cs.TCaplusValueFlag_ALLOLDVALUE == iResultFlagForFail &&
					res.pkg.Body.ListReplaceRes.ResultInfo.ElementNum > 0) {
				return 1
			}
		}
	} else {
		//老版本的result flag 通过ResultFlag判断
		if 0 == res.pkg.Body.ListReplaceRes.Result ||
			terror.SVR_ERR_FAIL_INVALID_VERSION == int(res.pkg.Body.ListReplaceRes.Result) {
			return 1
		}
	}
	return 0
}

func (res *listReplaceResponse) FetchRecord() (*record.Record, error) {
	data := res.pkg.Body.ListReplaceRes.ResultInfo
	if res.GetRecordCount() < 1 {
		logger.ERR("resp has no record")
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}

	if res.idx >= int32(data.ElementNum) {
		logger.ERR("resp fetch record over, current idx: %d, ", res.idx)
		return nil, &terror.ErrorCode{Code: terror.API_ERR_NO_MORE_RECORD}
	}

	rec := &record.Record{
		AppId:       uint64(res.pkg.Head.RouterInfo.AppID),
		ZoneId:      uint32(res.pkg.Head.RouterInfo.ZoneID),
		TableName:   string(res.pkg.Head.RouterInfo.TableName[0:res.pkg.Head.RouterInfo.TableNameLen]),
		Cmd:         int(res.pkg.Head.Cmd),
		KeyMap:      make(map[string][]byte),
		ValueMap:    make(map[string][]byte),
		Version:     -1,
		KeySet:      res.pkg.Head.KeyInfo,
		ValueSet:    nil,
		UpdFieldSet: nil,
	}

	//unpack
	if err := rec.UnPackKey(); err != nil {
		logger.ERR("record unpack key failed, app %d zone %d table %s ,err %s",
			rec.AppId, rec.ZoneId, rec.TableName, err.Error())
		return nil, err
	}

	read_bytes := uint32(0)
	err := unpack_element_buff(data.ElementsBuff, uint32(res.offset), data.ElementsBuffLen, &rec.Index,
		&read_bytes, rec.ValueMap)
	res.idx += 1
	res.offset += int32(read_bytes)
	res.record = rec
	return rec, err
}

func (res *listReplaceResponse) GetUserBuffer() []byte {
	return res.pkg.Head.UserBuff
}

func (res *listReplaceResponse) GetSeq() int32 {
	return res.pkg.Head.Seq
}

func (res *listReplaceResponse) HaveMoreResPkgs() int {
	return 0
}

func (res *listReplaceResponse) GetRecordMatchCount() int {
	return terror.API_ERR_OPERATION_TYPE_NOT_MATCH
}
