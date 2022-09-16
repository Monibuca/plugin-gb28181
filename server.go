package gb28181

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/pion/rtp/v2"
	"go.uber.org/zap"
	"m7s.live/engine/v4/util"
	"m7s.live/plugin/gb28181/v4/utils"

	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
)

var srv gosip.Server

type PortManager struct {
	recycle chan uint16
	max     uint16
	pos     uint16
	Valid   bool
}

func (pm *PortManager) Init(start, end uint16) {
	pm.pos = start
	pm.max = end
	if pm.pos > 0 && pm.max > pm.pos {
		pm.Valid = true
		pm.recycle = make(chan uint16, pm.Range())
	}
}

func (pm *PortManager) Range() uint16 {
	return pm.max - pm.pos
}

func (pm *PortManager) Recycle(p uint16) (err error) {
	select {
	case pm.recycle <- p:
		return nil
	default:
		return io.EOF //TODO: 换一个Error
	}
}

func (pm *PortManager) GetPort() (p uint16, err error) {
	select {
	case p = <-pm.recycle:
		return
	default:
		if pm.Range() > 0 {
			pm.pos++
			p = pm.pos
			return
		} else {
			return 0, io.EOF //TODO: 换一个Error
		}
	}
}

type Server struct {
	Ignores    map[string]struct{}
	publishers util.Map[uint32, *GBPublisher]
	tcpPorts   PortManager
	udpPorts   PortManager
}

const MaxRegisterCount = 3

func FindChannel(deviceId string, channelId string) (c *Channel) {
	if v, ok := Devices.Load(deviceId); ok {
		d := v.(*Device)
		d.channelMutex.RLock()
		c = d.channelMap[channelId]
		d.channelMutex.RUnlock()
	}
	return
}

var levelMap = map[string]log.Level{
	"trace": log.TraceLevel,
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
	"panic": log.PanicLevel,
}

func (config *GB28181Config) startServer() {
	config.publishers.Init()
	addr := "0.0.0.0:" + strconv.Itoa(int(config.SipPort))

	logger := utils.NewZapLogger(plugin.Logger, "GB SIP Server", nil)
	logger.SetLevel(levelMap[config.LogLevel])
	// logger := log.NewDefaultLogrusLogger().WithPrefix("GB SIP Server")
	srvConf := gosip.ServerConfig{}
	if config.SipIP != "" {
		srvConf.Host = config.SipIP
	}
	srv = gosip.NewServer(srvConf, nil, nil, logger)
	srv.OnRequest(sip.REGISTER, config.OnRegister)
	srv.OnRequest(sip.MESSAGE, config.OnMessage)
	srv.OnRequest(sip.NOTIFY, config.OnNotify)
	srv.OnRequest(sip.BYE, config.onBye)
	err := srv.Listen(strings.ToLower(config.SipNetwork), addr)
	if err != nil {
		plugin.Logger.Error("gb28181 server listen", zap.Error(err))
	} else {
		plugin.Info(fmt.Sprint(aurora.Green("Server gb28181 start at"), aurora.BrightBlue(addr)))
	}

	go config.startMediaServer()

	if config.Username != "" || config.Password != "" {
		go removeBanDevice(config)
	}
}

func (config *GB28181Config) startMediaServer() {
	if config.MediaNetwork == "tcp" {
		config.tcpPorts.Init(config.MediaPortMin, config.MediaPortMax)
		if !config.tcpPorts.Valid {
			config.listenMediaTCP()
		}
	} else {
		config.udpPorts.Init(config.MediaPortMin, config.MediaPortMax)
		if !config.udpPorts.Valid {
			config.listenMediaUDP()
		}
	}
}

func processTcpMediaConn(config *GB28181Config, conn net.Conn) {
	var rtpPacket rtp.Packet
	reader := bufio.NewReader(conn)
	lenBuf := make([]byte, 2)
	defer conn.Close()
	var err error
	for err == nil {
		if _, err = io.ReadFull(reader, lenBuf); err != nil {
			return
		}
		ps := make([]byte, binary.BigEndian.Uint16(lenBuf))
		if _, err = io.ReadFull(reader, ps); err != nil {
			return
		}
		if err := rtpPacket.Unmarshal(ps); err != nil {
			plugin.Error("gb28181 decode rtp error:", zap.Error(err))
		} else if publisher := config.publishers.Get(rtpPacket.SSRC); publisher != nil && publisher.Publisher.Err() == nil {
			publisher.PushPS(&rtpPacket)
		}
	}
}

func (config *GB28181Config) listenMediaTCP() {
	addr := ":" + strconv.Itoa(int(config.MediaPort))
	mediaAddr, _ := net.ResolveTCPAddr("tcp", addr)
	listen, err := net.ListenTCP("tcp", mediaAddr)

	if err != nil {
		plugin.Error("listen media server tcp err", zap.String("addr", addr), zap.Error(err))
		return
	}
	plugin.Info("Media tcp server start.", zap.Uint16("port", config.MediaPort))
	defer listen.Close()
	defer plugin.Info("Media tcp server stop", zap.Uint16("port", config.MediaPort))

	for {
		conn, err := listen.Accept()
		if err != nil {
			plugin.Error("Accept err=", zap.Error(err))
		}
		go processTcpMediaConn(config, conn)
	}
}

func (config *GB28181Config) listenMediaUDP() {
	var rtpPacket rtp.Packet
	networkBuffer := 1048576

	addr := ":" + strconv.Itoa(int(config.MediaPort))
	mediaAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, err := net.ListenUDP("udp", mediaAddr)

	if err != nil {
		plugin.Error("listen media server udp err", zap.String("addr", addr), zap.Error(err))
		return
	}
	bufUDP := make([]byte, networkBuffer)
	plugin.Info("Media udp server start.", zap.Uint16("port", config.MediaPort))
	defer plugin.Info("Media udp server stop", zap.Uint16("port", config.MediaPort))
	dumpLen := make([]byte, 6)
	for n, _, err := conn.ReadFromUDP(bufUDP); err == nil; n, _, err = conn.ReadFromUDP(bufUDP) {
		ps := bufUDP[:n]
		if err := rtpPacket.Unmarshal(ps); err != nil {
			plugin.Error("Decode rtp error:", zap.Error(err))
		}
		if publisher := config.publishers.Get(rtpPacket.SSRC); publisher != nil && publisher.Publisher.Err() == nil {
			if publisher.dumpFile != nil {
				util.PutBE(dumpLen[:4], n)
				if publisher.lastReceive.IsZero() {
					util.PutBE(dumpLen[4:], 0)
				} else {
					util.PutBE(dumpLen[4:], uint16(time.Since(publisher.lastReceive).Milliseconds()))
				}
				publisher.lastReceive = time.Now()
				publisher.dumpFile.Write(dumpLen)
				publisher.dumpFile.Write(ps)
			}
			publisher.PushPS(&rtpPacket)
		}
	}
}

// func queryCatalog(config *transaction.Config) {
// 	t := time.NewTicker(time.Duration(config.CatalogInterval) * time.Second)
// 	for range t.C {
// 		Devices.Range(func(key, value interface{}) bool {
// 			device := value.(*Device)
// 			if time.Since(device.UpdateTime) > time.Duration(config.RegisterValidity)*time.Second {
// 				Devices.Delete(key)
// 			} else if device.Channels != nil {
// 				go device.Catalog()
// 			}
// 			return true
// 		})
// 	}
// }

func removeBanDevice(config *GB28181Config) {
	t := time.NewTicker(time.Duration(config.RemoveBanInterval) * time.Second)
	for range t.C {
		DeviceRegisterCount.Range(func(key, value interface{}) bool {
			if value.(int) > MaxRegisterCount {
				DeviceRegisterCount.Delete(key)
			}
			return true
		})
	}
}
