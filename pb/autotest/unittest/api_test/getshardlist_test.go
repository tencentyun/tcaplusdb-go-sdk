package api

import (
	"fmt"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/autotest/unittest/tools"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/common"
	"git.code.com/gcloud_storage_group/tcaplus-go-api/protocol/cmd"
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

	fmt.Println(common.CovertToJson(resp.GetTcaplusPackagePtr().Body.GetShardListRes))

}
