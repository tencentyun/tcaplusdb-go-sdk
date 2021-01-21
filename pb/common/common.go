package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/pb/protocol/tcaplus_protocol_cs"
)

const (
	// 每个server处理响应的协程数
	ConfigProcRespRoutineNum = 4

	// 处理响应的协程管道深度
	ConfigProcRespDepth = 1024

	// 处理写请求的协程管道深度，用于合并写请求
	ConfigProcReqDepth = 1024

	// 处理用户同步请求的协程数
	ConfigProcRouterRoutineNum = 4

	// 处理用户同步请求的协程管道深度
	ConfigProcRouterDepth = 1024
)

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
	return fmt.Sprintf("{ Result:%d Magic:%d Version:%d HeadLen:%d BodyLen:%d AsynID:%d Seq:%d Cmd:%d" +
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
