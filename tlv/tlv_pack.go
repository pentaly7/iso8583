package tlv

func (t *Data) Pack() []byte {
	if t.list == nil || len(t.list) == 0 {
		return nil
	}
	result := make([]byte, 0)
	for _, v := range t.list {
		length := len(v.value)
		result = append(result, uint32ToBytes(v.tag)...)
		result = append(result, byte(length))
		result = append(result, v.value...)
	}

	return result
}
