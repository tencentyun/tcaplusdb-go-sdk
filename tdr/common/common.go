package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"net"
	"sync"
	"time"
)

const (
	// 每个server处理响应的协程数
	ConfigProcRespRoutineNum = 4

	// 管道深度大一些可以防止瞬间并发太多请求导致管道满了而阻塞住

	// 处理响应的协程管道深度
	ConfigProcRespDepth = 10000

	// 处理写请求的协程管道深度，用于合并写请求
	ConfigProcReqDepth = 10000

	// 处理用户同步请求的协程数
	ConfigProcRouterRoutineNum = 4

	// 处理用户同步请求的协程管道深度
	ConfigProcRouterDepth = 10000

	// 收包的最大包大小
	PkgBufferMaxLength = 10 * 1024 * 1024

	// 读写网络报文的超时时间
	ConfigReadWriteTimeOut = 30
)

var PublicIP string

var TimeNow time.Time = time.Now()

var bytesPool = sync.Pool{
	New: func() interface{} {
		return &PKGManager{data: make([]byte, PkgBufferMaxLength)}
	},
}

var pkgBufferPool = sync.Pool{
	New: func() interface{} {
		return &PKGBuffer{}
	},
}

func GetPKGManager(pkg *PKGManager) *PKGManager {
	newPkg := bytesPool.Get().(*PKGManager)
	if pkg != nil && pkg.parseOffset < pkg.readOffset {
		newPkg.readOffset = copy(newPkg.data, pkg.data[pkg.parseOffset:pkg.readOffset])
		go putPKGManager(pkg)
	}
	return newPkg
}

func putPKGManager(pkg *PKGManager) {
	// 等待所有buffer都失效再回收
	pkg.wg.Wait()
	pkg.readOffset = 0
	pkg.parseOffset = 0
	bytesPool.Put(pkg)
}

type PKGManager struct {
	data        []byte
	wg          sync.WaitGroup
	readOffset  int
	parseOffset int
}

func (m *PKGManager) Read(conn net.Conn) (int, error) {
	conn.SetReadDeadline(TimeNow.Add(ConfigReadWriteTimeOut * time.Second))
	n, err := conn.Read(m.data[m.readOffset:])
	m.readOffset += n
	return n, err
}

func (m *PKGManager) ValidBuffer() []byte {
	return m.data[m.parseOffset:m.readOffset]
}

func (m *PKGManager) ValidLength() int {
	return m.readOffset - m.parseOffset
}

func (m *PKGManager) GetPkgBuffer(size int) *PKGBuffer {
	pkgBuffer := pkgBufferPool.Get().(*PKGBuffer)
	pkgBuffer.wg = &m.wg
	pkgBuffer.data = m.data[m.parseOffset : m.parseOffset+size]
	m.parseOffset += size
	m.wg.Add(1)
	return pkgBuffer
}

func (m *PKGManager) BufferIsFull() bool {
	return m.readOffset == PkgBufferMaxLength
}

type PKGBuffer struct {
	data []byte
	wg   *sync.WaitGroup
}

func (b *PKGBuffer) GetData() []byte {
	return b.data
}

func (b *PKGBuffer) Done() {
	b.data = nil
	if b.wg != nil {
		b.wg.Done()
		b.wg = nil
	}
	pkgBufferPool.Put(b)
}

func StringToCByte(str string) []byte {
	b := []byte(str)
	b = append(b, 0)
	return b
}

func CsHeadVisualize(head *tcaplus_protocol_cs.TCaplusPkgHead) string {
	var keyInfo string
	for i := 0; i < int(head.KeyInfo.FieldNum); i++ {
		keyInfo += head.KeyInfo.Fields[i].FieldName
		keyInfo += ":"
		keyInfo += fmt.Sprintf("%v", head.KeyInfo.Fields[i].FieldBuff[0:head.KeyInfo.Fields[i].FieldLen])
	}
	return fmt.Sprintf("{ Result:%d Magic:%d Version:%d HeadLen:%d BodyLen:%d AsynID:%d Seq:%d Cmd:%d"+
		" SubCmd:%d Flags:%d AppID:%d ZoneId:%d ShardID:%d Table:%s RecVersion:%d KeyFieldNum:%d KeyInfo:{%s} }",
		head.Result, head.Magic, head.Version, head.HeadLen, head.BodyLen, head.AsynID, head.Seq, head.Cmd, head.SubCmd,
		head.Flags, head.RouterInfo.AppID, head.RouterInfo.ZoneID, head.RouterInfo.ShardID,
		string(head.RouterInfo.TableName[0:head.RouterInfo.TableNameLen]),
		head.KeyInfo.Version, head.KeyInfo.FieldNum, keyInfo)
}

func CovertToJson(v interface{}) string {
	data, _ := json.Marshal(v)
	buf := &bytes.Buffer{}
	json.Indent(buf, data, "", "\t")
	return buf.String()
}
