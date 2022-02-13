package fastimage

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var errInsufficientData = errors.New("Insufficient data")

type imageWEBP struct{}

func (p imageWEBP) Type() ImageType {
	return WEBP
}

func (p imageWEBP) Detect(buffer []byte) bool {
	firstTwoBytes := string(buffer[:2])
	return firstTwoBytes == "RI"
}

func (p imageWEBP) GetSize(buffer []byte) (*ImageSize, error) {
	if len(buffer) < 12 {
		return nil, errInsufficientData
	}

	buffer = buffer[12:]
	for ; len(buffer) >= 8; {
		var l uint32
		lr := bytes.NewReader(buffer[4:8])
		binary.Read(lr, binary.LittleEndian, &l)
		if l % 2 == 1 {
			l++
		}
		chunkType := string(buffer[:4])
		switch chunkType {
		case "VP8X":
			if len(buffer) < 18 {
				return nil, errInsufficientData
			}
			var w = (uint32(buffer[12]) | (uint32(buffer[13]) << 8)) | (uint32(buffer[14]) << 16) + 1
			var h = (uint32(buffer[15]) | (uint32(buffer[16]) << 8)) | (uint32(buffer[17]) << 16) + 1
			return &ImageSize{
				Width:  w,
				Height: h,
			}, nil
		case "VP8 ":
			if len(buffer) < 18 {
				return nil, errInsufficientData
			}
			var w = (uint32(buffer[14]) | (uint32(buffer[15]) << 8)) & 0x3fff
			var h = (uint32(buffer[16]) | (uint32(buffer[17]) << 8)) & 0x3fff
			return &ImageSize{
				Width:  w,
				Height: h,
			}, nil
		case "VP8L":
			if len(buffer) <= 12 {
				return nil, errInsufficientData
			}
			var w = uint32(buffer[9]) + ((uint32(buffer[10]) & 0b00111111) << 8) + 1
			var h = (uint32(buffer[10] & 0b11000000) >> 6) + (uint32(buffer[11]) << 2) + (uint32(buffer[12]) & 0b00001111 << 10) + 1
			return &ImageSize{
				Width:  w,
				Height: h,
			}, nil
		}
		if len(buffer) < 8+int(l) {
			break
		}
		buffer = buffer[8+int(l):]
	}
	return nil, errInsufficientData
}

func init() {
	register(&imageWEBP{})
}
