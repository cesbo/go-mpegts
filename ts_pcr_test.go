package mpegts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPCR_Set(t *testing.T) {
	type testData struct {
		pcr      PCR
		expected TS
	}

	list := []testData{
		{
			pcr:      86405647,
			expected: TS{0x10, 0x00, 0x02, 0x32, 0x89, 0x7E, 0xF7},
		},
		{
			pcr:      2268366350823,
			expected: TS{0x10, 0xE1, 0x57, 0x8A, 0x18, 0xFE, 0x7B},
		},
	}

	assert := assert.New(t)

	for i, x := range list {
		packet := NewTS(256)
		packet[3] |= 0x20 // has AF
		packet[4] = 7     // AF length

		packet.SetPCR(x.pcr)

		assert.Equal(x.expected, packet[5:12], "#%d pcr=%d", i+1, x.pcr)
	}
}

func TestTS_PCR(t *testing.T) {
	assert := assert.New(t)

	packet := NewTS(256)
	packet[3] |= 0x20 // has AF
	packet[4] = 7

	copy(packet[5:], []byte{0x10, 0xE1, 0x57, 0x8A, 0x18, 0xFE, 0x7B})

	assert.Equal(PCR(2268366350823), packet.PCR())
}

func TestPCR_EstimatedPCR(t *testing.T) {
	previousPCR := PCR(354923263808)
	currentPCR := PCR(354924281094)
	lastBlock := uint64(7708)
	currentBlock := uint64(7520)

	expected := PCR(354925273568)
	estimatedPCR := currentPCR.EstimatedPCR(previousPCR, lastBlock, currentBlock)

	assert.Equal(t, expected, estimatedPCR)
}

func TestPCR_Add(t *testing.T) {
	assert.Equal(t, PCR(2), MaxPcr.Add(3))
}
