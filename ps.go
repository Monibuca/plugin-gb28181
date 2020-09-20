package gb28181

import (
	. "github.com/Monibuca/engine/v2"
	"github.com/Monibuca/engine/v2/avformat"
	"github.com/Monibuca/engine/v2/util"

	"github.com/Monibuca/plugin-gb28181/packet"
	"github.com/Monibuca/plugin-gb28181/transaction"
	rtp "github.com/Monibuca/plugin-rtp"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	G711Mu  = 0x90
	G7221AUDIOTYPE = 0x92
	G7231AUDIOTYPE = 0x93
	G729AUDIOTYPE  = 0x99
)

func onPublish(device *transaction.Device) (port int) {
	rtpPublisher := new(rtp.RTP)
	if !rtpPublisher.Publish("gb28181/" + device.ID) {
		return
	}
	rtpPublisher.Type = "GB28181"
	var parser packet.DecPSPackage
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
		Printf("udp server start listen video port[%d]", port)
		defer Printf("udp server stop listen video port[%d]", port)
		timer := time.Unix(0, 0)
		var psPacket []byte
		var psRtp rtp.RTPPack
		for {
			if n, _, err := conn.ReadFromUDP(bufUDP); err == nil {
				elapsed := time.Now().Sub(timer)
				if elapsed >= 30*time.Second {
					Printf("Package recv from VConn.len:%d\n", n)
					timer = time.Now()
				}
				if err := psRtp.Unmarshal(bufUDP[:n]); err != nil {
					Println(err)
					continue
				}
				if len(psRtp.Payload) >= 4 && util.BigEndian.Uint32(psRtp.Payload) == packet.StartCodePS {
					if psPacket != nil{
						if err := parser.Read(psPacket); err == nil {
							for _, payload := range avformat.SplitH264(parser.VideoPayload) {
								rtpPublisher.WriteNALU(psRtp.Timestamp, payload)
							}
							if parser.AudioPayload != nil{
								//TODO: 需要增加一个字节的头
								//rtpPublisher.PushAudio(psRtp.Timestamp, parser.AudioPayload)
							}
						} else {
							Print(err)
						}
						psPacket = nil
					}
					psPacket = append(psPacket, psRtp.Payload...)
				} else if psPacket != nil {
					psPacket = append(psPacket, psRtp.Payload...)
				}
			} else {
				Println("udp server read video pack error", err)
				continue
			}
		}
	}()
	return
}
