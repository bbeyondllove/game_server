## kk_server说明 

这是一个kk_core的一个demo，单机版本

#### 协议说明
消息的包头和包体结构设计如下：
包头：msglen(4字节) + cmd（2个字节） + length（4字节）
包体：protobuf字节数组

为同时支持TCP和WebSocket，包头加了msglen， 因为是demo，没有设计加密解密机制，包体采用proto数据交换格式，经编码后转码为字节数组


