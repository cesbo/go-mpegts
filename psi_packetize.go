package mpegts

import (
	"github.com/cesbo/go-mpegts/crc32"
)

type psiSection interface {
	// Calculate section size from element i.
	// i == -1: initial section: header, descriptors, items, checksum
	// i >= 0: next sections: header, items, checksum
	// If no more data to packetize should return 0
	sectionSize(i int) int

	// Get section header without descriptors
	// i == -1: initial section. You may set section_number to 0
	// i >= 0: next section. You may increment section_number by 1
	sectionHeader(i int) []byte

	// Get section item.
	// i == -1: requests section descriptor. If no descriptor return empty slice
	// i >= 0: requests section items
	// If no more items return nil
	sectionItem(i int) []byte
}

// PsiPacketizer is a helper to splits PSI section into multiple TS packets.
// If data more than fits into one section, it will be split into multiple sections.
type PsiPacketizer struct {
	inner       psiSection
	sectionItem int
	sectionSize int
	sectionFill int
	skip        int
	crc         uint32
}

func newPsiPacketizer(inner psiSection) *PsiPacketizer {
	return &PsiPacketizer{
		inner:       inner,
		sectionItem: -1,
	}
}

func (p *PsiPacketizer) Next(ts TS) bool {
	packetFill := 4

	// start new section
	if p.sectionFill == p.sectionSize {
		p.sectionSize = p.inner.sectionSize(p.sectionItem)
		p.sectionFill = 0

		// no more data to packetize
		if p.sectionSize == 0 {
			return false
		}

		// prepare first TS packet
		ts.SetPUSI()
		ts.SetPayload()
		ts[4] = 0 // Set PUSI Pointer
		packetFill = 5

		header := p.inner.sectionHeader(p.sectionItem)

		// set length
		s := uint16(p.sectionSize - PsiHeaderSize)
		header[1] = (header[1] & 0xF0) | (uint8((s >> 8) & 0x0F))
		header[2] = uint8(s & 0xFF)

		p.crc = crc32.Checksum(0xFFFFFFFF, header)

		// write section header (without descriptors)
		n := copy(ts[packetFill:], header)
		p.sectionFill += n
		packetFill += n
	} else {
		ts.ClearPUSI()
	}

	for {
		// current section finished. set checksum
		if (p.sectionFill + crc32.Size) == p.sectionSize {
			for p.skip < crc32.Size {
				shift := uint32(24 - (8 * p.skip))
				ts[packetFill] = byte(p.crc >> shift)
				packetFill += 1
				p.skip += 1

				if packetFill == PacketSize {
					return true
				}
			}

			p.sectionFill = p.sectionSize
			p.skip = 0

			// empty section
			if p.sectionItem == -1 {
				p.sectionItem = 0
			}

			break
		}

		// put item
		data := p.inner.sectionItem(p.sectionItem)
		if data == nil {
			break
		}

		if len(data) == 0 {
			p.sectionItem += 1
			continue
		}

		n := copy(ts[packetFill:], data[p.skip:])
		packetFill += n
		p.skip += n

		if p.skip == len(data) {
			p.sectionItem += 1
			p.sectionFill += p.skip
			p.skip = 0
		}

		if packetFill == PacketSize {
			return true
		}
	}

	if packetFill < PacketSize {
		copy(ts[packetFill:], NullTS[packetFill:])
	}

	return true
}
