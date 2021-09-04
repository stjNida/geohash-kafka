package geohash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var geohashTests = []struct {
	name      string
	latitude  float64
	longitude float64
	geohash   string
}{
	{"deneme-1", 48.6687, -4.3293, "gbsuv7z"},
	{"deneme-2", 56.878, 28.373, "udh5tdk"},
	{"deneme-3", 12.878, -3.373, "efm970r"},
}

// Geohash checking!
func TestEncode(t *testing.T) {
	for _, v := range geohashTests {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.geohash, Encode(v.latitude, v.longitude, len(v.geohash)))

		})

	}
}

func TestOutOfBounds(t *testing.T) {
	// Min
	assert.Equal(t, 0.0, fixOutOfBounds(-2.0, -1.0, 1.0))
	// Max
	assert.Equal(t, 0.0, fixOutOfBounds(2.0, -1.0, 1.0))
}
