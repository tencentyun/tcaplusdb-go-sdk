package statistics

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/dir"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	tcaplusCmd "github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/tcaplus_protocol_cs"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/terror"

	"sync/atomic"
	"time"
)

type StatInterface interface {
	// OutPutZoneSummaryStat 输出zone的汇总统计，请实现非阻塞接口
	OutPutZoneSummaryStat(zoneId uint32, stat *StatNode)

	// OutPutZoneProxyStat 输出zone的单个proxy统计，请实现非阻塞接口
	OutPutZoneProxyStat(zoneId uint32, proxyUrl string, stat *StatNode)
}

type DelayRange struct {
	Cnt               uint32       // 总数
	DelaySum          uint64       // 时延总计(平均时延= DelaySum/Cnt/1000)
	MaxDelayMs        uint32       // 最大时延ms
	MaxDelayProxy     atomic.Value // 最大时延url,string
	MaxDelayTimestamp uint64       // 最大时延ms时间点
	MaxDelayCmd       uint32       // 最大时延命令字
	DelayLowZeroCnt   uint32       // <0的数量
	DelayLow10msCnt   uint32       // <=10ms的数量
	DelayLow100msCnt  uint32       // <=100ms的数量
	DelayLow500msCnt  uint32       // <=500ms的数量
	DelayLow1000msCnt uint32       // <=1000ms的数量
	DelayHigh1sCnt    uint32       // >1s的数量
}

func (d *DelayRange) GetMaxDelayProxy() string {
	if d.MaxDelayProxy.Load() != nil {
		return d.MaxDelayProxy.Load().(string)
	}
	return ""
}

func (d *DelayRange) CalcDelay(startUs uint64, endUs uint64, curUs uint64, url string, cmd uint32) {
	if endUs < startUs {
		atomic.AddUint32(&d.DelayLowZeroCnt, 1)
		return
	}
	delayUs := endUs - startUs
	delayMs := uint32(delayUs / 1000)
	atomic.AddUint64(&d.DelaySum, uint64(delayUs))
	atomic.AddUint32(&d.Cnt, 1)
	if delayMs > atomic.LoadUint32(&d.MaxDelayMs) {
		atomic.StoreUint32(&d.MaxDelayMs, delayMs)
		atomic.StoreUint64(&d.MaxDelayTimestamp, curUs/1000)
		d.MaxDelayProxy.Store(url)
		atomic.StoreUint32(&d.MaxDelayCmd, cmd)
	}
	if delayMs <= 10 {
		atomic.AddUint32(&d.DelayLow10msCnt, 1)
		return
	}
	if delayMs <= 100 {
		atomic.AddUint32(&d.DelayLow100msCnt, 1)
		return
	}
	if delayMs <= 500 {
		atomic.AddUint32(&d.DelayLow500msCnt, 1)
		return
	}
	if delayMs <= 1000 {
		atomic.AddUint32(&d.DelayLow1000msCnt, 1)
		return
	}
	atomic.AddUint32(&d.DelayHigh1sCnt, 1)
	return
}

type PerfDelay struct {
	PerfCnt         uint32     // 携带perf响应的数量
	ApiTotalDelay   DelayRange // api总时延区间
	ApiToProxyDelay DelayRange // api->proxy时延区间
	ProxyToSvrDelay DelayRange // proxy->svr时延区间
	SvrHandleDelay  DelayRange // svr处理时延区间
	SvrToProxyDelay DelayRange // svr->proxy时延区间
	ProxyTotalDelay DelayRange // proxy total时延区间
	ProxyToApiDelay DelayRange // proxy->api时延区间
}

type StatNode struct {
	ReqRouteFailCnt uint32 // 路由失败次数
	ResTimeoutCnt   uint32 // 协程响应超时数

	SimpleReqTotalCnt uint32 // 请求总数
	SimpleReqSucCnt   uint32 // 请求发送成功数
	SimpleReqFailCnt  uint32 // 请求发送失败数

	SimpleResTotalCnt   uint32    // 响应总数
	SimpleResSucCnt     uint32    // 响应成功数
	SimpleResFailCnt    uint32    // 响应失败数
	SimpleResWarningCnt uint32    // 响应告警数
	SimplePerf          PerfDelay // 简单请求时延区间

	ComplexReqTotalCnt uint32 // 请求总数
	ComplexReqSucCnt   uint32 // 请求发送成功数
	ComplexReqFailCnt  uint32 // 请求发送失败数

	ComplexResTotalCnt   uint32 // 响应总数
	ComplexResSucCnt     uint32 // 响应成功数
	ComplexResFailCnt    uint32 // 响应失败数
	ComplexResWarningCnt uint32 // 响应告警数

	IsoProxyCnt uint32 // 隔离proxy次数
}

func (s *StatNode) IncIsoProxyCnt() {
	atomic.AddUint32(&s.IsoProxyCnt, 1)
}

func (s *StatNode) IncResTimeout() {
	atomic.AddUint32(&s.ResTimeoutCnt, 1)
}

func (s *StatNode) IncSimpleReqSuc() {
	atomic.AddUint32(&s.SimpleReqSucCnt, 1)
	atomic.AddUint32(&s.SimpleReqTotalCnt, 1)
}

func (s *StatNode) IncSimpleReqFail() {
	atomic.AddUint32(&s.SimpleReqFailCnt, 1)
	atomic.AddUint32(&s.SimpleReqTotalCnt, 1)
}

func (s *StatNode) IncComplexReqSuc() {
	atomic.AddUint32(&s.ComplexReqSucCnt, 1)
	atomic.AddUint32(&s.ComplexReqTotalCnt, 1)
}

func (s *StatNode) IncComplexReqFail() {
	atomic.AddUint32(&s.ComplexReqFailCnt, 1)
	atomic.AddUint32(&s.ComplexReqTotalCnt, 1)
}

func (s *StatNode) IncSimpleResSuc() {
	atomic.AddUint32(&s.SimpleResSucCnt, 1)
	atomic.AddUint32(&s.SimpleResTotalCnt, 1)
}

func (s *StatNode) IncSimpleResFail() {
	atomic.AddUint32(&s.SimpleResFailCnt, 1)
	atomic.AddUint32(&s.SimpleResTotalCnt, 1)
}

func (s *StatNode) IncSimpleResWarn() {
	atomic.AddUint32(&s.SimpleResWarningCnt, 1)
	atomic.AddUint32(&s.SimpleResTotalCnt, 1)
}

func (s *StatNode) IncComplexResSuc() {
	atomic.AddUint32(&s.ComplexResSucCnt, 1)
	atomic.AddUint32(&s.ComplexResTotalCnt, 1)
}

func (s *StatNode) IncComplexResFail() {
	atomic.AddUint32(&s.ComplexResFailCnt, 1)
	atomic.AddUint32(&s.ComplexResTotalCnt, 1)
}

func (s *StatNode) IncComplexResWarn() {
	atomic.AddUint32(&s.ComplexResWarningCnt, 1)
	atomic.AddUint32(&s.ComplexResTotalCnt, 1)
}

func (s *StatNode) IncRouteFail() {
	atomic.AddUint32(&s.ReqRouteFailCnt, 1)
}

func GetAvgDelayMs(totalUs uint64, cnt uint64) uint64 {
	if cnt <= 0 {
		return 0
	}
	return totalUs / cnt / 1000
}

func (s *StatNode) Print(title string) {
	logger.INFO("%s\n\t[ResTimeoutCnt]:%d\n\t[ReqRouteFailCnt]:%d\n\t[SimpleReqTotalCnt]:%d\n\t[SimpleReqSucCnt]:%d\n\t"+
		"[SimpleReqFailCnt]:%d\n\t[SimpleResTotalCnt]:%d\n\t[SimpleResSucCnt]:%d\n\t[SimpleResFailCnt]:%d\n\t"+
		"[SimpleResWarningCnt]:%d\n\t[ComplexReqTotalCnt]:%d\n\t[ComplexReqSucCnt]:%d\n\t[ComplexReqFailCnt]:%d\n\t"+
		"[ComplexResTotalCnt]:%d\n\t[ComplexResSucCnt]:%d\n\t[ComplexResFailCnt]:%d\n\t[ComplexResWarningCnt]:%d\n\t"+
		"[SimpleResPerf]:%d\n\t"+
		"[ApiTotal  ][<0]:%d [0-10ms]:%d [10-100ms]:%d [100-500ms]:%d [500-1000ms]:%d [1000ms+]:%d [AvgDelay(ms)]:%d [MaxDelay(ms)]:%d(time:%d,%s,cmd:%d)\n\t"+
		"[ApiToProxy][<0]:%d [0-10ms]:%d [10-100ms]:%d [100-500ms]:%d [500-1000ms]:%d [1000ms+]:%d [AvgDelay(ms)]:%d [MaxDelay(ms)]:%d(time:%d,%s,cmd:%d)\n\t"+
		"[ProxyToSvr][<0]:%d [0-10ms]:%d [10-100ms]:%d [100-500ms]:%d [500-1000ms]:%d [1000ms+]:%d [AvgDelay(ms)]:%d [MaxDelay(ms)]:%d(time:%d,%s,cmd:%d)\n\t"+
		"[SvrTotal  ][<0]:%d [0-10ms]:%d [10-100ms]:%d [100-500ms]:%d [500-1000ms]:%d [1000ms+]:%d [AvgDelay(ms)]:%d [MaxDelay(ms)]:%d(time:%d,%s,cmd:%d)\n\t"+
		"[SvrToProxy][<0]:%d [0-10ms]:%d [10-100ms]:%d [100-500ms]:%d [500-1000ms]:%d [1000ms+]:%d [AvgDelay(ms)]:%d [MaxDelay(ms)]:%d(time:%d,%s,cmd:%d)\n\t"+
		"[ProxyTotal][<0]:%d [0-10ms]:%d [10-100ms]:%d [100-500ms]:%d [500-1000ms]:%d [1000ms+]:%d [AvgDelay(ms)]:%d [MaxDelay(ms)]:%d(time:%d,%s,cmd:%d)\n\t"+
		"[ProxyToApi][<0]:%d [0-10ms]:%d [10-100ms]:%d [100-500ms]:%d [500-1000ms]:%d [1000ms+]:%d [AvgDelay(ms)]:%d [MaxDelay(ms)]:%d(time:%d,%s,cmd:%d)\n",
		title,
		atomic.LoadUint32(&s.ResTimeoutCnt),
		atomic.LoadUint32(&s.ReqRouteFailCnt),
		atomic.LoadUint32(&s.SimpleReqTotalCnt),
		atomic.LoadUint32(&s.SimpleReqSucCnt),
		atomic.LoadUint32(&s.SimpleReqFailCnt),
		atomic.LoadUint32(&s.SimpleResTotalCnt),
		atomic.LoadUint32(&s.SimpleResSucCnt),
		atomic.LoadUint32(&s.SimpleResFailCnt),
		atomic.LoadUint32(&s.SimpleResWarningCnt),
		atomic.LoadUint32(&s.ComplexReqTotalCnt),
		atomic.LoadUint32(&s.ComplexReqSucCnt),
		atomic.LoadUint32(&s.ComplexReqFailCnt),
		atomic.LoadUint32(&s.ComplexResTotalCnt),
		atomic.LoadUint32(&s.ComplexResSucCnt),
		atomic.LoadUint32(&s.ComplexResFailCnt),
		atomic.LoadUint32(&s.ComplexResWarningCnt),
		atomic.LoadUint32(&s.SimplePerf.PerfCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiTotalDelay.DelayLowZeroCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiTotalDelay.DelayLow10msCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiTotalDelay.DelayLow100msCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiTotalDelay.DelayLow500msCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiTotalDelay.DelayLow1000msCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiTotalDelay.DelayHigh1sCnt),
		GetAvgDelayMs(atomic.LoadUint64(&s.SimplePerf.ApiTotalDelay.DelaySum), uint64(atomic.LoadUint32(&s.SimplePerf.ApiTotalDelay.Cnt))),
		atomic.LoadUint32(&s.SimplePerf.ApiTotalDelay.MaxDelayMs),
		atomic.LoadUint64(&s.SimplePerf.ApiTotalDelay.MaxDelayTimestamp),
		s.SimplePerf.ApiTotalDelay.GetMaxDelayProxy(),
		atomic.LoadUint32(&s.SimplePerf.ApiTotalDelay.MaxDelayCmd),
		atomic.LoadUint32(&s.SimplePerf.ApiToProxyDelay.DelayLowZeroCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiToProxyDelay.DelayLow10msCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiToProxyDelay.DelayLow100msCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiToProxyDelay.DelayLow500msCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiToProxyDelay.DelayLow1000msCnt),
		atomic.LoadUint32(&s.SimplePerf.ApiToProxyDelay.DelayHigh1sCnt),
		GetAvgDelayMs(atomic.LoadUint64(&s.SimplePerf.ApiToProxyDelay.DelaySum), uint64(atomic.LoadUint32(&s.SimplePerf.ApiToProxyDelay.Cnt))),
		atomic.LoadUint32(&s.SimplePerf.ApiToProxyDelay.MaxDelayMs),
		atomic.LoadUint64(&s.SimplePerf.ApiToProxyDelay.MaxDelayTimestamp),
		s.SimplePerf.ApiToProxyDelay.GetMaxDelayProxy(),
		atomic.LoadUint32(&s.SimplePerf.ApiToProxyDelay.MaxDelayCmd),
		atomic.LoadUint32(&s.SimplePerf.ProxyToSvrDelay.DelayLowZeroCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToSvrDelay.DelayLow10msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToSvrDelay.DelayLow100msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToSvrDelay.DelayLow500msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToSvrDelay.DelayLow1000msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToSvrDelay.DelayHigh1sCnt),
		GetAvgDelayMs(atomic.LoadUint64(&s.SimplePerf.ProxyToSvrDelay.DelaySum), uint64(atomic.LoadUint32(&s.SimplePerf.ProxyToSvrDelay.Cnt))),
		atomic.LoadUint32(&s.SimplePerf.ProxyToSvrDelay.MaxDelayMs),
		atomic.LoadUint64(&s.SimplePerf.ProxyToSvrDelay.MaxDelayTimestamp),
		s.SimplePerf.ProxyToSvrDelay.GetMaxDelayProxy(),
		atomic.LoadUint32(&s.SimplePerf.ProxyToSvrDelay.MaxDelayCmd),
		atomic.LoadUint32(&s.SimplePerf.SvrHandleDelay.DelayLowZeroCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrHandleDelay.DelayLow10msCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrHandleDelay.DelayLow100msCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrHandleDelay.DelayLow500msCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrHandleDelay.DelayLow1000msCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrHandleDelay.DelayHigh1sCnt),
		GetAvgDelayMs(atomic.LoadUint64(&s.SimplePerf.SvrHandleDelay.DelaySum), uint64(atomic.LoadUint32(&s.SimplePerf.SvrHandleDelay.Cnt))),
		atomic.LoadUint32(&s.SimplePerf.SvrHandleDelay.MaxDelayMs),
		atomic.LoadUint64(&s.SimplePerf.SvrHandleDelay.MaxDelayTimestamp),
		s.SimplePerf.SvrHandleDelay.GetMaxDelayProxy(),
		atomic.LoadUint32(&s.SimplePerf.SvrHandleDelay.MaxDelayCmd),
		atomic.LoadUint32(&s.SimplePerf.SvrToProxyDelay.DelayLowZeroCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrToProxyDelay.DelayLow10msCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrToProxyDelay.DelayLow100msCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrToProxyDelay.DelayLow500msCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrToProxyDelay.DelayLow1000msCnt),
		atomic.LoadUint32(&s.SimplePerf.SvrToProxyDelay.DelayHigh1sCnt),
		GetAvgDelayMs(atomic.LoadUint64(&s.SimplePerf.SvrToProxyDelay.DelaySum), uint64(atomic.LoadUint32(&s.SimplePerf.SvrToProxyDelay.Cnt))),
		atomic.LoadUint32(&s.SimplePerf.SvrToProxyDelay.MaxDelayMs),
		atomic.LoadUint64(&s.SimplePerf.SvrToProxyDelay.MaxDelayTimestamp),
		s.SimplePerf.SvrToProxyDelay.GetMaxDelayProxy(),
		atomic.LoadUint32(&s.SimplePerf.SvrToProxyDelay.MaxDelayCmd),
		atomic.LoadUint32(&s.SimplePerf.ProxyTotalDelay.DelayLowZeroCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyTotalDelay.DelayLow10msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyTotalDelay.DelayLow100msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyTotalDelay.DelayLow500msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyTotalDelay.DelayLow1000msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyTotalDelay.DelayHigh1sCnt),
		GetAvgDelayMs(atomic.LoadUint64(&s.SimplePerf.ProxyTotalDelay.DelaySum), uint64(atomic.LoadUint32(&s.SimplePerf.ProxyTotalDelay.Cnt))),
		atomic.LoadUint32(&s.SimplePerf.ProxyTotalDelay.MaxDelayMs),
		atomic.LoadUint64(&s.SimplePerf.ProxyTotalDelay.MaxDelayTimestamp),
		s.SimplePerf.ProxyTotalDelay.GetMaxDelayProxy(),
		atomic.LoadUint32(&s.SimplePerf.ProxyTotalDelay.MaxDelayCmd),
		atomic.LoadUint32(&s.SimplePerf.ProxyToApiDelay.DelayLowZeroCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToApiDelay.DelayLow10msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToApiDelay.DelayLow100msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToApiDelay.DelayLow500msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToApiDelay.DelayLow1000msCnt),
		atomic.LoadUint32(&s.SimplePerf.ProxyToApiDelay.DelayHigh1sCnt),
		GetAvgDelayMs(atomic.LoadUint64(&s.SimplePerf.ProxyToApiDelay.DelaySum), uint64(atomic.LoadUint32(&s.SimplePerf.ProxyToApiDelay.Cnt))),
		atomic.LoadUint32(&s.SimplePerf.ProxyToApiDelay.MaxDelayMs),
		atomic.LoadUint64(&s.SimplePerf.ProxyToApiDelay.MaxDelayTimestamp),
		s.SimplePerf.ProxyToApiDelay.GetMaxDelayProxy(),
		atomic.LoadUint32(&s.SimplePerf.ProxyToApiDelay.MaxDelayCmd))
}

type ConnStat struct {
	stat                atomic.Value
	mgr                 *StatMgr
	url                 string
	statOutPutInterface StatInterface
}

func (c *ConnStat) CalcPerf(stat *StatNode, summary *StatNode, perfData []byte, cmd uint32) {
	perf := tcaplus_protocol_cs.NewPerfTest()
	if err := perf.Unpack(tcaplus_protocol_cs.TCaplusPkgCurrentVersion, perfData); err != nil {
		return
	}

	atomic.AddUint32(&stat.SimplePerf.PerfCnt, 1)
	atomic.AddUint32(&summary.SimplePerf.PerfCnt, 1)
	curUs := uint64(time.Now().UnixMicro())

	//ApiTotalDelay
	stat.SimplePerf.ApiTotalDelay.CalcDelay(perf.ApiSendTime, curUs, curUs, c.url, cmd)
	summary.SimplePerf.ApiTotalDelay.CalcDelay(perf.ApiSendTime, curUs, curUs, c.url, cmd)

	//ApiToProxyDelay
	stat.SimplePerf.ApiToProxyDelay.CalcDelay(perf.ApiSendTime, perf.ProxyFrontForwardTime, curUs, c.url, cmd)
	summary.SimplePerf.ApiToProxyDelay.CalcDelay(perf.ApiSendTime, perf.ProxyFrontForwardTime, curUs, c.url, cmd)

	// ProxyToSvrDelay
	stat.SimplePerf.ProxyToSvrDelay.CalcDelay(perf.ProxyEndForwardTime, perf.SvrTbusppForwardTime, curUs, c.url, cmd)
	summary.SimplePerf.ProxyToSvrDelay.CalcDelay(perf.ProxyEndForwardTime, perf.SvrTbusppForwardTime, curUs, c.url, cmd)

	// SvrHandleDelay
	stat.SimplePerf.SvrHandleDelay.CalcDelay(perf.SvrTbusppForwardTime, perf.SvrWorkerSendTime, curUs, c.url, cmd)
	summary.SimplePerf.SvrHandleDelay.CalcDelay(perf.SvrTbusppForwardTime, perf.SvrWorkerSendTime, curUs, c.url, cmd)

	// SvrToProxyDelay
	stat.SimplePerf.SvrToProxyDelay.CalcDelay(perf.SvrWorkerSendTime, perf.ProxyFrontBackwardTime, curUs, c.url, cmd)
	summary.SimplePerf.SvrToProxyDelay.CalcDelay(perf.SvrWorkerSendTime, perf.ProxyFrontBackwardTime, curUs, c.url, cmd)

	// ProxyTotalDelay
	stat.SimplePerf.ProxyTotalDelay.CalcDelay(perf.ProxyFrontForwardTime, perf.ProxyEndBackwardTime, curUs, c.url, cmd)
	summary.SimplePerf.ProxyTotalDelay.CalcDelay(perf.ProxyFrontForwardTime, perf.ProxyEndBackwardTime, curUs, c.url, cmd)

	// ProxyToApiDelay
	stat.SimplePerf.ProxyToApiDelay.CalcDelay(perf.ProxyEndBackwardTime, curUs, curUs, c.url, cmd)
	summary.SimplePerf.ProxyToApiDelay.CalcDelay(perf.ProxyEndBackwardTime, curUs, curUs, c.url, cmd)
}

func (c *ConnStat) ProcResPerf(cmd uint32, result int32, perfData []byte, perfLen uint32) {
	stat := c.stat.Load().(*StatNode)
	summaryStat := c.mgr.summaryStat.Load().(*StatNode)
	if tcaplusCmd.IsSimpleRes(cmd) {
		if result >= 0 {
			stat.IncSimpleResSuc()
			summaryStat.IncSimpleResSuc()
		} else if terror.IsWarningCode(result) {
			stat.IncSimpleResWarn()
			summaryStat.IncSimpleResWarn()
		} else {
			stat.IncSimpleResFail()
			summaryStat.IncSimpleResFail()
		}

		// calc perf
		if perfLen > 0 && perfData != nil {
			c.CalcPerf(stat, summaryStat, perfData, cmd)
		}
		return
	}

	if tcaplusCmd.IsComplexRes(cmd) {
		if result >= 0 {
			stat.IncComplexResSuc()
			summaryStat.IncComplexResSuc()
		} else if terror.IsWarningCode(result) {
			stat.IncComplexResWarn()
			summaryStat.IncComplexResWarn()
		} else {
			stat.IncComplexResFail()
			summaryStat.IncComplexResFail()
		}
		return
	}
}

func (c *ConnStat) Init(mgr *StatMgr, url string, statOutPutInterface StatInterface) {
	c.url = url
	c.mgr = mgr
	c.statOutPutInterface = statOutPutInterface
	node := &StatNode{}
	c.stat.Store(node)
}

func (c *ConnStat) IncResTimeout() {
	stat := c.stat.Load().(*StatNode)
	summaryStat := c.mgr.summaryStat.Load().(*StatNode)
	stat.IncResTimeout()
	summaryStat.IncResTimeout()
	return
}

func (c *ConnStat) IncIsoProxyCnt() {
	summaryStat := c.mgr.summaryStat.Load().(*StatNode)
	summaryStat.IncIsoProxyCnt()
	return
}

func (c *ConnStat) IncReqSuc(cmd uint32) {
	stat := c.stat.Load().(*StatNode)
	summaryStat := c.mgr.summaryStat.Load().(*StatNode)
	if tcaplusCmd.IsSimpleReq(cmd) {
		stat.IncSimpleReqSuc()
		summaryStat.IncSimpleReqSuc()
		return
	}

	if tcaplusCmd.IsComplexReq(cmd) {
		stat.IncComplexReqSuc()
		summaryStat.IncComplexReqSuc()
		return
	}
}

func (c *ConnStat) IncReqFail(cmd uint32) {
	stat := c.stat.Load().(*StatNode)
	summaryStat := c.mgr.summaryStat.Load().(*StatNode)
	if tcaplusCmd.IsSimpleReq(cmd) {
		stat.IncSimpleReqFail()
		summaryStat.IncSimpleReqFail()
		return
	}

	if tcaplusCmd.IsComplexReq(cmd) {
		stat.IncComplexReqFail()
		summaryStat.IncComplexReqFail()
		return
	}
}

func (c *ConnStat) Reset() {
	stat := c.stat.Load().(*StatNode)
	newNode := &StatNode{}
	c.stat.Store(newNode)
	if stat != nil && atomic.LoadUint32(&stat.SimplePerf.ApiTotalDelay.MaxDelayMs) > 1000 {
		// proxy太多，只打印超过1s的proxy
		title := fmt.Sprintf("App %d zone %d Proxy %s stat info", c.mgr.app, c.mgr.zone, c.url)
		stat.Print(title)
	}
	if c.statOutPutInterface != nil {
		c.statOutPutInterface.OutPutZoneProxyStat(c.mgr.zone, c.url, stat)
	}
}

type StatMgr struct {
	summaryStat         atomic.Value
	app                 uint32
	zone                uint32
	statOutPutInterface StatInterface
	dirReportStat       *dir.ReportStat
}

func (s *StatMgr) Init(app uint32, zone uint32, outPutInterface StatInterface, dirReportStat *dir.ReportStat) {
	s.app = app
	s.zone = zone
	s.statOutPutInterface = outPutInterface
	s.dirReportStat = dirReportStat
	node := &StatNode{}
	s.summaryStat.Store(node)
}

func (s *StatMgr) Reset() {
	stat := s.summaryStat.Load().(*StatNode)
	newNode := &StatNode{}
	s.summaryStat.Store(newNode)
	if stat != nil {
		title := fmt.Sprintf("App %d zone %d SummaryStat info:", s.app, s.zone)
		stat.Print(title)
		s.SetDirReportStr(stat)
	}
	if s.statOutPutInterface != nil {
		s.statOutPutInterface.OutPutZoneSummaryStat(s.zone, stat)
	}
}

func (s *StatMgr) SetDirReportStr(stat *StatNode) {
	//${ZoneId}-${请求成功数}-${请求失败数}-${响应成功数}-${响应WARN数}-${响应失败数}-${超时数}-${路由失败数}-${平均时延}-${最大时延}-${最大时延proxy}-${最大时延时间点}-${最大时延CMD}-${最长收包时间间隔}-${300ms隔离proxy次数}
	reportStr := fmt.Sprintf("%d-%d-%d-%d-%d-%d-%d-%d-%d-%d-%s-%d-%d-0-%d", s.zone,
		atomic.LoadUint32(&stat.SimpleReqSucCnt)+atomic.LoadUint32(&stat.ComplexReqSucCnt),
		atomic.LoadUint32(&stat.SimpleReqFailCnt)+atomic.LoadUint32(&stat.ComplexReqFailCnt),
		atomic.LoadUint32(&stat.SimpleResSucCnt)+atomic.LoadUint32(&stat.ComplexResSucCnt),
		atomic.LoadUint32(&stat.SimpleResWarningCnt)+atomic.LoadUint32(&stat.ComplexResWarningCnt),
		atomic.LoadUint32(&stat.SimpleResFailCnt)+atomic.LoadUint32(&stat.ComplexResFailCnt),
		atomic.LoadUint32(&stat.ResTimeoutCnt),
		atomic.LoadUint32(&stat.ReqRouteFailCnt),
		GetAvgDelayMs(atomic.LoadUint64(&stat.SimplePerf.ApiTotalDelay.DelaySum), uint64(atomic.LoadUint32(&stat.SimplePerf.ApiTotalDelay.Cnt))),
		atomic.LoadUint32(&stat.SimplePerf.ApiTotalDelay.MaxDelayMs),
		stat.SimplePerf.ApiTotalDelay.GetMaxDelayProxy(),
		atomic.LoadUint64(&stat.SimplePerf.ApiTotalDelay.MaxDelayTimestamp),
		atomic.LoadUint32(&stat.SimplePerf.ApiTotalDelay.MaxDelayCmd),
		atomic.LoadUint32(&stat.IsoProxyCnt))
	s.dirReportStat.SetZoneReportStat(s.zone, reportStr)
}

func (s *StatMgr) NewConnStat(proxyUrl string) *ConnStat {
	stat := &ConnStat{}
	stat.Init(s, proxyUrl, s.statOutPutInterface)
	return stat
}

func (s *StatMgr) IncRouteFail() {
	stat := s.summaryStat.Load().(*StatNode)
	atomic.AddUint32(&stat.ReqRouteFailCnt, 1)
}
