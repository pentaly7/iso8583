package iso8583

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

func (m *Message) parseBitmap(bitmap []byte) (err error) {
	startBit := 0
	if m.hasSecondBitmap {
		startBit = 64
	}

	totalBits := len(bitmap) * 8

	for i := startBit; i < totalBits; i++ {

		// Find which byte and which bit
		byteIndex := i / 8
		bitIndex := i % 8

		if bitmap[byteIndex]&(1<<(7-bitIndex)) == 0 {
			continue
		}

		bitNum := i + 1
		m.activeBit = append(m.activeBit, bitNum)
		if bitNum == 1 {
			m.secondBitmap = m.byteData[m.cursor : m.cursor+16]
			m.cursor += 16
			continue
		}

		var value []byte

		length, err := m.parseBitLength(bitNum)
		if err != nil {
			return err
		}

		if length <= 0 {
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
		err := m.parseBitmap(m.secondBitmap)
		return err
	}

	return nil
}

// ValidateMessage validate ISO Message
func (m *Message) ValidateMessage() error {
	if m.MTI.Equal(MTINMMRequestByte) || m.MTI.Equal(MTINMMResponseByte) {
		mandatoryBits := []int{7, 11, 70}
		for _, bit := range mandatoryBits {
			if _, ok := m.isoMessageMap[bit]; ok == false {
				return fmt.Errorf("missing mandatory bit %d", bit)
			}
		}
		if m.MTI.Equal(MTINMMRequestByte) {
			return nil
		}
		if _, ok := m.isoMessageMap[39]; ok == false {
			return fmt.Errorf("missing mandatory bit 39 for response")
		}
		return nil
	}

	for _, bit := range m.packager.MandatoryBit {
		if _, ok := m.isoMessageMap[bit]; ok == false {
			return fmt.Errorf("missing mandatory bit %d", bit)
		}
	}

	if m.MTI.Equal(MTIReversalRequestByte) ||
		m.MTI.Equal(MTIReversalResponseByte) ||
		m.MTI.Equal(MTIRepeatedReversalRequestByte) {
		if _, ok := m.isoMessageMap[90]; ok == false {
			return fmt.Errorf("missing mandatory bit 90 for reversal")
		}
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
}

// PackISO to get Create Message ISO in string
func (m *Message) PackISO() error {
	m.MapActiveBit()

	m.dataLength = 0

	if m.packager.HasHeader {
		m.dataLength += m.packager.HeaderLength
	}

	if m.MTI == nil {
		return ErrNoMtiToPack
	}
	m.dataLength += len(m.MTI)

	bitmap := make([]byte, 16) // max 128 bits = 16 bytes

	for _, bit := range m.activeBit {
		if bit < 1 || bit > 128 {
			return ErrInvalidBitNumber
		}

		length, err := m.getTotalBitLength(bit)
		if err != nil {
			return err
		}

		m.dataLength += length

		// bit position (ISO8583 is Big Endian, MSB first)
		byteIndex := (bit - 1) / 8
		bitIndex := (bit - 1) % 8
		bitmap[byteIndex] |= 1 << (7 - bitIndex)
	}

	m.firstBitmap = bitmap[:8]
	m.dataLength += len(m.firstBitmap)

	// check second bitmap
	if !bytes.Equal(bitmap[8:], make([]byte, 8)) {
		bitmap[0] |= 0x80 // set first bit to indicate second bitmap
		m.hasSecondBitmap = true
		m.secondBitmap = bitmap[8:16]
		m.dataLength += len(m.secondBitmap)
	}

	result, err := m.processPackIso()
	if err != nil {
		return err
	}

	m.byteData = result
	return nil
}

func (m *Message) processPackIso() (result []byte, err error) {
	result = make([]byte, 0, m.dataLength)

	if m.packager.HasHeader {
		result = append(result, m.header...)
	}

	result = append(result, m.MTI...)
	result = append(result, m.firstBitmap...)

	if m.hasSecondBitmap {
		result = append(result, m.secondBitmap...)
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

		encPrefixLen := lenType.encodeLen(prefixLen, len(m.isoMessageMap[i]))
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
