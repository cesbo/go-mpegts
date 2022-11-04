package mpegts

import (
	"encoding/binary"
)

const (
	SyncByte   uint8  = 0x47
	PacketSize int    = 188
	MaxPID     uint16 = 8191
)

// ISO/IEC 13818-1 : 2.4.3 Specification of the Transport Stream syntax and semantics
type TS []byte

type ScramblingControl byte

const (
	NotScrambled     ScramblingControl = 0
	ScrambledEvenKey ScramblingControl = 2 // 10
	ScrambledOddKey  ScramblingControl = 3 // 11
)

var NullTS = TS{
	0x47, 0x1F, 0xFF, 0x10, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
}

// NewTS allocates new TS packet. Sets sync byte and PID
func NewTS(pid uint16) TS {
	ts := make(TS, PacketSize)
	ts[0] = SyncByte
	ts.SetPID(pid)
	return ts
}

// PID returns packet identifier value.
// PID is a 13-bit field that identifies the payload carried in the packet.
func (p TS) PID() uint16 {
	return binary.BigEndian.Uint16(p[1:]) & 0x1FFF
}

// SetPID sets PID in packet
func (p TS) SetPID(pid uint16) {
	pid = (pid & 0x1FFF) | (uint16(p[1]&0xE0) << 8)
	binary.BigEndian.PutUint16(p[1:], pid)
}

// CC returns continuity counter value.
// Continuity Counter is a 4-bit field incrementing by 1 for each packet
// with payload on the same PID.
func (p TS) CC() uint8 {
	return p[3] & 0x0F
}

// SetCC sets continuity counter in packet
func (p TS) SetCC(cc uint8) {
	p[3] = (p[3] & 0xF0) | (cc & 0x0F)
}

// IncrementCC increments continuity counter in packet
func (p TS) IncrementCC() {
	cc := p.CC() + 1
	p.SetCC(cc)
}

// CheckCC checks continuity counter
// Returns current CC and true if CC is equal to expected value
func (p TS) CheckCC(previous uint8) (uint8, bool) {
	cc := p.CC()
	expected := (previous + 1) & 0x0F
	return cc, cc == expected
}

// HasPayload checks is packet has payload
func (p TS) HasPayload() bool {
	return (p[3] & 0x10) != 0
}

// SetPayload sets Payload bit
func (p TS) SetPayload() {
	p[3] |= 0x10
}

// HasPUSI checks is payload starts in the packet (Payload Unit Start Indicator)
func (p TS) HasPUSI() bool {
	return (p[1] & 0x40) != 0
}

// SetPUSI sets Payload Unit Start Indicator bit
func (p TS) SetPUSI() {
	p[1] |= 0x40
}

// ClearPUSI clears Payload Unit Start Indicator bit
func (p TS) ClearPUSI() {
	p[1] &^= 0x40
}

// HeaderSize returns size of packet header
func (p TS) HeaderSize() int {
	if !p.HasAF() {
		return 4
	} else {
		return 4 + 1 + int(p[4])
	}
}

// Payload returns payload if available
func (p TS) Payload() []byte {
	if !p.HasPayload() || p.HasTEI() {
		return nil
	}

	s := p.HeaderSize()
	if s >= PacketSize {
		return nil
	}

	return p[s:PacketSize]
}

// HasTEI checks the Transport Error Indicator bit
func (p TS) HasTEI() bool {
	return (p[1] & 0x80) != 0
}

// HasAF checks the Adaptation Field bit
func (p TS) HasAF() bool {
	return (p[3] & 0x20) != 0
}

// SetAF sets Adaptation Field bit
func (p TS) SetAF() {
	p[3] |= 0x20
}

// ClearAF clears the Adaptation Field bit
func (p TS) ClearAF() {
	p[3] &^= 0x20
}

// TSC returns 2-bit Transport Scrambling Control field
func (p TS) TSC() ScramblingControl {
	if (p[3] & 0x80) != 0 {
		return ScramblingControl((p[3] & 0xC0) >> 6)
	} else {
		return NotScrambled
	}
}

// Fill fills incomplete TS packet with adaptation field stuffing bytes
func (p TS) Fill(size int) {
	headerSize := p.HeaderSize()

	payloadSize := size - headerSize
	offset := PacketSize - payloadSize
	copy(p[offset:], p[headerSize:size])

	if headerSize == 4 {
		// Set adaptation field
		p[3] |= 0x20
		headerSize += 1

		if headerSize < offset {
			p[5] = 0x00
			headerSize += 1
		}
	}

	for i := headerSize; i < offset; i++ {
		p[i] = 0xFF
	}

	p[4] = byte(PacketSize - 4 - 1 - payloadSize)
}
