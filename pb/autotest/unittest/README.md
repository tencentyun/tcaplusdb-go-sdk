# Tcaplus Go API UnitTest

## 1 各个目录介绍
* api_test为测试用例
* cfg是tcaplus的环境信息和日志配置文件
* table下是被测试表，测试其中的table_generic，需要提前将table_common.xml中的table_generic加入到tcaplus中

## 2 编译unittest
1. 提前将table下的xml加到tcaplus，加table_generic表
2. 配置cfg下api_cfg.xml中的tcaplus信息
3. 到api_test执行make
4. 执行生成的二进制文件./apiTest -test.v   会显示每个case的执行情况，并且最后一行打印PASS表示所有用例通过

## 3 other
./apiTest -help可查看单测的其他相关命令
