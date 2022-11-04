package mpegts

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/cesbo/go-mpegts/crc32"
)

const (
	PmtHeaderSize  = 12
	PmtMaximumSize = 1024
	PmtItemSize    = 5
)

// PMT is Program Map Table
type PMT struct {
	header []byte
	Items  []*PmtItem
}

// PMT Item contains information about elementary stream
type PmtItem struct {
	header []byte
}

var (
	ErrPmtFormat = errors.New("pmt: invalid format")
)

var (
	emptyPmt = []byte{
		0x02,        // table_id
		0x80 | 0x30, // section_length_1
		0x00,        // section_length 2
		0x00,        // program_number 1
		0x00,        // program_number 2
		0xC0 | 0x01, // version
		0x00,        // section_number
		0x00,        // last_section_number
		0xE0,        // pcr_pid 1
		0x00,        // pcr_pid 2
		0xF0,        // program_info_length 1
		0x00,        // program_info_length 2
	}
	emptyPmtItem = []byte{
		0x00, // stream_type
		0xE0, // elementary_pid 1
		0x00, // elementary_pid 2
		0xF0, // es_info_length 1
		0x00, // es_info_length 2
	}
)

func NewPmt() *PMT {
	p := new(PMT)
	p.header = make([]byte, len(emptyPmt))
	copy(p.header, emptyPmt)

	return p
}

func (p *PMT) ParsePmtSection(b []byte) error {
	if len(b) < (PmtHeaderSize + crc32.Size) {
		return ErrPmtFormat
	}

	next := PmtHeaderSize
	end := len(b) - crc32.Size

	pmtDescLen := binary.BigEndian.Uint16(b[10:]) & 0x0FFF
	if pmtDescLen > 0 {
		if pmtDescLen > 0x03FF {
			return ErrPmtFormat
		}

		next += int(pmtDescLen)
		if next > end {
			return ErrPmtFormat
		}

		pmtDesc := Descriptors(b[PmtHeaderSize:next])
		if err := pmtDesc.Check(); err != nil {
			return fmt.Errorf("pmt: %w", err)
		}
	}

	// copy header only from first section
	if b[6] == 0 {
		p.header = make([]byte, next)
		copy(p.header, b)
	}

	skip := next

	for skip < end {
		next += PmtItemSize
		if next > end {
			return ErrPmtFormat
		}

		esDescLen := binary.BigEndian.Uint16(b[skip+3:]) & 0x0FFF
		if esDescLen > 0 {
			if esDescLen > 0x03FF {
				return ErrPmtFormat
			}

			next += int(esDescLen)
			if next > end {
				return ErrPmtFormat
			}

			esDesc := Descriptors(b[skip+PmtItemSize : next])
			if err := esDesc.Check(); err != nil {
				return fmt.Errorf("pmt: %w", err)
			}
		}

		item := new(PmtItem)
		item.header = make([]byte, next-skip)
		copy(item.header, b[skip:])

		p.Items = append(p.Items, item)

		skip = next
	}

	return nil
}

func (p *PMT) Version() uint8 {
	return (p.header[5] & 0x3E) >> 1
}

func (p *PMT) SetVersion(version uint8) {
	p.header[5] &^= 0x3E
	p.header[5] |= (version << 1) & 0x3E
}

func (p *PMT) PNR() uint16 {
	return binary.BigEndian.Uint16(p.header[3:])
}

func (p *PMT) SetPNR(pnr uint16) {
	binary.BigEndian.PutUint16(p.header[3:], pnr)
}

func (p *PMT) PCR() uint16 {
	return binary.BigEndian.Uint16(p.header[8:]) & 0x1FFF
}

func (p *PMT) SetPCR(pcr uint16) {
	binary.BigEndian.PutUint16(p.header[8:], 0xE000|pcr)
}

func (p *PMT) Descriptors() Descriptors {
	return Descriptors(p.header[PmtHeaderSize:])
}

func (p *PMT) AppendDescriptors(desc Descriptors) {
	if len(desc) == 0 {
		return
	}

	ds := len(p.header) - PmtHeaderSize + len(desc)
	if ds > 0x03FF {
		panic("failed to append pmt descriptors: size limit")
	}
	p.header = append(p.header, desc...)

	binary.BigEndian.PutUint16(p.header[10:], 0xF000|uint16(ds))
}

// Calculates LastSectionNumber
func (p *PMT) Finalize() {
	p.header[6] = 0
	p.header[7] = 0

	size := len(p.header)
	remain := PmtMaximumSize - size - crc32.Size

	for _, item := range p.Items {
		is := len(item.header)
		if is > remain {
			remain = PmtMaximumSize - PmtHeaderSize - crc32.Size
			p.header[7] += 1
		}
		remain -= is
	}
}

// Packetizer returns a new PsiPacketizer to get TS packets from PMT
func (p *PMT) Packetizer() *PsiPacketizer {
	return newPsiPacketizer(p)
}

func (p *PMT) sectionSize(i int) int {
	if i == len(p.Items) {
		return 0
	}

	size := crc32.Size

	if i == -1 {
		size += len(p.header)
		i = 0
	} else {
		size += PmtHeaderSize
	}

	for i < len(p.Items) {
		is := len(p.Items[i].header)
		if (size + is) > PmtMaximumSize {
			break
		} else {
			size += is
			i += 1
		}
	}

	return size
}

func (p *PMT) sectionHeader(i int) []byte {
	if i == -1 {
		p.header[6] = 0
		s := uint16(len(p.header) - PmtHeaderSize)
		p.header[10] = 0xF0 | byte(s>>8)
		p.header[11] = byte(s)
	} else {
		p.header[6] += 1
		p.header[10] = 0xF0
		p.header[11] = 0x00
	}

	return p.header[:PmtHeaderSize]
}

func (p *PMT) sectionItem(i int) []byte {
	if i == -1 {
		return p.header[PmtHeaderSize:]
	}

	if i < len(p.Items) {
		return p.Items[i].header
	}

	return nil
}

func NewPmtItem() *PmtItem {
	p := new(PmtItem)
	p.header = make([]byte, len(emptyPmtItem))
	copy(p.header, emptyPmtItem)

	return p
}

func (p *PmtItem) Type() uint8 {
	return p.header[0]
}

func (p *PmtItem) SetType(ty uint8) {
	p.header[0] = ty
}

func (p *PmtItem) PID() uint16 {
	return binary.BigEndian.Uint16(p.header[1:]) & 0x1FFF
}

func (p *PmtItem) SetPID(pid uint16) {
	binary.BigEndian.PutUint16(p.header[1:], 0xE000|pid)
}

func (p *PmtItem) Descriptors() Descriptors {
	return Descriptors(p.header[PmtItemSize:])
}

func (p *PmtItem) AppendDescriptors(desc Descriptors) {
	if len(desc) == 0 {
		return
	}

	ds := len(p.header) - PmtItemSize + len(desc)
	if ds > 0x03FF {
		panic("failed to append pmt item descriptors: size limit")
	}
	p.header = append(p.header, desc...)

	binary.BigEndian.PutUint16(p.header[3:], 0xF000|uint16(ds))
}

func (p *PmtItem) checkData05() PacketType {
	d := p.Descriptors()

	for len(d) != 0 {
		if d[0] == 0x6F {
			return PACKET_AIT
		}

		d = d.Next()
	}

	return PACKET_DATA
}

func (p *PmtItem) checkData06() PacketType {
	d := p.Descriptors()

	for len(d) != 0 {
		switch d[0] {
		case 0x56:
			return PACKET_TTX
		case 0x59:
			return PACKET_SUB
		case 0x6A, 0x81:
			return PACKET_AUDIO_AC_3
		case 0x7A:
			return PACKET_AUDIO_EAC_3
		}

		d = d.Next()
	}

	return PACKET_DATA
}

// Get PacketType by element ID and related descriptors
func (p *PmtItem) GetType() PacketType {
	var ty PacketType

	switch p.Type() {
	// Video
	case 0x01:
		ty = PACKET_VIDEO_H261
	case 0x02:
		ty = PACKET_VIDEO_H262
	case 0x10:
		ty = PACKET_VIDEO_H263
	case 0x1B:
		ty = PACKET_VIDEO_H264
	case 0x24:
		ty = PACKET_VIDEO_H265
	// Audio
	case 0x03:
		ty = PACKET_AUDIO_MP2
	case 0x04:
		ty = PACKET_AUDIO_MP3
	case 0x0F:
		ty = PACKET_AUDIO_AAC
	case 0x11:
		ty = PACKET_AUDIO_LATM
	case 0x81:
		ty = PACKET_AUDIO_AC_3
	case 0x87:
		ty = PACKET_AUDIO_EAC_3
	// Data
	case 0x05:
		ty = p.checkData05()
	case 0x06:
		ty = p.checkData06()
	default:
		ty = PACKET_DATA
	}

	return ty
}

func (p *PmtItem) Clone() *PmtItem {
	c := new(PmtItem)
	c.header = make([]byte, len(p.header))
	copy(c.header, p.header)

	return c
}
