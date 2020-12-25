# Tcaplus Go PB API 3.46.0
[TOC]
## 1 使用准备
本API是Tcaplus API的Go封装，支持Generic表的增删改查
### 1.1 各代码目录
* pack为打包脚本
* example为示例
* aurotest为测试工具
* 其他目录为Tcaplus Go PB API的源码
* vendor 中为依赖库的代码，需要使用git submodule init和git submodule update拉取

### 1.2 编译example
example 展示了对game_players表的insert get replace update delete操作
1. Go环境安装https://golang.org/
2. 将game_players.proto加入Tcaplus
3. 修改example.go的开头的AppId ZoneId DirUrl Signature为相应的Tcaplus配置信息
4. make之后执行

### 1.3 打包脚本
pack/pack.sh展示了对源码及依赖库的打包，方便用户移植到无法使用go mod的场景
1. cd pack && sh pack.sh

## 2 API的使用
### 2.1 vendor方式使用
Tcaplus API的依赖库及其源码都在打包后的src/vendor目录下，用户只需将vendor放入自己的工程目录即可使用Tcaplus Go API的接口

vendor依赖介绍：
* github.com/tencentyun/tcaplusdb-go-sdk/pb是Tcaplus Go API源码
* git.code.oa.com/ProtoCodeGen是proto工具，可将proto转换为go源码
* github.com/tencentyun/tsf4g/tdrcom是tdr go源码打解包的依赖库
* go.uber.org/zap是日志库
* github.com/natefinch/lumberjack是日志文件切割库
* google.golang.org/protobuf是protobuf库，用于protobuf打解包，源码小幅度修改

### 2.2 mod 方式使用
mod 模式需要在能连内网及公网环境下使用
* 在工程中建立go.mod
* 开启module模式
* 执行命令go clean --modcache 
* 执行命令 go mod edit -require="github.com/tencentyun/tcaplusdb-go-sdk/pb@v0.1.0"
* 出现tlinux无法download的错误可以参考:[Golang git.code.oa.com 的 go get、go mod 踩坑之旅](http://km.oa.com/group/29073/articles/show/376902?kmref=search&from_page=1&no=1#-%20%E9%94%99%E8%AF%AF-x509-%20certificate%20signed%20by%20unknown%20authority)

## 3 接口列表
```
支持命令：

//Generic表插入请求
TcaplusApiInsertReq = 0x0001

//Generic表替换/插入请求
TcaplusApiReplaceReq = 0x0003

//Generic表单条查询请求
TcaplusApiGetReq = 0x0007

//Generic表删除请求
TcaplusApiDeleteReq = 0x0009

//Generic表更新请求
TcaplusApiUpdateReq = 0x001d

//批量查询请求
TcaplusApiBatchGetReq = 0x0017

//部分Key查询请求
TcaplusApiGetByPartkeyReq = 0x0019

//表遍历请求
TcaplusApiTableTraverseReq = 0x0045

//protobuf部分字段获取请求
TcaplusApiPBFieldGetReq = 0x0067

//protobuf部分字段更新请求
TcaplusApiPBFieldUpdateReq = 0x0069

//protobuf部分字段自增请求
TcaplusApiPBFieldIncreaseReq = 0x006b

//索引查询请求
TcaplusApiSqlReq = 0x0081

```
### 3.1 Client接口
#### 3.1.1 创建tcaplus pbclient
```
/**
   @brief 创建一个tcaplus api客户端
   @retval 返回客户端指针
**/
func NewPBClient() *PBClient 
```
#### 3.1.2 日志配置接口
创建client之后，立刻调用（不调用此接口日志将会打到控制台）
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

#### 3.1.3 连接tcaplus接口
```
/**
   @brief 连接tcaplue函数
   @param [IN] appId         appId，在网站注册相应服务以后，你可以得到该appId
   @param [IN] zoneList      需要操作表的区服ID列表，操作的表在多个不同的zone，填zoneId列表；操作的表在一个zone，zone列表填一个zoneId
   @param [IN] signature     签名/密码，在网站注册相应服务以后，你可以得到该字符串
   @param [IN] dirUrl        目录服务器的url，形如"tcp://172.25.40.181:10600"
   @param [IN] timeout       second, 连接所有表对应的tcaplus proxy服务器。若所有的proxy连通且鉴权通过，则立即返回成功；
                                若到达超时时间，只要有一个proxy连通且鉴权通过，也会返回成功；否则返回超时错误。
   @param [IN] zoneTable     将会用到的pb表(zone:tables)
   @retval                   错误码
**/
func (c *Client) Dial(appId uint64, zoneList []uint32, dirUrl string, signature string, timeout uint32, zoneTable map[uint32][]string{}) error
```

#### 3.1.4 创建tcaplus请求
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

#### 3.1.5 发送tcaplus请求
```
/**
    @brief 发送tcaplus请求
    @param [IN] req       tcaplus请求
    @retval error         错误码
*/
func (c *Client) SendRequest(req request.TcaplusRequest) error
```

#### 3.1.6 异步接收tcaplus响应
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

#### 3.1.7 发送tcaplus同步请求并接受响应
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

#### 3.1.8 发送tcaplus同步请求并接受多个响应
与3.1.7的区别为：3.1.7请求只会有一个响应，3.1.8请求会有多个响应，例如：
TcaplusApiBatchGetReq TcaplusApiGetByPartkeyReq TcaplusApiSqlReq
```
/**
    @brief 发送tcaplus同步请求并接受响应
	@param [IN] req tcaplus请求
	@param [IN] timeout 超时时间
    @retval []response.TcaplusResponse tcaplus响应
    @retval error 错误码
            error nil，response nil 成功但当前无响应消息
            error nil, response 非nil，成功获取响应消息
            error 非nil，response 非nil 接收部分回包正确，但是收到了错误包或者超时退出
**/
func (c *client) DoMore(req request.TcaplusRequest, timeout time.Duration) ([]response.TcaplusResponse, error)
```

#### 3.1.9 设置默认zoneId
连接数据库后会将传入的zoneTables的第一个zone作为默认zoneId(根据需要调用，非必须)
```
/**
    @brief 设置默认zoneId
    @param [IN] zoneId zoneID
    @retval error 错误码，如果未dial调用此接口将会返错 ClientNotDial
**/
func (c *PBClient) SetDefaultZoneId(zoneId uint32) error
```

#### 3.1.10 设置默认超时时间
默认超时时间5s(根据需要调用，非必须)
```
/**
    @brief 设置默认超时时间
    @param [IN] t time.Duration
    @retval error 错误码，如果未dial调用此接口将会返错 ClientNotDial
**/
func (c *PBClient) SetDefaultTimeOut(t time.Duration) error
```

#### 3.1.11 插入记录
```
/**
    @brief 插入记录，可以使用 SetDefaultZoneId 来设置zoneid； SetDefaultTimeOut 设置超时时间
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (c *PBClient) Insert(msg proto.Message) error

/**
    @brief 指定zone插入记录，当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) InsertWithZone(msg proto.Message, zoneId uint32) error
```
#### 3.1.12 替换记录
调用成功msg将带回此次替换前的记录
如果记录不存在，将此条记录插入。替换的是整条记录，只需要替换部分字段不要使用此接口。
```
/**
    @brief 替换记录，记录不存在时插入
    @param [IN] msg proto.Message 由proto文件生成的记录结构体，调用成功msg将带回替换前的值
    @retval error 错误码
**/
func (c *PBClient) Replace(msg proto.Message) error

/**
    @brief 替换记录，记录不存在时插入。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msg proto.Message 由proto文件生成的记录结构体，调用成功msg将带回替换前的值
    @param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) ReplaceWithZone(msg proto.Message, zoneId uint32) error
```
#### 3.1.13 修改记录
调用成功msg将带回此次修改前的记录
如果记录不存在，将返回错误。修改的是整条记录，只需要修改部分字段不要使用此接口。
```
/**
    @brief 修改记录，记录不存在时返错
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (c *PBClient) Update(msg proto.Message) error

/**
    @brief 修改记录，记录不存在时返错。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) UpdateWithZone(msg proto.Message, zoneId uint32) error
```
#### 3.1.11 删除记录
调用成功msg将带回此次删除的记录
```
/**
    @brief 删除记录
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (c *PBClient) Delete(msg proto.Message) error

/**
    @brief 删除记录。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) DeleteWithZone(msg proto.Message, zoneId uint32) error
```
#### 3.1.12 获取记录
调用成功msg会带回此次获取到的记录
```
/**
    @brief 获取记录
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (c *PBClient) Get(msg proto.Message) error

/**
    @brief 获取记录。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) GetWithZone(msg proto.Message, zoneId uint32) error
```
#### 3.1.13 批量获取记录
调用成功msgs会带回此次获取到的所有记录
```
/**
    @brief 批量获取记录
    @param [IN] msgs []proto.Message 需获取的记录列表
    @retval error 错误码
**/
func (c *PBClient) BatchGet(msgs []proto.Message) error

/**
    @brief 批量获取记录。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msgs []proto.Message 需获取的记录列表
    @param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) BatchGetWithZone(msgs []proto.Message, zoneId uint32) error
```
#### 3.1.14 部分key获取记录
```
/**
    @brief 批量获取记录
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] keys []string 部分key
    @retval []proto.Message 返回记录，可能匹配到多条记录
    @retval error 错误码
**/
func (c *PBClient) GetByPartKey(msg proto.Message, keys []string) ([]proto.Message, error)

/**
    @brief 部分key获取记录。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
	@param [IN] msg proto.Message 由proto文件生成的记录结构体
	@param [IN] keys []string 部分key，根据 proto 文件中的 index 选择填写
	@param [IN] zoneId 指定表所在zone
	@retval []proto.Message 返回记录，可能匹配到多条记录
    @retval error 错误码
**/
func (c *PBClient) GetByPartKeyWithZone(msg proto.Message, keys []string, zoneId uint32) ([]proto.Message, error)
```
#### 3.1.15 获取部分value
调用成功msg会带回此次获取到的记录
```
/**
    @brief 获取记录部分字段value
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] values []string 部分字段名，根据需要选择填写
    @retval error 错误码
**/
func (c *PBClient) FieldGet(msg proto.Message, values []string) error

/**
    @brief 获取记录部分字段value。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] values []string 部分字段名，根据需要选择填写
    @param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) FieldGetWithZone(msg proto.Message, values []string, zoneId uint32) error
```
#### 3.1.16 更新部分value
调用成功msg会带回此次更新后的记录
```
/**
    @brief 更新记录部分字段value
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] values []string 部分字段名，根据需要选择填写
    @retval error 错误码
**/
func (c *PBClient) FieldUpdate(msg proto.Message, values []string) error

/**
    @brief 更新记录部分字段value。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] values []string 部分字段名，根据需要选择填写
    @param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) FieldUpdateWithZone(msg proto.Message, values []string, zoneId uint32) error
```
#### 3.1.17 自增部分value（仅支持整型）
调用成功msg会带回此次自增后的记录
```
/**
    @brief 自增记录部分字段value
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] values []string 部分字段名，根据需要选择填写
    @retval error 错误码
**/
func (c *PBClient) FieldIncrease(msg proto.Message, values []string) error

/**
    @brief 自增记录部分字段value。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] values []string 部分字段名，根据需要选择填写
    @param [IN] zoneId 指定表所在zone
    @retval error 错误码
**/
func (c *PBClient) FieldIncreaseWithZone(msg proto.Message, values []string, zoneId uint32) error
```
#### 3.1.18 二级索引查询
```
/**
    @brief 分布式索引查询
    @param [IN] query sql 查询语句 详情见 https://iwiki.woa.com/pages/viewpage.action?pageId=419645505
    @retval []proto.Message 非聚合查询结果
    @retval []string 聚合查询结果
    @retval error 错误码
**/
func (c *PBClient) IndexQuery(query string) ([]proto.Message, []string, error)

/**
    @brief 自增记录部分字段value
    @param [IN] query sql 查询语句 详情见 https://iwiki.woa.com/pages/viewpage.action?pageId=419645505
    @param [IN] zoneId 指定表所在zone
    @retval []proto.Message 非聚合查询结果
    @retval []string 聚合查询结果
    @retval error 错误码
**/
func (c *PBClient) IndexQueryWithZone(query string, zoneId uint32) ([]proto.Message, []string, error)
```
#### 3.1.19 获取遍历器
```
/**
    @brief 获取遍历器（存在则直接获取，不存在则新建一个）
    @param [IN] zoneId tcaplus请求
    @param [IN] table 超时时间
    @retval *traverser.Traverser 遍历器，一个client最多分配8个遍历器，超过将会返回 nil
**/
func (c *client) GetTraverser(zoneId uint32, table string) *traverser.Traverser
```

### 3.2 TcaplusRequest接口
#### 3.2.1 添加记录
```
/**
  @brief  向请求中添加一条记录。
  @param [IN] index         用于List操作(目前不支持)，通常>=0，表示该Record在所属List中的Index；对于Generic操作，index无意义，设0即可
  @retval record.Record     返回记录指针
  @retval error   			错误码
**/
AddRecord(index int32) (*record.Record, error)
```
#### 3.2.2 设置请求异步ID
```
/**
    @brief  设置请求的异步事务ID，api会将其值不变地通过对应的响应消息带回来
    @param  [IN] asyncId  请求对应的异步事务ID
**/
SetAsyncId(id uint64)
```
#### 3.2.3 设置版本校验规则
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
#### 3.2.4 设置响应标志
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
#### 3.2.5 设置用户缓存
```
/**
    @brief 设置用户缓存，响应消息将携带返回
    @param [IN] userBuffer  用户缓存
    @retval error           错误码
**/
SetUserBuff(userBuffer []byte) error
```
#### 3.2.6 返回记录条数限制
```
/**
    @brief  如果此请求会返回多条记录，通过此接口对返回的记录做一些限制
    @param [IN] limit       需要查询的记录条数, limit若等于-1表示操作或返回所有匹配的数据记录.
    @param [IN] offset      记录起始编号；若设置为负值(-N, N>0)，则从倒数第N个记录开始返回结果
    @retval 0               设置成功
    @retval <0              设置失败，具体错误参见 \link ErrorCode \endlink
    @note 对于Generic类型的部分Key查询，limit表示所要获取Record的条数，offset表示所要获取Record的开始下标；
          对于List类型的GetAll操作，limit表示所要获取Record的条数，offset表示所要获取Record的开始下标，
          在当前版本中这些Record一定属于同一个List.
          该函数仅仅对于GET_BY_PARTKEY(Generic类型的部分Key查询), UPDATE_BY_PARTKEY,
          DELETE_BY_PARTKEY, LIST_GETALL(List类型的GetAll操)这4种操作类型有效。
*/
SetResultLimit(limit int32, offset int32) int32
```
#### 3.2.7 设置分包
```
/**
    @brief  设置是否允许一个请求包可以自动响应多个应答包，仅对ListGetAll和BatchGet协议有效。
    @param [IN] multi_flag   多响应包标示，1表示允许一个请求包可以自动响应多个应答包, 
                           0表示不允许一个请求包自动响应多个应答包
    @retval 0                设置成功
    @retval <0               设置失败，具体错误参见 \link ErrorCode \endlink
    @note	分包应答，目前只支持ListGetAll和BatchGet操作；其他操作设置该值是没有意义的，
            函数会返回<0的错误码。
*/
SetMultiResponseFlag(multi_flag byte) int32
```
#### 3.2.8 设置sql语句
```
/*
    @brief  添加LIST记录的元素索引值。该函数只对于 TcaplusApiSqlReq 有效
    @param  query sql语句
    @retval 0                 设置成功
    @retval 非0               设置失败，具体错误参见 \link ErrorCode \endlink
*/
SetSql(query string) int
```

### 3.3 Record接口
#### 3.3.1 SetPBData和GetPBData接口
通过PB Message，对记录进行赋值(请求消息)和获取(响应消息)
```
/**
    @brief  基于 PB Message 设置record数据
    @param [IN] data  PB Message
    @retval []byte 记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应
    @retval error     错误码
*/
func (r *Record) SetPBData(message proto.Message) ([]byte, error)

/**
    @brief  基于 PB Message 读取record数据
    @param [IN] data   PB Message
    @retval []byte 记录的keybuf，用来唯一确定一条记录，多用于请求与响应记录相对应
    @retval error      错误码
**/
func (r *Record) GetPBData(data proto.Message) ([]byte, error)
```
#### 3.3.2 设置记录版本号
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

#### 3.3.3 获取记录版本号
```
/**
    @brief  获取记录版本号
    @retval 记录版本号
**/
func (r *Record) GetVersion() int32
```
#### 3.3.4 获取部分记录值
```
/**
    @brief 获取部分记录值, 专用于 field 方法，TcaplusApiPBFieldGetReq TcaplusApiPBFieldUpdateReq TcaplusApiPBFieldIncreaseReq
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (r *Record) GetPBFieldValues(message proto.Message) error
```
#### 3.3.5 获取记录key编码值
```
/**
    @brief 获取记录key编码值 
    @retval []byte 由记录key字段编码生成，由于多条记录的响应记录是无序的，可以用这个值来匹配记录
    @retval error 错误码
**/
func (r *Record) GetPBKey() ([]byte, error)
```
#### 3.3.6 设置部分key字段
```
/**
    @brief 设置部分key字段，专用于partkey操作，TcaplusApiGetByPartkeyReq
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] keys []string 指定本次设置的 key 字段
    @retval []byte 由记录key字段编码生成，由于多条记录的响应记录是无序的，可以用这个值来匹配记录
    @retval error 错误码
**/
func (r *Record) SetPBPartKeys(message proto.Message, keys []string) ([]byte, error)
```
#### 3.3.7 设置部分value字段
```
/**
    @brief 设置部分value字段，专用于field操作，TcaplusApiPBFieldGetReq TcaplusApiPBFieldUpdateReq TcaplusApiPBFieldIncreaseReq
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] values []string 指定本次设置的 value 字段
    @retval []byte 由记录key字段编码生成，由于多条记录的响应记录是无序的，可以用这个值来匹配记录
    @retval error 错误码
**/
func (r *Record) SetPBFieldValues(message proto.Message, values []string) ([]byte, error)
```

### 3.4 TcaplusResponse接口
#### 3.4.1 获取响应结果
```
/*
    @brief  获取响应结果
    @retval int tcaplus api自定义错误码。 0，表示请求成功；非0,有错误码，可从terror.GetErrMsg(int)得到错误消息
*/
GetResult() int
```
#### 3.4.2 获取表名
```
/*
    @brief  获取响应表名
    @retval string 响应消息对应的表名称
*/
GetTableName() string
```
#### 3.4.3 获取appId
```
/*
    @brief  获取响应appId
    @retval uint64 响应消息对应的appId
*/
GetAppId() uint64
```
#### 3.4.4 获取zoneId
```
/*
    @brief  获取响应zoneId
    @retval uint32 响应消息对应的zoneId
*/
GetZoneId() uint32
```
#### 3.4.5 获取响应命令字
```
/*
    @brief  获取响应命令
    @retval int 响应消息命令字，cmd包中的响应命令字
*/
GetCmd() int
```
#### 3.4.6 获取响应异步ID
```
/*
    @brief  获取响应异步id，和请求对应
    @retval uint64 响应消息对应的异步id和请求对应
*/
GetAsyncId() uint64
```
#### 3.4.7 获取响应中记录数
```
/*
    @brief  获取本响应中结果记录条数
    @retval int 响应中结果记录条数
*/
GetRecordCount() int
```
#### 3.4.8 获取响应中一条记录
```
/*
    @brief  从结果中获取一条记录
    @retval *record.Record 记录指针
    @retval error 错误码
*/
FetchRecord() (*record.Record, error)
```
#### 3.4.9 获取响应中用户缓存信息
```
/**
    @brief  获取响应消息中的用户缓存信息
    @retval []byte 用户缓存二进制，和请求消息中的buffer内容一致
*/
GetUserBuffer() []byte
```
#### 3.4.10 获取响应中的序列号
```
/**
    @brief 获取响应消息中的序列号
**/
GetSeq() int32
```
#### 3.4.11 获取分布式索引结果
```
/*
    @brief 该函数仅用于索引查询类型为聚合查询时获取聚合结果
    @retval(2) 返回聚合查询结果，返回NULL表示获取失败
    @note 如果索引查询类型为记录查询，该函数将返回NULL
            如果索引查询类型为聚合查询，则根据SqlResult类提供的函数来获取查询结果
*/
FetchSqlResult() (*sqlResult, error)

/*
    仅聚合操作可以使用，请在使用前判断索引类型
    处理聚合分布式索引响应并返回处理后的结果，切片长度为结果行数，以逗号作为字段分割符。
*/
ProcAggregationSqlQueryType() ([]string, error)

/*
    获取分布式索引类型 2 为聚合操作  1 为非聚合操作  0 为无效操作
*/
GetSqlType() int
```
#### 3.4.12 判断是否有更多的回包
```
/*
    @判断是否有更多的回包
    @retval  1 有， 0 没有
*/
HaveMoreResPkgs() int
```
#### 3.4.13 获取整个结果中的记录条数
```
/**
    @brief  获取整个结果中的记录条数。既包括本响应返回的记录数，也包括本响应未返回的记录数。
    @retval 记录条数
    @note 该函数只能用于以下请求：
	(1）TCAPLUS_API_GET_BY_PARTKEY_REQ（部分key查询）；
	(2）TCAPLUS_API_LIST_GETALL_REQ（list表查询所有匹配的记录）；
	(3）TCAPLUS_API_BATCH_GET_REQ（批量查询）

    @note 该函数的作用是：
       (1）当使用了SetResultLimit()来限制返回的记录时，
               使用该函数可以获取所有匹配的记录的个数，
               包括本响应返回的记录数和本响应未返回的记录数，
               而GetRecordCount()函数只能获取本响应返回的记录的个数；
       (2）当所返回的记录很多时，需要分包的时候，
               使用该函数可以获取总共的记录数，
               即多个分包所有记录数的总和，
               而GetRecordCount()函数只能返回单个分包中的(本响应中的)记录数.
*/
GetRecordMatchCount() int
```
### 3.5 遍历
从3.1.19获取遍历器
#### 3.5.1 限制条件（非必须）
```
/**
    @brief 设定本次遍历多少条记录，默认遍历所有
    @param [IN] limit 限制条数
    @retval error 错误码
**/
func (t *Traverser) SetLimit(limit int64) error

/**
    @brief 设置异步id
    @param [IN] id 异步id
    @retval error 错误码
**/
func (t *Traverser) SetAsyncId(id uint64) error


/**
    @brief 设置仅从slave获取记录，默认false
    @param [IN] flag bool
    @retval error 错误码
**/
func (t *Traverser) SetOnlyReadFromSlave(flag bool) error

/**
    @brief 设置用户缓存，响应消息将携带返回
    @param [IN] userBuffer  用户缓存
    @retval error           错误码
**/
func (t *Traverser) SetUserBuff(buf []byte) error
```
#### 3.5.2 开始遍历
```
// 开始遍历，仅当状态为TraverseStateReady可调用
func (t *Traverser) Start() error
// 结束遍历
func (t *Traverser) Stop() error
// 恢复遍历，仅当状态为TraverseStateRecoverable可调用
func (t *Traverser) Resume() error

// 获取状态 t.State()
TraverseStateIdle          = 1      // 结束状态（遍历完毕）
TraverseStateReady         = 2      // 准备状态（初始化成功，可以start）
TraverseStateNormal        = 4      // 遍历中
TraverseStateStop          = 8      // 停止状态（处于此状态会被回收）
TraverseStateRecoverable   = 16     // 可恢复状态（某个响应出问题，可以恢复继续遍历）
TraverseStateUnRecoverable = 32     // 不可恢复状态（获取shardlist出错，或者发生了主备切换）

// 注：遍历器上限最多 8 个，请在用完后调用 t.Stop() 来回收，否则可能导致 6.3失败。
```
