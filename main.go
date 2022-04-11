package gb28181

import (
	"net"
	"strings"

	"github.com/pion/rtp"
	"go.uber.org/zap"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/config"
	"m7s.live/plugin/gb28181/v4/transaction"
)

type GB28181Config struct {
	AutoInvite     bool
	AutoCloseAfter int
	PreFetchRecord bool
	UdpCacheSize   int
	config.Publish
	Server
	transaction.Config
}

func (c *GB28181Config) OnEvent(event any) {
	switch event.(type) {
	case FirstConfig:
		c.startServer()
	}
}

func (c *GB28181Config) ServeUDP(conn *net.UDPConn) {
	var rtpPacket rtp.Packet
	networkBuffer := 1048576
	bufUDP := make([]byte, networkBuffer)

	for n, _, err := conn.ReadFromUDP(bufUDP); err == nil; n, _, err = conn.ReadFromUDP(bufUDP) {
		ps := bufUDP[:n]
		plugin.Info("get udp package:", zap.Any("buffer", ps))
		if err := rtpPacket.Unmarshal(ps); err != nil {
			plugin.Error("gb28181 decode rtp error:", zap.Error(err))
		}
		// if publisher := publishers.Get(rtpPacket.SSRC); publisher != nil && publisher.Err() == nil {
		// 	publisher.PushPS(&rtpPacket)
		// }
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
	Server: Server{
		MediaNetwork: "udp",
	},
	Config: transaction.Config{
		SipNetwork: "udp",
		SipIP:      "127.0.0.1",
		SipPort:    5060,
		Serial:     "34020000002000000001",
		Realm:      "3402000000",
		Username:   "",
		Password:   "",

		AckTimeout:        10,
		RegisterValidity:  60,
		RegisterInterval:  60,
		HeartbeatInterval: 60,
		HeartbeatRetry:    3,

		MediaIP:          "127.0.0.1",
		MediaPort:        58200,
		MediaIdleTimeout: 30,

		RemoveBanInterval: 600,
		UdpCacheSize:      0,
		LogVerbose:        false,
		AudioEnable:       true,
		WaitKeyFrame:      true,
	},
}

var plugin = InstallPlugin(conf)
