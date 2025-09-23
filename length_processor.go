package iso8583

import (
	"errors"
	"fmt"
)

func (m *Message) parseBitLength(b []byte, bitNum, cursor int) (length, prefixLen int, err error) {

	prefixLen = m.packager.PrefixLengths[bitNum]

	if prefixLen == 0 {
		length = m.packager.MaxLengths[bitNum]
		if length == 0 {
			return -1, prefixLen, fmt.Errorf("packager not found for bit %d", bitNum)
		}
		return length, prefixLen, nil
	}

	if len(b[cursor:]) < prefixLen {
		msg := fmt.Errorf("insufficient data for bit %d length: need %d, have %d", bitNum, prefixLen, len(b[cursor:]))
		return length, prefixLen, errors.Join(msg, ErrFailedToParseBitmapData)
	}
	length, err = asciiBytesToInt(b[cursor : cursor+prefixLen])
	if err != nil {
		msg := fmt.Errorf("failed to parse length for bit %d", bitNum)
		return length, prefixLen, errors.Join(msg, ErrFailedToParseBitmapData)
	}

	return
}

var digitLookup = [256]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4,
	'5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
}

// asciiBytesToInt converts an ASCII byte array to an integer
func asciiBytesToInt(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	n := 0
	for _, c := range b {
		digit := digitLookup[c]
		if digit == 0 && c != '0' {
			return 0, fmt.Errorf("invalid digit %q", c)
		}
		n = n*10 + digit
	}
	return n, nil
}

// getTotalBitLength returns the total length of the bit
// for LLVar, LLLVar, LLLLVar it returns the length of the data + the length of the length data
func (m *Message) getTotalBitLength(bitNum int) (length int, err error) {

	packager := m.packager.IsoPackagerConfig[bitNum]
	if packager.Type == "" {
		return 0, fmt.Errorf("packager not found for bit %d", bitNum)
	}

	prefixLen := packager.Length.Type.GetPrefixLen()
	if prefixLen == 0 {
		if len(m.isoMessageMap[bitNum]) != packager.Length.Max {
			return 0, fmt.Errorf(
				"invalid bit length for bit %d: expected %d, got %d",
				bitNum,
				packager.Length.Max,
				len(m.isoMessageMap[bitNum]),
			)
		}
		return packager.Length.Max, nil
	}

	length = len(m.isoMessageMap[bitNum])
	if length > packager.Length.Max {
		return 0, fmt.Errorf(
			"invalid bit length for bit %d: max %d, got %d",
			bitNum,
			packager.Length.Max,
			len(m.isoMessageMap[bitNum]),
		)
	}

	return length + prefixLen, nil
}
