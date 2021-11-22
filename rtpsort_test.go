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

// 测试rtp序号数据
var items = []uint16{
	65526, 65530, 65524, 65525, 65527, 65528, 65529,
	0, 65533, 65531, 65532, 65534, 65535, 1,
	3, 6, 5, 4, 2, 8, 7,
}

func _pushPsWithCache(p *Publisher, rtp *rtp.Packet) {
	if config.UdpCacheSize > 0 && !config.TCP {
		//序号小于第一个包的丢弃,rtp包序号达到65535后会从0开始，所以这里需要判断一下
		if rtp.SequenceNumber < p.lastSeq && p.lastSeq-rtp.SequenceNumber < utils.MaxRtpDiff {
			return
		}
		p.udpCache.Push(*rtp)
		rtpTmp, _ := p.udpCache.Pop()
		rtp = &rtpTmp
	}

	if p.lastSeq != 0 {
		// rtp序号不连续，丢弃PS
		if p.lastSeq+1 != rtp.SequenceNumber {
			if config.UdpCacheSize > 0 && !config.TCP {
				if p.udpCache.Len() < config.UdpCacheSize {
					p.udpCache.Push(*rtp)
					return
				} else {
					p.udpCache.Empty()
				}
			}
			p.parser.Reset()
		}
	}

	p.lastSeq = rtp.SequenceNumber
	fmt.Println("rtp.SequenceNumber:", rtp.SequenceNumber)

}

// 如果运行失败可以关闭gc,go test -gcflags=all=-l -v
func TestRtpSort(t *testing.T) {
	publisher := Publisher{
		Stream: &engine.Stream{
			StreamPath: "live/test",
		},
		udpCache: utils.NewPqRtp(),
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
	pq := utils.NewPqRtp()
	for i := 0; i < len(items); i++ {
		rtpPacket := rtp.Packet{Header: rtp.Header{SequenceNumber: items[i]}}
		pq.Push(rtpPacket)
	}

	for pq.Len() > 0 {
		rtpPacket, _ := pq.Pop()
		fmt.Println("packet seq:", rtpPacket.SequenceNumber)
	}
}
