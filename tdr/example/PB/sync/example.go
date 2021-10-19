package main

/*******************************************************************************************************************************************
* author : Tcaplus
* note :本例将演示TcaplusDB PB API的同步调用使用方法, 假定用户已经通过 game_players.proto 在自己的TcaplusDB应用中创建了名为 game_players 的表
创建表格、获取访问点信息的指引请参考 https://cloud.tencent.com/document/product/596/38807。
********************************************************************************************************************************************/

import (
	"fmt"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr"
	"github.com/tencentyun/tcaplusdb-go-sdk/tdr/example/table/tcaplusservice"
	"google.golang.org/protobuf/proto"
	"time"
)

//TcaplusDB GO PB API的连接参数
const (
	//集群访问地址，本地docker版：配置docker部署机器IP, 端口默认:9999, 腾讯云线上环境配置为连接地址IP和端口
	DirUrl = "tcp://x.x.x.x:xxxx"
	//集群接入ID, 默认为3，本地docker版：直接填3，云上版本：根据实际集群接入ID填写
	AppId = 3
	//集群访问密码，本地docker版：登录tcaplusdb web运维平台查看(账号/密码:tcaplus/tcaplus)，业务管理->业务维护->查看pb_app业务对应密码; 云上版本：根据实际申请集群详情页查看
	Signature = "xxxxx"
	//表格组ID，替换为自己创建的表格组ID
	ZoneId = 2
	//表名称
	TableName = "game_players"
)

//声明一个TcaplusDB连接客户端
var client *tcaplus.PBClient

//初始化客户端连接
func initClient() {
	//通过指定接入ID(AppId), 表格组id表表(zoneList), 接入地址(DirUrl), 集群密码(Signature) 参数创建TcaplusClient的对象client
	//通过client对象可以访问集群下的所有大区和表
	//创建表格、获取访问点信息的指引请参考 https://cloud.tencent.com/document/product/596/38807
	client = tcaplus.NewPBClient()
	//设置log配置，log级别默认可设置为ERROR或INFO
	if err := client.SetLogCfg("../cfg/logconf.xml"); err != nil {
		fmt.Println(err.Error())
		return
	}

	zoneList := []uint32{ZoneId}
	zoneTable := make(map[uint32][]string)
	//构造Map对象存储对应表格组下所有的表
	zoneTable[ZoneId] = []string{TableName}
	//建立到对应集群的连接客户端
	err := client.Dial(AppId, zoneList, DirUrl, Signature, 30, zoneTable)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

// 插入记录
func insertRecord() {

	record := &tcaplusservice.GamePlayers{
		PlayerId:        10805514,
		PlayerName:      "Calvin",
		PlayerEmail:     "calvin@test.com",
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
	err := client.Insert(record)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case Insert:")
	fmt.Printf("message:%+v\n", record)
}

// 获取记录
func getRecord() {

	record := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}
	err := client.Get(record)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case Get:")
	fmt.Printf("message:%+v\n", record)
}

// 替换记录（记录不存在则插入）
func replaceRecord() {

	record := &tcaplusservice.GamePlayers{
		PlayerId:        10805514,
		PlayerName:      "Calvin",
		PlayerEmail:     "calvin@test.com",
		GameServerId:    12,
		LoginTimestamp:  []string{"2019-12-12 15:00:00"},
		LogoutTimestamp: []string{"2019-12-12 16:00:00"},
		IsOnline:        false,
		Pay: &tcaplusservice.Payment{
			PayId:  10102,
			Amount: 1002,
			Method: 2,
		},
	}
	err := client.Replace(record)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case Replace:")
	fmt.Printf("message:%+v\n", record)
}

// 修改记录 （记录不存在则报错）
func updateRecord() {

	record := &tcaplusservice.GamePlayers{
		PlayerId:        10805514,
		PlayerName:      "Calvin",
		PlayerEmail:     "calvin@test.com",
		GameServerId:    12,
		LoginTimestamp:  []string{"2019-12-12 15:00:00"},
		LogoutTimestamp: []string{"2019-12-12 16:00:00"},
		IsOnline:        false,
		Pay: &tcaplusservice.Payment{
			PayId:  10104,
			Amount: 1004,
			Method: 4,
		},
	}
	err := client.Update(record)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case Update:")
	fmt.Printf("message:%+v\n", record)
}

// 获取部分value
func fieldGetRecord() {

	record := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}
	err := client.FieldGet(record, []string{"pay", "pay.pay_id"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case FieldGet:")
	fmt.Printf("message:%+v\n", record)
}

// 更新部分value
func fieldUpdateRecord() {

	record := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
		Pay: &tcaplusservice.Payment{
			PayId:  10102,
			Amount: 1002,
		},
	}
	err := client.FieldUpdate(record, []string{"pay.amount", "pay.pay_id"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case FieldUpdate:")
	fmt.Printf("message:%+v\n", record)
}

// 部分value自增
func fieldIncreaseRecord() {

	record := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
		Pay: &tcaplusservice.Payment{
			PayId:  10102,
			Amount: 1002,
		},
	}
	err := client.FieldIncrease(record, []string{"pay.amount", "pay.pay_id"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case FieldIncrease:")
	fmt.Printf("message:%+v\n", record)
}

// 删除记录
func deleteRecord() {

	record := &tcaplusservice.GamePlayers{
		PlayerId:    10805514,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}
	err := client.Delete(record)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case Delete:")
	fmt.Printf("message:%+v\n", record)
}

// 批量获取记录
func batchGetRecord() {

	key := &tcaplusservice.GamePlayers{
		PlayerId:    10805510,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}
	key2 := &tcaplusservice.GamePlayers{
		PlayerId:    10805511,
		PlayerName:  "Calvin",
		PlayerEmail: "calvin@test.com",
	}

	msgs := []proto.Message{key, key2}
	err := client.BatchGet(msgs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case BatchGet:")
	fmt.Printf("message:%+v\n", msgs)
}

// 部分key字段获取记录
func partkeyGetRecord() {

	record := &tcaplusservice.GamePlayers{
		PlayerId:   10805514,
		PlayerName: "Calvin",
	}
	msgs, err := client.GetByPartKey(record, []string{"player_id", "player_name"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case GetByPartKey:")
	fmt.Printf("message:%+v\n", msgs)
}

// 二级索引查询, 需设置索引才能使用
func indexQuery() {

	// 非聚合查询
	query := fmt.Sprintf("select pay.pay_id, pay.amount from game_players where player_id=10805514")
	msgs, _, err := client.IndexQuery(query)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case IndexQuery:")
	fmt.Printf("message:%+v\n", msgs)

	// 聚合查询
	query = fmt.Sprintf("select count(pay) from game_players where player_id=10805514")
	_, res, err := client.IndexQuery(query)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case IndexQuery:")
	fmt.Printf("message:%+v\n", res)
}

// 遍历记录
func traverse() {

	record := &tcaplusservice.GamePlayers{}
	// 遍历时间可能比较长超时时间设长一些
	client.SetDefaultTimeOut(30 * time.Second)
	msgs, err := client.Traverse(record)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case Traverse:")
	fmt.Printf("message:%+v\n", msgs)
}

// 获取表记录总数
func count() {
	count, err := client.GetTableCount("game_players")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Case Count:")
	fmt.Printf("Count:%d\n", count)
}

func main() {
	initClient()
	//insertRecord()
	//getRecord()
	//replaceRecord()
	//updateRecord()
	//batchGetRecord()
	//partkeyGetRecord()
	//fieldGetRecord()
	//fieldUpdateRecord()
	//fieldIncreaseRecord()
	//deleteRecord()

	//batchGetRecord()		// 使用前请插入需要查询的记录
	//partkeyGetRecord()	// 使用前请插入需要查询的记录
	//indexQuery()			// 使用前请设置索引
	//traverse()			// 使用前请先随便插入几条记录
	//count()				// 使用前请先随便插入几条记录
}
