package main

import (
	"errors"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/logger"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/response"
	"time"
)

const (
	AppId     = uint64(2)
	ZoneId    = uint32(3)
	DirUrl    = "tcp://x.x.x.x:xxxx"
	Signature = "xxxxxxxxxxxxxxxxx"
	TableName = "service_info"
)

var client *tcaplus.Client
var respChan chan response.TcaplusResponse

//阻塞接收，需要一个channel不停RecvResponse，使用channel将响应传出
func recvResponse(client *tcaplus.Client) (response.TcaplusResponse, error) {
	//5s超时
	timeOutChan := time.After(5 * time.Second)
	for {
		select {
		case <-timeOutChan:
			return nil, errors.New("5s timeout")
		case res := <-respChan:
			return res, nil
		}
	}
}

func main() {
	client = tcaplus.NewClient()
	//日志配置，不配置则debug打印到控制台
	if err := client.SetLogCfg("./logconf.xml"); err != nil {
		fmt.Println(err.Error())
		return
	}
	//client连接tcaplus
	err := client.Dial(AppId, []uint32{ZoneId}, DirUrl, Signature, 60)
	if err != nil {
		fmt.Printf("init failed %v\n", err.Error())
		return
	}
	fmt.Printf("Dial finish\n")

	respChan = make(chan response.TcaplusResponse)
	go func() {
		for {
			// resp err 均为 nil 说明响应池中没有任何响应
			resp, err := client.RecvResponse()
			if err != nil {
				logger.ERR("RecvResponse error:%s", err)
				continue
			} else if resp == nil {
				time.Sleep(time.Microsecond * 5)
				continue
			}
			// 找到对应的响应
			respChan <- resp
		}
	}()

	getPartKeyExample()
	deleteByPartKeyExample()
	insertExample()
	getExample()
	updateExample()
	replaceExample()
	deleteExample()
}
