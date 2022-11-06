package mpegts

import "time"

// PCR is Program Clock Reference
type PCR uint64

const (
	NonPcr PCR = (1 << 33) * 300
	MaxPcr PCR = NonPcr - 1
)

const (
	ProgramClock = 27000000 // 27MHz
)

// HasPCR returns true if PCR flag is set in the Adaptation Field.
func (p TS) HasPCR() bool {
	return (p[5] & 0x10) != 0
}

// SetPCR sets PCR flag and PCR value in the Adaptation Field.
func (p TS) SetPCR(value PCR) {
	p[5] |= 0x10 // PCR_flag

	pcrBase := value / 300
	pcrExt := value - (pcrBase * 300)

	p[6] = byte(pcrBase >> 25)
	p[7] = byte(pcrBase >> 17)
	p[8] = byte(pcrBase >> 9)
	p[9] = byte(pcrBase >> 1)
	p[10] = (byte((pcrBase << 7) & 0x80)) | 0x7E | (byte((pcrExt >> 8) & 0x01))
	p[11] = byte(pcrExt)
}

// PCR returns PCR value from the Adaptation Field.
// Packet should be with Adaptation Field
func (p TS) PCR() PCR {
	pcrBase := (PCR(p[6]) << 25) |
		(PCR(p[7]) << 17) |
		(PCR(p[8]) << 9) |
		(PCR(p[9]) << 1) |
		(PCR(p[10]) >> 7)
	pcrExt := (PCR((p[10] & 1)) << 8) | PCR(p[11])

	return (pcrBase * 300) + pcrExt
}

// Delta returns the difference p-u considering value overflow
func (p PCR) Delta(u PCR) PCR {
	if p >= u {
		return p - u
	} else {
		return NonPcr - u + p
	}
}

// Bitrate returns bitrate in bits per second for delta PCR.
func (p PCR) Bitrate(bytes int) int {
	return int((uint64(bytes) * 8 * ProgramClock) / uint64(p))
}

// Add returns the timestamp p+u
func (p PCR) Add(u PCR) PCR {
	return (p + u) & MaxPcr
}

// EstimatedPCR returns estimated PCR value
//
//	| time -->
//	| X---------X---------X
//	|  \         \         \
//	|   \         \         estimated PCR
//	|    \         current PCR
//	|     previous PCR
//
// - lastBlock - bytes between PCR(previous) and PCR(current)
// - currentBlock - bytes between PCR(current) and PCR(estimated)
func (p PCR) EstimatedPCR(previous PCR, lastBlock, currentBlock uint64) PCR {
	delta := uint64(p.Delta(previous))
	stc := PCR(delta * currentBlock / lastBlock)
	return stc.Add(p)
}

// Jitter returns the difference between two PCR values in nanoseconds
func (p PCR) Jitter(previous PCR) time.Duration {
	delta := p.Delta(previous)
	return time.Duration(delta * 1000 / 27)
}
