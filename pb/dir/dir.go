package dir

import (
	"encoding/binary"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/common"
	log "github.com/tencentyun/tcaplusdb-go-sdk/pb/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcapdir_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/version"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/tnet"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	NotSignUp     = 0
	SignUpIng     = 1
	SignUpSuccess = 2
	SignUpFail    = 3
)

// dir服务管理
type DirServer struct {
	appId     uint64
	zoneList  []uint32
	url       string
	signature string

	MsgPipe chan *tcapdir_protocol_cs.TCapdirCSPkg //dir消息通道
	conn    *tnet.Conn

	oldDirIndex uint32 //随机选中的dir
	curDirIndex uint32 //当前选中的dir
	urlList     []string

	signUpFlag        uint32
	signUpTime        time.Time
	AllDirConnectFail bool

	//心跳时间间隔s, 10s
	heartbeatInterval time.Duration
	lastHeartbeatTime time.Time
	error error
}

func (dir *DirServer) Init(appId uint64, zoneList []uint32, dirUrl string, signature string) error {
	if len(zoneList) > int(tcapdir_protocol_cs.TCAPDIR_MAX_TABLE) {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "zoneList large 256"}
	}
	dir.heartbeatInterval = 10
	dir.lastHeartbeatTime = time.Now()
	dir.appId = appId
	dir.zoneList = make([]uint32, len(zoneList))
	copy(dir.zoneList, zoneList)
	dir.url = dirUrl
	dir.signature = signature
	dir.conn = nil
	dir.signUpFlag = NotSignUp
	dir.AllDirConnectFail = false
	dir.MsgPipe = make(chan *tcapdir_protocol_cs.TCapdirCSPkg, 1)
	if err := dir.domainConvert(); err != nil {
		return err
	}

	//随机选个连接
	rand.Seed(time.Now().UnixNano())
	dir.curDirIndex = uint32(rand.Intn(len(dir.urlList)))
	dir.oldDirIndex = dir.curDirIndex
	return nil
}

func (dir *DirServer) SendRequest(buf []byte) error {
	if dir.conn != nil {
		return dir.conn.Send(buf)
	}
	log.ERR("dir.conn is nil")
	return &terror.ErrorCode{Code: terror.SendRequestFail, Message: "dir not connected"}
}

func (dir *DirServer) Update() {
	dir.connect()
	if nil == dir.conn {
		return
	}

	if dir.conn.GetStat() == tnet.Connected && dir.signUpFlag == SignUpSuccess {
		//发送心跳
		diff := time.Now().Sub(dir.lastHeartbeatTime)
		if diff > dir.heartbeatInterval*time.Second {
			dir.SendHeartbeat()
			dir.lastHeartbeatTime = time.Now()
		}
	}
}

func (dir *DirServer) DisConnect() {
	dir.signUpFlag = NotSignUp
	if dir.conn != nil {
		dir.conn.Close()
		dir.conn = nil
	}
}

//初始化超时，获取dir的error信息
func (dir *DirServer) GetError() error {
	return  dir.error
}

//从列表中选择一个连接dir
func (dir *DirServer) connect() error {
	if dir.conn == nil {
		//连接dir, 3s超时
		for i := 0; i < len(dir.urlList); i++ {
			var err error
			dir.conn, err = tnet.NewConn(dir.urlList[dir.curDirIndex], 3*time.Second, ParseDirPkgLen,
				DirCallBackFunc, dir, 0)
			if err == nil {
				break
			}
			dir.error = &terror.ErrorCode{Code: terror.API_ERR_DIR_CONNECT_FAILED, Message: "connect dir failed:" + dir.urlList[dir.curDirIndex]}
			dir.curDirIndex++
			dir.curDirIndex = dir.curDirIndex % uint32(len(dir.urlList))
			if dir.curDirIndex == dir.oldDirIndex {
				dir.AllDirConnectFail = true
				log.WARN("AllDirConnectFail dirList:%v", dir.urlList)
			}
			log.ERR("new conn failed %v", err)
		}
	} else {
		if dir.conn.GetStat() == tnet.Connected {
			//连接成功,认证
			if dir.signUpFlag == NotSignUp {
				dir.signUpFlag = SignUpIng
				log.INFO("start sign up dir %v", dir.urlList[dir.curDirIndex])
				dir.signUp()
			} else if dir.signUpFlag != SignUpSuccess && time.Now().Sub(dir.signUpTime).Seconds() > 3 {
				//认证超时
				log.ERR("sign up dir %v timeout(3s), conn stat %v", dir.urlList[dir.curDirIndex], dir.conn.GetStat())
				dir.error = &terror.ErrorCode{Code: terror.API_ERR_DIR_SIGNUP_FAILED, Message: "sign up dir timeout(3s):" + dir.urlList[dir.curDirIndex]}
				dir.DisConnect()
				dir.curDirIndex++
				dir.curDirIndex = dir.curDirIndex % uint32(len(dir.urlList))
				if dir.curDirIndex == dir.oldDirIndex {
					dir.AllDirConnectFail = true
					log.WARN("AllDirConnectFail dirList:%v", dir.urlList)
				}
			}
			return nil
		} else if dir.conn.GetStat() == tnet.Connecting {
			//连接中
			return nil
		} else {
			dir.error = &terror.ErrorCode{Code: terror.API_ERR_DIR_CONNECT_FAILED, Message: "connect dir failed:" + dir.urlList[dir.curDirIndex]}
			//连接失败，重连下个连接
			log.ERR("connect dir %v failed, conn stat %v",
				dir.urlList[dir.curDirIndex%uint32(len(dir.urlList))], dir.conn.GetStat())
			dir.DisConnect()
			dir.curDirIndex++
			dir.curDirIndex = dir.curDirIndex % uint32(len(dir.urlList))
			if dir.curDirIndex == dir.oldDirIndex {
				dir.AllDirConnectFail = true
				log.WARN("AllDirConnectFail dirList:%v", dir.urlList)
			}
			return nil
		}
	}
	return nil
}

func (dir *DirServer) ProcessSignUpRes(res int) {
	if 0 == res {
		dir.signUpFlag = SignUpSuccess
		log.INFO("dir %v signUp success", dir.urlList[dir.curDirIndex%uint32(len(dir.urlList))])
	} else {
		dir.signUpFlag = SignUpFail
		log.ERR("dir %v signUp failed, %d", dir.urlList[dir.curDirIndex%uint32(len(dir.urlList))], res)
	}
}

//判断是否收到完整dir包
//TcapdirCSHead = Magic(2) + Cmd(2) + Version(2) + HeadLen(2) + BodyLen(4) + AppID(8) = 20
func ParseDirPkgLen(buf []byte) int {
	if len(buf) >= 20 {
		headLen := binary.BigEndian.Uint16(buf[6:8])
		bodyLen := binary.BigEndian.Uint32(buf[8:12])
		return int(headLen) + int(bodyLen)
	}
	return 0
}

/*
@brief tcp回调函数，解析消息，并触发client协程进行处理
@param url 触发响应包的url
@param  buf 响应包
@param cbPara 回调参数，此处为*DirServer
@retVal error
*/
func DirCallBackFunc(url *string, pkg *tnet.PKG) error {
	buf := pkg.GetData()
	dir, ok := pkg.GetCbPara().(*DirServer)
	if !ok {
		log.ERR("RecvCallBackFunc cbPara type invalid")
		return nil
	}

	resp := tcapdir_protocol_cs.NewTCapdirCSPkg()
	err := resp.Unpack(tcapdir_protocol_cs.TCapdirCSPkgCurrentVersion, buf)
	// 之后这个pkg不会再被用到，回收到对象池中
	pkg.Done()
	if err != nil {
		log.ERR("Unpack dir msg failed, url %v err %v", *url, err.Error())
		return err
	}

	log.DEBUG("recv dir msg %+v", *resp.Head)
	dir.MsgPipe <- resp
	return nil
}

func (dir *DirServer) domainConvert() error {

	if len(dir.url) <= 0 {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dirUrl is invalid"}
	}

	//解析url tcp://a.com:9999;tcp://b.com:9999
	list := strings.Split(dir.url, ";")
	if len(list) < 1 {
		return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dirUrl is invalid, use ';' split"}
	}

	//解析域名
	urlMap := map[string]byte{}
	for _, url := range list {
		network, domain, port, err := tnet.ParseUrl(&url)
		if err != nil {
			return err
		}

		if network != "tcp" {
			return &terror.ErrorCode{Code: terror.ParameterInvalid, Message: "dirUrl network must be tcp"}
		}

		ipList, err := net.LookupHost(domain)
		if err != nil {
			return err
		}

		for _, ip := range ipList {
			newUrl := network + "://" + ip + ":" + port
			if _, ok := urlMap[newUrl]; !ok {
				//not exist
				urlMap[newUrl] = 0
				dir.urlList = append(dir.urlList, newUrl)
			}
		}
	}

	return nil
}

func (dir *DirServer) signUp() error {
	dir.signUpTime = time.Now()
	req := tcapdir_protocol_cs.NewTCapdirCSPkg()
	//head
	req.Head.Magic = uint16(tcapdir_protocol_cs.TCAPLUS_PROTOCOL_MAGIC_DIR_CS)
	req.Head.Cmd = uint16(tcapdir_protocol_cs.TCAPDIR_CS_CMD_SIGNUP_REQ)
	req.Head.Version = 0
	req.Head.HeadLen = 0
	req.Head.BodyLen = 0
	req.Head.AppID = dir.appId
	req.Body.Init(int64(req.Head.Cmd))

	//body
	req.Body.ReqSignUpApp.Signature = dir.signature
	req.Body.ReqSignUpApp.Type = 0
	req.Body.ReqSignUpApp.TableCount = int16(len(dir.zoneList))
	req.Body.ReqSignUpApp.TableList = make([]*tcapdir_protocol_cs.TableInfo, len(dir.zoneList))
	for i := 0; i < len(dir.zoneList); i++ {
		tableInfo := tcapdir_protocol_cs.NewTableInfo()
		tableInfo.ZoneID = int32(dir.zoneList[i])
		tableInfo.Name = ""
		req.Body.ReqSignUpApp.TableList[i] = tableInfo
	}

	req.Body.ReqSignUpApp.ClientInfo.ApiVersion = version.MAJOR
	req.Body.ReqSignUpApp.ClientInfo.DetailVer = version.GetModuleName()
	req.Body.ReqSignUpApp.ClientInfo.Version = version.Version
	req.Body.ReqSignUpApp.ClientInfo.GitSHA1 = version.GitCommitId
	req.Body.ReqSignUpApp.ClientInfo.GitBranch = version.GitBranch
	req.Body.ReqSignUpApp.ClientInfo.Platform = int16(tcapdir_protocol_cs.TCAPDIR_PLATFORM_LINUX64)
	req.Body.ReqSignUpApp.ClientInfo.TableCount = 0
	req.Body.ReqSignUpApp.ClientInfo.TraitBits = 0
	req.Body.ReqSignUpApp.ClientInfo.HostTime = uint64(time.Now().Unix())

	//pack
	if buf, err := req.Pack(tcapdir_protocol_cs.TCapdirCSPkgCurrentVersion); err != nil {
		log.ERR("signUp pack failed %v", err.Error())
		return err
	} else {
		/*
			resp := tcapdir_protocol_cs.NewTCapdirCSPkg()
			if err := resp.Unpack(tcapdir_protocol_cs.TCapdirCSPkgCurrentVersion, buf); err != nil {
				log.ERR("Unpack dir msg failed, err %v", err.Error())
			}*/
		log.DEBUG("msg:%+v signUp pack len %v", *req.Head, len(buf))
		dir.SendRequest(buf)
	}
	log.DEBUG("signUp send ")
	return nil
}

func (dir *DirServer) GetDirList() error {
	if dir.signUpFlag != SignUpSuccess {
		log.INFO("dir not signUp, %v ", dir.signUpFlag)
		return nil
	}

	if nil == dir.conn {
		log.INFO("dir conn is nil")
		return nil
	}

	if dir.conn.GetStat() != tnet.Connected {
		log.INFO("dir not connected")
		return nil
	}

	req := tcapdir_protocol_cs.NewTCapdirCSPkg()
	//head
	req.Head.Magic = uint16(tcapdir_protocol_cs.TCAPLUS_PROTOCOL_MAGIC_DIR_CS)
	req.Head.Cmd = uint16(tcapdir_protocol_cs.TCAPDIR_CS_CMD_GET_DIR_SERVER_LIST_REQ)
	req.Head.Version = 0
	req.Head.HeadLen = 0
	req.Head.BodyLen = 0
	req.Head.AppID = dir.appId

	//pack
	if buf, err := req.Pack(tcapdir_protocol_cs.TCapdirCSPkgCurrentVersion); err != nil {
		log.ERR("GetDirList pack failed %v", err.Error())
		return err
	} else {
		log.DEBUG("msg:%+v GetDirList pack len %v", *req.Head, len(buf))
		dir.SendRequest(buf)
	}
	log.DEBUG("GetDirList send ")
	return nil
}

func (dir *DirServer) GetAccessProxy() error {
	if dir.signUpFlag != SignUpSuccess {
		log.INFO("GetAccessProxy dir not signUp, %v ", dir.signUpFlag)
		return nil
	}

	if nil == dir.conn {
		log.INFO("dir conn is nil")
		return nil
	}

	if dir.conn.GetStat() != tnet.Connected {
		log.INFO("GetAccessProxy dir not connected")
		return nil
	}

	for i := 0; i < len(dir.zoneList); i++ {
		req := tcapdir_protocol_cs.NewTCapdirCSPkg()
		//head
		req.Head.Magic = uint16(tcapdir_protocol_cs.TCAPLUS_PROTOCOL_MAGIC_DIR_CS)
		req.Head.Cmd = uint16(tcapdir_protocol_cs.TCAPDIR_CS_CMD_GET_TABLES_AND_ACCESS_REQ)
		req.Head.Version = 0
		req.Head.HeadLen = 0
		req.Head.BodyLen = 0
		req.Head.AppID = dir.appId
		req.Body.Init(int64(req.Head.Cmd))

		//body
		req.Body.ReqGetTablesAndAccess.ZoneID = int32(dir.zoneList[i])
		req.Body.ReqGetTablesAndAccess.Signature = dir.signature
		req.Body.ReqGetTablesAndAccess.Version = version.Version

		//pack
		if buf, err := req.Pack(tcapdir_protocol_cs.TCapdirCSPkgCurrentVersion); err != nil {
			log.ERR("GetAccessProxy pack failed %v", err.Error())
			return err
		} else {
			log.DEBUG("msg:%+v GetAccessProxy pack len %v", *req.Head, len(buf))
			dir.SendRequest(buf)
		}
		log.DEBUG("GetAccessProxy send, zone %v", dir.zoneList[i])
	}

	return nil
}

//根据dir的响应消息判断dir列表是否发生变化
func (dir *DirServer) ProcessDirListRes(res *tcapdir_protocol_cs.ResGetDirServerList) {
	if res.Result != 0 {
		log.ERR("DirListRes err %d", res.Result)
		return
	}

	if res.DirServerCount <= 0 {
		log.INFO("DirListRes DirServerCount invalid %d", res.DirServerCount)
		return
	}

	if common.PublicIP != "" {
		for index := range res.DirServer[0:res.DirServerCount] {
			urlNet, _, urlPort, err := tnet.ParseUrl(&res.DirServer[index])
			if err != nil {
				log.ERR("proxy url is invalid %s", res.DirServer[index])
				return
			}
			// 变更IP
			res.DirServer[index] = fmt.Sprintf("%s://%s:%s", urlNet, common.PublicIP, urlPort)
		}
	}

	//比较数量
	if len(dir.urlList) != int(res.DirServerCount) {
		//重新连接
		log.INFO("dirList changed old %v new %v", dir.urlList, res.DirServer[0:res.DirServerCount])
		dir.urlList = make([]string, res.DirServerCount)
		copy(dir.urlList, res.DirServer[0:res.DirServerCount])
		dir.curDirIndex = uint32(rand.Intn(len(dir.urlList)))
		dir.oldDirIndex = dir.curDirIndex
		dir.DisConnect()
		dir.connect()
		return
	}

	//比较内容
	oldUrlMap := map[string]byte{}
	for _, v := range dir.urlList {
		oldUrlMap[v] = 0
	}

	for i := 0; i < int(res.DirServerCount); i++ {
		if _, exist := oldUrlMap[res.DirServer[i]]; !exist {
			//重新连接
			log.INFO("dirList changed old %v new %v", dir.urlList, res.DirServer[0:res.DirServerCount])
			dir.urlList = make([]string, res.DirServerCount)
			copy(dir.urlList, res.DirServer[0:res.DirServerCount])
			dir.curDirIndex = uint32(rand.Intn(len(dir.urlList)))
			dir.oldDirIndex = dir.curDirIndex
			dir.DisConnect()
			dir.connect()
			return
		}
	}
	log.DEBUG("dirList not changed")
}

func (dir *DirServer) SetHeartbeatInterval(heartbeatInterval time.Duration) {
	if heartbeatInterval != dir.heartbeatInterval && heartbeatInterval > 0 {
		dir.heartbeatInterval = heartbeatInterval
	}
}

func (dir *DirServer) SendHeartbeat() {
	req := tcapdir_protocol_cs.NewTCapdirCSPkg()
	//head
	req.Head.Magic = uint16(tcapdir_protocol_cs.TCAPLUS_PROTOCOL_MAGIC_DIR_CS)
	req.Head.Cmd = uint16(tcapdir_protocol_cs.TCAPDIR_CS_CMD_HEARTBEAT_REQ)
	req.Head.Version = 0
	req.Head.HeadLen = 0
	req.Head.BodyLen = 0
	req.Head.AppID = dir.appId
	req.Body.Init(int64(req.Head.Cmd))

	req.Body.ReqHeartBeat.HostTime = uint64(time.Now().Unix())
	req.Body.ReqHeartBeat.WithQos = 0

	//pack
	if buf, err := req.Pack(tcapdir_protocol_cs.TCapdirCSPkgCurrentVersion); err != nil {
		log.ERR("SendHeartbeat pack failed %v", err.Error())
		return
	} else {
		log.DEBUG("msg:%+v SendHeartbeat pack len %v", *req.Head, len(buf))
		dir.SendRequest(buf)
	}
	log.DEBUG("SendHeartbeat send, dir %v", dir.url)
}
