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

// EncodeISO8859_1 converts an UTF-8 string into ISO-8859-1 (Western European)
func EncodeISO8859_1(src string) []byte {
	return encodeISO8859(src, hi_map_1, encode_map_1)
}

// EncodeISO8859_2 converts an UTF-8 string into ISO-8859-2 (Central European)
func EncodeISO8859_2(src string) []byte {
	return encodeISO8859(src, hi_map_2, encode_map_2)
}

// EncodeISO8859_3 converts an UTF-8 string into ISO-8859-3 (South European)
func EncodeISO8859_3(src string) []byte {
	return encodeISO8859(src, hi_map_3, encode_map_3)
}

// EncodeISO8859_4 converts an UTF-8 string into ISO-8859-4 (North European)
func EncodeISO8859_4(src string) []byte {
	return encodeISO8859(src, hi_map_4, encode_map_4)
}

// EncodeISO8859_5 converts an UTF-8 string into ISO-8859-5 (Cyrillic)
func EncodeISO8859_5(src string) []byte {
	return encodeISO8859(src, hi_map_5, encode_map_5)
}

// EncodeISO8859_6 converts an UTF-8 string into ISO-8859-6 (Arabic)
func EncodeISO8859_6(src string) []byte {
	return encodeISO8859(src, hi_map_6, encode_map_6)
}

// EncodeISO8859_7 converts an UTF-8 string into ISO-8859-7 (Greek)
func EncodeISO8859_7(src string) []byte {
	return encodeISO8859(src, hi_map_7, encode_map_7)
}

// EncodeISO8859_8 converts an UTF-8 string into ISO-8859-8 (Hebrew)
func EncodeISO8859_8(src string) []byte {
	return encodeISO8859(src, hi_map_8, encode_map_8)
}

// EncodeISO8859_9 converts an UTF-8 string into ISO-8859-9 (Turkish)
func EncodeISO8859_9(src string) []byte {
	return encodeISO8859(src, hi_map_9, encode_map_9)
}

// EncodeISO8859_10 converts an UTF-8 string into ISO-8859-10 (Nordic)
func EncodeISO8859_10(src string) []byte {
	return encodeISO8859(src, hi_map_10, encode_map_10)
}

// EncodeISO8859_11 converts an UTF-8 string into ISO-8859-11 (Thai)
func EncodeISO8859_11(src string) []byte {
	return encodeISO8859(src, hi_map_11, encode_map_11)
}

// EncodeISO8859_13 converts an UTF-8 string into ISO-8859-13 (Baltic Rim)
func EncodeISO8859_13(src string) []byte {
	return encodeISO8859(src, hi_map_13, encode_map_13)
}

// EncodeISO8859_14 converts an UTF-8 string into ISO-8859-14 (Celtic)
func EncodeISO8859_14(src string) []byte {
	return encodeISO8859(src, hi_map_14, encode_map_14)
}

// EncodeISO8859_15 converts an UTF-8 string into ISO-8859-15 (Western European)
func EncodeISO8859_15(src string) []byte {
	return encodeISO8859(src, hi_map_15, encode_map_15)
}

// EncodeISO8859_16 converts an UTF-8 string into ISO-8859-16 (Romanian)
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

// DecodeISO8859_1 converts ISO-8859-1 into UTF-8
func DecodeISO8859_1(src []byte) string {
	return decodeISO8859(src, decode_map_1)
}

// DecodeISO8859_2 converts ISO-8859-2 into UTF-8
func DecodeISO8859_2(src []byte) string {
	return decodeISO8859(src, decode_map_2)
}

// DecodeISO8859_3 converts ISO-8859-3 into UTF-8
func DecodeISO8859_3(src []byte) string {
	return decodeISO8859(src, decode_map_3)
}

// DecodeISO8859_4 converts ISO-8859-4 into UTF-8
func DecodeISO8859_4(src []byte) string {
	return decodeISO8859(src, decode_map_4)
}

// DecodeISO8859_5 converts ISO-8859-5 into UTF-8
func DecodeISO8859_5(src []byte) string {
	return decodeISO8859(src, decode_map_5)
}

// DecodeISO8859_6 converts ISO-8859-6 into UTF-8
func DecodeISO8859_6(src []byte) string {
	return decodeISO8859(src, decode_map_6)
}

// DecodeISO8859_7 converts ISO-8859-7 into UTF-8
func DecodeISO8859_7(src []byte) string {
	return decodeISO8859(src, decode_map_7)
}

// DecodeISO8859_8 converts ISO-8859-8 into UTF-8
func DecodeISO8859_8(src []byte) string {
	return decodeISO8859(src, decode_map_8)
}

// DecodeISO8859_9 converts ISO-8859-9 into UTF-8
func DecodeISO8859_9(src []byte) string {
	return decodeISO8859(src, decode_map_9)
}

// DecodeISO8859_10 converts ISO-8859-10 into UTF-8
func DecodeISO8859_10(src []byte) string {
	return decodeISO8859(src, decode_map_10)
}

// DecodeISO8859_11 converts ISO-8859-11 into UTF-8
func DecodeISO8859_11(src []byte) string {
	return decodeISO8859(src, decode_map_11)
}

// DecodeISO8859_13 converts ISO-8859-13 into UTF-8
func DecodeISO8859_13(src []byte) string {
	return decodeISO8859(src, decode_map_13)
}

// DecodeISO8859_14 converts ISO-8859-14 into UTF-8
func DecodeISO8859_14(src []byte) string {
	return decodeISO8859(src, decode_map_14)
}

// DecodeISO8859_15 converts ISO-8859-15 into UTF-8
func DecodeISO8859_15(src []byte) string {
	return decodeISO8859(src, decode_map_15)
}

// DecodeISO8859_16 converts ISO-8859-16 into UTF-8
func DecodeISO8859_16(src []byte) string {
	return decodeISO8859(src, decode_map_16)
}
