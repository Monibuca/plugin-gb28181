package gb28181

import (
	"github.com/Monibuca/engine/v3"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
	. "github.com/Monibuca/utils/v3"
	"github.com/pion/rtp"
)

type Publisher struct {
	*engine.Stream
	parser    utils.DecPSPackage
	pushVideo func(uint32, uint32, []byte)
	pushAudio func(uint32, []byte)
	lastSeq   uint16
}

func (p *Publisher) PushVideo(ts uint32, cts uint32, payload []byte) {
	p.pushVideo(ts, cts, payload)
}
func (p *Publisher) PushAudio(ts uint32, payload []byte) {
	p.pushAudio(ts, payload)
}
func (p *Publisher) Publish() (result bool) {
	if result = p.Stream.Publish(); result {
		p.pushVideo = func(ts uint32, cts uint32, payload []byte) {
			var vt *engine.VideoTrack
			switch p.parser.VideoStreamType {
			case utils.StreamTypeH264:
				vt = p.Stream.NewVideoTrack(7)
			case utils.StreamTypeH265:
				vt = p.Stream.NewVideoTrack(12)
			default:
				return
			}
			vt.PushAnnexB(ts, cts, payload)
			p.pushVideo = vt.PushAnnexB
		}
		p.pushAudio = func(ts uint32, payload []byte) {
			switch p.parser.AudioStreamType {
			case utils.G711A:
				at := p.Stream.NewAudioTrack(7)
				at.SoundRate = 8000
				at.SoundSize = 16
				at.Channels = 1
				at.ExtraData = []byte{(at.CodecID << 4) | (1 << 1)}
				at.PushRaw(ts, payload)
				p.pushAudio = at.PushRaw
				// case utils.G711U:
				// 	at := p.Stream.NewAudioTrack(8)
				// 	at.SoundRate = 8000
				// 	at.SoundSize = 16
				// 	asc := at.CodecID << 4
				// 	asc = asc + 1<<1
				// 	at.ExtraData = []byte{asc}
				// 	at.PushRaw(pack)
				// 	p.pushAudio = at.PushRaw
			}
		}
	}
	return
}
func (p *Publisher) PushPS(rtp *rtp.Packet) {
	ps := rtp.Payload
	if p.lastSeq != 0 {
		// rtp序号不连续，丢弃PS
		if p.lastSeq+1 != rtp.SequenceNumber {
			p.parser.Reset()
		}
	}
	p.lastSeq = rtp.SequenceNumber
	p.Update()
	if len(ps) >= 4 && BigEndian.Uint32(ps) == utils.StartCodePS {
		if p.parser.Len() > 0 {
			p.parser.Uint32()
			p.parser.Read(rtp.Timestamp, p)
			p.parser.Reset()
		}
		p.parser.Write(ps)
	} else if p.parser.Len() > 0 {
		p.parser.Write(ps)
	}
}
