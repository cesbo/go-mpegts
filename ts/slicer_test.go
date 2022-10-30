package ts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleSlicer() {
	slicer := Slicer{}
	buffer := []byte{}

	for packet := slicer.Begin(buffer); packet != nil; packet = slicer.Next() {
		fmt.Println(len(packet) == PacketSize)
		fmt.Println(packet[0] == SyncByte)
	}
}

func TestSlicer(t *testing.T) {
	var packet TS
	var err error

	// Generate test data

	testPackets := make([]byte, PacketSize*2)
	testHeader := []byte{0x47, 0x1F, 0xFF, 0x10}
	copy(testPackets[0:], testHeader)
	copy(testPackets[PacketSize:], testHeader)
	for i := 4; i < 188; i++ {
		testPackets[i] = byte(i - 4)
		testPackets[i+188] = byte(187 - i)
	}

	slicer := Slicer{}

	t.Run("single packet", func(t *testing.T) {
		assert := assert.New(t)

		totalCount := 0
		buffer := testPackets[:PacketSize]
		for packet = slicer.Begin(buffer); packet != nil; packet = slicer.Next() {
			s1 := totalCount * PacketSize
			s2 := s1 + PacketSize

			assert.Equal(
				TS(testPackets[s1:s2]),
				packet,
			)

			totalCount += 1
		}
		assert.Equal(1, totalCount)

		err = slicer.Err()
		assert.NoError(err)
	})

	t.Run("two packets", func(t *testing.T) {
		assert := assert.New(t)

		totalCount := 0
		for packet = slicer.Begin(testPackets); packet != nil; packet = slicer.Next() {
			s1 := totalCount * PacketSize
			s2 := s1 + PacketSize

			assert.Equal(
				TS(testPackets[s1:s2]),
				packet,
			)

			totalCount += 1
		}
		assert.Equal(2, totalCount)

		assert.NoError(slicer.Err())
	})

	t.Run("append partially. 2 parts", func(t *testing.T) {
		assert := assert.New(t)

		var packet TS

		packet = slicer.Begin(testPackets[:50])
		assert.Nil(packet)

		assert.NoError(slicer.Err())

		packet = slicer.Begin(testPackets[50:PacketSize])

		if assert.NotNil(packet) {
			assert.Equal(
				TS(testPackets[:PacketSize]),
				packet,
			)
		}

		assert.Nil(slicer.Next())

		assert.NoError(slicer.Err())
	})

	t.Run("append partially. 3 parts", func(t *testing.T) {
		assert := assert.New(t)

		var packet TS

		packet = slicer.Begin(testPackets[:50])
		assert.Nil(packet)

		assert.NoError(slicer.Err())

		packet = slicer.Begin(testPackets[50:100])
		assert.Nil(packet)

		assert.NoError(slicer.Err())

		packet = slicer.Begin(testPackets[100:PacketSize])
		if assert.NotNil(packet) {
			assert.Equal(
				TS(testPackets[:PacketSize]),
				packet,
			)
		}

		assert.Nil(slicer.Next())

		assert.NoError(slicer.Err())
	})

	t.Run("append partially. Part after full packet", func(t *testing.T) {
		assert := assert.New(t)

		var packet TS

		packet = slicer.Begin(testPackets[:PacketSize+4])

		if assert.NotNil(packet) {
			assert.Equal(
				TS(testPackets[:PacketSize]),
				packet,
			)
		}

		assert.Nil(slicer.Next())

		assert.NoError(slicer.Err())

		// skip 4 bytes (second packet header in testPacket)
		packet = slicer.Begin(testPackets[PacketSize+4:])

		if assert.NotNil(packet) {
			assert.Equal(
				TS(testPackets[PacketSize:]),
				packet,
			)
		}

		assert.Nil(slicer.Next())

		assert.NoError(slicer.Err())
	})

	t.Run("skip unexpected bytes", func(t *testing.T) {
		assert := assert.New(t)

		totalCount := 0
		buffer := testPackets[PacketSize-4:]
		for packet = slicer.Begin(buffer); packet != nil; packet = slicer.Next() {
			s1 := (totalCount * PacketSize) + PacketSize
			s2 := s1 + PacketSize

			assert.Equal(
				TS(testPackets[s1:s2]),
				packet,
			)

			totalCount += 1
		}
		assert.Equal(1, totalCount)

		err = slicer.Err()
		assert.NoError(err)
	})
}
