package gb28181

import (
	"strings"

	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/config"
)

type GB28181Config struct {
	AutoInvite     bool
	AutoCloseAfter int
	PreFetchRecord bool

	//sip服务器的配置
	SipNetwork    string //传输协议，默认UDP，可选TCP
	SipIP         string //sip 服务器公网IP
	SipPort       uint16 //sip 服务器端口，默认 5060
	SipExtendPort uint16
	Serial        string //sip 服务器 id, 默认 34020000002000000001
	Realm         string //sip 服务器域，默认 3402000000
	Username      string //sip 服务器账号
	Password      string //sip 服务器密码

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

	AudioEnable       bool //是否开启音频
	LogVerbose        bool
	WaitKeyFrame      bool //是否等待关键帧，如果等待，则在收到第一个关键帧之前，忽略所有媒体流
	RemoveBanInterval int  //移除禁止设备间隔
	UdpCacheSize      int  //udp缓存大小

	config.Publish
	Server
}

func (c *GB28181Config) OnEvent(event any) {
	switch event.(type) {
	case FirstConfig:
		c.startServer()
	}
}

func (c *GB28181Config) IsMediaNetworkTCP() bool {
	return strings.ToLower(c.MediaNetwork) == "tcp"
}

var conf = &GB28181Config{
	AutoInvite:     true,
	AutoCloseAfter: -1,
	PreFetchRecord: false,
	UdpCacheSize:   0,
	SipNetwork:     "udp",
	SipIP:          "127.0.0.1",
	SipPort:        5060,
	SipExtendPort:  45060,
	Serial:         "34020000002000000001",
	Realm:          "3402000000",
	Username:       "",
	Password:       "",

	AckTimeout:        10,
	RegisterValidity:  60,
	RegisterInterval:  60,
	HeartbeatInterval: 60,
	HeartbeatRetry:    3,

	MediaIP:          "127.0.0.1",
	MediaPort:        58200,
	MediaIdleTimeout: 30,
	MediaNetwork:     "udp",

	RemoveBanInterval: 600,
	LogVerbose:        false,
	AudioEnable:       true,
	WaitKeyFrame:      true,
}

var plugin = InstallPlugin(conf)
