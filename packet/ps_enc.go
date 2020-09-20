package packet

type encPSPacket struct {
	crc32 uint64
}

/*
	https://github.com/videolan/vlc/tree/master/modules/mux/mpeg
*/

func (enc *encPSPacket) encPackHeader(pts uint64) []byte {
	pack := make([]byte, PSHeaderLength)
	bits := bitsInit(PSHeaderLength, pack)
	bitsWrite(bits, 32, StartCodePS)
	bitsWrite(bits, 2, 0x01)
	bitsWrite(bits, 3, (pts>>30)&0x07)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 15, (pts>>15)&0x7fff)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 15, (pts & 0x7fff))
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 9, 0)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 22, 255&0x3fffff)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 5, 0x1f)
	bitsWrite(bits, 3, 0)
	return bits.pData
}

func (enc *encPSPacket) encSystemHeader(data []byte, vrates, arates int) []byte {
	pack := make([]byte, SystemHeaderLength)
	bits := bitsInit(SystemHeaderLength, pack)
	bitsWrite(bits, 32, StartCodeSYS)
	bitsWrite(bits, 16, uint64(SystemHeaderLength-6))
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 22, 40960)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 6, 1)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 5, 1)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 7, 0xff)

	// video stream bound
	bitsWrite(bits, 8, 0xe0)
	bitsWrite(bits, 2, 3)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 13, uint64(vrates))

	// audio stream bound
	bitsWrite(bits, 8, 0xc0) // 0xc0 音频
	bitsWrite(bits, 2, 0x03)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 13, uint64(arates))

	return append(data, bits.pData...)
}

func (enc *encPSPacket) encProgramStreamMap(data []byte) []byte {

	pack := make([]byte, MAPHeaderLength)
	bits := bitsInit(MAPHeaderLength, pack)
	bitsWrite(bits, 32, StartCodeMAP)
	bitsWrite(bits, 16, uint64(MAPHeaderLength-6))
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 2, 0xf)
	bitsWrite(bits, 5, 0)
	bitsWrite(bits, 7, 0xff)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 16, 0)
	bitsWrite(bits, 16, 8)

	//video
	bitsWrite(bits, 8, StreamTypeH264)
	bitsWrite(bits, 8, StreamIDVideo)
	bitsWrite(bits, 16, 0)
	// audio
	bitsWrite(bits, 8, StreamTypeAAC) // 0x90 G711
	bitsWrite(bits, 8, StreamIDAudio) // 0x0c0 音频取值（0xc0-0xdf），通常为0xc0
	bitsWrite(bits, 16, 0)

	bitsWrite(bits, 8, enc.crc32>>24) // CRC_32 : (32) CRC 32字段
	bitsWrite(bits, 8, (enc.crc32>>16)&0xFF)
	bitsWrite(bits, 8, (enc.crc32>>8)&0xFF)
	bitsWrite(bits, 8, enc.crc32&0xFF)
	return append(data, bits.pData...)
}

func (enc *encPSPacket) encPESPacket(data []byte, streamtype int, payloadlen int, pts, dts uint64) []byte {

	pack := make([]byte, PESHeaderLength)
	bits := bitsInit(PESHeaderLength, pack)
	bitsWrite(bits, 24, 0x01)
	bitsWrite(bits, 8, uint64(streamtype))

	bitsWrite(bits, 16, uint64(payloadlen+13))
	bitsWrite(bits, 2, 2)
	bitsWrite(bits, 2, 0)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)

	bitsWrite(bits, 2, 0x03)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)

	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)
	bitsWrite(bits, 1, 0)

	bitsWrite(bits, 8, 10)
	// pts dts //
	bitsWrite(bits, 4, 3)
	bitsWrite(bits, 3, (pts>>30)&0x07)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 15, (pts>>15)&0x7fff)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 15, pts&0x7fff)
	bitsWrite(bits, 1, 1)

	bitsWrite(bits, 4, 1)
	bitsWrite(bits, 3, (dts>>30)&0x07)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 15, (dts>>15)&0x7fff)
	bitsWrite(bits, 1, 1)
	bitsWrite(bits, 15, dts&0x7fff)
	bitsWrite(bits, 1, 1)
	return append(bits.pData, data...)
}
