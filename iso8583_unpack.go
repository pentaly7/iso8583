package iso8583

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

const BitmapLength = 16

// Unpack to Single Data Element
func (m *Message) Unpack(b []byte) error {
	m.byteData = b
	m.cursor = 0

	if len(m.byteData) < 4 {
		return ErrInsufficientDataMti
	}
	// check ISO header

	if m.packager.HasHeader {

		m.header = m.byteData[:m.packager.HeaderLength]
		m.cursor += m.packager.HeaderLength
	} else if bytes.Equal(m.byteData[:3], isoHeader) {
		m.header = isoHeader
		m.cursor += 3
	}

	if len(m.byteData[m.cursor:]) < 4 {
		return ErrInsufficientDataMti
	}

	// check MTI
	mti := MTITypeByte(m.byteData[0:4])
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

	if len(m.byteData[m.cursor:]) < BitmapLength {
		return ErrInsufficientDataFirstBitmap
	}

	if m.isoMessageMap == nil {
		m.isoMessageMap = make(map[int][]byte)
	}

	if err := m.parseBitmap(); err != nil {
		return err
	}

	// remove byte data to free memory
	m.byteData = nil

	return nil
}

func (m *Message) parseBitmap() (err error) {
	startBit := 0
	var bitmap []byte
	var asciiHex []byte

	// process first bitmap if the bitmap still empty
	if m.firstBitmap == nil {
		// bitmap is in hex string, it needs to be converted to byte
		asciiHex = m.byteData[m.cursor : m.cursor+BitmapLength]
		bitmap, err = hex.DecodeString(string(asciiHex))
		if err != nil {
			return ErrInvalidBitMap
		}
		// assign first bitmap to prevent parsing it again if second bitmap active
		m.firstBitmap = bitmap
		m.cursor += BitmapLength
	} else if m.secondBitmap != nil { // check if second bitmap not nil
		m.secondBitmapFlag = true
		bitmap = m.secondBitmap
		startBit = 64
	}

	// 1 byte = 8 bit
	// so it need to x8 to get the total bit
	totalBits := len(bitmap) * 8

	for i := 0; i < totalBits; i++ {

		// Find which byte and which bit
		byteIndex := i / 8 // BYTE index
		bitIndex := i % 8  // BIT index

		// use AND operator to check if the bit is active
		// ex checking bit 4(i=3 because start from 0) :
		// byte index = 3 / 8 = 0
		// bit index = 3 % 8 = 4
		// 00011110 AND (00000001 LSHIFT 7-4) 	== 0
		// 00011110 AND (00000001 LSHIFT 3) 	== 0
		// 00011110 AND (00010000)				== 0
		// 00010000 							== 0 is false
		// so the bit 4 is active
		if bitmap[byteIndex]&(1<<(7-bitIndex)) == 0 {
			continue
		}

		// need to add 1, because start from 0
		// start bit is for secondary bit
		bitNum := i + 1 + startBit
		m.activeBit = append(m.activeBit, bitNum)

		// checking secondary bit
		if bitNum == 1 {
			// convert HEX string to []byte
			asciiHex = m.byteData[m.cursor : m.cursor+BitmapLength]
			m.secondBitmap, err = hex.DecodeString(string(asciiHex))
			if err != nil {
				return ErrInvalidBitMap
			}
			m.cursor += BitmapLength
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

	// check if there is a second bitmap
	// mark it to has second bitmap to prevent infinite loop
	// second bitmap flag will mark the second bitmap already parsed or not
	if m.secondBitmap != nil && !m.secondBitmapFlag {
		// call parse bitmap again to parse the second bitmap
		err := m.parseBitmap()
		return err
	}

	return nil
}
