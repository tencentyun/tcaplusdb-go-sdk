package cs_pool

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"sync"
)

//cs协议缓存池
var tcaplusCSPkgPool = sync.Pool{
	New: func() interface{} {
		return tcaplus_protocol_cs.NewTCaplusPkg()
	},
}

func GetTcaplusCSPkg(cmd uint32) *tcaplus_protocol_cs.TCaplusPkg {
	pkg := tcaplusCSPkgPool.Get().(*tcaplus_protocol_cs.TCaplusPkg)
	pkg.Head.SubCmd = 0
	pkg.Head.Flags = 0
	pkg.Head.RouterInfo.TableNameLen = 0
	pkg.Head.UserBuffLen = 0
	pkg.Head.PerfTestLen = 0
	pkg.Head.KeyInfo.Version = 0
	pkg.Head.KeyInfo.FieldNum = 0
	pkg.Head.ReqBodyCompressType = 0
	pkg.Head.RespBodyCompressType = 0
	pkg.Head.Result = 0
	pkg.Head.SplitTableKeyBuffLen = 0
	if pkg.Head.Cmd == cmd {
		return pkg
	}
	pkg.Head.Cmd = cmd
	pkg.Body.Init(int64(cmd))
	return pkg
}

func PutTcaplusCSPkg(pkg *tcaplus_protocol_cs.TCaplusPkg) {
	tcaplusCSPkgPool.Put(pkg)
}

var tcaplusCSResPkgPool = sync.Pool{
	New: func() interface{} {
		return tcaplus_protocol_cs.NewTCaplusPkg()
	},
}

func GetTcaplusCSResPkg() *tcaplus_protocol_cs.TCaplusPkg {
	pkg := tcaplusCSResPkgPool.Get().(*tcaplus_protocol_cs.TCaplusPkg)
	pkg.Head.Init()
	pkg.Body = nil
	return pkg
}

func PutTcaplusCSResPkg(pkg *tcaplus_protocol_cs.TCaplusPkg) {
	tcaplusCSResPkgPool.Put(pkg)
}
