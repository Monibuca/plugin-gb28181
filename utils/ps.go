package utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	UDPTransfer        int = 0
	TCPTransferActive  int = 1
	TCPTransferPassive int = 2
	LocalCache         int = 3

	StartCodePS        = 0x000001ba
	StartCodeSYS       = 0x000001bb
	StartCodeMAP       = 0x000001bc
	StartCodeVideo     = 0x000001e0
	StartCodeAudio     = 0x000001c0
	PrivateStreamCode  = 0x000001bd
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
	PrintDump(string)
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
	Payload     []byte
	videoBuffer []byte
	audioBuffer []byte
	PTS         uint32
	DTS         uint32
	Pusher
}

func NewDecPSPackage(p Pusher) *DecPSPackage {
	p.PrintDump("<tr><td>")
	return &DecPSPackage{
		Pusher: p,
	}
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

// Drop 由于丢包引起的必须丢弃的数据
func (dec *DecPSPackage) Drop() {
	dec.Reset()
	dec.videoBuffer = nil
	dec.audioBuffer = nil
	dec.Payload = nil
}

func (dec *DecPSPackage) Feed(ps []byte) (err error) {
	if ps[0] == 0 && ps[1] == 0 && ps[2] == 1 {
		defer dec.Write(ps)
		if dec.Len() >= 4 {
			//说明需要处理PS包，处理完后，清空缓存
			defer dec.Reset()
		} else {
			return
		}
	} else {
		// 说明是中间数据，直接写入缓存，否则数据不合法需要丢弃
		if dec.Len() > 0 {
			dec.Write(ps)
		}
		return nil
	}
	for dec.Len() >= 4 {
		code, _ := dec.Uint32()
		// println("code:", code)
		switch code {
		case StartCodePS:
			dec.PrintDump("</td></tr><tr><td>")
			if len(dec.audioBuffer) > 0 {
				dec.PushAudio(dec.PTS, dec.audioBuffer)
				dec.audioBuffer = nil
			}
			if len(dec.videoBuffer) > 0 {
				dec.PushVideo(dec.PTS, dec.DTS, dec.videoBuffer)
				dec.videoBuffer = nil
			}
			if err := dec.Skip(9); err != nil {
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
		case StartCodeSYS:
			dec.PrintDump("</td><td>[sys]")
			dec.ReadPayload()
		case StartCodeMAP:
			dec.decProgramStreamMap()
			dec.PrintDump("</td><td>[map]")
		case StartCodeVideo:
			if dec.videoBuffer == nil {
				dec.PrintDump("</td><td>")
			}
			if err = dec.decPESPacket(); err == nil {
				dec.videoBuffer = append(dec.videoBuffer, dec.Payload...)
			} else {
				fmt.Println("video", err)
			}
			dec.PrintDump("[video]")
		case StartCodeAudio:
			if dec.audioBuffer == nil {
				dec.PrintDump("</td><td>")
			}
			if err = dec.decPESPacket(); err == nil {
				dec.audioBuffer = append(dec.audioBuffer, dec.Payload...)
				dec.PrintDump("[audio]")
			} else {
				fmt.Println("audio", err)
			}
		case PrivateStreamCode:
			dec.ReadPayload()
			dec.PrintDump("</td></tr><tr><td>[ac3]")
		case MEPGProgramEndCode:
			dec.PrintDump("</td></tr>")
			return io.EOF
		default:
			fmt.Println("unknow code", code)
			return ErrParsePakcet
		}
	}
	return nil
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
	l := len(psm)
	index := 2
	programStreamInfoLen := binary.BigEndian.Uint16(psm[index:])
	index += 2
	index += int(programStreamInfoLen)
	programStreamMapLen := binary.BigEndian.Uint16(psm[index:])
	index += 2
	for programStreamMapLen > 0 {
		if l <= index+1 {
			break
		}
		streamType := psm[index]
		index++
		elementaryStreamID := psm[index]
		index++
		if elementaryStreamID >= 0xe0 && elementaryStreamID <= 0xef {
			dec.VideoStreamType = uint32(streamType)
		} else if elementaryStreamID >= 0xc0 && elementaryStreamID <= 0xdf {
			dec.AudioStreamType = uint32(streamType)
		}
		if l <= index+1 {
			break
		}
		elementaryStreamInfoLength := binary.BigEndian.Uint16(psm[index:])
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
		dec.PTS = pts
		if dtsFlag && len(extraData) > 9 {
			dts = uint32(extraData[5]&0b0000_1110) << 29
			dts += uint32(extraData[6]) << 22
			dts += uint32(extraData[7]&0b1111_1110) << 14
			dts += uint32(extraData[8]) << 7
			dts += uint32(extraData[9]) >> 1
			dec.DTS = dts
		}
	}
	dec.Payload = payload[pesHeaderDataLen:]
	return err
}
