package textcode

func EncodeISO6937(src string) []byte {
	var result []byte

	for _, r := range src {
		c := uint16(r)

		if c <= 0x7F {
			result = append(result, byte(c))
			continue
		}

		hi := int(c >> 8)
		lo := int(c & 0xFF)

		pos := (hi_map_iso6937[hi] * 0x100) + lo
		code := encode_map_iso6937[pos]

		switch {
		case code > 0xFF:
			result = append(result, (byte(code >> 8)), (byte(code & 0xFF)))
		case code > 0:
			result = append(result, byte(code))
		default:
			result = append(result, '?')
		}
	}

	return result
}

func DecodeISO6937(src []byte) string {
	var result []rune

	skip := 0
	for skip < len(src) {
		b := src[skip]
		skip += 1

		if b <= 0x7F {
			result = append(result, rune(b))
			continue
		}

		var mapSkip int

		// diactrics
		if (b >= 0xC1) && (b <= 0xCF) {
			if skip >= len(src) {
				result = append(result, '�')
				break
			}

			// 96 - bytes before 0xC1 section
			// 58 - bytes for each diactric section
			mapSkip = 96 + ((int(b - 0xC1)) * 58)
			b = src[skip]
			skip += 1

			if (b >= 'A') && (b <= ('A' + 58 - 1)) {
				mapSkip += int(b - 'A')
			} else {
				mapSkip = 6 // position of 0x0000
			}
		} else {
			mapSkip = int(b - 0xA0)
		}

		if r := decode_map_iso6937[mapSkip]; r != 0x0000 {
			result = append(result, rune(r))
		} else {
			result = append(result, '�')
		}
	}

	return string(result)
}
