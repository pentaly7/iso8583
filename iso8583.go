package iso8583

import (
	"bytes"
	"encoding/hex"
)

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

func (m *Message) SetByte(bit int, value []byte) {
	m.isoMessageMap[bit] = value
}

func (m *Message) SetString(bit int, value string) {
	m.isoMessageMap[bit] = []byte(value)
}

func (m *Message) GetString(bit int) string {
	return string(m.isoMessageMap[bit])
}

func (m *Message) GetByte(bit int) []byte {
	return m.isoMessageMap[bit]
}

// Unpack to Single Data Element
func (m *Message) Unpack(b []byte) error {
	m.byteData = b
	m.cursor = 0

	data := m.byteData
	if len(data) < 4 {
		return ErrInsufficientDataMti
	}
	// check ISO header
	if m.packager.HasHeader || bytes.Equal(data[:3], isoHeader) {
		m.header = data[:m.packager.HeaderLength]
		m.cursor += m.packager.HeaderLength
	}

	if len(data[m.cursor:]) < 4 {
		return ErrInsufficientDataMti
	}

	// check MTI
	mti := MTITypeByte(data[0:4])
	switch {
	case mti.Equal(MTICardProcessingRequestByte),
		mti.Equal(MTICardProcessingResponseByte),
		mti.Equal(MTIFinancialRequestByte),
		mti.Equal(MTIFinancialResponseByte),
		mti.Equal(MTIReversalRequestByte),
		mti.Equal(MTIReversalResponseByte),
		mti.Equal(MTIRepeatedReversalRequestByte),
		mti.Equal(MTINMMRequestByte),
		mti.Equal(MTINMMResponseByte):
		m.MTI = mti
	default:
		return ErrNotDefaultMti
	}
	m.cursor += 4

	if len(data[m.cursor:]) < 16 {
		return ErrInsufficientDataFirstBitmap
	}

	asciiHex := data[m.cursor : m.cursor+16]
	bitmap, err := hex.DecodeString(string(asciiHex))
	if err != nil {
		return ErrInvalidBitMap
	}
	m.firstBitmap = bitmap
	m.cursor += 16

	if m.isoMessageMap == nil {
		m.isoMessageMap = make(map[int][]byte)
	}

	if err := m.parseBitmap(bitmap); err != nil {
		return err
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

// CreateReturnISO to create Message Error
// TODO : MODIFY TO MAKE IT MORE PROPER
//func (m *Message) CreateReturnISO(i Message, RC string) (Message, error) {
//	// create msg
//	err := m.SetMTIResponse()
//	if err != nil {
//		return i, err
//	}
//	i.[39] = RC
//	err = i.PackISO()
//	if err != nil {
//		return i, err
//	}
//	return i, nil
//}

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
