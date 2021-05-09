package gb28181

import (
	"fmt"

	"github.com/Monibuca/engine/v3"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
	. "github.com/Monibuca/utils/v3"
	"github.com/Monibuca/utils/v3/codec"
)

type Publisher struct {
	engine.Publisher
	psPacket  []byte
	parser    utils.DecPSPackage
	pushVideo func(engine.VideoPack)
	pushAudio func(uint32, []byte)
}

func (p *Publisher) Publish(c *Channel, start string) (result bool) {
	defer func() {
		if result {
			vt := engine.NewVideoTrack()
			p.pushVideo = func(pack engine.VideoPack) {
				vt.Push(pack)
				if vt.CodecID == 0 && vt.RtmpTag != nil {
					vt.CodecID = 7
					p.SetOriginVT(vt)
					p.pushVideo = vt.Push
				}
			}
			at := engine.NewAudioTrack()
			p.pushAudio = func(ts uint32, payload []byte) {
				switch p.parser.AudioStreamType {
				case utils.G711A:
					at.SoundFormat = 7
					at.SoundRate = 8000
					at.SoundSize = 16
					asc := at.SoundFormat << 4
					asc = asc + 1<<1
					at.RtmpTag = []byte{asc}
					at.Push(ts, payload)
					p.SetOriginAT(at)
					p.pushAudio = at.Push
				// case utils.G711U:
				// 	at.SoundFormat = 8
				// 	at.SoundRate = 8000
				// 	at.SoundSize = 16
				// 	asc := at.SoundFormat << 4
				// 	asc = asc + 1<<1
				// 	at.RtmpTag = []byte{asc}
				// 	at.Push(ts, payload)
				// 	p.SetOriginAT(at)
				// 	p.pushAudio = at.Push
				}
			}
		}
	}()
	if start == "0" {
		result = p.Publisher.Publish(fmt.Sprintf("%s/%s", c.device.ID, c.DeviceID))
	} else {
		result = p.Publisher.Publish(fmt.Sprintf("%s/%s", c.DeviceID, start))
	}
	return
}
func (p *Publisher) PushPS(ps []byte, ts uint32) {
	if len(ps) >= 4 && BigEndian.Uint32(ps) == utils.StartCodePS {
		if p.psPacket != nil {
			if err := p.parser.Read(p.psPacket); err == nil {
				for _, payload := range codec.SplitH264(p.parser.VideoPayload) {
					p.pushVideo(engine.VideoPack{Timestamp: ts / 90, Payload: payload})
				}
				if p.parser.AudioPayload != nil {
					p.pushAudio(ts / 90, p.parser.AudioPayload)
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
