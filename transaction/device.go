package transaction

import (
	"fmt"
	"github.com/Monibuca/plugin-gb28181/sip"
	"github.com/Monibuca/plugin-gb28181/utils"
	"time"
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
	host         string
	port         string
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
			Host:      d.host,
			Port:      d.port,
			Params: map[string]string{
				"branch": fmt.Sprintf("z9hG4bK%s", utils.RandNumString(8)),
				"rport":  "-1", //only key,no-value
			},
		}, From: d.from,
		To: d.to, CSeq: &sip.CSeq{
			ID:     1,
			Method: Method,
		}, CallID: utils.RandNumString(8),
	}
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
func (d *Device) Control(PTZCmd string) int {
	requestMsg := d.CreateMessage(sip.MESSAGE)
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
func (d *Device) Invite(channelIndex int) {
	channel := &d.Channels[channelIndex]
	port := d.core.OnInvite(channel)
	if port == 0 {
		return
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
a=rtpmap:98 H264/90000`, d.core.config.Serial, d.core.config.MediaIP, d.core.config.MediaIP, port)
	invite := d.CreateMessage(sip.INVITE)
	invite.StartLine.Uri = sip.NewURI(channel.DeviceID + "@" + d.to.Uri.Domain())
	invite.To = &sip.Contact{
		Uri: invite.StartLine.Uri,
	}
	invite.ContentType = "application/sdp"
	invite.Contact = &sip.Contact{
		Nickname: d.core.config.Serial,
		Uri:      sip.NewURI(fmt.Sprintf("%s:%d", d.core.config.MediaIP, d.core.config.MediaPort)),
	}
	invite.Body = sdp
	invite.ContentLength = len(sdp)
	fmt.Printf("invite response statuscode: %d\n", d.core.SendMessage(invite).Code)
}
