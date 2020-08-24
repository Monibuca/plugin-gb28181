package gb28181

import (
	"fmt"
	"gitee.com/xiaochengtech/sdp"
	"gitee.com/xiaochengtech/sip"
	. "github.com/Monibuca/engine/v2"
	pluginrtp "github.com/Monibuca/plugin-rtp"
	"math/rand"

	//"github.com/max-min/streams/packet"
	"github.com/looplab/fsm"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var Devices sync.Map
var config = struct {
	SipID      string
	ListenAddr string
	Expires    int
	AutoInvite bool
	MediaPort int
}{"34020000002000000001", ":5060", 3600, true,6000}

type DeviceInfo struct {
	sip.User
	Self sip.User
	Status     string
	LastUpdate time.Time
}

type DevicePublisher struct {
	pluginrtp.RTP
	*fsm.FSM
	conn *net.UDPConn
	addr *net.UDPAddr
	DeviceInfo
}

func NewDevicePublisher(conn *net.UDPConn, header sip.Header) {
	var result DevicePublisher
	result.addr,_=net.ResolveUDPAddr("udp",header.From.URI.Domain)
	result.conn = conn
	result.User = header.From
	uri,_:=sip.NewURI(conn.LocalAddr().String())
	result.Self = sip.User{
		config.SipID,
		uri,
		sip.Args{},
	}
	Devices.Store(result.User.URI.String(), &result)
	result.FSM = fsm.NewFSM(sip.MethodRegister, []fsm.EventDesc{
		{"invite", []string{sip.MethodRegister}, sip.MethodInvite},
		{"ok", []string{sip.MethodInvite}, sip.MethodAck},
	}, map[string]fsm.Callback{
		"invite": result.invite,
		"ok":     result.onOK,
	})
	if config.AutoInvite {
		result.Event("invite")
	}
}
var  r = rand.New(rand.NewSource(time.Now().Unix()))
func GenerateNonce(length int) string{
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
func (p *DevicePublisher) invite(e *fsm.Event) {
	p.Status = e.Dst
	var invite sip.Message
	invite.IsRequest = true
	invite.Header = sip.Header{
		From:   p.Self,
		To:     p.User,
		CSeq:   sip.CSeq{1, sip.MethodInvite},
		CallID: p.Self.URI.String(),
		ContentType: "application/sdp",
	}
	invite.Header.Via.Add(fmt.Sprintf("SIP/2.0/UDP %s;branch=z9hG4bK%s,",p.Self.URI.Domain,GenerateNonce(29)))
	invite.Header.MaxForwards.Reset()
	invite.RequestLine, _ = sip.NewRequestLine(fmt.Sprintf("%s %s %s", sip.MethodInvite, p.Self.URI.String(), "SIP/2.0"))
	invite.Body = &sdp.Message{
		Version:     "0",
		Basic:       sdp.Basic{p.User.Username(), "0", "0", "IN", "IP4", p.Self.URI.Domain},
		SessionName: "Play",
		Connection:  &sdp.Connection{"IN", "IP4", p.Self.URI.Domain},
		Media: []sdp.Media{
			{Info: sdp.MediaInfo{"video", strconv.Itoa(config.MediaPort), "RTP/AVP", "96 98 97"}, UnsupportLine: []string{
				"a=recvonly", "a=rtpmap:96 PS/90000", "a=rtpmap:97 MPEG4/90000", "a=rtpmap:98 H264/90000",
			}},
		},
	}
	p.conn.WriteToUDP([]byte(invite.String()),p.addr)
}
func (p *DevicePublisher) onOK(e *fsm.Event) {
	p.Status = e.Dst
	if p.Publish(p.User.URI.String()) {
		var ack sip.Message
		ack.IsRequest = true
		ack.RequestLine, _ = sip.NewRequestLine(fmt.Sprintf("%s %s %s", sip.MethodAck, p.Self.URI.String(), "SIP/2.0"))
		ack.Header = sip.Header{
			From:  p.Self,
			To:     p.User,
			CSeq:   sip.CSeq{1, sip.MethodAck},
			CallID: p.Self.URI.String(),
		}
		ack.Header.Via.Add(fmt.Sprintf("SIP/2.0/UDP %s;branch=z9hG4bK%s,",p.Self.URI.Domain,GenerateNonce(29)))
		ack.Header.MaxForwards.Reset()
		p.conn.WriteToUDP([]byte(ack.String()),p.addr)
	} else {

	}
}
func init() {
	InstallPlugin(&PluginConfig{
		Name:   "GB28181",
		Config: &config,
		Type:   PLUGIN_PUBLISHER,
		Run:    run,
	})
}
func resolvePS(){
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d",config.MediaPort))
	if err != nil {
		log.Fatal(err)
	}
	if listener,err:=net.ListenTCP("tcp",addr);err==nil{
		for {
			l,err:=listener.AcceptTCP()
			if err!=nil{
				log.Fatal(err)
				return
			}
			//parser:=packet.NewRtpParsePacket()
			go func() {
				b:=make([]byte,1024)
				l.Read(b)
			}()
		}
	}else{
		log.Fatal(err)
	}
}

func run() {
	go resolvePS()
	addr, err := net.ResolveUDPAddr("udp", config.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		sipMsg, err := sip.NewMessage(listener)
		if err != nil {
			Print(err)
			continue
		}
		var target *DevicePublisher
		if sipMsg.IsRequest {
			targetUdp,_:=net.ResolveUDPAddr("udp",sipMsg.Header.From.URI.Domain)
			if d, ok := Devices.Load(sipMsg.Header.From.URI.String()); ok {
				target = d.(*DevicePublisher)
				target.LastUpdate = time.Now()
			}
			switch sipMsg.RequestLine.Method {
			case sip.MethodRegister:
				if sipMsg.Header.Authorization == "" {
					res := sip.NewResponse(sip.StatusUnauthorized, &sipMsg)
					h := &res.Header
					h.WWWAuthenticate = "Digest realm=\"3402000000\",nonce=\"1677f194104d46aea6c9f8aebe507017\""
					listener.WriteToUDP([]byte(res.String()),targetUdp)
				} else {
					res := sip.NewResponse(sip.StatusOK, &sipMsg)
					listener.WriteToUDP([]byte(res.String()),targetUdp)
					NewDevicePublisher(listener, sipMsg.Header)
				}
			default:
				listener.WriteToUDP([]byte(sip.NewResponse(sip.StatusOK, &sipMsg).String()),targetUdp)
				if target == nil{
					NewDevicePublisher(listener, sipMsg.Header)
				}
			}
		} else {
			if d, ok := Devices.Load(sipMsg.Header.To.URI.String()); ok {
				target = d.(*DevicePublisher)
				switch sipMsg.ResponseLine.StatusCode {
				case 200:
					target.Event("ok")
				}
			}
		}
	}
}
