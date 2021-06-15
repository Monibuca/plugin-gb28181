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
	pushVideo func(engine.VideoPack)
	pushAudio func(engine.AudioPack)
}

func (p *Publisher) Publish() (result bool) {
	if result = p.Stream.Publish(); result {
		p.pushVideo = func(pack engine.VideoPack) {
			var vt *engine.VideoTrack
			switch p.parser.VideoStreamType {
			case 0x1B:
				vt = p.Stream.NewVideoTrack(7)
			case 0x24:
				vt = p.Stream.NewVideoTrack(12)
			default:
				return
			}
			vt.PushAnnexB(pack)
			p.pushVideo = vt.PushAnnexB
		}
		p.pushAudio = func(pack engine.AudioPack) {
			switch p.parser.AudioStreamType {
			case utils.G711A:
				at := p.Stream.NewAudioTrack(7)
				at.SoundRate = 8000
				at.SoundSize = 16
				at.ExtraData = []byte{(at.CodecID << 4) | (1 << 1)}
				at.PushRaw(pack)
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
			if err := p.parser.Read(p.psPacket); err == nil {
				p.pushVideo(engine.VideoPack{Timestamp: ts / 90, Payload: p.parser.VideoPayload})
				if p.parser.AudioPayload != nil {
					p.pushAudio(engine.AudioPack{Timestamp: ts / 90, Raw: p.parser.AudioPayload})
				}
			} else {
				Print(err)
			}
			p.psPacket = nil
		}
		p.psPacket = append(p.psPacket, ps...)
	} else if p.psPacket != nil {
		p.psPacket = append(p.psPacket, ps...)
	}
}
