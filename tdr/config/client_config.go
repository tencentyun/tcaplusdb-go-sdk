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

type ClientOption struct {
	ProxyConnOption    ProxyConnOption
	PackRoutineCount   int
	UnPackRoutineCount int
}

type ClientCtrl struct {
	Option *ClientOption
	sync.WaitGroup
}

func NewDefaultClientOption() *ClientOption {
	cpuNum := runtime.NumCPU()
	return &ClientOption{
		PackRoutineCount:   cpuNum,
		UnPackRoutineCount: cpuNum,
		ProxyConnOption: ProxyConnOption{
			BufSizePerCon: 10 * 1024 * 1024, //设置读写缓冲区10MB
			ProxyMaxCount: 200,              //限制单个zone的proxy连接数
			ConTimeout:    15 * time.Second, //15秒读写包失败，将连接置位不可用
		},
	}
}
