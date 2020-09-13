package gb28181

import (
	. "github.com/Monibuca/engine/v2"
	"github.com/Monibuca/engine/v2/util"

	"github.com/Monibuca/plugin-gb28181/transaction"
	rtp "github.com/Monibuca/plugin-rtp"
	"github.com/mask-pp/rtp-ps/packet"
	"net"
	"strconv"
	"strings"
	"time"
)

func init() {
	println(packet.NewRtpParsePacket())
}

func onPublish(device *transaction.Device) (port int) {
	rtpPublisher := new(rtp.RTP)
	if !rtpPublisher.Publish("gb28181/" + device.ID) {
		return
	}
	parser := packet.NewRtpParsePacket()
	addr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	networkBuffer := 1048576
	if err := conn.SetReadBuffer(networkBuffer); err != nil {
		Printf("udp server video conn set read buffer error, %v", err)
	}
	if err := conn.SetWriteBuffer(networkBuffer); err != nil {
		Printf("udp server video conn set write buffer error, %v", err)
	}
	la := conn.LocalAddr().String()
	strPort := la[strings.LastIndex(la, ":")+1:]
	if port, err = strconv.Atoi(strPort); err != nil {
		return
	}
	go func() {
		bufUDP := make([]byte, 1048576)
		Printf("udp server start listen video port[%d]",port)
		defer Printf("udp server stop listen video port[%d]", port)
		timer := time.Unix(0, 0)
		var psPackets []byte
		var timestamp uint32
		for {
			if n, _, err := conn.ReadFromUDP(bufUDP); err == nil {
				elapsed := time.Now().Sub(timer)
				if elapsed >= 30*time.Second {
					Printf("Package recv from VConn.len:%d\n", n)
					timer = time.Now()
				}
				pack := &rtp.RTPPack{
					Type: rtp.RTP_TYPE_VIDEO,
				}
				pack.Unmarshal(bufUDP[:n])
				if len(pack.Payload)>=4 && util.BigEndian.Uint32(pack.Payload) == packet.StartCodePS {
					if psPackets != nil {
						if nalus, err := parser.Read(psPackets); err == nil {
							for _, nalu := range nalus {
								rtpPublisher.WriteNALU(timestamp, nalu)
							}
						}
					}
					timestamp = pack.Timestamp
					psPackets = append([]byte{}, pack.Payload...)
				} else {
					psPackets = append(psPackets, pack.Payload...)
				}
				//s.HandleRTP(pack)
			} else {
				Println("udp server read video pack error", err)
				continue
			}
		}
	}()
	return
}
