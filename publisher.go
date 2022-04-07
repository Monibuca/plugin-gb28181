package gb28181

import (
	"encoding/binary"

	"github.com/pion/rtp"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/common"
	. "m7s.live/engine/v4/track"
	"m7s.live/plugin/gb28181/v4/utils"
)

type GBPublisher struct {
	Publisher
	StreamPath string
	parser     *utils.DecPSPackage
	pushVideo  func(uint32, uint32, common.AnnexBFrame)
	pushAudio  func(uint32, common.AVCCFrame)
	OnClose    func()
	lastSeq    uint16
	udpCache   *utils.PriorityQueueRtp
	config     *GB28181Config
}

func (p *GBPublisher) PushVideo(ts uint32, cts uint32, payload []byte) {
	//p.VideoTrack.WriteAnnexB(ts, cts, payload)
	p.pushVideo(ts, cts, payload)
}
func (p *GBPublisher) PushAudio(ts uint32, payload []byte) {
	//p.VideoTrack.WriteAVCC(ts, payload)
	p.pushAudio(ts, payload)
}

func (p *GBPublisher) Publish() (result bool) {
	if err := plugin.Publish(p.StreamPath, p); err == nil {
		p.pushVideo = func(pts uint32, dts uint32, frame common.AnnexBFrame) {
			var vt common.VideoTrack
			switch p.parser.VideoStreamType {
			case utils.StreamTypeH264:
				vt = NewH264(p.Stream)
			case utils.StreamTypeH265:
				vt = NewH265(p.Stream)
			default:
				return
			}
			vt.WriteAnnexB(pts, dts, frame)
			p.pushVideo = vt.WriteAnnexB
		}
		p.pushAudio = func(ts uint32, payload common.AVCCFrame) {
			switch p.parser.AudioStreamType {
			case utils.G711A:
			case utils.G711A + 1:
				at := NewG711(p.Stream, true)
				at.SampleRate = 8000
				at.SampleSize = 16
				at.Channels = 1
				// at.ExtraData = []byte{(at.CodecID << 4) | (1 << 1)}
				at.WriteAVCC(ts, payload)
				p.pushAudio = at.WriteAVCC
			}
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
	// p.Update()
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
