package tlv

import (
	"fmt"
	"slices"
)

func (t *Data) Append(tag uint32, v []byte) error {
	if tag == 0 {
		return fmt.Errorf("tag cannot be 0")
	}
	if v == nil {
		return fmt.Errorf("append data cannot be nil")
	}

	t.list = append(t.list, TagData{
		tag:   tag,
		value: v,
	})

	return nil
}

func (t *Data) Remove(tag uint32) error {
	i := slices.IndexFunc(t.list, func(d TagData) bool {
		if d.tag == tag {
			return true
		}
		return false
	})
	if i == -1 {
		return fmt.Errorf("tag not found")
	}

	t.list = append(t.list[:i], t.list[i+1:]...)

	return nil
}
