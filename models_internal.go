package iso8583

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func (m *Message) parseBitmap(bitmap uint64, offset int) (err error) {
	for i := 0; i < 64; i++ {
		// using "AND" operator by shifting the bit and check if its 0
		if bitmap&(1<<(63-i)) == 0 {
			continue
		}
		bitNum := i + offset + 1

		if bitNum == 1 {
			m.actSecondBitmap = true
			m.secondBitmap = string(m.runeData[:16])
			copy(m.runeData, m.runeData[16:])
			continue
		}

		var value string

		length, err := m.parseBitLength(bitNum)
		if err != nil {
			return err
		}

		if length <= 0 {
			continue
		}

		if len(m.runeData) < length {
			msg := fmt.Errorf("insufficient data for bit %d: need %d, have %d", bitNum, length, len(m.runeData))
			return errors.Join(msg, ErrInsufficientDataBitmap)
		}

		value = string(m.runeData[:length])
		copy(m.runeData, m.runeData[length:])

		m.IsoMessageMap[bitNum] = value
	}
	return nil
}

func (m *Message) parseBitLength(bitNum int) (length int, err error) {

	packager := m.IsoPack.IsoPackager[bitNum-1]

	switch {
	case packager > 0:
		length = packager
	case packager < 0:
		lengthSize := -packager
		if len(m.runeData) < lengthSize {
			msg := fmt.Errorf("insufficient data for bit %d length: need %d, have %d", bitNum, lengthSize, len(m.runeData))
			return length, errors.Join(msg, ErrFailedToParseBitmapData)
		}
		length, err = strconv.Atoi(string(m.runeData[:lengthSize]))
		if err != nil {
			msg := fmt.Errorf("failed to parse length for bit %d", bitNum)
			return length, errors.Join(msg, ErrFailedToParseBitmapData)
		}
		copy(m.runeData, m.runeData[lengthSize:])
	default:
		return
	}

	return
}

// GetKeyActive to get Bit Active in Message ISO
func (m *Message) GetKeyActive() error {
	keys := make([]int, 0, len(m.IsoMessageMap))
	for k := range m.IsoMessageMap {
		keys = append(keys, k)
	}
	m.BitActive = keys
	sort.Ints(m.BitActive)

	if m.MTI == "0800" || m.MTI == "0810" {
		return nil
	}
	for _, bit := range m.IsoPack.MandatoryBits {
		if _, ok := m.IsoMessageMap[bit]; ok == false {
			return fmt.Errorf("missing mandatory bit %d", bit)
		}
	}
	if m.MTI == "0400" || m.MTI == "0401" {
		if _, ok := m.IsoMessageMap[90]; ok == false {
			return fmt.Errorf("missing mandatory bit 90 for reversal")
		}
	}
	return nil
}

// PackISO to get Create Message ISO in string
func (m *Message) PackISO() error {
	err := m.GetKeyActive()
	if err != nil {
		return err
	}

	if m.MTI == "" {
		return ErrNoMtiToPack
	}

	bitmap := make([]byte, 128)
	for _, bit := range m.BitActive {
		if bit < 1 || bit > 128 {
			err := ErrInvalidBitNumber
			return err
		}
		bitmap[bit-1] = '1'
	}

	if bytes.Contains(bitmap[64:], []byte{'1'}) {
		bitmap[0] = '1'
		m.actSecondBitmap = true
	}

	m.firstBitmap = formatBitmap(bitmap[:64])
	if m.actSecondBitmap {
		m.secondBitmap = formatBitmap(bitmap[64:])
	}

	result, err := m.processPackIso()
	if err != nil {
		return err
	}

	m.AsString = result.String()
	return nil
}

func (m *Message) processPackIso() (result strings.Builder, err error) {

	result.Grow(len(m.MTI) + len(m.firstBitmap) + len(m.secondBitmap) + 1000) // Preallocate some space

	result.WriteString(m.MTI)
	result.WriteString(m.firstBitmap)
	if m.actSecondBitmap {
		result.WriteString(m.secondBitmap)
	}

	for _, i := range m.BitActive {
		packager := m.IsoPack.IsoPackager[i-1]
		value := []rune(m.IsoMessageMap[i])

		switch {
		case packager > 0:
			result.WriteString(string(value))
		case packager == -2:
			_, err := fmt.Fprintf(&result, "%02d%s", len(value), string(value))
			if err != nil {
				return result, ErrInvalidPackager
			}
		case packager == -3:
			_, err := fmt.Fprintf(&result, "%03d%s", len(value), string(value))
			if err != nil {
				return result, ErrInvalidPackager
			}
		default:
			err := fmt.Errorf("invalid packager value for bit %d: %d", i, packager)
			return result, errors.Join(err, ErrInvalidPackager)
		}
	}

	return
}

func formatBitmap(bits []byte) string {
	var bitmap uint64
	for i, bit := range bits {
		if bit == '1' {
			bitmap |= 1 << (63 - i)
		}
	}
	return strings.ToUpper(fmt.Sprintf("%016x", bitmap))
}

// checkMessageISO for checking Message, MTI or  assume there is ISO Header
func (m *Message) checkMessageISO() (string, error) {
	data := []rune(m.AsString)
	// check ISO header
	if m.IsoPack.IsoHeaderActive || string(data[0:3]) == "ISO" {
		data = data[m.IsoPack.IsoHeaderLength:]
		return string(data), nil
	}
	// check MTI
	mti := string(data[0:4])
	mtiSlices := []string{"0200", "0210", "0400", "0410", "0800", "0810"}
	for _, list := range mtiSlices {
		if mti == list {
			return m.AsString, nil
		}
		continue
	}
	return m.AsString, ErrNotDefaultMti
}
