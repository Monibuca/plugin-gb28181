package gb28181

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/ghettovoice/gosip/sip"
	"github.com/pion/rtp/v2"
	"go.uber.org/zap"
	. "m7s.live/engine/v4"
	. "m7s.live/engine/v4/track"
	"m7s.live/engine/v4/util"
	"m7s.live/plugin/gb28181/v4/utils"
)

type GBPublisher struct {
	Publisher
	InviteOptions
	channel     *Channel
	inviteRes   sip.Response
	parser      *utils.DecPSPackage
	lastSeq     uint16
	udpCache    *utils.PriorityQueueRtp
	dumpFile    *os.File
	dumpPrint   io.Writer
	lastReceive time.Time
}

func (p *GBPublisher) PrintDump(s string) {
	if p.dumpPrint != nil {
		p.dumpPrint.Write([]byte(s))
	}
}

func (p *GBPublisher) OnEvent(event any) {
	if p.channel == nil {
		p.IO.OnEvent(event)
		return
	}
	switch event.(type) {
	case IPublisher:
		if p.IsLive() {
			p.Type = "GB28181 Live"
			p.channel.LivePublisher = p
		} else {
			p.Type = "GB28181 Playback"
			p.channel.RecordPublisher = p
		}
		conf.publishers.Add(p.SSRC, p)
		if err := error(nil); p.dump != "" {
			if p.dumpFile, err = os.OpenFile(p.dump, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
				p.Error("open dump file failed", zap.Error(err))
			}
		}
	case SEwaitPublish:
		//掉线自动重新拉流
		if p.IsLive() {
			if p.channel.LivePublisher != nil {
				p.channel.LivePublisher = nil
				p.channel.liveInviteLock.Unlock()
			}
			go p.channel.Invite(InviteOptions{})
		}
	case SEclose, SEKick:
		if p.IsLive() {
			if p.channel.LivePublisher != nil {
				p.channel.LivePublisher = nil
				p.channel.liveInviteLock.Unlock()
			}
		} else {
			p.channel.RecordPublisher = nil
		}
		conf.publishers.Delete(p.SSRC)
		if p.dumpFile != nil {
			p.dumpFile.Close()
		}
		p.Bye()
	}
	p.Publisher.OnEvent(event)
}

func (p *GBPublisher) Bye() int {
	res := p.inviteRes
	if res == nil {
		return 404
	}
	defer p.Stop()
	p.inviteRes = nil
	bye := p.channel.CreateRequst(sip.BYE)
	from, _ := res.From()
	to, _ := res.To()
	callId, _ := res.CallID()
	bye.ReplaceHeaders(from.Name(), []sip.Header{from})
	bye.ReplaceHeaders(to.Name(), []sip.Header{to})
	bye.ReplaceHeaders(callId.Name(), []sip.Header{callId})
	resp, err := p.channel.device.SipRequestForResponse(bye)
	if err != nil {
		p.Error("Bye", zap.Error(err))
		return 500
	}
	return int(resp.StatusCode())
}

func (p *GBPublisher) PushVideo(pts uint32, dts uint32, payload []byte) {
	if p.VideoTrack == nil {
		switch p.parser.VideoStreamType {
		case utils.StreamTypeH264:
			p.VideoTrack = NewH264(p.Publisher.Stream)
		case utils.StreamTypeH265:
			p.VideoTrack = NewH265(p.Publisher.Stream)
		default:
			return
		}
	}
	p.PrintDump(fmt.Sprintf("<td>pts:%d dts:%d data: % 2X</td>", pts, dts, payload[:10]))
	p.VideoTrack.WriteAnnexB(pts, dts, payload)
}
func (p *GBPublisher) PushAudio(ts uint32, payload []byte) {
	if p.AudioTrack == nil {
		switch p.parser.AudioStreamType {
		case utils.G711A:
			at := NewG711(p.Publisher.Stream, true)
			at.Audio.SampleRate = 8000
			at.Audio.SampleSize = 16
			at.Channels = 1
			at.AVCCHead = []byte{(byte(at.CodecID) << 4) | (1 << 1)}
			p.AudioTrack = at
		case utils.G711A + 1:
			at := NewG711(p.Publisher.Stream, false)
			at.Audio.SampleRate = 8000
			at.Audio.SampleSize = 16
			at.Channels = 1
			at.AVCCHead = []byte{(byte(at.CodecID) << 4) | (1 << 1)}
			p.AudioTrack = at
		default:
			return
		}
	}
	p.AudioTrack.WriteAVCC(ts, payload)
}

func (p *GBPublisher) PushPS(rtp *rtp.Packet) {
	originRtp := *rtp
	if conf.UdpCacheSize > 0 && !conf.IsMediaNetworkTCP() {
		//序号小于第一个包的丢弃,rtp包序号达到65535后会从0开始，所以这里需要判断一下
		if rtp.SequenceNumber < p.lastSeq && p.lastSeq-rtp.SequenceNumber < utils.MaxRtpDiff {
			return
		}
		p.udpCache.Push(*rtp)
		rtpTmp, _ := p.udpCache.Pop()
		rtp = &rtpTmp
	}
	ps := rtp.Payload
	if p.lastSeq != 0 {
		// rtp序号不连续，丢弃PS
		if p.lastSeq+1 != rtp.SequenceNumber {
			if conf.UdpCacheSize > 0 && !conf.IsMediaNetworkTCP() {
				if p.udpCache.Len() < conf.UdpCacheSize {
					p.udpCache.Push(*rtp)
					return
				} else {
					p.udpCache.Empty()
					rtp = &originRtp // 还原rtp包，而不是使用缓存中，避免rtp序号断裂
				}
			}
			p.parser.Reset()
		}
	}
	p.lastSeq = rtp.SequenceNumber
	if p.parser == nil {
		p.parser = new(utils.DecPSPackage)
	}
	if len(ps) >= 4 && binary.BigEndian.Uint32(ps) == utils.StartCodePS {
		if p.parser.Len() > 0 {
			p.parser.Skip(4)
			p.PrintDump("</td></tr>")
			p.PrintDump("<tr>")
			p.parser.Read(rtp.Timestamp, p)
			p.PrintDump("</tr>")
			p.PrintDump("<tr  class=gray><td colspan=12>")
			p.parser.Reset()
		}
		p.parser.Write(ps)
	} else if p.parser.Len() > 0 {
		p.parser.Write(ps)
	}
}
func (p *GBPublisher) Replay(f *os.File) (err error) {
	var rtpPacket rtp.Packet
	defer f.Close()
	if p.dumpPrint != nil {
		p.PrintDump(`<style  type="text/css">
		.gray {
			color: gray;
		}
		</style>
		`)
		p.PrintDump("<table>")
		defer p.PrintDump("</table>")
	}
	p.PrintDump("<tr class=gray><td colspan=12>")
	var t uint16
	for l := make([]byte, 6); !p.IsClosed(); time.Sleep(time.Millisecond * time.Duration(t)) {
		_, err = f.Read(l)
		if err != nil {
			return
		}
		payload := make([]byte, util.ReadBE[int](l[:4]))
		t = util.ReadBE[uint16](l[4:])
		p.PrintDump(fmt.Sprintf("[<b>%d</b> %d]", t, len(payload)))
		_, err = f.Read(payload)
		if err != nil {
			return
		}
		rtpPacket.Unmarshal(payload)
		p.PushPS(&rtpPacket)
	}
	return
}

func (p *GBPublisher) ListenUDP() (port uint16, err error) {
	var rtpPacket rtp.Packet
	networkBuffer := 1048576
	port, err = conf.udpPorts.GetPort()
	if err != nil {
		return
	}
	addr := fmt.Sprintf(":%d", port)
	mediaAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, err := net.ListenUDP("udp", mediaAddr)
	if err != nil {
		plugin.Error("listen media server udp err", zap.String("addr", addr), zap.Error(err))
		return 0, err
	}
	p.SetIO(conn)
	go func() {
		bufUDP := make([]byte, networkBuffer)
		plugin.Info("Media udp server start.", zap.Uint16("port", port))
		defer plugin.Info("Media udp server stop", zap.Uint16("port", port))
		defer conf.udpPorts.Recycle(port)
		dumpLen := make([]byte, 6)
		for n, _, err := conn.ReadFromUDP(bufUDP); err == nil; n, _, err = conn.ReadFromUDP(bufUDP) {
			ps := bufUDP[:n]
			if err := rtpPacket.Unmarshal(ps); err != nil {
				plugin.Error("Decode rtp error:", zap.Error(err))
			}
			if p.dumpFile != nil {
				util.PutBE(dumpLen[:4], n)
				if p.lastReceive.IsZero() {
					util.PutBE(dumpLen[4:], 0)
				} else {
					util.PutBE(dumpLen[4:], uint16(time.Since(p.lastReceive).Milliseconds()))
				}
				p.lastReceive = time.Now()
				p.dumpFile.Write(dumpLen)
				p.dumpFile.Write(ps)
			}
			p.PushPS(&rtpPacket)
		}
	}()
	return
}


func (p *GBPublisher) ListenTCP() (port uint16, err error) {
	port, err = conf.tcpPorts.GetPort()
	if err != nil {
		return
	}
	addr := fmt.Sprintf(":%d", port)
	mediaAddr, _ := net.ResolveTCPAddr("tcp", addr)
	listen, err := net.ListenTCP("tcp", mediaAddr)
	if err != nil {
		plugin.Error("listen media server tcp err", zap.String("addr", addr), zap.Error(err))
		return 0, err
	}
	// p.SetIO(conn)
	go func() {
		plugin.Info("Media tcp server start.", zap.Uint16("port", port))
		defer listen.Close()
		defer conf.tcpPorts.Recycle(port)
		defer plugin.Info("Media tcp server stop", zap.Uint16("port", port))
		for {
			conn, err := listen.Accept()
			if err != nil {
				plugin.Error("Accept err=", zap.Error(err))
			}
			go processTcpMediaConn(conf, conn)
		}
		
	}()
	return
}
