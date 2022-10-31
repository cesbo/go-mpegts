package mpegts

import (
	"encoding/binary"
	"errors"

	"github.com/cesbo/go-mpegts/crc32"
)

const (
	// First 3 bytes of the PSI packet. Contains Table ID and Section Length
	PsiHeaderSize = 3

	// The maximum number of bytes in a section of
	// a ITU-T Rec. H.222.0 | ISO/IEC 13818-1 defined PSI table is 1024 bytes.
	// The maximum number of bytes in a private_section is 4096 bytes.
	// Includes PsiHeaderSize
	PsiMaximumSize = 4096
)

// Program Specific Information (ISO 13818-1 / 2.4.4)
type PSI struct {
	TableID           uint8
	Version           uint8
	SectionNumber     uint8
	LastSectionNumber uint8
	CRC               uint32

	cc byte

	buffer [PsiMaximumSize]byte // PSI buffer
	skip   int                  // bytes in buffer
	size   int                  // actual PSI size
}

var (
	ErrCC          = errors.New("psi: discontinuity received")
	ErrPUSI        = errors.New("psi: pointer field out of range")
	ErrAssemblePSI = errors.New("psi: assemble failed")
	ErrCRC         = errors.New("psi: checksum not match")
	ErrPsiFormat   = errors.New("psi: invalid format")
)

// PSI assembler callback
type AssembleFn func(error)

// Clears buffer
func (p *PSI) Clear() {
	p.skip = 0
	p.size = 0
}

// Returns PSI payload and assembling error status. Should be used in AssembleFn.
// Buffer length is equal to PsiHeaderSize + Section Length.
func (p *PSI) Payload() []byte {
	return p.buffer[:p.size]
}

// Base checks of PSI section:
// minSize - packet header size;
// maxSize - total size limit;
// crc - checksum validation.
func (p *PSI) commonCheck(minSize, maxSize int, crc bool) error {
	if crc {
		minSize += crc32.Size
	}

	if p.size < minSize {
		return ErrPsiFormat
	}

	if p.size > maxSize {
		return ErrPsiFormat
	}

	if crc {
		skip := p.size - crc32.Size
		crcActual := crc32.Checksum(0xFFFFFFFF, p.buffer[:skip])
		crcExpected := binary.BigEndian.Uint32(p.buffer[skip:])

		if crcActual != crcExpected {
			return ErrCRC
		}

		p.CRC = crcActual
	}

	p.TableID = p.buffer[0]
	p.Version = (p.buffer[5] >> 1) & 0x1F
	p.SectionNumber = p.buffer[6]
	p.LastSectionNumber = p.buffer[7]

	if p.SectionNumber > p.LastSectionNumber {
		return ErrPsiFormat
	}

	return nil
}

func (p *PSI) assembleCheck() error {
	switch p.buffer[0] {
	case 0x00: // PAT
		return p.commonCheck(PatHeaderSize, PatMaximumSize, true)
	case 0x01: // CAT
		return p.commonCheck(12, 1024, true)
	case 0x02: // PMT
		return p.commonCheck(PmtHeaderSize, PmtMaximumSize, true)
	case 0x42: // SDT Actual
		return p.commonCheck(SdtHeaderSize, SdtMaximumSize, true)
	case 0x46: // SDT Other
		return p.commonCheck(SdtHeaderSize, SdtMaximumSize, true)
	}

	return nil
}

func (p *PSI) callAssembleFn(fn AssembleFn, err error) {
	if err == nil {
		err = p.assembleCheck()
	}

	fn(err)

	p.Clear()
}

func (p *PSI) assembleStep(payload []byte) error {
	if p.size == 0 {
		// p.skip less than PsiHeaderSize if p.size == 0
		skip := copy(p.buffer[p.skip:PsiHeaderSize], payload)
		p.skip += skip

		if p.skip == PsiHeaderSize {
			p.size = p.getSectionLength()
			if p.size > PsiMaximumSize {
				return ErrAssemblePSI
			}
			payload = payload[skip:]
		} else if skip == p.skip {
			// after pointer field less than 3 bytes
			return nil
		} else {
			// 1 byte in buffer and 1 byte in payload
			return ErrAssemblePSI
		}
	}

	p.skip += copy(p.buffer[p.skip:p.size], payload)
	return nil
}

func (p *PSI) getSectionLength() int {
	_ = p.buffer[2]
	return PsiHeaderSize + int(binary.BigEndian.Uint16(p.buffer[1:])&0x0FFF)
}

// Assembles TS packets into single PSI.
// Calls fn when PSI is ready or error occurs
func (p *PSI) Assemble(packet TS, fn AssembleFn) {
	payload := packet.Payload()
	if payload == nil {
		return
	}

	if packet.IsPUSI() {
		remain := int(payload[0])
		payload = payload[1:]

		if remain >= len(payload) {
			p.callAssembleFn(fn, ErrPUSI)
			return
		}

		if p.skip != 0 {
			if packet.CC() != (p.cc+1)&0x0F {
				p.callAssembleFn(fn, ErrCC)
			} else if err := p.assembleStep(payload[:remain]); err != nil {
				p.callAssembleFn(fn, err)
			} else if p.skip == p.size {
				p.callAssembleFn(fn, nil)
			} else {
				p.callAssembleFn(fn, ErrAssemblePSI)
			}
		}

		payload = payload[remain:]
	} else {
		if p.skip == 0 {
			return
		}

		if packet.CC() != (p.cc+1)&0x0F {
			p.callAssembleFn(fn, ErrCC)
			return
		}
	}

	if err := p.assembleStep(payload); err != nil {
		p.callAssembleFn(fn, err)
		return
	}

	if p.skip == p.size {
		p.callAssembleFn(fn, nil)
	}

	p.cc = packet.CC()
}
