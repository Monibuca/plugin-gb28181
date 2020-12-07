package transaction

import (
	"fmt"
	"strings"
	"time"

	"github.com/Monibuca/plugin-gb28181/sip"
	"github.com/Monibuca/plugin-gb28181/utils"
)

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
	device       *Device
	inviteRes    *sip.Message
	Connected    bool
}
type Device struct {
	ID           string
	RegisterTime time.Time
	UpdateTime   time.Time
	Status       string
	Channels     []Channel
	core         *Core
	sn           int
	from         *sip.Contact
	to           *sip.Contact
	Addr         string
}

func (c *Core) RemoveDead() {
	c.Devices.Range(func(k, v interface{}) bool {
		device := v.(*Device)
		if device.UpdateTime.Sub(device.RegisterTime) > time.Duration(c.config.RegisterValidity)*time.Second {
			c.Devices.Delete(k)
		}
		return true
	})
}
func (d *Device) UpdateChannels(list []Channel) {
	for _, c := range list {
		c.device = d
		have := false
		for i, o := range d.Channels {
			if o.DeviceID == c.DeviceID {
				c.inviteRes = o.inviteRes
				c.Connected = o.inviteRes != nil
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
func (c *Channel) CreateMessage(Method sip.Method) (requestMsg *sip.Message) {
	requestMsg = c.device.CreateMessage(Method)
	requestMsg.StartLine.Uri = sip.NewURI(c.DeviceID + "@" + c.device.to.Uri.Domain())
	requestMsg.To = &sip.Contact{
		Uri: requestMsg.StartLine.Uri,
	}
	requestMsg.From = &sip.Contact{
		Uri:    sip.NewURI(c.device.core.config.Serial + "@" + c.device.core.config.Realm),
		Params: map[string]string{"tag": utils.RandNumString(9)},
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
			Host:      d.core.config.SipIP,
			Port:      fmt.Sprintf("%d", d.core.config.SipPort),
			Params: map[string]string{
				"branch": fmt.Sprintf("z9hG4bK%s", utils.RandNumString(8)),
				"rport":  "-1", //only key,no-value
			},
		}, From: d.from,
		To: d.to, CSeq: &sip.CSeq{
			ID:     1,
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
	return d.core.SendMessage(requestMsg).Code
}
func (d *Device) Control(channelIndex int, PTZCmd string) int {
	channel := &d.Channels[channelIndex]
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
	return d.core.SendMessage(requestMsg).Code
}
func (d *Device) Invite(channelIndex int) int {
	channel := &d.Channels[channelIndex]
	port := d.core.OnInvite(channel)
	if port == 0 {
		channel.Connected = true
		return 304
	}
	sdp := fmt.Sprintf(`v=0
o=%s 0 0 IN IP4 %s
s=Play
c=IN IP4 %s
t=0 0
m=video %d RTP/AVP 96 98 97
a=recvonly
a=rtpmap:96 PS/90000
a=rtpmap:97 MPEG4/90000
a=rtpmap:98 H264/90000
y=0200000001
`, d.core.config.Serial, d.core.config.MediaIP, d.core.config.MediaIP, port)
	sdp = strings.ReplaceAll(sdp, "\n", "\r\n")
	invite := channel.CreateMessage(sip.INVITE)
	invite.ContentType = "application/sdp"
	invite.Contact = &sip.Contact{
		Uri: sip.NewURI(fmt.Sprintf("%s@%s:%d", d.core.config.Serial, d.core.config.SipIP, d.core.config.SipPort)),
	}
	invite.Body = sdp
	invite.ContentLength = len(sdp)
	invite.Subject = fmt.Sprintf("%s:0200000001,34020000002020000001:0", channel.DeviceID)
	response := d.core.SendMessage(invite)
	fmt.Printf("invite response statuscode: %d\n", response.Code)
	if response.Code == 200 {
		channel.inviteRes = response.Data
		channel.Connected = true
		channel.Ack()
	}
	return response.Code
}
func (d *Device) Bye(channelIndex int) int {
	channel := &d.Channels[channelIndex]
	defer func() {
		channel.inviteRes = nil
		channel.Connected = false
	}()
	return channel.Bye().Code
}
func (c *Channel) Ack() {
	ack := c.CreateMessage(sip.ACK)
	ack.From = c.inviteRes.From
	ack.To = c.inviteRes.To
	ack.CallID = c.inviteRes.CallID
	go c.device.core.Send(ack)
}
func (c *Channel) Bye() *Response {
	bye := c.CreateMessage(sip.BYE)
	bye.From = c.inviteRes.From
	bye.To = c.inviteRes.To
	bye.CallID = c.inviteRes.CallID
	return c.device.core.SendMessage(bye)
}
