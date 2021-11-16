package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/cfg"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestClose(t *testing.T) {
	if err := cfg.ReadApiCfg("../cfg/api_cfg.xml"); err != nil {
		t.Errorf("ReadApiCfg fail %s", err.Error())
		return
	}

	client := tcaplus.NewClient()
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		t.Errorf("excepted SetLogCfg success")
		return
	}

	err := client.Dial(cfg.ApiConfig.AppId, []uint32{cfg.ApiConfig.ZoneId}, cfg.ApiConfig.DirUrl, cfg.ApiConfig.Signature, 30)
	if err != nil {
		t.Errorf("excepted dial success, %s", err.Error())
		return
	}

	buff := make([]byte, 10240)
	runtime.Stack(buff, true)
	// 启动1个client，将会拉起4个协程处理用户的同步请求和1个控制网络包的协程
	// 连接1个proxy将会产生1个发送请求的协程与1个接收响应的协程以及4个解包cs协议的协程
	// 连接1个dir将会产生1个发送请求的协程与1个接收响应的协程

	// 假设client数为i  proxy数为j  dir数为k
	// tnet.(*Conn).process = i * (j + k)
	fmt.Println("tnet.(*Conn).process", strings.Count(string(buff), "tnet.(*Conn).process"))
	// tnet.(*Conn).mergeSend = i * (j + k)
	fmt.Println("tnet.(*Conn).mergeSend", strings.Count(string(buff), "tnet.(*Conn).mergeSend"))
	// router.(*Router).processSyncOperate.func1 = 4 * i
	fmt.Println("router.(*Router).processSyncOperate.func1",
		strings.Count(string(buff), "router.(*Router).processSyncOperate.func1"))
	// (*netServer).netPkgProcess = i
	fmt.Println("(*netServer).netPkgProcess", strings.Count(string(buff), "(*netServer).netPkgProcess"))
	// router.(*server).initRecv.func1 = 4 * i * j
	fmt.Println("router.(*server).initRecv.func1",
		strings.Count(string(buff), "router.(*server).initRecv.func1"))

	// close之后client就不能再重用了，需要重新调用tcaplus.NewClient()生成
	client.Close()

	// 注：关闭为异步关闭，不保证调用close之后所有协程就立马停止
	time.Sleep(time.Second)

	buff = make([]byte, 10240)
	runtime.Stack(buff, true)
	if strings.Count(string(buff), "tnet.(*Conn).process") != 0 {
		t.Errorf("tnet.(*Conn).process != 0")
		return
	}
	if strings.Count(string(buff), "tnet.(*Conn).mergeSend") != 0 {
		t.Errorf("tnet.(*Conn).mergeSend != 0")
		return
	}
	if strings.Count(string(buff), "router.(*Router).processSyncOperate.func1") != 0 {
		t.Errorf("router.(*Router).processSyncOperate.func1 != 0")
		return
	}
	if strings.Count(string(buff), "(*netServer).netPkgProcess") != 0 {
		t.Errorf("(*netServer).netPkgProcess != 0")
		return
	}
	if strings.Count(string(buff), "router.(*server).initRecv.func1") != 0 {
		t.Errorf("router.(*server).initRecv.func1 != 0")
		return
	}
}
