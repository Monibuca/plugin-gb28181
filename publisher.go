package gb28181

import (
	"github.com/Monibuca/engine/v3"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
	. "github.com/Monibuca/utils/v3"
)

type Publisher struct {
	*engine.Stream
	psPacket  []byte
	parser    utils.DecPSPackage
	pushVideo func(uint32, uint32, []byte)
	pushAudio func(uint32, []byte)
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
func (p *Publisher) PushPS(ps []byte, ts uint32) {
	if len(ps) >= 4 && BigEndian.Uint32(ps) == utils.StartCodePS {
		if p.psPacket != nil {
			p.parser.Read(p.psPacket)
			if p.parser.VideoPayload != nil {
				if p.parser.DTS != 0 {
					ts = p.parser.DTS
					p.pushVideo(ts/90, (p.parser.PTS/90 - p.parser.DTS/90), p.parser.VideoPayload)
				} else {
					p.pushVideo(ts/90, 0, p.parser.VideoPayload)
				}
			}
			if p.parser.AudioPayload != nil {
				p.pushAudio(ts/8, p.parser.AudioPayload)
			}
			p.psPacket = nil
		}
		p.psPacket = append(p.psPacket, ps...)
	} else if p.psPacket != nil {
		p.psPacket = append(p.psPacket, ps...)
	}
}
