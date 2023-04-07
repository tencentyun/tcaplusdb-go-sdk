package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/PB/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"strings"
	"sync"
	"time"
)

//TcaplusDB GO PB API的连接参数
const (
	//集群访问地址，本地docker版：配置docker部署机器IP, 端口默认:9999, 腾讯云线上环境配置为连接地址IP和端口
	DirUrl = "tcp://x.x.x.x:9999"
	//集群接入ID, 默认为3，本地docker版：直接填3，云上版本：根据实际集群接入ID填写
	AppId = 2
	//集群访问密码，本地docker版：登录tcaplusdb web运维平台查看(账号/密码:tcaplus/tcaplus)，业务管理->业务维护->查看pb_app业务对应密码; 云上版本：根据实际申请集群详情页查看
	Signature = "xxxxxxxxxx"
	//表格组ID，替换为自己创建的表格组ID
	Zone = 3
	//表名称
	TableName = "game_players"
)

//初始化客户端连接
func NewPBClient() *tcaplus.PBClient {
	//通过指定接入ID(AppId), 表格组id表表(zoneList), 接入地址(DirUrl), 集群密码(Signature) 参数创建TcaplusClient的对象client
	//通过client对象可以访问集群下的所有大区和表
	//创建表格、获取访问点信息的指引请参考 https://cloud.tencent.com/document/product/596/38807
	client := tcaplus.NewPBClient()
	zoneList := []uint32{Zone}
	zoneTable := make(map[uint32][]string)
	//构造Map对象存储对应表格组下所有的表
	zoneTable[Zone] = []string{TableName}
	//建立到对应集群的连接客户端
	err := client.Dial(AppId, zoneList, DirUrl, Signature, 30, zoneTable)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	client.SetDefaultZoneId(Zone)
	return client
}

type clieninf interface {
	RecvResponse() (response.TcaplusResponse, error)
}

//同步接收
func RecvResponse(client clieninf) (response.TcaplusResponse, error) {
	//recv response
	timeOutChan := time.After(5 * time.Second)
	for {
		select {
		case <-timeOutChan:
			return nil, errors.New("5s timeout")
		default:
			resp, err := client.RecvResponse()
			if err != nil {
				return nil, err
			} else if resp == nil {
				time.Sleep(time.Microsecond * 1)
			} else {
				return resp, nil
			}
		}
	}
}

var pbclient *tcaplus.PBClient
var once sync.Once
var ZoneId uint32

func InitPBSyncClient() *tcaplus.PBClient {
	var err error
	once.Do(func() {
		err = cfg.ReadApiCfg("../cfg/api_cfg.xml")
		if err != nil {
			fmt.Printf("ReadApiCfg fail %s", err.Error())
			return
		}

		pbclient = tcaplus.NewPBClient()
		err = pbclient.SetLogCfg("../cfg/logconf.xml")
		if err != nil {
			fmt.Printf("excepted SetLogCfg success")
			return
		}

		ZoneId = cfg.ApiConfig.ZoneId

		tables := strings.Split(cfg.ApiConfig.Table, ",")
		zoneTable := map[uint32][]string{cfg.ApiConfig.ZoneId: tables}
		err = pbclient.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30, zoneTable)
		if err != nil {
			fmt.Printf("excepted dial success, %s", err.Error())
			return
		}
	})
	if err != nil {
		return nil
	}
	return pbclient
}

func ConvertToJson(v interface{}) string {
	body, _ := json.Marshal(v)
	buf := &bytes.Buffer{}
	json.Indent(buf, body, "", "\t")
	return buf.String()
}
