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
	for i := 0; i < m.activeCount; i++ {
		bit := m.activeBits[i]

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

func (m *Message) processLength(bitNum int) {

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
	for i := 0; i < m.activeCount; i++ {
		bitNum := m.activeBits[i]
		prefixLen := m.packager.PrefixLengths[bitNum]
		value := m.isoMessageMap[bitNum]

		if prefixLen == 0 {
			return nil, fmt.Errorf("packager not found for bit %d", bitNum)
		}

		if prefixLen != FixedLength {
			length := len(value)
			if length > m.packager.MaxLengths[bitNum] {
				return nil, fmt.Errorf(
					"invalid bit length for bit %d: max %d, got %d",
					bitNum,
					m.packager.MaxLengths[bitNum],
					length,
				)
			}
			encodeLenInto(length, prefixLen, m.byteData[pos:pos+prefixLen])
			pos += prefixLen
		}

		pos += copy(m.byteData[pos:], value)
	}

	// Trim to actual used size
	final := m.byteData[:pos]

	return final, nil
}

// encodeHexUpper encodes the source byte slice into hexadecimal representation
// and stores the result in the destination byte slice.
//
// The function uses uppercase hexadecimal characters (0-9, A-F) for encoding.
//
// Parameters:
//   - dst: the destination byte slice where the encoded hexadecimal values will be stored
//   - src: the source byte slice containing the original data to be encoded
//
// Note: The destination slice must have sufficient capacity (at least 2*len(src)) to
// accommodate the encoded result, as each byte in src produces two hexadecimal characters.
func encodeHexUpper(dst []byte, src []byte) {
	for i, b := range src {
		dst[i*2], dst[i*2+1] = hexTable[b][0], hexTable[b][1]
	}
}
