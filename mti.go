package iso8583

import "unsafe"

type MTIType string

const (
	MTICardProcessingRequest   MTIType = "0100"
	MTICardProcessingResponse  MTIType = "0110"
	MTIFinancialRequest        MTIType = "0200"
	MTIFinancialResponse       MTIType = "0210"
	MTIReversalRequest         MTIType = "0400"
	MTIReversalResponse        MTIType = "0410"
	MTIRepeatedReversalRequest MTIType = "0401"
	MTINMMRequest              MTIType = "0800"
	MTINMMResponse             MTIType = "0810"
)

type MTITypeByte [4]byte

var (
	MTICardProcessingRequestByte   MTITypeByte = [4]byte{0x30, 0x31, 0x30, 0x30}
	MTICardProcessingResponseByte  MTITypeByte = [4]byte{0x30, 0x31, 0x31, 0x30}
	MTIFinancialRequestByte        MTITypeByte = [4]byte{0x30, 0x32, 0x30, 0x30}
	MTIFinancialResponseByte       MTITypeByte = [4]byte{0x30, 0x32, 0x30, 0x30}
	MTIReversalRequestByte         MTITypeByte = [4]byte{0x30, 0x34, 0x30, 0x30}
	MTIReversalResponseByte        MTITypeByte = [4]byte{0x30, 0x34, 0x31, 0x30}
	MTIRepeatedReversalRequestByte MTITypeByte = [4]byte{0x30, 0x34, 0x30, 0x31}
	MTINMMRequestByte              MTITypeByte = [4]byte{0x30, 0x38, 0x30, 0x30}
	MTINMMResponseByte             MTITypeByte = [4]byte{0x30, 0x34, 0x31, 0x30}
)

// Pre-computed valid MTI lookup for O(1) validation
var validMTIs = map[[4]byte]struct{}{
	MTICardProcessingRequestByte:   {},
	MTICardProcessingResponseByte:  {},
	MTIFinancialRequestByte:        {},
	MTIFinancialResponseByte:       {},
	MTIReversalRequestByte:         {},
	MTIReversalResponseByte:        {},
	MTIRepeatedReversalRequestByte: {},
	MTINMMRequestByte:              {},
	MTINMMResponseByte:             {},
}

func (m MTITypeByte) ToMtiString() MTIType {
	return MTIType(m[:])
}
func (m MTITypeByte) String() string {
	return unsafe.String(unsafe.SliceData(m[:]), len(m[:]))
}

func (m MTITypeByte) Equal(other MTITypeByte) bool {
	return m == other
}

func (ms MTIType) ToMtiByte() MTITypeByte {
	var b MTITypeByte
	copy(b[:], ms) // string → []byte → [4]byte
	return b
}

func isValidMti(mti MTITypeByte) bool {
	_, ok := validMTIs[mti]
	return ok
}
