// Package geohash provides encoding/decoding of base32 geohashes into coordinate pairs.
// From: https://en.wikipedia.org/wiki/Geohash
package geohash

import (
	"bytes"
)

var (
	//Base32 is the dictionary of characters for generating hashes
	base32 = []byte("0123456789bcdefghjkmnpqrstuvwxyz")
	// Bitmask positions for 5 bit base32 encoding
	// []int{ 0b10000, 0b01000, 0b00100, 0b00010, 0b00001 }
	bits = []int{16, 8, 4, 2, 1}
)

// Location is a coordinate pair of latitude and longitude (y, x)
type Location struct {
	lat, lon float64
}

// Encode a latitude/longitude pair into a geohash with the given precision.
func Encode(latitude, longitude float64, precision int) string {
	minLatitude, maxLatitude := -90.0, 90.0
	minLongitude, maxLongitude := -180.0, 180.0
	latitude = fixOutOfBounds(latitude, minLatitude, maxLatitude)
	longitude = fixOutOfBounds(longitude, minLongitude, maxLongitude)
	char, bit := 0, 0
	even := true
	var geohash bytes.Buffer
	// Encode to the given precision
	for geohash.Len() < precision {
		if even { // LONGITUDE
			mid := (minLongitude + maxLongitude) / 2
			if longitude > mid { // EAST
				char |= bits[bit]
				minLongitude = mid
			} else { // WEST
				maxLongitude = mid
			}
		} else { // LATITUDE
			mid := (minLatitude + maxLatitude) / 2
			if latitude > mid { // NORTH
				char |= bits[bit]
				minLatitude = mid
			} else { //SOUTH
				maxLatitude = mid
			}
		}
		even = !even // toggle lat/lon

		// Every 5 bits, encode a character and reset
		if bit < 4 {
			bit++
		} else {
			geohash.WriteByte(base32[char])
			char, bit = 0, 0
		}
	}
	return geohash.String()
}

// Rotates the map for out of bound coordinates
func fixOutOfBounds(num, min, max float64) float64 {
	if num < min {
		return max + (num - min)
	}
	if num > max {
		return min + (num - max)
	}
	return num
}
