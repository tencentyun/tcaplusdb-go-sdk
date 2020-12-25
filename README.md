# TcaplusDB GO SDK
本SDK主要提供GO语言进行TcaplusDB表数据操作，支持三种协议：
* __PB(protobuf)协议__：支持Google Protobuf进行数据序列化和反序列化操作，具体请参考`pb`子目录相关文档说明
* __TDR(tencent data presentation)协议__：TcaplusDB内部私有化协议进行数据序列化和反序列化，具体请参考`tdr`子目录相关文档说明
* __RESTful协议__: 通用的RESTful协议进行数据操作，具体请参考`restful`子目录相关文档说明
