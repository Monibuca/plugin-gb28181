
#### 介绍

transport 包括客户端和服务器端，仅实现tcp和udp的传输层，不负责具体消息的解析和处理。不负责粘包、半包、消息解析、心跳处理、状态管理等工作。

比如设备关闭或者离线，要修改缓存状态、数据库状态、发送离线通知、执行离线回调等等，都在上层处理。tcp server 和 udp server、tcp client 和 udp client ， 消息的接收和发送都在外面处理。

tcp是流传输，需要注意粘包和半包的处理。在上层处理tcp包的时候，可以尝试使用 ring buffer


#### usage

参考 example.go

#### TODO

- sip协议的传输层，TCP和UDP有所不同，比如重传以及超时的错误信息等。所以需要在 transaction上面，处理消息重传、错误上报等。具体参考RFC3261。

