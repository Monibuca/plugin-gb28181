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
	Ignores    = make(map[string]struct{})
	publishers Publishers
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
	AutoCloseAfter    int
	Ignore            []string
	TCP               bool
	TCPMediaPortNum   uint16
	RemoveBanInterval int
	PreFetchRecord    bool
	Username          string
	Password          string
	UdpCacheSize      int //udp排序缓存
}{"34020000002000000001", "3402000000", "127.0.0.1:5060", 3600, 58200, false, -1, nil, false, 1, 600, false, "", "", 0}

func init() {
	pc := engine.PluginConfig{
		Name:   "GB28181",
		Config: &config,
	}
	pc.Install(run)
	publishers.data = make(map[uint32]*Publisher)
}
func onBye(req *sip.Request, tx *sip.GBTx) {
	response := &sip.Response{req.BuildOK()}
	_ = tx.Respond(response)
}
func storeDevice(id string, mediaIP string, s *transaction.Core, req *sip.Message) {
	var d *Device

	if _d, loaded := Devices.LoadOrStore(id, &Device{
		ID:           id,
		RegisterTime: time.Now(),
		UpdateTime:   time.Now(),
		Status:       string(sip.REGISTER),
		Core:         s,
		from:         &sip.Contact{Uri: req.StartLine.Uri, Params: make(map[string]string)},
		to:           req.To,
		Addr:         req.Via.GetSendBy(),
		SipIP:        mediaIP,
		channelMap:   make(map[string]*Channel),
	}); loaded {
		d = _d.(*Device)
		d.UpdateTime = time.Now()
		d.from = &sip.Contact{Uri: req.StartLine.Uri, Params: make(map[string]string)}
		d.to = req.To
		d.Addr = req.Via.GetSendBy()

		//TODO: Should we send  GetDeviceInf request?
		//message := d.CreateMessage(sip.MESSAGE)
		//message.Body = sip.GetDeviceInfoXML(d.ID)

		//request := &sip.Request{Message: message}
		//if newTx, err := s.Request(request); err == nil {
		//	if _, err = newTx.SipResponse(); err != nil {
		//		Println("notify device after register,", err)
		//		return
		//	}
		//}

	}
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
	useTCP := config.TCP
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
		UdpCacheSize:      config.UdpCacheSize,
	}

	s := transaction.NewCore(config)
	s.OnRegister = func(req *sip.Request, tx *sip.GBTx) {
		id := req.From.Uri.UserInfo()

		passAuth := false
		// 不需要密码情况
		if config.Username == "" && config.Password == "" {
			passAuth = true
		} else {
			// 需要密码情况 设备第一次上报，返回401和加密算法
			if req.Authorization != nil || req.Authorization.GetUsername() != "" {
				// 有些摄像头没有配置用户名的地方，用户名就是摄像头自己的国标id
				var username string
				if req.Authorization.GetUsername() == id {
					username = id
				} else {
					username = config.Username
				}

				if DeviceRegisterCount[id] >= MaxRegisterCount {
					var response sip.Response
					response.Message = req.BuildResponse(http.StatusForbidden)
					_ = tx.Respond(&response)
					return
				} else {
					// 设备第二次上报，校验
					if req.Authorization.Verify(username, config.Password, config.Realm, DeviceNonce[id]) {
						passAuth = true
					} else {
						DeviceRegisterCount[id]++
					}
				}
			}

		}
		if passAuth {
			storeDevice(id, config.MediaIP, s, req.Message)
			delete(DeviceNonce, id)
			delete(DeviceRegisterCount, id)
			m := req.BuildOK()
			resp := &sip.Response{Message: m}
			_ = tx.Respond(resp)
		} else {
			var response sip.Response
			response.Message = req.BuildResponseWithPhrase(401, "Unauthorized")
			if DeviceNonce[id] == "" {
				nonce := utils.RandNumString(32)
				DeviceNonce[id] = nonce
			}
			response.WwwAuthenticate = sip.NewWwwAuthenticate(s.Realm, DeviceNonce[id], sip.DIGEST_ALGO_MD5)
			response.SourceAdd = req.DestAdd
			response.DestAdd = req.SourceAdd
			_ = tx.Respond(&response)
		}
	}
	s.OnMessage = func(req *sip.Request, tx *sip.GBTx) {

		if v, ok := Devices.Load(req.From.Uri.UserInfo()); ok {
			d := v.(*Device)
			if d.Status == string(sip.REGISTER) {
				d.Status = "ONLINE"
				go d.Query(req)
			}
			d.UpdateTime = time.Now()
			temp := &struct {
				XMLName    xml.Name
				CmdType    string
				DeviceID   string
				DeviceList []*Channel `xml:"DeviceList>Item"`
				RecordList []*Record  `xml:"RecordList>Item"`
			}{}
			decoder := xml.NewDecoder(bytes.NewReader([]byte(req.Body)))
			decoder.CharsetReader = charset.NewReaderLabel
			err := decoder.Decode(temp)
			if err != nil {
				err = utils.DecodeGbk(temp, []byte(req.Body))
				if err != nil {
					log.Printf("decode catelog err: %s", err)
				}
			}
			switch temp.CmdType {
			case "Keeyalive":
				if d.subscriber.CallID != "" && time.Now().After(d.subscriber.Timeout) {
					go d.Subscribe(req)
				}
				d.CheckSubStream()
			case "Catalog":
				d.UpdateChannels(temp.DeviceList)

			case "RecordInfo":
				d.UpdateRecord(temp.DeviceID, temp.RecordList)

			}
			response := &sip.Response{req.BuildOK()}
			tx.Respond(response)
		}
	}
	s.RegistHandler(sip.REGISTER, s.OnRegister)
	s.RegistHandler(sip.MESSAGE, s.OnMessage)
	s.RegistHandler(sip.BYE, onBye)

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
	if useTCP {
		listenMediaTCP()
	} else {
		go listenMediaUDP()
	}
	// go queryCatalog(config)
	if config.Username != "" || config.Password != "" {
		go removeBanDevice(config)
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

	s.StartAndWait()
}
func listenMediaTCP() {
	for i := uint16(0); i < config.TCPMediaPortNum; i++ {
		addr := ":" + strconv.Itoa(int(config.MediaPort+i))
		go ListenTCP(addr, func(conn net.Conn) {
			var rtpPacket rtp.Packet
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
	}
}
func listenMediaUDP() {
	var rtpPacket rtp.Packet
	networkBuffer := 1048576
	addr := ":" + strconv.Itoa(int(config.MediaPort))
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
