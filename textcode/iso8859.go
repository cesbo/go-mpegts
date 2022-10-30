package textcode

func encodeISO8859(src string, hiMap []int, m []uint8) []byte {
	var result []byte

	for _, r := range src {
		c := uint16(r)

		if c <= 0x7F {
			result = append(result, byte(c))
			continue
		}

		hi := int(c >> 8)
		lo := int(c & 0xFF)
		pos := (hiMap[hi] * 0x100) + lo

		code := uint8(0)
		if pos < len(m) {
			code = m[pos]
		}

		if code != 0x00 {
			result = append(result, code)
		} else {
			result = append(result, '?')
		}
	}

	return result
}

func EncodeISO8859_1(src string) []byte {
	return encodeISO8859(src, hi_map_1, encode_map_1)
}

func EncodeISO8859_2(src string) []byte {
	return encodeISO8859(src, hi_map_2, encode_map_2)
}

func EncodeISO8859_3(src string) []byte {
	return encodeISO8859(src, hi_map_3, encode_map_3)
}

func EncodeISO8859_4(src string) []byte {
	return encodeISO8859(src, hi_map_4, encode_map_4)
}

func EncodeISO8859_5(src string) []byte {
	return encodeISO8859(src, hi_map_5, encode_map_5)
}

func EncodeISO8859_6(src string) []byte {
	return encodeISO8859(src, hi_map_6, encode_map_6)
}

func EncodeISO8859_7(src string) []byte {
	return encodeISO8859(src, hi_map_7, encode_map_7)
}

func EncodeISO8859_8(src string) []byte {
	return encodeISO8859(src, hi_map_8, encode_map_8)
}

func EncodeISO8859_9(src string) []byte {
	return encodeISO8859(src, hi_map_9, encode_map_9)
}

func EncodeISO8859_10(src string) []byte {
	return encodeISO8859(src, hi_map_10, encode_map_10)
}

func EncodeISO8859_11(src string) []byte {
	return encodeISO8859(src, hi_map_11, encode_map_11)
}

func EncodeISO8859_13(src string) []byte {
	return encodeISO8859(src, hi_map_13, encode_map_13)
}

func EncodeISO8859_14(src string) []byte {
	return encodeISO8859(src, hi_map_14, encode_map_14)
}

func EncodeISO8859_15(src string) []byte {
	return encodeISO8859(src, hi_map_15, encode_map_15)
}

func EncodeISO8859_16(src string) []byte {
	return encodeISO8859(src, hi_map_16, encode_map_16)
}

func decodeISO8859(src []byte, m []uint16) string {
	var result []rune

	for _, c := range src {
		if c <= 0x7F {
			result = append(result, rune(c))
			continue
		}

		if c >= 0xA0 {
			c -= 0xA0

			if r := m[c]; r != 0 {
				result = append(result, rune(r))
				continue
			}
		}

		result = append(result, '�')
	}

	return string(result)
}

func DecodeISO8859_1(src []byte) string {
	return decodeISO8859(src, decode_map_1)
}

func DecodeISO8859_2(src []byte) string {
	return decodeISO8859(src, decode_map_2)
}

func DecodeISO8859_3(src []byte) string {
	return decodeISO8859(src, decode_map_3)
}

func DecodeISO8859_4(src []byte) string {
	return decodeISO8859(src, decode_map_4)
}

func DecodeISO8859_5(src []byte) string {
	return decodeISO8859(src, decode_map_5)
}

func DecodeISO8859_6(src []byte) string {
	return decodeISO8859(src, decode_map_6)
}

func DecodeISO8859_7(src []byte) string {
	return decodeISO8859(src, decode_map_7)
}

func DecodeISO8859_8(src []byte) string {
	return decodeISO8859(src, decode_map_8)
}

func DecodeISO8859_9(src []byte) string {
	return decodeISO8859(src, decode_map_9)
}

func DecodeISO8859_10(src []byte) string {
	return decodeISO8859(src, decode_map_10)
}

func DecodeISO8859_11(src []byte) string {
	return decodeISO8859(src, decode_map_11)
}

func DecodeISO8859_13(src []byte) string {
	return decodeISO8859(src, decode_map_13)
}

func DecodeISO8859_14(src []byte) string {
	return decodeISO8859(src, decode_map_14)
}

func DecodeISO8859_15(src []byte) string {
	return decodeISO8859(src, decode_map_15)
}

func DecodeISO8859_16(src []byte) string {
	return decodeISO8859(src, decode_map_16)
}
