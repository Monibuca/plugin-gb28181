package gb28181

import (
	"bytes"
	"embed"
	"encoding/xml"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Monibuca/plugin-gb28181/sip"
	"golang.org/x/net/html/charset"

	. "github.com/Monibuca/engine/v2"
	"github.com/Monibuca/engine/v2/util"
	"github.com/Monibuca/plugin-gb28181/transaction"
	rtp "github.com/Monibuca/plugin-rtp"
	. "github.com/logrusorgru/aurora"
)

var Devices sync.Map
var config = struct {
	Serial       string
	Realm        string
	ListenAddr   string
	Expires      int
	AutoInvite   bool
	MediaPortMin uint16
	MediaPortMax uint16
}{"34020000002000000001", "3402000000", "127.0.0.1:5060", 3600, true, 58200, 58300}

//go:embed ui/*
//go:embed README.md
var ui embed.FS

func init() {
	InstallPlugin(&PluginConfig{
		Name:   "GB28181",
		Config: &config,
		Type:   PLUGIN_PUBLISHER,
		Run:    run,
		UIFile: &ui,
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
		MediaPortMin:     config.MediaPortMin,
		MediaPortMax:     config.MediaPortMax,
		MediaIdleTimeout: 30,
	}

	http.HandleFunc("/gb28181/query/records", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := r.URL.Query().Get("id")
		channel, err := strconv.Atoi(r.URL.Query().Get("channel"))
		if err != nil {
			w.WriteHeader(404)
		}
		startTime := r.URL.Query().Get("startTime")
		endTime := r.URL.Query().Get("endTime")
		if v, ok := Devices.Load(id); ok {
			w.WriteHeader(v.(*Device).QueryRecord(channel, startTime, endTime))
		} else {
			w.WriteHeader(404)
		}
	})
	http.HandleFunc("/gb28181/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		sse := util.NewSSE(w, r.Context())
		for {
			var list []*Device
			Devices.Range(func(key, value interface{}) bool {
				device := value.(*Device)
				if time.Since(device.UpdateTime) > time.Duration(config.RegisterValidity)*time.Second {
					Devices.Delete(key)
				} else {
					list = append(list, device)
				}
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
		channel, err := strconv.Atoi(r.URL.Query().Get("channel"))
		if err != nil {
			w.WriteHeader(404)
		}
		ptzcmd := r.URL.Query().Get("ptzcmd")
		if v, ok := Devices.Load(id); ok {
			w.WriteHeader(v.(*Device).Control(channel, ptzcmd))
		} else {
			w.WriteHeader(404)
		}
	})
	http.HandleFunc("/gb28181/invite", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		query := r.URL.Query()
		id := query.Get("id")
		channel, err := strconv.Atoi(query.Get("channel"))
		startTime := query.Get("startTime")
		endTime := query.Get("endTime")
		f := query.Get("f")
		if startTime == "" {
			startTime = "0"
		}
		if endTime == "" {
			endTime = "0"
		}
		if err != nil {
			w.WriteHeader(404)
		}
		if v, ok := Devices.Load(id); ok {
			w.WriteHeader(v.(*Device).Invite(channel, startTime, endTime, f))
		} else {
			w.WriteHeader(404)
		}
	})
	http.HandleFunc("/gb28181/bye", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := r.URL.Query().Get("id")
		channel, err := strconv.Atoi(r.URL.Query().Get("channel"))
		if err != nil {
			w.WriteHeader(404)
		}
		if v, ok := Devices.Load(id); ok {
			w.WriteHeader(v.(*Device).Bye(channel))
		} else {
			w.WriteHeader(404)
		}
	})
	s := transaction.NewCore(config)
	s.OnRegister = func(msg *sip.Message) {
		Devices.Store(msg.From.Uri.UserInfo(), &Device{
			ID:           msg.From.Uri.UserInfo(),
			RegisterTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       string(sip.REGISTER),
			Core:         s,
			from:         &sip.Contact{Uri: msg.StartLine.Uri, Params: make(map[string]string)},
			to:           msg.To,
			Addr:         msg.Via.GetSendBy(),
			SipIP:        config.MediaIP,
		})
	}
	s.OnMessage = func(msg *sip.Message) bool {
		if v, ok := Devices.Load(msg.From.Uri.UserInfo()); ok {
			d := v.(*Device)
			if d.Status == string(sip.REGISTER) {
				d.Status = "ONLINE"
			}
			d.UpdateTime = time.Now()
			temp := &struct {
				XMLName    xml.Name
				CmdType    string
				DeviceID   string
				DeviceList []*Channel `xml:"DeviceList>Item"`
				RecordList []*Record  `xml:"RecordList>Item"`
			}{}
			decoder := xml.NewDecoder(bytes.NewReader([]byte(msg.Body)))
			decoder.CharsetReader = charset.NewReaderLabel
			decoder.Decode(temp)
			switch temp.XMLName.Local {
			case "Notify":
				go d.Query()
			case "Response":
				switch temp.CmdType {
				case "Catalog":
					d.UpdateChannels(temp.DeviceList)
				case "RecordInfo":
					d.UpdateRecord(temp.DeviceID, temp.RecordList)
				}
			}
			return true
		}
		return false
	}
	//OnStreamClosedHooks.AddHook(func(stream *Stream) {
	//	Devices.Range(func(key, value interface{}) bool {
	//		device:=value.(*Device)
	//		for _,channel := range device.Channels {
	//			if stream.StreamPath == channel.RecordSP {
	//
	//			}
	//		}
	//	})
	//})
	s.Start()
}

func (d *Device) publish(name string) (port int, publisher *rtp.RTP_PS) {
	publisher = new(rtp.RTP_PS)
	if !publisher.Publish(name) {
		return
	}
	defer func() {
		if port == 0 {
			publisher.Close()
		}
	}()
	publisher.Type = "GB28181"
	publisher.AutoUnPublish = true
	var conn *net.UDPConn
	var err error
	rang := int(config.MediaPortMax - config.MediaPortMin)
	for count := rang; count > 0; count-- {
		randNum := rand.Intn(rang)
		port = int(config.MediaPortMin) + randNum
		addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
		conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			continue
		} else {
			break
		}
	}
	if err != nil {
		return
	}
	networkBuffer := 1048576
	if err = conn.SetReadBuffer(networkBuffer); err != nil {
		Printf("udp server video conn set read buffer error, %v", err)
	}
	if err = conn.SetWriteBuffer(networkBuffer); err != nil {
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
		for publisher.Err() == nil {
			if err = conn.SetReadDeadline(time.Now().Add(time.Second * 30)); err != nil {
				return
			}
			if n, _, err := conn.ReadFromUDP(bufUDP); err == nil {
				publisher.PushPS(bufUDP[:n])
			} else {
				Println("udp server read video pack error", err)
				publisher.Close()
				if !publisher.AutoUnPublish {
					for _, channel := range d.Channels {
						if channel.LiveSP == name {
							channel.LiveSP = ""
							channel.Connected = false
							channel.Bye(channel.inviteRes)
							break
						}
					}
				}
			}
		}
		conn.Close()
		if publisher.AutoUnPublish {
			for _, channel := range d.Channels {
				if channel.RecordSP == name {
					channel.RecordSP = ""
					channel.Bye(channel.recordInviteRes)
					break
				}
			}
		}
	}()
	return
}
