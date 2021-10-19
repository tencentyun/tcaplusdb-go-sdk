# Tcaplus Go API Example

## 1 各个目录介绍
* service_info.xml测试所用的表
* service_info目录是使用tdr工具将service_info.xml转换成的go源码(cd vendor/git.code.oa.com(v0.2.4,v0.1.9后用git.woa.com)/tsf4g/TdrCodeGen/;python tdr.py service_info.xml得到service_info的go代码)
* logconf.xml为api日志配置文件
* main.go为示例代码

## 2 编译example
1. 提前将service_info.xml加入tcaplus
2. 配置main.go开头的tcaplus信息
```
const (
	AppId = uint64(2)
	ZoneId = uint32(3)
	DirUrl = "tcp://x.x.x.x:9999"
	Signature = "xxxxxx"
	TableName = "service_info"
)
```
3. make进行编译
4. ./async 执行example

## 3 注意
设置record有两套接口，切记不可混用：
1. SetKey SetValue接口设置的数据，只能通过GetKey，GetValue接口读取
2. SetData接口设置的数据，只能通过GetData读取
