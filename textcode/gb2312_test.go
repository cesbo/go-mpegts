package textcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGB2312(t *testing.T) {
	utf8 := "天地玄黄 宇宙洪荒 日月盈昃 辰宿列张 寒来暑往 秋收冬藏"

	coded := []byte{
		0xCC, 0xEC, 0xB5, 0xD8, 0xD0, 0xFE, 0xBB, 0xC6, 0x20, 0xD3, 0xEE, 0xD6,
		0xE6, 0xBA, 0xE9, 0xBB, 0xC4, 0x20, 0xC8, 0xD5, 0xD4, 0xC2, 0xD3, 0xAF,
		0xEA, 0xBE, 0x20, 0xB3, 0xBD, 0xCB, 0xDE, 0xC1, 0xD0, 0xD5, 0xC5, 0x20,
		0xBA, 0xAE, 0xC0, 0xB4, 0xCA, 0xEE, 0xCD, 0xF9, 0x20, 0xC7, 0xEF, 0xCA,
		0xD5, 0xB6, 0xAC, 0xB2, 0xD8,
	}

	t.Run("encode", func(t *testing.T) {
		result := EncodeGB2312(utf8)
		assert.Equal(t, coded, result)
	})

	t.Run("decode", func(t *testing.T) {
		result := DecodeGB2312(coded)
		assert.Equal(t, utf8, result)
	})

	t.Run("out of bond", func(t *testing.T) {
		result := DecodeGB2312([]byte{0xF7, 0xFF})
		assert.Equal(t, "�", result)
	})
}
