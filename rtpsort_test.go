package gb28181

import (
	"fmt"
	"github.com/Monibuca/engine/v3"
	"github.com/Monibuca/plugin-gb28181/v3/utils"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/pion/rtp"
	"reflect"
	"testing"
)

var items = []uint16{
	63401, 63397, 63395, 63399, 63398, 63400, 63396,
	63402, 63410, 63403, 63404, 63405, 63406, 63407,
	63408, 63409, 63420, 63411, 63412, 63413, 63414,
	63415, 63416, 63417, 63418, 63419, 63421, 63428,
}

func _pushPsWithCache(p *Publisher, rtp *rtp.Packet) {
	if config.UdpCacheSize > 0 {
		if rtp.SequenceNumber < p.lastSeq { //序号小于第一个包的丢弃
			return
		}

		p.UdpCache.Push(*rtp)
		rtpTmp, _ := p.UdpCache.Pop()
		rtp = &rtpTmp
	}

	if p.lastSeq != 0 {

		// rtp序号不连续，丢弃PS
		if p.lastSeq+1 != rtp.SequenceNumber {
			if config.UdpCacheSize > 0 && p.UdpCache.Len() < config.UdpCacheSize {
				p.UdpCache.Push(*rtp)
				return
			}
			p.parser.Reset()
		}
	}

	p.lastSeq = rtp.SequenceNumber

	fmt.Println("rtp.SequenceNumber:", rtp.SequenceNumber)

}

// run with go test -gcflags=all=-l -v
func TestRtpSort(t *testing.T) {
	publisher := Publisher{
		Stream: &engine.Stream{
			StreamPath: "live/test",
		},
		UdpCache: utils.NewPq(),
	}
	config.UdpCacheSize = 7

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&publisher), "PushPS", _pushPsWithCache)
	defer patches.Reset()

	for i := 0; i < len(items); i++ {
		rtpPacket := &rtp.Packet{Header: rtp.Header{SequenceNumber: items[i]}}
		publisher.PushPS(rtpPacket)

	}

}

func TestPqSort(t *testing.T) {
	pq := utils.NewPq()
	for i := 0; i < len(items); i++ {
		rtpPacket := rtp.Packet{Header: rtp.Header{SequenceNumber: items[i]}}
		pq.Push(rtpPacket)
	}

	for pq.Len() > 0 {
		rtpPacket, _ := pq.Pop()
		fmt.Println("packet seq:", rtpPacket.SequenceNumber)
	}
}
