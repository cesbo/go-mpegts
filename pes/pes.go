package pes

import (
	"encoding/binary"
)

type PES []byte

// Checks is PES prefix is 0x000001
func (p PES) CheckPrefix() bool {
	return (p[0] == 0) && (p[1] == 0) && (p[2] == 1)
}

// Specifies the type and number of the elementary stream.
// Audio streams (0xC0-0xDF), Video streams (0xE0-0xEF)
func (p PES) StreamID() uint8 {
	return p[3]
}

// Returns true if PES has an Elementary Stream data
func (p PES) IsES() bool {
	switch p.StreamID() {
	case 0xBC: // program_stream_map
		return false
	case 0xBE: // padding_stream
		return false
	case 0xBF: // private_stream_2
		return false
	case 0xF0: // ECM
		return false
	case 0xF1: // EMM
		return false
	case 0xF2: // DSMCC_stream
		return false
	case 0xF8: // ITU-T Rec. H.222.1 type E
		return false
	case 0xFF: // program_stream_directory
		return false
	}

	return true
}

// SetLength sets PES_packet_length field.
// Value specifying the number of bytes in the PES packet following the last byte
// of the field. 0 allowed only for video elementary stream.
func (p PES) SetLength(value int) {
	binary.BigEndian.PutUint16(p[4:], uint16(value))
}

// IsPTS checks is a Presentation Time Stamp (PTS) defined in the PES header.
// PTS field presents only for elementary streams.
func (p PES) IsPTS() bool {
	return (p[7] & 0x80) != 0
}

// GetPTS returns PTS value.
func (p PES) GetPTS() Timestamp {
	return (Timestamp(p[9]&0x0E) << 29) |
		(Timestamp(p[10]) << 22) |
		(Timestamp(p[11]&0xFE) << 14) |
		(Timestamp(p[12]) << 7) |
		(Timestamp(p[13]) >> 1)
}

// SetPTS sets PTS value.
func (p PES) SetPTS(value Timestamp) {
	_ = p[13]
	value &= MaxTimestamp

	p[9] = 0x20 | byte(value>>29) | 0x01
	p[10] = byte(value >> 22)
	p[11] = byte(value>>14) | 0x01
	p[12] = byte(value >> 7)
	p[13] = byte(value<<1) | 0x01
}

// IsDTS checks a Decoding TimeStamp (DTS) is defined in the PES header.
// DTS field presents only in pair with PTS field.
// If DTS field is not presented than DTS value equal to PTS.
func (p PES) IsDTS() bool {
	return (p[7] & 0x40) != 0
}

// GetDTS returns DTS value.
func (p PES) GetDTS() Timestamp {
	return (Timestamp(p[14]&0x0E) << 29) |
		(Timestamp(p[15]) << 22) |
		(Timestamp(p[16]&0xFE) << 14) |
		(Timestamp(p[17]) << 7) |
		(Timestamp(p[18]) >> 1)
}

// SetDTS sets DTS value.
func (p PES) SetDTS(value Timestamp) {
	_ = p[18]
	value &= MaxTimestamp

	p[9] |= 0x10

	p[14] = 0x10 | byte(value>>29) | 0x01
	p[15] = byte(value >> 22)
	p[16] = byte(value>>14) | 0x01
	p[17] = byte(value >> 7)
	p[18] = byte(value<<1) | 0x01
}
