package gb28181

import (
	"encoding/binary"
	"sync/atomic"

	"github.com/ghettovoice/gosip/sip"
	"github.com/pion/rtp/v2"
	"go.uber.org/zap"
	. "m7s.live/engine/v4"
	. "m7s.live/engine/v4/track"
	"m7s.live/plugin/gb28181/v4/utils"
)

type GBPublisher struct {
	Publisher
	Start     string
	End       string
	SSRC      uint32
	channel   *Channel
	inviteRes *sip.Response
	parser    *utils.DecPSPackage
	lastSeq   uint16
	udpCache  *utils.PriorityQueueRtp
}

func (p *GBPublisher) OnEvent(event any) {
	switch v := event.(type) {
	case IPublisher:
		if p.IsLive() {
			p.Type = "GB28181 Live"
			p.channel.LivePublisher = p
		} else {
			p.Type = "GB28181 Playback"
			p.channel.RecordPublisher = p
		}
		conf.publishers.Add(p.SSRC, p)
		if p.Equal(v) { //第一任

		} else {
			//删除前任
			conf.publishers.Delete(v.(*GBPublisher).SSRC)
			p.Publisher.OnEvent(v)
		}
	case SEwaitPublish:
		//掉线自动重新拉流
		if p.IsLive() {
			atomic.StoreInt32(&p.channel.state, 0)
			p.channel.LivePublisher = nil
			go p.channel.Invite("", "")
		}
	case SEclose, SEKick:
		if p.IsLive() {
			p.channel.LivePublisher = nil
		} else {
			p.channel.RecordPublisher = nil
		}
		p.Publisher.OnEvent(v)
		conf.publishers.Delete(p.SSRC)
		p.Bye()
	default:
		p.Publisher.OnEvent(v)
	}
}

func (p *GBPublisher) Bye() int {
	res := p.inviteRes
	if res == nil {
		return 404
	}
	defer p.Stop()
	defer atomic.StoreInt32(&p.channel.state, 0)
	p.inviteRes = nil
	bye := p.channel.CreateRequst(sip.BYE)
	from, _ := (*res).From()
	to, _ := (*res).To()
	callId, _ := (*res).CallID()
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

func (p *GBPublisher) IsLive() bool {
	return p.Start == ""
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
			p.parser.Read(rtp.Timestamp, p)
			p.parser.Reset()
		}
		p.parser.Write(ps)
	} else if p.parser.Len() > 0 {
		p.parser.Write(ps)
	}
}
