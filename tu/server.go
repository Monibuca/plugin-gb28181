package tu

import (
	"sync"

	"github.com/Monibuca/plugin-gb28181/v3/transaction"
)

//TODO:参考http服务，使用者仅需要根据需要实现某些handler，替换某些header fileds or body信息。其他的处理都由库来实现。
type Server struct {
	*transaction.Core          //SIP transaction manager
	registers         sync.Map //管理所有已经注册的设备端
	//routers:TODO:消息路由，应用层可以处理消息体，或者针对某些消息的callback
}

//提供config参数
func NewServer(config *transaction.Config) *Server {
	return &Server{
		Core: transaction.NewCore(config),
	}
}

//运行一个sip server
func RunServer() {
	config := &transaction.Config{
		SipIP:      "192.168.1.102",
		SipPort:    5060,
		SipNetwork: "UDP",
		Serial:     "34020000002000000001",
		Realm:      "3402000000",
		AckTimeout: 10,

		RegisterValidity:  3600,
		RegisterInterval:  60,
		HeartbeatInterval: 60,
		HeartbeatRetry:    3,

		AudioEnable:      true,
		WaitKeyFrame:     true,
		MediaPortMin:     58200,
		MediaPortMax:     58300,
		MediaIdleTimeout: 30,
	}
	s := NewServer(config)

	s.StartAndWait()

	select {}
}
