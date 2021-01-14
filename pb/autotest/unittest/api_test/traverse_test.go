package api

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/table/tcaplusservice"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/terror"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/traverser"
	"testing"
	"time"
)

func TestPBTraverse(t *testing.T) {
	client := tools.InitPBSyncClient()
	tra := client.GetTraverser(tools.ZoneId, "game_players")
	defer tra.Stop()
	err := tra.Start()
	if err != nil {
		t.Errorf("start error:%s", err)
	}

	msg := &tcaplusservice.GamePlayers{}
	msg.PlayerId = 233
	msg.PlayerName = "TestPBTraverse"
	msg.PlayerEmail = "dsf"
	for i:=0;i<500;i++{
		msg.PlayerId = 233 * int64(i)
		client.Insert(msg)
		defer func(id int64) {
			msg.PlayerId = 233 * id
			client.Delete(msg)
		}(int64(i))
	}

	count := 0

	for {
		resp, err := client.RecvResponse()
		if err != nil {
			t.Errorf("RecvResponse fail, %s", err.Error())
			return
		} else if resp == nil {
			if tra.State() != traverser.TraverseStateNormal {
				fmt.Println(tra.State())
				break
			} else {
				time.Sleep(time.Microsecond * 10)
				continue
			}
		}

		if err := resp.GetResult(); err != terror.GEN_ERR_SUC {
			t.Errorf("GetResult fail, %d %s", err, terror.GetErrMsg(err))
			return
		}

		count += resp.GetRecordCount()

		for i := 0; i < resp.GetRecordCount(); i++ {
			record, err := resp.FetchRecord()
			if err != nil {
				t.Errorf("FetchRecord fail, %s", err.Error())
				return
			}

			msg := &tcaplusservice.GamePlayers{}
			err = record.GetPBData(msg)
			if err != nil {
				t.Errorf("GetPBData fail, %s", err.Error())
				return
			}

			newJson := tools.StToJson(msg)
			fmt.Println(newJson)
		}
	}

	fmt.Println(count)
}

func TestPBSyncTraverse(t *testing.T) {
	client := tools.InitPBSyncClient()

	msg := &tcaplusservice.GamePlayers{}
	msg.PlayerId = 233
	msg.PlayerName = "TestPBTraverse"
	msg.PlayerEmail = "dsf"
	// 插入500条数据
	for i:=0;i<500;i++{
		msg.PlayerId = 233 * int64(i)
		client.Insert(msg)
		defer func(id int64) {
			msg.PlayerId = 233 * id
			client.Delete(msg)
		}(int64(i))
	}


	// 遍历，参数为定义的 proto message， 返回  message 列表与错误
	table := &tcaplusservice.GamePlayers{}
	client.SetDefaultTimeOut(30*time.Second)
	msgs, err := client.Traverse(table)
	if err != nil {
		t.Errorf("start error:%s", err)
	}

	fmt.Println(len(msgs))
}

func TestPBSyncTraverse2(t *testing.T) {
	client := tools.InitPBSyncClient()

	msg := &tcaplusservice.GamePlayers{}
	msg.PlayerId = 233
	msg.PlayerName = "TestPBTraverse"
	msg.PlayerEmail = "dsf"
	// 插入500条数据
	for i:=0;i<5;i++{
		msg.PlayerId = 233 * int64(i)
		client.Insert(msg)
		defer func(id int64) {
			msg.PlayerId = 233 * id
			client.Delete(msg)
		}(int64(i))
	}


	// 遍历，参数为定义的 proto message， 返回  message 列表与错误
	table := &tcaplusservice.GamePlayers{}
	client.SetDefaultTimeOut(30*time.Second)

	for i:=0;i<10;i++ {
		start := time.Now()
		msgs, err := client.Traverse(table)
		if err != nil {
			t.Errorf("start error:%s", err)
		}

		fmt.Println(len(msgs))
		fmt.Println(time.Since(start))
		time.Sleep(10*time.Millisecond)
	}

}
