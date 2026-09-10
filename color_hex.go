package colors

import (
	"unsafe"

	"github.com/thmalt/colors/internal/hexlut"
	"github.com/thmalt/colors/space"
)

var nibbleToFloat = [16]float64{
	0 / 15.0, 1 / 15.0, 2 / 15.0, 3 / 15.0,
	4 / 15.0, 5 / 15.0, 6 / 15.0, 7 / 15.0,
	8 / 15.0, 9 / 15.0, 10 / 15.0, 11 / 15.0,
	12 / 15.0, 13 / 15.0, 14 / 15.0, 15 / 15.0,
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
		r, g, b, a = c.toRgba8()
	}

	var out [9]byte
	out[0] = '#'

	base := unsafe.Pointer(unsafe.StringData(hexlut.Enc))

	v := *(*[2]uint8)(unsafe.Add(base, uintptr(r)<<1))
	out[1], out[2] = v[0], v[1]

	v = *(*[2]uint8)(unsafe.Add(base, uintptr(g)<<1))
	out[3], out[4] = v[0], v[1]

	v = *(*[2]uint8)(unsafe.Add(base, uintptr(b)<<1))
	out[5], out[6] = v[0], v[1]

	v = *(*[2]uint8)(unsafe.Add(base, uintptr(a)<<1))
	out[7], out[8] = v[0], v[1]

	n := 9
	if a == maxUint8 {
		n = 7
	}

	return string(out[:n])
}

// Hex returns an sRGB color from a hexadecimal color string.
func Hex(s string) Color {
	c, _ := TryHex(s)
	return c
}

// TryHex returns an sRGB color from a hexadecimal color string
// and reports whether the string was successfully parsed.
func TryHex(s string) (Color, bool) {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}

	switch len(s) {
	case 3:
		x0 := hexlut.Dec[s[0]]
		x1 := hexlut.Dec[s[1]]
		x2 := hexlut.Dec[s[2]]

		if x0|x1|x2 >= hexlut.NibbleLimit {
			return Color{}, false
		}

		base := unsafe.Pointer(&nibbleToFloat[0])
		r := *(*float64)(unsafe.Add(base, uintptr(x0)*8))
		g := *(*float64)(unsafe.Add(base, uintptr(x1)*8))
		b := *(*float64)(unsafe.Add(base, uintptr(x2)*8))

		return Srgb(r, g, b), true
	case 4:
		x0 := hexlut.Dec[s[0]]
		x1 := hexlut.Dec[s[1]]
		x2 := hexlut.Dec[s[2]]
		x3 := hexlut.Dec[s[3]]

		if x0|x1|x2|x3 >= hexlut.NibbleLimit {
			return Color{}, false
		}

		base := unsafe.Pointer(&nibbleToFloat[0])
		r := *(*float64)(unsafe.Add(base, uintptr(x0)*8))
		g := *(*float64)(unsafe.Add(base, uintptr(x1)*8))
		b := *(*float64)(unsafe.Add(base, uintptr(x2)*8))
		alpha := *(*float64)(unsafe.Add(base, uintptr(x3)*8))

		return SrgbAlpha(r, g, b, alpha), true
	case 6:
		x0, x1 := hexlut.Dec[s[0]], hexlut.Dec[s[1]]
		x2, x3 := hexlut.Dec[s[2]], hexlut.Dec[s[3]]
		x4, x5 := hexlut.Dec[s[4]], hexlut.Dec[s[5]]

		if x0|x1|x2|x3|x4|x5 >= hexlut.NibbleLimit {
			return Color{}, false
		}

		r := float64(x0<<4|x1) * invMaxUint8
		g := float64(x2<<4|x3) * invMaxUint8
		b := float64(x4<<4|x5) * invMaxUint8

		return Srgb(r, g, b), true
	case 8:
		x0, x1 := hexlut.Dec[s[0]], hexlut.Dec[s[1]]
		x2, x3 := hexlut.Dec[s[2]], hexlut.Dec[s[3]]
		x4, x5 := hexlut.Dec[s[4]], hexlut.Dec[s[5]]
		x6, x7 := hexlut.Dec[s[6]], hexlut.Dec[s[7]]

		if x0|x1|x2|x3|x4|x5|x6|x7 >= hexlut.NibbleLimit {
			return Color{}, false
		}

		r := float64(x0<<4|x1) * invMaxUint8
		g := float64(x2<<4|x3) * invMaxUint8
		b := float64(x4<<4|x5) * invMaxUint8
		alpha := float64(x6<<4|x7) * invMaxUint8

		return SrgbAlpha(r, g, b, alpha), true
	default:
		return Color{}, false
	}
}
