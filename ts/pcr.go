package ts

type PCR uint64

const (
	NonPcr PCR = (1 << 33) * 300
	MaxPcr PCR = NonPcr - 1
)

func SetPCR(dst TS, value PCR) {
	dst[5] |= 0x10 // PCR_flag

	pcrBase := value / 300
	pcrExt := value % 300

	dst[6] = byte(pcrBase >> 25)
	dst[7] = byte(pcrBase >> 17)
	dst[8] = byte(pcrBase >> 9)
	dst[9] = byte(pcrBase >> 1)
	dst[10] = (byte((pcrBase << 7) & 0x80)) | 0x7E | (byte((pcrExt >> 8) & 0x01))
	dst[11] = byte(pcrExt)
}

func GetPCR(src TS) PCR {
	pcrBase := (PCR(src[6]) << 25) |
		(PCR(src[7]) << 17) |
		(PCR(src[8]) << 9) |
		(PCR(src[9]) << 1) |
		(PCR(src[10]) >> 7)
	pcrExt := (PCR((src[10] & 1)) << 8) | PCR(src[11])

	return (pcrBase * 300) + pcrExt
}
