package iso8583

import "bytes"

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

type MTITypeByte []byte

var (
	MTICardProcessingRequestByte   MTITypeByte = []byte("0100")
	MTICardProcessingResponseByte  MTITypeByte = []byte("0110")
	MTIFinancialRequestByte        MTITypeByte = []byte("0200")
	MTIFinancialResponseByte       MTITypeByte = []byte("0210")
	MTIReversalRequestByte         MTITypeByte = []byte("0400")
	MTIReversalResponseByte        MTITypeByte = []byte("0410")
	MTIRepeatedReversalRequestByte MTITypeByte = []byte("0401")
	MTINMMRequestByte              MTITypeByte = []byte("0800")
	MTINMMResponseByte             MTITypeByte = []byte("0810")
)

func (m MTITypeByte) String() string {
	return string(m)
}

func (m MTITypeByte) Equal(other MTITypeByte) bool {
	return bytes.Equal(m, other)
}
