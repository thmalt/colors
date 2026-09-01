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

// Hex returns an sRGB color from a hexadecimal color string.
func Hex(s string) Color {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}

	switch len(s) {
	case 3:
		x0 := hexLUT[s[0]]
		x1 := hexLUT[s[1]]
		x2 := hexLUT[s[2]]

		if x0 == maxUint8 || x1 == maxUint8 || x2 == maxUint8 {
			return Color{}
		}

		r := float64(x0<<4|x0) * invMaxUint8
		g := float64(x1<<4|x1) * invMaxUint8
		b := float64(x2<<4|x2) * invMaxUint8

		return Srgb(r, g, b)
	case 4:
		x0 := hexLUT[s[0]]
		x1 := hexLUT[s[1]]
		x2 := hexLUT[s[2]]
		x3 := hexLUT[s[3]]

		if x0 == maxUint8 || x1 == maxUint8 || x2 == maxUint8 || x3 == maxUint8 {
			return Color{}
		}

		r := float64(x0<<4|x0) * invMaxUint8
		g := float64(x1<<4|x1) * invMaxUint8
		b := float64(x2<<4|x2) * invMaxUint8
		alpha := float64(x3<<4|x3) * invMaxUint8

		return SrgbAlpha(r, g, b, alpha)
	case 6:
		x0, x1 := hexLUT[s[0]], hexLUT[s[1]]
		x2, x3 := hexLUT[s[2]], hexLUT[s[3]]
		x4, x5 := hexLUT[s[4]], hexLUT[s[5]]

		if x0 == maxUint8 || x1 == maxUint8 || x2 == maxUint8 ||
			x3 == maxUint8 || x4 == maxUint8 || x5 == maxUint8 {
			return Color{}
		}

		r := float64(x0<<4|x1) * invMaxUint8
		g := float64(x2<<4|x3) * invMaxUint8
		b := float64(x4<<4|x5) * invMaxUint8

		return Srgb(r, g, b)
	case 8:
		x0, x1 := hexLUT[s[0]], hexLUT[s[1]]
		x2, x3 := hexLUT[s[2]], hexLUT[s[3]]
		x4, x5 := hexLUT[s[4]], hexLUT[s[5]]
		x6, x7 := hexLUT[s[6]], hexLUT[s[7]]

		if x0 == maxUint8 || x1 == maxUint8 || x2 == maxUint8 || x3 == maxUint8 ||
			x4 == maxUint8 || x5 == maxUint8 || x6 == maxUint8 || x7 == maxUint8 {
			return Color{}
		}

		r := float64(x0<<4|x1) * invMaxUint8
		g := float64(x2<<4|x3) * invMaxUint8
		b := float64(x4<<4|x5) * invMaxUint8
		alpha := float64(x6<<4|x7) * invMaxUint8

		return SrgbAlpha(r, g, b, alpha)
	default:
		return Color{}
	}
}
