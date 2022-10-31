package mpegts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp_Scale(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(Timestamp(22500), Scale(250, 1000))
}

func TestTimestamp_Delta(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		assert := assert.New(t)

		pts1 := Timestamp(30)
		pts2 := Timestamp(91)

		assert.Equal(Timestamp(61), pts2.Delta(pts1))
	})

	t.Run("overflow", func(t *testing.T) {
		assert := assert.New(t)

		pts1 := MaxTimestamp - 30
		pts2 := Timestamp(30)

		assert.Equal(Timestamp(61), pts2.Delta(pts1))
	})
}

func TestTimestamp_Add(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		assert := assert.New(t)

		pts1 := Timestamp(30)
		delta := Timestamp(61)

		assert.Equal(Timestamp(91), pts1.Add(delta))
	})

	t.Run("overflow", func(t *testing.T) {
		assert := assert.New(t)

		pts1 := MaxTimestamp - 30
		delta := Timestamp(61)

		assert.Equal(Timestamp(30), pts1.Add(delta))
	})
}
