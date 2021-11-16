package api

import (
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/autotest/unittest/tools"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/cmd"
	"testing"
)

func TestShardListSuccess(t *testing.T) {
	client, req := tools.InitPBClientAndReqWithTableName(cmd.TcaplusApiGetShardListReq, "game_players")

	pkg := req.GetTcaplusPackagePtr()
	pkg.Body.GetShardListReq.BeginIndex = -1
	pkg.Body.GetShardListReq.EndIndex = -1

	err := client.SendRequest(req)
	if err != nil {
		t.Errorf("SendRequest fail, %s", err.Error())
		return
	}

	//recv resp
	resp, err := tools.RecvResponse(client)
	if err != nil {
		t.Errorf("recvResponse fail, %s", err.Error())
		return
	}

	if resp.GetResult() != 0 {
		t.Errorf("resp.GetResult() != 0")
		return
	}
	//fmt.Println(common.CovertToJson(resp.GetTcaplusPackagePtr().Body.GetShardListRes))
}
