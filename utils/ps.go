package utils

import (
	"errors"

	"github.com/Monibuca/utils/v3"
)

//
const (
	UDPTransfer        int = 0
	TCPTransferActive  int = 1
	TCPTransferPassive int = 2
	LocalCache         int = 3

	StreamTypeH264 = 0x1b
	StreamTypeH265 = 0x24
	G711A          = 0x90 //PCMA
	G7221AUDIOTYPE = 0x92
	G7231AUDIOTYPE = 0x93
	G729AUDIOTYPE  = 0x99

	StreamIDVideo = 0xe0
	StreamIDAudio = 0xc0

	StartCodePS        = 0x000001ba
	StartCodeSYS       = 0x000001bb
	StartCodeMAP       = 0x000001bc
	StartCodeVideo     = 0x000001e0
	StartCodeAudio     = 0x000001c0
	HaiKangCode        = 0x000001bd
	MEPGProgramEndCode = 0x000001b9

	RTPHeaderLength    int = 12
	PSHeaderLength     int = 14
	SystemHeaderLength int = 18
	MAPHeaderLength    int = 24
	PESHeaderLength    int = 19
	RtpLoadLength      int = 1460
	PESLoadLength      int = 0xFFFF
	MAXFrameLen        int = 1024 * 1024 * 2
)

var (
	ErrNotFoundStartCode = errors.New("not found the need start code flag")
	ErrMarkerBit         = errors.New("marker bit value error")
	ErrFormatPack        = errors.New("not package standard")
	ErrParsePakcet       = errors.New("parse ps packet error")
)

type Pusher interface {
	PushVideo(uint32, uint32, []byte)
	PushAudio(uint32, []byte)
}

/*
 This implement from VLC source code
 notes: https://github.com/videolan/vlc/blob/master/modules/mux/mpeg/bits.h
*/

//bitsBuffer bits buffer
// type bitsBuffer struct {
// 	iSize int
// 	iData int
// 	iMask uint8
// 	pData []byte
// }

// func bitsInit(isize int, buffer []byte) *bitsBuffer {

// 	bits := &bitsBuffer{
// 		iSize: isize,
// 		iData: 0,
// 		iMask: 0x80,
// 		pData: buffer,
// 	}
// 	if bits.pData == nil {
// 		bits.pData = make([]byte, isize)
// 	}
// 	return bits
// }

// func bitsAlign(bits *bitsBuffer) {

// 	if bits.iMask != 0x80 && bits.iData < bits.iSize {
// 		bits.iMask = 0x80
// 		bits.iData++
// 		bits.pData[bits.iData] = 0x00
// 	}
// }
// func bitsWrite(bits *bitsBuffer, count int, src uint64) *bitsBuffer {

// 	for count > 0 {
// 		count--
// 		if ((src >> uint(count)) & 0x01) != 0 {
// 			bits.pData[bits.iData] |= bits.iMask
// 		} else {
// 			bits.pData[bits.iData] &= ^bits.iMask
// 		}
// 		bits.iMask >>= 1
// 		if bits.iMask == 0 {
// 			bits.iData++
// 			bits.iMask = 0x80
// 		}
// 	}

// 	return bits
// }

/*
https://github.com/videolan/vlc/blob/master/modules/demux/mpeg
*/
type DecPSPackage struct {
	systemClockReferenceBase      uint64
	systemClockReferenceExtension uint64
	programMuxRate                uint32

	VideoStreamType uint32
	AudioStreamType uint32
	IOBuffer
	Payload []byte
	PTS     uint32
	DTS     uint32
}

func (dec *DecPSPackage) clean() {
	dec.systemClockReferenceBase = 0
	dec.systemClockReferenceExtension = 0
	dec.programMuxRate = 0
	dec.Payload = nil
	dec.PTS = 0
	dec.DTS = 0
}

func (dec *DecPSPackage) ReadPayload() (payload []byte, err error) {
	payloadlen, err := dec.Uint16()
	if err != nil {
		return
	}
	return dec.ReadN(int(payloadlen))
}

func (dec *DecPSPackage) Read(data []byte, ts uint32, pusher Pusher) error {
	var nextStartCode uint32
	dec.clean()
	dec.Reset()
	// 加载数据
	dec.Write(data)

	if err := dec.Skip(5); err != nil {
		return err
	}

	psl, err := dec.ReadByte()
	if err != nil {
		return err
	}
	psl &= 0x07
	if err = dec.Skip(int(psl)); err != nil {
		return err
	}
	var video []byte
	var videoTs, videoCts uint32
	defer func() {
		if video != nil {
			pusher.PushVideo(videoTs, videoCts, video)
			video = nil
		}
		if nextStartCode == StartCodePS {
			err = dec.Read(dec.Bytes(), ts, pusher)
		}
	}()
	for err == nil {
		nextStartCode, err = dec.Uint32()
		if err != nil {
			return err
		}
		switch nextStartCode {
		case StartCodeSYS:
			dec.ReadPayload()
			//err = dec.decSystemHeader()
		case StartCodeMAP:
			err = dec.decProgramStreamMap()
		case StartCodeVideo:
			var cts uint32
			if err = dec.decPESPacket(); err == nil {
				if video == nil {
					if dec.PTS == 0 {
						dec.PTS = ts
					}
					if dec.DTS != 0 {
						cts = dec.PTS - dec.DTS
					} else {
						dec.DTS = dec.PTS
					}
					videoTs = dec.DTS / 90
					videoCts = cts / 90
				}
				video = append(video, dec.Payload...)
			} else {
				utils.Println("video", err)
			}
		case StartCodeAudio:
			if err = dec.decPESPacket(); err == nil {
				if dec.PTS != 0 {
					pusher.PushAudio(dec.PTS/8, dec.Payload)
				} else {
					pusher.PushAudio(ts/8, dec.Payload)
				}
			} else {
				utils.Println("audio", err)
			}
		case StartCodePS:
			return nil
		default:
			dec.ReadPayload()
		}
	}
	return err
}

/*
func (dec *DecPSPackage) decSystemHeader() error {
	syslens, err := dec.Uint16()
	if err != nil {
		return err
	}
	// drop rate video audio bound and lock flag
	syslens -= 6
	if err = dec.Skip(6); err != nil {
		return err
	}

	// ONE WAY: do not to parse the stream  and skip the buffer
	//br.Skip(syslen * 8)

	// TWO WAY: parse every stream info
	for syslens > 0 {
		if nextbits, err := dec.Uint8(); err != nil {
			return err
		} else if (nextbits&0x80)>>7 != 1 {
			break
		}
		if err = dec.Skip(2); err != nil {
			return err
		}
		syslens -= 3
	}
	return nil
}
*/
func (dec *DecPSPackage) decProgramStreamMap() error {
	psm, err := dec.ReadPayload()
	if err != nil {
		return err
	}
	index := 2
	programStreamInfoLen := utils.BigEndian.Uint16(psm[index:])
	index += 2
	index += int(programStreamInfoLen)
	programStreamMapLen := utils.BigEndian.Uint16(psm[index:])
	index += 2
	for programStreamMapLen > 0 {
		streamType := psm[index]
		index++
		elementaryStreamID := psm[index]
		index++
		if elementaryStreamID >= 0xe0 && elementaryStreamID <= 0xef {
			dec.VideoStreamType = uint32(streamType)
		} else if elementaryStreamID >= 0xc0 && elementaryStreamID <= 0xdf {
			dec.AudioStreamType = uint32(streamType)
		}

		elementaryStreamInfoLength := utils.BigEndian.Uint16(psm[index:])
		index += 2
		index += int(elementaryStreamInfoLength)
		programStreamMapLen -= 4 + elementaryStreamInfoLength
	}
	return nil
}

func (dec *DecPSPackage) decPESPacket() error {
	payload, err := dec.ReadPayload()
	if err != nil {
		return err
	}
	if len(payload) < 4 {
		return errors.New("not enough data")
	}
	//data_alignment_indicator := (payload[0]&0b0001_0000)>>4 == 1
	flag := payload[1]
	ptsFlag := flag>>7 == 1
	dtsFlag := (flag&0b0100_0000)>>6 == 1
	var pts, dts uint32
	pesHeaderDataLen := payload[2]
	payload = payload[3:]
	extraData := payload[:pesHeaderDataLen]
	if ptsFlag && len(extraData) > 4 {
		pts = uint32(extraData[0]&0b0000_1110) << 29
		pts += uint32(extraData[1]) << 22
		pts += uint32(extraData[2]&0b1111_1110) << 14
		pts += uint32(extraData[3]) << 7
		pts += uint32(extraData[4]) >> 1
		if dtsFlag && len(extraData) > 9 {
			dts = uint32(extraData[5]&0b0000_1110) << 29
			dts += uint32(extraData[6]) << 22
			dts += uint32(extraData[7]&0b1111_1110) << 14
			dts += uint32(extraData[8]) << 7
			dts += uint32(extraData[9]) >> 1
		}
	}
	dec.PTS = pts
	dec.DTS = dts
	dec.Payload = payload[pesHeaderDataLen:]
	return err
}
