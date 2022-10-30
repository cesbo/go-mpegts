package textcode

func EncodeGB2312(src string) []byte {
	var result []byte

	for _, r := range src {
		c := uint16(r)

		if c <= 0x7F {
			result = append(result, byte(c))
			continue
		}

		hi := int(c >> 8)
		lo := int(c & 0xFF)
		pos := (hi_map_gb2312[hi] * 0x100) + lo
		code := encode_map_gb2312[pos]

		if code != 0x0000 {
			result = append(result, (byte(code >> 8)), (byte(code & 0xFF)))
		} else {
			result = append(result, '?')
		}
	}

	return result
}

func DecodeGB2312(src []byte) string {
	var result []rune

	const loSize = 0x7F - 0x21

	skip := 0
	for skip < len(src) {
		b := src[skip]
		skip += 1

		if b <= 0x7F {
			result = append(result, rune(b))
			continue
		}

		if skip >= len(src) {
			result = append(result, '�')
			break
		}

		hi := int(b & 0x7F)
		lo := int(src[skip] & 0x7F)
		skip += 1

		if (lo < 0x21) || (hi < 0x21) {
			result = append(result, '�')
			continue
		}

		mapSkip := ((hi - 0x21) * loSize) + (lo - 0x21)
		if mapSkip >= len(decode_map_gb2312) {
			mapSkip = 94 // position of 0x0000
		}

		if r := decode_map_gb2312[mapSkip]; r != 0x0000 {
			result = append(result, rune(r))
		} else {
			result = append(result, '�')
		}
	}

	return string(result)
}
