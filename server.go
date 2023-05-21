package gb28181

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"go.uber.org/zap"
	"m7s.live/plugin/gb28181/v4/utils"

	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
)

var srv gosip.Server

const MaxRegisterCount = 3

func FindChannel(deviceId string, channelId string) (c *Channel) {
	if v, ok := Devices.Load(deviceId); ok {
		d := v.(*Device)
		if v, ok := d.channelMap.Load(channelId); ok {
			return v.(*Channel)
		}
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

func GetSipServer(transport string) gosip.Server {
	return srv
}

var sn = 0

func CreateRequest(exposedId string, Method sip.RequestMethod, recipient *sip.Address, netAddr string) (req sip.Request) {

	sn++

	callId := sip.CallID(utils.RandNumString(10))
	userAgent := sip.UserAgentHeader("Monibuca")
	cseq := sip.CSeq{
		SeqNo:      uint32(sn),
		MethodName: Method,
	}
	port := sip.Port(conf.SipPort)
	serverAddr := sip.Address{
		//DisplayName: sip.String{Str: d.config.Serial},
		Uri: &sip.SipUri{
			FUser: sip.String{Str: exposedId},
			FHost: conf.SipIP,
			FPort: &port,
		},
		Params: sip.NewParams().Add("tag", sip.String{Str: utils.RandNumString(9)}),
	}
	req = sip.NewRequest(
		"",
		Method,
		recipient.Uri,
		"SIP/2.0",
		[]sip.Header{
			serverAddr.AsFromHeader(),
			recipient.AsToHeader(),
			&callId,
			&userAgent,
			&cseq,
			serverAddr.AsContactHeader(),
		},
		"",
		nil,
	)

	req.SetTransport(conf.SipNetwork)
	req.SetDestination(netAddr)
	//fmt.Printf("构建请求参数:%s", *&req)
	// requestMsg.DestAdd, err2 = d.ResolveAddress(requestMsg)
	// if err2 != nil {
	// 	return nil
	// }
	//intranet ip , let's resolve it with public ip
	// var deviceIp, deviceSourceIP net.IP
	// switch addr := requestMsg.DestAdd.(type) {
	// case *net.UDPAddr:
	// 	deviceIp = addr.IP
	// case *net.TCPAddr:
	// 	deviceIp = addr.IP
	// }

	// switch addr2 := d.SourceAddr.(type) {
	// case *net.UDPAddr:
	// 	deviceSourceIP = addr2.IP
	// case *net.TCPAddr:
	// 	deviceSourceIP = addr2.IP
	// }
	// if deviceIp.IsPrivate() && !deviceSourceIP.IsPrivate() {
	// 	requestMsg.DestAdd = d.SourceAddr
	// }
	return
}
func RequestForResponse(transport string, request sip.Request,
	options ...gosip.RequestWithContextOption) (sip.Response, error) {
	return (GetSipServer(transport)).RequestWithContext(context.Background(), request, options...)
}

func (c *GB28181Config) startServer() {
	addr := c.ListenAddr + ":" + strconv.Itoa(int(c.SipPort))

	logger := utils.NewZapLogger(GB28181Plugin.Logger, "GB SIP Server", nil)
	logger.SetLevel(levelMap[c.LogLevel])
	// logger := log.NewDefaultLogrusLogger().WithPrefix("GB SIP Server")
	srvConf := gosip.ServerConfig{}
	if c.SipIP != "" {
		srvConf.Host = c.SipIP
	}
	srv = gosip.NewServer(srvConf, nil, nil, logger)
	srv.OnRequest(sip.REGISTER, c.OnRegister)
	srv.OnRequest(sip.MESSAGE, c.OnMessage)
	srv.OnRequest(sip.NOTIFY, c.OnNotify)
	srv.OnRequest(sip.BYE, c.OnBye)
	err := srv.Listen(strings.ToLower(c.SipNetwork), addr)
	if err != nil {
		GB28181Plugin.Logger.Error("gb28181 server listen", zap.Error(err))
	} else {
		GB28181Plugin.Info(fmt.Sprint(aurora.Green("Server gb28181 start at"), aurora.BrightBlue(addr)))
	}

	if c.MediaNetwork == "tcp" {
		c.tcpPorts.Init(c.MediaPortMin, c.MediaPortMax)
	} else {
		c.udpPorts.Init(c.MediaPortMin, c.MediaPortMax)
	}

	if c.Username != "" || c.Password != "" {
		go c.removeBanDevice()
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

func (c *GB28181Config) removeBanDevice() {
	t := time.NewTicker(c.RemoveBanInterval)
	for range t.C {
		DeviceRegisterCount.Range(func(key, value interface{}) bool {
			if value.(int) > MaxRegisterCount {
				DeviceRegisterCount.Delete(key)
			}
			return true
		})
	}
}
