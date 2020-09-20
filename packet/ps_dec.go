package packet

import (
	"github.com/mask-pp/rtp-ps/buffer"
)

/*
https://github.com/videolan/vlc/blob/master/modules/demux/mpeg
*/
type DecPSPackage struct {
	systemClockReferenceBase      uint64
	systemClockReferenceExtension uint64
	programMuxRate                uint32

	VideoStreamType uint32
	AudioStreamType uint32
	buffer.RawBuffer
	VideoPayload []byte
	AudioPayload []byte
}
// data包含 接受到完整一帧数据后，所有的payload, 解析出去后是一阵完整的raw数据
func (dec *DecPSPackage) Read(data []byte) error{
	return dec.decPackHeader(append(data, 0x00, 0x00, 0x01, 0xb9))
}

func (dec *DecPSPackage) clean() {
	dec.systemClockReferenceBase = 0
	dec.systemClockReferenceExtension = 0
	dec.programMuxRate = 0
	dec.VideoPayload = nil
	dec.AudioPayload = nil
}

func (dec *DecPSPackage) decPackHeader(data []byte) error {
	dec.clean()

	// 加载数据
	dec.LoadBuffer(data)

	if startcode, err := dec.Uint32(); err != nil {
		return err
	} else if startcode != StartCodePS {
		return ErrNotFoundStartCode
	}

	if err := dec.Skip(9); err != nil {
		return err
	}

	psl, err := dec.Uint8()
	if err != nil {
		return err
	}
	psl &= 0x07
	if err = dec.Skip(int(psl)); err != nil {
		return err
	}

	for {
		nextStartCode, err := dec.Uint32()
		if err != nil {
			return err
		}

		switch nextStartCode {
		case StartCodeSYS:
			if err := dec.decSystemHeader(); err != nil {
				return err
			}
		case StartCodeMAP:
			if err := dec.decProgramStreamMap(); err != nil {
				return err
			}
		case StartCodeVideo,StartCodeAudio:
			if err := dec.decPESPacket(nextStartCode); err != nil {
				return err
			}
		case HaiKangCode, MEPGProgramEndCode:
			return nil
		}
	}
}

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

func (dec *DecPSPackage) decProgramStreamMap() error {
	psm, err := dec.Uint16()
	if err != nil {
		return err
	}
	//drop psm version infor
	if err = dec.Skip(2); err != nil {
		return err
	}
	psm -= 2

	programStreamInfoLen, err := dec.Uint16()
	if err != nil {
		return err
	}
	if err = dec.Skip(int(programStreamInfoLen)); err != nil {
		return err
	}
	psm -= programStreamInfoLen + 2

	programStreamMapLen, err := dec.Uint16()
	if err != nil {
		return err
	}
	psm -= 2 + programStreamMapLen

	for programStreamMapLen > 0 {
		streamType, err := dec.Uint8()
		if err != nil {
			return err
		}

		elementaryStreamID, err := dec.Uint8()
		if err != nil {
			return err
		}

		if elementaryStreamID >= 0xe0 && elementaryStreamID <= 0xef {
			dec.VideoStreamType = uint32(streamType)
		} else if elementaryStreamID >= 0xc0 && elementaryStreamID <= 0xdf {
			dec.AudioStreamType = uint32(streamType)
		}

		elementaryStreamInfoLength, err := dec.Uint16()
		if err != nil {
			return err
		}
		if err = dec.Skip(int(elementaryStreamInfoLength)); err != nil {
			return err
		}
		programStreamMapLen -= 4 + elementaryStreamInfoLength
	}

	// crc 32
	if psm != 4 {
		return ErrFormatPack
	}
	if err = dec.Skip(4); err != nil {
		return err
	}
	return nil
}

func (dec *DecPSPackage) decPESPacket(t uint32) error {
	payloadlen, err := dec.Uint16()
	if err != nil {
		return err
	}
	if err = dec.Skip(2); err != nil {
		return err
	}

	payloadlen -= 2
	pesHeaderDataLen, err := dec.Uint8()
	if err != nil {
		return err
	}
	payloadlen -= uint16(pesHeaderDataLen) + 1

	if err = dec.Skip(int(pesHeaderDataLen)); err != nil {
		return err
	}

	if payload, err := dec.Bytes(int(payloadlen)); err != nil {
		return err
	} else {
		if StartCodeVideo == t {
			dec.VideoPayload = append(dec.VideoPayload, payload...)
		}else{
			dec.AudioPayload = append(dec.AudioPayload, payload...)
		}
	}
	return nil
}
