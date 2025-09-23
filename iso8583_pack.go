package iso8583

import (
	"bytes"
	"fmt"
	"sort"
)

var EmptyBitmap [8]byte
var EmptyMti [4]byte

// PackISO to get Create Message ISO in string
func (m *Message) PackISO() ([]byte, error) {
	sort.Ints(m.activeBits[:m.activeCount])
	// count data len to correctly allocate memory
	dataLength := 0

	if m.packager.HasHeader {
		dataLength += m.packager.HeaderLength
	}

	if m.MTI == EmptyMti {
		return nil, ErrNoMtiToPack
	}
	dataLength += len(m.MTI)

	//bitmap := make([]byte, 16) // max 128 bits = 16 bytes

	bitmap := [16]byte{}

	//for i, v := range m.isoMessageMap {
	for i := 0; i <= m.activeCount; i++ {
		bit := m.activeBits[i]
		if m.isoMessageMap[bit] == nil {
			continue
		}

		length, err := m.getTotalBitLength(bit)
		if err != nil {
			return nil, err
		}

		dataLength += length

		// bit position (ISO8583 is Big Endian, MSB first)
		byteIndex := (bit - 1) / 8
		bitIndex := (bit - 1) % 8

		// use OR operation to keep the already turned on bitmap
		// and turn the current bitmap ON
		// the logic exactly the same as the unpacking
		bitmap[byteIndex] |= 1 << (7 - bitIndex)
	}

	dataLength += BitmapLength

	// check second bitmap
	if !bytes.Equal(bitmap[8:], EmptyBitmap[:]) {
		// set first bit to indicate second bitmap is on
		// 0x80 is 10000000, and use OR operation
		bitmap[0] |= 0x80
		dataLength += BitmapLength
	}

	result, err := m.processPackIso(bitmap, dataLength)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Message) processPackIso(bitmap [16]byte, dataLength int) ([]byte, error) {

	if cap(m.byteData) < dataLength {
		m.byteData = make([]byte, dataLength) // exact size
	} else {
		m.byteData = m.byteData[:dataLength]
	}

	// Write offset instead of appending
	pos := 0

	// Header
	if m.packager.HasHeader {
		pos += copy(m.byteData[pos:], m.header)
	}

	// MTI
	pos += copy(m.byteData[pos:], m.MTI[:])

	// --- First bitmap directly into m.byteData ---
	encodeHexUpper(m.byteData[pos:], bitmap[:8])
	pos += BitmapLength

	// --- Second bitmap if exists ---
	if bitmap[0]&0x80 != 0 {
		encodeHexUpper(m.byteData[pos:], bitmap[8:])
		pos += BitmapLength
	}

	// --- Fields ---
	for i := 2; i <= 128; i++ {
		if m.isoMessageMap[i] == nil {
			continue
		}
		packager := m.packager.IsoPackagerConfig[i]
		if packager.Type == "" {
			return nil, fmt.Errorf("packager not found for bit %d", i)
		}

		prefixLen := packager.Length.Type.GetPrefixLen()

		if prefixLen > 0 {
			packager.Length.Type.encodeLenInto(len(m.isoMessageMap[i]), m.byteData[pos:pos+prefixLen])
			pos += prefixLen
		}

		pos += copy(m.byteData[pos:], m.isoMessageMap[i])
	}

	// Trim to actual used size
	final := m.byteData[:pos]

	return final, nil
}

var hexUpper = "0123456789ABCDEF"

func encodeHexUpper(dst []byte, src []byte) {
	for i, b := range src {
		dst[i*2] = hexUpper[b>>4]
		dst[i*2+1] = hexUpper[b&0x0f]
	}
}
