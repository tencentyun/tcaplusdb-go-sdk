package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/example/PB/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/record"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/request"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var pbclient *tcaplus.PBClient
var once sync.Once

type ApiCfg struct {
	XMLName   xml.Name `xml:"ApiTstCfg"` // 指定最外层的标签为ApiTstCfg
	DirUrl    string   `xml:"dir_addr"`
	AppId     uint64   `xml:"app_id"`
	ZoneId    uint32   `xml:"zone_id"`
	Signature string   `xml:"signature"`
	PBTable   string   `xml:"pb_table"`
}

var ApiConfig *ApiCfg = nil

func ReadApiCfg(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("ReadFile " + filePath + " err:" + err.Error())
		return err
	}

	ApiConfig = new(ApiCfg)
	err = xml.Unmarshal(data, ApiConfig)
	if err != nil {
		fmt.Println("xml Unmarshal " + filePath + " err:" + err.Error())
		return err
	}

	fmt.Printf("API CONF: %+v\n", *ApiConfig)
	return nil
}

func InitPBClient() *tcaplus.PBClient {
	var err error
	once.Do(func() {
		err = ReadApiCfg("api_cfg.xml")
		if err != nil {
			fmt.Printf("ReadApiCfg fail %s", err.Error())
			return
		}

		pbclient = tcaplus.NewPBClient()
		err = pbclient.SetLogCfg("logconf.xml")
		if err != nil {
			fmt.Printf("excepted SetLogCfg success")
			return
		}

		tables := strings.Split(ApiConfig.PBTable, ",")
		zoneTable := map[uint32][]string{ApiConfig.ZoneId: tables}
		err = pbclient.Dial(ApiConfig.AppId, []uint32{ApiConfig.ZoneId}, ApiConfig.DirUrl, ApiConfig.Signature, 30, zoneTable)
		if err != nil {
			fmt.Printf("excepted dial success, %s", err.Error())
			return
		}
	})
	if err != nil {
		fmt.Println("\ninit fail. please check tcaplus config")
		os.Exit(-1)
		return nil
	}
	return pbclient
}

var (
	ttt      = flag.Int("t", 5, "route num")
	nnn      = flag.Int("n", 2000, "num")
	size     = flag.Int("size", 0, "G")
	core     = flag.Int("core", 4, "core num")
	mode     = flag.String("mode", "async", "sync|async")
	fff      = flag.String("f", "add", "select func add|get|mix|insert")
	ppp      = flag.String("path", "", "")
	interval = flag.Int("interval", 10, "time interval")
)

func main() {
	flag.Parse()
	if *core < runtime.NumCPU() {
		runtime.GOMAXPROCS(*core)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if *fff == "insert" {
		TcaplusInsert(*size)
	} else if *mode == "sync" {
		TcaplusSyncTest(*ttt, *nnn, *fff, *ppp, *interval)
	} else {
		TcaplusAsyncTest(*ttt, *nnn, *fff, *ppp, *interval)
	}
}

func TcaplusInsert(size int) {
	c := InitPBClient()

	records := make([]*tcaplusservice.GamePlayers, 1000)
	var playerId int64
	for i := 0; i < 1000; i++ {
		records[i] = &tcaplusservice.GamePlayers{
			PlayerId:        0,
			PlayerName:      "12345",
			PlayerEmail:     "12345@qq.com",
			GameServerId:    10,
			LoginTimestamp:  []string{Bytes500B()},
			LogoutTimestamp: []string{Bytes500B()},
			IsOnline:        false,
			Pay: &tcaplusservice.Payment{
				PayId:  10101,
				Amount: 1000,
				Method: 1,
			},
		}
	}

	reqs := make([]request.TcaplusRequest, 0, 1000)
	recs := make([]*record.Record, 0, 1000)
	startCh := make([]chan struct{}, 1000)
	num := 1000 * size
	var count uint64

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		req, err := c.NewRequest(ApiConfig.ZoneId, ApiConfig.PBTable, cmd.TcaplusApiReplaceReq)
		if err != nil {
			fmt.Println(err)
			return
		}
		rec, err := req.AddRecord(0)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.SetResultFlagForSuccess(0)
		reqs = append(reqs, req)
		recs = append(recs, rec)
		startCh[i] = make(chan struct{}, 1)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for n := 0; n < num; n++ {
				records[i].PlayerId = atomic.AddInt64(&playerId, 1)
				if _, err = recs[i].SetPBData(records[i]); err != nil {
					fmt.Println(err)
					return
				}
				resp, err := c.Do(reqs[i], time.Minute)
				if atomic.AddUint64(&count, 1)%100000 == 0 {
					fmt.Printf("%s insert 100M\n", time.Now().Format("2006/1/2 15:04:05"))
				}
				if err != nil || resp.GetResult() != 0 {
					fmt.Println(err, resp.GetResult())
				}
			}
		}(i)
	}
	wg.Wait()

	fmt.Println("insert ", size, "G data success")
}

func Bytes500B() string {
	buf := make([]byte, 500)
	for i, _ := range buf {
		buf[i] = 100
	}
	return string(buf)
}

func TcaplusAsyncTest(tCount int, num int, f, p string, itv int) {
	cmdId := cmd.TcaplusApiGetReq
	flagId := 0
	if f == "add" {
		cmdId = cmd.TcaplusApiReplaceReq
		flagId = 0
	} else if f == "mix" {
		cmdId = cmd.TcaplusApiReplaceReq
		flagId = 3
	}

	c := InitPBClient()

	records := make([]*tcaplusservice.GamePlayers, tCount*num)
	for i := 0; i < tCount*num; i++ {
		randseed := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(int64(i)*12345))))
		records[i] = &tcaplusservice.GamePlayers{
			PlayerId:        int64(i) * 12345,
			PlayerName:      string(randseed[:]),
			PlayerEmail:     string(randseed[:]),
			GameServerId:    10,
			LoginTimestamp:  []string{"2019-12-12 15:00:00"},
			LogoutTimestamp: []string{"2019-12-12 16:00:00"},
			IsOnline:        false,
			Pay: &tcaplusservice.Payment{
				PayId:  10101,
				Amount: 1000,
				Method: 1,
			},
		}
	}

	allcount := 0
	var countTime uint64
	var maxTime uint64
	var failCount int
	gch := make(chan uint64)

	go func() {
		for {
			resp, _ := c.RecvResponse()
			if resp != nil {
				if resp.GetResult() != 0 {
					failCount++
				}
				curCost := uint64(time.Now().UnixNano()) - binary.LittleEndian.Uint64(resp.GetUserBuffer())
				if curCost > maxTime {
					maxTime = curCost
				}
				countTime += maxTime
				allcount++
			} else {
				time.Sleep(time.Nanosecond * 1)
			}
			if allcount == tCount*num {
				gch <- countTime / uint64(allcount)
				allcount = 0
				countTime = 0
			}
		}
	}()

	reqs := make([]request.TcaplusRequest, 0, tCount)
	recs := make([]*record.Record, 0, tCount)
	startCh := make([]chan struct{}, tCount)
	for i := 0; i < tCount; i++ {
		req, err := c.NewRequest(ApiConfig.ZoneId, ApiConfig.PBTable, cmdId)
		if err != nil {
			fmt.Println(err)
			return
		}
		rec, err := req.AddRecord(0)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.SetResultFlagForSuccess(byte(flagId))
		reqs = append(reqs, req)
		recs = append(recs, rec)
		startCh[i] = make(chan struct{}, 1)
		go func(i int) {
			tmp := make([]byte, 8)
			for {
				select {
				case <-startCh[i]:
					for n := 0; n < num; n++ {
						begin := time.Now()
						if _, err = recs[i].SetPBData(records[i*num+n]); err != nil {
							fmt.Println(err)
							return
						}
						binary.LittleEndian.PutUint64(tmp, uint64(begin.UnixNano()))
						req.SetUserBuff(tmp)
						if err = c.SendRequest(reqs[i]); err != nil {
							fmt.Println(err)
							return
						}
					}
				}
			}
		}(i)
	}

	var pMaxTime float64
	var pFailCount int
	var pqps float64
	var pavg float64
	var pcostAll float64

	if p != "" {
		f, _ := os.Create(p)
		defer f.Close()
		f.WriteString("cur_time,cost_all(ms),max(ms),qps(w/s),avg(ms),count_all,count_success,count_fail\n")
		timer := time.NewTicker(time.Duration(itv) * time.Second)
		go func() {
			for {
				select {
				case <-timer.C:
					t := fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f,%d,%d,%d\n",
						time.Now().Format("2006/1/2 15:04:05"), pcostAll, pMaxTime, pqps, pavg, tCount*num, tCount*num-pFailCount, pFailCount)
					f.WriteString(t)
				}
			}
		}()
	}

	for {
		maxTime = 0
		failCount = 0
		start := time.Now()
		for i := 0; i < tCount; i++ {
			startCh[i] <- struct{}{}
		}
		avg := <-gch
		cost := time.Now().UnixNano() - start.UnixNano()
		qps := float64(tCount*num) * 1e9 / float64(cost)
		if p != "" {
			pMaxTime = float64(maxTime) / 1e6
			pcostAll = float64(cost) / 1e6
			pqps = qps / 1e4
			pavg = float64(avg) / 1e6
			pFailCount = failCount
		}
		fmt.Printf("insert cost: %.4fms; max: %.4fms; qps: %.4fw/s; avg: %.4fms\n",
			float64(cost)/1e6, float64(maxTime)/1e6, qps/1e4, float64(avg)/1e6)
		fmt.Printf("count all: %d; success: %d; fail: %d\n", tCount*num, tCount*num-failCount, failCount)
		time.Sleep(time.Microsecond * 1000)
	}
}

func TcaplusSyncTest(tCount int, num int, f, p string, itv int) {
	cmdId := cmd.TcaplusApiGetReq
	flagId := 0
	if f == "add" {
		cmdId = cmd.TcaplusApiReplaceReq
		flagId = 0
	} else if f == "mix" {
		cmdId = cmd.TcaplusApiReplaceReq
		flagId = 3
	}

	c := InitPBClient()

	records := make([]*tcaplusservice.GamePlayers, tCount*num)
	for i := 0; i < tCount*num; i++ {
		randseed := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(int64(i)*12345))))
		records[i] = &tcaplusservice.GamePlayers{
			PlayerId:        int64(i) * 12345,
			PlayerName:      string(randseed[:]),
			PlayerEmail:     string(randseed[:]),
			GameServerId:    10,
			LoginTimestamp:  []string{"2019-12-12 15:00:00"},
			LogoutTimestamp: []string{"2019-12-12 16:00:00"},
			IsOnline:        false,
			Pay: &tcaplusservice.Payment{
				PayId:  10101,
				Amount: 1000,
				Method: 1,
			},
		}
	}

	reqs := make([]request.TcaplusRequest, 0, tCount)
	recs := make([]*record.Record, 0, tCount)
	startCh := make([]chan struct{}, tCount)

	maxCost := time.Duration(0)
	avgCost := time.Duration(0)
	ch := make(chan time.Duration, tCount)
	failCount := 0

	for i := 0; i < tCount; i++ {
		req, err := c.NewRequest(ApiConfig.ZoneId, ApiConfig.PBTable, cmdId)
		if err != nil {
			fmt.Println(err)
			return
		}
		rec, err := req.AddRecord(0)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.SetResultFlagForSuccess(byte(flagId))
		reqs = append(reqs, req)
		recs = append(recs, rec)
		startCh[i] = make(chan struct{}, 1)
		go func(i int) {
			for {
				select {
				case <-startCh[i]:
					countCost := time.Duration(0)
					for n := 0; n < num; n++ {
						begin := time.Now()
						if _, err = recs[i].SetPBData(records[i*num+n]); err != nil {
							fmt.Println(err)
							return
						}
						if resp, err := c.Do(reqs[i], 60*time.Second); err != nil || resp.GetResult() != 0 {
							failCount++
							continue
						}
						end := time.Since(begin)
						countCost += end
						if end > maxCost {
							maxCost = end
						}
					}
					ch <- countCost
				}
			}
		}(i)
	}

	var pMaxTime float64
	var pFailCount int
	var pqps float64
	var pavg float64
	var pcostAll float64
	if p != "" {
		f, _ := os.Create(p)
		defer f.Close()
		f.WriteString("cur_time,cost_all(ms),max(ms),qps(w/s),avg(ms),count_all,count_success,count_fail\n")
		timer := time.NewTicker(time.Duration(itv) * time.Second)
		go func() {
			for {
				select {
				case <-timer.C:
					t := fmt.Sprintf("%s,%.4f,%.4f,%.4f,%.4f,%d,%d,%d\n",
						time.Now().Format("2006/1/2 15:04:05"), pcostAll, pMaxTime, pqps, pavg, tCount*num, tCount*num-pFailCount, pFailCount)
					f.WriteString(t)
				}
			}
		}()
	}

	for {
		maxCost = time.Duration(0)
		avgCost = time.Duration(0)
		start := time.Now()
		failCount = 0
		for i := 0; i < tCount; i++ {
			startCh[i] <- struct{}{}
		}
		for i := 0; i < tCount; i++ {
			avgCost += <-ch
		}
		cost := time.Now().UnixNano() - start.UnixNano()
		qps := float64(tCount*num) * 1e9 / float64(cost)
		avg := float64(avgCost.Nanoseconds()) / float64(tCount*num)

		if p != "" {
			pMaxTime = float64(maxCost.Nanoseconds()) / 1e6
			pcostAll = float64(cost) / 1e6
			pqps = qps / 1e4
			pavg = float64(avg) / 1e6
			pFailCount = failCount
		}
		fmt.Printf("insert cost: %.4fms; max: %.4fms; qps: %.4fw/s; avg: %.4fms\n",
			float64(cost)/1e6, float64(maxCost.Nanoseconds())/1e6, qps/1e4, avg/1e6)
		fmt.Printf("count all: %d; success: %d; fail: %d\n", tCount*num, tCount*num-failCount, failCount)
		time.Sleep(time.Microsecond * 1000)
	}
}
