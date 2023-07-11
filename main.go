package gb28181

import (
	"os"
	"strings"
	"sync"
	"time"

	myip "github.com/husanpao/ip"
	"go.uber.org/zap"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/util"
)

type GB28181PositionConfig struct {
	AutosubPosition bool          //是否自动订阅定位
	Expires         time.Duration `default:"3600s"` //订阅周期(单位：秒)
	Interval        time.Duration `default:"6s"`    //订阅间隔（单位：秒）
}

type GB28181Config struct {
	InviteMode int    `default:"1"` //邀请模式，0:手动拉流，1:预拉流，2:按需拉流
	InviteIDs  string //按照国标gb28181协议允许邀请的设备类型:132 摄像机 NVR
	ListenAddr string `default:"0.0.0.0"`
	//sip服务器的配置
	SipNetwork string   `default:"udp"` //传输协议，默认UDP，可选TCP
	SipIP      string   //sip 服务器公网IP
	SipPort    uint16   `default:"5060"`                 //sip 服务器端口，默认 5060
	Serial     string   `default:"34020000002000000001"` //sip 服务器 id, 默认 34020000002000000001
	Realm      string   `default:"3402000000"`           //sip 服务器域，默认 3402000000
	Username   string   //sip 服务器账号
	Password   string   //sip 服务器密码
	Port       struct { // 新配置方式
		Sip   string `default:"udp:5060"`
		Media string `default:"tcp:58200-59200"`
	}
	// AckTimeout        uint16 //sip 服务应答超时，单位秒
	RegisterValidity time.Duration `default:"3600s"` //注册有效期，单位秒，默认 3600
	// RegisterInterval  int    //注册间隔，单位秒，默认 60
	HeartbeatInterval time.Duration `default:"60s"` //心跳间隔，单位秒，默认 60
	// HeartbeatRetry    int    //心跳超时次数，默认 3

	//媒体服务器配置
	MediaIP      string //媒体服务器地址
	MediaPort    uint16 `default:"58200"` //媒体服务器端口
	MediaNetwork string `default:"tcp"`   //媒体传输协议，默认UDP，可选TCP
	MediaPortMin uint16 `default:"58200"`
	MediaPortMax uint16 `default:"59200"`
	// MediaIdleTimeout uint16 //推流超时时间，超过则断开链接，让设备重连

	// WaitKeyFrame      bool //是否等待关键帧，如果等待，则在收到第一个关键帧之前，忽略所有媒体流
	RemoveBanInterval time.Duration `default:"600s"` //移除禁止设备间隔
	// UdpCacheSize      int           //udp缓存大小
	LogLevel string `default:"info"` //trace, debug, info, warn, error, fatal, panic
	routes   map[string]string
	DumpPath string //dump PS流本地文件路径
	Ignores  map[string]struct{}
	tcpPorts PortManager
	udpPorts PortManager

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
	GB28181Plugin.Info("LocalAndInternalIPs", zap.Any("routes", c.routes))
}

func (c *GB28181Config) OnEvent(event any) {
	switch e := event.(type) {
	case FirstConfig:
		if c.Port.Sip != "udp:5060" {
			protocol, ports := util.Conf2Listener(c.Port.Sip)
			c.SipNetwork = protocol
			c.SipPort = ports[0]
		}
		if c.Port.Media != "tcp:58200-59200" {
			protocol, ports := util.Conf2Listener(c.Port.Media)
			c.MediaNetwork = protocol
			if len(ports) > 1 {
				c.MediaPortMin = ports[0]
				c.MediaPortMax = ports[1]
			} else {
				c.MediaPortMin = 0 
				c.MediaPortMax = 0
				c.MediaPort = ports[0]
			}
		}
		os.MkdirAll(c.DumpPath, 0766)
		c.ReadDevices()
		go c.initRoutes()
		c.startServer()
	case *Stream:
		if c.InviteMode == INVIDE_MODE_ONSUBSCRIBE {
			//流可能是回放流，stream path是device/channel/start-end形式
			streamNames := strings.Split(e.StreamName, "/")
			if channel := FindChannel(e.AppName, streamNames[0]); channel != nil {
				opt := InviteOptions{}
				if len(streamNames) > 1 {
					last := len(streamNames) - 1
					timestr := streamNames[last]
					trange := strings.Split(timestr, "-")
					if len(trange) == 2 {
						startTime := trange[0]
						endTime := trange[1]
						opt.Validate(startTime, endTime)
					}
				}
				channel.TryAutoInvite(&opt)
			}
		}
	case SEpublish:
		if channel := FindChannel(e.Target.AppName, strings.TrimSuffix(e.Target.StreamName, "/rtsp")); channel != nil {
			channel.LiveSubSP = e.Target.Path
		}
	case SEclose:
		if channel := FindChannel(e.Target.AppName, strings.TrimSuffix(e.Target.StreamName, "/rtsp")); channel != nil {
			channel.LiveSubSP = ""
		}
		if v, ok := PullStreams.LoadAndDelete(e.Target.Path); ok {
			go v.(*PullStream).Bye()
		}
	}
}

func (c *GB28181Config) IsMediaNetworkTCP() bool {
	return strings.ToLower(c.MediaNetwork) == "tcp"
}

var conf GB28181Config

var GB28181Plugin = InstallPlugin(&conf)
var PullStreams sync.Map //拉流
