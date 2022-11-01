package mpegts

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/cesbo/go-mpegts/crc32"
)

const (
	SdtHeaderSize  = 11
	SdtMaximumSize = 1024
	SdtItemSize    = 5
)

type SDT struct {
	header []byte
	Items  []*SdtItem
}

type SdtItem struct {
	header []byte
}

var (
	ErrSdtFormat = errors.New("sdt: invalid format")
)

var (
	emptySdt = []byte{
		0x42,        // table_id
		0xF0,        // section_length 1
		0x00,        // section_length 2
		0x00,        // transport_stream_id 1
		0x00,        // transport_stream_id 2
		0xC0 | 0x01, // version
		0x00,        // section_number
		0x00,        // last_section_number
		0x00,        // original_network_id 1
		0x00,        // original_network_id 2
		0xFF,        // reserved
	}
	emptySdtItem = []byte{
		0x00, // service_id 1
		0x00, // service_id 2
		0xFC, // EIT_schedule_flag | EIT_present_following_flag
		0x00, // running_status | free_CA_mode | descriptors_loop_length 1
		0x00, // descriptors_loop_length
	}
)

func NewSdt() *SDT {
	s := new(SDT)
	s.header = make([]byte, len(emptySdt))
	copy(s.header, emptySdt)

	return s
}

func (s *SDT) ParseSdtSection(b []byte) error {
	if len(b) < (SdtHeaderSize + crc32.Size) {
		return ErrSdtFormat
	}

	next := SdtHeaderSize
	end := len(b) - crc32.Size

	// copy header only from first section
	if b[6] == 0 {
		s.header = make([]byte, next)
		copy(s.header, b)
	}

	skip := next

	for skip < end {
		next += SdtItemSize
		if next > end {
			return ErrSdtFormat
		}

		descLen := binary.BigEndian.Uint16(b[skip+3:]) & 0x0FFF
		if descLen > 0 {
			next += int(descLen)
			if next > end {
				return ErrSdtFormat
			}

			desc := Descriptors(b[skip+SdtItemSize : next])
			if err := desc.Check(); err != nil {
				return fmt.Errorf("sdt: %w", err)
			}
		}

		item := new(SdtItem)
		item.header = make([]byte, next-skip)
		copy(item.header, b[skip:])

		s.Items = append(s.Items, item)

		skip = next
	}

	return nil
}

func (s *SDT) Actual() bool {
	return s.header[0] == 0x42
}

func (s *SDT) Version() uint8 {
	return (s.header[5] & 0x3E) >> 1
}

func (s *SDT) SetVersion(version uint8) {
	s.header[5] &^= 0x3E
	s.header[5] |= (version << 1) & 0x3E
}

func (s *SDT) TSID() uint16 {
	return binary.BigEndian.Uint16(s.header[3:5])
}

func (s *SDT) SetTSID(tsid uint16) {
	binary.BigEndian.PutUint16(s.header[3:], tsid)
}

func (s *SDT) ONID() uint16 {
	return binary.BigEndian.Uint16(s.header[8:10])
}

func (s *SDT) SetONID(onid uint16) {
	binary.BigEndian.PutUint16(s.header[8:], onid)
}

// Calculates LastSectionNumber
func (s *SDT) Finalize() {
	s.header[6] = 0
	s.header[7] = 0

	size := len(s.header)
	remain := SdtMaximumSize - size - crc32.Size

	for _, item := range s.Items {
		is := len(item.header)
		if is > remain {
			remain = SdtMaximumSize - SdtHeaderSize - crc32.Size
			s.header[7] += 1
		}
		remain -= is
	}
}

// Packetize splits SDT to TS packets
// packet - is a buffer to store data, on packet is ready,
// it is passed to the callback as byte slice
func (s *SDT) Packetize(packet TS, fn func([]byte)) error {
	return psiPacketize(s, packet, fn)
}

func (s *SDT) sectionSize(i int) int {
	if i == len(s.Items) {
		return 0
	}

	if i == -1 {
		i = 0
	}

	size := SdtHeaderSize + crc32.Size

	for i < len(s.Items) {
		is := len(s.Items[i].header)
		if (size + is) > SdtMaximumSize {
			break
		} else {
			size += is
			i += 1
		}
	}

	return size
}

func (s *SDT) sectionHeader(i int) []byte {
	if i == -1 {
		s.header[6] = 0
	} else {
		s.header[6] += 1
	}

	return s.header[:SdtHeaderSize]
}

func (s *SDT) sectionItem(i int) []byte {
	if i == -1 {
		return []byte{}
	}

	if i < len(s.Items) {
		return s.Items[i].header
	}

	return nil
}

func NewSdtItem() *SdtItem {
	p := new(SdtItem)
	p.header = make([]byte, len(emptySdtItem))
	copy(p.header, emptySdtItem)

	return p
}

func (s *SdtItem) PNR() uint16 {
	return binary.BigEndian.Uint16(s.header)
}

func (s *SdtItem) SetPNR(pnr uint16) {
	binary.BigEndian.PutUint16(s.header, pnr)
}

// Returns true if schedule information is present in the stream
// EIT_schedule_flag
func (s *SdtItem) IsSchedule() bool {
	return (s.header[2] & 0x02) != 0
}

func (s *SdtItem) SetSchedule(flag bool) {
	if flag {
		s.header[2] |= 0x02
	} else {
		s.header[2] &^= 0x02
	}
}

// Returns true if present-following information is present in the stream
// EIT_present_following_flag
func (s *SdtItem) IsPresentFollowing() bool {
	return (s.header[2] & 0x01) != 0
}

func (s *SdtItem) SetPresentFollowing(flag bool) {
	if flag {
		s.header[2] |= 0x01
	} else {
		s.header[2] &^= 0x01
	}
}

func (s *SdtItem) RunningStatus() uint8 {
	return (s.header[3] & 0xE0) >> 5
}

func (s *SdtItem) SetRunningStatus(status uint8) {
	s.header[3] &^= 0xE0
	s.header[3] |= (status << 5) & 0xE0
}

// Returns false if stream is not scrambled
// Returns true if one or more streams may be controlled by a CA system
// free_CA_mode
func (s *SdtItem) IsScrambled() bool {
	return (s.header[3] & 0x10) != 0
}

func (s *SdtItem) SetScrambled(flag bool) {
	if flag {
		s.header[3] |= 0x10
	} else {
		s.header[3] &^= 0x10
	}
}

func (s *SdtItem) Descriptors() Descriptors {
	return Descriptors(s.header[SdtItemSize:])
}

// Appends descriptors to the last added item
func (s *SdtItem) AppendDescriptors(desc Descriptors) {
	if len(desc) == 0 {
		return
	}

	ds := len(s.header) - SdtItemSize + len(desc)
	if ds > 0x0FFF {
		panic("failed to append sdt item descriptors: size limit")
	}
	s.header = append(s.header, desc...)

	b := uint16(s.header[3]&0xF0) << 8
	b |= uint16(ds)
	binary.BigEndian.PutUint16(s.header[3:], b)
}
