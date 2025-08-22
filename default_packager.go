package iso8583

func DefaultPackager() *IsoPackager {
	packager := &IsoPackager{
		HasHeader:    false,
		HeaderLength: 0,
		MessageKey:   []int{2, 7, 11, 12, 13, 41, 37},
		IsoPackagerConfig: map[int]BitConfig{
			1:   NewBitConfigFixed(true, BitTypeB, 16),
			2:   NewBitConfigLLVar(true, BitTypeN, 19),
			3:   NewBitConfigFixed(false, BitTypeANS, 6),
			4:   NewBitConfigFixed(false, BitTypeANS, 12),
			5:   NewBitConfigFixed(false, BitTypeANS, 12),
			6:   NewBitConfigFixed(false, BitTypeANS, 12),
			7:   NewBitConfigFixed(true, BitTypeANS, 10),
			8:   NewBitConfigFixed(false, BitTypeANS, 8),
			9:   NewBitConfigFixed(false, BitTypeANS, 8),
			10:  NewBitConfigFixed(false, BitTypeANS, 8),
			11:  NewBitConfigFixed(true, BitTypeANS, 6),
			12:  NewBitConfigFixed(true, BitTypeANS, 6),
			13:  NewBitConfigFixed(true, BitTypeANS, 4),
			14:  NewBitConfigFixed(false, BitTypeANS, 4),
			15:  NewBitConfigFixed(false, BitTypeANS, 4),
			16:  NewBitConfigFixed(false, BitTypeANS, 4),
			17:  NewBitConfigFixed(false, BitTypeANS, 4),
			18:  NewBitConfigFixed(false, BitTypeANS, 4),
			19:  NewBitConfigFixed(false, BitTypeANS, 4),
			20:  NewBitConfigFixed(false, BitTypeANS, 4),
			21:  NewBitConfigFixed(false, BitTypeANS, 3),
			22:  NewBitConfigFixed(false, BitTypeANS, 3),
			23:  NewBitConfigFixed(false, BitTypeANS, 3),
			24:  NewBitConfigFixed(false, BitTypeANS, 3),
			25:  NewBitConfigFixed(false, BitTypeANS, 2),
			26:  NewBitConfigFixed(false, BitTypeANS, 2),
			27:  NewBitConfigFixed(false, BitTypeANS, 3),
			28:  NewBitConfigFixed(false, BitTypeANS, 9),
			29:  NewBitConfigFixed(false, BitTypeANS, 3),
			30:  NewBitConfigFixed(false, BitTypeANS, 3),
			31:  NewBitConfigLLVar(false, BitTypeANS, 99),
			32:  NewBitConfigLLVar(false, BitTypeANS, 99),
			33:  NewBitConfigLLVar(false, BitTypeANS, 99),
			34:  NewBitConfigLLVar(false, BitTypeANS, 99),
			35:  NewBitConfigLLVar(false, BitTypeZ, 99),
			36:  NewBitConfigLLVar(false, BitTypeANS, 99),
			37:  NewBitConfigFixed(true, BitTypeANS, 12),
			38:  NewBitConfigFixed(false, BitTypeANS, 6),
			39:  NewBitConfigFixed(false, BitTypeANS, 2),
			40:  NewBitConfigLLVar(false, BitTypeANS, 99),
			41:  NewBitConfigFixed(false, BitTypeANS, 16),
			42:  NewBitConfigFixed(false, BitTypeANS, 15),
			43:  NewBitConfigFixed(false, BitTypeANS, 40),
			44:  NewBitConfigLLVar(false, BitTypeANS, 99),
			45:  NewBitConfigLLVar(false, BitTypeANS, 99),
			46:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			47:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			48:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			49:  NewBitConfigFixed(false, BitTypeANS, 3),
			50:  NewBitConfigFixed(false, BitTypeANS, 3),
			51:  NewBitConfigFixed(false, BitTypeANS, 3),
			52:  NewBitConfigFixed(false, BitTypeANS, 16),
			53:  NewBitConfigFixed(false, BitTypeANS, 16),
			54:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			55:  NewBitConfigLLLVar(false, BitTypeB, 999),
			56:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			57:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			58:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			59:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			60:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			61:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			62:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			63:  NewBitConfigLLLVar(false, BitTypeANS, 999),
			64:  NewBitConfigFixed(false, BitTypeANS, 16),
			65:  NewBitConfigFixed(false, BitTypeANS, 1),
			66:  NewBitConfigFixed(false, BitTypeANS, 1),
			67:  NewBitConfigFixed(false, BitTypeANS, 2),
			68:  NewBitConfigFixed(false, BitTypeANS, 3),
			69:  NewBitConfigFixed(false, BitTypeANS, 3),
			70:  NewBitConfigFixed(false, BitTypeANS, 3),
			71:  NewBitConfigFixed(false, BitTypeANS, 4),
			72:  NewBitConfigFixed(false, BitTypeANS, 4),
			73:  NewBitConfigFixed(false, BitTypeANS, 6),
			74:  NewBitConfigFixed(false, BitTypeANS, 10),
			75:  NewBitConfigFixed(false, BitTypeANS, 10),
			76:  NewBitConfigFixed(false, BitTypeANS, 10),
			77:  NewBitConfigFixed(false, BitTypeANS, 10),
			78:  NewBitConfigFixed(false, BitTypeANS, 10),
			79:  NewBitConfigFixed(false, BitTypeANS, 10),
			80:  NewBitConfigFixed(false, BitTypeANS, 10),
			81:  NewBitConfigFixed(false, BitTypeANS, 10),
			82:  NewBitConfigFixed(false, BitTypeANS, 12),
			83:  NewBitConfigFixed(false, BitTypeANS, 12),
			84:  NewBitConfigFixed(false, BitTypeANS, 12),
			85:  NewBitConfigFixed(false, BitTypeANS, 12),
			86:  NewBitConfigFixed(false, BitTypeANS, 16),
			87:  NewBitConfigFixed(false, BitTypeANS, 16),
			88:  NewBitConfigFixed(false, BitTypeANS, 16),
			89:  NewBitConfigFixed(false, BitTypeANS, 16),
			90:  NewBitConfigFixed(false, BitTypeANS, 42),
			91:  NewBitConfigFixed(false, BitTypeANS, 1),
			92:  NewBitConfigFixed(false, BitTypeANS, 2),
			93:  NewBitConfigFixed(false, BitTypeANS, 5),
			94:  NewBitConfigFixed(false, BitTypeANS, 7),
			95:  NewBitConfigFixed(false, BitTypeANS, 42),
			96:  NewBitConfigFixed(false, BitTypeANS, 16),
			97:  NewBitConfigFixed(false, BitTypeANS, 16),
			98:  NewBitConfigFixed(false, BitTypeANS, 25),
			99:  NewBitConfigLLVar(false, BitTypeANS, 99),
			100: NewBitConfigLLVar(false, BitTypeANS, 99),
			101: NewBitConfigLLVar(false, BitTypeANS, 99),
			102: NewBitConfigLLVar(false, BitTypeANS, 99),
			103: NewBitConfigLLLVar(false, BitTypeANS, 999),
			104: NewBitConfigLLLVar(false, BitTypeANS, 999),
			105: NewBitConfigLLLVar(false, BitTypeANS, 999),
			106: NewBitConfigLLLVar(false, BitTypeANS, 999),
			107: NewBitConfigLLLVar(false, BitTypeANS, 999),
			108: NewBitConfigLLLVar(false, BitTypeANS, 999),
			109: NewBitConfigLLLVar(false, BitTypeANS, 999),
			110: NewBitConfigLLLVar(false, BitTypeANS, 999),
			111: NewBitConfigLLLVar(false, BitTypeANS, 999),
			112: NewBitConfigLLLVar(false, BitTypeANS, 999),
			113: NewBitConfigLLLVar(false, BitTypeANS, 999),
			114: NewBitConfigLLLVar(false, BitTypeANS, 999),
			115: NewBitConfigLLLVar(false, BitTypeANS, 999),
			116: NewBitConfigLLLVar(false, BitTypeANS, 999),
			117: NewBitConfigLLLVar(false, BitTypeANS, 999),
			118: NewBitConfigLLLVar(false, BitTypeANS, 999),
			119: NewBitConfigLLLVar(false, BitTypeANS, 999),
			120: NewBitConfigLLLVar(false, BitTypeANS, 999),
			121: NewBitConfigLLLVar(false, BitTypeANS, 999),
			122: NewBitConfigLLLVar(false, BitTypeANS, 999),
			123: NewBitConfigLLLVar(false, BitTypeANS, 999),
			124: NewBitConfigLLLVar(false, BitTypeANS, 999),
			125: NewBitConfigLLLVar(false, BitTypeANS, 999),
			126: NewBitConfigLLLVar(false, BitTypeANS, 999),
			127: NewBitConfigLLLVar(false, BitTypeANS, 999),
			128: NewBitConfigFixed(false, BitTypeANS, 16),
		},
	}

	packager.MandatoryBit = packager.GetMandatoryBitsFromConfig()

	return packager
}

func NewBitConfigFixed(isMandatory bool, bitType BitType, length int) BitConfig {
	return BitConfig{
		IsMandatory: isMandatory,
		Type:        bitType,
		Length: BitLength{
			Type: LengthTypeFixed,
			Max:  length,
		},
	}
}

func NewBitConfigLLVar(isMandatory bool, bitType BitType, length int) BitConfig {
	if length > 99 {
		panic("length must be less than 99")
	}
	return BitConfig{
		IsMandatory: isMandatory,
		Type:        bitType,
		Length: BitLength{
			Type: LengthTypeLLVar,
			Max:  length,
		},
	}
}

func NewBitConfigLLLVar(isMandatory bool, bitType BitType, length int) BitConfig {
	if length > 999 {
		panic("length must be less than 999")
	}
	return BitConfig{
		IsMandatory: isMandatory,
		Type:        bitType,
		Length: BitLength{
			Type: LengthTypeLLLVar,
			Max:  length,
		},
	}
}

func NewBitConfigLLLLVar(isMandatory bool, bitType BitType, length int) BitConfig {
	if length > 9999 {
		panic("length must be less than 9999")
	}
	return BitConfig{
		IsMandatory: isMandatory,
		Type:        bitType,
		Length: BitLength{
			Type: LengthTypeLLLLVar,
			Max:  length,
		},
	}
}
