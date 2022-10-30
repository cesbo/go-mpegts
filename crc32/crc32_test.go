package crc32

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Checksum(t *testing.T) {
	crc := Checksum(0xFFFFFFFF, []byte("123456789"))

	assert := assert.New(t)
	assert.Equal(uint32(0x0376E6E7), crc)
}
