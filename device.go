package gb28181

import (
	"encoding/json"
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

func (c *Channel) MarshalJSON() ([]byte, error) {
	var data = map[string]interface{}{
		"DeviceID":     c.DeviceID,
		"Name":         c.Name,
		"Manufacturer": c.Manufacturer,
		"Address":      c.Address,
		"Status":       c.Status,
		"RecordSP":     c.RecordSP,
		"LiveSP":       c.LiveSP,
		"Records":      c.Records,
		"Connected":    c.Connected,
	}
	return json.Marshal(data)
}

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
func (d *Device) Invite(channelIndex int, start, end string) int {
	channel := d.Channels[channelIndex]
	port, publisher := d.publish(channel.GetPublishStreamPath(start))
	if port == 0 {
		channel.Connected = true
		return 304
	}
	ssrc := "0200000001"
	sdpInfo := []string{"v=0", fmt.Sprintf("o=%s 0 0 IN IP4 %s", d.Serial, d.SipIP), "s=Play", "c=IN IP4 " + d.SipIP, fmt.Sprintf("t=%s %s", start, end), fmt.Sprintf("m=video %d RTP/AVP 96 98 97", port), "a=recvonly", "a=rtpmap:96 PS/90000", "a=rtpmap:97 MPEG4/90000", "a=rtpmap:98 H264/90000", "y=" + ssrc}
	if start != "0" {
		sdpInfo[2] = "s=Playback"
		publisher.AutoUnPublish = true
		channel.RecordSP = publisher.StreamPath
	} else {
		channel.LiveSP = publisher.StreamPath
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
