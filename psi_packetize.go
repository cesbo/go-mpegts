package mpegts

import (
	"encoding/binary"

	"github.com/cesbo/go-mpegts/crc32"
)

type PsiPacketizer interface {
	Packetize(TS, func([]byte)) error
}

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

func psiPacketize(b psiSection, packet TS, fn func([]byte)) error {
	item := -1

	var (
		sectionSize int
		sectionFill int
		packetFill  int
	)

	var (
		crc       uint32
		crcBuffer [crc32.Size]byte
	)

	send := func(data []byte) {
		skip := 0
		for skip < len(data) {
			n := copy(packet[packetFill:PacketSize], data[skip:])
			skip += n
			packetFill += n
			if packetFill == PacketSize {
				fn(packet)
				packetFill = 4
				packet.ClearPUSI()
				packet.IncrementCC()
			}
		}
		sectionFill += skip
	}

	for {
		// start new section.
		sectionSize = b.sectionSize(item)
		sectionFill = 0

		// no more data to packetize
		if sectionSize == 0 {
			return nil
		}

		// prepare first TS packet
		packet.SetPUSI()
		packet.SetPayload()
		packet[4] = 0 // Set PUSI Pointer
		packetFill = 5

		header := b.sectionHeader(item)

		// set length
		binary.BigEndian.PutUint16(
			header[1:],
			(((uint16(header[1] & 0xF0)) << 8) | ((uint16(sectionSize - PsiHeaderSize)) & 0x0FFF)),
		)

		// write PSI header
		crc = crc32.Checksum(0xFFFFFFFF, header)
		send(header)

		sectionSize -= crc32.Size

		// write PSI data
		for sectionFill < sectionSize {
			// put item
			data := b.sectionItem(item)
			if data == nil {
				break
			}

			item += 1

			if len(data) == 0 {
				continue
			}

			crc = crc32.Checksum(crc, data)
			send(data)
		}

		// current section finished. set checksum
		binary.BigEndian.PutUint32(crcBuffer[:], crc)
		send(crcBuffer[:])

		if packetFill < PacketSize {
			send(NullTS[packetFill:PacketSize])
		}

		// empty section without items
		if item == -1 {
			item = 0
		}
	}
}
