package image

import (
	"encoding/binary"
	"math"
)

func readJPEGDPI(data []byte) DPI {
	dpi := DPI{Horizontal: DefaultDPI, Vertical: DefaultDPI}
	offset := 2

	for offset+4 <= len(data) {
		if data[offset] != 0xFF {
			offset++
			continue
		}

		marker := data[offset+1]

		if marker == 0x01 || marker == 0xD0 || marker == 0xD1 ||
			marker == 0xD2 || marker == 0xD3 || marker == 0xD4 ||
			marker == 0xD5 || marker == 0xD6 || marker == 0xD7 ||
			marker == 0xD8 || marker == 0xD9 {
			offset += 2
			continue
		}

		if offset+4 > len(data) {
			break
		}

		segLen := int(binary.BigEndian.Uint16(data[offset+2 : offset+4]))
		if segLen < 2 || offset+2+segLen > len(data) {
			break
		}

		segStart := offset + 2
		segEnd := segStart + segLen

		switch marker {
		case 0xE0:
			result := parseJFIF(data[segStart:segEnd])
			if result.Horizontal != DefaultDPI {
				return result
			}
		case 0xE1:
			result := parseExifAPP1(data[segStart:segEnd])
			if result.Horizontal != DefaultDPI {
				return result
			}
		}

		offset = segEnd

		if marker == 0xDA {
			break
		}
	}

	return dpi
}

func parseJFIF(seg []byte) DPI {
	dpi := DPI{Horizontal: DefaultDPI, Vertical: DefaultDPI}
	if len(seg) < 14 {
		return dpi
	}
	if string(seg[2:7]) != "JFIF\x00" {
		return dpi
	}

	densityUnits := seg[9]
	xDensity := int(binary.BigEndian.Uint16(seg[10:12]))
	yDensity := int(binary.BigEndian.Uint16(seg[12:14]))

	switch densityUnits {
	case 1:
		dpi.Horizontal = xDensity
		dpi.Vertical = yDensity
	case 2:
		dpi.Horizontal = int(math.Round(float64(xDensity) * 2.54))
		dpi.Vertical = int(math.Round(float64(yDensity) * 2.54))
	}

	return dpi
}

func parseExifAPP1(seg []byte) DPI {
	dpi := DPI{Horizontal: DefaultDPI, Vertical: DefaultDPI}
	if len(seg) < 14 {
		return dpi
	}
	if string(seg[2:8]) != "Exif\x00\x00" {
		return dpi
	}

	tiffData := seg[8:]
	if len(tiffData) < 8 {
		return dpi
	}

	var bo binary.ByteOrder
	switch string(tiffData[0:2]) {
	case "II":
		bo = binary.LittleEndian
	case "MM":
		bo = binary.BigEndian
	default:
		return dpi
	}

	if bo.Uint16(tiffData[2:4]) != 0x002A {
		return dpi
	}

	ifdOffset := int(bo.Uint32(tiffData[4:8]))
	if ifdOffset+2 > len(tiffData) {
		return dpi
	}

	numEntries := int(bo.Uint16(tiffData[ifdOffset : ifdOffset+2]))
	ifdOffset += 2

	var xResolution, yResolution float64
	resolutionUnit := 2

	for i := 0; i < numEntries; i++ {
		entryOff := ifdOffset + i*12
		if entryOff+12 > len(tiffData) {
			break
		}

		tag := bo.Uint16(tiffData[entryOff : entryOff+2])
		fieldType := bo.Uint16(tiffData[entryOff+2 : entryOff+4])
		_ = bo.Uint32(tiffData[entryOff+4 : entryOff+8])
		valueOffset := bo.Uint32(tiffData[entryOff+8 : entryOff+12])

		switch tag {
		case 0x011A:
			xResolution = parseTiffRational(bo, tiffData, int(valueOffset))
		case 0x011B:
			yResolution = parseTiffRational(bo, tiffData, int(valueOffset))
		case 0x0128:
			if fieldType == 3 {
				resolutionUnit = int(valueOffset)
			}
		}
	}

	if xResolution > 0 && yResolution > 0 {
		switch resolutionUnit {
		case 2:
			dpi.Horizontal = int(math.Round(xResolution))
			dpi.Vertical = int(math.Round(yResolution))
		case 3:
			dpi.Horizontal = int(math.Round(xResolution * 2.54))
			dpi.Vertical = int(math.Round(yResolution * 2.54))
		}
	}

	return dpi
}

func parseTiffRational(bo binary.ByteOrder, data []byte, offset int) float64 {
	if offset+8 > len(data) {
		return 0
	}
	num := float64(bo.Uint32(data[offset : offset+4]))
	den := float64(bo.Uint32(data[offset+4 : offset+8]))
	if den == 0 {
		return 0
	}
	return num / den
}
