# Tcaplus Go API 3.46.0
[TOC]
## 1 使用准备
本API是Tcaplus API的Go封装，支持Generic表的增删改查
### 1.1 各代码目录
* pack为打包脚本
* example为示例
* aurotest为测试工具
* 其他目录为Tcaplus Go API的源码
* vendor 中为依赖库的代码，需要使用git submodule init和git submodule update拉取

### 1.2 编译example
example/generic_table展示了对service_info表的insert get replace update delete操作
1. Go环境安装https://golang.org/
2. 将service_info.xml加入Tcaplus
3. 修改main.go的开头的AppId ZoneId DirUrl Signature为相应的Tcaplus配置信息
4. make之后执行

### 1.3 打包脚本
pack/pack.sh展示了对源码及依赖库的打包，方便用户移植到无法使用go mod的场景
1. cd pack && sh pack.sh

## 2 API的使用
### 2.1 vendor方式使用
Tcaplus API的依赖库及其源码都在打包后的src/vendor目录下，用户只需将vendor放入自己的工程目录即可使用Tcaplus Go API的接口

vendor依赖介绍：
* git.code.oa.com/gcloud_storage_group/tcaplus-go-api是Tcaplus Go API源码
* git.code.oa.com/tsf4g/TdrCodeGen是tdr工具，可将tdr的xml转换为go源码
* git.code.oa.com/tsf4g/tdrcom是tdr go源码打解包的依赖库
* go.uber.org/zap是日志库
* github.com/natefinch/lumberjack是日志文件切割库

### 2.2 mod 方式使用
mod 模式需要在能连内网及公网环境下使用
* 在工程中建立go.mod
* 开启module模式
* 执行命令go clean --modcache 
* 执行命令 go mod edit -require="git.code.oa.com/gcloud_storage_group/tcaplus-go-api@v0.1.0"
* 出现tlinux无法download的错误可以参考:[Golang git.code.oa.com 的 go get、go mod 踩坑之旅](http://km.oa.com/group/29073/articles/show/376902?kmref=search&from_page=1&no=1#-%20%E9%94%99%E8%AF%AF-x509-%20certificate%20signed%20by%20unknown%20authority)

## 3 接口使用步骤
对表中的record的操作有两套接口:
* 一套使用SetKey SetValue接口对record赋值，由用户指定key字段和value字段的内容，响应消息只能通过GetKey，GetValue接口读取
* 另一套使用SetData接口对record赋值，用户赋值Tdr结构体，SetData通过反射对record赋值，响应消息只能通过GetData接口读取

### 3.1 SetKey和SetValue方式使用
1 通过tcaplus.NewClient()创建一个tcaplus客户端指针
```
client := tcaplus.NewClient()
```
2 指定操作表的AppId， ZoneIdList， DirUrl，Signature，Timeout(秒)连接tcaplus
```
err := client.Dial(2, []uint32{3,4}, "tcp://x.x.x.x:9999", "xxxx",60)
if err != nil {
        log.ERR("dial failed %s", err.Error())
        return
}

```

3 指定zoneId，表名，命令字，client.NewRequest创建一个请求
```
req, err := client.NewRequest(3, "service_info", cmd.TcaplusApiInsertReq)
if err != nil {
        log.ERR("NewRequest TcaplusApiInsertReq failed %s", err.Error())
        return
}
```

4 req.AddRecord为request添加一条记录record，（index为list表的记录所在编号，generic不支持设为 0 即可）
```
rec, err := req.AddRecord(0)
if err != nil {
        log.ERR("AddRecord failed %s", err.Error())
        return
}
```

5 通过record的SetKey和SetValue接口对记录进行赋值
```
err := rec.SetKeyInt8("keyName", int8(1))
if err != nil {
        log.ERR("SetKeyInt8 failed %s", err.Error())
        return
}

err := rec.SetValueInt8("valueName", int8(1))
if err != nil {
        log.ERR("SetKeyInt8 failed %s", err.Error())
        return
}
```

6 client.SendRequest将请求发送出去
```
if err := client.SendRequest(req); err != nil {
        log.ERR("SendRequest failed %s", err.Error())
        return
}
```

7 client.RecvResponse为异步接收请求响应接口，通过如下方式可阻塞接收响应
```
func recvResponse(client *tcaplus.Client) (response.TcaplusResponse, error){
        //5s超时
        timeOutChan := time.After(5 * time.Second)
        for {
                select {
                case <-timeOutChan:
                        return nil, errors.New("5s timeout")
                default:
                        resp,err := client.RecvResponse()
                        if err != nil {
                                return nil, err
                        } else if resp == nil {
                                time.Sleep(time.Microsecond * 1)
                        } else {
                                return  resp, nil
                        }
                }
        }
}
```

8 操作response的GetResult获取响应结果，GetRecordCount,FetchRecord获取响应消息中的记录record，通过record的GetKey，GetValue接口获取响应记录的字段信息
```
tcapluserr := resp.GetResult()
if tcapluserr != 0 {
	fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
	return
}

for i := 0; i < resp.GetRecordCount(); i++ {
                record, err := resp.FetchRecord()
                if err != nil {
                        log.ERR("FetchRecord failed %s", err.Error())
                        return
                }
                
                keyNameData, err := rec.GetKeyInt8("keyName")
                if err != nil {
                        log.ERR("GetKeyInt8 failed %s", err.Error())
                        return
                }
                
                valueNameData, err := rec.GetValueInt8("valueName")
                if err != nil {
                        log.ERR("GetValueInt8 failed %s", err.Error())
                        return
                }
        }

```

### 3.2 TDR SetData方式使用

1 将要操作的tdr的表xml转换成GO源码
```
cd vendor/git.code.oa.com/tsf4g/TdrCodeGen/
python tdr.py table.xml
得到相应表的go源码目录table/table.go
将table放到自己的go的工程目录即可使用
```

2 通过tcaplus.NewClient()创建一个tcaplus客户端指针
```
client := tcaplus.NewClient()
```

3 指定操作表的AppId， ZoneIdList， DirUrl，Signature，Timeout(秒)连接tcaplus
```
err := client.Dial(2, []uint32{3,4}, "tcp://x.x.x.x:9999", "xxxx",60)
if err != nil {
        log.ERR("dial failed %s", err.Error())
        return
}

```

4 指定zoneId，表名，命令字，client.NewRequest创建一个请求
```
req, err := client.NewRequest(3, "service_info", cmd.TcaplusApiInsertReq)
if err != nil {
        log.ERR("NewRequest TcaplusApiInsertReq failed %s", err.Error())
        return
}
```

5 req.AddRecord为request添加一条记录record，（index为list表的记录所在编号，generic不支持设为 0 即可）
```
rec, err := req.AddRecord(0)
if err != nil {
        log.ERR("AddRecord failed %s", err.Error())
        return
}
```

6 通过第一步中tdr的go源码的New接口创建一个结构体并且赋值,通过SetData接口对record进行赋值
```
data := service_info.NewService_Info()
data.Gameid = "dev"
data.Envdata = "oa"
data.Name = "com"
data.Filterdata = time.Now().Format("2006-01-02T15:04:05.000000Z")
data.Updatetime = uint64(time.Now().UnixNano())
data.Inst_Max_Num = 2
data.Inst_Min_Num = 3
//数组类型为slice需要准确赋值长度，与refer保持一致
data.Routeinfo_Len = uint32(len(route))
data.Routeinfo = []byte("test")

//将tdr的数据设置到请求的记录中
if err := rec.SetData(data); err != nil {
        log.ERR("SetData failed %v", err.Error())
        return
}

```

7 client.SendRequest将请求发送出去
```
if err := client.SendRequest(req); err != nil {
        log.ERR("SendRequest failed %s", err.Error())
        return
}
```

8 client.RecvResponse为异步接收请求响应接口，通过如下方式可阻塞接收响应
```
func recvResponse(client *tcaplus.Client) (response.TcaplusResponse, error){
        //5s超时
        timeOutChan := time.After(5 * time.Second)
        for {
                select {
                case <-timeOutChan:
                        return nil, errors.New("5s timeout")
                default:
                        resp,err := client.RecvResponse()
                        if err != nil {
                                return nil, err
                        } else if resp == nil {
                                time.Sleep(time.Microsecond * 1)
                        } else {
                                return  resp, nil
                        }
                }
        }
}
```

9 操作response的GetResult获取响应结果，GetRecordCount,FetchRecord获取响应消息中的记录record，通过record的GetData接口获取响应记录
```
tcapluserr := resp.GetResult()
if tcapluserr != 0 {
	fmt.Printf("response ret errCode: %d, errMsg: %s", tcapluserr, terror.GetErrMsg(tcapluserr))
	return
}

for i := 0; i < resp.GetRecordCount(); i++ {
                record, err := resp.FetchRecord()
                if err != nil {
                        log.ERR("FetchRecord failed %s", err.Error())
                        return
                }
                
                //通过GetData获取响应记录
                data := service_info.NewService_Info()
                if err := record.GetData(data); err != nil {
                        log.ERR("record.GetData failed %s", err.Error())
                        return
                }
        }
```

## 4 接口列表
```
支持命令：

//Generic表插入请求
TcaplusApiInsertReq = 0x0001

//Generic表替换/插入请求
TcaplusApiReplaceReq = 0x0003

//Generic表增量更新请求
TcaplusApiIncreaseReq = 0x0005

//Generic表单条查询请求
TcaplusApiGetReq = 0x0007

//Generic表删除请求
TcaplusApiDeleteReq = 0x0009

//Generic表删除应答
TcaplusApiDeleteRes = 0x000a

//Generic表更新请求
TcaplusApiUpdateReq = 0x001d

//List表查询所有元素请求
TcaplusApiListGetAllReq = 0x000b

//Generic表批量查询请求
TcaplusApiBatchGetReq = 0x0017

//Generic表按索引查询请求
TcaplusApiGetByPartkeyReq  = 0x0019

//Generic表按索引更新请求
TcaplusApiUpdateByPartkeyReq = 0x004d

//Generic表按索引删除请求
TcaplusApiDeleteByPartkeyReq = 0x004f

```
### 4.1 Client接口
#### 4.1.1 创建tcaplus client
```
/**
   @brief 创建一个tcaplus api客户端
   @retval 返回客户端指针
**/
func NewClient() *Client 
```
#### 4.1.2 日志配置接口
创建client之后，立刻调用
```
/**
   @brief                   设置API日志配置文件全路径log.conf(json格式，example下有示例)，请在client.Dial之前调用
   @param [IN] cfgPath      日志配置文件全路径log.conf
   @retval                  错误码
   @note                    Api日志默认使用的zap，用户也可自行实现日志接口logger.LogInterface，调用SetLogger进行设置
**/
func (c *Client) SetLogCfg(cfgPath string) error 

/**
   @brief                   自定义API日志接口,用户实现logger.LogInterface日志接口，日志将打印到用户的日志接口中，请在client.Dial之前调用
   @param [IN] handle       logger.LogInterface类型的日志接口
   @retval                  错误码
**/
func (c *Client) SetLogger(handle logger.LogInterface) 
```

#### 4.1.3 连接tcaplus接口
```
/**
   @brief 连接tcaplue函数
   @param [IN] appId         appId，在网站注册相应服务以后，你可以得到该appId
   @param [IN] zoneList      需要操作表的区服ID列表，操作的表在多个不同的zone，填zoneId列表；操作的表在一个zone，zone列表填一个zoneId
   @param [IN] signature     签名/密码，在网站注册相应服务以后，你可以得到该字符串
   @param [IN] dirUrl        目录服务器的url，形如"tcp://172.25.40.181:10600"
   @param [IN] timeout       second, 连接所有表对应的tcaplus proxy服务器。若所有的proxy连通且鉴权通过，则立即返回成功；
                                若到达超时时间，只要有一个proxy连通且鉴权通过，也会返回成功；否则返回超时错误。
   @retval                   错误码
**/
func (c *Client) Dial(appId uint64, zoneList []uint32, dirUrl string, signature string, timeout uint32) error
```

#### 4.1.4 创建tcaplus请求
```
/**
	@brief 创建指定分区指定表的tcaplus请求
	@param [IN] zoneId              区服ID
	@param [IN] tableName           表名
	@param [IN] cmd                 命令字(cmd包中cmd.TcaplusApiGetReq等)
	@retval request.TcaplusRequest  tcaplus请求
	@retval error                   错误码
*/
func (c *Client) NewRequest(zoneId uint32, tableName string, cmd int) (request.TcaplusRequest, error) 
```

#### 4.1.5 发送tcaplus请求
```
/**
	@brief 发送tcaplus请求
	@param [IN] req       tcaplus请求
	@retval error         错误码
*/
func (c *Client) SendRequest(req request.TcaplusRequest) error
```

#### 4.1.6 异步接收tcaplus响应
```
/**
    @brief 异步接收tcaplus响应
    @retval response.TcaplusResponse tcaplus响应
    @retval error 错误码
            error nil，response nil 成功但当前无响应消息
            error nil, response 非nil，成功获取响应消息
            error 非nil，接收响应出错
*/
func (c *Client) RecvResponse() (response.TcaplusResponse, error)
```

#### 4.1.7 发送tcaplus同步请求并接受响应
```
/**
    @brief 发送tcaplus同步请求并接受响应
    @param [IN] req tcaplus请求
    @param [IN] timeout 超时时间
    @retval response.TcaplusResponse tcaplus响应
    @retval error 错误码
            error nil，response nil 成功但当前无响应消息
            error nil, response 非nil，成功获取响应消息
            error 非nil，接收响应出错
**/
func (c *Client) Do(req request.TcaplusRequest, timeout time.Duration) (response.TcaplusResponse, error)
```

### 4.2 TcaplusRequest接口
#### 4.2.1 添加记录
```
/**
  @brief  向请求中添加一条记录。
  @param [IN] index         用于List操作(目前不支持)，通常>=0，表示该Record在所属List中的Index；对于Generic操作，index无意义，设0即可
  @retval record.Record     返回记录指针
  @retval error   			错误码
**/
AddRecord(index int32) (*record.Record, error)
```
#### 4.2.2 设置请求异步ID
```
/**
    @brief  设置请求的异步事务ID，api会将其值不变地通过对应的响应消息带回来
    @param  [IN] asyncId  请求对应的异步事务ID
**/
SetAsyncId(id uint64)
```
#### 4.2.3 设置版本校验规则
```
/**
    @brief  设置记录版本的检查类型，用于乐观锁
    @param [IN] type   版本检测类型，取值可以为(policy包中):
                        CheckDataVersionAutoIncrease: 表示检测记录版本号，只有当record.SetVersion函数传入的参数version的值>0,并且版本号与服务器端的版本号相同时，请求才会成功同时在服务器端该版本号会自增1；如果record.SetVersion的version <=0，则仍然表示不关心版本号
                        NoCheckDataVersionOverwrite: 表示不检测记录版本号。当record.SetVersion函数传入的参数version的值>0,覆盖服务端的版本号；如果record.SetVersion的version <=0，则仍然表示不关心版本号
                        NoCheckDataVersionAutoIncrease: 表示不检测记录版本号，将服务器端的数据记录版本号自增1，若服务器端新写入数据记录则新写入的数据记录的版本号为1
    @retval error      错误码
    @note 此函数适合Replace, Update操作
**/
SetVersionPolicy(p uint8) error
```
#### 4.2.4 设置响应标志
```
/**
    @brief  设置响应标志。主要用于Generic表的insert、replace、update、delete操作。
    @param  [IN] flag  请求标志:
                            0表示: 只需返回操作执行成功与否
                            1表示: 操作成功，响应返回与请求字段一致
                            2表示: 操作成功，响应返回变更记录的所有字段最新数据
                            3表示: 操作成功，响应返回变更记录的所有字段旧数据
    @retval error      错误码
**/
SetResultFlag(flag int) error
```
#### 4.2.5 部分字段查询和更新
```
/**
    @brief 设置需要查询或更新的Value字段名称列表，即部分Value字段查询和更新，可用于get、replace、update操作。
    @param [IN] valueNameList   需要查询或更新的字段名称列表
    @retval error               错误码
    @note  在使用该函数设置字段名时，字段名只能包含value字段名，不能包含key字段名；对于数组类型的字段，refer字段和数组字段要同时设置或者同时不设置，否则容易数据错乱
**/
SetFieldNames(valueNameList []string) error
```
#### 4.2.6 设置用户缓存
```
/**
    @brief 设置用户缓存，响应消息将携带返回
    @param [IN] userBuffer  用户缓存
    @retval error           错误码
**/
SetUserBuff(userBuffer []byte) error
```

### 4.3 Record接口
#### 4.3.1 SetKey/SetValue和GetKey/GetValue接口
通过KV接口，对记录进行赋值(请求消息)和获取(响应消息)
```
/**
    @brief  key字段内容设置
    @param  [in] name   字段名称，最大长度32
    @param  [in] data   字段内容
    @retval error       错误码
*/
func (r *Record) SetKeyInt8(name string, data int8) error
func (r *Record) SetKeyInt16(name string, data int16) error 
func (r *Record) SetKeyInt32(name string, data int32) error 
func (r *Record) SetKeyInt64(name string, data int64) error 
func (r *Record) SetKeyFloat32(name string, data float32) error 
func (r *Record) SetKeyFloat64(name string, data float64) error 
func (r *Record) SetKeyStr(name string, data string) error 
func (r *Record) SetKeyBlob(name string, data []byte) error 

/**
    @brief  value字段内容设置
    @param  [in] name   字段名称，最大长度32
    @param  [in] data   字段内容
    @retval error       错误码
*/
func (r *Record) SetValueInt8(name string, data int8) error
func (r *Record) SetValueInt16(name string, data int16) error 
func (r *Record) SetValueInt32(name string, data int32) error 
func (r *Record) SetValueInt64(name string, data int64) error 
func (r *Record) SetValueFloat32(name string, data float32) error 
func (r *Record) SetValueFloat64(name string, data float64) error 
func (r *Record) SetValueStr(name string, data string) error 
func (r *Record) SetValueBlob(name string, data []byte) error

/**
    @brief  key字段内容获取
    @param  [in] name   字段名称，最大长度32
    @retval data        字段内容
    @retval error       错误码
*/
func (r *Record) GetKeyInt8(name string) (int8, error)
func (r *Record) GetKeyInt16(name string) (int16, error)
func (r *Record) GetKeyInt32(name string) (int32, error)
func (r *Record) GetKeyInt64(name string) (int64, error)
func (r *Record) GetKeyFloat32(name string) (float32, error)
func (r *Record) GetKeyFloat64(name string) (float64, error) 
func (r *Record) GetKeyStr(name string) (string, error) 
func (r *Record) GetKeyBlob(name string) ([]byte, error)

/**
    @brief  value字段内容获取
    @param  [in] name   字段名称，最大长度32
    @retval data        字段内容
    @retval error       错误码
*/
func (r *Record) GetValueInt8(name string) (int8, error)
func (r *Record) GetValueInt16(name string) (int16, error)
func (r *Record) GetValueInt32(name string) (int32, error)
func (r *Record) GetValueInt64(name string) (int64, error)
func (r *Record) GetValueFloat32(name string) (float32, error)
func (r *Record) GetValueFloat64(name string) (float64, error) 
func (r *Record) GetValueStr(name string) (string, error) 
func (r *Record) GetValueBlob(name string) ([]byte, error)
```
#### 4.3.2 SetData和GetData接口
通过TDR结构体，对记录进行赋值(请求消息)和获取(响应消息)
```
/**
	@brief  基于TDR描述设置record数据
	@param [IN] data  基于TDR描述record接口数据，tdr的xml通过工具生成的go结构体，包含的TdrTableSt接口的一系列方法
	@retval error     错误码
*/
func (r *Record) SetData(data TdrTableSt) error

/**
	@brief  基于TDR描述读取record数据
	@param [IN] data   基于TDR描述record接口数据，tdr的xml通过工具生成的go结构体，包含的TdrTableSt接口的一系列方法
	@retval error      错误码
**/
func (r *Record) GetData(data TdrTableSt) error
```
#### 4.3.3 设置记录版本号
```
/**
    @brief  设置记录版本号
    @param [IN] v     数据记录的版本号:  <=0 表示不关注版本号不关心版本号。具体含义如下。
                当CHECKDATAVERSION_AUTOINCREASE时: 表示检测记录版本号。
					如果Version的值<=0,则仍然表示不关心版本号不关注版本号；
					如果Version的值>0，那么只有当该版本号与服务器端的版本号相同时，
					Replace, Update, Increase, ListAddAfter, ListDelete, ListReplace, ListDeleteBatch操作才会成功同时在服务器端该版本号会自增1。
                当NOCHECKDATAVERSION_OVERWRITE时: 表示不检测记录版本号。
					如果Version的值<=0,则会把版本号1写入服务端的数据记录版本号(服务器端成功写入的数据记录的版本号最少为1)；
					如果Version的值>0，那么会把该版本号写入服务端的数据记录版本号。
                当NOCHECKDATAVERSION_AUTOINCREASE时: 表示不检测记录版本号，将服务器端的数据记录版本号自增1，若服务器端新写入数据记录则新写入的数据记录的版本号为1。
**/
func (r *Record) SetVersion(v int32) 
```

#### 4.3.4 获取记录版本号
```
/**
	@brief  获取记录版本号
	@retval 记录版本号
**/
func (r *Record) GetVersion() int32
```

### 4.4 TcaplusResponse接口
#### 4.4.1 获取响应结果
```
/*
    @brief  获取响应结果
    @retval int tcaplus api自定义错误码。 0，表示请求成功；非0,有错误码，可从terror.GetErrMsg(int)得到错误消息
*/
GetResult() int
```
#### 4.4.2 获取表名
```
/*
    @brief  获取响应表名
    @retval string 响应消息对应的表名称
*/
GetTableName() string
```
#### 4.4.3 获取appId
```
/*
    @brief  获取响应appId
    @retval uint64 响应消息对应的appId
*/
GetAppId() uint64
```
#### 4.4.4 获取zoneId
```
/*
    @brief  获取响应zoneId
    @retval uint32 响应消息对应的zoneId
*/
GetZoneId() uint32
```
#### 4.4.5 获取响应命令字
```
/*
    @brief  获取响应命令
    @retval int 响应消息命令字，cmd包中的响应命令字
*/
GetCmd() int
```
#### 4.4.6 获取响应异步ID
```
/*
    @brief  获取响应异步id，和请求对应
    @retval uint64 响应消息对应的异步id和请求对应
*/
GetAsyncId() uint64
```
#### 4.4.7 获取响应中记录数
```
/*
    @brief  获取本响应中结果记录条数
    @retval int 响应中结果记录条数
*/
GetRecordCount() int
```
#### 4.4.8 获取响应中一条记录
```
/*
    @brief  从结果中获取一条记录
    @retval *record.Record 记录指针
    @retval error 错误码
*/
FetchRecord() (*record.Record, error)
```
#### 4.4.9 获取响应中用户缓存信息
```
/**
    @brief  获取响应消息中的用户缓存信息
    @retval []byte 用户缓存二进制，和请求消息中的buffer内容一致
*/
GetUserBuffer() []byte
```
#### 4.4.10 获取响应中的序列号
```
/**
    @brief 获取响应消息中的序列号
**/
GetSeq() int32
```

# NOTE
v0.1.3：修复同步接口加锁慢，导致响应消息未收到的问题
v0.1.4: 增加partkey接口，increase，batchget接口
v0.1.5: 增加listgetAll接口
v0.1.6: 优化反射