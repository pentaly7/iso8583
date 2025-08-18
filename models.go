package iso8583

import (
	"fmt"
	"strconv"
)

type (
	// IsoPackage for PackagerISO message
	IsoPackage struct {
		IsoHeaderActive bool        `json:"IsoHeaderActive"`
		IsoHeaderLength int         `json:"IsoHeaderLength"`
		IsoKey          []int       `json:"IsoKey"`
		IsoPackager     []int       `json:"IsoPackager"`
		IsoPackagerMap  map[int]int `json:"IsoPackagerMap"`
		MandatoryBits   []int       `json:"MandatoryBit"`
	}

	// Message for Component Message
	Message struct {
		IsoHeader       string         // iso header
		MTI             string         // MTI
		firstBitmap     string         // First bitmap
		actSecondBitmap bool           // Second Bitmap activation
		secondBitmap    string         // Second bitmap
		AsString        string         // ASCII full data
		ErrorIso        error          // info Error to take
		IsoMessageMap   map[int]string // Get Element of Iso Message in Map
		BitActive       []int          // Find Bit Active for check Mandatory Bit
		IsoPack         IsoPackage     // Set Iso Package for each ISO Specification
		runeData        []rune
	}
)

func New(packager IsoPackage) *Message {
	return &Message{
		IsoPack:       packager,
		IsoMessageMap: make(map[int]string),
	}
}

// DefaultISOPackage for set default ISO Package if no packager used
func DefaultISOPackage() (i IsoPackage) {
	var (
		isoHeaderLength = 12
		isoKey          = []int{2, 11, 12, 13, 37, 41}
		isoPackager     = []int{ // Default Package ISO if no Package
			// 	1   2   3   4   5   6   7   8   9   10
			16, -2, 6, 12, 12, 12, 10, 8, 8, 8, // bit 1 - 10
			6, 6, 4, 4, 4, 4, 4, 4, 4, 4, // bit 11 - 20
			3, 3, 3, 3, 2, 2, 3, 9, 3, 3, // bit 21 - 30
			-2, -2, -2, -2, -2, -2, 12, 6, 2, -2, // bit 31 - 40
			16, 15, 40, -2, -2, -3, -3, -3, 3, 3, // bit 41 - 50
			3, 16, 16, -3, -3, -3, -3, -3, -3, -3, // bit 51 - 60
			-3, -3, -3, 16, 1, 1, 2, 3, 3, 3, // bit 61 - 70
			4, 4, 6, 10, 10, 10, 10, 10, 10, 10, // bit 71 - 80
			10, 12, 12, 12, 12, 16, 16, 16, 16, 42, // bit 81 - 90
			1, 2, 5, 7, 42, 16, 16, 25, -2, -2, // bit 91 - 100
			-2, -2, -2, -3, -3, -3, -3, -3, -3, -3, // bit 101 - 110
			-3, -3, -3, -3, -3, -3, -3, -3, -3, -3, // bit 111 - 120
			-3, -3, -3, -3, -3, -3, -3, 16, // bit 121 - 128
		}
		mandatoryBits = []int{2, 3, 4, 7, 11, 12, 13, 15, 18, 32, 37, 42, 48, 49, 63}
	)
	i.IsoHeaderActive = false
	i.IsoHeaderLength = isoHeaderLength
	i.IsoKey = isoKey
	i.IsoPackager = isoPackager
	i.MandatoryBits = mandatoryBits
	return i
}

func FlattenIsoPackagerConfig(m map[int]int) []int {
	// Find max key to determine slice length
	maxKey := 0
	for k := range m {
		if k > maxKey {
			maxKey = k
		}
	}

	// Initialize slice with zero values
	s := make([]int, maxKey)

	// Fill slice based on map keys
	for k, v := range m {
		s[k-1] = v
	}

	return s
}

func (m *Message) Set(bit int, value string) {
	m.IsoMessageMap[bit] = value
}

func (m *Message) Get(bit int) string {
	return m.IsoMessageMap[bit]
}

// ParseISO to Single Data Element
func (m *Message) ParseISO() error {
	isoMsg, err := m.checkMessageISO()
	if err != nil {
		return ErrNotIsoMessage
	}

	m.runeData = []rune(isoMsg)
	if len(m.runeData) < 4 {
		return ErrInsufficientDataMti
	}
	m.MTI = string(m.runeData[:4])
	copy(m.runeData, m.runeData[4:])

	if len(m.runeData) < 16 {
		return ErrInsufficientDataFirstBitmap
	}
	m.firstBitmap = string(m.runeData[:16])

	copy(m.runeData, m.runeData[16:])

	firstBitmap, err := strconv.ParseUint(m.firstBitmap, 16, 64)
	if err != nil {
		return ErrParsingFirstBitmap
	}

	m.IsoMessageMap = make(map[int]string)

	if err := m.parseBitmap(firstBitmap, 0); err != nil {
		return err
	}

	if m.actSecondBitmap {
		secondBitmap, err := strconv.ParseUint(m.secondBitmap, 16, 64)
		if err != nil {
			return ErrParsingSecondBitmap
		}

		if err := m.parseBitmap(secondBitmap, 64); err != nil {
			return err
		}
	}

	err = m.GetKeyActive()
	if err != nil {
		return err
	}
	return nil
}

// CreateTraceISO is to Trace ISO Message created for Tracing ISO Message Respon
func CreateTraceISO(isoData Message) string {
	var result string
	// combination ISO => MTI(respon), PAN , STAN, RRN, TransmissionDateTime
	if IsRequest(isoData) {
		isoData, _ = SetResponseMTI(isoData)
		result = isoData.MTI
	} else {
		result = isoData.MTI
	}
	isomap := isoData.IsoMessageMap
	for _, v := range isoData.IsoPack.IsoKey {
		result += isomap[v]
	}
	return result
}

// CreateReturnISOGeneral to create Message Error
func CreateReturnISOGeneral(i Message, RC string) (Message, error) {
	// create msg
	i, err := SetResponseMTI(i)
	if err != nil {
		return i, err
	}
	i.IsoMessageMap[39] = RC
	err = i.PackISO()
	if err != nil {
		return i, err
	}
	return i, nil
}

// SetResponseMTI is create MTI Response for Response Message
func SetResponseMTI(i Message) (Message, error) {
	if i.MTI == "0200" {
		i.MTI = "0210"
		return i, nil
	} else if i.MTI == "0400" || i.MTI == "0401" {
		i.MTI = "0410"
		return i, nil
	} else if i.MTI == "0800" {
		i.MTI = "0810"
		return i, nil
	} else if i.MTI == "0210" || i.MTI == "0410" || i.MTI == "0810" {
		return i, nil
	}
	return i, fmt.Errorf("cant found mti type")
}

// IsRequest for check MTI is Request or Response
func IsRequest(i Message) bool {
	if i.MTI == "0200" || i.MTI == "0400" || i.MTI == "0800" {
		return true
	}
	return false
}

// IsTransactional for check MTI message is Transaction or NMM
func IsTransactional(i Message) bool {
	if i.MTI == "0200" || i.MTI == "0400" || i.MTI == "0210" || i.MTI == "0410" {
		return true
	}
	return false
}

// IsNMM for check MTI message is Transaction or NMM
func IsNMM(i Message) bool {
	if i.MTI == "0800" || i.MTI == "0810" {
		return true
	}
	return false
}

// IsResponse for check MTI is Request or Response
func IsResponse(i Message) bool {
	if i.MTI == "0210" || i.MTI == "0410" || i.MTI == "0810" {
		return true
	}
	return false
}
