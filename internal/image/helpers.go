package image

import "encoding/binary"

func ReadBigEndianInt(data []byte) int {
	if len(data) < 4 {
		return 0
	}
	return int(binary.BigEndian.Uint32(data))
}

func ReadBigEndianShort(data []byte) int {
	if len(data) < 2 {
		return 0
	}
	return int(binary.BigEndian.Uint16(data))
}

func ReadLittleEndianInt(data []byte) int {
	if len(data) < 4 {
		return 0
	}
	return int(binary.LittleEndian.Uint32(data))
}

func ReadLittleEndianShort(data []byte) int {
	if len(data) < 2 {
		return 0
	}
	return int(binary.LittleEndian.Uint16(data))
}
