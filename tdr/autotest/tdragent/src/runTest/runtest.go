package runTest

import (
	"errors"
	"fmt"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/tdragent/src/cfg"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/autotest/tdragent/src/logger"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/record"
	"git.code.oa.com/gcloud_storage_group/tcaplus-go-api/terror"
	"math/rand"
	"strconv"
	"time"
)

const (
	TCAPLUS_TEST_TYPE_INT8   = 1
	TCAPLUS_TEST_TYPE_UINT8  = 2
	TCAPLUS_TEST_TYPE_INT16  = 3
	TCAPLUS_TEST_TYPE_UINT16 = 4
	TCAPLUS_TEST_TYPE_INT32  = 5
	TCAPLUS_TEST_TYPE_UINT32 = 6
	TCAPLUS_TEST_TYPE_INT64  = 7
	TCAPLUS_TEST_TYPE_UINT64 = 8
	TCAPLUS_TEST_TYPE_FLOAT  = 9
	TCAPLUS_TEST_TYPE_DOUBLE = 10
	TCAPLUS_TEST_TYPE_STRING = 11
	TCAPLUS_TEST_TYPE_BINARY = 12
)

const (
	TCAPLUS_TEST_GENERIC_RAND_REQ = 0x0FF1
	TCAPLUS_TEST_INSERT_REQ       = 0x0001
	TCAPLUS_TEST_REPLACE_REQ      = 0x0003
	TCAPLUS_TEST_GET_REQ          = 0x0007
	TCAPLUS_TEST_DELETE_REQ       = 0x0009
	TCAPLUS_TEST_UPDATE_REQ       = 0x001d
)

type runFunc func() error

type FieldInfo struct {
	FieldName    string
	FieldTypeInt int
	FieldBuff    interface{}
}

type RecKeyInfo struct {
	SendTime time.Time
	KeyList  []FieldInfo
}

type RunTest struct {
	TcaplusClient *tcaplus.Client // 指定最外层的标签为TcaplusCase
	TcaplusCase   *cfg.TcaplusCase

	//key，value随机变化
	keyIncrease   []int
	valueIncrease []int

	//运行的函数
	runFunc     runFunc
	randRunFunc []runFunc

	//异步id
	asyncId    uint64
	sendKeyMap map[uint64]*RecKeyInfo
	keyCache   map[string]*RecKeyInfo //svr 端存在的key

	//速度控制
	curSpeed int64 //当前速度， 每秒发送包量
	minSpeed int64 //起始速度，每秒最大包量
	maxSpeed int64 //最大速度，每秒最大包量

	onePeriod         time.Duration //一个发送时间片的大小
	onePeriodNeedSend int64         //一个时间片应该发送的包量
	curPeriodSend     int64         //当前时间片已经发送的包量
	secondSliceNum    int64         //1s应该切片的数量，1s切成200份，每份5ms
	curPeriodOverNum  int64         //当前周期过载数量

	//统计
	lastCalTime time.Time
	//当前一秒的统计
	sendNum  int64
	recvNum  int64
	delaySum int64
	aveDelay int64
	maxDelay int64
	overNum  int64
	errNum   int64

	//累计统计
	totalSendNum  int64
	totalRecvNum  int64
	totalDelaySum int64
	totalAveDelay int64
	totalMaxDelay int64
	totalOverNum  int64
	totalErrNum   int64

	//随机字符串
	longRandStr []byte
	lenRandStr  int32
}

func (r *RunTest) init() error {
	r.asyncId = r.TcaplusCase.Head.AsyncID

	//将cmdType转为int
	switch r.TcaplusCase.Head.Cmd.CmdType {
	case "TCAPLUS_TEST_GENERIC_RAND_REQ":
		r.runFunc = nil
		r.randRunFunc = make([]runFunc, 0, len(r.TcaplusCase.Head.Cmd.RandCmd))
		r.TcaplusCase.Head.Cmd.CmdTypeInt = TCAPLUS_TEST_GENERIC_RAND_REQ
		randCmd := make([]int, 0, len(r.TcaplusCase.Head.Cmd.RandCmd))
		for _, subCmd := range r.TcaplusCase.Head.Cmd.RandCmd {
			switch subCmd {
			case "TCAPLUS_TEST_INSERT_REQ":
				randCmd = append(randCmd, TCAPLUS_TEST_INSERT_REQ)
				r.randRunFunc = append(r.randRunFunc, r.TestInsert)
			case "TCAPLUS_TEST_REPLACE_REQ":
				randCmd = append(randCmd, TCAPLUS_TEST_REPLACE_REQ)
				r.randRunFunc = append(r.randRunFunc, r.TestReplace)
			case "TCAPLUS_TEST_GET_REQ":
				randCmd = append(randCmd, TCAPLUS_TEST_GET_REQ)
				r.randRunFunc = append(r.randRunFunc, r.TestGet)
			case "TCAPLUS_TEST_DELETE_REQ":
				randCmd = append(randCmd, TCAPLUS_TEST_DELETE_REQ)
				r.randRunFunc = append(r.randRunFunc, r.TestDelete)
			case "TCAPLUS_TEST_UPDATE_REQ":
				randCmd = append(randCmd, TCAPLUS_TEST_UPDATE_REQ)
				r.randRunFunc = append(r.randRunFunc, r.TestUpdate)
			default:
				logger.ERR("invalid randCmd type %s", subCmd)
				return errors.New("invalid randCmd type")
			}
			r.TcaplusCase.Head.Cmd.RandCmdInt = randCmd
		}
	case "TCAPLUS_TEST_INSERT_REQ":
		r.TcaplusCase.Head.Cmd.CmdTypeInt = TCAPLUS_TEST_INSERT_REQ
		r.runFunc = r.TestInsert

	case "TCAPLUS_TEST_REPLACE_REQ":
		r.TcaplusCase.Head.Cmd.CmdTypeInt = TCAPLUS_TEST_REPLACE_REQ
		r.runFunc = r.TestReplace

	case "TCAPLUS_TEST_GET_REQ":
		r.TcaplusCase.Head.Cmd.CmdTypeInt = TCAPLUS_TEST_GET_REQ
		r.runFunc = r.TestGet

	case "TCAPLUS_TEST_DELETE_REQ":
		r.TcaplusCase.Head.Cmd.CmdTypeInt = TCAPLUS_TEST_DELETE_REQ
		r.runFunc = r.TestDelete

	case "TCAPLUS_TEST_UPDATE_REQ":
		r.TcaplusCase.Head.Cmd.CmdTypeInt = TCAPLUS_TEST_UPDATE_REQ
		r.runFunc = r.TestUpdate

	default:
		logger.ERR("invalid CmdType type %s", r.TcaplusCase.Head.Cmd.CmdType)
		return errors.New("invalid CmdType type")
	}

	//将keyfield type转为int
	for i, keyInfo := range r.TcaplusCase.Body.KeyInfoList {
		switch keyInfo.FieldType {
		case "INT8":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_INT8
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "UINT8":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_UINT8
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "INT16":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_INT16
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "UINT16":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_UINT16
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "INT32":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_INT32
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "UINT32":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_UINT32
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "INT64":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_INT64
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "UINT64":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_UINT64
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "FLOAT":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_FLOAT
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "DOUBLE":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_DOUBLE
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "STRING":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_STRING
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		case "BINARY":
			keyInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_BINARY
			r.TcaplusCase.Body.KeyInfoList[i] = keyInfo
		default:
			logger.ERR("invalid key type %s", keyInfo.FieldType)
			return errors.New("invalid key type")
		}
	}

	//将valuefield转int
	for i, valueInfo := range r.TcaplusCase.Body.ValueInfoList {
		switch valueInfo.FieldType {
		case "INT8":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_INT8
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "UINT8":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_UINT8
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "INT16":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_INT16
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "UINT16":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_UINT16
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "INT32":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_INT32
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "UINT32":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_UINT32
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "INT64":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_INT64
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "UINT64":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_UINT64
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "FLOAT":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_FLOAT
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "DOUBLE":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_DOUBLE
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "STRING":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_STRING
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		case "BINARY":
			valueInfo.FieldTypeInt = TCAPLUS_TEST_TYPE_BINARY
			r.TcaplusCase.Body.ValueInfoList[i] = valueInfo
		default:
			logger.ERR("invalid value type %s", valueInfo.FieldType)
			return errors.New("invalid value type")
		}
	}

	r.asyncId = r.TcaplusCase.Head.AsyncID
	r.keyIncrease = make([]int, len(r.TcaplusCase.Body.KeyInfoList), len(r.TcaplusCase.Body.KeyInfoList))
	r.valueIncrease = make([]int, len(r.TcaplusCase.Body.ValueInfoList), len(r.TcaplusCase.Body.ValueInfoList))
	r.curSpeed = int64(r.TcaplusCase.Head.SpeedControl.StartSpeed)
	r.minSpeed = int64(r.TcaplusCase.Head.SpeedControl.StartSpeed)
	r.maxSpeed = int64(r.TcaplusCase.Head.SpeedControl.SpeedLimit)
	r.lastCalTime = time.Now()
	if r.maxSpeed < r.minSpeed {
		return errors.New("r.maxSpeed < r.minSpeed")
	}
	r.sendKeyMap = make(map[uint64]*RecKeyInfo)
	r.keyCache = make(map[string]*RecKeyInfo, 1000)
	//产生随机数串
	fmt.Println("begin to generate random string")
	r.lenRandStr = 100000000
	r.longRandStr = make([]byte, r.lenRandStr, r.lenRandStr)
	seed := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
	for i := int32(0); i < r.lenRandStr; i++ {
		randomnum := rand.Intn(len(seed))
		r.longRandStr[i] = seed[randomnum]
	}
	fmt.Println("generate random string end")
	return nil
}

func (r *RunTest) Run() {
	if r.TcaplusCase == nil || r.TcaplusClient == nil {
		fmt.Println("err: r.TcaplusCase == nil || r.TcaplusClient == nil")
		return
	}

	if err := r.init(); err != nil {
		fmt.Println(err.Error())
		return
	}

	rand.Seed(time.Now().UnixNano())

	//将1s切成200份, 每份5ms
	r.secondSliceNum = 200
	if r.minSpeed < 200 {
		r.minSpeed = 200
		r.curSpeed = 200
		logger.INFO("minSpeed auto change to 200")
	}
	if r.maxSpeed < 200 {
		r.maxSpeed = 200
		logger.INFO("maxSpeed auto change to 200")
	}
	r.onePeriod = 1 * time.Second / time.Duration(r.secondSliceNum)
	r.onePeriodNeedSend = r.curSpeed / r.secondSliceNum

	//change speed tick 每20s调整下发包速度
	if r.TcaplusCase.Head.SpeedControl.SpeedChangePeriod <= 0 {
		logger.ERR("SpeedChangePeriod invalid")
		return
	}

	fmt.Printf("TcaplusCase: %+v\n", *r.TcaplusCase)
	fmt.Println("")
	fmt.Printf("-->RunTest Info(Speed:%d-%d):\n", r.minSpeed, r.maxSpeed)
	fmt.Println("")

	changeSpeedTick := time.NewTicker(time.Duration(r.TcaplusCase.Head.SpeedControl.SpeedChangePeriod) * time.Second)
	sendTick := time.NewTicker(r.onePeriod)
	//每ms收一次包
	//recvTick := time.NewTicker(time.Millisecond)
	for {
		select {
		//定时器驱动发包
		case <-sendTick.C:
			r.curPeriodSend = 0
			if r.runFunc != nil {
				//固定命令
				for ; r.curPeriodSend < r.onePeriodNeedSend; r.curPeriodSend++ {
					r.asyncId++
					if err := r.runFunc(); err != nil {
						logger.ERR("test runFunc error, %s", err.Error())
					} else {
						r.sendNum++
						r.totalSendNum++
					}
				}
			} else {
				//随机命令
				for ; r.curPeriodSend < r.onePeriodNeedSend; r.curPeriodSend++ {
					r.asyncId++
					if err := r.randRunFunc[rand.Intn(len(r.randRunFunc))](); err != nil {
						logger.ERR("test randRunFunc error, %s", err.Error())
					} else {
						r.sendNum++
						r.totalSendNum++
					}
				}
			}

			//rsp
			r.recvResponse()
			r.staticsInfo()
		//定时进行速度调节
		case <-changeSpeedTick.C:
			r.adjustSpeed()

		//rsp
		//case <-recvTick.C:
		//	r.recvResponse()
		//	r.staticsInfo()
		default:
			if 0 == r.recvResponse() {
				time.Sleep(time.Microsecond * 1)
			}
			r.staticsInfo()
		}
	}
}

func (r *RunTest) staticsInfo() {
	curTime := time.Now()
	diff := curTime.Sub(r.lastCalTime)
	if diff < 1*time.Second {
		return
	}
	r.lastCalTime = curTime

	if r.recvNum > 0 {
		r.aveDelay = r.delaySum / r.recvNum
	}

	if r.totalRecvNum > 0 {
		r.totalAveDelay = r.totalDelaySum / r.totalRecvNum
	}

	//print
	logger.INFO("CurSpeed %d CurSend %d CurRecv %d CurAveDelay %dus CurMaxDelay %dus CurOver %d CurErr %d "+
		"TotalSend %d TotalRecv %d TotalAveDelay %dus TotalMaxDelay %dus TotalOver %d TotalErr %d",
		r.curSpeed, r.sendNum, r.recvNum, r.aveDelay/1000, r.maxDelay/1000, r.overNum, r.errNum,
		r.totalSendNum, r.totalRecvNum, r.totalAveDelay/1000, r.totalMaxDelay/1000, r.totalOverNum, r.totalErrNum)

	fmt.Printf("CurSpeed %d CurSend %d CurRecv %d CurAveDelay %dus CurMaxDelay %dus CurOver %d CurErr %d "+
		"TotalSend %d TotalRecv %d TotalAveDelay %dus TotalMaxDelay %dus TotalOver %d TotalErr %d\n",
		r.curSpeed, r.sendNum, r.recvNum, r.aveDelay/1000, r.maxDelay/1000, r.overNum, r.errNum,
		r.totalSendNum, r.totalRecvNum, r.totalAveDelay/1000, r.totalMaxDelay/1000, r.totalOverNum, r.totalErrNum)

	//reset
	r.sendNum = 0
	r.recvNum = 0
	r.delaySum = 0
	r.aveDelay = 0
	r.maxDelay = 0
	r.overNum = 0
	r.errNum = 0
}

//20s 一次
func (r *RunTest) adjustSpeed() {
	oldSpeed := r.curSpeed
	if r.curPeriodOverNum > r.TcaplusCase.Head.SpeedControl.AllowErrorNum {
		//睡眠加减速 TODO sleep
		r.curSpeed -= r.TcaplusCase.Head.SpeedControl.Step
		if r.curSpeed < r.minSpeed {
			r.curSpeed = r.minSpeed
		}
	} else {
		//加速
		r.curSpeed += r.TcaplusCase.Head.SpeedControl.Step
		if r.curSpeed > r.maxSpeed {
			r.curSpeed = r.maxSpeed
		}
	}
	oldPeriodNeedSend := r.onePeriodNeedSend
	r.onePeriodNeedSend = r.curSpeed / r.secondSliceNum
	logger.INFO("OverNum %d, Change Speed %d -> %d, onePeriodNeedSend %d -> %d", r.curPeriodOverNum, oldSpeed, r.curSpeed,
		oldPeriodNeedSend, r.onePeriodNeedSend)
	r.curPeriodOverNum = 0
}

func (r *RunTest) recvResponse() int {
	count := 0
	for {
		resp, err := r.TcaplusClient.RecvResponse()
		if err != nil {
			logger.ERR("rsp err %s", err.Error())
			return count
		} else if resp == nil {
			//没有包
			return count
		} else {
			//有包
			asyncId := resp.GetAsyncId()
			recKey, exist := r.sendKeyMap[asyncId]
			if !exist {
				logger.ERR("recv invalid asyncId %d", asyncId)
				continue
			}
			delete(r.sendKeyMap, asyncId)
			r.recvNum++
			r.totalRecvNum++
			count++

			//delay
			curTime := time.Now()
			delay := int64(curTime.Sub(recKey.SendTime))
			r.delaySum += delay
			r.totalDelaySum += delay

			//max delay
			if delay > r.maxDelay {
				r.maxDelay = delay
			}
			if delay > r.totalMaxDelay {
				r.totalMaxDelay = delay
			}
			if delay > int64(10*time.Millisecond) {
				logger.INFO("asyncId:%d cost %dus > 10ms", asyncId, delay/1000)
			}

			ret := resp.GetResult()
			if ret != 0 {
				switch ret {
				case terror.SVR_ERR_FAIL_SYSTEM_BUSY, terror.SVR_ERR_FAIL_OVERLOAD, terror.SVR_ERR_FAIL_TIMEOUT,
					terror.SVR_ERR_FAIL_ROUTE, terror.PROXY_ERR_SEND_MSG, terror.PROXY_ERR_REQUEST_ACCESS_CTRL_REJECT:
					r.overNum++
					r.curPeriodOverNum++
					r.totalOverNum++
					logger.ERR("cmd %d rsp overload %s", resp.GetCmd(), terror.GetErrMsg(ret))
				case terror.SVR_ERR_FAIL_RECORD_EXIST:
					logger.ERR("cmd %d rsp already exist %s", resp.GetCmd(), terror.GetErrMsg(ret))
				case terror.TXHDB_ERR_RECORD_NOT_EXIST:
					logger.ERR("cmd %d rsp not exist %s", resp.GetCmd(), terror.GetErrMsg(ret))
				case terror.PROXY_ERR_DIRECT_RESPONSE:
					logger.DEBUG("cmd %d rsp err %s", resp.GetCmd(), terror.GetErrMsg(ret))
				default:
					logger.ERR("cmd %d rsp err %s", resp.GetCmd(), terror.GetErrMsg(ret))
					r.errNum++
					r.totalErrNum++
				}
			} else {
				switch resp.GetCmd() {
				case cmd.TcaplusApiInsertRes, cmd.TcaplusApiReplaceRes:
					//将数据更新到cache
					if len(r.keyCache) < 1000 {
						r.UpdateCache(recKey)
					}
				}
			}
			return count
		}
	}
}

func (r *RunTest) UpdateCache(recKey *RecKeyInfo) {
	//cacheMap
	keyStr := ""
	for _, data := range recKey.KeyList {
		keyStr += fmt.Sprintf("%s%v_", data.FieldName, data.FieldBuff)

	}
	r.keyCache[keyStr] = recKey
}

func (r *RunTest) MakeKeyFromCacheRec(rec *record.Record, recKey *RecKeyInfo) (*RecKeyInfo, error) {
	copyRecKey := &RecKeyInfo{
		KeyList: make([]FieldInfo, len(recKey.KeyList), len(recKey.KeyList)),
	}
	//set key
	for i, keyInfo := range recKey.KeyList {
		//set
		if err := rec.SetKey(keyInfo.FieldName, keyInfo.FieldBuff); err != nil {
			logger.ERR("%s", err.Error())
			return nil, err
		}
		copyRecKey.KeyList[i].FieldName = keyInfo.FieldName
		copyRecKey.KeyList[i].FieldTypeInt = keyInfo.FieldTypeInt
		copyRecKey.KeyList[i].FieldBuff = keyInfo.FieldBuff
	}
	return copyRecKey, nil
}

func (r *RunTest) MakeKey(rec *record.Record) (*RecKeyInfo, error) {
	recKey := &RecKeyInfo{
		KeyList: make([]FieldInfo, len(r.TcaplusCase.Body.KeyInfoList), len(r.TcaplusCase.Body.KeyInfoList)),
	}
	//set key
	for i, keyInfo := range r.TcaplusCase.Body.KeyInfoList {
		recKey.KeyList[i].FieldName = keyInfo.FieldName
		recKey.KeyList[i].FieldTypeInt = keyInfo.FieldTypeInt

		switch keyInfo.FieldTypeInt {
		case TCAPLUS_TEST_TYPE_INT8, TCAPLUS_TEST_TYPE_UINT8:
			k, err := strconv.Atoi(keyInfo.FieldBuff)
			if err != nil {
				logger.ERR("Atoi %s", err.Error())
				return nil, err
			}

			//random
			data := int8(k) + int8(r.keyIncrease[i])
			r.keyIncrease[i] += keyInfo.KeyStep
			if r.keyIncrease[i] > keyInfo.KeyRange {
				r.keyIncrease[i] = 0
			}

			//set
			if err := rec.SetKeyInt8(keyInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return nil, err
			}
			recKey.KeyList[i].FieldBuff = data

		case TCAPLUS_TEST_TYPE_INT16, TCAPLUS_TEST_TYPE_UINT16:
			k, err := strconv.Atoi(keyInfo.FieldBuff)
			if err != nil {
				logger.ERR("Atoi %s", err.Error())
				return nil, err
			}

			//random
			data := int16(k) + int16(r.keyIncrease[i])
			r.keyIncrease[i] += keyInfo.KeyStep
			if r.keyIncrease[i] > keyInfo.KeyRange {
				r.keyIncrease[i] = 0
			}

			//set
			if err := rec.SetKeyInt16(keyInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return nil, err
			}
			recKey.KeyList[i].FieldBuff = data

		case TCAPLUS_TEST_TYPE_INT32, TCAPLUS_TEST_TYPE_UINT32:
			k, err := strconv.ParseInt(keyInfo.FieldBuff, 10, 64)
			if err != nil {
				logger.ERR("ParseInt %s", err.Error())
				return nil, err
			}

			//random
			data := int32(k) + int32(r.keyIncrease[i])
			r.keyIncrease[i] += keyInfo.KeyStep
			if r.keyIncrease[i] > keyInfo.KeyRange {
				r.keyIncrease[i] = 0
			}

			//set
			if err := rec.SetKeyInt32(keyInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return nil, err
			}
			recKey.KeyList[i].FieldBuff = data

		case TCAPLUS_TEST_TYPE_INT64, TCAPLUS_TEST_TYPE_UINT64:
			k, err := strconv.ParseInt(keyInfo.FieldBuff, 10, 64)
			if err != nil {
				logger.ERR("ParseInt %s", err.Error())
				return nil, err
			}

			//random
			data := int64(k) + int64(r.keyIncrease[i])
			r.keyIncrease[i] += keyInfo.KeyStep
			if r.keyIncrease[i] > keyInfo.KeyRange {
				r.keyIncrease[i] = 0
			}

			//set
			if err := rec.SetKeyInt64(keyInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return nil, err
			}
			recKey.KeyList[i].FieldBuff = data

		case TCAPLUS_TEST_TYPE_FLOAT:
			k, err := strconv.ParseFloat(keyInfo.FieldBuff, 32)
			if err != nil {
				logger.ERR("ParseInt %s", err.Error())
				return nil, err
			}

			//random
			data := float32(k) + float32(r.keyIncrease[i])
			r.keyIncrease[i] += keyInfo.KeyStep
			if r.keyIncrease[i] > keyInfo.KeyRange {
				r.keyIncrease[i] = 0
			}

			//set
			if err := rec.SetKeyFloat32(keyInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return nil, err
			}
			recKey.KeyList[i].FieldBuff = data

		case TCAPLUS_TEST_TYPE_DOUBLE:
			k, err := strconv.ParseFloat(keyInfo.FieldBuff, 64)
			if err != nil {
				logger.ERR("ParseInt %s", err.Error())
				return nil, err
			}

			//random
			data := float64(k) + float64(r.keyIncrease[i])
			r.keyIncrease[i] += keyInfo.KeyStep
			if r.keyIncrease[i] > keyInfo.KeyRange {
				r.keyIncrease[i] = 0
			}

			//set
			if err := rec.SetKeyFloat64(keyInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return nil, err
			}
			recKey.KeyList[i].FieldBuff = data

		case TCAPLUS_TEST_TYPE_STRING:
			data := keyInfo.FieldBuff

			if r.keyIncrease[i] != 0 {
				data = fmt.Sprintf("%s%09d", data, r.keyIncrease[i])
			}
			r.keyIncrease[i] += keyInfo.KeyStep
			if r.keyIncrease[i] > keyInfo.KeyRange {
				r.keyIncrease[i] = 0
			}

			//set
			if err := rec.SetKeyStr(keyInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return nil, err
			}
			recKey.KeyList[i].FieldBuff = data

		case TCAPLUS_TEST_TYPE_BINARY:
			data := keyInfo.FieldBuff
			if r.keyIncrease[i] != 0 {
				data = fmt.Sprintf("%s%09d", data, r.keyIncrease[i])
			}

			r.keyIncrease[i] += keyInfo.KeyStep
			if r.keyIncrease[i] > keyInfo.KeyRange {
				r.keyIncrease[i] = 0
			}

			//set
			if err := rec.SetKeyBlob(keyInfo.FieldName, []byte(data)); err != nil {
				logger.ERR("%s", err.Error())
				return nil, err
			}
			recKey.KeyList[i].FieldBuff = data

		default:
			logger.ERR("not support type %d", keyInfo.FieldTypeInt)
			return nil, errors.New("not support type")
		}
	}
	return recKey, nil
}

func (r *RunTest) MakeValue(rec *record.Record) error {
	//set key
	for i, valueInfo := range r.TcaplusCase.Body.ValueInfoList {
		switch valueInfo.FieldTypeInt {
		case TCAPLUS_TEST_TYPE_INT8, TCAPLUS_TEST_TYPE_UINT8:
			k, err := strconv.Atoi(valueInfo.FieldBuff)
			if err != nil {
				logger.ERR("Atoi %s", err.Error())
				return err
			}

			//random
			data := int8(k) + int8(r.valueIncrease[i])
			r.valueIncrease[i] += valueInfo.ValueStep
			if r.valueIncrease[i] > valueInfo.ValueRange {
				r.valueIncrease[i] = 0
			}

			//set
			if err := rec.SetValueInt8(valueInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return err
			}

		case TCAPLUS_TEST_TYPE_INT16, TCAPLUS_TEST_TYPE_UINT16:
			k, err := strconv.Atoi(valueInfo.FieldBuff)
			if err != nil {
				logger.ERR("Atoi %s", err.Error())
				return err
			}

			//random
			data := int16(k) + int16(r.valueIncrease[i])
			r.valueIncrease[i] += valueInfo.ValueStep
			if r.valueIncrease[i] > valueInfo.ValueRange {
				r.valueIncrease[i] = 0
			}

			//set
			if err := rec.SetValueInt16(valueInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return err
			}

		case TCAPLUS_TEST_TYPE_INT32, TCAPLUS_TEST_TYPE_UINT32:
			k, err := strconv.ParseInt(valueInfo.FieldBuff, 10, 64)
			if err != nil {
				logger.ERR("ParseInt %s", err.Error())
				return err
			}

			//random
			data := int32(k) + int32(r.valueIncrease[i])
			r.valueIncrease[i] += valueInfo.ValueStep
			if r.valueIncrease[i] > valueInfo.ValueRange {
				r.valueIncrease[i] = 0
			}

			//set
			if err := rec.SetValueInt32(valueInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return err
			}

		case TCAPLUS_TEST_TYPE_INT64, TCAPLUS_TEST_TYPE_UINT64:
			k, err := strconv.ParseInt(valueInfo.FieldBuff, 10, 64)
			if err != nil {
				logger.ERR("ParseInt %s", err.Error())
				return err
			}

			//random
			data := int64(k) + int64(r.valueIncrease[i])
			r.valueIncrease[i] += valueInfo.ValueStep
			if r.valueIncrease[i] > valueInfo.ValueRange {
				r.valueIncrease[i] = 0
			}

			//set
			if err := rec.SetValueInt64(valueInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return err
			}

		case TCAPLUS_TEST_TYPE_FLOAT:
			k, err := strconv.ParseFloat(valueInfo.FieldBuff, 32)
			if err != nil {
				logger.ERR("ParseInt %s", err.Error())
				return err
			}

			//random
			data := float32(k) + float32(r.valueIncrease[i])
			r.valueIncrease[i] += valueInfo.ValueStep
			if r.valueIncrease[i] > valueInfo.ValueRange {
				r.valueIncrease[i] = 0
			}

			//set
			if err := rec.SetValueFloat32(valueInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return err
			}
		case TCAPLUS_TEST_TYPE_DOUBLE:
			k, err := strconv.ParseFloat(valueInfo.FieldBuff, 64)
			if err != nil {
				logger.ERR("ParseInt %s", err.Error())
				return err
			}

			//random
			data := float64(k) + float64(r.valueIncrease[i])
			r.valueIncrease[i] += valueInfo.ValueStep
			if r.valueIncrease[i] > valueInfo.ValueRange {
				r.valueIncrease[i] = 0
			}

			//set
			if err := rec.SetValueFloat64(valueInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return err
			}

		case TCAPLUS_TEST_TYPE_STRING:
			//字符串长度为[ValueStep, ValueRange)
			data := valueInfo.FieldBuff
			if valueInfo.ValueStep == 0 && valueInfo.ValueRange > valueInfo.ValueStep {
				data = r.RandStr(valueInfo.ValueRange)
			} else if valueInfo.ValueStep > 0 || valueInfo.ValueRange > valueInfo.ValueStep {
				strLen := valueInfo.ValueStep + rand.Intn(valueInfo.ValueRange-valueInfo.ValueStep)
				data = r.RandStr(strLen)
			}
			//set
			if err := rec.SetValueStr(valueInfo.FieldName, data); err != nil {
				logger.ERR("%s", err.Error())
				return err
			}
		case TCAPLUS_TEST_TYPE_BINARY:
			//字符串长度为[ValueStep, ValueRange)
			data := valueInfo.FieldBuff
			if valueInfo.ValueStep == 0 && valueInfo.ValueRange > valueInfo.ValueStep {
				data = r.RandStr(valueInfo.ValueRange)
			} else if valueInfo.ValueStep > 0 || valueInfo.ValueRange > valueInfo.ValueStep {
				strLen := valueInfo.ValueStep + rand.Intn(valueInfo.ValueRange-valueInfo.ValueStep)
				data = r.RandStr(strLen)
			}
			//set
			if err := rec.SetValueBlob(valueInfo.FieldName, []byte(data)); err != nil {
				logger.ERR("%s", err.Error())
				return err
			}

		default:
			logger.ERR("not support type %d", valueInfo.FieldTypeInt)
			return errors.New("not support type")
		}
	}
	return nil
}

func (r *RunTest) RandStr(strLen int) string {
	randomnum := time.Now().UnixNano() % int64((r.lenRandStr - int32(strLen)))
	return string(r.longRandStr[randomnum : randomnum+int64(strLen)])
	//ret := make([]byte, strLen, strLen)
	//randomnum := rand.Int31n(r.lenRandStr)
	//if randomnum+int32(strLen) < r.lenRandStr {
	//	copy(ret, r.longRandStr[randomnum:])
	//} else {
	//	copy(ret[0:r.lenRandStr-randomnum], r.longRandStr[randomnum:])
	//	copy(ret[r.lenRandStr-randomnum:strLen], r.longRandStr[0:r.lenRandStr-randomnum])
	//}
	//return string(ret)
}
