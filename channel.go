package gb28181

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sync/atomic"

	"github.com/ghettovoice/gosip/sip"
	"go.uber.org/zap"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/log"
	"m7s.live/plugin/gb28181/v4/utils"
	"m7s.live/plugin/ps/v4"
)

type PullStream struct {
	opt       *InviteOptions
	channel   *Channel
	inviteRes sip.Response
}

func (p *PullStream) Bye() int {
	res := p.inviteRes
	bye := p.channel.CreateRequst(sip.BYE)
	from, _ := res.From()
	to, _ := res.To()
	callId, _ := res.CallID()
	bye.ReplaceHeaders(from.Name(), []sip.Header{from})
	bye.ReplaceHeaders(to.Name(), []sip.Header{to})
	bye.ReplaceHeaders(callId.Name(), []sip.Header{callId})
	resp, err := p.channel.device.SipRequestForResponse(bye)
	if p.opt.IsLive() {
		p.channel.status.Store(0)
		// defer p.channel.TryAutoInvite(p.opt)
	}
	if p.opt.recyclePort != nil {
		p.opt.recyclePort(p.opt.MediaPort)
	}
	if err != nil {
		return ServerInternalError
	}
	return int(resp.StatusCode())
}

type ChannelEx struct {
	device          *Device      // 所属设备
	status          atomic.Int32 // 通道状态,0:空闲,1:正在invite,2:正在播放
	LiveSubSP       string       // 实时子码流，通过rtsp
	Records         []*Record
	RecordStartTime string
	RecordEndTime   string
	recordStartTime time.Time
	recordEndTime   time.Time
	GpsTime         time.Time //gps时间
	Longitude       string    //经度
	Latitude        string    //纬度
	*log.Logger     `json:"-" yaml:"-"`
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
	Port         int
	Parental     int
	SafetyWay    int
	RegisterWay  int
	Secrecy      int
	Status       string
	Children     []*Channel `json:"-" yaml:"-"`
	ChannelEx               //自定义属性
}

func (channel *Channel) CreateRequst(Method sip.RequestMethod) (req sip.Request) {
	d := channel.device
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
		//DisplayName: sip.String{Str: d.serverConfig.Serial},
		Uri: &sip.SipUri{
			FUser: sip.String{Str: conf.Serial},
			FHost: d.sipIP,
			FPort: &port,
		},
		Params: sip.NewParams().Add("tag", sip.String{Str: utils.RandNumString(9)}),
	}
	//非同一域的目标地址需要使用@host
	host := conf.Realm
	if channel.DeviceID[0:9] != host {
		if channel.Port != 0 {
			deviceIp := d.NetAddr
			deviceIp = deviceIp[0:strings.LastIndex(deviceIp, ":")]
			host = fmt.Sprintf("%s:%d", deviceIp, channel.Port)
		} else {
			host = d.NetAddr
		}
	}

	channelAddr := sip.Address{
		//DisplayName: sip.String{Str: d.serverConfig.Serial},
		Uri: &sip.SipUri{FUser: sip.String{Str: channel.DeviceID}, FHost: host},
	}
	req = sip.NewRequest(
		"",
		Method,
		channelAddr.Uri,
		"SIP/2.0",
		[]sip.Header{
			serverAddr.AsFromHeader(),
			channelAddr.AsToHeader(),
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
	return req
}
func (channel *Channel) QueryRecord(startTime, endTime string) int {
	d := channel.device
	channel.RecordStartTime = startTime
	channel.RecordEndTime = endTime
	channel.recordStartTime, _ = time.Parse(TIME_LAYOUT, startTime)
	channel.recordEndTime, _ = time.Parse(TIME_LAYOUT, endTime)
	channel.Records = nil
	request := d.CreateRequest(sip.MESSAGE)
	contentType := sip.ContentType("Application/MANSCDP+xml")
	request.AppendHeader(&contentType)
	body := fmt.Sprintf(`<?xml version="1.0"?>
		<Query>
		<CmdType>RecordInfo</CmdType>
		<SN>%d</SN>
		<DeviceID>%s</DeviceID>
		<StartTime>%s</StartTime>
		<EndTime>%s</EndTime>
		<Secrecy>0</Secrecy>
		<Type>all</Type>
		</Query>`, d.sn, channel.DeviceID, startTime, endTime)
	request.SetBody(body, true)
	resp, err := d.SipRequestForResponse(request)
	if err != nil {
		return http.StatusRequestTimeout
	}
	return int(resp.StatusCode())
}
func (channel *Channel) Control(PTZCmd string) int {
	d := channel.device
	request := d.CreateRequest(sip.MESSAGE)
	contentType := sip.ContentType("Application/MANSCDP+xml")
	request.AppendHeader(&contentType)
	body := fmt.Sprintf(`<?xml version="1.0"?>
		<Control>
		<CmdType>DeviceControl</CmdType>
		<SN>%d</SN>
		<DeviceID>%s</DeviceID>
		<PTZCmd>%s</PTZCmd>
		</Control>`, d.sn, channel.DeviceID, PTZCmd)
	request.SetBody(body, true)
	resp, err := d.SipRequestForResponse(request)
	if err != nil {
		return http.StatusRequestTimeout
	}
	return int(resp.StatusCode())
}

// Invite 发送Invite报文 invites a channel to play
// 注意里面的锁保证不同时发送invite报文，该锁由channel持有
/***
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

func (channel *Channel) Invite(opt *InviteOptions) (code int, err error) {
	if opt.IsLive() {
		if !channel.status.CompareAndSwap(0, 1) {
			return 304, nil
		}
		defer func() {
			if err != nil {
				GB28181Plugin.Error("Invite", zap.Error(err))
				channel.status.Store(0)
				if conf.InviteMode == 1 {
					// 5秒后重试
					time.AfterFunc(time.Second*5, func() {
						channel.Invite(opt)
					})
				}
			} else {
				channel.status.Store(2)
			}
		}()
	}
	d := channel.device
	streamPath := fmt.Sprintf("%s/%s", d.ID, channel.DeviceID)
	s := "Play"
	opt.CreateSSRC()
	if opt.Record() {
		s = "Playback"
		streamPath = fmt.Sprintf("%s/%s/%d-%d", d.ID, channel.DeviceID, opt.Start, opt.End)
	}
	if opt.StreamPath != "" {
		streamPath = opt.StreamPath
	}
	if opt.dump == "" {
		opt.dump = conf.DumpPath
	}
	protocol := ""
	networkType := "udp"
	resuePort := true
	if conf.IsMediaNetworkTCP() {
		networkType = "tcp"
		protocol = "TCP/"
		if conf.tcpPorts.Valid {
			opt.MediaPort, err = conf.tcpPorts.GetPort()
			opt.recyclePort = conf.tcpPorts.Recycle
			resuePort = false
		}
	} else {
		if conf.udpPorts.Valid {
			opt.MediaPort, err = conf.udpPorts.GetPort()
			opt.recyclePort = conf.udpPorts.Recycle
			resuePort = false
		}
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if opt.MediaPort == 0 {
		opt.MediaPort = conf.MediaPort
	}

	sdpInfo := []string{
		"v=0",
		fmt.Sprintf("o=%s 0 0 IN IP4 %s", channel.DeviceID, d.mediaIP),
		"s=" + s,
		"u=" + channel.DeviceID + ":0",
		"c=IN IP4 " + d.mediaIP,
		opt.String(),
		fmt.Sprintf("m=video %d %sRTP/AVP 96", opt.MediaPort, protocol),
		"a=recvonly",
		"a=rtpmap:96 PS/90000",
		"y=" + opt.ssrc,
	}
	if conf.IsMediaNetworkTCP() {
		sdpInfo = append(sdpInfo, "a=setup:passive", "a=connection:new")
	}
	invite := channel.CreateRequst(sip.INVITE)
	contentType := sip.ContentType("application/sdp")
	invite.AppendHeader(&contentType)

	invite.SetBody(strings.Join(sdpInfo, "\r\n")+"\r\n", true)

	subject := sip.GenericHeader{
		HeaderName: "Subject", Contents: fmt.Sprintf("%s:%s,%s:0", channel.DeviceID, opt.ssrc, conf.Serial),
	}
	invite.AppendHeader(&subject)
	inviteRes, err := d.SipRequestForResponse(invite)
	if err != nil {
		channel.Error("invite", zap.Error(err), zap.String("msg", invite.String()))
		return http.StatusInternalServerError, err
	}
	code = int(inviteRes.StatusCode())
	channel.Info("invite response", zap.Int("status code", code))

	if code == OK {
		ds := strings.Split(inviteRes.Body(), "\r\n")
		for _, l := range ds {
			if ls := strings.Split(l, "="); len(ls) > 1 {
				if ls[0] == "y" && len(ls[1]) > 0 {
					if _ssrc, err := strconv.ParseInt(ls[1], 10, 0); err == nil {
						opt.SSRC = uint32(_ssrc)
					} else {
						channel.Error("read invite response y ", zap.Error(err))
					}
					break
				}
			}
		}
		err = ps.Receive(streamPath, opt.dump, fmt.Sprintf("%s:%d", networkType, opt.MediaPort), opt.SSRC, resuePort)
		if err == nil {
			PullStreams.Store(streamPath, &PullStream{
				opt:       opt,
				channel:   channel,
				inviteRes: inviteRes,
			})
			err = srv.Send(sip.NewAckRequest("", invite, inviteRes, "", nil))
		}
	}
	return
}

func (channel *Channel) Bye(streamPath string) int {
	d := channel.device
	if streamPath == "" {
		streamPath = fmt.Sprintf("%s/%s", d.ID, channel.DeviceID)
	}
	if s, loaded := PullStreams.LoadAndDelete(streamPath); loaded {
		s.(*PullStream).Bye()
		if s := Streams.Get(streamPath); s != nil {
			s.Close()
		}
		return http.StatusOK
	}
	return http.StatusNotFound
}

func (channel *Channel) TryAutoInvite(opt *InviteOptions) {
	if conf.InviteMode == 1 && channel.CanInvite() {
		go channel.Invite(opt)
	}
}

func (channel *Channel) CanInvite() bool {
	if channel.status.Load() != 0 || len(channel.DeviceID) != 20 || channel.Status == "OFF" {
		return false
	}

	if conf.InviteIDs == "" {
		return true
	}

	// 11～13位是设备类型编码
	typeID := channel.DeviceID[10:13]

	// format: start-end,type1,type2
	tokens := strings.Split(conf.InviteIDs, ",")
	for _, tok := range tokens {
		if first, second, ok := strings.Cut(tok, "-"); ok {
			if typeID >= first && typeID <= second {
				return true
			}
		} else {
			if typeID == first {
				return true
			}
		}
	}

	return false
}
