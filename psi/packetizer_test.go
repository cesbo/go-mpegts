package psi

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cesbo/go-mpegts/crc32"
	"github.com/cesbo/go-mpegts/ts"
)

type psiMock int

const (
	psiMockHeaderSize = 10
)

var psiMockHeader [psiMockHeaderSize]byte

func (p psiMock) SectionSize(i int) int {
	switch p {
	case 0:
		// Test_Packetizer_SingleTS
		switch i {
		case -1:
			return psiMockHeaderSize + 5 + 10 + crc32.Size
		case 1:
			return 0
		}
	case 1:
		// Test_Packetizer_TwoTS
		switch i {
		case -1:
			return psiMockHeaderSize + (183 - psiMockHeaderSize - 2) + crc32.Size
		case 1:
			return 0
		}
	case 2:
		// Test_Packetizer_TwoSections
		switch i {
		case -1:
			return psiMockHeaderSize + 200 + crc32.Size
		case 1:
			return psiMockHeaderSize + 200 + crc32.Size
		case 2:
			return 0
		}
	}

	panic("unreachable")
}

func (p psiMock) SectionHeader(i int) []byte {
	if i == -1 {
		for i := 0; i < len(psiMockHeader); i++ {
			psiMockHeader[i] = byte(i + 10)
		}

		psiMockHeader[1] = 0xF0 // section_length 1
		psiMockHeader[2] = 0x00 // section_length 2
		psiMockHeader[6] = 0    // section_number
		psiMockHeader[7] = 0    // last_section_number

		switch p {
		case 0:
			// Test_Packetizer_SingleTS
			psiMockHeader[8] = 0xF0
			psiMockHeader[9] = 0x05
		case 1:
			// Test_Packetizer_TwoTS
			psiMockHeader[8] = 0xF0
			psiMockHeader[9] = 0x00
		case 2:
			// Test_Packetizer_TwoSections
			psiMockHeader[7] = 1
			psiMockHeader[8] = 0xF0
			psiMockHeader[9] = 0x00
		}
	} else {
		psiMockHeader[6] += 1
		psiMockHeader[8] = 0xF0
		psiMockHeader[9] = 0x00
	}

	return psiMockHeader[:]
}

func (p psiMock) SectionItem(i int) []byte {
	switch p {
	case 0:
		// Test_Packetizer_SingleTS
		switch i {
		case -1:
			return []byte{0x40, 0x03, 0xAA, 0xBB, 0xCC}
		case 0:
			h := make([]byte, 10)
			for z := 0; z < len(h); z++ {
				h[z] = 0xA0 + byte(z)
			}
			return h
		default:
			panic("unreachable")
		}
	case 1:
		// Test_Packetizer_TwoTS
		switch i {
		case -1:
			return []byte{}
		case 0:
			h := make([]byte, 183-psiMockHeaderSize-2)
			for z := 0; z < len(h); z++ {
				h[z] = 0x50 + byte(z)
			}
			return h
		default:
			panic("unreachable")
		}
	case 2:
		// Test_Packetizer_TwoSections
		switch i {
		case -1:
			return []byte{}
		case 0:
			h := make([]byte, 200)
			for z := 0; z < len(h); z++ {
				h[z] = 0x10 + byte(z)
			}
			return h
		case 1:
			h := make([]byte, 200)
			for z := 0; z < len(h); z++ {
				h[z] = 0x20 + byte(z)
			}
			return h
		}
	}
	panic("unknown test id")
}

func TestPacketizer_SingleTS(t *testing.T) {
	assert := assert.New(t)

	packetizer := NewPacketizer(psiMock(0))
	packet := ts.NewPacket(4455)
	packet.SetCC(1)

	if !assert.True(packetizer.NextPacket(packet)) {
		return
	}

	expected := ts.TS{
		0x47, 0x40 | 0x11, 0x67, 0x10 | 1,
		0x00,
		10, 0xF0, 0x00, 13, 14, 15, 0, 0, 0xF0, 0x05,
		0x40, 0x03, 0xAA, 0xBB, 0xCC,
		0xA0, 0xA1, 0xA2, 0xA3, 0xA4,
		0xA5, 0xA6, 0xA7, 0xA8, 0xA9,
		0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	totalLength := psiMockHeaderSize + 5 + 10 + crc32.Size
	expected[7] = byte(totalLength - PsiHeaderSize)
	crc := crc32.Checksum(0xFFFFFFFF, expected[5:5+totalLength-crc32.Size])
	binary.BigEndian.PutUint32(expected[5+totalLength-crc32.Size:], crc)

	assert.Equal(expected, packet[:len(expected)])

	assert.False(packetizer.NextPacket(packet))
}

func TestPacketizer_TwoTS(t *testing.T) {
	assert := assert.New(t)

	packetizer := NewPacketizer(psiMock(1))
	packet := ts.NewPacket(4455)

	expected1 := ts.TS{
		0x47, 0x40 | 0x11, 0x67, 0x10 | 1,
		0x00,
		10, 0xF0, 0x00, 13, 14, 15, 0, 0, 0xF0, 0x00,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F,
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F,
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F,
		0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF,
		0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7, 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF,
		0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF,
		0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9, 0xDA, 0xDB, 0xDC, 0xDD, 0xDE, 0xDF,
		0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF,
		0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA,
		0x00, 0x00,
	}

	totalLength := psiMockHeaderSize + (183 - psiMockHeaderSize - 2) + crc32.Size
	expected1[7] = byte(totalLength - PsiHeaderSize)
	crc := crc32.Checksum(0xFFFFFFFF, expected1[5:5+totalLength-crc32.Size])
	expected1[186] = byte(crc >> 24)
	expected1[187] = byte(crc >> 16)

	packet.SetCC(1)
	if !assert.True(packetizer.NextPacket(packet)) {
		return
	}
	assert.Equal(expected1, packet)

	expected2 := ts.TS{
		0x47, 0x11, 0x67, 0x10 | 2,
		0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	expected2[4] = byte(crc >> 8)
	expected2[5] = byte(crc)

	packet.IncrementCC()
	if !assert.True(packetizer.NextPacket(packet)) {
		return
	}
	assert.Equal(expected2, packet[:len(expected2)])

	assert.False(packetizer.NextPacket(packet))
}

func TestPacketizer_TwoSections(t *testing.T) {
	assert := assert.New(t)

	var crc uint32

	packetizer := NewPacketizer(psiMock(2))
	packet := ts.NewPacket(4455)

	expected11 := ts.TS{
		0x47, 0x40 | 0x11, 0x67, 0x10 | 1,
		0x00,
		10, 0xF0, 0x00, 13, 14, 15, 0, 1, 0xF0, 0x00,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F,
		0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F,
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F,
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F,
		0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF,
		0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7, 0xB8, 0xB9, 0xBA, 0xBB, 0xBC,
	}

	expected12 := ts.TS{
		0x47, 0x11, 0x67, 0x10 | 2,
		0xBD, 0xBE, 0xBF, 0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC,
		0xCD, 0xCE, 0xCF, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7,
		0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	expected11[7] = byte(psiMockHeaderSize + 200 + crc32.Size - PsiHeaderSize)
	crc = crc32.Checksum(0xFFFFFFFF, expected11[5:])
	crc = crc32.Checksum(crc, expected12[4:4+27])
	binary.BigEndian.PutUint32(expected12[4+27:], crc)

	packet.SetCC(1)
	if !assert.True(packetizer.NextPacket(packet)) {
		return
	}
	assert.Equal(expected11, packet)

	packet.IncrementCC()
	if !assert.True(packetizer.NextPacket(packet)) {
		return
	}
	assert.Equal(expected12, packet[:len(expected12)])

	// second section

	expected21 := ts.TS{
		0x47, 0x40 | 0x11, 0x67, 0x10 | 3,
		0x00,
		10, 0xF0, 0x00, 13, 14, 15, 1, 1, 0xF0, 0x00,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F,
		0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F,
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F,
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F,
		0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF,
		0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7, 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF,
		0xC0, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC,
	}

	expected22 := ts.TS{
		0x47, 0x11, 0x67, 0x10 | 4,
		0xCD, 0xCE, 0xCF, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9, 0xDA, 0xDB, 0xDC,
		0xDD, 0xDE, 0xDF, 0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7,
		0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	expected21[7] = byte(psiMockHeaderSize + 200 + crc32.Size - PsiHeaderSize)
	crc = crc32.Checksum(0xFFFFFFFF, expected21[5:])
	crc = crc32.Checksum(crc, expected22[4:4+27])
	binary.BigEndian.PutUint32(expected22[4+27:], crc)

	packet.IncrementCC()
	if !assert.True(packetizer.NextPacket(packet)) {
		return
	}
	assert.Equal(expected21, packet)

	packet.IncrementCC()
	if !assert.True(packetizer.NextPacket(packet)) {
		return
	}
	assert.Equal(expected22, packet[:len(expected22)])

	// end

	assert.False(packetizer.NextPacket(packet))
}
