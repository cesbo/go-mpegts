package ts

import (
	"bytes"
	"errors"
)

// Slicer is a tool for split TS buffer to TS packets.
// If buffer length is not multiple to TS packet size it store remain
// data in the internal buffer and will be used on the next iteration.
type Slicer struct {
	packet [188]byte
	fill   int

	buffer []byte
	skip   int

	err error
}

var (
	ErrSyncTS      = errors.New("ts slicer: sync error")
	ErrNotComplete = errors.New("ts slicer: not complete")
)

// Prepares buffer and get first packet
func (s *Slicer) Begin(buffer []byte) TS {
	s.buffer = buffer
	s.skip = 0
	s.err = nil

	// some data remain in the buffer after previous iteration
	if s.fill != 0 {
		size := copy(s.packet[s.fill:], s.buffer)
		s.fill += size
		s.skip += size

		if s.fill != PacketSize {
			return nil
		}

		s.fill = 0
		return TS(s.packet[:])
	}

	if len(s.buffer) == 0 {
		return nil
	}

	// check TS sync byte
	s.skip = bytes.IndexByte(s.buffer, SyncByte)
	if s.skip == -1 {
		s.err = ErrSyncTS
		return nil
	}

	return s.Next()
}

// Get next packet
func (s *Slicer) Next() TS {
	next := s.skip + PacketSize
	if len(s.buffer) >= next {
		p := s.buffer[s.skip:next]
		s.skip = next
		return TS(p)
	}

	if len(s.buffer) > s.skip {
		s.fill = copy(s.packet[:], s.buffer[s.skip:])
		s.skip += s.fill
	}

	return nil
}

// Returns number of bytes processed and error if happens
func (s *Slicer) Err() error {
	if s.err != nil {
		return s.err
	}

	if s.skip != len(s.buffer) {
		return ErrNotComplete
	}

	return nil
}
