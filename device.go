package gb28181

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Monibuca/plugin-gb28181/transaction"

	"github.com/Monibuca/plugin-gb28181/sip"
	"github.com/Monibuca/plugin-gb28181/utils"
)

type ChannelEx struct {
	device          *Device
	inviteRes       *sip.Message
	recordInviteRes *sip.Message
	RecordSP        string //正在播放录像的StreamPath
	LiveSP          string //实时StreamPath
	Connected       bool
	Records         []*Record
}

// Channel 通道
type Channel struct {
	DeviceID     string
	Name         string
	Manufacturer string
	Model        string
	Owner        string
	CivilCode    string
	Address      string
	Parental     int
	SafetyWay    int
	RegisterWay  int
	Secrecy      int
	Status       string
	ChannelEx    //自定义属性
}

// func (c *Channel) MarshalJSON() ([]byte, error) {
// 	var data = map[string]interface{}{
// 		"DeviceID":     c.DeviceID,
// 		"Name":         c.Name,
// 		"Manufacturer": c.Manufacturer,
// 		"Address":      c.Address,
// 		"Status":       c.Status,
// 		"RecordSP":     c.RecordSP,
// 		"LiveSP":       c.LiveSP,
// 		"Records":      c.Records,
// 		"Connected":    c.Connected,
// 	}
// 	return json.Marshal(data)
// }

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

type Device struct {
	*transaction.Core `json:"-"`
	ID                string
	RegisterTime      time.Time
	UpdateTime        time.Time
	Status            string
	Channels          []*Channel
	sn                int
	from              *sip.Contact
	to                *sip.Contact
	Addr              string
	SipIP             string //暴露的IP
}

func (d *Device) UpdateChannels(list []*Channel) {
	for _, c := range list {
		c.device = d
		have := false
		for i, o := range d.Channels {
			if o.DeviceID == c.DeviceID {
				c.ChannelEx = o.ChannelEx
				d.Channels[i] = c
				have = true
				break
			}
		}
		if !have {
			d.Channels = append(d.Channels, c)
		}
	}
}
func (d *Device) UpdateRecord(channelId string, list []*Record) {
	for _, c := range d.Channels {
		if c.DeviceID == channelId {
			c.Records = list
			//for _, o := range list {
			//	o.channel = c
			//}
			break
		}
	}
}
func (c *Channel) CreateMessage(Method sip.Method) (requestMsg *sip.Message) {
	requestMsg = c.device.CreateMessage(Method)
	requestMsg.StartLine.Uri = sip.NewURI(c.DeviceID + "@" + c.device.to.Uri.Domain())
	requestMsg.To = &sip.Contact{
		Uri: requestMsg.StartLine.Uri,
	}
	requestMsg.From = &sip.Contact{
		Uri:    sip.NewURI(config.Serial + "@" + config.Realm),
		Params: map[string]string{"tag": utils.RandNumString(9)},
	}
	return
}
func (c *Channel) GetPublishStreamPath(start string) string {
	if start == "0" {
		return fmt.Sprintf("%s/%s", c.device.ID, c.DeviceID)
	}
	return fmt.Sprintf("%s/%s", c.DeviceID, start)
}

func parseStreamPath(streamPath string) (parentID string, deviceID string, start string, err error) {
	s := strings.Split(streamPath, "/")
	if len(s) < 2 {
		err = errors.New("bad stream path")
		return
	}
	// 国标历史回放时间戳
	// Invite点播时SDP中t=使用Unix时间戳
	// 录像查询时返回消息体中使用2006-01-02T15:04:05时间戳, 19位
	// 不等于国标定义的18位或者20位编码
	if len(s[1]) == len(config.Serial) {
		// 实时预览
		parentID = s[0]
		deviceID = s[1]
	} else {
		// 历史回放
		deviceID = s[0]
		start = s[1]
	}
	return
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
		}, From: d.from,
		To: d.to, CSeq: &sip.CSeq{
			ID:     uint32(d.sn),
			Method: Method,
		}, CallID: utils.RandNumString(10),
		Addr: d.Addr,
	}
	requestMsg.From.Params["tag"] = utils.RandNumString(9)
	return
}
func (d *Device) Query() int {
	requestMsg := d.CreateMessage(sip.MESSAGE)
	requestMsg.ContentType = "Application/MANSCDP+xml"
	requestMsg.Body = fmt.Sprintf(`<?xml version="1.0"?>
<Query>
<CmdType>Catalog</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
</Query>`, d.sn, requestMsg.To.Uri.UserInfo())
	requestMsg.ContentLength = len(requestMsg.Body)
	response := d.SendMessage(requestMsg)
	if response.Data != nil && response.Data.Via.Params["received"] != "" {
		d.SipIP = response.Data.Via.Params["received"]
	}
	return response.Code
}
func (d *Device) QueryRecord(channelIndex int, startTime, endTime string) int {
	channel := d.Channels[channelIndex]
	requestMsg := channel.CreateMessage(sip.MESSAGE)
	requestMsg.ContentType = "Application/MANSCDP+xml"
	requestMsg.Body = fmt.Sprintf(`<?xml version="1.0"?>
<Query>
<CmdType>RecordInfo</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
<StartTime>%s</StartTime>
<EndTime>%s</EndTime>
<Secrecy>0</Secrecy>
<Type>time</Type>
</Query>`, d.sn, requestMsg.To.Uri.UserInfo(), startTime, endTime)
	requestMsg.ContentLength = len(requestMsg.Body)
	return d.SendMessage(requestMsg).Code
}
func (d *Device) Control(channelIndex int, PTZCmd string) int {
	channel := d.Channels[channelIndex]
	requestMsg := channel.CreateMessage(sip.MESSAGE)
	requestMsg.ContentType = "Application/MANSCDP+xml"
	requestMsg.Body = fmt.Sprintf(`<?xml version="1.0"?>
<Control>
<CmdType>DeviceControl</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
<PTZCmd>%s</PTZCmd>
</Control>`, d.sn, requestMsg.To.Uri.UserInfo(), PTZCmd)
	requestMsg.ContentLength = len(requestMsg.Body)
	return d.SendMessage(requestMsg).Code
}

/*
f字段： f = v/编码格式/分辨率/帧率/码率类型/码率大小a/编码格式/码率大小/采样率
各项具体含义：
    v：后续参数为视频的参数；各参数间以 “/”分割；
编码格式：十进制整数字符串表示
1 –MPEG-4 2 –H.264 3 – SVAC 4 –3GP
    分辨率：十进制整数字符串表示
1 – QCIF 2 – CIF 3 – 4CIF 4 – D1 5 –720P 6 –1080P/I
帧率：十进制整数字符串表示 0～99
码率类型：十进制整数字符串表示
1 – 固定码率（CBR）     2 – 可变码率（VBR）
码率大小：十进制整数字符串表示 0～100000（如 1表示1kbps）
    a：后续参数为音频的参数；各参数间以 “/”分割；
编码格式：十进制整数字符串表示
1 – G.711    2 – G.723.1     3 – G.729      4 – G.722.1
码率大小：十进制整数字符串
音频编码码率： 1 — 5.3 kbps （注：G.723.1中使用）
   2 — 6.3 kbps （注：G.723.1中使用）
   3 — 8 kbps （注：G.729中使用）
   4 — 16 kbps （注：G.722.1中使用）
   5 — 24 kbps （注：G.722.1中使用）
   6 — 32 kbps （注：G.722.1中使用）
   7 — 48 kbps （注：G.722.1中使用）
   8 — 64 kbps（注：G.711中使用）
采样率：十进制整数字符串表示
	1 — 8 kHz（注：G.711/ G.723.1/ G.729中使用）
	2—14 kHz（注：G.722.1中使用）
	3—16 kHz（注：G.722.1中使用）
	4—32 kHz（注：G.722.1中使用）
	注1：字符串说明
本节中使用的“十进制整数字符串”的含义为“0”～“4294967296” 之间的十进制数字字符串。
注2：参数分割标识
各参数间以“/”分割，参数间的分割符“/”不能省略；
若两个分割符 “/”间的某参数为空时（即两个分割符 “/”直接将相连时）表示无该参数值；
注3：f字段说明
使用f字段时，应保证视频和音频参数的结构完整性，即在任何时候，f字段的结构都应是完整的结构：
f = v/编码格式/分辨率/帧率/码率类型/码率大小a/编码格式/码率大小/采样率
若只有视频时，音频中的各参数项可以不填写，但应保持 “a///”的结构:
f = v/编码格式/分辨率/帧率/码率类型/码率大小a///
若只有音频时也类似处理，视频中的各参数项可以不填写，但应保持 “v/”的结构：
f = v/a/编码格式/码率大小/采样率
f字段中视、音频参数段之间不需空格分割。
可使用f字段中的分辨率参数标识同一设备不同分辨率的码流。
*/
func (d *Device) Invite(channelIndex int, start, end string, f string) int {
	channel := d.Channels[channelIndex]
	port, publisher := d.publish(channel.GetPublishStreamPath(start))
	if port == 0 {
		channel.Connected = true
		return 304
	}
	ssrc := "0200000001"
	// size := 1
	// fps := 15
	// bitrate := 200
	// fmt.Sprintf("f=v/2/%d/%d/1/%da///", size, fps, bitrate)
	s := "Play"
	if start != "0" {
		s = "Playback"
		publisher.AutoUnPublish = true
		channel.RecordSP = publisher.StreamPath
	} else {
		channel.LiveSP = publisher.StreamPath
	}
	sdpInfo := []string{
		"v=0",
		fmt.Sprintf("o=%s 0 0 IN IP4 %s", d.Serial, d.SipIP),
		"s=" + s,
		"u=" + channel.DeviceID + ":0",
		"c=IN IP4 " + d.SipIP,
		fmt.Sprintf("t=%s %s", start, end),
		fmt.Sprintf("m=video %d RTP/AVP 96 97 98", port),
		"a=recvonly",
		"a=rtpmap:96 PS/90000",
		"a=rtpmap:97 MPEG4/90000",
		"a=rtpmap:98 H264/90000",
		"y=" + ssrc,
		"f=" + f,
	}

	invite := channel.CreateMessage(sip.INVITE)
	invite.ContentType = "application/sdp"
	invite.Contact = &sip.Contact{
		Uri: sip.NewURI(fmt.Sprintf("%s@%s:%d", d.Serial, d.SipIP, d.SipPort)),
	}
	invite.Body = strings.Join(sdpInfo, "\r\n") + "\r\n"
	invite.ContentLength = len(invite.Body)
	invite.Subject = fmt.Sprintf("%s:%s,%s:0", channel.DeviceID, ssrc, config.Serial)
	response := d.SendMessage(invite)
	fmt.Printf("invite response statuscode: %d\n", response.Code)
	if response.Code == 200 {
		if start == "0" {
			channel.inviteRes = response.Data
			channel.Connected = true
		} else {
			channel.recordInviteRes = response.Data
		}
		ack := d.CreateMessage(sip.ACK)
		ack.StartLine = &sip.StartLine{
			Uri:    sip.NewURI(channel.DeviceID + "@" + d.to.Uri.Domain()),
			Method: sip.ACK,
		}
		ack.From = response.Data.From
		ack.To = response.Data.To
		ack.CallID = response.Data.CallID
		ack.CSeq.ID = invite.CSeq.ID
		go d.Send(ack)
	}
	return response.Code
}
func (d *Device) Bye(channelIndex int) int {
	channel := d.Channels[channelIndex]
	defer func() {
		channel.inviteRes = nil
		channel.Connected = false
	}()
	return channel.Bye(channel.inviteRes).Code
}
func (c *Channel) Bye(res *sip.Message) *transaction.Response {
	if res == nil {
		return nil
	}
	bye := c.device.CreateMessage(sip.BYE)
	bye.StartLine = &sip.StartLine{
		Uri:    sip.NewURI(c.DeviceID + "@" + c.device.to.Uri.Domain()),
		Method: sip.BYE,
	}
	bye.From = res.From
	bye.To = res.To
	bye.CallID = res.CallID
	return c.device.SendMessage(bye)
}
