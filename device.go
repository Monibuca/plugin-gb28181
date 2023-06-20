package gb28181

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"m7s.live/engine/v4"
	"m7s.live/engine/v4/log"
	"m7s.live/plugin/gb28181/v4/utils"

	// . "github.com/logrusorgru/aurora"
	"github.com/ghettovoice/gosip/sip"
	myip "github.com/husanpao/ip"
)

const TIME_LAYOUT = "2006-01-02T15:04:05"

// Record 录像
type Record struct {
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

type DeviceStatus string

const (
	DeviceRegisterStatus = "REGISTER"
	DeviceRecoverStatus  = "RECOVER"
	DeviceOnlineStatus   = "ONLINE"
	DeviceOfflineStatus  = "OFFLINE"
	DeviceAlarmedStatus  = "ALARMED"
)

type Device struct {
	//*transaction.Core `json:"-" yaml:"-"`
	ID              string
	Name            string
	Manufacturer    string
	Model           string
	Owner           string
	RegisterTime    time.Time
	UpdateTime      time.Time
	LastKeepaliveAt time.Time
	Status          DeviceStatus
	sn              int
	addr            sip.Address
	sipIP           string //设备对应网卡的服务器ip
	mediaIP         string //设备对应网卡的服务器ip
	NetAddr         string
	channelMap      sync.Map
	subscriber      struct {
		CallID  string
		Timeout time.Time
	}
	lastSyncTime time.Time
	GpsTime      time.Time //gps时间
	Longitude    string    //经度
	Latitude     string    //纬度
	*log.Logger  `json:"-" yaml:"-"`
}

func (d *Device) MarshalJSON() ([]byte, error) {
	type Alias Device
	data := &struct {
		Channels []*ChannelInfo
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	d.channelMap.Range(func(key, value interface{}) bool {
		c := value.(*Channel)
		data.Channels = append(data.Channels, &c.ChannelInfo)
		return true
	})
	return json.Marshal(data)
}
func (c *GB28181Config) RecoverDevice(d *Device, req sip.Request) {
	from, _ := req.From()
	d.addr = sip.Address{
		DisplayName: from.DisplayName,
		Uri:         from.Address,
	}
	deviceIp := req.Source()
	servIp := req.Recipient().Host()
	//根据网卡ip获取对应的公网ip
	sipIP := c.routes[servIp]
	//如果相等，则服务器是内网通道.海康摄像头不支持...自动获取
	if strings.LastIndex(deviceIp, ".") != -1 && strings.LastIndex(servIp, ".") != -1 {
		if servIp[0:strings.LastIndex(servIp, ".")] == deviceIp[0:strings.LastIndex(deviceIp, ".")] || sipIP == "" {
			sipIP = servIp
		}
	}
	//如果用户配置过则使用配置的
	if c.SipIP != "" {
		sipIP = c.SipIP
	} else if sipIP == "" {
		sipIP = myip.InternalIPv4()
	}
	mediaIp := sipIP
	if c.MediaIP != "" {
		mediaIp = c.MediaIP
	}
	d.Info("RecoverDevice", zap.String("deviceIp", deviceIp), zap.String("servIp", servIp), zap.String("sipIP", sipIP), zap.String("mediaIp", mediaIp))
	d.Status = DeviceRegisterStatus
	d.sipIP = sipIP
	d.mediaIP = mediaIp
	d.NetAddr = deviceIp
	d.UpdateTime = time.Now()
}

func (c *GB28181Config) StoreDevice(id string, req sip.Request) (d *Device) {
	from, _ := req.From()
	deviceAddr := sip.Address{
		DisplayName: from.DisplayName,
		Uri:         from.Address,
	}
	deviceIp := req.Source()
	if _d, loaded := Devices.Load(id); loaded {
		d = _d.(*Device)
		d.UpdateTime = time.Now()
		d.NetAddr = deviceIp
		d.addr = deviceAddr
		d.Debug("UpdateDevice", zap.String("netaddr", d.NetAddr))
	} else {
		servIp := req.Recipient().Host()
		//根据网卡ip获取对应的公网ip
		sipIP := c.routes[servIp]
		//如果相等，则服务器是内网通道.海康摄像头不支持...自动获取
		if strings.LastIndex(deviceIp, ".") != -1 && strings.LastIndex(servIp, ".") != -1 {
			if servIp[0:strings.LastIndex(servIp, ".")] == deviceIp[0:strings.LastIndex(deviceIp, ".")] || sipIP == "" {
				sipIP = servIp
			}
		}
		//如果用户配置过则使用配置的
		if c.SipIP != "" {
			sipIP = c.SipIP
		} else if sipIP == "" {
			sipIP = myip.InternalIPv4()
		}
		mediaIp := sipIP
		if c.MediaIP != "" {
			mediaIp = c.MediaIP
		}
		d = &Device{
			ID:           id,
			RegisterTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       DeviceRegisterStatus,
			addr:         deviceAddr,
			sipIP:        sipIP,
			mediaIP:      mediaIp,
			NetAddr:      deviceIp,
			Logger:       GB28181Plugin.With(zap.String("id", id)),
		}
		d.Info("StoreDevice", zap.String("deviceIp", deviceIp), zap.String("servIp", servIp), zap.String("sipIP", sipIP), zap.String("mediaIp", mediaIp))
		Devices.Store(id, d)
		c.SaveDevices()
	}
	return
}
func (c *GB28181Config) ReadDevices() {
	if f, err := os.OpenFile("devices.json", os.O_RDONLY, 0644); err == nil {
		defer f.Close()
		var items []*Device
		if err = json.NewDecoder(f).Decode(&items); err == nil {
			for _, item := range items {
				if time.Since(item.UpdateTime) < conf.RegisterValidity {
					item.Status = "RECOVER"
					item.Logger = GB28181Plugin.With(zap.String("id", item.ID))
					Devices.Store(item.ID, item)
				}
			}
		}
	}
}
func (c *GB28181Config) SaveDevices() {
	var item []any
	Devices.Range(func(key, value any) bool {
		item = append(item, value)
		return true
	})
	if f, err := os.OpenFile("devices.json", os.O_WRONLY|os.O_CREATE, 0644); err == nil {
		defer f.Close()
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", " ")
		encoder.Encode(item)
	}
}

func (d *Device) addOrUpdateChannel(info ChannelInfo) (c *Channel) {
	if old, ok := d.channelMap.Load(info.DeviceID); ok {
		c = old.(*Channel)
		c.ChannelInfo = info
	} else {
		c = &Channel{
			device:      d,
			ChannelInfo: info,
			Logger:      d.Logger.With(zap.String("channel", info.DeviceID)),
		}
		if s := engine.Streams.Get(fmt.Sprintf("%s/%s/rtsp", c.device.ID, c.DeviceID)); s != nil {
			c.LiveSubSP = s.Path
		} else {
			c.LiveSubSP = ""
		}
		d.channelMap.Store(info.DeviceID, c)
	}
	return
}

func (d *Device) deleteChannel(DeviceID string) {
	d.channelMap.Delete(DeviceID)
}

func (d *Device) UpdateChannels(list ...ChannelInfo) {
	for _, c := range list {
		if _, ok := conf.Ignores[c.DeviceID]; ok {
			continue
		}
		//当父设备非空且存在时、父设备节点增加通道
		if c.ParentID != "" {
			path := strings.Split(c.ParentID, "/")
			parentId := path[len(path)-1]
			//如果父ID并非本身所属设备，一般情况下这是因为下级设备上传了目录信息，该信息通常不需要处理。
			// 暂时不考虑级联目录的实现
			if d.ID != parentId {
				if v, ok := Devices.Load(parentId); ok {
					parent := v.(*Device)
					parent.addOrUpdateChannel(c)
					continue
				} else {
					c.Model = "Directory " + c.Model
					c.Status = "NoParent"
				}
			}
		}
		//本设备增加通道
		channel := d.addOrUpdateChannel(c)

		if conf.InviteMode == INVIDE_MODE_AUTO {
			channel.TryAutoInvite(&InviteOptions{})
		}
		if s := engine.Streams.Get("sub/" + c.DeviceID); s != nil {
			channel.LiveSubSP = s.Path
		} else {
			channel.LiveSubSP = ""
		}
	}
}

func (d *Device) CreateRequest(Method sip.RequestMethod) (req sip.Request) {
	d.sn++

	callId := sip.CallID(utils.RandNumString(10))
	userAgent := sip.UserAgentHeader("Monibuca")
	maxForwards := sip.MaxForwards(70) //增加max-forwards为默认值 70
	cseq := sip.CSeq{
		SeqNo:      uint32(d.sn),
		MethodName: Method,
	}
	port := sip.Port(conf.SipPort)
	serverAddr := sip.Address{
		//DisplayName: sip.String{Str: d.config.Serial},
		Uri: &sip.SipUri{
			FUser: sip.String{Str: conf.Serial},
			FHost: d.sipIP,
			FPort: &port,
		},
		Params: sip.NewParams().Add("tag", sip.String{Str: utils.RandNumString(9)}),
	}
	req = sip.NewRequest(
		"",
		Method,
		d.addr.Uri,
		"SIP/2.0",
		[]sip.Header{
			serverAddr.AsFromHeader(),
			d.addr.AsToHeader(),
			&callId,
			&userAgent,
			&cseq,
			&maxForwards,
			serverAddr.AsContactHeader(),
		},
		"",
		nil,
	)

	req.SetTransport(conf.SipNetwork)
	req.SetDestination(d.NetAddr)
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
		if response.StatusCode() == http.StatusOK {
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
	//os.Stdout.Write(debug.Stack())
	request := d.CreateRequest(sip.MESSAGE)
	expires := sip.Expires(3600)
	d.subscriber.Timeout = time.Now().Add(time.Second * time.Duration(expires))
	contentType := sip.ContentType("Application/MANSCDP+xml")

	request.AppendHeader(&contentType)
	request.AppendHeader(&expires)
	request.SetBody(BuildCatalogXML(d.sn, d.ID), true)
	// 输出Sip请求设备通道信息信令
	GB28181Plugin.Sugar().Debugf("SIP->Catalog:%s", request)
	resp, err := d.SipRequestForResponse(request)
	if err == nil && resp != nil {
		GB28181Plugin.Sugar().Debugf("SIP<-Catalog Response: %s", resp.String())
		return int(resp.StatusCode())
	} else if err != nil {
		GB28181Plugin.Error("SIP<-Catalog error:", zap.Error(err))
	}
	return http.StatusRequestTimeout
}

func (d *Device) QueryDeviceInfo() {
	for i := time.Duration(5); i < 100; i++ {

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
			d.Info("QueryDeviceInfo", zap.Uint16("status code", uint16(response.StatusCode())))
			if response.StatusCode() == http.StatusOK {
				break
			}
		}
	}
}

func (d *Device) SipRequestForResponse(request sip.Request) (sip.Response, error) {
	return srv.RequestWithContext(context.Background(), request)
}

// MobilePositionSubscribe 移动位置订阅
func (d *Device) MobilePositionSubscribe(id string, expires time.Duration, interval time.Duration) (code int) {
	mobilePosition := d.CreateRequest(sip.SUBSCRIBE)
	if d.subscriber.CallID != "" {
		callId := sip.CallID(utils.RandNumString(10))
		mobilePosition.ReplaceHeaders(callId.Name(), []sip.Header{&callId})
	}
	expiresHeader := sip.Expires(expires / time.Second)
	d.subscriber.Timeout = time.Now().Add(expires)
	contentType := sip.ContentType("Application/MANSCDP+xml")
	mobilePosition.AppendHeader(&contentType)
	mobilePosition.AppendHeader(&expiresHeader)

	mobilePosition.SetBody(BuildDevicePositionXML(d.sn, id, int(interval/time.Second)), true)

	response, err := d.SipRequestForResponse(mobilePosition)
	if err == nil && response != nil {
		if response.StatusCode() == http.StatusOK {
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
	if v, ok := d.channelMap.Load(channelId); ok {
		c := v.(*Channel)
		c.GpsTime = time.Now() //时间取系统收到的时间，避免设备时间和格式问题
		c.Longitude = lng
		c.Latitude = lat
		c.Debug("update channel position success")
	} else {
		//如果未找到通道，则更新到设备上
		d.GpsTime = time.Now() //时间取系统收到的时间，避免设备时间和格式问题
		d.Longitude = lng
		d.Latitude = lat
		d.Debug("update device position success", zap.String("channelId", channelId))
	}
}

// UpdateChannelStatus 目录订阅消息处理：新增/移除/更新通道或者更改通道状态
func (d *Device) UpdateChannelStatus(deviceList []*notifyMessage) {
	for _, v := range deviceList {
		switch v.Event {
		case "ON":
			d.Debug("receive channel online notify")
			d.channelOnline(v.DeviceID)
		case "OFF":
			d.Debug("receive channel offline notify")
			d.channelOffline(v.DeviceID)
		case "VLOST":
			d.Debug("receive channel video lost notify")
			d.channelOffline(v.DeviceID)
		case "DEFECT":
			d.Debug("receive channel video defect notify")
			d.channelOffline(v.DeviceID)
		case "ADD":
			d.Debug("receive channel add notify")
			channel := ChannelInfo{
				DeviceID:     v.DeviceID,
				ParentID:     v.ParentID,
				Name:         v.Name,
				Manufacturer: v.Manufacturer,
				Model:        v.Model,
				Owner:        v.Owner,
				CivilCode:    v.CivilCode,
				Address:      v.Address,
				Port:         v.Port,
				Parental:     v.Parental,
				SafetyWay:    v.SafetyWay,
				RegisterWay:  v.RegisterWay,
				Secrecy:      v.Secrecy,
				Status:       ChannelStatus(v.Status),
			}
			d.addOrUpdateChannel(channel)
		case "DEL":
			//删除
			d.Debug("receive channel delete notify")
			d.deleteChannel(v.DeviceID)
		case "UPDATE":
			d.Debug("receive channel update notify")
			// 更新通道
			channel := ChannelInfo{
				DeviceID:     v.DeviceID,
				ParentID:     v.ParentID,
				Name:         v.Name,
				Manufacturer: v.Manufacturer,
				Model:        v.Model,
				Owner:        v.Owner,
				CivilCode:    v.CivilCode,
				Address:      v.Address,
				Port:         v.Port,
				Parental:     v.Parental,
				SafetyWay:    v.SafetyWay,
				RegisterWay:  v.RegisterWay,
				Secrecy:      v.Secrecy,
				Status:       ChannelStatus(v.Status),
			}
			d.UpdateChannels(channel)
		}
	}
}

func (d *Device) channelOnline(DeviceID string) {
	if v, ok := d.channelMap.Load(DeviceID); ok {
		c := v.(*Channel)
		c.Status = ChannelOnStatus
		c.Debug("channel online", zap.String("channelId", DeviceID))
	} else {
		d.Debug("update channel status failed, not found", zap.String("channelId", DeviceID))
	}
}

func (d *Device) channelOffline(DeviceID string) {
	if v, ok := d.channelMap.Load(DeviceID); ok {
		c := v.(*Channel)
		c.Status = ChannelOffStatus
		c.Debug("channel offline", zap.String("channelId", DeviceID))
	} else {
		d.Debug("update channel status failed, not found", zap.String("channelId", DeviceID))
	}
}
