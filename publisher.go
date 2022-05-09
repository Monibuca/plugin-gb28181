package gb28181

import (
	"encoding/binary"

	"github.com/pion/rtp/v2"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/common"
	"m7s.live/engine/v4/track"
	. "m7s.live/engine/v4/track"
	"m7s.live/plugin/gb28181/v4/utils"
)

type GBPublisher struct {
	Publisher
	StreamPath string
	parser     *utils.DecPSPackage
	OnClose    func()
	lastSeq    uint16
	udpCache   *utils.PriorityQueueRtp
	config     *GB28181Config
	vt         common.VideoTrack
	at         *track.G711
}

func (p *GBPublisher) PushVideo(pts uint32, dts uint32, payload []byte) {
	if p.vt == nil {
		switch p.parser.VideoStreamType {
		case utils.StreamTypeH264:
			p.vt = NewH264(p.Stream)
		case utils.StreamTypeH265:
			p.vt = NewH265(p.Stream)
		default:
			return
		}
	}
	p.vt.WriteAnnexB(pts, dts, payload)
}
func (p *GBPublisher) PushAudio(ts uint32, payload []byte) {
	if p.at == nil {
		switch p.parser.AudioStreamType {
		case utils.G711A:
			at := NewG711(p.Stream, true)
			at.SampleRate = 8000
			at.SampleSize = 16
			at.Channels = 1
			at.AVCCHead = []byte{(byte(at.CodecID) << 4) | (1 << 1)}
			p.at = at
		case utils.G711A + 1:
			at := NewG711(p.Stream, false)
			at.SampleRate = 8000
			at.SampleSize = 16
			at.Channels = 1
			at.AVCCHead = []byte{(byte(at.CodecID) << 4) | (1 << 1)}
			p.at = at
		default:
			return
		}
	}
	p.at.WriteAVCC(ts, payload)
}

func (p *GBPublisher) Publish() (result bool) {
	if err := plugin.Publish(p.StreamPath, p); err == nil {
		if p.vt != nil {
			p.vt.Detach()
			p.vt = nil
		}
		if p.at != nil {
			p.at.Detach()
			p.at = nil
		}
		return true
	}
	return false
}

func (p *GBPublisher) PushPS(rtp *rtp.Packet) {
	originRtp := *rtp
	if p.config.UdpCacheSize > 0 && !p.config.IsMediaNetworkTCP() {
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
			if p.config.UdpCacheSize > 0 && !p.config.IsMediaNetworkTCP() {
				if p.udpCache.Len() < p.config.UdpCacheSize {
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
