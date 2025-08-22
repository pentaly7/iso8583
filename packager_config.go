package iso8583

import (
	"encoding/json"
	"fmt"
	"strings"
)

type BitLength struct {
	Type LengthType `json:"mode"` // "FIXED", "LLVAR", "LLLVAR", "LLLLVAR"
	Max  int        `json:"max"`  // max length (or exact length if FIXED)
}

type LengthType string

const (
	LengthTypeFixed   LengthType = "FIXED"
	LengthTypeLLVar   LengthType = "LLVAR"
	LengthTypeLLLVar  LengthType = "LLLVAR"
	LengthTypeLLLLVar LengthType = "LLLLVAR"
)

func (lt *LengthType) GetPrefixLen() int {
	switch *lt {
	case LengthTypeFixed:
		return 0
	case LengthTypeLLVar:
		return 2
	case LengthTypeLLLVar:
		return 3
	case LengthTypeLLLLVar:
		return 4
	default:
		return -1
	}
}

func (lt *LengthType) encodeLen(n, size int) []byte {
	return []byte(fmt.Sprintf("%0*d", size, n))
}

// UnmarshalJSON Implement json.Unmarshaler
func (lt *LengthType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.ToUpper(s)
	lenType := LengthType(s)
	switch lenType {
	case LengthTypeFixed, LengthTypeLLVar, LengthTypeLLLVar, LengthTypeLLLLVar:
		*lt = lenType
	default:
		return ErrInvalidBitType
	}
	return nil
}

type BitType string

const (
	BitTypeN   BitType = "n"   // numeric
	BitTypeAN  BitType = "an"  // alphanumeric
	BitTypeANS BitType = "ans" // alphanumeric + special
	BitTypeB   BitType = "b"   // binary
	BitTypeZ   BitType = "z"   // track data
)

// UnmarshalJSON Implement json.Unmarshaler
func (bt *BitType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	bitType := BitType(s)
	switch bitType {
	case BitTypeN, BitTypeAN, BitTypeANS, BitTypeB, BitTypeZ:
		*bt = BitType(s)
	default:
		return ErrInvalidBitType
	}
	return nil
}
