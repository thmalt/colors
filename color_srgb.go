package colors

import (
	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/space"
)

// Srgb8 returns a [Color] from 8-bit sRGB components.
func Srgb8(r, g, b uint8) Color {
	return Srgb(
		float64(r)*invMaxUint8,
		float64(g)*invMaxUint8,
		float64(b)*invMaxUint8,
	)
}

// Srgb8 returns a [Color] from 8-bit sRGB components and alpha.
func Srgba8(r, g, b, a uint8) Color {
	return SrgbAlpha(
		float64(r)*invMaxUint8,
		float64(g)*invMaxUint8,
		float64(b)*invMaxUint8,
		float64(a)*invMaxUint8,
	)
}

// Srgb8ToLinear returns a [Color] in linear sRGB from 8-bit sRGB components.
//
// It provides fast sRGB-to-linear conversion.
func Srgb8ToLinear(r, g, b uint8) Color {
	return LinearSrgb(convert.Srgb8ToLinearSrgb(r, g, b))
}

// Srgba8ToLinear returns a [Color] in linear sRGB from 8-bit sRGB components and alpha.
//
// It provides fast sRGB-to-linear conversion.
func Srgba8ToLinear(r, g, b, a uint8) Color {
	lr, lg, lb := convert.Srgb8ToLinearSrgb(r, g, b)
	return LinearSrgbAlpha(lr, lg, lb, float64(a)*invMaxUint8)
}

// Srgb8 returns the sRGB components of the [Color] as 8-bit values in the range [0, 255].
func (c Color) Srgb8() (r, g, b uint8) {
	return c.srgb8()
}

// Srgba8 returns the sRGB components and alpha of the [Color] as 8-bit values in the range [0, 255].
func (c Color) Srgba8() (r, g, b, a uint8) {
	return c.srgba8()
}

func (c *Color) srgb8() (r, g, b uint8) {
	switch c.space {
	case space.Srgb:
		r = uint8(clamp01(c.c1)*maxUint8 + 0.5)
		g = uint8(clamp01(c.c2)*maxUint8 + 0.5)
		b = uint8(clamp01(c.c3)*maxUint8 + 0.5)
		return
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.srgb()
		r = uint8(clamp01(fr)*maxUint8 + 0.5)
		g = uint8(clamp01(fg)*maxUint8 + 0.5)
		b = uint8(clamp01(fb)*maxUint8 + 0.5)
		return
	case space.LinearSrgb:
		r = convert.Srgb8Encode(c.c1)
		g = convert.Srgb8Encode(c.c2)
		b = convert.Srgb8Encode(c.c3)
		return
	default:
		fr, fg, fb := c.linearSrgb()
		r = convert.Srgb8Encode(fr)
		g = convert.Srgb8Encode(fg)
		b = convert.Srgb8Encode(fb)
		return
	}
}

func (c *Color) srgba8() (r, g, b, a uint8) {
	a = uint8(clamp01(c.alpha)*maxUint8 + 0.5)
	switch c.space {
	case space.Srgb:
		r = uint8(clamp01(c.c1)*maxUint8 + 0.5)
		g = uint8(clamp01(c.c2)*maxUint8 + 0.5)
		b = uint8(clamp01(c.c3)*maxUint8 + 0.5)
		return
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.srgb()
		r = uint8(clamp01(fr)*maxUint8 + 0.5)
		g = uint8(clamp01(fg)*maxUint8 + 0.5)
		b = uint8(clamp01(fb)*maxUint8 + 0.5)
		return
	case space.LinearSrgb:
		r = convert.Srgb8Encode(c.c1)
		g = convert.Srgb8Encode(c.c2)
		b = convert.Srgb8Encode(c.c3)
		return
	default:
		fr, fg, fb := c.linearSrgb()
		r = convert.Srgb8Encode(fr)
		g = convert.Srgb8Encode(fg)
		b = convert.Srgb8Encode(fb)
		return
	}
}

// Srgb16 returns a [Color] from 16-bit sRGB components.
func Srgb16(r, g, b uint16) Color {
	return Srgb(
		float64(r)*invMaxUint16,
		float64(g)*invMaxUint16,
		float64(b)*invMaxUint16,
	)
}

// Srgba16 returns a [Color] from 16-bit sRGB components and alpha.
func Srgba16(r, g, b, a uint16) Color {
	return SrgbAlpha(
		float64(r)*invMaxUint16,
		float64(g)*invMaxUint16,
		float64(b)*invMaxUint16,
		float64(a)*invMaxUint16,
	)
}

// Srgb16ToLinear returns a [Color] in linear sRGB from 16-bit sRGB components.
//
// It provides fast sRGB-to-linear conversion.
func Srgb16ToLinear(r, g, b uint16) Color {
	return LinearSrgb(convert.Srgb16ToLinearSrgb(r, g, b))
}

// Srgba16ToLinear returns a [Color] in linear sRGB from 16-bit sRGB components and alpha.
//
// It provides fast sRGB-to-linear conversion.
func Srgba16ToLinear(r, g, b, a uint16) Color {
	lr, lg, lb := convert.Srgb16ToLinearSrgb(r, g, b)
	return LinearSrgbAlpha(lr, lg, lb, float64(a)*invMaxUint16)
}

// Srgb16 returns the sRGB components of the [Color] as 16-bit values in the range [0, 65535].
func (c Color) Srgb16() (r, g, b uint16) {
	return c.srgb16()
}

// Srgba16 returns the sRGB components and alpha of the [Color] as 16-bit values in the range [0, 65535].
func (c Color) Srgba16() (r, g, b, a uint16) {
	return c.srgba16()
}

func (c *Color) srgb16() (r, g, b uint16) {
	switch c.space {
	case space.Srgb:
		r = uint16(clamp01(c.c1)*maxUint16 + 0.5)
		g = uint16(clamp01(c.c2)*maxUint16 + 0.5)
		b = uint16(clamp01(c.c3)*maxUint16 + 0.5)
		return
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.srgb()
		r = uint16(clamp01(fr)*maxUint16 + 0.5)
		g = uint16(clamp01(fg)*maxUint16 + 0.5)
		b = uint16(clamp01(fb)*maxUint16 + 0.5)
		return
	case space.LinearSrgb:
		r = convert.Srgb16Encode(c.c1)
		g = convert.Srgb16Encode(c.c2)
		b = convert.Srgb16Encode(c.c3)
		return
	default:
		fr, fg, fb := c.linearSrgb()
		r = convert.Srgb16Encode(fr)
		g = convert.Srgb16Encode(fg)
		b = convert.Srgb16Encode(fb)
		return
	}
}

func (c *Color) srgba16() (r, g, b, a uint16) {
	a = uint16(clamp01(c.alpha)*maxUint16 + 0.5)

	switch c.space {
	case space.Srgb:
		r = uint16(clamp01(c.c1)*maxUint16 + 0.5)
		g = uint16(clamp01(c.c2)*maxUint16 + 0.5)
		b = uint16(clamp01(c.c3)*maxUint16 + 0.5)
		return
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.srgb()
		r = uint16(clamp01(fr)*maxUint16 + 0.5)
		g = uint16(clamp01(fg)*maxUint16 + 0.5)
		b = uint16(clamp01(fb)*maxUint16 + 0.5)
		return
	case space.LinearSrgb:
		r = convert.Srgb16Encode(c.c1)
		g = convert.Srgb16Encode(c.c2)
		b = convert.Srgb16Encode(c.c3)
		return
	default:
		fr, fg, fb := c.linearSrgb()
		r = convert.Srgb16Encode(fr)
		g = convert.Srgb16Encode(fg)
		b = convert.Srgb16Encode(fb)
		return
	}
}
