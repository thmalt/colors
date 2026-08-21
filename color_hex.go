package colors

import (
	"github.com/thmalt/colors/space"
)

const hexDigits = "0123456789abcdef"

func encodeHexByte(dst []byte, b byte) {
	dst[0] = hexDigits[b>>4]
	dst[1] = hexDigits[b&0x0f]
}

// Hex returns the hexadecimal representation of the color.
func (c Color) Hex() string {
	var r, g, b, a uint8

	if c.space == space.Srgb {
		r = uint8(clamp01(c.c1)*maxUint8 + 0.5)
		g = uint8(clamp01(c.c2)*maxUint8 + 0.5)
		b = uint8(clamp01(c.c3)*maxUint8 + 0.5)
		a = uint8(clamp01(c.alpha)*maxUint8 + 0.5)
	} else {
		r, g, b, a = c.ToRgba8()
	}

	var out [9]byte
	out[0] = '#'

	encodeHexByte(out[1:3], r)
	encodeHexByte(out[3:5], g)
	encodeHexByte(out[5:7], b)

	if a == maxUint8 {
		return string(out[:7])
	}

	encodeHexByte(out[7:9], a)

	return string(out[:])
}
