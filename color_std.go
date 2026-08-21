package colors

import (
	"image/color"
	"math"

	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/space"
)

// RGBA implements the [color.Color] interface
func (c Color) RGBA() (r, g, b, a uint32) {
	alpha16 := clamp01(c.alpha) * maxUint16
	a = uint32(alpha16 + 0.5)

	switch c.space {
	case space.Srgb:
		r = uint32(clamp01(c.c1)*alpha16 + 0.5)
		g = uint32(clamp01(c.c2)*alpha16 + 0.5)
		b = uint32(clamp01(c.c3)*alpha16 + 0.5)
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint32(clamp01(fr)*alpha16 + 0.5)
		g = uint32(clamp01(fg)*alpha16 + 0.5)
		b = uint32(clamp01(fb)*alpha16 + 0.5)
	case space.LinearSrgb:
		r = uint32(lsrgb(clamp01(c.c1))*alpha16 + 0.5)
		g = uint32(lsrgb(clamp01(c.c2))*alpha16 + 0.5)
		b = uint32(lsrgb(clamp01(c.c3))*alpha16 + 0.5)
	default:
		fr, fg, fb := c.LinearSrgb()
		r = uint32(lsrgb(clamp01(fr))*alpha16 + 0.5)
		g = uint32(lsrgb(clamp01(fg))*alpha16 + 0.5)
		b = uint32(lsrgb(clamp01(fb))*alpha16 + 0.5)
	}

	return
}

// ToRGBA64 converts the color to an sRGB [color.RGBA64].
func (c Color) ToRGBA64() color.RGBA64 {
	var r, g, b uint16
	alpha16 := clamp01(c.alpha) * maxUint16
	a := uint16(alpha16 + 0.5)

	switch c.space {
	case space.Srgb:
		r = uint16(clamp01(c.c1)*alpha16 + 0.5)
		g = uint16(clamp01(c.c2)*alpha16 + 0.5)
		b = uint16(clamp01(c.c3)*alpha16 + 0.5)
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint16(clamp01(fr)*alpha16 + 0.5)
		g = uint16(clamp01(fg)*alpha16 + 0.5)
		b = uint16(clamp01(fb)*alpha16 + 0.5)
	case space.LinearSrgb:
		r = uint16(lsrgb(clamp01(c.c1))*alpha16 + 0.5)
		g = uint16(lsrgb(clamp01(c.c2))*alpha16 + 0.5)
		b = uint16(lsrgb(clamp01(c.c3))*alpha16 + 0.5)
	default:
		fr, fg, fb := c.LinearSrgb()
		r = uint16(lsrgb(clamp01(fr))*alpha16 + 0.5)
		g = uint16(lsrgb(clamp01(fg))*alpha16 + 0.5)
		b = uint16(lsrgb(clamp01(fb))*alpha16 + 0.5)
	}

	return color.RGBA64{R: r, G: g, B: b, A: a}
}

// ToNRGBA64 converts the color to an sRGB [color.NRGBA64].
func (c Color) ToNRGBA64() color.NRGBA64 {
	var r, g, b, a uint16
	a = uint16(clamp01(c.alpha)*maxUint16 + 0.5)

	switch c.space {
	case space.Srgb:
		r = uint16(clamp01(c.c1)*maxUint16 + 0.5)
		g = uint16(clamp01(c.c2)*maxUint16 + 0.5)
		b = uint16(clamp01(c.c3)*maxUint16 + 0.5)
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint16(clamp01(fr)*maxUint16 + 0.5)
		g = uint16(clamp01(fg)*maxUint16 + 0.5)
		b = uint16(clamp01(fb)*maxUint16 + 0.5)
	case space.LinearSrgb:
		r = uint16(lsrgb(clamp01(c.c1))*maxUint16 + 0.5)
		g = uint16(lsrgb(clamp01(c.c2))*maxUint16 + 0.5)
		b = uint16(lsrgb(clamp01(c.c3))*maxUint16 + 0.5)
	default:
		fr, fg, fb := c.LinearSrgb()
		r = uint16(lsrgb(clamp01(fr))*maxUint16 + 0.5)
		g = uint16(lsrgb(clamp01(fg))*maxUint16 + 0.5)
		b = uint16(lsrgb(clamp01(fb))*maxUint16 + 0.5)
	}

	return color.NRGBA64{R: r, G: g, B: b, A: a}
}

// ToRGBA converts the color to an sRGB [color.RGBA].
func (c Color) ToRGBA() color.RGBA {
	var r, g, b, a uint8
	alpha := clamp01(c.alpha)
	alpha8 := alpha * maxUint8
	a = uint8(alpha8 + 0.5)

	switch c.space {
	case space.Srgb:
		r = uint8(clamp01(c.c1)*alpha8 + 0.5)
		g = uint8(clamp01(c.c2)*alpha8 + 0.5)
		b = uint8(clamp01(c.c3)*alpha8 + 0.5)
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint8(clamp01(fr)*alpha8 + 0.5)
		g = uint8(clamp01(fg)*alpha8 + 0.5)
		b = uint8(clamp01(fb)*alpha8 + 0.5)
	case space.LinearSrgb:
		r, g, b = convert.LinearSrgbToRgb8(c.c1, c.c2, c.c3)
		r = uint8(float64(r)*alpha + 0.5)
		g = uint8(float64(g)*alpha + 0.5)
		b = uint8(float64(b)*alpha + 0.5)
	default:
		r, g, b = convert.LinearSrgbToRgb8(c.LinearSrgb())
		r = uint8(float64(r)*alpha + 0.5)
		g = uint8(float64(g)*alpha + 0.5)
		b = uint8(float64(b)*alpha + 0.5)
	}

	return color.RGBA{R: r, G: g, B: b, A: a}
}

// ToNRGBA converts the color to an sRGB [color.NRGBA].
func (c Color) ToNRGBA() color.NRGBA {
	var r, g, b, a uint8
	a = uint8(clamp01(c.alpha)*maxUint8 + 0.5)

	switch c.space {
	case space.Srgb:
		r = uint8(clamp01(c.c1)*maxUint8 + 0.5)
		g = uint8(clamp01(c.c2)*maxUint8 + 0.5)
		b = uint8(clamp01(c.c3)*maxUint8 + 0.5)
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint8(clamp01(fr)*maxUint8 + 0.5)
		g = uint8(clamp01(fg)*maxUint8 + 0.5)
		b = uint8(clamp01(fb)*maxUint8 + 0.5)
	case space.LinearSrgb:
		r, g, b = convert.LinearSrgbToRgb8(c.c1, c.c2, c.c3)
	default:
		r, g, b = convert.LinearSrgbToRgb8(c.LinearSrgb())
	}

	return color.NRGBA{R: r, G: g, B: b, A: a}
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
	invA := 1 / a

	r := float64(R) * invA
	g := float64(G) * invA
	b := float64(B) * invA
	a *= invMaxUint16

	return SrgbAlpha(r, g, b, a)
}

func fromAlpha(A uint8) Color {
	a := float64(A) * invMaxUint8
	return SrgbAlpha(1, 1, 1, a)
}

func fromAlpha16(A uint16) Color {
	a := float64(A) * invMaxUint16
	return SrgbAlpha(1, 1, 1, a)
}

func fromGray(Y uint8) Color {
	y := float64(Y) * invMaxUint8
	return SrgbAlpha(y, y, y, 1)
}

func fromGray16(Y uint16) Color {
	y := float64(Y) * invMaxUint16
	return SrgbAlpha(y, y, y, 1)
}

func fromNRGBA(R, G, B, A uint8) Color {
	r := float64(R) * invMaxUint8
	g := float64(G) * invMaxUint8
	b := float64(B) * invMaxUint8
	a := float64(A) * invMaxUint8

	return SrgbAlpha(r, g, b, a)
}

func fromNRGBA64(R, G, B, A uint16) Color {
	r := float64(R) * invMaxUint16
	g := float64(G) * invMaxUint16
	b := float64(B) * invMaxUint16
	a := float64(A) * invMaxUint16

	return SrgbAlpha(r, g, b, a)
}

func fromRGBA(R, G, B, A uint8) Color {
	if A == 0 {
		return SrgbAlpha(0, 0, 0, 0)
	}

	a := float64(A)
	invA := 1 / a

	r := float64(R) * invA
	g := float64(G) * invA
	b := float64(B) * invA
	a *= invMaxUint8

	return SrgbAlpha(r, g, b, a)
}

func fromRGBA64(R, G, B, A uint16) Color {
	if A == 0 {
		return SrgbAlpha(0, 0, 0, 0)
	}

	a := float64(A)
	invA := 1 / a

	r := float64(R) * invA
	g := float64(G) * invA
	b := float64(B) * invA
	a *= invMaxUint16

	return SrgbAlpha(r, g, b, a)
}

// lsrgb converts a clamped linear sRGB component to sRGB.
// x must be clamped to [0, 1], and the result is in [0, 1].
// It uses Log/Exp instead of Pow for better performance.
func lsrgb(x float64) float64 {
	const inv24 = 1 / 2.4

	if x <= 0.0031308 {
		return x * 12.92
	} else {
		// Use Log/Exp instead of Pow for better performance.
		return 1.055*math.Exp(math.Log(x)*inv24) - 0.055
	}
}
