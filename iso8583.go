package iso8583

import "fmt"

type (
	// Message for Component Message
	Message struct {
		packager        *IsoPackager
		header          []byte         // iso header
		MTI             MTITypeByte    // MTI
		firstBitmap     []byte         // First bitmap
		hasSecondBitmap bool           // Second Bitmap activation
		secondBitmap    []byte         // Second bitmap
		isoMessageMap   map[int][]byte // Get Element of Iso Message in Map
		activeBit       []int          // Find Bit Active for check Mandatory Bit
		byteData        []byte
		cursor          int
		dataLength      int
	}
)

func NewMessage(packager *IsoPackager) *Message {

	return &Message{
		packager:      packager,
		isoMessageMap: make(map[int][]byte),
	}
}

func (m *Message) SetByte(bit int, value []byte) *Message {
	m.isoMessageMap[bit] = value
	return m
}

func (m *Message) SetString(bit int, value string) *Message {
	m.isoMessageMap[bit] = []byte(value)
	return m
}
func (m *Message) SetMTI(value MTIType) *Message {
	m.MTI = []byte(value)
	return m
}

func (m *Message) GetString(bit int) string {
	return string(m.isoMessageMap[bit])
}

func (m *Message) GetByte(bit int) []byte {
	return m.isoMessageMap[bit]
}

func (m *Message) HasBit(bit int) bool {
	_, ok := m.isoMessageMap[bit]
	return ok
}

// ValidateMessage validate ISO Message
func (m *Message) ValidateMessage() error {
	if m.MTI.Equal(MTINMMRequestByte) || m.MTI.Equal(MTINMMResponseByte) {
		mandatoryBits := []int{7, 11, 70}
		for _, bit := range mandatoryBits {
			if _, ok := m.isoMessageMap[bit]; ok == false {
				return fmt.Errorf("missing mandatory bit %d", bit)
			}
		}
		if m.MTI.Equal(MTINMMRequestByte) {
			return nil
		}
		if _, ok := m.isoMessageMap[39]; ok == false {
			return fmt.Errorf("missing mandatory bit 39 for response")
		}
		return nil
	}

	for _, bit := range m.packager.MandatoryBit {
		if _, ok := m.isoMessageMap[bit]; ok == false {
			return fmt.Errorf("missing mandatory bit %d", bit)
		}
	}

	if m.MTI.Equal(MTIReversalRequestByte) ||
		m.MTI.Equal(MTIReversalResponseByte) ||
		m.MTI.Equal(MTIRepeatedReversalRequestByte) {
		if _, ok := m.isoMessageMap[90]; ok == false {
			return fmt.Errorf("missing mandatory bit 90 for reversal")
		}
	}

	return nil
}

// GetMessageKey is to Trace ISO Message created for Tracing ISO Message Respon
func (m *Message) GetMessageKey() string {
	var result string
	// combination ISO => MTI(respon), PAN , STAN, RRN, TransmissionDateTime
	if m.IsRequest() {
		mtiRes, _ := m.GetMTIResponse()
		result = string(mtiRes)
	} else {
		result = string(m.MTI)
	}
	isomap := m.isoMessageMap
	for _, v := range m.packager.MessageKey {
		result += string(isomap[v])
	}
	return result
}

func (m *Message) GetMTIResponse() (mti MTITypeByte, err error) {
	switch {
	case m.MTI.Equal(MTICardProcessingRequestByte):
		mti = MTICardProcessingResponseByte
	case m.MTI.Equal(MTIFinancialRequestByte):
		mti = MTIFinancialResponseByte
	case m.MTI.Equal(MTIReversalRequestByte), m.MTI.Equal(MTIRepeatedReversalRequestByte):
		mti = MTIReversalResponseByte
	case m.MTI.Equal(MTINMMRequestByte):
		mti = MTINMMResponseByte
	default:
		return nil, ErrNotDefaultMti
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
		m.MTI.Equal(MTINMMResponseByte):
		return true
	default:
		return false
	}
}

// CreateResponseISO create response ISO Message
func CreateResponseISO(i Message, rc string) ([]byte, error) {
	// create msg
	msg := NewMessage(i.packager)
	msg.MTI = i.MTI
	for bit, v := range i.isoMessageMap {
		msg.SetByte(bit, v)
	}
	err := msg.SetMTIResponse()
	if err != nil {
		return nil, err
	}
	msg.SetString(39, rc)
	result, err := i.PackISO()
	if err != nil {
		return nil, err
	}
	return result, nil
}
