package gb28181

import (
	"fmt"
	"sync"
	"time"

	"github.com/Monibuca/plugin-gb28181/v3/sip"
	"github.com/Monibuca/plugin-gb28181/v3/transaction"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
	// . "github.com/Monibuca/utils/v3"
	// . "github.com/logrusorgru/aurora"
)

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
	channelMap        map[string]*Channel
	channelMutex      sync.RWMutex
}

func (d *Device) UpdateChannels(list []*Channel) {
	d.channelMutex.Lock()
	defer d.channelMutex.Unlock()
	for _, c := range list {
		c.device = d
		var oldList []*Channel
		if c.ParentID != "" {
			if parent, ok := d.channelMap[c.ParentID]; ok {
				oldList = parent.Children
				parent.Children = list
			}
		} else {
			oldList = d.Channels
			d.Channels = list
		}
		if len(oldList) > 0 {
			for _, o := range oldList {
				if o.DeviceID == c.DeviceID {
					c.ChannelEx = o.ChannelEx
					break
				}
			}
		}
		d.channelMap[c.DeviceID] = c
	}

	//单通道代码
	// inviteFunc := func() {
	// 	Print(Green("total count of channels is"), BrightBlue(len(d.Channels)))
	// 	if config.AutoInvite {
	// 		clen := len(d.Channels)
	// 		for i := 0; i < clen; i++ {
	// 			resultCode := d.Invite(i, "", "")
	// 			Print(Green("invite result is"), resultCode, Green("current index is "), i)
	// 			if resultCode == 200 {
	// 				break
	// 			}
	// 		}
	// 	}
	// }
	// go once.Do(inviteFunc)
	//firstChannel.channelMutex.Lock()
	//if firstChannel.firstChannel {
	//	firstChannel.firstChannel = false
	//	firstChannel.channelMutex.Unlock()
	//	if len(d.Channels) > 0 {
	//		go d.Invite(0, "", "")
	//	}
	//} else {
	//	firstChannel.channelMutex.Unlock()
	//}
	//多通道代码
	//for i := range d.Channels {
	//	if config.AutoInvite {
	//		go d.Invite(i, "", "")
	//	}
	//}
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

