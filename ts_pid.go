package mpegts

import (
	"encoding/binary"
)

// PID is Packet Identifier
type PID uint16

const (
	NonPid PID = 8192
	MaxPid PID = NonPid - 1
)

// getPID returns PID from bytes array
func getPID(b []byte) PID {
	return PID(binary.BigEndian.Uint16(b) & 0x1FFF)
}

// setPID sets PID in bytes array and keep first three bits
func setPID(b []byte, pid PID) {
	v := uint16(pid & 0x1FFF)
	v |= (uint16(b[0] & 0xE0)) << 8
	binary.BigEndian.PutUint16(b, v)
}
