package tcaplus

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"
	"time"
)

type Client struct {
	*client
}

// 兼容老接口，保留
func NewClient() *Client {
	c := new(Client)
	c.client = newClient(false)
	c.defTimeout = 5 * time.Second
	return c
}

func NewTDRClient() *Client {
	c := new(Client)
	c.client = newClient(false)
	c.defTimeout = 5 * time.Second
	return c
}

// CheckTdrMetaVersion 对比本地和svr端的tdr表的version兼容性，防止忘记在svr端改表，导致数据写乱
// tdrVersion取值为表的go文件中的开头常量部分，类似：表名+CurrentVersion
func (c *Client) CheckTdrMetaVersion(zone uint32, table string, tdrVersion uint32) error {
	req, err := c.NewRequest(zone, table, cmd.TcaplusApiMetadataGetReq)
	if err != nil {
		terr := terror.MakeError(terror.ParameterInvalid,
			fmt.Sprintf("zone %d table %s NewRequest error:%s", zone, table, err.Error()))
		logger.ERR(terr.Error())
		return terr
	}
	reqPkg := req.GetTcaplusPackagePtr()
	reqPkg.Body.MetadataGetReq.MetadataVersion = tdrVersion

	resp, err := c.Do(req, c.defTimeout)
	if err != nil {
		terr := terror.MakeError(terror.ParameterInvalid,
			fmt.Sprintf("zone %d table %s get pb meta Do request error:%s", zone, table, err))
		logger.ERR(terr.Error())
		return terr
	}
	if r := resp.GetResult(); r != 0 {
		terr := terror.MakeError(terror.ParameterInvalid,
			fmt.Sprintf("get zone %d table %s metadata error:%s", zone, table, terror.GetErrMsg(r)))
		logger.ERR(terr.Error())
		return terr
	}
	if resp.GetTcaplusPackagePtr() == nil {
		terr := terror.MakeError(terror.ParameterInvalid,
			fmt.Sprintf("get zone %d table %s metadata error:response pkg is nil", zone, table))
		logger.ERR(terr.Error())
		return terr
	}

	metaRes := resp.GetTcaplusPackagePtr().Body.MetadataGetRes
	if metaRes.IdlType == 2 {
		terr := terror.MakeError(terror.ParameterInvalid,
			fmt.Sprintf("get zone %d table %s metadata error:table type %d is proto",
				zone, table, metaRes.IdlType))
		logger.ERR(terr.Error())
		return terr
	}

	// 本地的version > 服务端的version
	if uint32(metaRes.TdrMetaCurrentVersion) < tdrVersion {
		terr := terror.MakeError(terror.ParameterInvalid,
			fmt.Sprintf("zone %d table %s meta local version %d > svr version %d, change table on OMS first",
				zone, table, tdrVersion, metaRes.TdrMetaCurrentVersion))
		logger.ERR(terr.Error())
		return terr
	}

	logger.DEBUG("zone %d table %s meta check success local version %d, svr version %d",
		zone, table, tdrVersion, metaRes.TdrMetaCurrentVersion)
	return nil
}
