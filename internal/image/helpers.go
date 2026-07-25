// Package image provides image decoding and DPI extraction for formats supported by OOXML.
package image

import "encoding/binary"

// ReadBigEndianInt reads a 4-byte big-endian integer from data. Returns 0 if fewer than 4 bytes are available.
func ReadBigEndianInt(data []byte) int {
	if len(data) < 4 {
		return 0
	}
	return int(binary.BigEndian.Uint32(data))
}

// ReadBigEndianShort reads a 2-byte big-endian integer from data. Returns 0 if fewer than 2 bytes are available.
func ReadBigEndianShort(data []byte) int {
	if len(data) < 2 {
		return 0
	}
	return int(binary.BigEndian.Uint16(data))
}

// ReadLittleEndianInt reads a 4-byte little-endian integer from data. Returns 0 if fewer than 4 bytes are available.
func ReadLittleEndianInt(data []byte) int {
	if len(data) < 4 {
		return 0
	}
	return int(binary.LittleEndian.Uint32(data))
}

// ReadLittleEndianShort reads a 2-byte little-endian integer from data. Returns 0 if fewer than 2 bytes are available.
func ReadLittleEndianShort(data []byte) int {
	if len(data) < 2 {
		return 0
	}
	return int(binary.LittleEndian.Uint16(data))
}
