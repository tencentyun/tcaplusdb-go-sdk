# Tcaplus Go API TDR Example

注意：
协议限制，SetKey SetValue和SetData接口不可混用;
## 各个目录介绍
* sync是使用go 协程的同步请求，请求方法和C++ API比较类似
  * 注意：
  * 单记录返回的使用Client.Do接口发送请求
  * 多记录返回的，request->SetMultiResponseFlag设置多包返回，并且使用Client.DoMore接口发送请求
* sync2.0是针对sync的进一步封装，Do+option的使用方法，将所有标记统一放在option中，统一设置，使用更加简便
* async异步使用方式，较少场景使用，用法不便，根据业务场景采用，Go的同步协程是最方便的
  * 针对遍历推荐使用异步方法，一边遍历一边处理响应记录，如果是同步遍历，会占较多内存