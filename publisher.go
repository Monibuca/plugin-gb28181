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
	psPacket []byte
	parser utils.DecPSPackage
}

func NewPublisher() *Publisher {
	var p Publisher
	p.Type = "GB28181"
	p.AutoUnPublish = true

	return &p
}
func (p *Publisher) Publish(c *Channel, start string) bool {
	if start == "0" {
		return p.Publisher.Publish(fmt.Sprintf("%s/%s", c.device.ID, c.DeviceID))
	}
	return p.Publisher.Publish(fmt.Sprintf("%s/%s", c.DeviceID, start))
}
func (p *Publisher)PushPS(ps []byte,ts uint32) {
	p.Update()
	if len(ps) >= 4 && BigEndian.Uint32(ps) == utils.StartCodePS {
		if p.psPacket != nil{
			if err := p.parser.Read(p.psPacket); err == nil {
					for _, payload := range codec.SplitH264(p.parser.VideoPayload) {
						p.OriginVideoTrack.Push(engine.VideoPack{Timestamp:ts / 90, Payload: payload})
						if p.OriginVideoTrack.CodecID == 0 && p.OriginVideoTrack.RtmpTag != nil {
							p.OriginVideoTrack.CodecID = 7
							p.SetOriginVT(p.OriginVideoTrack)
						}
					}
					if p.parser.AudioPayload != nil {
						// switch parser.AudioStreamType {
						// case G711A:
						// 	rtp.AudioInfo.SoundFormat = 7
						// 	rtp.AudioInfo.SoundRate = 8000
						// 	rtp.AudioInfo.SoundSize = 16
						// 	asc := rtp.AudioInfo.SoundFormat << 4
						// 	asc = asc + 1<<1
						// 	rtp.PushAudio(rtp.Timestamp, append([]byte{asc}, parser.AudioPayload...))
						// }
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