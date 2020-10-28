package gb28181

import (
	. "github.com/Monibuca/engine/v2"
	"github.com/Monibuca/engine/v2/util"
	"github.com/Monibuca/plugin-gb28181/transaction"
	rtp "github.com/Monibuca/plugin-rtp"
	. "github.com/logrusorgru/aurora"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Devices sync.Map
var config = struct {
	Serial     string
	Realm      string
	ListenAddr string
	Expires    int
	AutoInvite bool
}{"34020000002000000001", "3402000000", "127.0.0.1:5060", 3600, true}

func init() {
	InstallPlugin(&PluginConfig{
		Name:   "GB28181",
		Config: &config,
		Type:   PLUGIN_PUBLISHER,
		Run:    run,
	})
}

func run() {
	ipAddr, err := net.ResolveUDPAddr("", config.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	Print(Green("server gb28181 start at"), BrightBlue(config.ListenAddr))
	config := &transaction.Config{
		SipIP:             ipAddr.IP.String(),
		SipPort:           uint16(ipAddr.Port),
		SipNetwork:        "UDP",
		Serial:            config.Serial,
		Realm:             config.Realm,
		AckTimeout:        10,
		MediaIP:           ipAddr.IP.String(),
		RegisterValidity:  config.Expires,
		RegisterInterval:  60,
		HeartbeatInterval: 60,
		HeartbeatRetry:    3,

		AudioEnable:      true,
		WaitKeyFrame:     true,
		MediaPortMin:     58200,
		MediaPortMax:     58300,
		MediaIdleTimeout: 30,
	}
	s := transaction.NewCore(config)
	s.OnInvite = onPublish
	http.HandleFunc("/gb28181/list", func(w http.ResponseWriter, r *http.Request) {
		sse := util.NewSSE(w, r.Context())
		for {
			var list []*transaction.Device
			s.Devices.Range(func(key, value interface{}) bool {
				list = append(list, value.(*transaction.Device))
				return true
			})
			sse.WriteJSON(list)
			select {
			case <-time.After(time.Second * 5):
			case <-sse.Done():
				return
			}
		}
	})
	http.HandleFunc("/gb28181/control", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := r.URL.Query().Get("id")
		channel ,err:= strconv.Atoi(r.URL.Query().Get("channel"))
		if err!=nil{
			w.WriteHeader(404)
		}
		ptzcmd := r.URL.Query().Get("ptzcmd")
		if v, ok := s.Devices.Load(id); ok {
			w.WriteHeader(v.(*transaction.Device).Control(channel,ptzcmd))
		} else {
			w.WriteHeader(404)
		}
	})
	http.HandleFunc("/gb28181/invite", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := r.URL.Query().Get("id")
		channel ,err:= strconv.Atoi(r.URL.Query().Get("channel"))
		if err != nil {
			w.WriteHeader(404)
		}
		if v, ok := s.Devices.Load(id); ok {
			w.WriteHeader(v.(*transaction.Device).Invite(channel))
		} else {
			w.WriteHeader(404)
		}
	})
	http.HandleFunc("/gb28181/bye", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := r.URL.Query().Get("id")
		channel ,err:= strconv.Atoi(r.URL.Query().Get("channel"))
		if err != nil {
			w.WriteHeader(404)
		}
		if v, ok := s.Devices.Load(id); ok {
			w.WriteHeader(v.(*transaction.Device).Bye(channel))
		} else {
			w.WriteHeader(404)
		}
	})
	s.Start()
}

func onPublish(channel *transaction.Channel) (port int) {
	rtpPublisher := new(rtp.RTP_PS)
	if !rtpPublisher.Publish("gb28181/" + channel.DeviceID) {
		return
	}
	rtpPublisher.Type = "GB28181"
	addr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	networkBuffer := 1048576
	if err := conn.SetReadBuffer(networkBuffer); err != nil {
		Printf("udp server video conn set read buffer error, %v", err)
	}
	if err := conn.SetWriteBuffer(networkBuffer); err != nil {
		Printf("udp server video conn set write buffer error, %v", err)
	}
	la := conn.LocalAddr().String()
	strPort := la[strings.LastIndex(la, ":")+1:]
	if port, err = strconv.Atoi(strPort); err != nil {
		return
	}
	go func() {
		bufUDP := make([]byte, 1048576)
		Printf("udp server start listen video port[%d]", port)
		defer Printf("udp server stop listen video port[%d]", port)
		for rtpPublisher.Err() == nil {
			if err = conn.SetReadDeadline(time.Now().Add(time.Second*30));err!=nil{
				return
			}
			if n, _, err := conn.ReadFromUDP(bufUDP); err == nil {
				rtpPublisher.PushPS(bufUDP[:n])
			} else {
				Println("udp server read video pack error", err)
				rtpPublisher.Close()
			}
		}
	}()
	return
}
