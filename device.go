package gb28181

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"m7s.live/engine/v4"
	"m7s.live/plugin/gb28181/v4/utils"

	// . "github.com/logrusorgru/aurora"
	"github.com/ghettovoice/gosip/sip"
)

const TIME_LAYOUT = "2006-01-02T15:04:05"

// Record 录像
type Record struct {
	//channel   *Channel
	DeviceID  string
	Name      string
	FilePath  string
	Address   string
	StartTime string
	EndTime   string
	Secrecy   int
	Type      string
}

func (r *Record) GetPublishStreamPath() string {
	return fmt.Sprintf("%s/%s", r.DeviceID, r.StartTime)
}

var (
	Devices             sync.Map
	DeviceNonce         sync.Map //保存nonce防止设备伪造
	DeviceRegisterCount sync.Map //设备注册次数
)

type Device struct {
	//*transaction.Core `json:"-"`
	config          *GB28181Config
	ID              string
	Name            string
	Manufacturer    string
	Model           string
	Owner           string
	RegisterTime    time.Time
	UpdateTime      time.Time
	LastKeepaliveAt time.Time
	Status          string
	Channels        []*Channel
	sn              int
	from            *sip.FromHeader
	to              *sip.ToHeader
	tx              *sip.ServerTransaction
	Addr            string
	SipIP           string //暴露的IP
	MediaIP         string //Media Server 暴露的IP
	SourceAddr      net.Addr
	channelMap      map[string]*Channel
	channelMutex    sync.RWMutex
	subscriber      struct {
		CallID  string
		Timeout time.Time
	}
}

func (config *GB28181Config) StoreDevice(id string, req sip.Request, tx *sip.ServerTransaction) {
	var d *Device
	plugin.Debug("StoreDevice", zap.String("id", id))
	if _d, loaded := Devices.Load(id); loaded {
		d = _d.(*Device)
		d.UpdateTime = time.Now()
		d.from, _ = req.From()
		d.to, _ = req.To()
		d.Addr = req.Source()
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
	} else {
		from, _ := req.From()
		to, _ := req.To()
		d = &Device{
			ID:           id,
			RegisterTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       string(sip.REGISTER),
			from:         from,
			to:           to,
			tx:           tx,
			Addr:         req.Source(),
			SipIP:        config.SipIP,
			MediaIP:      config.MediaIP,
			channelMap:   make(map[string]*Channel),
			config:       config,
		}
		Devices.Store(id, d)
		go d.Catalog()
	}
}

func (d *Device) addChannel(channel *Channel) {
	for _, c := range d.Channels {
		if c.DeviceID == channel.DeviceID {
			return
		}
	}
	d.Channels = append(d.Channels, channel)
}

func (d *Device) CheckSubStream() {
	d.channelMutex.Lock()
	defer d.channelMutex.Unlock()
	for _, c := range d.Channels {
		if s := engine.Streams.Get("sub/" + c.DeviceID); s != nil {
			c.LiveSubSP = s.Path
		} else {
			c.LiveSubSP = ""
		}
	}
}
func (d *Device) UpdateChannels(list []*Channel) {
	d.channelMutex.Lock()
	defer d.channelMutex.Unlock()
	for _, c := range list {
		if c.ParentID != "" {
			path := strings.Split(c.ParentID, "/")
			parentId := path[len(path)-1]
			if parent, ok := d.channelMap[parentId]; ok {
				if c.DeviceID != parentId {
					parent.Children = append(parent.Children, c)
				}
			} else {
				d.addChannel(c)
			}
		} else {
			d.addChannel(c)
		}
		if old, ok := d.channelMap[c.DeviceID]; ok {
			c.ChannelEx = old.ChannelEx
			if d.config.PreFetchRecord {
				n := time.Now()
				n = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local)
				if len(c.Records) == 0 || (n.Format(TIME_LAYOUT) == c.RecordStartTime &&
					n.Add(time.Hour*24-time.Second).Format(TIME_LAYOUT) == c.RecordEndTime) {
					go c.QueryRecord(n.Format(TIME_LAYOUT), n.Add(time.Hour*24-time.Second).Format(TIME_LAYOUT))
				}
			}
			if d.config.AutoInvite &&
				(c.LivePublisher == nil) {
				c.Invite("", "")
			}

		} else {
			c.ChannelEx = &ChannelEx{
				device: d,
			}
			if d.config.AutoInvite {
				c.Invite("", "")
			}
		}
		if s := engine.Streams.Get("sub/" + c.DeviceID); s != nil {
			c.LiveSubSP = s.Path
		} else {
			c.LiveSubSP = ""
		}
		d.channelMap[c.DeviceID] = c
	}
}
func (d *Device) UpdateRecord(channelId string, list []*Record) {
	d.channelMutex.RLock()
	if c, ok := d.channelMap[channelId]; ok {
		c.Records = append(c.Records, list...)
	}
	d.channelMutex.RUnlock()
}

func (d *Device) CreateRequest(Method sip.RequestMethod) (req sip.Request) {
	d.sn++

	callId := sip.CallID(utils.RandNumString(10))
	userAgent := sip.UserAgentHeader("Monibuca")
	cseq := sip.CSeq{
		SeqNo:      uint32(d.sn),
		MethodName: Method,
	}
	via := sip.ViaHeader{
		&sip.ViaHop{
			ProtocolName:    "SIP",
			ProtocolVersion: "2.0",
			Transport:       "UDP",
			Host:            d.SipIP,
			Port:            (*sip.Port)(&d.config.SipPort),
			Params:          sip.NewParams(),
		},
	}
	contact := sip.Address{
		DisplayName: sip.String{d.ID},
		Uri: &sip.SipUri{
			FUser: sip.String{d.ID},
			FHost: d.SipIP,
			FPort: (*sip.Port)(&d.config.SipPort),
		},
	}

	req = sip.NewRequest(
		"",
		Method,
		d.from.Address,
		"SIP/2.0",
		[]sip.Header{
			d.to,
			d.from,
			&callId,
			&userAgent,
			&cseq,
			&via,
			contact.AsContactHeader(),
		},
		"",
		nil,
	)

	req.SetDestination(d.Addr)
	req.SetRecipient(d.from.Address)

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

func (d *Device) Subscribe() int {
	request := d.CreateRequest(sip.SUBSCRIBE)
	if d.subscriber.CallID != "" {
		callId := sip.CallID(utils.RandNumString(10))
		request.AppendHeader(&callId)
	}
	expires := sip.Expires(3600)
	d.subscriber.Timeout = time.Now().Add(time.Second * time.Duration(expires))
	contentType := sip.ContentType("Application/MANSCDP+xml")
	request.AppendHeader(&contentType)
	request.AppendHeader(&expires)

	request.SetBody(BuildCatalogXML(d.sn, d.ID), true)

	response, err := d.SipRequestForResponse(request)
	if err == nil && response != nil {
		if response.StatusCode() == 200 {
			callId, _ := request.CallID()
			d.subscriber.CallID = string(*callId)
		} else {
			d.subscriber.CallID = ""
		}
		return int(response.StatusCode())
	}
	return http.StatusRequestTimeout
}

func (d *Device) Catalog() int {
	request := d.CreateRequest(sip.MESSAGE)
	expires := sip.Expires(3600)
	d.subscriber.Timeout = time.Now().Add(time.Second * time.Duration(expires))
	contentType := sip.ContentType("Application/MANSCDP+xml")
	request.AppendHeader(&contentType)
	request.AppendHeader(&expires)

	request.SetBody(BuildCatalogXML(d.sn, d.ID), true)

	resp, err := d.SipRequestForResponse(request)

	if err == nil && resp != nil {
		return int(resp.StatusCode())
	}
	return http.StatusRequestTimeout
}

func (d *Device) QueryDeviceInfo(req *sip.Request) {
	for i := time.Duration(5); i < 100; i++ {

		plugin.Info(fmt.Sprintf("QueryDeviceInfo:%s ipaddr:%s", d.ID, d.Addr))
		time.Sleep(time.Second * i)
		request := d.CreateRequest(sip.MESSAGE)
		contentType := sip.ContentType("Application/MANSCDP+xml")
		request.AppendHeader(&contentType)
		request.SetBody(BuildDeviceInfoXML(d.sn, d.ID), true)

		response, _ := d.SipRequestForResponse(request)
		if response != nil {
			via, _ := response.ViaHop()

			if via != nil && via.Params.Has("received") {
				received, _ := via.Params.Get("received")
				d.SipIP = received.String()
			}
			if response.StatusCode() != 200 {
				plugin.Error(fmt.Sprintf("device %s send Catalog : %d\n", d.ID, response.StatusCode()))
			} else {
				d.Subscribe()
				break
			}
		}
	}
}

func (d *Device) SipRequestForResponse(request sip.Request) (sip.Response, error) {
	return (*GetSipServer()).RequestWithContext(context.Background(), request)
}
