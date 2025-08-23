package iso8583

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

// PackISO to get Create Message ISO in string
func (m *Message) PackISO() ([]byte, error) {
	m.MapActiveBit()

	m.dataLength = 0

	if m.packager.HasHeader {
		m.dataLength += m.packager.HeaderLength
	}

	if m.MTI == nil {
		return nil, ErrNoMtiToPack
	}
	m.dataLength += len(m.MTI)

	bitmap := make([]byte, 16) // max 128 bits = 16 bytes

	for _, bit := range m.activeBit {
		if bit < 1 || bit > 128 {
			return nil, ErrInvalidBitNumber
		}

		length, err := m.getTotalBitLength(bit)
		if err != nil {
			return nil, err
		}

		m.dataLength += length

		// bit position (ISO8583 is Big Endian, MSB first)
		byteIndex := (bit - 1) / 8
		bitIndex := (bit - 1) % 8
		bitmap[byteIndex] |= 1 << (7 - bitIndex)
	}

	m.firstBitmap = bitmap[:8]
	m.dataLength += len(m.firstBitmap) * 2

	// check second bitmap
	if !bytes.Equal(bitmap[8:], make([]byte, 8)) {
		bitmap[0] |= 0x80 // set first bit to indicate second bitmap
		m.hasSecondBitmap = true
		m.secondBitmap = bitmap[8:16]
		m.dataLength += len(m.secondBitmap) * 2
	}

	result, err := m.processPackIso()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Message) processPackIso() (result []byte, err error) {
	result = make([]byte, 0, m.dataLength)

	if m.packager.HasHeader {
		result = append(result, m.header...)
	}

	result = append(result, m.MTI...)

	// need to convert bitmaps to hex string before packing
	// []byte -> hex string -> convert the string to []byte
	// the size will be double
	result = append(result, []byte(strings.ToUpper(hex.EncodeToString(m.firstBitmap)))...)

	if m.hasSecondBitmap {
		result = append(result, []byte(strings.ToUpper(hex.EncodeToString(m.secondBitmap)))...)
	}

	for _, i := range m.activeBit {
		packager, ok := m.packager.IsoPackagerConfig[i]
		if !ok {
			return nil, fmt.Errorf("packager not found for bit %d", i)
		}
		lenType := packager.Length.Type

		prefixLen := lenType.GetPrefixLen()
		if prefixLen == 0 {
			result = append(result, m.isoMessageMap[i]...)
			continue
		}

		encPrefixLen := lenType.encodeLen(len(m.isoMessageMap[i]), prefixLen)
		result = append(result, encPrefixLen...)
		result = append(result, m.isoMessageMap[i]...)
	}

	return result, nil
}

func formatBitmap(bits []byte) []byte {
	var bitmap uint64
	var buff []byte
	for i, bit := range bits {
		if bit == '1' {
			bitmap |= 1 << (63 - i)
		}
	}
	binary.BigEndian.PutUint64(buff, bitmap)
	return buff
}
