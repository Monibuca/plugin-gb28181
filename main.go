package gb28181

import (
	"fmt"
	myip "github.com/husanpao/ip"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/config"
	"m7s.live/plugin/gb28181/v4/utils"
	"regexp"
	"strings"
)

// InviteOnSubscribe 按需Invite过滤器。
type InviteOnSubscribe struct {
	On bool // 是否开启
	DaemonAppList []string // 设备号列表
	Filter     string //正则表达式
	filterReg  *regexp.Regexp
}

type GB28181Config struct {
	AutoInvite              bool
	PreFetchRecord          bool
	InviteOnSubscribe      // 按需拉流过滤器
	//sip服务器的配置
	SipNetwork string //传输协议，默认UDP，可选TCP
	SipIP      string //sip 服务器公网IP
	SipPort    uint16 //sip 服务器端口，默认 5060
	Serial     string //sip 服务器 id, 默认 34020000002000000001
	Realm      string //sip 服务器域，默认 3402000000
	Username   string //sip 服务器账号
	Password   string //sip 服务器密码

	AckTimeout        uint16 //sip 服务应答超时，单位秒
	RegisterValidity  int    //注册有效期，单位秒，默认 3600
	RegisterInterval  int    //注册间隔，单位秒，默认 60
	HeartbeatInterval int    //心跳间隔，单位秒，默认 60
	HeartbeatRetry    int    //心跳超时次数，默认 3

	//媒体服务器配置
	MediaIP          string //媒体服务器地址
	MediaPort        uint16 //媒体服务器端口
	MediaNetwork     string //媒体传输协议，默认UDP，可选TCP
	MediaPortMin     uint16
	MediaPortMax     uint16
	MediaIdleTimeout uint16 //推流超时时间，超过则断开链接，让设备重连

	// WaitKeyFrame      bool //是否等待关键帧，如果等待，则在收到第一个关键帧之前，忽略所有媒体流
	RemoveBanInterval int //移除禁止设备间隔
	UdpCacheSize      int //udp缓存大小

	config.Publish
	Server
	LogLevel string //trace, debug, info, warn, error, fatal, panic
	routes   map[string]string
	DumpPath string //dump PS流本地文件路径
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
	switch v := event.(type) {
	case FirstConfig:
		ReadDevices()
		go c.initRoutes()
		c.startServer()
		// 按需拉流开启
		if c.On {
			if c.Filter != "" {
				c.filterReg = regexp.MustCompile(c.Filter)
			}
		}
	case UnsubscribeEvent:
		// 最后一个订阅者离开。
		stream := v.Subscriber.GetSubscriber().Stream
		if stream != nil && len(stream.Subscribers) == 0 {
			if conf.On {
				// 按需拉流的设备通道在最后一个订阅者离开后，主动断开推流。
				if (len(conf.DaemonAppList) > 0 && utils.In(stream.AppName, conf.DaemonAppList)) ||
					(conf.filterReg != nil && conf.filterReg.MatchString(stream.Path)) {
					stream.Receive(ACTION_LASTLEAVE)
				}
			}
		}
	case *Stream:
		if v.AppName != "" && v.StreamName != ""{
			channel := FindChannel(v.AppName, v.StreamName)
			if channel != nil && channel.LivePublisher == nil {
				channel.Invite(InviteOptions{InviteOnSubscribe: true})
			}
		}
	}
}

func (c *GB28181Config) IsMediaNetworkTCP() bool {
	return strings.ToLower(c.MediaNetwork) == "tcp"
}

var conf = &GB28181Config{
	AutoInvite:        true,
	PreFetchRecord:    false,
	UdpCacheSize:      0,
	SipNetwork:        "udp",
	SipIP:             "",
	SipPort:           5060,
	Serial:            "34020000002000000001",
	Realm:             "3402000000",
	Username:          "",
	Password:          "",
	AckTimeout:        10,
	RegisterValidity:  60,
	RegisterInterval:  60,
	HeartbeatInterval: 60,
	HeartbeatRetry:    3,

	MediaIP:          "",
	MediaPort:        58200,
	MediaIdleTimeout: 30,
	MediaNetwork:     "udp",

	RemoveBanInterval: 600,
	LogLevel:          "info",
	// WaitKeyFrame:      true,
}

var plugin = InstallPlugin(conf)
