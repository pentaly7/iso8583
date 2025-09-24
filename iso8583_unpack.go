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

	cursor := 0

	if len(b) < 4 {
		return ErrInsufficientDataMti
	}
	// check ISO header

	if m.packager.HasHeader {
		m.header = b[:m.packager.HeaderLength]
		cursor += m.packager.HeaderLength
	} else if bytes.Equal(b[:3], isoHeader) {
		m.header = isoHeader
		cursor += 3
	}

	if len(b[cursor:]) < 4 {
		return ErrInsufficientDataMti
	}

	// check MTI
	mti := MTITypeByte(b[0:4])
	if !isValidMti(mti) {
		return ErrNotDefaultMti
	}
	m.MTI = mti
	cursor += 4

	if len(b[cursor:]) < BitmapLength {
		return ErrInsufficientDataFirstBitmap
	}

	if err := m.parseBitmap(b, cursor); err != nil {
		return err
	}

	return nil
}
func (m *Message) parseBitmap(b []byte, cursor int) error {
	// Ensure enough data for at least a primary bitmap
	if len(b[cursor:]) < BitmapLength {
		return errors.Join(fmt.Errorf("insufficient data for bitmap: need %d, have %d", BitmapLength, len(b[cursor:])), ErrInsufficientDataBitmap)
	}

	// ----- parse primary bitmap -----
	// asciiHex is a slice of the hex characters in b (no allocation)
	asciiHex := b[cursor : cursor+BitmapLength]

	// get a buffer from pool and decode into it
	bitmap := [16]byte{}

	// bmpBuf has length 8 (as constructed above). Use that directly.
	if _, err := hex.Decode(bitmap[:8], asciiHex); err != nil {
		return ErrInvalidBitMap
	}

	cursor += BitmapLength

	maxBits := 8
	// If bit 1 (first bit) is set, there is a secondary bitmap to parse later.
	if bitmap[0]&(0x80) != 0 { // bit index 0 -> bit 1
		if len(b[cursor:]) < BitmapLength {
			return errors.Join(fmt.Errorf("insufficient data for second bitmap: need %d, have %d", BitmapLength, len(b[cursor:])), ErrInsufficientDataBitmap)
		}

		asciiHex2 := b[cursor : cursor+BitmapLength]

		if _, err := hex.Decode(bitmap[8:], asciiHex2); err != nil {
			return ErrInvalidBitMap
		}
		cursor += BitmapLength
		maxBits = 16
		// flip the bit 1
		bitmap[0] &= 0x7F
	}

	// Process primary bitmap bits
	for byteIdx := 0; byteIdx < maxBits; byteIdx++ {
		v := bitmap[byteIdx]
		if v == 0 {
			continue
		}

		count := bitCounts[v]
		for i := 0; i < count; i++ {
			bitNum := byteIdx*8 + bitPositions[v][i]

			length, prefixLen, err := m.parseBitLength(b, bitNum, cursor)
			if err != nil {
				return err
			}
			if length < 0 {
				continue
			}

			cursor += prefixLen
			if len(b[cursor:]) < length {
				msg := fmt.Errorf("insufficient data for bit %d: need %d, have %d", bitNum, length, len(b[cursor:]))
				return errors.Join(msg, ErrInsufficientDataBitmap)
			}

			value := b[cursor : cursor+length]
			cursor += length
			m.isoMessageMap[bitNum] = value
			m.appendBit(bitNum)
		}

	}

	return nil
}
