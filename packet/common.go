package packet

import "errors"

//
const (
	UDPTransfer        int = 0
	TCPTransferActive  int = 1
	TCPTransferPassive int = 2
	LocalCache         int = 3
)

//
const (
	StreamTypeH264 = 0x1b
	StreamTypeH265 = 0x24
	StreamTypeAAC  = 0x90
)

//
const (
	StreamIDVideo = 0xe0
	StreamIDAudio = 0xc0
)

//
const (
	StartCodePS        = 0x000001ba
	StartCodeSYS       = 0x000001bb
	StartCodeMAP       = 0x000001bc
	StartCodeVideo     = 0x000001e0
	StartCodeAudio     = 0x000001c0
	HaiKangCode        = 0x000001bd
	MEPGProgramEndCode = 0x000001b9
)

//... len limit
const (
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

/*
 This implement from VLC source code
 notes: https://github.com/videolan/vlc/blob/master/modules/mux/mpeg/bits.h
*/

//bitsBuffer bits buffer
type bitsBuffer struct {
	iSize int
	iData int
	iMask uint8
	pData []byte
}

func bitsInit(isize int, buffer []byte) *bitsBuffer {

	bits := &bitsBuffer{
		iSize: isize,
		iData: 0,
		iMask: 0x80,
		pData: buffer,
	}
	if bits.pData == nil {
		bits.pData = make([]byte, isize)
	}
	return bits
}

func bitsAlign(bits *bitsBuffer) {

	if bits.iMask != 0x80 && bits.iData < bits.iSize {
		bits.iMask = 0x80
		bits.iData++
		bits.pData[bits.iData] = 0x00
	}
}
func bitsWrite(bits *bitsBuffer, count int, src uint64) *bitsBuffer {

	for count > 0 {
		count--
		if ((src >> uint(count)) & 0x01) != 0 {
			bits.pData[bits.iData] |= bits.iMask
		} else {
			bits.pData[bits.iData] &= ^bits.iMask
		}
		bits.iMask >>= 1
		if bits.iMask == 0 {
			bits.iData++
			bits.iMask = 0x80
		}
	}

	return bits
}
