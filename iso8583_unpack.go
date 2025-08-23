package iso8583

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
)

// Unpack to Single Data Element
func (m *Message) Unpack(b []byte) error {
	m.byteData = b
	m.cursor = 0

	data := m.byteData
	if len(data) < 4 {
		return ErrInsufficientDataMti
	}
	// check ISO header
	if m.packager.HasHeader || bytes.Equal(data[:3], isoHeader) {
		m.header = data[:m.packager.HeaderLength]
		m.cursor += m.packager.HeaderLength
	}

	if len(data[m.cursor:]) < 4 {
		return ErrInsufficientDataMti
	}

	// check MTI
	mti := MTITypeByte(data[0:4])
	switch {
	case mti.Equal(MTICardProcessingRequestByte),
		mti.Equal(MTICardProcessingResponseByte),
		mti.Equal(MTIFinancialRequestByte),
		mti.Equal(MTIFinancialResponseByte),
		mti.Equal(MTIReversalRequestByte),
		mti.Equal(MTIReversalResponseByte),
		mti.Equal(MTIRepeatedReversalRequestByte),
		mti.Equal(MTINMMRequestByte),
		mti.Equal(MTINMMResponseByte):
		m.MTI = mti
	default:
		return ErrNotDefaultMti
	}
	m.cursor += 4

	if len(data[m.cursor:]) < 16 {
		return ErrInsufficientDataFirstBitmap
	}

	if m.isoMessageMap == nil {
		m.isoMessageMap = make(map[int][]byte)
	}

	if err := m.parseBitmap(); err != nil {
		return err
	}

	return nil
}

func (m *Message) parseBitmap() (err error) {
	startBit := 0
	var bitmap []byte
	var asciiHex []byte

	if m.firstBitmap == nil {
		asciiHex = m.byteData[m.cursor : m.cursor+16]
		bitmap, err = hex.DecodeString(string(asciiHex))
		if err != nil {
			return ErrInvalidBitMap
		}
		m.firstBitmap = bitmap
		m.cursor += 16
	}
	if m.hasSecondBitmap {
		bitmap = m.secondBitmap
		startBit = 64
	}

	totalBits := len(bitmap) * 8

	for i := 0; i < totalBits; i++ {

		// Find which byte and which bit
		byteIndex := i / 8
		bitIndex := i % 8

		if bitmap[byteIndex]&(1<<(7-bitIndex)) == 0 {
			continue
		}

		bitNum := i + 1 + startBit
		m.activeBit = append(m.activeBit, bitNum)
		if bitNum == 1 {
			asciiHex = m.byteData[m.cursor : m.cursor+16]
			m.secondBitmap, err = hex.DecodeString(string(asciiHex))
			if err != nil {
				return ErrInvalidBitMap
			}
			m.cursor += 16
			continue
		}

		var value []byte

		length, err := m.parseBitLength(bitNum)
		if err != nil {
			return err
		}

		if length < 0 {
			continue
		}

		if len(m.byteData[m.cursor:]) < length {
			msg := fmt.Errorf("insufficient data for bit %d: need %d, have %d", bitNum, length, len(m.byteData[m.cursor:]))
			return errors.Join(msg, ErrInsufficientDataBitmap)
		}

		value = m.byteData[m.cursor : m.cursor+length]
		m.cursor += length

		m.isoMessageMap[bitNum] = value
	}

	if m.secondBitmap != nil && !m.hasSecondBitmap {
		m.hasSecondBitmap = true
		err := m.parseBitmap()
		return err
	}

	return nil
}

func (m *Message) MapActiveBit() {
	for k, _ := range m.isoMessageMap {
		if k == 0 {
			continue
		}
		m.activeBit = append(m.activeBit, k)
	}
	sort.Ints(m.activeBit)
}
