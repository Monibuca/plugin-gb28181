package tu

import (
	"fmt"
	"github.com/Monibuca/plugin-gb28181/transaction"
)

//sip server和client的配置，可以得到sip URI：sip
//格式：user:password@host:port;uri-parameters?headers
//在这些URI里边包含了足够的信息来发起和维持到这个资源的一个通讯会话。
//client静态配置
type ClientStatic struct {
	LocalIP   string //设备本地IP
	LocalPort uint16 //客户端SIP端口
	Username  string //SIP用户名，一般是取通道ID，默认 34020000001320000001
	AuthID    string //SIP用户认证ID，一般是通道ID， 默认 34020000001320000001
	Password  string //密码
}

//client运行时信息
type ClientRuntime struct {
	RemoteAddress string //设备的公网的IP和端口，格式x.x.x.x:x
	Online        bool   //在线状态
	Branch        string //branch
	Cseq          int    //消息序列号，发送消息递增, uint32
	FromTag       string //from tag
	ToTag         string //to tag
	Received      string //remote ip
	Rport         string //remote port
}

type Client struct {
	*transaction.Core                //transaction manager
	static            *ClientStatic  //静态配置
	runtime           *ClientRuntime //运行时信息
}

//config:sip信令服务器配置
//static:sip客户端配置
func NewClient(config *transaction.Config, static *ClientStatic) *Client {
	return &Client{
		Core:    transaction.NewCore(config),
		static:  static,
		runtime: &ClientRuntime{},
	}
}

//TODO：对于一个TU，开启之后
//运行一个sip client
func RunClient() {
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
	static := &ClientStatic{
		LocalIP:   "192.168.1.65",
		LocalPort: 5060,
		Username:  "34020000001320000001",
		AuthID:    "34020000001320000001",
		Password:  "123456",
	}
	c := NewClient(config, static)

	go c.Start()

	//TODO：先发起注册
	//TODO:build sip message
	msg := BuildMessageRequest("", "", "", "", "", "",
		0, 0, 0, "")
	resp := c.SendMessage(msg)
	if resp.Code != 0 {
		fmt.Println("request failed")
	}
	fmt.Println("response: ", resp.Data)

	select {}
}
