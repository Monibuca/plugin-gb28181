package gb28181

import (
	"fmt"
	"strings"
	"time"

	myip "github.com/husanpao/ip"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/config"
)

type GB28181PositionConfig struct {
	AutosubPosition bool          //是否自动订阅定位
	Expires         time.Duration //订阅周期(单位：秒)
	Interval        time.Duration //订阅间隔（单位：秒）
}

type GB28181Config struct {
	AutoInvite     bool
	InviteIDs      string
	PreFetchRecord bool

	//sip服务器的配置
	SipNetwork string //传输协议，默认UDP，可选TCP
	SipIP      string //sip 服务器公网IP
	SipPort    uint16 //sip 服务器端口，默认 5060
	Serial     string //sip 服务器 id, 默认 34020000002000000001
	Realm      string //sip 服务器域，默认 3402000000
	Username   string //sip 服务器账号
	Password   string //sip 服务器密码

	// AckTimeout        uint16 //sip 服务应答超时，单位秒
	RegisterValidity time.Duration //注册有效期，单位秒，默认 3600
	// RegisterInterval  int    //注册间隔，单位秒，默认 60
	HeartbeatInterval time.Duration //心跳间隔，单位秒，默认 60
	// HeartbeatRetry    int    //心跳超时次数，默认 3

	//媒体服务器配置
	MediaIP      string //媒体服务器地址
	MediaPort    uint16 //媒体服务器端口
	MediaNetwork string //媒体传输协议，默认UDP，可选TCP
	MediaPortMin uint16
	MediaPortMax uint16
	// MediaIdleTimeout uint16 //推流超时时间，超过则断开链接，让设备重连

	// WaitKeyFrame      bool //是否等待关键帧，如果等待，则在收到第一个关键帧之前，忽略所有媒体流
	RemoveBanInterval time.Duration //移除禁止设备间隔
	UdpCacheSize      int           //udp缓存大小

	config.Publish
	Server
	LogLevel string //trace, debug, info, warn, error, fatal, panic
	routes   map[string]string
	DumpPath string //dump PS流本地文件路径

	Position GB28181PositionConfig //关于定位的配置参数
}

func (c *GB28181Config) initRoutes() {
	c.routes = make(map[string]string)
	tempIps := myip.LocalAndInternalIPs()
	for k, v := range tempIps {
		c.routes[k] = v
		if lastdot := strings.LastIndex(k, "."); lastdot >= 0 {
			c.routes[k[0:lastdot]] = k
		}
	}
	plugin.Info(fmt.Sprintf("LocalAndInternalIPs detail: %s", c.routes))
}
func (c *GB28181Config) OnEvent(event any) {
	switch event.(type) {
	case FirstConfig:
		c.ReadDevices()
		go c.initRoutes()
		c.startServer()
	}
}

func (c *GB28181Config) IsMediaNetworkTCP() bool {
	return strings.ToLower(c.MediaNetwork) == "tcp"
}

func (c *GB28181Config) CanInvite(deviceID string) bool {
	if len(deviceID) != 20 {
		return false
	}

	if c.InviteIDs == "" {
		return true
	}

	// 11～13位是设备类型编码
	typeID := deviceID[10:13]

	// format: start-end,type1,type2
	tokens := strings.Split(c.InviteIDs, ",")
	for _, tok := range tokens {
		if first, second, ok := strings.Cut(tok, "-"); ok {
			if typeID >= first && typeID <= second {
				return true
			}
		} else {
			if typeID == first {
				return true
			}
		}
	}

	return false
}

var conf = &GB28181Config{
	AutoInvite:     true,
	PreFetchRecord: false,
	UdpCacheSize:   0,
	SipNetwork:     "udp",
	SipIP:          "",
	SipPort:        5060,
	Serial:         "34020000002000000001",
	Realm:          "3402000000",
	Username:       "",
	Password:       "",

	// AckTimeout:        10,
	RegisterValidity: 60 * time.Second,
	// RegisterInterval:  60,
	HeartbeatInterval: 60 * time.Second,
	// HeartbeatRetry:    3,

	MediaIP:   "",
	MediaPort: 58200,
	// MediaIdleTimeout: 30,
	MediaNetwork: "udp",

	RemoveBanInterval: 600 * time.Second,
	LogLevel:          "info",
	// WaitKeyFrame:      true,
	Position: GB28181PositionConfig{AutosubPosition: false, Expires: 3600 * time.Second, Interval: 6 * time.Second},
}

var plugin = InstallPlugin(conf)
