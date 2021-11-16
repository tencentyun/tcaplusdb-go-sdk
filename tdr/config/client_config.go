package config

import (
	"runtime"
	"sync"
)

type ProxyConnOption struct {
	BufSizePerCon int
	ProxyMaxCount int
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
			BufSizePerCon: 10 * 1024 * 1024, //10MB
			ProxyMaxCount: 200,
		},
	}
}
