package version

import (
	"os"
	"strconv"
	"strings"
)

const (
	MAJOR       = 3
	MINOR       = 55
	REV         = 0
	GitBranch   = "TcaplusGoApi3.55.0"
	GitCommitId = "v0.6.30"                      //每次tag，必须更新
	Version     = "3.55.0.000029.GoApi_20250922" //每次tag，必须更新版本和时间
)

func GetModuleName() string {
	procName := os.Getenv("_")
	if len(procName) == 0 {
		procName = "unknownGo." + strconv.Itoa(os.Getpid())
	} else {
		tmpSlice := strings.Split(procName, "/")
		size := len(tmpSlice)
		if size > 0 {
			if len(tmpSlice[size-1]) > 0 {
				procName = tmpSlice[size-1]
			}
		}
	}
	if hostname, err := os.Hostname(); err == nil {
		procName = hostname + ":" + procName
	}
	// max len 64
	if len(procName) > 63 {
		procName = procName[:63]
	}
	return procName
}
