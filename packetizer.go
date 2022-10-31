package mpegts

import (
	"encoding/binary"

	"github.com/cesbo/go-mpegts/crc32"
)

type SectionBuilder interface {
	// Calculate section size from element i.
	// i == -1: initial section: header, descriptors, items, checksum
	// i >= 0: next sections: header, items, checksum
	// If no more data to packetize should return 0
	SectionSize(i int) int

	// Get section header without descriptors
	// i == -1: initial section. You may set section_number to 0
	// i >= 0: next section. You may increment section_number by 1
	SectionHeader(i int) []byte

	// Get section item.
	// i == -1: requests section descriptor. If no descriptor return empty slice
	// i >= 0: requests section items
	// If no more items return nil
	SectionItem(i int) []byte
}

// Packetizer splits information tables into TS packets.
type Packetizer struct {
	inner SectionBuilder

	item int
	skip int // skip bytes in item

	fill int // section fill
	size int // section size

	crc uint32
}

func NewPacketizer(b SectionBuilder) *Packetizer {
	return &Packetizer{
		inner: b,
		item:  -1,
	}
}

func (p *Packetizer) NextPacket(packet TS) bool {
	packetSkip := 4

	if p.fill == p.size {
		p.size = p.inner.SectionSize(p.item)
		p.fill = 0

		// no more data to packetize
		if p.size == 0 {
			return false
		}

		// start first packet
		packet.SetPUSI()
		packet.SetPayload()
		packet[4] = 0 // PUSI Pointer
		packetSkip += 1

		header := p.inner.SectionHeader(p.item)

		// set length
		b := uint16(header[1]&0xF0) << 8
		s := uint16(p.size-PsiHeaderSize) & 0x0FFF
		binary.BigEndian.PutUint16(header[1:], s|b)

		n := copy(packet[packetSkip:], header)
		p.crc = crc32.Checksum(0xFFFFFFFF, header)
		p.fill += n
		packetSkip += n
	} else {
		packet.ClearPUSI()
	}

	for {
		// current section finished. set checksum
		if (p.fill + crc32.Size) == p.size {
			for p.skip < crc32.Size {
				shift := uint32(24 - (8 * p.skip))
				packet[packetSkip] = byte(p.crc >> shift)
				packetSkip += 1
				p.skip += 1
				if packetSkip == PacketSize {
					return true
				}
			}
			p.fill = p.size
			p.skip = 0

			// empty section
			if p.item == -1 {
				p.item = 0
			}

			break
		}

		// put item
		data := p.inner.SectionItem(p.item)
		if data == nil {
			break
		}

		if len(data) == 0 {
			p.item += 1
			continue
		}

		if p.skip == 0 {
			p.crc = crc32.Checksum(p.crc, data)
		}

		n := copy(packet[packetSkip:], data[p.skip:])
		p.skip += n
		p.fill += n
		packetSkip += n

		if p.skip == len(data) {
			p.skip = 0
			p.item += 1
		}

		if packetSkip == PacketSize {
			return true
		}
	}

	if packetSkip < PacketSize {
		copy(packet[packetSkip:], NullTS[packetSkip:])
	}

	return true
}
