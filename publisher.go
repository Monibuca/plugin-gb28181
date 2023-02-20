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
	. "m7s.live/engine/v4/codec"
	"m7s.live/engine/v4/codec/mpegps"
	"m7s.live/engine/v4/codec/mpegts"
	. "m7s.live/engine/v4/track"
	"m7s.live/engine/v4/util"
	"m7s.live/plugin/gb28181/v4/utils"
)

type GBPublisher struct {
	Publisher
	InviteOptions
	channel     *Channel
	inviteRes   sip.Response
	parser      mpegps.MpegPsStream
	lastSeq     uint16
	udpCache    *utils.PriorityQueueRtp
	dumpFile    *os.File
	dumpPrint   io.Writer
	lastReceive time.Time
	reorder     util.RTPReorder[*rtp.Packet]
}

func (p *GBPublisher) PrintDump(s string) {
	if p.dumpPrint != nil {
		p.dumpPrint.Write([]byte(s))
	}
}

func (p *GBPublisher) OnEvent(event any) {
	if p.channel == nil {
		// p.parser.EsHandler = p
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
		// p.parser.EsHandler = p
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

func (p *GBPublisher) ReceiveVideo(es mpegps.MpegPsEsStream) {
	if p.VideoTrack == nil {
		switch es.Type {
		case mpegts.STREAM_TYPE_H264:
			p.VideoTrack = NewH264(p.Publisher.Stream)
		case mpegts.STREAM_TYPE_H265:
			p.VideoTrack = NewH265(p.Publisher.Stream)
		default:
			//推测编码类型
			var maybe264 H264NALUType
			maybe264 = maybe264.Parse(es.Buffer[4])
			switch maybe264 {
			case NALU_Non_IDR_Picture,
				NALU_IDR_Picture,
				NALU_SEI,
				NALU_SPS,
				NALU_PPS,
				NALU_Access_Unit_Delimiter:
				p.VideoTrack = NewH264(p.Publisher.Stream)
			default:
				p.Info("maybe h265", zap.Uint8("type", maybe264.Byte()))
				p.VideoTrack = NewH265(p.Publisher.Stream)
			}
		}
	}
	payload, pts, dts := es.Buffer, es.PTS, es.DTS
	if len(payload) > 10 {
		p.PrintDump(fmt.Sprintf("<td>pts:%d dts:%d data: % 2X</td>", pts, dts, payload[:10]))
	} else {
		p.PrintDump(fmt.Sprintf("<td>pts:%d dts:%d data: % 2X</td>", pts, dts, payload))
	}
	if dts == 0 {
		dts = pts
	}
	// if binary.BigEndian.Uint32(payload) != 1 {
	// 	panic("not annexb")
	// }
	p.WriteAnnexB(pts, dts, payload)
}
func (p *GBPublisher) ReceiveAudio(es mpegps.MpegPsEsStream) {
	ts, payload := es.PTS, es.Buffer
	if p.AudioTrack == nil {
		switch es.Type {
		case mpegts.STREAM_TYPE_G711A:
			p.AudioTrack = NewG711(p.Publisher.Stream, true)
		case mpegts.STREAM_TYPE_G711U:
			p.AudioTrack = NewG711(p.Publisher.Stream, false)
		case mpegts.STREAM_TYPE_AAC:
			p.AudioTrack = NewAAC(p.Publisher.Stream)
			p.WriteADTS(ts, payload)
		case 0: //推测编码类型
			if payload[0] == 0xff && payload[1]>>4 == 0xf {
				p.AudioTrack = NewAAC(p.Publisher.Stream)
				p.WriteADTS(ts, payload)
			}
		default:
			p.Error("audio type not supported yet", zap.Uint8("type", es.Type))
		}
	} else if es.Type == mpegts.STREAM_TYPE_AAC {
		p.WriteADTS(ts, payload)
	} else {
		p.WriteRaw(ts, payload)
	}
}

// 解析rtp封装 https://www.ietf.org/rfc/rfc2250.txt
func (p *GBPublisher) PushPS(rtp *rtp.Packet) {
	if p.parser.EsHandler == nil {
		p.parser.EsHandler = p
		p.lastSeq = rtp.SequenceNumber - 1
	}
	if conf.IsMediaNetworkTCP() {
		p.parser.Feed(rtp.Payload)
		p.lastSeq = rtp.SequenceNumber
	} else {
		for rtp = p.reorder.Push(rtp.SequenceNumber, rtp); rtp != nil; rtp = p.reorder.Pop() {
			if rtp.SequenceNumber != p.lastSeq+1 {
				fmt.Println("drop", rtp.SequenceNumber, p.lastSeq)
				p.parser.Drop()
				if p.VideoTrack != nil {
					p.SetLostFlag()
				}
			}
			p.parser.Feed(rtp.Payload)
			p.lastSeq = rtp.SequenceNumber
		}
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
		conf.udpPorts.Recycle(port)
		plugin.Error("listen media server udp err", zap.String("addr", addr), zap.Error(err))
		return 0, err
	}
	p.SetIO(conn)
	go func() {
		defer conn.Close()
		bufUDP := make([]byte, networkBuffer)
		plugin.Info("Media udp server start.", zap.Uint16("port", port))
		defer plugin.Info("Media udp server stop", zap.Uint16("port", port))
		defer conf.udpPorts.Recycle(port)
		dumpLen := make([]byte, 6)
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
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
			conn.SetReadDeadline(time.Now().Add(time.Second * 10))
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
		defer conf.tcpPorts.Recycle(port)
		plugin.Error("listen media server tcp err", zap.String("addr", addr), zap.Error(err))
		return 0, err
	}
	go func() {
		plugin.Info("Media tcp server start.", zap.Uint16("port", port))
		defer conf.tcpPorts.Recycle(port)
		defer plugin.Info("Media tcp server stop", zap.Uint16("port", port))
		conn, err := listen.Accept()
		listen.Close()
		p.SetIO(conn)
		if err != nil {
			plugin.Error("Accept err=", zap.Error(err))
			return
		}
		var rtpPacket rtp.Packet
		lenBuf := make([]byte, 2)
		defer conn.Close()
		for err == nil {
			if _, err = io.ReadFull(conn, lenBuf); err != nil {
				return
			}
			ps := make([]byte, binary.BigEndian.Uint16(lenBuf))
			if _, err = io.ReadFull(conn, ps); err != nil {
				return
			}
			if err := rtpPacket.Unmarshal(ps); err != nil {
				plugin.Error("gb28181 decode rtp error:", zap.Error(err))
			} else if !p.IsClosed() {
				p.PushPS(&rtpPacket)
			}
		}
	}()
	return
}
