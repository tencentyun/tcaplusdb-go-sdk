# Tcaplus Go API AutoTest

## 1 各个目录介绍
* case测试所用的案例，详情见xml注释
* cfg/apiLogConf.xml为api日志配置文件
* src下为agent源码

## 2 编译
1. 提前将测试表加入tcaplus
2. 配置case中的tcaplus信息和表信息
3. src目录下make
4. tdrAgent -h可看help
5. ./tdrAgent -case ../case/case.xml -apiLog ../cfg/apiLogConf.xml

## 3 说明
* 目前只支持generic表的get insert replace update delete命令
* 日志打印在src/log下
* 控制台每秒打印 当前速度，当前发包，当前平均时延，当前最大时延等统计信息
