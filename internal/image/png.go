// Package image provides image decoding and DPI extraction for formats supported by OOXML.
package image

import (
	"encoding/binary"
	"math"
)

func readPNGDPI(data []byte) DPI {
	dpi := DPI{Horizontal: DefaultDPI, Vertical: DefaultDPI}
	if len(data) < 8 {
		return dpi
	}

	offset := 8
	for offset+8 <= len(data) {
		chunkLen := int(binary.BigEndian.Uint32(data[offset : offset+4]))
		chunkType := string(data[offset+4 : offset+8])

		if chunkType == "IEND" {
			break
		}

		if chunkType == "pHYs" && offset+8+chunkLen <= len(data) {
			chunkData := data[offset+8 : offset+8+chunkLen]
			if len(chunkData) >= 9 {
				ppuX := binary.BigEndian.Uint32(chunkData[0:4])
				ppuY := binary.BigEndian.Uint32(chunkData[4:8])
				unit := chunkData[8]
				if unit == 1 && ppuX > 0 && ppuY > 0 {
					dpi.Horizontal = int(math.Round(float64(ppuX) * 0.0254))
					dpi.Vertical = int(math.Round(float64(ppuY) * 0.0254))
				}
			}
			break
		}

		offset += 12 + chunkLen
	}

	return dpi
}
