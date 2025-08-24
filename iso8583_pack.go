package iso8583

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// PackISO to get Create Message ISO in string
func (m *Message) PackISO() ([]byte, error) {
	m.mapActiveBit()

	// count data len to correctly allocate memory
	m.dataLength = 0

	if m.packager.HasHeader {
		m.dataLength += m.packager.HeaderLength
	}

	if m.MTI == nil {
		return nil, ErrNoMtiToPack
	}
	m.dataLength += len(m.MTI)

	bitmap := make([]byte, 16) // max 128 bits = 16 bytes
	m.firstBitmap = nil
	m.secondBitmap = nil
	for _, bit := range m.activeBit {
		if bit < 2 || bit > 128 {
			errMsg := fmt.Errorf("invalid bit %d", bit)
			return nil, errors.Join(ErrInvalidBitNumber, errMsg)
		}

		length, err := m.getTotalBitLength(bit)
		if err != nil {
			return nil, err
		}

		m.dataLength += length

		// bit position (ISO8583 is Big Endian, MSB first)
		byteIndex := (bit - 1) / 8
		bitIndex := (bit - 1) % 8

		// use OR operation to keep the already turned on bitmap
		// and turn the current bitmap ON
		// the logic exactly the same as the unpacking
		bitmap[byteIndex] |= 1 << (7 - bitIndex)
	}

	// first bitmap is 8 byte, later converted to hex string to be 16 digit hex string
	m.firstBitmap = bitmap[:8]
	m.dataLength += BitmapLength

	// check second bitmap
	if !bytes.Equal(bitmap[8:], make([]byte, 8)) {
		// set first bit to indicate second bitmap is on
		// 0x80 is 10000000, and use OR operation
		m.firstBitmap[0] |= 0x80
		m.secondBitmap = bitmap[8:16]
		m.dataLength += BitmapLength
	}

	result, err := m.processPackIso()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Message) mapActiveBit() {
	for bitNum := range m.isoMessageMap {
		if bitNum == 0 {
			continue
		}
		m.activeBit = append(m.activeBit, bitNum)
	}

	// sort active bitmap to correctly handling the packing
	sort.Ints(m.activeBit)
}

func (m *Message) processPackIso() (result []byte, err error) {
	result = make([]byte, 0, m.dataLength)

	if m.packager.HasHeader {
		result = append(result, m.header...)
	}

	result = append(result, m.MTI...)

	// need to convert bitmaps to hex string before packing
	// []byte -> uppercase hex string -> convert the string to []byte
	// the size will be doubled
	result = append(result, []byte(strings.ToUpper(hex.EncodeToString(m.firstBitmap)))...)

	if m.secondBitmap != nil {
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
