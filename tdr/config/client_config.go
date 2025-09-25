package config

import (
	"runtime"
	"sync"
	"time"
)

type ProxyConnOption struct {
	BufSizePerCon int
	ProxyMaxCount int
	ConTimeout    time.Duration
}

type SubscribeOption struct {
	ReSubscribeIntervalTime time.Duration //定时扫描需要续订的topic时间间隔，默认50ms
}

type ClientOption struct {
	ProxyConnOption        ProxyConnOption
	PackRoutineCount       int
	UnPackRoutineCount     int
	PbMetaInitRoutineCount int
	SwitchProxyPercent     int
	SubscribeOption        SubscribeOption
}

type ClientCtrl struct {
	Option *ClientOption
	sync.WaitGroup
}

func NewDefaultClientOption() *ClientOption {
	cpuNum := runtime.NumCPU()
	return &ClientOption{
		PackRoutineCount:       cpuNum, // 运行过程中，打包协程数
		UnPackRoutineCount:     cpuNum, // 运行过程中，解包协程数
		PbMetaInitRoutineCount: cpuNum, // pb初始化时，拉取meta的协程数控制
		SwitchProxyPercent:     10,     // 部分proxy算法，每次proxy列表切换，按百分比例切换，防止内存抖动过大
		ProxyConnOption: ProxyConnOption{
			BufSizePerCon: 256 * 1024,       //设置读写缓冲区256kb
			ProxyMaxCount: 200,              //限制单个zone的proxy连接数
			ConTimeout:    15 * time.Second, //15秒读写包失败，将连接置位不可用
		},
		SubscribeOption: SubscribeOption{
			ReSubscribeIntervalTime: 50 * time.Millisecond,
		},
	}
}
