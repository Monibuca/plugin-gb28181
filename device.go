package gb28181

import (
	"context"
	"fmt"
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
	addr            sip.Address
	tx              *sip.ServerTransaction
	NetAddr         string
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

	from, _ := req.From()
	deviceAddr := sip.Address{
		DisplayName: from.DisplayName,
		Uri:         from.Address,
	}

	if _d, loaded := Devices.Load(id); loaded {
		d = _d.(*Device)
		d.UpdateTime = time.Now()
		d.NetAddr = req.Source()
		d.addr = deviceAddr
	} else {
		d = &Device{
			ID:           id,
			RegisterTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       string(sip.REGISTER),
			addr:         deviceAddr,
			tx:           tx,
			NetAddr:      req.Source(),
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
		if _, ok := d.config.Ignores[c.DeviceID]; ok {
			continue
		}
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
	serverAddr := sip.Address{
		DisplayName: sip.String{Str: d.config.Serial},
		Uri: &sip.SipUri{
			FUser: sip.String{Str: d.config.Serial},
			FHost: d.config.Realm,
		},
	}

	req = sip.NewRequest(
		"",
		Method,
		d.addr.Uri,
		"SIP/2.0",
		[]sip.Header{
			d.addr.AsToHeader(),
			serverAddr.AsFromHeader(),
			&callId,
			&userAgent,
			&cseq,
			serverAddr.AsContactHeader(),
		},
		"",
		nil,
	)

	req.SetTransport(d.config.SipNetwork)
	req.SetDestination(d.NetAddr)

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

		plugin.Info(fmt.Sprintf("QueryDeviceInfo:%s ipaddr:%s", d.ID, d.NetAddr))
		time.Sleep(time.Second * i)
		request := d.CreateRequest(sip.MESSAGE)
		contentType := sip.ContentType("Application/MANSCDP+xml")
		request.AppendHeader(&contentType)
		request.SetBody(BuildDeviceInfoXML(d.sn, d.ID), true)

		response, _ := d.SipRequestForResponse(request)
		if response != nil {
			// via, _ := response.ViaHop()

			// if via != nil && via.Params.Has("received") {
			// 	received, _ := via.Params.Get("received")
			// 	d.SipIP = received.String()
			// }
			if response.StatusCode() != 200 {
				plugin.Sugar().Errorf("device %s send Catalog : %d\n", d.ID, response.StatusCode())
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

// MobilePositionSubscribe 移动位置订阅
func (d *Device) MobilePositionSubscribe(id string, expires int, interval int) (code int) {
	mobilePosition := d.CreateRequest(sip.SUBSCRIBE)
	if d.subscriber.CallID != "" {
		callId := sip.CallID(utils.RandNumString(10))
		mobilePosition.ReplaceHeaders(callId.Name(), []sip.Header{&callId})
	}
	expiresHeader := sip.Expires(expires)
	d.subscriber.Timeout = time.Now().Add(time.Second * time.Duration(expires))
	contentType := sip.ContentType("Application/MANSCDP+xml")
	mobilePosition.AppendHeader(&contentType)
	mobilePosition.AppendHeader(&expiresHeader)

	mobilePosition.SetBody(BuildDevicePositionXML(d.sn, id, interval), true)

	response, err := d.SipRequestForResponse(mobilePosition)
	if err == nil && response != nil {
		if response.StatusCode() == 200 {
			callId, _ := mobilePosition.CallID()
			d.subscriber.CallID = callId.String()
		} else {
			d.subscriber.CallID = ""
		}
		return int(response.StatusCode())
	}
	return http.StatusRequestTimeout
}

// UpdateChannelPosition 更新通道GPS坐标
func (d *Device) UpdateChannelPosition(channelId string, gpsTime string, lng string, lat string) {
	if c, ok := d.channelMap[channelId]; ok {
		c.ChannelEx.GpsTime, _ = time.ParseInLocation("2006-01-02 15:04:05", gpsTime, time.Local)
		c.ChannelEx.Longitude = lng
		c.ChannelEx.Latitude = lat
		plugin.Sugar().Debugf("更新通道[%s]坐标成功\n", c.Name)
	} else {
		plugin.Sugar().Debugf("更新失败，未找到通道[%s]\n", channelId)
	}
}

// UpdateChannelStatus 目录订阅消息处理：新增/移除/更新通道或者更改通道状态
func (d *Device) UpdateChannelStatus(deviceList []*notifyMessage) {
	for _, v := range deviceList {
		switch v.Event {
		case "ON":
			plugin.Debug("收到通道上线通知")
			d.channelOnline(v.DeviceID)
		case "OFF":
			plugin.Debug("收到通道离线通知")
			d.channelOffline(v.DeviceID)
		case "VLOST":
			plugin.Debug("收到通道视频丢失通知")
			d.channelOffline(v.DeviceID)
		case "DEFECT":
			plugin.Debug("收到通道故障通知")
			d.channelOffline(v.DeviceID)
		case "ADD":
			plugin.Debug("收到通道新增通知")
			channel := Channel{
				DeviceID:     v.DeviceID,
				ParentID:     v.ParentID,
				Name:         v.Name,
				Manufacturer: v.Manufacturer,
				Model:        v.Model,
				Owner:        v.Owner,
				CivilCode:    v.CivilCode,
				Address:      v.Address,
				Parental:     v.Parental,
				SafetyWay:    v.SafetyWay,
				RegisterWay:  v.RegisterWay,
				Secrecy:      v.Secrecy,
				Status:       v.Status,
			}
			d.addChannel(&channel)
		case "DEL":
			//删除
			plugin.Debug("收到通道删除通知")
			delete(d.channelMap, v.DeviceID)
		case "UPDATE":
			plugin.Debug("收到通道更新通知")
			// 更新通道
			channel := &Channel{
				DeviceID:     v.DeviceID,
				ParentID:     v.ParentID,
				Name:         v.Name,
				Manufacturer: v.Manufacturer,
				Model:        v.Model,
				Owner:        v.Owner,
				CivilCode:    v.CivilCode,
				Address:      v.Address,
				Parental:     v.Parental,
				SafetyWay:    v.SafetyWay,
				RegisterWay:  v.RegisterWay,
				Secrecy:      v.Secrecy,
				Status:       v.Status,
			}
			channels := []*Channel{channel}
			d.UpdateChannels(channels)
		}
	}
}

func (d *Device) channelOnline(DeviceID string) {
	if c, ok := d.channelMap[DeviceID]; ok {
		c.Status = "ON"
		plugin.Sugar().Debugf("通道[%s]在线\n", c.Name)
	} else {
		plugin.Sugar().Debugf("更新通道[%s]状态失败，未找到\n", DeviceID)
	}
}

func (d *Device) channelOffline(DeviceID string) {
	if c, ok := d.channelMap[DeviceID]; ok {
		c.Status = "OFF"
		plugin.Sugar().Debugf("通道[%s]离线\n", c.Name)
	} else {
		plugin.Sugar().Debugf("更新通道[%s]状态失败，未找到\n", DeviceID)
	}
}
