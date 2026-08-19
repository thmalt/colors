package colors

import (
	"image/color"
	"math"

	"github.com/thmalt/colors/space"
)

// RGBA implements the [color.Color] interface
func (c Color) RGBA() (r, g, b, a uint32) {
	const max = 65535

	alpha := clamp01(c.alpha)
	alphaMax := alpha * max
	a = uint32(alphaMax + 0.5)

	switch c.space {
	case space.Srgb:
		r = uint32(clamp01(c.c1)*alphaMax + 0.5)
		g = uint32(clamp01(c.c2)*alphaMax + 0.5)
		b = uint32(clamp01(c.c3)*alphaMax + 0.5)
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint32(clamp01(fr)*alphaMax + 0.5)
		g = uint32(clamp01(fg)*alphaMax + 0.5)
		b = uint32(clamp01(fb)*alphaMax + 0.5)
	case space.LinearSrgb:
		r = uint32(linearSrgbToRgb16(c.c1, alphaMax))
		g = uint32(linearSrgbToRgb16(c.c2, alphaMax))
		b = uint32(linearSrgbToRgb16(c.c3, alphaMax))
	default:
		fr, fg, fb := c.LinearSrgb()
		r = uint32(linearSrgbToRgb16(fr, alphaMax))
		g = uint32(linearSrgbToRgb16(fg, alphaMax))
		b = uint32(linearSrgbToRgb16(fb, alphaMax))
	}

	return
}

// ToRGBA64 converts the color to an sRGB [color.RGBA64].
func (c Color) ToRGBA64() color.RGBA64 {
	const max = 65535

	var r, g, b uint16
	alpha := clamp01(c.alpha)
	alphaMax := alpha * max
	a := uint16(alphaMax + 0.5)

	switch c.space {
	case space.Srgb:
		r = uint16(clamp01(c.c1)*alphaMax + 0.5)
		g = uint16(clamp01(c.c2)*alphaMax + 0.5)
		b = uint16(clamp01(c.c3)*alphaMax + 0.5)
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint16(clamp01(fr)*alphaMax + 0.5)
		g = uint16(clamp01(fg)*alphaMax + 0.5)
		b = uint16(clamp01(fb)*alphaMax + 0.5)
	case space.LinearSrgb:
		r = linearSrgbToRgb16(c.c1, alphaMax)
		g = linearSrgbToRgb16(c.c2, alphaMax)
		b = linearSrgbToRgb16(c.c3, alphaMax)
	default:
		fr, fg, fb := c.LinearSrgb()
		r = linearSrgbToRgb16(fr, alphaMax)
		g = linearSrgbToRgb16(fg, alphaMax)
		b = linearSrgbToRgb16(fb, alphaMax)
	}

	return color.RGBA64{R: r, G: g, B: b, A: a}
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

// linearSrgbToRgb16 converts a linear sRGB component to a premultiplied
// 16-bit sRGB value using Log/Exp for the power function.
func linearSrgbToRgb16(x float64, alphaMax float64) uint16 {
	const inv24 = 1 / 2.4

	x = clamp01(x)

	if x <= 0.0031308 {
		x *= 12.92
	} else {
		x = 1.055*math.Exp(math.Log(x)*inv24) - 0.055
	}

	return uint16(x*alphaMax + 0.5)
}
