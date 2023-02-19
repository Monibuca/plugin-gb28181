package gb28181

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"m7s.live/engine/v4"
	"m7s.live/plugin/gb28181/v4/utils"

	// . "github.com/logrusorgru/aurora"
	"github.com/ghettovoice/gosip/sip"
	myip "github.com/husanpao/ip"
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
	ID              string
	Name            string
	Manufacturer    string
	Model           string
	Owner           string
	RegisterTime    time.Time
	UpdateTime      time.Time
	LastKeepaliveAt time.Time
	Status          string
	sn              int
	addr            sip.Address
	sipIP           string //设备对应网卡的服务器ip
	mediaIP         string //设备对应网卡的服务器ip
	NetAddr         string
	channelMap      map[string]*Channel
	channelMutex    sync.RWMutex
	subscriber      struct {
		CallID  string
		Timeout time.Time
	}
	lastSyncTime time.Time
	GpsTime      time.Time //gps时间
	Longitude    string    //经度
	Latitude     string    //纬度
}

func (d *Device) MarshalJSON() ([]byte, error) {
	type Alias Device
	return json.Marshal(&struct {
		Channels []*Channel
		*Alias
	}{
		Channels: maps.Values(d.channelMap),
		Alias:    (*Alias)(d),
	})
}
func (config *GB28181Config) RecoverDevice(d *Device, req sip.Request) {
	from, _ := req.From()
	d.addr = sip.Address{
		DisplayName: from.DisplayName,
		Uri:         from.Address,
	}
	deviceIp := req.Source()
	servIp := req.Recipient().Host()
	//根据网卡ip获取对应的公网ip
	sipIP := config.routes[servIp]
	//如果相等，则服务器是内网通道.海康摄像头不支持...自动获取
	if strings.LastIndex(deviceIp, ".") != -1 && strings.LastIndex(servIp, ".") != -1 {
		if servIp[0:strings.LastIndex(servIp, ".")] == deviceIp[0:strings.LastIndex(deviceIp, ".")] || sipIP == "" {
			sipIP = servIp
		}
	}
	//如果用户配置过则使用配置的
	if config.SipIP != "" {
		sipIP = config.SipIP
	} else if sipIP == "" {
		sipIP = myip.InternalIPv4()
	}
	mediaIp := sipIP
	if config.MediaIP != "" {
		mediaIp = config.MediaIP
	}
	plugin.Info("RecoverDevice", zap.String("id", d.ID), zap.String("deviceIp", deviceIp), zap.String("servIp", servIp), zap.String("sipIP", sipIP), zap.String("mediaIp", mediaIp))
	d.Status = string(sip.REGISTER)
	d.sipIP = sipIP
	d.mediaIP = mediaIp
	d.NetAddr = deviceIp
	d.UpdateTime = time.Now()
	d.channelMap = make(map[string]*Channel)
}
func (config *GB28181Config) StoreDevice(id string, req sip.Request) *Device {
	var d *Device
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
		plugin.Debug("UpdateDevice", zap.String("id", id), zap.String("netaddr", d.NetAddr))
	} else {
		servIp := req.Recipient().Host()
		//根据网卡ip获取对应的公网ip
		sipIP := config.routes[servIp]
		//如果相等，则服务器是内网通道.海康摄像头不支持...自动获取
		if strings.LastIndex(deviceIp, ".") != -1 && strings.LastIndex(servIp, ".") != -1 {
			if servIp[0:strings.LastIndex(servIp, ".")] == deviceIp[0:strings.LastIndex(deviceIp, ".")] || sipIP == "" {
				sipIP = servIp
			}
		}
		//如果用户配置过则使用配置的
		if config.SipIP != "" {
			sipIP = config.SipIP
		} else if sipIP == "" {
			sipIP = myip.InternalIPv4()
		}
		mediaIp := sipIP
		if config.MediaIP != "" {
			mediaIp = config.MediaIP
		}
		plugin.Info("StoreDevice", zap.String("id", id), zap.String("deviceIp", deviceIp), zap.String("servIp", servIp), zap.String("sipIP", sipIP), zap.String("mediaIp", mediaIp))
		d = &Device{
			ID:           id,
			RegisterTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       string(sip.REGISTER),
			addr:         deviceAddr,
			sipIP:        sipIP,
			mediaIP:      mediaIp,
			NetAddr:      deviceIp,
			channelMap:   make(map[string]*Channel),
		}
		Devices.Store(id, d)
		SaveDevices()
	}
	return d
}
func ReadDevices() {
	if f, err := os.OpenFile("devices.json", os.O_RDONLY, 0644); err == nil {
		defer f.Close()
		var items []*Device
		if err = json.NewDecoder(f).Decode(&items); err == nil {
			for _, item := range items {
				if time.Since(item.UpdateTime) < conf.RegisterValidity {
					item.Status = "RECOVER"
					Devices.Store(item.ID, item)
				}
			}
		}
	}
}
func SaveDevices() {
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

func (d *Device) addOrUpdateChannel(channel *Channel) {
	d.channelMutex.Lock()
	defer d.channelMutex.Unlock()
	channel.device = d
	var oldLock *sync.Mutex
	if old, ok := d.channelMap[channel.DeviceID]; ok {
		//复制锁指针
		oldLock = old.liveInviteLock
	}
	if oldLock == nil {
		channel.liveInviteLock = &sync.Mutex{}
	} else {
		channel.liveInviteLock = oldLock
	}
	d.channelMap[channel.DeviceID] = channel
}

func (d *Device) deleteChannel(DeviceID string) {
	d.channelMutex.Lock()
	defer d.channelMutex.Unlock()
	delete(d.channelMap, DeviceID)
}

func (d *Device) CheckSubStream() {
	d.channelMutex.Lock()
	defer d.channelMutex.Unlock()
	for _, c := range d.channelMap {
		if s := engine.Streams.Get("sub/" + c.DeviceID); s != nil {
			c.LiveSubSP = s.Path
		} else {
			c.LiveSubSP = ""
		}
	}
}
func (d *Device) UpdateChannels(list []*Channel) {

	for _, c := range list {
		if _, ok := conf.Ignores[c.DeviceID]; ok {
			continue
		}
		//当父设备非空且存在时、父设备节点增加通道
		if c.ParentID != "" {
			path := strings.Split(c.ParentID, "/")
			parentId := path[len(path)-1]
			if c.DeviceID != parentId {
				if v, ok := Devices.Load(parentId); ok {
					parent := v.(*Device)
					parent.addOrUpdateChannel(c)
					continue
				}
			}
		}
		//本设备增加通道
		d.addOrUpdateChannel(c)

		//预取和邀请
		if conf.PreFetchRecord {
			n := time.Now()
			n = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local)
			if len(c.Records) == 0 || (n.Format(TIME_LAYOUT) == c.RecordStartTime &&
				n.Add(time.Hour*24-time.Second).Format(TIME_LAYOUT) == c.RecordEndTime) {
				go c.QueryRecord(n.Format(TIME_LAYOUT), n.Add(time.Hour*24-time.Second).Format(TIME_LAYOUT))
			}
		}
		if conf.AutoInvite && (c.LivePublisher == nil) {
			go c.Invite(InviteOptions{})
		}
		if s := engine.Streams.Get("sub/" + c.DeviceID); s != nil {
			c.LiveSubSP = s.Path
		} else {
			c.LiveSubSP = ""
		}
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
	//os.Stdout.Write(debug.Stack())
	request := d.CreateRequest(sip.MESSAGE)
	expires := sip.Expires(3600)
	d.subscriber.Timeout = time.Now().Add(time.Second * time.Duration(expires))
	contentType := sip.ContentType("Application/MANSCDP+xml")

	request.AppendHeader(&contentType)
	request.AppendHeader(&expires)
	request.SetBody(BuildCatalogXML(d.sn, d.ID), true)
	// 输出Sip请求设备通道信息信令
	plugin.Sugar().Debugf("SIP->Catalog:%s", request)
	resp, err := d.SipRequestForResponse(request)
	if err == nil && resp != nil {
		return int(resp.StatusCode())
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
			plugin.Info(fmt.Sprintf("QueryDeviceInfo:%s ipaddr:%s response code:%d", d.ID, d.NetAddr, response.StatusCode()))
			if response.StatusCode() == 200 {
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
		c.ChannelEx.GpsTime = time.Now() //时间取系统收到的时间，避免设备时间和格式问题
		c.ChannelEx.Longitude = lng
		c.ChannelEx.Latitude = lat
		plugin.Sugar().Debugf("更新通道[%s]坐标成功\n", c.Name)
	} else {
		//如果未找到通道，则更新到设备上
		d.GpsTime = time.Now() //时间取系统收到的时间，避免设备时间和格式问题
		d.Longitude = lng
		d.Latitude = lat
		plugin.Sugar().Debugf("未找到通道[%s]，更新设备[%s]坐标成功\n", channelId, d.ID)
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
				Port:         v.Port,
				Parental:     v.Parental,
				SafetyWay:    v.SafetyWay,
				RegisterWay:  v.RegisterWay,
				Secrecy:      v.Secrecy,
				Status:       v.Status,
			}
			d.addOrUpdateChannel(&channel)
		case "DEL":
			//删除
			plugin.Debug("收到通道删除通知")
			d.deleteChannel(v.DeviceID)
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
				Port:         v.Port,
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
