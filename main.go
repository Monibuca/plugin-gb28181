package gb28181

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Monibuca/engine/v3"
	"github.com/Monibuca/plugin-gb28181/v3/sip"
	"github.com/Monibuca/plugin-gb28181/v3/transaction"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
	. "github.com/Monibuca/utils/v3"
	. "github.com/logrusorgru/aurora"
	"github.com/pion/rtp"
	"golang.org/x/net/html/charset"
)

var (
	Devices             sync.Map
	DeviceNonce         = make(map[string]string) //保存nonce防止设备伪造
	DeviceRegisterCount = make(map[string]int)    //设备注册次数
	Ignores             = make(map[string]struct{})
	publishers          Publishers
)

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

type Publishers struct {
	data map[uint32]*Publisher
	sync.RWMutex
}

func (p *Publishers) Add(key uint32, pp *Publisher) {
	p.Lock()
	p.data[key] = pp
	p.Unlock()
}
func (p *Publishers) Remove(key uint32) {
	p.Lock()
	delete(p.data, key)
	p.Unlock()
}
func (p *Publishers) Get(key uint32) *Publisher {
	p.RLock()
	defer p.RUnlock()
	return p.data[key]
}

var config = struct {
	Serial            string
	Realm             string
	ListenAddr        string
	Expires           int
	MediaPort         uint16
	AutoInvite        bool
	AutoUnPublish     bool
	Ignore            []string
	TCP               bool
	RemoveBanInterval int
	PreFetchRecord    bool
	Username          string
	Password          string
}{"34020000002000000001", "3402000000", "127.0.0.1:5060", 3600, 58200, false, true, nil, false, 600, false, "", ""}

func init() {
	engine.InstallPlugin(&engine.PluginConfig{
		Name:   "GB28181",
		Config: &config,
		Run:    run,
	})
	publishers.data = make(map[uint32]*Publisher)
}

func run() {
	ipAddr, err := net.ResolveUDPAddr("", config.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	Print(Green("server gb28181 start at"), BrightBlue(config.ListenAddr))
	for _, id := range config.Ignore {
		Ignores[id] = struct{}{}
	}
	config := &transaction.Config{
		SipIP:             ipAddr.IP.String(),
		SipPort:           uint16(ipAddr.Port),
		SipNetwork:        "UDP",
		Serial:            config.Serial,
		Realm:             config.Realm,
		Username:          config.Username,
		Password:          config.Password,
		AckTimeout:        10,
		MediaIP:           ipAddr.IP.String(),
		RegisterValidity:  config.Expires,
		RegisterInterval:  60,
		HeartbeatInterval: 60,
		HeartbeatRetry:    3,
		AudioEnable:       true,
		WaitKeyFrame:      true,
		MediaIdleTimeout:  30,
		RemoveBanInterval: config.RemoveBanInterval,
	}
	http.HandleFunc("/api/gb28181/query/records", func(w http.ResponseWriter, r *http.Request) {
		CORS(w, r)
		id := r.URL.Query().Get("id")
		channel := r.URL.Query().Get("channel")
		startTime := r.URL.Query().Get("startTime")
		endTime := r.URL.Query().Get("endTime")
		if c := FindChannel(id, channel); c != nil {
			w.WriteHeader(c.QueryRecord(startTime, endTime))
		} else {
			w.WriteHeader(404)
		}
	})
	http.HandleFunc("/api/gb28181/list", func(w http.ResponseWriter, r *http.Request) {
		CORS(w, r)
		sse := NewSSE(w, r.Context())
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
	http.HandleFunc("/api/gb28181/control", func(w http.ResponseWriter, r *http.Request) {
		CORS(w, r)
		id := r.URL.Query().Get("id")
		channel := r.URL.Query().Get("channel")
		ptzcmd := r.URL.Query().Get("ptzcmd")
		if c := FindChannel(id, channel); c != nil {
			w.WriteHeader(c.Control(ptzcmd))
		} else {
			w.WriteHeader(404)
		}
	})
	http.HandleFunc("/api/gb28181/invite", func(w http.ResponseWriter, r *http.Request) {
		CORS(w, r)
		query := r.URL.Query()
		id := query.Get("id")
		channel := r.URL.Query().Get("channel")
		startTime := query.Get("startTime")
		endTime := query.Get("endTime")
		if c := FindChannel(id, channel); c != nil {
			if startTime == "" && c.LivePublisher != nil {
				w.WriteHeader(304) //直播流已存在
			} else {
				w.WriteHeader(c.Invite(startTime, endTime))
			}
		} else {
			w.WriteHeader(404)
		}
	})
	http.HandleFunc("/api/gb28181/bye", func(w http.ResponseWriter, r *http.Request) {
		CORS(w, r)
		id := r.URL.Query().Get("id")
		channel := r.URL.Query().Get("channel")
		live := r.URL.Query().Get("live")
		if c := FindChannel(id, channel); c != nil {
			w.WriteHeader(c.Bye(live != "false"))
		} else {
			w.WriteHeader(404)
		}
	})
	s := transaction.NewCore(config)
	s.OnRegister = func(msg *sip.Message) {
		id := msg.From.Uri.UserInfo()
		storeDevice := func() {
			var d *Device

			if _d, loaded := Devices.LoadOrStore(id, &Device{
				ID:           id,
				RegisterTime: time.Now(),
				UpdateTime:   time.Now(),
				Status:       string(sip.REGISTER),
				Core:         s,
				from:         &sip.Contact{Uri: msg.StartLine.Uri, Params: make(map[string]string)},
				to:           msg.To,
				Addr:         msg.Via.GetSendBy(),
				SipIP:        config.MediaIP,
				channelMap:   make(map[string]*Channel),
			}); loaded {
				d = _d.(*Device)
				d.UpdateTime = time.Now()
				d.from = &sip.Contact{Uri: msg.StartLine.Uri, Params: make(map[string]string)}
				d.to = msg.To
				d.Addr = msg.Via.GetSendBy()
			}
		}
		// 不需要密码情况
		if config.Username == "" && config.Password == "" {
			storeDevice()
			return
		}
		sendUnauthorized := func() {
			response := msg.BuildResponseWithPhrase(401, "Unauthorized")
			if DeviceNonce[id] == "" {
				nonce := utils.RandNumString(32)
				DeviceNonce[id] = nonce
			}
			response.WwwAuthenticate = sip.NewWwwAuthenticate(s.Realm, DeviceNonce[id], sip.DIGEST_ALGO_MD5)
			s.Send(response)
		}
		// 需要密码情况 设备第一次上报，返回401和加密算法
		if msg.Authorization == nil || msg.Authorization.GetUsername() == "" {
			sendUnauthorized()
			return
		}
		// 有些摄像头没有配置用户名的地方，用户名就是摄像头自己的国标id
		username := config.Username
		if msg.Authorization.GetUsername() == id {
			username = id
		}

		if DeviceRegisterCount[id] >= MaxRegisterCount {
			s.Send(msg.BuildResponse(403))
			return
		}

		// 设备第二次上报，校验
		if !msg.Authorization.Verify(username, config.Password, config.Realm, DeviceNonce[id]) {
			sendUnauthorized()
			DeviceRegisterCount[id] += 1
			return
		}
		storeDevice()
		delete(DeviceNonce, id)
		delete(DeviceRegisterCount, id)
	}
	s.OnMessage = func(msg *sip.Message) bool {
		if v, ok := Devices.Load(msg.From.Uri.UserInfo()); ok {
			d := v.(*Device)
			if d.Status == string(sip.REGISTER) {
				d.Status = "ONLINE"
				go d.Query()
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
			err := decoder.Decode(temp)
			if err != nil {
				err = utils.DecodeGbk(temp, []byte(msg.Body))
				if err != nil {
					log.Printf("decode catelog err: %s", err)
				}
			}
			switch temp.XMLName.Local {
			case "Notify":
				switch temp.CmdType {
				case "Keeyalive":
					if d.subscriber.CallID != "" && time.Now().After(d.subscriber.Timeout) {
						go d.Subscribe()
					}
					d.CheckSubStream()
				case "Catalog":
					d.UpdateChannels(temp.DeviceList)
				}
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
	go listenMedia()
	// go queryCatalog(config)
	if config.Username != "" || config.Password != "" {
		go removeBanDevice(config)
	}
	s.Start()
}
func listenMedia() {
	networkBuffer := 1048576
	addr := ":" + strconv.Itoa(int(config.MediaPort))
	var rtpPacket rtp.Packet
	if config.TCP {
		ListenTCP(addr, func(conn net.Conn) {
			reader := bufio.NewReader(conn)
			lenBuf := make([]byte, 2)
			defer conn.Close()
			var err error
			for err == nil {
				if _, err = io.ReadFull(reader, lenBuf); err != nil {
					return
				}
				ps := make([]byte, BigEndian.Uint16(lenBuf))
				if _, err = io.ReadFull(reader, ps); err != nil {
					return
				}
				if err := rtpPacket.Unmarshal(ps); err != nil {
					Println("gb28181 decode rtp error:", err)
				} else if publisher := publishers.Get(rtpPacket.SSRC); publisher != nil && publisher.Err() == nil {
					publisher.PushPS(&rtpPacket)
				}
			}
		})
	} else {
		conn, err := ListenUDP(addr, networkBuffer)
		if err != nil {
			Printf("listen udp %s err: %v", addr, err)
			return
		}
		bufUDP := make([]byte, networkBuffer)
		Printf("udp server start listen video port[%d]", config.MediaPort)
		defer Printf("udp server stop listen video port[%d]", config.MediaPort)
		for n, _, err := conn.ReadFromUDP(bufUDP); err == nil; n, _, err = conn.ReadFromUDP(bufUDP) {
			ps := bufUDP[:n]
			if err := rtpPacket.Unmarshal(ps); err != nil {
				Println("gb28181 decode rtp error:", err)
			}
			if publisher := publishers.Get(rtpPacket.SSRC); publisher != nil && publisher.Err() == nil {
				publisher.PushPS(&rtpPacket)
			}
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
// 				go device.Subscribe()
// 			}
// 			return true
// 		})
// 	}
// }

func removeBanDevice(config *transaction.Config) {
	t := time.NewTicker(time.Duration(config.RemoveBanInterval) * time.Second)
	for range t.C {
		for id, cnt := range DeviceRegisterCount {
			if cnt >= MaxRegisterCount {
				delete(DeviceRegisterCount, id)
			}
		}
	}
}
