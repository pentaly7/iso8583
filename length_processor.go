package iso8583

import (
	"errors"
	"fmt"
)

func (m *Message) parseBitLength(bitNum int) (length int, err error) {

	packager, ok := m.packager.IsoPackagerConfig[bitNum]
	if !ok {
		return 0, fmt.Errorf("packager not found for bit %d", bitNum)
	}

	prefixLen := packager.Length.Type.GetPrefixLen()

	if prefixLen == 0 {
		length = packager.Length.Max
		return length, nil
	}

	if len(m.byteData[m.cursor:]) < prefixLen {
		msg := fmt.Errorf("insufficient data for bit %d length: need %d, have %d", bitNum, prefixLen, len(m.byteData[m.cursor:]))
		return length, errors.Join(msg, ErrFailedToParseBitmapData)
	}
	length, err = asciiBytesToInt(m.byteData[m.cursor : m.cursor+prefixLen])
	if err != nil {
		msg := fmt.Errorf("failed to parse length for bit %d", bitNum)
		return length, errors.Join(msg, ErrFailedToParseBitmapData)
	}
	m.cursor += prefixLen

	return
}

// asciiBytesToInt converts an ASCII byte array to an integer
func asciiBytesToInt(b []byte) (int, error) {
	n := 0
	// loop the byte array, then check the digit
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid digit %q", c)
		}
		// convert to int
		n = n*10 + int(c-'0')
	}
	return n, nil
}

// getTotalBitLength returns the total length of the bit
// for LLVar, LLLVar, LLLLVar it returns the length of the data + the length of the length data
func (m *Message) getTotalBitLength(bitNum int) (length int, err error) {

	packager, ok := m.packager.IsoPackagerConfig[bitNum]
	if !ok {
		return 0, fmt.Errorf("packager not found for bit %d", bitNum)
	}

	prefixLen := packager.Length.Type.GetPrefixLen()
	if prefixLen == 0 {
		if len(m.isoMessageMap[bitNum]) != packager.Length.Max {
			return 0, fmt.Errorf("invalid bit length for bit %d: expected %d, got %d", bitNum, packager.Length.Max, len(m.isoMessageMap[bitNum]))
		}
		return packager.Length.Max, nil
	}

	length = len(m.isoMessageMap[bitNum])
	if length > packager.Length.Max {
		return 0, fmt.Errorf("invalid bit length for bit %d: max %d, got %d", bitNum, packager.Length.Max, len(m.isoMessageMap[bitNum]))
	}

	return length + prefixLen, nil
}
