package tlv

func bytesToUint32(b []byte) uint32 {
	var n uint32
	for _, v := range b {
		n = (n << 8) | uint32(v)
	}
	return n
}

func uint32ToBytes(n uint32) []byte {
	if n == 0 {
		return nil
	}
	//  EXAMPLE
	// n := uint32(40706) //  0x00009F02
	// 00 00 9F 02
	// b3 := byte(n >> 24) //  0x00009F02 >> 24 = 0x00 → 0x00 (0)
	// b2 := byte(n >> 16) //  0x00009F02 >> 16 = 0x00 → 0x00 (0)
	// b1 := byte(n >> 8)  //  0x00009F02 >>  8 = 0x9F → 0x9F (159)
	// b0 := byte(n)       //                        = 0x02 (2)
	//
	// bytes := []byte{b3, b2, b1, b0} //  []byte{0x00, 0x00, 0x9F, 0x02}
	//

	// makes 4 byte of uint32
	b := []byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}

	// strip leading zeros
	i := 0
	for i < len(b) && b[i] == 0 {
		i++
	}
	return b[i:]

}
