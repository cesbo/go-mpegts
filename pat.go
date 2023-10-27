package mpegts

import (
	"encoding/binary"
	"errors"

	"github.com/cesbo/go-mpegts/crc32"
)

const (
	PatHeaderSize  = 8
	PatMaximumSize = 1024
	PatItemSize    = 4
)

// PAT is Program Association Table
type PAT struct {
	header []byte
	Items  []*PatItem
}

type PatItem struct {
	header []byte
}

var (
	ErrPatFormat = errors.New("pat: invalid format")
)

var (
	emptyPat = []byte{
		0x00,        // table_id
		0x80 | 0x30, // section_length 1
		0x00,        // section_length 2
		0x00,        // transport_stream_id 1
		0x00,        // transport_stream_id 2
		0xC0 | 0x01, // version
		0x00,        // section_number
		0x00,        // last_section_number
	}
	emptyPatItem = []byte{
		0x00, // program_number 1
		0x00, // program_number 2
		0xE0, // program_map_pid 1
		0x00, // program_map_pid 2
	}
)

func NewPat() *PAT {
	p := new(PAT)
	p.header = make([]byte, len(emptyPat))
	copy(p.header, emptyPat)

	return p
}

func (p *PAT) ParsePatSection(b []byte) error {
	if len(b) < (PatHeaderSize + crc32.Size) {
		return ErrPatFormat
	}

	next := PatHeaderSize
	end := len(b) - crc32.Size

	// copy header only from first section
	if b[6] == 0 {
		p.header = make([]byte, next)
		copy(p.header, b)
	}

	skip := next

	for skip < end {
		next += PatItemSize
		if next > end {
			return ErrPatFormat
		}

		item := new(PatItem)
		item.header = make([]byte, next-skip)
		copy(item.header, b[skip:])

		p.Items = append(p.Items, item)

		skip = next
	}

	return nil
}

func (p *PAT) Version() uint8 {
	return (p.header[5] & 0x3E) >> 1
}

func (p *PAT) SetVersion(version uint8) {
	p.header[5] &^= 0x3E
	p.header[5] |= (version << 1) & 0x3E
}

// Returns Transport Stream ID
func (p *PAT) TSID() uint16 {
	return binary.BigEndian.Uint16(p.header[3:])
}

func (p *PAT) SetTSID(tsid uint16) {
	binary.BigEndian.PutUint16(p.header[3:], tsid)
}

// Calculates LastSectionNumber
func (p *PAT) Finalize() {
	p.header[6] = 0
	p.header[7] = 0

	size := len(p.header)
	remain := PatMaximumSize - size - crc32.Size

	for _, item := range p.Items {
		is := len(item.header)
		if is > remain {
			remain = PatMaximumSize - PatHeaderSize - crc32.Size
			p.header[7] += 1
		}
		remain -= is
	}
}

// Packetizer returns a new PsiPacketizer to get TS packets from PAT
func (p *PAT) Packetizer() *PsiPacketizer {
	return newPsiPacketizer(p)
}

func (p *PAT) sectionSize(i int) int {
	if i == len(p.Items) {
		return 0
	}

	if i == -1 {
		i = 0
	}

	size := PatHeaderSize + crc32.Size

	for i < len(p.Items) {
		is := len(p.Items[i].header)
		if (size + is) > PatMaximumSize {
			break
		} else {
			size += is
			i += 1
		}
	}

	return size
}

func (p *PAT) sectionHeader(i int) []byte {
	if i == -1 {
		p.header[6] = 0
	} else {
		p.header[6] += 1
	}

	return p.header[:PatHeaderSize]
}

func (p *PAT) sectionItem(i int) []byte {
	if i == -1 {
		return []byte{}
	}

	if i < len(p.Items) {
		return p.Items[i].header
	}

	return nil
}

func NewPatItem() *PatItem {
	p := new(PatItem)
	p.header = make([]byte, len(emptyPatItem))
	copy(p.header, emptyPatItem)

	return p
}

// Returns Program Number
func (p *PatItem) PNR() uint16 {
	return binary.BigEndian.Uint16(p.header[0:])
}

func (p *PatItem) SetPNR(pnr uint16) {
	binary.BigEndian.PutUint16(p.header[0:], pnr)
}

// Returns PID of the Program Map Table
func (p *PatItem) PID() PID {
	return getPID(p.header[2:])
}

func (p *PatItem) SetPID(pid PID) {
	setPID(p.header[2:], pid)
}
