package main

import (
	"flag"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/tdragent/src/cfg"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/tdragent/src/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/tdragent/src/runTest"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////GO TDR AGENT////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	casePath     = flag.String("case", "../case/case.xml", "case file path")
	apiLogConf   = flag.String("apiLog", "../cfg/apiLogConf.xml", "tcaplus api log config xml")
	agentLogConf = flag.String("agentLog", "", "tcaplus api log config xml")
	help         = flag.Bool("h", false, "help")
)

func main() {
	//Parse Flag
	flag.Parse()
	if *help {
		fmt.Println("\n./tdrAgent -case ../case/case.xml -apiLog ../cfg/apiLogConf.xml")
		flag.Usage()
		return
	}

	//初始化日志
	if err := logger.Init(*agentLogConf); err != nil {
		fmt.Println("logger.Init " + err.Error())
		return
	}

	//解析case配置
	tcaplusCase, err := cfg.ParseCase(*casePath)
	if err != nil {
		fmt.Println("cfg.ParseCase " + err.Error())
		return
	}

	//初始化Tcaplus go api
	client := tcaplus.NewClient()
	if err := client.SetLogCfg(*apiLogConf); err != nil {
		fmt.Println("client.SetLogCfg " + err.Error())
		return
	}

	if err := client.Dial(uint64(tcaplusCase.Head.AppID), []uint32{uint32(tcaplusCase.Head.ZoneID)},
		tcaplusCase.Head.DirUrl, tcaplusCase.Head.AppSignUp, 60); err != nil {
		fmt.Println("client.Dial " + err.Error())
		return
	}

	fmt.Println("client dial success")

	//开始主循环
	testCase := &runTest.RunTest{
		TcaplusClient: client,
		TcaplusCase:   tcaplusCase,
	}
	testCase.Run()
}
