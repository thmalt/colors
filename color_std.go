package colors

import (
	"image/color"
)

// RGBA implements the [color.Color] interface
func (c Color) RGBA() (r, g, b, a uint32) {
	red, green, blue, alpha := c.ToRgba8()

	a = uint32(alpha) * 0x101
	r = uint32(red) * a / 0xff
	g = uint32(green) * a / 0xff
	b = uint32(blue) * a / 0xff

	return
}

// ToRGBA64 converts the color to an sRGB [color.RGBA64].
func (c Color) ToRGBA64() color.RGBA64 {
	red, green, blue, alpha := c.ToRgba8()

	a := uint32(alpha) * 0x101
	r := uint32(red) * a / 0xff
	g := uint32(green) * a / 0xff
	b := uint32(blue) * a / 0xff

	return color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

// FromStd converts a [color.Color] to a [Color].
func FromStd(c color.Color) Color {
	switch c := c.(type) {
	case color.Alpha:
		return fromAlpha(c.A)
	case *color.Alpha:
		return fromAlpha(c.A)

	case color.Alpha16:
		return fromAlpha16(c.A)
	case *color.Alpha16:
		return fromAlpha16(c.A)

	case color.Gray:
		return fromGray(c.Y)
	case *color.Gray:
		return fromGray(c.Y)

	case color.Gray16:
		return fromGray16(c.Y)
	case *color.Gray16:
		return fromGray16(c.Y)

	case color.NRGBA:
		return fromNRGBA(c.R, c.G, c.B, c.A)
	case *color.NRGBA:
		return fromNRGBA(c.R, c.G, c.B, c.A)

	case color.NRGBA64:
		return fromNRGBA64(c.R, c.G, c.B, c.A)
	case *color.NRGBA64:
		return fromNRGBA64(c.R, c.G, c.B, c.A)

	case color.RGBA:
		return fromRGBA(c.R, c.G, c.B, c.A)
	case *color.RGBA:
		return fromRGBA(c.R, c.G, c.B, c.A)

	case color.RGBA64:
		return fromRGBA64(c.R, c.G, c.B, c.A)
	case *color.RGBA64:
		return fromRGBA64(c.R, c.G, c.B, c.A)

	default:
		return fromColorRGBA(c.RGBA())
	}
}

func fromColorRGBA(R, G, B, A uint32) Color {
	if A == 0 {
		return SrgbAlpha(0, 0, 0, 0)
	}

	a := float64(A)
	r := float64(R) / a
	g := float64(G) / a
	b := float64(B) / a
	a /= 0xffff

	return SrgbAlpha(r, g, b, a)
}

func fromAlpha(A uint8) Color {
	a := float64(A) / 0xff
	return SrgbAlpha(1, 1, 1, a)
}

func fromAlpha16(A uint16) Color {
	a := float64(A) / 0xffff
	return SrgbAlpha(1, 1, 1, a)
}

func fromGray(Y uint8) Color {
	y := float64(Y) / 0xff
	return SrgbAlpha(y, y, y, 1)
}

func fromGray16(Y uint16) Color {
	y := float64(Y) / 0xffff
	return SrgbAlpha(y, y, y, 1)
}

func fromNRGBA(R, G, B, A uint8) Color {
	r := float64(R) / 0xff
	g := float64(G) / 0xff
	b := float64(B) / 0xff
	a := float64(A) / 0xff

	return SrgbAlpha(r, g, b, a)
}

func fromNRGBA64(R, G, B, A uint16) Color {
	r := float64(R) / 0xffff
	g := float64(G) / 0xffff
	b := float64(B) / 0xffff
	a := float64(A) / 0xffff

	return SrgbAlpha(r, g, b, a)
}

func fromRGBA(R, G, B, A uint8) Color {
	if A == 0 {
		return SrgbAlpha(0, 0, 0, 0)
	}

	a := float64(A)
	r := float64(R) / a
	g := float64(G) / a
	b := float64(B) / a
	a /= 0xff

	return SrgbAlpha(r, g, b, a)
}

func fromRGBA64(R, G, B, A uint16) Color {
	if A == 0 {
		return SrgbAlpha(0, 0, 0, 0)
	}

	a := float64(A)
	r := float64(R) / a
	g := float64(G) / a
	b := float64(B) / a
	a /= 0xffff

	return SrgbAlpha(r, g, b, a)
}
