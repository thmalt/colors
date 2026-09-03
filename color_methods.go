package colors

import (
	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/dither"
	"github.com/thmalt/colors/space"
)

// To converts the color to the destination color space.
func (c Color) To(dst space.Space) Color {
	c.mutTo(dst)
	return c
}

// TryTo converts the color to destination color space
// and reports whether the conversion succeeded.
func (c Color) TryTo(dst space.Space) (Color, bool) {
	ok := c.mutTo(dst)
	return c, ok
}

// Mix is shorthand for [Mix](c, other, t).
func (c Color) Mix(other Color, t float64) Color {
	return Mix(c, other, t)
}

// MixWith is shorthand for [MixWith](c, other, t, opts).
func (c Color) MixWith(other Color, t float64, opts InterpOptions) Color {
	return MixWith(c, other, t, opts)
}

// Dither applies ordered dithering to the sRGB color channels at pixel position (x, y).
// The dithering offset is scaled to the normalized sRGB [0, 1] range;
// the alpha channel is left unchanged. The returned color is always in [space.Srgb].
func (c Color) Dither(x, y int) Color {
	var r, g, b float64

	switch c.space {
	case space.Srgb:
		r, g, b = c.c1, c.c2, c.c3
	case space.Hsl, space.Hsv, space.Hwb:
		r, g, b = c.Srgb()
	case space.LinearSrgb:
		r = convert.SrgbEncodeExp(c.c1)
		g = convert.SrgbEncodeExp(c.c2)
		b = convert.SrgbEncodeExp(c.c3)
	default:
		r, g, b = c.LinearSrgb()
		r = convert.SrgbEncodeExp(r)
		g = convert.SrgbEncodeExp(g)
		b = convert.SrgbEncodeExp(b)
	}

	d := dither.Offset(x, y) * invMaxUint8
	return SrgbAlpha(clamp01(r+d), clamp01(g+d), clamp01(b+d), c.alpha)
}

// Rgb returns the color components in the RGB color space.
// Components are in the range [0, 255].
func (c Color) Rgb() (r, g, b float64) {
	return srgbToRgb(c.Srgb())
}

// Rgb returns a [Color] from 8-bit RGB components in [0, 255].
//
//	r: [0, 255]
//	g: [0, 255]
//	b: [0, 255]
func Rgb(r, g, b float64) Color {
	return Srgb(rgbToSrgb(r, g, b))
}

// RgbAlpha returns a [Color] from 8-bit RGB components with alpha.
//
//	r: [0, 255]
//	g: [0, 255]
//	b: [0, 255]
//	alpha: [0, 1]
func RgbAlpha(r, g, b, alpha float64) Color {
	r, g, b = rgbToSrgb(r, g, b)
	return SrgbAlpha(r, g, b, alpha)
}

// FromRgb8 creates a color from 8-bit sRGB components in linear sRGB.
func FromRgb8(r, g, b uint8) Color {
	return LinearSrgb(convert.Rgb8ToLinearSrgb(r, g, b))
}

// FromRgba8 creates a color from 8-bit sRGB components and alpha in linear sRGB.
func FromRgba8(r, g, b, a uint8) Color {
	c := LinearSrgb(convert.Rgb8ToLinearSrgb(r, g, b))
	c.alpha = float64(a) * invMaxUint8
	return c
}

// ToRgb8 converts the color to 8-bit sRGB components.
func (c Color) ToRgb8() (r, g, b uint8) {
	switch c.space {
	case space.Srgb:
		r = uint8(clamp01(c.c1)*maxUint8 + 0.5)
		g = uint8(clamp01(c.c2)*maxUint8 + 0.5)
		b = uint8(clamp01(c.c3)*maxUint8 + 0.5)
		return
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint8(clamp01(fr)*maxUint8 + 0.5)
		g = uint8(clamp01(fg)*maxUint8 + 0.5)
		b = uint8(clamp01(fb)*maxUint8 + 0.5)
		return
	case space.LinearSrgb:
		r = convert.LinearSrgbToU8(c.c1)
		g = convert.LinearSrgbToU8(c.c2)
		b = convert.LinearSrgbToU8(c.c3)
		return
	default:
		fr, fg, fb := c.LinearSrgb()
		r = convert.LinearSrgbToU8(fr)
		g = convert.LinearSrgbToU8(fg)
		b = convert.LinearSrgbToU8(fb)
		return
	}
}

// ToRgba8 converts the color to 8-bit sRGB components with alpha.
func (c Color) ToRgba8() (r, g, b, a uint8) {
	a = uint8(clamp01(c.alpha)*maxUint8 + 0.5)
	switch c.space {
	case space.Srgb:
		r = uint8(clamp01(c.c1)*maxUint8 + 0.5)
		g = uint8(clamp01(c.c2)*maxUint8 + 0.5)
		b = uint8(clamp01(c.c3)*maxUint8 + 0.5)
		return
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint8(clamp01(fr)*maxUint8 + 0.5)
		g = uint8(clamp01(fg)*maxUint8 + 0.5)
		b = uint8(clamp01(fb)*maxUint8 + 0.5)
		return
	case space.LinearSrgb:
		r = convert.LinearSrgbToU8(c.c1)
		g = convert.LinearSrgbToU8(c.c2)
		b = convert.LinearSrgbToU8(c.c3)
		return
	default:
		fr, fg, fb := c.LinearSrgb()
		r = convert.LinearSrgbToU8(fr)
		g = convert.LinearSrgbToU8(fg)
		b = convert.LinearSrgbToU8(fb)
		return
	}
}

func rgbToSrgb(r, g, b float64) (float64, float64, float64) {
	return r * invMaxUint8, g * invMaxUint8, b * invMaxUint8
}

func srgbToRgb(r, g, b float64) (float64, float64, float64) {
	return r * maxUint8, g * maxUint8, b * maxUint8
}
