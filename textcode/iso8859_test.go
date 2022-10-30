package textcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestISO8859_5(t *testing.T) {
	utf8 := "Привет!"
	coded := []byte{
		0xBF, 0xE0, 0xD8, 0xD2, 0xD5, 0xE2, 0x21,
	}

	t.Run("encode", func(t *testing.T) {
		result := EncodeISO8859_5(utf8)
		assert.Equal(t, coded, result)
	})

	t.Run("decode", func(t *testing.T) {
		result := DecodeISO8859_5(coded)
		assert.Equal(t, utf8, result)
	})
}

func TestISO8859_6(t *testing.T) {
	utf8 := "مرحبا"
	coded := []byte{
		0xE5, 0xD1, 0xCD, 0xC8, 0xC7,
	}

	t.Run("encode", func(t *testing.T) {
		result := EncodeISO8859_6(utf8)
		assert.Equal(t, coded, result)
	})

	t.Run("decode", func(t *testing.T) {
		result := DecodeISO8859_6(coded)
		assert.Equal(t, utf8, result)
	})
}

func TestISO8859_11(t *testing.T) {
	utf8 := "มีวันที่ดี!"
	coded := []byte{
		0xC1, 0xD5, 0xC7, 0xD1, 0xB9, 0xB7, 0xD5, 0xE8, 0xB4, 0xD5, 0x21,
	}

	t.Run("encode", func(t *testing.T) {
		result := EncodeISO8859_11(utf8)
		assert.Equal(t, coded, result)
	})

	t.Run("decode", func(t *testing.T) {
		result := DecodeISO8859_11(coded)
		assert.Equal(t, utf8, result)
	})
}
