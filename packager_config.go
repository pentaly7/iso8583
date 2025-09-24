package iso8583

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

type BitLength struct {
	Type LengthType `json:"type"` // "FIXED", "LLVAR", "LLLVAR", "LLLLVAR"
	Max  int        `json:"max"`  // max length (or exact length if FIXED)
}

type LengthType string

const (
	LengthTypeFixed   LengthType = "FIXED"
	LengthTypeLLVar   LengthType = "LLVAR"
	LengthTypeLLLVar  LengthType = "LLLVAR"
	LengthTypeLLLLVar LengthType = "LLLLVAR"
)

const (
	FixedLength   = 1
	LLVarLength   = 2
	LLLVarLength  = 3
	LLLLVarLength = 4
)

func (lt *LengthType) GetPrefixLen() int {
	switch *lt {
	case LengthTypeFixed:
		return FixedLength
	case LengthTypeLLVar:
		return LLVarLength
	case LengthTypeLLLVar:
		return LLLVarLength
	case LengthTypeLLLLVar:
		return LLLLVarLength
	default:
		return 0
	}
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

var (
	ErrInvalidValue = errors.New("value does not match bit type")

	// precompiled regex for performance
	reNumeric   = regexp.MustCompile(`^[0-9]+$`)
	reAlphaNum  = regexp.MustCompile(`^[A-Za-z0-9]+$`)
	reTrackData = regexp.MustCompile(`^[0-9D=]+$`) // track 2 chars
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
