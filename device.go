package gb28181

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"m7s.live/plugin/gb28181/v4/sip"
	"m7s.live/plugin/gb28181/v4/transaction"
	"m7s.live/plugin/gb28181/v4/utils"
	// . "github.com/logrusorgru/aurora"
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
	*transaction.Core `json:"-"`
	config            *GB28181Config
	ID                string
	Name              string
	Manufacturer      string
	Model             string
	Owner             string
	RegisterTime      time.Time
	UpdateTime        time.Time
	LastKeepaliveAt   time.Time
	Status            string
	Channels          []*Channel
	sn                int
	from              *sip.Contact
	to                *sip.Contact
	Addr              string
	SipIP             string //暴露的IP
	MediaIP           string //Media Server 暴露的IP
	SourceAddr        net.Addr
	channelMap        map[string]*Channel
	channelMutex      sync.RWMutex
	subscriber        struct {
		CallID  string
		Timeout time.Time
	}
}

func (config *GB28181Config) StoreDevice(id string, s *transaction.Core, req *sip.Message) {
	var d *Device
	plugin.Debug("StoreDevice", zap.String("id", id))
	if _d, loaded := Devices.Load(id); loaded {
		d = _d.(*Device)
		d.UpdateTime = time.Now()
		d.from = &sip.Contact{Uri: req.StartLine.Uri, Params: make(map[string]string)}
		d.to = req.To
		d.Addr = req.SourceAdd.String()
		//TODO: Should we send  GetDeviceInf request?
		//message := d.CreateMessage(sip.MESSAGE)
		//message.Body = sip.GetDeviceInfoXML(d.ID)

		//request := &sip.Request{Message: message}
		//if newTx, err := s.Request(request); err == nil {
		//	if _, err = newTx.SipResponse(); err != nil {
		//		plugin.Debug("notify device after register,", err)
		//		return
		//	}
		//}
	} else {
		d = &Device{
			ID:           id,
			RegisterTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       string(sip.REGISTER),
			Core:         s,
			from:         &sip.Contact{Uri: req.StartLine.Uri, Params: make(map[string]string)},
			to:           req.To,
			Addr:         req.SourceAdd.String(),
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

func (d *Device) CreateMessage(Method sip.Method) (requestMsg *sip.Message) {
	d.sn++
	requestMsg = &sip.Message{
		Mode:        sip.SIP_MESSAGE_REQUEST,
		MaxForwards: 70,
		UserAgent:   "Monibuca",
		StartLine: &sip.StartLine{
			Method: Method,
			Uri:    d.to.Uri,
		}, Via: &sip.Via{
			Transport: "UDP",
			Host:      d.Core.SipIP,
			Port:      fmt.Sprintf("%d", d.SipPort),
			Params: map[string]string{
				"branch": fmt.Sprintf("z9hG4bK%s", utils.RandNumString(8)),
				"rport":  "-1", //only key,no-value
			},
		}, From: &sip.Contact{Uri: d.from.Uri, Params: map[string]string{"tag": utils.RandNumString(9)}},
		To: d.to, CSeq: &sip.CSeq{
			ID:     uint32(d.sn),
			Method: Method,
		}, CallID: utils.RandNumString(10),
		Addr: d.Addr,
	}
	var err2 error
	requestMsg.DestAdd, err2 = d.ResolveAddress(requestMsg)
	if err2 != nil {
		return nil
	}
	//intranet ip , let's resolve it with public ip
	var deviceIp, deviceSourceIP net.IP
	switch addr := requestMsg.DestAdd.(type) {
	case *net.UDPAddr:
		deviceIp = addr.IP
	case *net.TCPAddr:
		deviceIp = addr.IP
	}

	switch addr2 := d.SourceAddr.(type) {
	case *net.UDPAddr:
		deviceSourceIP = addr2.IP
	case *net.TCPAddr:
		deviceSourceIP = addr2.IP
	}
	if deviceIp.IsPrivate() && !deviceSourceIP.IsPrivate() {
		requestMsg.DestAdd = d.SourceAddr
	}
	return
}
func (d *Device) Subscribe() int {
	requestMsg := d.CreateMessage(sip.SUBSCRIBE)
	if d.subscriber.CallID != "" {
		requestMsg.CallID = d.subscriber.CallID
	}
	requestMsg.Expires = 3600
	requestMsg.Event = "Catalog"
	d.subscriber.Timeout = time.Now().Add(time.Second * time.Duration(requestMsg.Expires))
	requestMsg.ContentType = "Application/MANSCDP+xml"
	requestMsg.Contact = &sip.Contact{
		Uri: sip.NewURI(fmt.Sprintf("%s@%s:%d", d.Serial, d.SipIP, d.SipPort)),
	}
	requestMsg.Body = sip.BuildCatalogXML(d.sn, requestMsg.To.Uri.UserInfo())
	requestMsg.ContentLength = len(requestMsg.Body)

	request := &sip.Request{Message: requestMsg}
	response, err := d.Core.SipRequestForResponse(request)
	if err == nil && response != nil {
		if response.GetStatusCode() == 200 {
			d.subscriber.CallID = requestMsg.CallID
		} else {
			d.subscriber.CallID = ""
		}
		return response.GetStatusCode()
	}
	return http.StatusRequestTimeout
}

func (d *Device) Catalog() int {
	requestMsg := d.CreateMessage(sip.MESSAGE)
	requestMsg.Expires = 3600
	requestMsg.Event = "Catalog"
	d.subscriber.Timeout = time.Now().Add(time.Second * time.Duration(requestMsg.Expires))
	requestMsg.ContentType = "Application/MANSCDP+xml"
	requestMsg.Body = sip.BuildCatalogXML(d.sn, requestMsg.To.Uri.UserInfo())
	requestMsg.ContentLength = len(requestMsg.Body)

	request := &sip.Request{Message: requestMsg}
	response, err := d.Core.SipRequestForResponse(request)
	if err == nil && response != nil {
		return response.GetStatusCode()
	}
	return http.StatusRequestTimeout
}
func (d *Device) QueryDeviceInfo(req *sip.Request) {
	for i := time.Duration(5); i < 100; i++ {

		plugin.Info(fmt.Sprintf("QueryDeviceInfo:%s ipaddr:%s", d.ID, d.Addr))
		time.Sleep(time.Second * i)
		requestMsg := d.CreateMessage(sip.MESSAGE)
		requestMsg.ContentType = "Application/MANSCDP+xml"
		requestMsg.Body = sip.BuildDeviceInfoXML(d.sn, requestMsg.To.Uri.UserInfo())
		requestMsg.ContentLength = len(requestMsg.Body)
		request := &sip.Request{Message: requestMsg}

		response, _ := d.Core.SipRequestForResponse(request)
		if response != nil {

			if response.Via != nil && response.Via.Params["received"] != "" {
				d.SipIP = response.Via.Params["received"]
			}
			if response.GetStatusCode() != 200 {
				plugin.Error(fmt.Sprintf("device %s send Catalog : %d\n", d.ID, response.GetStatusCode()))
			} else {
				d.Subscribe()
				break
			}
		}
	}
}

// MobilePositionSubscribe 移动位置订阅
func (d *Device) MobilePositionSubscribe(id string, expires int, interval int) (code int) {
	mobilePosition := d.CreateMessage(sip.SUBSCRIBE)
	if d.subscriber.CallID != "" {
		mobilePosition.CallID = d.subscriber.CallID
	}
	mobilePosition.Expires = expires
	mobilePosition.Event = "presence"
	mobilePosition.Contact = &sip.Contact{
		Uri: sip.NewURI(fmt.Sprintf("%s@%s:%d", d.Serial, d.SipIP, d.SipPort)),
	}
	d.subscriber.Timeout = time.Now().Add(time.Second * time.Duration(mobilePosition.Expires))
	mobilePosition.ContentType = "Application/MANSCDP+xml"
	mobilePosition.Body = sip.BuildDevicePositionXML(d.sn, id, interval)
	mobilePosition.ContentLength = len(mobilePosition.Body)
	msg := &sip.Request{Message: mobilePosition}
	response, err := d.Core.SipRequestForResponse(msg)
	if err == nil && response != nil {
		if response.GetStatusCode() == 200 {
			d.subscriber.CallID = mobilePosition.CallID
		} else {
			d.subscriber.CallID = ""
		}
		return response.GetStatusCode()
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
