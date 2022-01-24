package tnet

import (
	"net"
	"sync"
	"sync/atomic"
)

const (
	BigPkgMemorySize = 11 * 1024 * 1024 //收包的最大包大小10MB,设置包内存块最大为11MB
)

var PkgMemorySize = 0 //小内存块大小1MB-11MB，当小内存不够用时，申请大内存

//tcp拆包，从大块内存中拆小块消息
//内存池
type PKGMemory struct {
	data   []byte
	refNum int64
	prev   int
	last   int
}

//memory缓存池
var pkgMemoryPool = sync.Pool{
	New: func() interface{} {
		return &PKGMemory{data: make([]byte, PkgMemorySize)}
	},
}

var bigPkgMemoryPool = sync.Pool{
	New: func() interface{} {
		return &PKGMemory{data: make([]byte, BigPkgMemorySize)}
	},
}

//消息，内存引用PKGMemory中的地址
type PKG struct {
	data   []byte
	mem    *PKGMemory
	cbPara interface{}
}

//pkg缓存池
var pkgPool = sync.Pool{
	New: func() interface{} {
		return &PKG{}
	},
}

func GetPKGMemory(oldMemory *PKGMemory) *PKGMemory {
	var newMemory *PKGMemory
	if oldMemory != nil &&
		(oldMemory.last-oldMemory.prev >= PkgMemorySize) {
		//申请大内存
		newMemory = bigPkgMemoryPool.Get().(*PKGMemory)
	} else {
		//申请小内存
		newMemory = pkgMemoryPool.Get().(*PKGMemory)
	}
	newMemory.refNum = 1
	newMemory.prev = 0
	newMemory.last = 0
	if oldMemory != nil && oldMemory.prev < oldMemory.last {
		newMemory.last = copy(newMemory.data, oldMemory.data[oldMemory.prev:oldMemory.last])
		PutPKGMemory(oldMemory)
	}
	return newMemory
}

func PutPKGMemory(mem *PKGMemory) {
	//判断引用计数，引用计数为空即可释放
	ref := atomic.AddInt64(&mem.refNum, -1)
	if ref == 0 {
		if cap(mem.data) == BigPkgMemorySize {
			bigPkgMemoryPool.Put(mem)
		} else {
			pkgMemoryPool.Put(mem)
		}
	}
}

func (m *PKGMemory) ReadFromNetConn(conn net.Conn) (int, error) {
	n, err := conn.Read(m.data[m.last:])
	m.last += n
	return n, err
}

func (m *PKGMemory) ValidBuffer() []byte {
	return m.data[m.prev:m.last]
}

func (m *PKGMemory) ValidLength() int {
	return m.last - m.prev
}

func (m *PKGMemory) GetPkg(size int) *PKG {
	pkg := pkgPool.Get().(*PKG)
	pkg.mem = m
	pkg.data = m.data[m.prev : m.prev+size]
	m.prev += size
	atomic.AddInt64(&m.refNum, 1)
	return pkg
}

func (m *PKGMemory) BufferIsFull() bool {
	return m.last == cap(m.data)
}

//pkg func
func (p *PKG) GetData() []byte {
	return p.data
}
func (p *PKG) GetCbPara() interface{} {
	return p.cbPara
}

func (p *PKG) Done() {
	p.data = nil
	p.cbPara = nil
	if p.mem != nil {
		PutPKGMemory(p.mem)
		p.mem = nil
	}
	pkgPool.Put(p)
}
