package iso8583

import (
	"strings"
	"unsafe"
)

// bitPositions[byte] = slice of positions (1–8) of set bits in that byte.
// Example: 0b10100000 (0xA0) -> []int{1, 3}
var (
	bitPositions   [256][8]int
	bitCounts      [256]int
	fourDigitTable [10000][4]byte
	hexTable       [256][2]byte
)

// init initializes lookup tables used for bit manipulation and data conversion.
// It pre-calculates bit positions, bit counts, hexadecimal representations,
// and four-digit decimal representations for improved performance.
func init() {
	initLookupTables()
}

func initLookupTables() {
	const digits = "0123456789ABCDEF"
	for b := 0; b < 10000; b++ {
		if b < 256 {
			pos := 0
			for i := 0; i < 8; i++ {
				if b&(1<<(7-i)) != 0 {
					bitPositions[b][pos] = i + 1 // bit numbers 1–8
					pos++
				}
			}
			bitCounts[b] = pos

			hexTable[b][0] = digits[b>>4]
			hexTable[b][1] = digits[b&0x0F]
		}

		fourDigitTable[b][0] = byte('0' + b/1000)
		fourDigitTable[b][1] = byte('0' + (b/100)%10)
		fourDigitTable[b][2] = byte('0' + (b/10)%10)
		fourDigitTable[b][3] = byte('0' + b%10)
	}
}

type (
	// Message for Component Message
	Message struct {
		MTI           MTITypeByte // MTI
		header        []byte      // iso header
		isoMessageMap [129][]byte // Get Element of Iso Message in Map
		activeBits    [129]int
		activeCount   int
		packager      *IsoPackager
		length        int
		err           error
	}
)

func NewMessage(packager *IsoPackager) *Message {

	return &Message{
		packager: packager,
	}
}

func (m *Message) SetByte(bit int, value []byte) *Message {
	if value == nil {
		return m
	}
	if m.isoMessageMap[bit] == nil { // only insert if new
		m.appendBit(bit)
	}
	m.isoMessageMap[bit] = value
	return m
}

func (m *Message) SetHeader(value []byte) *Message {
	if value == nil {
		return m
	}
	m.header = value
	return m
}

func (m *Message) SetHeaderString(value string) *Message {
	if value == "" {
		return m
	}
	m.header = []byte(value)
	return m
}

func (m *Message) SetString(bit int, s string) *Message {
	if m.isoMessageMap[bit] == nil { // only insert if new
		m.appendBit(bit)
	}
	if len(s) == 0 {
		m.isoMessageMap[bit] = []byte{}
		return m
	}
	m.isoMessageMap[bit] = unsafe.Slice(unsafe.StringData(s), len(s))
	return m
}
func (m *Message) SetMtiString(s MTIType) *Message {
	m.MTI = s.ToMtiByte()
	return m
}
func (m *Message) SetMTIByte(value MTITypeByte) *Message {
	m.MTI = value
	return m
}

func (m *Message) appendBit(bit int) {
	m.activeBits[m.activeCount] = bit
	m.activeCount++
}

func (m *Message) Unset(bit int) *Message {
	for i := 0; i < m.activeCount; i++ {
		if m.activeBits[i] == bit {
			m.isoMessageMap[bit] = nil

			// swap with last
			m.activeBits[i] = m.activeBits[m.activeCount-1]
			m.activeBits[m.activeCount-1] = 0 // optional cleanup
			m.activeCount--
			break
		}
	}
	return m
}

func (m *Message) GetString(bit int) string {
	b := m.isoMessageMap[bit]
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func (m *Message) GetByte(bit int) []byte {
	return m.isoMessageMap[bit]
}

func (m *Message) HasBit(bit int) bool {
	val := m.isoMessageMap[bit]
	if val != nil {
		return true
	}
	return false
}

// GetMessageKey is to Trace ISO Message created for Tracing ISO Message Respon
func (m *Message) GetMessageKey() string {
	var sb strings.Builder

	sb.Grow(128)
	// MTI (4 bytes)
	if m.IsRequest() {
		mtiRes, _ := m.GetMTIResponse()
		sb.Write(mtiRes[:])
	} else {
		sb.Write(m.MTI[:])
	}

	// Key fields
	for _, bitNum := range m.packager.MessageKey {
		sb.Write(m.isoMessageMap[bitNum])
	}

	return sb.String()
}

func (m *Message) GetMTIResponse() (mti MTITypeByte, err error) {
	switch {
	case m.MTI.Equal(MTICardProcessingRequestByte):
		mti = MTICardProcessingResponseByte
	case m.MTI.Equal(MTIFinancialRequestByte):
		mti = MTIFinancialResponseByte
	case m.MTI.Equal(MTIReversalRequestByte),
		m.MTI.Equal(MTIRepeatedReversalRequestByte):
		mti = MTIReversalResponseByte
	case m.MTI.Equal(MTIReversalAdviceRequestByte),
		m.MTI.Equal(MTIReversalAdviceRepeatedRequestByte):
		mti = MTIReversalAdviceResponseByte
	case m.MTI.Equal(MTINMMRequestByte):
		mti = MTINMMResponseByte
	default:
		return mti, ErrNotDefaultMti
	}
	return mti, nil
}

// SetMTIResponse is set MTI Response for Response Message
func (m *Message) SetMTIResponse() error {
	mti, err := m.GetMTIResponse()
	if err != nil {
		return err
	}
	m.MTI = mti
	return nil
}

// IsRequest for check MTI is Request or Response
func (m *Message) IsRequest() bool {
	switch {
	case m.MTI.Equal(MTICardProcessingRequestByte),
		m.MTI.Equal(MTIFinancialRequestByte),
		m.MTI.Equal(MTIReversalRequestByte),
		m.MTI.Equal(MTIRepeatedReversalRequestByte),
		m.MTI.Equal(MTIReversalAdviceRequestByte),
		m.MTI.Equal(MTIReversalAdviceRepeatedRequestByte),
		m.MTI.Equal(MTINMMRequestByte):
		return true
	default:
		return false
	}
}

// IsTransactional for check MTI message is Transaction or NMM
func (m *Message) IsTransactional() bool {
	switch {
	case m.MTI.Equal(MTICardProcessingRequestByte),
		m.MTI.Equal(MTICardProcessingResponseByte),
		m.MTI.Equal(MTIFinancialRequestByte),
		m.MTI.Equal(MTIFinancialResponseByte),
		m.MTI.Equal(MTIReversalRequestByte),
		m.MTI.Equal(MTIReversalResponseByte),
		m.MTI.Equal(MTIReversalAdviceRequestByte),
		m.MTI.Equal(MTIReversalAdviceResponseByte),
		m.MTI.Equal(MTIReversalAdviceRepeatedRequestByte),
		m.MTI.Equal(MTIRepeatedReversalRequestByte):
		return true
	default:
		return false
	}

}

// IsNMM for check MTI message is Transaction or NMM
func (m *Message) IsNMM() bool {
	switch {
	case m.MTI.Equal(MTINMMRequestByte),
		m.MTI.Equal(MTINMMResponseByte):
		return true
	default:
		return false
	}
}

// IsResponse for check MTI is Request or Response
func (m *Message) IsResponse() bool {
	switch {
	case m.MTI.Equal(MTICardProcessingResponseByte),
		m.MTI.Equal(MTIFinancialResponseByte),
		m.MTI.Equal(MTIReversalResponseByte),
		m.MTI.Equal(MTINMMResponseByte),
		m.MTI.Equal(MTIReversalAdviceResponseByte):
		return true
	default:
		return false
	}
}

// ClearEntries for clear all entries so this message can be reused
func (m *Message) ClearEntries() {
	m.MTI = EmptyMti
	m.header = nil
	for i := 0; i < m.activeCount; i++ {
		m.isoMessageMap[m.activeBits[i]] = nil
	}
	m.activeCount = 0
}

// CreateResponseISO create response ISO Message
func CreateResponseISO(m *Message, rc string) (*Message, error) {
	// create msg
	out := NewMessage(m.packager)
	out.MTI = m.MTI
	out.header = m.header
	for bit, v := range m.isoMessageMap {
		if v != nil {
			out.SetByte(bit, v)
		}
	}
	err := out.SetMTIResponse()
	if err != nil {
		return nil, err
	}
	out.SetString(39, rc)
	return out, nil
}

func CloneMessage(m *Message) *Message {
	out := NewMessage(m.packager)
	out.MTI = m.MTI
	out.header = m.header
	for bit, v := range m.isoMessageMap {
		if v != nil {
			out.SetByte(bit, v)
		}
	}
	return out
}
