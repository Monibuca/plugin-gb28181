package transport

import (
	"net"
	"time"
)

/*
transport层仅实现数据的读写，连接的关闭

TCP和UDP的区别

- TCP面向链接，在服务关闭的时候，要先close掉所有客户端连接。所以使用一个map简单做session管理，key 是 remote address。
- udp不需要管理session。
*/

//transport server and client interface
//对于面向连接的服务，需要有两个关闭接口：Close and CloseOne
//非面向连接的服务，不必实现
//TODO：心跳管理，使用timewheel

type ITransport interface {
	Name() string
	ReadPacketChan() <-chan *Packet //读消息，消息处理器需在循环中阻塞读取
	WritePacket(packet *Packet)     //写消息
	Start() error                   //开启连接，阻塞接收消息
	Close() error                   //关闭连接
	IsReliable() bool               //是否可靠传输
}

type IServer interface {
	ITransport
	CloseOne(addr string) //对于关闭某个客户端连接，比如没有鉴权的非法链接，心跳超时等
	IsKeepalive() bool    //persistent connection or not
}

//transport 需要实现的接口如下
type IClient interface {
	ITransport
	LocalAddr() net.Addr      //本地地址
	RemoteAddr() net.Addr     //远程地址
	Heartbeat(packet *Packet) //客户端需要定期发送心跳包到服务器端
}

type Packet struct {
	Type string //消息类型，预留字段，对于客户端主动关闭等消息的上报、心跳超时等。如果为空，则仅透传消息。
	Addr net.Addr
	Data []byte
}

//对于面向连接的（UDP或者TCP都可以面向连接，维持心跳即可），必须有session
type Connection struct {
	Addr           net.Addr
	Conn           net.Conn
	Online         bool
	ReconnectCount int64 //重连次数
}

func (s *Connection) Close() {
	//TODO：处理session的关闭，修改缓存状态、数据库状态、发送离线通知、执行离线回调等等
}

//通讯统计
type Statistic struct {
	startTime time.Time
	stopTime  time.Time
	recvCount int64
	sendCount int64
	errCount  int64
}
