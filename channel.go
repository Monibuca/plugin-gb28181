package gb28181

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Monibuca/engine/v3"
	. "github.com/Monibuca/plugin-gb28181/v3/sip"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
	. "github.com/Monibuca/utils/v3"
)

type ChannelEx struct {
	device          *Device   `json:"-"`
	inviteRes       *Response `json:"-"`
	recordInviteRes *Response `json:"-"`
	RecordPublisher *Publisher
	LivePublisher   *Publisher
	LiveSubSP       string //实时子码流
	Records         []*Record
	RecordStartTime string
	RecordEndTime   string
	recordStartTime time.Time
	recordEndTime   time.Time
	state           int32
	tcpPortIndex    uint16
}

// Channel 通道
type Channel struct {
	DeviceID     string
	ParentID     string
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
	Children     []*Channel `json:"-"`
	*ChannelEx              //自定义属性
}

func (c *Channel) CreateRequst(Method Method) (request *Request) {
	request = &Request{}
	request.Message = c.device.CreateMessage(Method)
	request.Message.StartLine.Uri = NewURI(c.DeviceID + "@" + c.device.to.Uri.Domain())
	request.Message.To = &Contact{
		Uri: request.Message.StartLine.Uri,
	}
	request.Message.From = &Contact{
		Uri:    NewURI(config.Serial + "@" + config.Realm),
		Params: map[string]string{"tag": utils.RandNumString(9)},
	}
	return
}
func (channel *Channel) QueryRecord(startTime, endTime string) int {
	d := channel.device
	channel.RecordStartTime = startTime
	channel.RecordEndTime = endTime
	channel.recordStartTime, _ = time.Parse(TIME_LAYOUT, startTime)
	channel.recordEndTime, _ = time.Parse(TIME_LAYOUT, endTime)
	channel.Records = nil
	requestMsg := channel.CreateRequst(MESSAGE)
	requestMsg.ContentType = "Application/MANSCDP+xml"
	requestMsg.Body = fmt.Sprintf(`<?xml version="1.0"?>
<Query>
<CmdType>RecordInfo</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
<StartTime>%s</StartTime>
<EndTime>%s</EndTime>
<Secrecy>0</Secrecy>
<Type>all</Type>
</Query>`, d.sn, requestMsg.To.Uri.UserInfo(), startTime, endTime)
	requestMsg.ContentLength = len(requestMsg.Body)
	resp, err := d.SipRequestForResponse(requestMsg)
	if err != nil {
		return http.StatusRequestTimeout
	}
	return resp.GetStatusCode()

}
func (channel *Channel) Control(PTZCmd string) int {
	d := channel.device
	requestMsg := channel.CreateRequst(MESSAGE)
	requestMsg.ContentType = "Application/MANSCDP+xml"
	requestMsg.Body = fmt.Sprintf(`<?xml version="1.0"?>
<Control>
<CmdType>DeviceControl</CmdType>
<SN>%d</SN>
<DeviceID>%s</DeviceID>
<PTZCmd>%s</PTZCmd>
</Control>`, d.sn, requestMsg.To.Uri.UserInfo(), PTZCmd)
	requestMsg.ContentLength = len(requestMsg.Body)
	resp, err := d.SipRequestForResponse(requestMsg)
	if err != nil {
		return http.StatusRequestTimeout
	}
	return resp.GetStatusCode()
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
func (channel *Channel) Invite(start, end string) (code int) {
	if start == "" {
		if !atomic.CompareAndSwapInt32(&channel.state, 0, 1) {
			return 304
		}
		defer func() {
			if code != 200 {
				atomic.StoreInt32(&channel.state, 0)
			}
		}()
		channel.Bye(true)
	} else {
		channel.Bye(false)
	}
	sint, err1 := strconv.ParseInt(start, 10, 0)
	eint, err2 := strconv.ParseInt(end, 10, 0)
	d := channel.device
	streamPath := fmt.Sprintf("%s/%s", d.ID, channel.DeviceID)
	s := "Play"
	ssrc := make([]byte, 10)
	if start != "" {
		if err1 != nil || err2 != nil {
			return 400
		}
		s = "Playback"
		ssrc[0] = '1'
		streamPath = fmt.Sprintf("%s/%s/%s-%s", d.ID, channel.DeviceID, start, end)
	} else {
		ssrc[0] = '0'
	}

	// size := 1
	// fps := 15
	// bitrate := 200
	// fmt.Sprintf("f=v/2/%d/%d/1/%da///", size, fps, bitrate)
	copy(ssrc[1:6], []byte(config.Serial[3:8]))
	randNum := rand.Intn(10000)
	copy(ssrc[6:], []byte(strconv.Itoa(randNum)))
	protocol := ""
	port := config.MediaPort
	if config.TCP {
		protocol = "TCP/"
		port = config.MediaPort + channel.tcpPortIndex
		if channel.tcpPortIndex++; channel.tcpPortIndex >= config.TCPMediaPortNum {
			channel.tcpPortIndex = 0
		}
	}
	sdpInfo := []string{
		"v=0",
		fmt.Sprintf("o=%s 0 0 IN IP4 %s", d.Serial, d.SipIP),
		"s=" + s,
		"u=" + channel.DeviceID + ":0",
		"c=IN IP4 " + d.SipIP,
		fmt.Sprintf("t=%d %d", sint, eint),
		fmt.Sprintf("m=video %d %sRTP/AVP 96", port, protocol),
		"a=recvonly",
		"a=rtpmap:96 PS/90000",
	}
	if config.TCP {
		sdpInfo = append(sdpInfo, "a=setup:passive", "a=connection:new")
	}
	invite := channel.CreateRequst(INVITE)
	invite.ContentType = "application/sdp"
	invite.Contact = &Contact{
		Uri: NewURI(fmt.Sprintf("%s@%s:%d", d.Serial, d.SipIP, d.SipPort)),
	}
	invite.Body = strings.Join(sdpInfo, "\r\n") + "\r\ny=" + string(ssrc) + "\r\n"
	invite.ContentLength = len(invite.Body)
	invite.Subject = fmt.Sprintf("%s:%s,%s:0", channel.DeviceID, ssrc, config.Serial)
	response, _ := d.Core.SipRequestForResponse(invite)
	if response == nil {
		return http.StatusRequestTimeout
	}
	Printf("Channel :%s invite response status code: %d\n", channel.DeviceID, response.GetStatusCode())

	if response.GetStatusCode() == 200 {
		ds := strings.Split(response.Body, "\r\n")
		_SSRC, _ := strconv.ParseInt(string(ssrc), 10, 0)
		SSRC := uint32(_SSRC)
		for _, l := range ds {
			if ls := strings.Split(l, "="); len(ls) > 1 {
				if ls[0] == "y" && len(ls[1]) > 0 {
					_SSRC, _ = strconv.ParseInt(ls[1], 10, 0)
					SSRC = uint32(_SSRC)
					break
				}
			}
		}
		publisher := &Publisher{
			Stream: &engine.Stream{
				StreamPath:     streamPath,
				AutoCloseAfter: &config.AutoCloseAfter,
			},
		}
		if config.UdpCacheSize > 0 && !config.TCP {
			publisher.udpCache = utils.NewPqRtp()
		}
		if start == "" {
			publisher.Type = "GB28181 Live"
			publisher.OnClose = func() {
				publishers.Remove(SSRC)
				channel.LivePublisher = nil
				channel.ByeBye((*Request)(channel.inviteRes))
				channel.inviteRes = nil
				atomic.StoreInt32(&channel.state, 0)
				if config.AutoInvite {
					go channel.Invite("", "")
				}
			}
		} else {
			publisher.Type = "GB28181 Record"
			publisher.OnClose = func() {
				publishers.Remove(SSRC)
				channel.RecordPublisher = nil
				channel.ByeBye((*Request)(channel.recordInviteRes))
				channel.recordInviteRes = nil
			}
		}
		if !publisher.Publish() {
			return 403
		}
		publishers.Add(SSRC, publisher)
		if start == "" {
			channel.inviteRes = response
			channel.LivePublisher = publisher
		} else {
			channel.RecordPublisher = publisher
			channel.recordInviteRes = response
		}
		ack := d.CreateMessage(ACK)
		ack.StartLine = &StartLine{
			Uri:    NewURI(channel.DeviceID + "@" + d.to.Uri.Domain()),
			Method: ACK,
		}
		ack.From = response.From
		ack.To = response.To
		ack.CallID = response.CallID
		ack.CSeq.ID = invite.CSeq.ID
		d.Respond(&Response{Message: ack})
	} else if start == "" && config.AutoInvite {
		time.AfterFunc(time.Second*5, func() {
			channel.Invite("", "")
		})
	}
	return response.GetStatusCode()
}
func (channel *Channel) Bye(live bool) int {
	if live && channel.inviteRes != nil {
		defer func() {
			channel.inviteRes = nil
			if channel.LivePublisher != nil {
				channel.LivePublisher.Close()
			}
		}()
		return channel.ByeBye((*Request)(channel.inviteRes)).GetStatusCode()
	}
	if !live && channel.recordInviteRes != nil {
		defer func() {
			channel.recordInviteRes = nil
			if channel.RecordPublisher != nil {
				channel.RecordPublisher.Close()
			}
		}()
		return channel.ByeBye((*Request)(channel.recordInviteRes)).GetStatusCode()
	}
	return 404
}
func (c *Channel) ByeBye(res *Request) *Response {
	if res == nil {
		return nil
	}
	d := c.device
	bye := c.device.CreateMessage(BYE)
	bye.StartLine = &StartLine{
		Uri:    NewURI(c.DeviceID + "@" + c.device.to.Uri.Domain()),
		Method: BYE,
	}
	bye.From = res.From
	bye.To = res.To
	bye.CallID = res.CallID
	req := &Request{}
	req.Message = bye

	resp, _ := d.SipRequestForResponse(req)
	return resp

}
