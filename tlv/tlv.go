package tlv

import (
	"encoding/hex"
	"fmt"
)

type Data struct {
	list []TagData
}

type TagData struct {
	tag   uint32
	value []byte
}

func New(data []byte) (*Data, error) {
	result := &Data{
		list: make([]TagData, 0),
	}

	if data == nil {
		return result, nil
	}

	i := 0

	for i < len(data) {
		// --- Parse Tag ---
		tag := []byte{data[i]}
		i++

		// Multi-byte tag check: if 5 LSBs of first tag byte are all 1s (0x1F)
		// In EMV, if the lower 5 bits of the first byte are all 1 (0x1F), the tag extends into more bytes.
		// 0x1F is 00011111 in binary
		// use AND operator to check the tag
		if tag[0]&0x1F == 0x1F {
			// Read continuation bytes until MSB = 0
			for {
				if i >= len(data) {
					return nil, fmt.Errorf("unexpected end while reading tag")
				}
				tag = append(tag, data[i])
				i++
				// 0x80 is 10000000 in binary
				// use AND operator if the value is 00000000 then break
				// we found the tag if its MSB = 0
				if data[i]&0x80 == 0 {
					break
				}
			}
		}

		// --- Parse Length ---
		if i >= len(data) {
			return nil, fmt.Errorf("unexpected end while reading length")
		}
		length := int(data[i])
		i++

		if length&0x80 != 0 { // Long form length
			numBytes := length & 0x7F
			if i+numBytes > len(data) {
				return nil, fmt.Errorf("invalid length encoding")
			}
			length = 0
			for j := 0; j < numBytes; j++ {
				length = (length << 8) | int(data[i])
				i++
			}
		}

		// --- Parse Value ---
		if i+length > len(data) {
			return nil, fmt.Errorf("value exceeds available data")
		}
		value := data[i : i+length]
		i += length

		// Store in map (tag as string)
		tagKey := bytesToUint32(tag)

		result.list = append(result.list, TagData{
			tag:   tagKey,
			value: value,
		})
	}

	return result, nil
}

func (t *Data) HasTag(tag uint32) bool {
	for _, k := range t.list {
		if k.tag == tag {
			return true
		}
	}
	return false
}

func (t *Data) GetByte(tag uint32) []byte {
	for _, k := range t.list {
		if k.tag == tag {
			return k.value
		}
	}
	return nil
}

func (t *Data) GetInt64(tag uint32) int64 {
	b := t.GetByte(tag)
	if b == nil {
		return 0
	}
	var result int64
	for _, v := range b {
		result = result*100 + int64((v>>4)&0xF)*10 + int64(v&0xF)
	}
	return result
}

func (t *Data) GetHexString(tag uint32) string {
	b := t.GetByte(tag)
	if b == nil {
		return ""
	}

	return hex.EncodeToString(b)
}
