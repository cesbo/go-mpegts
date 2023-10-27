package mpegts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleTS() {
	packet := TS([]byte{0x47, 0x40, 0x11, 0x15})
	fmt.Println("PID", packet.PID(), "CC", packet.CC())

	// Output:
	// PID 17 CC 5
}

func TestTS(t *testing.T) {
	assert := assert.New(t)
	packet := NewTS(256)

	assert.Equal(PacketSize, len(packet))
	assert.Equal([]byte{0x47, 0x01, 0x00, 0x00}, []byte(packet[:4]))

	packetNotScrambled := TS([]byte{0x47, 0x40, 0x11, 0x15})
	assert.Equal(NotScrambled, packetNotScrambled.TSC())

	packetScrambledEven := TS([]byte{0x47, 0x40, 0x11, 0x95})
	assert.Equal(ScrambledEvenKey, packetScrambledEven.TSC())

	packetScrambledOdd := TS([]byte{0x47, 0x40, 0x11, 0xD5})
	assert.Equal(ScrambledOddKey, packetScrambledOdd.TSC())
}

func TestTS_PID(t *testing.T) {
	assert := assert.New(t)
	packet := TS([]byte{0x47, 0x40, 0x11, 0x15})

	assert.Equal(PID(17), packet.PID())

	packet.SetPID(256)
	assert.Equal([]byte{0x47, 0x41, 0x00, 0x15}, []byte(packet))
	assert.Equal(PID(256), packet.PID())

	packet.SetPID(0xFFFF)
	assert.Equal([]byte{0x47, 0x5F, 0xFF, 0x15}, []byte(packet))
	assert.Equal(MaxPid, packet.PID())

	packet.SetPID(0)
	assert.Equal([]byte{0x47, 0x40, 0x00, 0x15}, []byte(packet))
	assert.Equal(PID(0), packet.PID())
}

func TestTS_CC(t *testing.T) {
	assert := assert.New(t)
	packet := TS([]byte{0x47, 0x40, 0x11, 0x15})

	assert.Equal(uint8(5), packet.CC())

	packet.SetCC(10)
	assert.Equal([]byte{0x47, 0x40, 0x11, 0x1A}, []byte(packet))
	assert.Equal(uint8(10), packet.CC())

	packet.SetCC(0xFF)
	assert.Equal([]byte{0x47, 0x40, 0x11, 0x1F}, []byte(packet))
	assert.Equal(uint8(15), packet.CC())

	packet.IncrementCC()
	assert.Equal([]byte{0x47, 0x40, 0x11, 0x10}, []byte(packet))
	assert.Equal(uint8(0), packet.CC())
}

func TestTS_Payload(t *testing.T) {
	assert := assert.New(t)
	packet := NewTS(17)
	packet[1] = 0x40
	packet[3] = 0x15

	for i := 4; i < 188; i++ {
		packet[i] = byte(i - 4)
	}

	assert.True(packet.HasPayload())
	assert.True(packet.HasPUSI())

	payload := packet.Payload()

	if assert.NotNil(payload) == true {
		assert.Equal([]byte(packet[4:]), payload)
	}

	// Check with Adaptation Field

	packet[3] |= 0x20
	packet[4] = 7    // AF length
	packet[5] = 0x10 // PCR flag

	payload = packet.Payload()

	if assert.NotNil(payload) == true {
		assert.Equal([]byte(packet[4+1+7:]), payload)
	}

	// Check with invalid adaptation field length

	packet[4] = 188

	payload = packet.Payload()

	assert.Nil(payload)
}

func TestTS_Fill(t *testing.T) {
	makePacketWithoutAf := func(size int) (packet, expected TS) {
		packet = NewTS(101)
		packet[3] |= 0x10 // with payload

		skip := 4

		for i := 0; i < size; i++ {
			packet[skip+i] = byte(i + 0x30)
		}

		packet.Fill(skip + size)

		expected = NewTS(101)
		expected[3] = 0x30
		expected[4] = byte(PacketSize - 4 - 1 - size)
		expected[5] = 0x00

		next := 4 + 1 + int(expected[4])

		for i := 6; i < next; i++ {
			expected[i] = 0xFF
		}

		for i := 0; i < size; i++ {
			expected[next+i] = byte(i + 0x30)
		}

		return
	}

	makePacketWithAf := func(size int) (packet, expected TS) {
		packet = NewTS(101)
		packet[3] |= 0x30 // with payload and adaptation field

		af := []byte{
			0x07, // AF length
			0x10, // AF flags (PCR)
			0x00, 0x02, 0x32, 0x89, 0x7E, 0xF7,
		}

		skip := 4 + copy(packet[4:], af)

		for i := 0; i < size; i++ {
			packet[skip+i] = byte(i + 0x30)
		}

		packet.Fill(skip + size)

		expected = NewTS(101)
		expected[3] = 0x30
		copy(expected[4:], af)
		expected[4] = byte(PacketSize - 4 - 1 - size)

		next := 4 + 1 + int(expected[4])

		for i := skip; i < next; i++ {
			expected[i] = 0xFF
		}

		for i := 0; i < size; i++ {
			expected[next+i] = byte(i + 0x30)
		}

		return
	}

	t.Run("without af", func(t *testing.T) {
		assert := assert.New(t)

		sizes := []int{
			20,  // small packet
			181, // adaptation filed with single stuffing byte
			182, // adaptation field without stuffing (only size and header)
			183, // adaptation field without header (only size)
		}

		for _, size := range sizes {
			packet, expected := makePacketWithoutAf(size)
			assert.Equal(expected, packet, size)
		}
	})

	t.Run("with af", func(t *testing.T) {
		assert := assert.New(t)

		sizes := []int{
			20,  // small packet
			175, // adaptation filed with single stuffing byte
			176, // adaptation field without stuffing (only size and header)
		}

		for _, size := range sizes {
			packet, expected := makePacketWithAf(size)
			assert.Equal(expected, packet, size)
		}
	})
}
