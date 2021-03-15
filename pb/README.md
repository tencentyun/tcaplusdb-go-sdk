# Tcaplus Go PB SDK 3.46.0
Table of Contents
=================

   * [Tcaplus Go PB SDK 3.46.0](#tcaplus-go-pb-sdk-3460)
      * [1 SDK 说明](#1-sdk-\xE8\xAF\xB4\xE6\x98\x8E)
      * [2 SDK 使用方式](#2-sdk-\xE4\xBD\xBF\xE7\x94\xA8\xE6\x96\xB9\xE5\xBC\x8F)
         * [2.1mod 方式使用](#21mod-\xE6\x96\xB9\xE5\xBC\x8F\xE4\xBD\xBF\xE7\x94\xA8)
      * [3 接口列表](#3-\xE6\x8E\xA5\xE5\x8F\xA3\xE5\x88\x97\xE8\xA1\xA8)
         * [3.1 Client 接口](#31-client-\xE6\x8E\xA5\xE5\x8F\xA3)
            * [3.1.1 创建 tcaplus pbclient](#311-\xE5\x88\x9B\xE5\xBB\xBA-tcaplus-pbclient)
            * [3.1.2 日志配置接口](#312-\xE6\x97\xA5\xE5\xBF\x97\xE9\x85\x8D\xE7\xBD\xAE\xE6\x8E\xA5\xE5\x8F\xA3)
            * [3.1.3 连接 tcaplus 接口](#313-\xE8\xBF\x9E\xE6\x8E\xA5-tcaplus-\xE6\x8E\xA5\xE5\x8F\xA3)
            * [3.1.4 创建 tcaplus 请求](#314-\xE5\x88\x9B\xE5\xBB\xBA-tcaplus-\xE8\xAF\xB7\xE6\xB1\x82)
            * [3.1.5 发送 tcaplus 请求](#315-\xE5\x8F\x91\xE9\x80\x81-tcaplus-\xE8\xAF\xB7\xE6\xB1\x82)
            * [3.1.6 异步接收 tcaplus 响应](#316-\xE5\xBC\x82\xE6\xAD\xA5\xE6\x8E\xA5\xE6\x94\xB6-tcaplus-\xE5\x93\x8D\xE5\xBA\x94)
            * [3.1.7 发送 tcaplus 同步请求并接受响应](#317-\xE5\x8F\x91\xE9\x80\x81-tcaplus-\xE5\x90\x8C\xE6\xAD\xA5\xE8\xAF\xB7\xE6\xB1\x82\xE5\xB9\xB6\xE6\x8E\xA5\xE5\x8F\x97\xE5\x93\x8D\xE5\xBA\x94)
            * [3.1.8 发送 tcaplus 同步请求并接受多个响应](#318-\xE5\x8F\x91\xE9\x80\x81-tcaplus-\xE5\x90\x8C\xE6\xAD\xA5\xE8\xAF\xB7\xE6\xB1\x82\xE5\xB9\xB6\xE6\x8E\xA5\xE5\x8F\x97\xE5\xA4\x9A\xE4\xB8\xAA\xE5\x93\x8D\xE5\xBA\x94)
            * [3.1.9 设置默认 zoneId (非必须)](#319-\xE8\xAE\xBE\xE7\xBD\xAE\xE9\xBB\x98\xE8\xAE\xA4-zoneid-\xE9\x9D\x9E\xE5\xBF\x85\xE9\xA1\xBB)
            * [3.1.10 设置默认超时时间](#3110-\xE8\xAE\xBE\xE7\xBD\xAE\xE9\xBB\x98\xE8\xAE\xA4\xE8\xB6\x85\xE6\x97\xB6\xE6\x97\xB6\xE9\x97\xB4)
            * [3.1.11 插入记录](#3111-\xE6\x8F\x92\xE5\x85\xA5\xE8\xAE\xB0\xE5\xBD\x95)
            * [3.1.12 替换记录](#3112-\xE6\x9B\xBF\xE6\x8D\xA2\xE8\xAE\xB0\xE5\xBD\x95)
            * [3.1.13 修改记录](#3113-\xE4\xBF\xAE\xE6\x94\xB9\xE8\xAE\xB0\xE5\xBD\x95)
            * [3.1.11 删除记录](#3111-\xE5\x88\xA0\xE9\x99\xA4\xE8\xAE\xB0\xE5\xBD\x95)
            * [3.1.12 获取记录](#3112-\xE8\x8E\xB7\xE5\x8F\x96\xE8\xAE\xB0\xE5\xBD\x95)
            * [3.1.13 批量获取记录](#3113-\xE6\x89\xB9\xE9\x87\x8F\xE8\x8E\xB7\xE5\x8F\x96\xE8\xAE\xB0\xE5\xBD\x95)
            * [3.1.14 部分 key 获取记录](#3114-\xE9\x83\xA8\xE5\x88\x86-key-\xE8\x8E\xB7\xE5\x8F\x96\xE8\xAE\xB0\xE5\xBD\x95)
            * [3.1.15 获取部分 value](#3115-\xE8\x8E\xB7\xE5\x8F\x96\xE9\x83\xA8\xE5\x88\x86-value)
            * [3.1.16 更新部分 value](#3116-\xE6\x9B\xB4\xE6\x96\xB0\xE9\x83\xA8\xE5\x88\x86-value)
            * [3.1.17 自增部分 value（仅支持整型）](#3117-\xE8\x87\xAA\xE5\xA2\x9E\xE9\x83\xA8\xE5\x88\x86-value\xE4\xBB\x85\xE6\x94\xAF\xE6\x8C\x81\xE6\x95\xB4\xE5\x9E\x8B)
            * [3.1.18 二级索引查询](#3118-\xE4\xBA\x8C\xE7\xBA\xA7\xE7\xB4\xA2\xE5\xBC\x95\xE6\x9F\xA5\xE8\xAF\xA2)
            * [3.1.19 获取表记录数](#3119-\xE8\x8E\xB7\xE5\x8F\x96\xE8\xA1\xA8\xE8\xAE\xB0\xE5\xBD\x95\xE6\x95\xB0)
            * [3.1.20 获取遍历器](#3120-\xE8\x8E\xB7\xE5\x8F\x96\xE9\x81\x8D\xE5\x8E\x86\xE5\x99\xA8)
         * [3.2 TcaplusRequest 接口](#32-tcaplusrequest-\xE6\x8E\xA5\xE5\x8F\xA3)
            * [3.2.1 添加记录](#321-\xE6\xB7\xBB\xE5\x8A\xA0\xE8\xAE\xB0\xE5\xBD\x95)
            * [3.2.2 设置请求异步 ID](#322-\xE8\xAE\xBE\xE7\xBD\xAE\xE8\xAF\xB7\xE6\xB1\x82\xE5\xBC\x82\xE6\xAD\xA5-id)
            * [3.2.3 设置版本校验规则](#323-\xE8\xAE\xBE\xE7\xBD\xAE\xE7\x89\x88\xE6\x9C\xAC\xE6\xA0\xA1\xE9\xAA\x8C\xE8\xA7\x84\xE5\x88\x99)
            * [3.2.4 设置响应标志](#324-\xE8\xAE\xBE\xE7\xBD\xAE\xE5\x93\x8D\xE5\xBA\x94\xE6\xA0\x87\xE5\xBF\x97)
            * [3.2.5 设置用户缓存](#325-\xE8\xAE\xBE\xE7\xBD\xAE\xE7\x94\xA8\xE6\x88\xB7\xE7\xBC\x93\xE5\xAD\x98)
            * [3.2.6 返回记录条数限制](#326-\xE8\xBF\x94\xE5\x9B\x9E\xE8\xAE\xB0\xE5\xBD\x95\xE6\x9D\xA1\xE6\x95\xB0\xE9\x99\x90\xE5\x88\xB6)
            * [3.2.7 设置分包](#327-\xE8\xAE\xBE\xE7\xBD\xAE\xE5\x88\x86\xE5\x8C\x85)
            * [3.2.8 设置 sql 语句](#328-\xE8\xAE\xBE\xE7\xBD\xAE-sql-\xE8\xAF\xAD\xE5\x8F\xA5)
         * [3.3 Record 接口](#33-record-\xE6\x8E\xA5\xE5\x8F\xA3)
            * [3.3.1 SetPBData 和 GetPBData 接口](#331-setpbdata-\xE5\x92\x8C-getpbdata-\xE6\x8E\xA5\xE5\x8F\xA3)
            * [3.3.2 设置记录版本号](#332-\xE8\xAE\xBE\xE7\xBD\xAE\xE8\xAE\xB0\xE5\xBD\x95\xE7\x89\x88\xE6\x9C\xAC\xE5\x8F\xB7)
            * [3.3.3 获取记录版本号](#333-\xE8\x8E\xB7\xE5\x8F\x96\xE8\xAE\xB0\xE5\xBD\x95\xE7\x89\x88\xE6\x9C\xAC\xE5\x8F\xB7)
            * [3.3.4 SetPBFieldValues 和 GetPBFieldValues 获取部分记录值](#334-setpbfieldvalues-\xE5\x92\x8C-getpbfieldvalues-\xE8\x8E\xB7\xE5\x8F\x96\xE9\x83\xA8\xE5\x88\x86\xE8\xAE\xB0\xE5\xBD\x95\xE5\x80\xBC)
            * [3.3.5 设置部分 key 字段](#335-\xE8\xAE\xBE\xE7\xBD\xAE\xE9\x83\xA8\xE5\x88\x86-key-\xE5\xAD\x97\xE6\xAE\xB5)
            * [3.3.6 获取记录 key 编码值](#336-\xE8\x8E\xB7\xE5\x8F\x96\xE8\xAE\xB0\xE5\xBD\x95-key-\xE7\xBC\x96\xE7\xA0\x81\xE5\x80\xBC)
         * [3.4 TcaplusResponse 接口](#34-tcaplusresponse-\xE6\x8E\xA5\xE5\x8F\xA3)
            * [3.4.1 获取响应结果](#341-\xE8\x8E\xB7\xE5\x8F\x96\xE5\x93\x8D\xE5\xBA\x94\xE7\xBB\x93\xE6\x9E\x9C)
            * [3.4.2 获取表名](#342-\xE8\x8E\xB7\xE5\x8F\x96\xE8\xA1\xA8\xE5\x90\x8D)
            * [3.4.3 获取 appId](#343-\xE8\x8E\xB7\xE5\x8F\x96-appid)
            * [3.4.4 获取 zoneId](#344-\xE8\x8E\xB7\xE5\x8F\x96-zoneid)
            * [3.4.5 获取响应命令字](#345-\xE8\x8E\xB7\xE5\x8F\x96\xE5\x93\x8D\xE5\xBA\x94\xE5\x91\xBD\xE4\xBB\xA4\xE5\xAD\x97)
            * [3.4.6 获取响应异步 ID](#346-\xE8\x8E\xB7\xE5\x8F\x96\xE5\x93\x8D\xE5\xBA\x94\xE5\xBC\x82\xE6\xAD\xA5-id)
            * [3.4.7 获取响应中记录数](#347-\xE8\x8E\xB7\xE5\x8F\x96\xE5\x93\x8D\xE5\xBA\x94\xE4\xB8\xAD\xE8\xAE\xB0\xE5\xBD\x95\xE6\x95\xB0)
            * [3.4.8 获取响应中一条记录](#348-\xE8\x8E\xB7\xE5\x8F\x96\xE5\x93\x8D\xE5\xBA\x94\xE4\xB8\xAD\xE4\xB8\x80\xE6\x9D\xA1\xE8\xAE\xB0\xE5\xBD\x95)
            * [3.4.9 获取响应中用户缓存信息](#349-\xE8\x8E\xB7\xE5\x8F\x96\xE5\x93\x8D\xE5\xBA\x94\xE4\xB8\xAD\xE7\x94\xA8\xE6\x88\xB7\xE7\xBC\x93\xE5\xAD\x98\xE4\xBF\xA1\xE6\x81\xAF)
            * [3.4.10 获取响应中的序列号](#3410-\xE8\x8E\xB7\xE5\x8F\x96\xE5\x93\x8D\xE5\xBA\x94\xE4\xB8\xAD\xE7\x9A\x84\xE5\xBA\x8F\xE5\x88\x97\xE5\x8F\xB7)
            * [3.4.11 获取分布式索引结果](#3411-\xE8\x8E\xB7\xE5\x8F\x96\xE5\x88\x86\xE5\xB8\x83\xE5\xBC\x8F\xE7\xB4\xA2\xE5\xBC\x95\xE7\xBB\x93\xE6\x9E\x9C)
            * [3.4.12 判断是否有更多的回包](#3412-\xE5\x88\xA4\xE6\x96\xAD\xE6\x98\xAF\xE5\x90\xA6\xE6\x9C\x89\xE6\x9B\xB4\xE5\xA4\x9A\xE7\x9A\x84\xE5\x9B\x9E\xE5\x8C\x85)
            * [3.4.13 获取整个结果中的记录条数](#3413-\xE8\x8E\xB7\xE5\x8F\x96\xE6\x95\xB4\xE4\xB8\xAA\xE7\xBB\x93\xE6\x9E\x9C\xE4\xB8\xAD\xE7\x9A\x84\xE8\xAE\xB0\xE5\xBD\x95\xE6\x9D\xA1\xE6\x95\xB0)
         * [3.5 遍历](#35-\xE9\x81\x8D\xE5\x8E\x86)
            * [3.5.1 限制条件（非必须）](#351-\xE9\x99\x90\xE5\x88\xB6\xE6\x9D\xA1\xE4\xBB\xB6\xE9\x9D\x9E\xE5\xBF\x85\xE9\xA1\xBB)
            * [3.5.2 开始遍历](#352-\xE5\xBC\x80\xE5\xA7\x8B\xE9\x81\x8D\xE5\x8E\x86)
      * [4. 错误码](#4-\xE9\x94\x99\xE8\xAF\xAF\xE7\xA0\x81)
      * [5.附录](#5\xE9\x99\x84\xE5\xBD\x95)
         * [5.1 条件查询](#51-\xE6\x9D\xA1\xE4\xBB\xB6\xE6\x9F\xA5\xE8\xAF\xA2)
         * [5.2 分页查询](#52-\xE5\x88\x86\xE9\xA1\xB5\xE6\x9F\xA5\xE8\xAF\xA2)
         * [5.3 聚合查询](#53-\xE8\x81\x9A\xE5\x90\x88\xE6\x9F\xA5\xE8\xAF\xA2)
         * [5.4 支持查询部分字段的值](#54-\xE6\x94\xAF\xE6\x8C\x81\xE6\x9F\xA5\xE8\xAF\xA2\xE9\x83\xA8\xE5\x88\x86\xE5\xAD\x97\xE6\xAE\xB5\xE7\x9A\x84\xE5\x80\xBC)
         * [5.5 不支持的 sql 查询语句](#55-\xE4\xB8\x8D\xE6\x94\xAF\xE6\x8C\x81\xE7\x9A\x84-sql-\xE6\x9F\xA5\xE8\xAF\xA2\xE8\xAF\xAD\xE5\x8F\xA5)
            * [5.5.1 不支持聚合查询与非聚合查询混用](#551-\xE4\xB8\x8D\xE6\x94\xAF\xE6\x8C\x81\xE8\x81\x9A\xE5\x90\x88\xE6\x9F\xA5\xE8\xAF\xA2\xE4\xB8\x8E\xE9\x9D\x9E\xE8\x81\x9A\xE5\x90\x88\xE6\x9F\xA5\xE8\xAF\xA2\xE6\xB7\xB7\xE7\x94\xA8)
            * [5.5.2 不支持 order by 查询](#552-\xE4\xB8\x8D\xE6\x94\xAF\xE6\x8C\x81-order-by-\xE6\x9F\xA5\xE8\xAF\xA2)
            * [5.5.3 不支持 group by 查询](#553-\xE4\xB8\x8D\xE6\x94\xAF\xE6\x8C\x81-group-by-\xE6\x9F\xA5\xE8\xAF\xA2)
            * [5.5.4 不支持 having 查询](#554-\xE4\xB8\x8D\xE6\x94\xAF\xE6\x8C\x81-having-\xE6\x9F\xA5\xE8\xAF\xA2)
            * [5.5.5 不支持多表联合查询](#555-\xE4\xB8\x8D\xE6\x94\xAF\xE6\x8C\x81\xE5\xA4\x9A\xE8\xA1\xA8\xE8\x81\x94\xE5\x90\x88\xE6\x9F\xA5\xE8\xAF\xA2)
            * [5.5.6 不支持嵌套 select 查询](#556-\xE4\xB8\x8D\xE6\x94\xAF\xE6\x8C\x81\xE5\xB5\x8C\xE5\xA5\x97-select-\xE6\x9F\xA5\xE8\xAF\xA2)
            * [5.5.7 不支持别名](#557-\xE4\xB8\x8D\xE6\x94\xAF\xE6\x8C\x81\xE5\x88\xAB\xE5\x90\x8D)
            * [5.5.8 不支持的其他查询](#558-\xE4\xB8\x8D\xE6\x94\xAF\xE6\x8C\x81\xE7\x9A\x84\xE5\x85\xB6\xE4\xBB\x96\xE6\x9F\xA5\xE8\xAF\xA2)
      * [6. 其它](#6-\xE5\x85\xB6\xE5\xAE\x83)

## 1 SDK 说明

本 SDK 支持通过 GO 来操作 TcaplusDB Protobuf 表的数据，共支持 12 个接口。包括：插入、替换、查询、删除、更新、批量查询、主键索引查询、遍历表、部分字段获取、部分字段更新、部分字段自增和二级索引查询。接口支持同步调用模式和异步调用模式。

- **同步模式**: 接口调用逻辑较简单，适合对性能要求不高场景
- **异步模式**: 接口调用逻辑稍微复杂，适合高吞吐、高并发业务场景

## 2 SDK 使用方式

目前 SDK 支持通过 go mod 方式来管理整个 package。在使用时可参考[SDK Example](https://github.com/tencentyun/tcaplusdb-go-examples.git)，有详细 SDK 接口示例说明。

### 2.1mod 方式使用

mod 模式需要在能连公网环境下使用。对于用户新建项目，可参考如下步骤引入 SDK 到项目中：

- 1.在工程中建立 go.mod
- 2.开启 module 模式
- 3.执行命令 go clean --modcache
- 4.执行命令 go mod edit -require="github.com/tencentyun/tcaplusdb-go-sdk@v0.0.7"
- 5.在代码中引入 sdk: import "github.com/tencentyun/tcaplusdb-go-sdk/pb"

## 3 接口列表

目前 SDK 接口以不同请求命令字方式来区分，具体如下：

```
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

//遍历全表请求
TcaplusApiTableTraverseReq = 0x0045

//table的记录总数请求
TcaplusApiGetTableRecordCountReq = 0x0053
```

### 3.1 Client 接口

业务在调用 SDK 接口时，需要先初始化连接客户端，步骤如下。

#### 3.1.1 创建 tcaplus pbclient

```
/**
   @brief 创建一个tcaplus api客户端
   @retval 返回客户端指针
**/
func NewPBClient() *PBClient
```

#### 3.1.2 日志配置接口

创建 client 之后，需要配置日志（备注：**若不调用此接口日志将会直接输出控制台**）。

```
/**
   @brief                   设置API日志配置文件全路径log.conf(json格式，example下有示例)，请在client.Dial之前调用
   @param [IN] cfgPath      日志配置文件全路径logconf.xml, 可在配置文件中配置日志级别如：ERROR | INFO | WARN | DEBUG
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

#### 3.1.3 连接 tcaplus 接口

在初始化客户端连接指针后，需要调用具体的连接接口建立与 TcalusDB 后端连接。

```
/**
   @brief 连接tcaplue函数
   @param [IN] appId         集群接入ID, 云环境：参考控制台集群详情页，本地Docker版：对于Protobuf业务默认创建了一个集群，默认为3，用户无需再创建
   @param [IN] zoneList      集群表格组id列表(区服ID), 支持指定1个或多个表格组id
   @param [IN] signature     集群访问密码，云环境：参考控制台集群访问密码，本地Docker版: 登录TcalusDB web运维平台(用户名/密码:tcaplus/tcaplus), 业务管理->业务维护->查看pb_app业务密码
   @param [IN] dirUrl        集群访问地址，形如"tcp://172.25.40.181:9999"
   @param [IN] timeout       访问超时设置，单位：秒, 连接所有表对应的tcaplus proxy服务器。若所有的proxy连通且鉴权通过，则立即返回成功；
                                若到达超时时间，只要有一个proxy连通且鉴权通过，也会返回成功；否则返回超时错误。
   @param [IN] zoneTable     以map对象映射表格组id和对应的表列表, key: 表格组id, value: 表名列表
   @retval                   错误码
**/
func (c *Client) Dial(appId uint64, zoneList []uint32, dirUrl string, signature string, timeout uint32, zoneTable map[uint32][]string{}) error
```

#### 3.1.4 创建 tcaplus 请求

```
/**
    @brief 创建指定分区指定表的tcaplus请求
    @param [IN] zoneId              区服ID,表格组id
    @param [IN] tableName           表名
    @param [IN] cmd                 命令字(cmd包中cmd.TcaplusApiGetReq等)
    @retval request.TcaplusRequest  tcaplus请求
    @retval error                   错误码
*/
func (c *Client) NewRequest(zoneId uint32, tableName string, cmd int) (request.TcaplusRequest, error)
```

#### 3.1.5 发送 tcaplus 请求

```
/**
    @brief 发送tcaplus请求
    @param [IN] req       tcaplus请求
    @retval error         错误码
*/
func (c *Client) SendRequest(req request.TcaplusRequest) error
```

#### 3.1.6 异步接收 tcaplus 响应

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

#### 3.1.7 发送 tcaplus 同步请求并接受响应

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

#### 3.1.8 发送 tcaplus 同步请求并接受多个响应

与 3.1.7 的区别为：3.1.7 请求只会有一个响应，3.1.8 请求会有多个响应，例如：
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

#### 3.1.9 设置默认 zoneId (非必须)

连接数据库后会将传入的 zoneTables 的第一个 zone 作为默认 zoneId(根据需要调用，非必须)

```
/**
    @brief 设置默认zoneId
    @param [IN] zoneId zoneID
    @retval error 错误码，如果未dial调用此接口将会返错 ClientNotDial
**/
func (c *PBClient) SetDefaultZoneId(zoneId uint32) error
```

#### 3.1.10 设置默认超时时间

默认超时时间 5s(根据需要调用，非必须)

```
/**
    @brief 设置默认超时时间
    @param [IN] t time.Duration
    @retval error 错误码，如果未dial调用此接口将会返错 ClientNotDial
**/
func (c *PBClient) SetDefaultTimeOut(t time.Duration) error
```

#### 3.1.11 插入记录

插入单条记录

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

调用成功 msg 将带回此次替换前的记录。
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

调用成功 msg 将带回此次修改前的记录。记录不存在会报错。
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

调用成功 msg 将带回此次删除的记录。根据主键删除单条记录。

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

调用成功 msg 会带回此次获取到的记录。一次返回单条记录。

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

调用成功 msgs 会带回此次获取到的所有记录。批量获取数据接口方便用户一次返回多条记录，比如一次返回多个玩家的记录用于在业务层作聚合操作。

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

#### 3.1.14 部分 key 获取记录

此接口主要作用于表定义的主键索引，TcaplusDB 支持最多`8`个联合主键字段，主键索引可支持建`4`个，每个主键索引可支持 1 个或多个主键字段构成，这样方便用户灵活根据业务场景进行组合，满足更多查询场景需要。
注意：

- **如果表没定义主键索引，此接口无效。**
- ** 表主键索引通过在 proto 文件中定义: option(tcaplusservice.tcaplus_index) = "index_1(pk_field_1, pk_field_2)"; 来实现，具体参考示例中的表定义文件**

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

#### 3.1.15 获取部分 value

调用成功 msg 会带回此次获取到的记录。通过此接口可实现只返回少数字段，避免返回整条记录，对于记录字段数的表效率尤其明显，可大大降低返回包的大小，及提高解析包的效率。

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

#### 3.1.16 更新部分 value

调用成功 msg 会带回此次更新后的记录。通过此接口研发同学可避免更新少数字段需要传整条记录的情况，大幅增加传输效率。

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

#### 3.1.17 自增部分 value（仅支持整型）

调用成功 msg 会带回此次自增后的记录。

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

二级索引查询支持通过 SQL 语法进行数据查询，主要基于 TcaplusDB 的全局二级索引字段进行 Select 查询，在 where 条件中可用指定为索引的字段进行范围查询，模糊查询，等值查询和聚合查询。
注意前提：**在云控制台已经针对表添加了全局二级索引**,　如若未添加是无法使用此接口的。

```
/**
    @brief 全局二级索引查询
    @param [IN] query sql 查询语句 详情见 附录
    @retval []proto.Message 非聚合查询结果
    @retval []string 聚合查询结果
    @retval error 错误码
**/
func (c *PBClient) IndexQuery(query string) ([]proto.Message, []string, error)

/**
    @brief 自增记录部分字段value
    @param [IN] query sql 查询语句,语法参考附录
    @param [IN] zoneId 指定表所在zone
    @retval []proto.Message 非聚合查询结果
    @retval []string 聚合查询结果
    @retval error 错误码
**/
func (c *PBClient) IndexQueryWithZone(query string, zoneId uint32) ([]proto.Message, []string, error)
```


#### 3.1.19 获取表记录数
```
/**
    @brief 获取表记录总数
    @param [IN] table string 表名
    @retval int 记录数，请求失败返回0
    @retval error 错误码
**/
func (c *PBClient) GetTableCount(table string) (int, error)

/**
    @brief 获取表记录总数。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] table string 表名
    @param [IN] zoneId 指定表所在zone
    @retval int 记录数，请求失败返回0
    @retval error 错误码
**/
func (c *PBClient) GetTableCountWithZone(table string, zoneId uint32) (int, error)
```
#### 3.1.20 获取遍历器

用于遍历全表接口使用。

```
/**
    @brief 获取遍历器（存在则直接获取，不存在则新建一个）
    @param [IN] zoneId tcaplus请求
    @param [IN] table 超时时间
    @retval *traverser.Traverser 遍历器，一个client最多分配8个遍历器，超过将会返回 nil
**/
func (c *client) GetTraverser(zoneId uint32, table string) *traverser.Traverser
```
#### 3.1.21 获取appId
```
/**
    @brief 获取本次连接的appId
    @retval int appId
**/
func (c *client) GetAppId() uint64
```
#### 3.1.22 关闭client
```
/**
    @brief 关闭client，释放资源
**/
func (c *client) Close()
```
#### 3.1.23 遍历表记录
```
/**
    @brief 遍历表
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval []proto.Message 查询结果列表
    @retval error 错误码
**/
func (c *PBClient) Traverse(msg proto.Message) ([]proto.Message, error)

/**
    @brief 遍历表。当并发时如果zoneId各不相同，无法通过 SetDefaultZoneId 来设置zoneid，需使用此接口
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @param [IN] zoneId 指定表所在zone
    @retval []proto.Message 查询结果列表
    @retval error 错误码
**/
func (c *PBClient) TraverseWithZone(msg proto.Message, zoneId uint32) ([]proto.Message, error)
```

### 3.2 TcaplusRequest 接口

#### 3.2.1 添加记录

一次请求支持添加多条需要操作的记录，通过 AddRecord 实现此逻辑，同时为兼容 TcaplusDB 的 List 类型表，支持添加记录到 List 记录的指定下标位置，相当于在指定数组下标下添加一条要操作的记录。本文档只介绍 Generic 表，所以对于 AddRecord 的下标索引默认为`0`即可。

```
/**
  @brief  向请求中添加一条记录。
  @param [IN] index         用于List表操作(目前不支持)，通常>=0，表示该Record在所属List中的Index；对于Generic表操作，index无意义，设0即可
  @retval record.Record     返回记录指针
  @retval error   			错误码
**/
AddRecord(index int32) (*record.Record, error)
```

#### 3.2.2 设置请求异步 ID

此接口主要是为映射发送请求体与响应请求体之间的关系，通过此 ID 来表示响应请求属于哪个发送请求。

```
/**
    @brief  设置请求的异步事务ID，api会将其值不变地通过对应的响应消息带回来
    @param  [IN] asyncId  请求对应的异步事务ID
**/
SetAsyncId(id uint64)
```

#### 3.2.3 设置版本校验规则

通过版本校验接口，可以灵活设置数据的版本号，也可设置严格的写入数据校验机制，避免数据写乱、写错，极大的保障了数据的一致性、安全性。

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

设置响应包返回的格式，如只返回响应成功与否、返回原始记录或只返回新的记录。主要用于比对发送的数据是否和接收的数据保持一致，可减少研发自身去判断此类逻辑的工作量。

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

可以简单理解是一种上下文 Context 机制。用户缓存主要用于一些全局变量场景，对于异步调用模式，处理响应数据是异步的，有一些数据在发送请求时用到，同时也希望在响应请求处理时用到，对于异步请求这个场景研发自己实现的话需要设置大量的全局变量来做，不好管理。有了 UserBuffer，就不用设置大量的全局变量来保存一些发送与接收请求都需要用到的数据，直接通过请求本身来传递此类数据。也大大节省了研发工作量。另一种场景是用于保存请求 id, 类似上面 SetAsyncId 接口，以实现上下文 Context 这种效果。

```
/**
    @brief 设置用户缓存，响应消息将携带返回
    @param [IN] userBuffer  用户缓存
    @retval error           错误码
**/
SetUserBuff(userBuffer []byte) error
```

#### 3.2.6 返回记录条数限制

此设置可以避免一次返回过多数据导致返回记录分包不正常。

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

#### 3.2.8 设置 sql 语句

用于 IndexQuery 接口，二级索引查询通过设置 SQL 来实现查询逻辑。SQL 语法参考`附录`。

```
/*
    @brief  添加LIST记录的元素索引值。该函数只对于 TcaplusApiSqlReq 有效
    @param  query sql语句
    @retval 0                 设置成功
    @retval 非0               设置失败，具体错误参见 \link ErrorCode \endlink
*/
SetSql(query string) int
```

### 3.3 Record 接口

#### 3.3.1 SetPBData 和 GetPBData 接口

通过 PB Message，对记录进行赋值(请求消息)和获取(响应消息)

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

#### 3.3.4 SetPBFieldValues 和 GetPBFieldValues 获取部分记录值

主要用于`FieldGet, FieldUpdate, FieldIncrease`三个接口。用于设置需要操作的记录部分字段情况。

```
/**
    @brief 设置部分value字段，专用于field操作，TcaplusApiPBFieldGetReq TcaplusApiPBFieldUpdateReq TcaplusApiPBFieldIncreaseReq
    @param [IN] msg proto.Message 由proto文件生成的记录结构体，这里通常为表主键字段值和需要操作的部分字段值
    @param [IN] values []string 指定本次设置的 value 字段
    @retval []byte 由记录key字段编码生成，由于多条记录的响应记录是无序的，可以用这个值来匹配记录
    @retval error 错误码
**/
func (r *Record) SetPBFieldValues(message proto.Message, values []string) ([]byte, error)
```

```
/**
    @brief 获取部分记录值, 专用于 field 方法，TcaplusApiPBFieldGetReq TcaplusApiPBFieldUpdateReq TcaplusApiPBFieldIncreaseReq
    @param [IN] msg proto.Message 由proto文件生成的记录结构体
    @retval error 错误码
**/
func (r *Record) GetPBFieldValues(message proto.Message) error
```

#### 3.3.5 设置部分 key 字段

用于根据表定义中的主键索引字段来查询数据。

```
/**
    @brief 设置部分key字段，专用于partkey操作，TcaplusApiGetByPartkeyReq
    @param [IN] msg proto.Message 由proto文件生成的记录结构体，即为PartKey字段对应的值
    @param [IN] keys []string 指定本次设置的 key 字段，即为PartKey字段名（主键索引对应的字段名）
    @retval []byte 由记录key字段编码生成，由于多条记录的响应记录是无序的，可以用这个值来匹配记录
    @retval error 错误码
**/
func (r *Record) SetPBPartKeys(message proto.Message, keys []string) ([]byte, error)
```

#### 3.3.6 获取记录 key 编码值

```
/**
    @brief 获取记录key编码值
    @retval []byte 由记录key字段编码生成，由于多条记录的响应记录是无序的，可以用这个值来匹配记录
    @retval error 错误码
**/
func (r *Record) GetPBKey() ([]byte, error)
```

### 3.4 TcaplusResponse 接口

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

#### 3.4.3 获取 appId

```
/*
    @brief  获取响应appId
    @retval uint64 响应消息对应的appId
*/
GetAppId() uint64
```

#### 3.4.4 获取 zoneId

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

#### 3.4.6 获取响应异步 ID

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
#### 3.4.14 获取表记录总数
```
/*
    @获取表的记录总数，只适用于TCAPLUS_API_GET_TABLE_RECORD_COUNT_REQ请求获取返回结果
    @retval  <0 出错  记录总数
*/
GetTableRecordCount() int
```

### 3.5 遍历

从 3.1.19 获取遍历器

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

## 4. 错误码

SDK 所的有错误码描述均在源码目录`terror/error.go`中，用户可自行参考错误码描述，错误码命名基本能反映一些错误的一些原因，如果有疑惑可随时 TcaplusDB 相关同学。下面是一些常见的错误码列表：

| 编号 | 错误码 | 描述                                                   |
| ---- | ------ | ------------------------------------------------------ |
| 1    | -1792  | 表处于只读模式，请检查 RCU,                            |
| 3    | 261    | 该记录不存在                                           |
| 4    | -525   | batchget 操作请求超时,                                 |
| 5    | -781   | batchget,                                              |
| 6    | -1037  | 系统繁忙，请联系管理员                                 |
| 7    | -1293  | 记录已存在，请不要重复插入                             |
| 8    | -1549  | 访问的表字段不存在                                     |
| 9    | -2061  | 表字段类型错误                                         |
| 10   | -3085  | SetFieldName 操作指定了错误的字段                      |
| 11   | -3341  | 字段值大小超过其定义类型的限制                         |
| 12   | -4109  | list 数据类型元素下标超过范围                          |
| 14   | -4621  | 请求缺少主键字段或索引字段                             |
| 15   | -6157  | list 表元素个数超过定义范围,请设置元素淘汰             |
| 16   | -6925  | result_flag 设置错误，请参考 SDK 中 result_flag 说明   |
| 17   | -7949  | 请检查乐观锁，请求记录版本号与实际记录版本号不一致     |
| 18   | -11277 | 操作表的方法不存在                                     |
| 19   | -16141 | PB 表 GetRecord 操作失败，请联系管理员                 |
| 20   | -16397 | PB 表非主键字段值超过限定大小(256KB)                   |
| 21   | -16653 | PB 表 FieldSetRecord 操作失败，请联系管理员            |
| 22   | -16909 | PB 表 FieldIncRecord 操作失败，请联系管理员            |
| 23   | -275   | 主键字段个数超过限制，Generic 表限制数为 4,            |
| 24   | -531   | 非主键字段个数超过限制，Generic 表限制数为 128,        |
| 25   | -787   | 字段名称大小超过限制（32B）                            |
| 26   | -1043  | 字段值指超过限制（256KB）                              |
| 27   | -1555  | 字段值的数据类型与其定义类型不匹配                     |
| 28   | -5395  | 请求中缺少主键                                         |
| 29   | -9235  | index 不存在                                           |
| 30   | -12307 | 请求发送失败，网络过载，请联系管理员。                 |
| 31   | -12819 | 表不存在                                               |
| 32   | -13843 | 后台网络异常，请求无法发送成功，如持续存在请联系管理员 |
| 33   | -14099 | 插入的记录超过大小限制（1MB）                          |
| 34   | －6673 | 请求参数无主键                                         |
| 35   | －6929 | 请求参数缺少主键字段                                   |

## 5.附录

主要介绍二级索引查询所支持的 SQL 语法,　注意前提：**在云控制台已经针对表添加了全局二级索引**,　如若未添加是无法使用此接口的。

### 5.1 条件查询

支持 =, >, >=, <, <=, !=, between, in, not in, like, not like, and, or , 比如:

```
select * from table where a > 100 and b < 1000;

select * from table where a between 1 and 100 and b < 1000;

select * from table where str like "test";

select * from table where a > 100 or b < 1000;
```

注意：between 查询时，between a and b，对应的查询范围为[a, b]，比如 between 1 and 100, 是会包含 1 和 100 这两个值的，即查询范围为[1,100]

注意：like 查询是支持模糊匹配，其中"%"通配符，匹配 0 个或者多个字符； “\_”通配符，匹配 1 个字符；
分页查询

### 5.2 分页查询

支持 limit offset 分页查询。
比如：

```
select * from table whre a > 100 limit 100 offset 0;
```

注意：当前 limit 必须与 offset 搭配使用，即不支持 limit 1 或者 limit 0,1 这种。

### 5.3 聚合查询

当前支持的聚合查询包括：sum, count, max, min, avg，比如：

```
select sum(a), count(*), max(a), min(a), avg(a) from table where a > 1000;
```

注意：聚合查询不支持 limit offset，即 limit offset 不生效；

注意：目前只有 count 支持 distinct，即 select count(distinct(a)) from table where a > 1000; 其他情况均不支持 distinct
部分字段查询

### 5.4 支持查询部分字段的值

```
select a, b from table where a > 1000;
```

对于 pb 表，还支持查询嵌套字段的值，用点分方式，类似：

```
select field1.field2.field3, a, b from table where a > 1000;
```

### 5.5 不支持的 sql 查询语句

#### 5.5.1 不支持聚合查询与非聚合查询混用

```
select *, a, b from table where a > 1000;

select sum(a), a, b from table where a  > 1000;

select count(*), * from table where a  > 1000;
```

#### 5.5.2 不支持 order by 查询

```
select * from table where a > 1000 order by a;
```

#### 5.5.3 不支持 group by 查询

```
select * from table where a > 1000 group by a;
```

#### 5.5.4 不支持 having 查询

```
select sum(a) from table where  a > 1000 group by a having sum(a) > 10000;
```

#### 5.5.5 不支持多表联合查询

```
select * from table1 where table1.a > 1000 and table1.a = table2.b;
```

#### 5.5.6 不支持嵌套 select 查询

```
select * from table where a > 1000 and b in (select b from table where b < 5000);
```

#### 5.5.7 不支持别名

```
select sum(a) as sum_a from table where a > 1000;
```

#### 5.5.8 不支持的其他查询

- 不支持 join 查询；
- 不支持 union 查询；
- 不支持类似 select a+b from table where a > 1000 的查询；
- 不支持类似 select \* from table where a+b > 1000 的查询；
- 不支持类似 select \* from table where a >= b 的查询；
- 不支持其他未提到的查询。

## 6. 其它

