package iso8583

import (
	"encoding/hex"
	"errors"
	"fmt"
)

// ValidateMandatoryBits validate ISO Message
func (m *Message) ValidateMandatoryBits() error {
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

func (m *Message) ValidateBitType() (err error) {
	for bit, v := range m.isoMessageMap {
		bitType := m.packager.IsoPackagerConfig[bit].Type
		val := string(v)
		switch bitType {
		case BitTypeN:
			if !reNumeric.MatchString(val) {
				err = errors.Join(err, ErrInvalidValue, fmt.Errorf("invalid type bit %d type %s got %s", bit, bitType, val))
			}
		case BitTypeAN:
			if !reAlphaNum.MatchString(val) {
				err = errors.Join(err, ErrInvalidValue, fmt.Errorf("invalid type bit %d type %s got %s", bit, bitType, val))
			}
		case BitTypeANS:
			// accept everything
			return nil
		case BitTypeB:
			_, errConv := hex.DecodeString(val)
			if errConv != nil {
				err = errors.Join(err, ErrInvalidValue, fmt.Errorf("invalid type bit %d type %s got %s", bit, bitType, val))
			}
		case BitTypeZ:
			if !reTrackData.MatchString(val) {
				err = errors.Join(err, ErrInvalidValue, fmt.Errorf("invalid type bit %d type %s got %s", bit, bitType, val))
			}
		default:
			err = errors.Join(err, ErrInvalidBitType)
		}
	}

	return err
}
