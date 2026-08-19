package colors

import (
	"math"

	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/dither"
	"github.com/thmalt/colors/space"
)

// MustTo converts the color to the destination color space and panics if the conversion fails.
func (c Color) MustTo(dst space.Space) Color {
	to, err := c.To(dst)
	if err != nil {
		panic(err)
	}

	return to
}

// Clamp is shorthand for [Clamp](c).
func (c Color) Clamp() Color {
	return Clamp(c)
}

// InGamut reports whether the color is within the gamut of its color space.
func (c Color) InGamut() bool {
	return InGamut(c)
}

// InGamutSpace reports whether the color is within the gamut of the specified color space.
func (c Color) InGamutSpace(dst space.Space) bool {
	return InGamutSpace(c, dst)
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
		r = linearSrgbToSrgb(c.c1)
		g = linearSrgbToSrgb(c.c2)
		b = linearSrgbToSrgb(c.c3)
	default:
		r, g, b = c.LinearSrgb()
		r = linearSrgbToSrgb(r)
		g = linearSrgbToSrgb(g)
		b = linearSrgbToSrgb(b)
	}

	d := dither.Offset(x, y) * (1 / 255.0)
	return SrgbAlpha(clamp01(r+d), clamp01(g+d), clamp01(b+d), c.alpha)
}

// Rgb returns the color components in the RGB color space.
// Components are in the range [0, 255].
func (c Color) Rgb() (r, g, b float64) {
	return convert.SrgbToRgb(c.Srgb())
}

// Rgb returns a [Color] from 8-bit RGB components in [0, 255].
//
//	r: [0, 255]
//	g: [0, 255]
//	b: [0, 255]
func Rgb(r, g, b float64) Color {
	return Srgb(convert.RgbToSrgb(r, g, b))
}

// Rgb returns a [Color] from 8-bit RGB components in [0, 255] with alpha.
//
//	r: [0, 255]
//	g: [0, 255]
//	b: [0, 255]
//	alpha: [0, 1]
func RgbAlpha(r, g, b, alpha float64) Color {
	r, g, b = convert.RgbToSrgb(r, g, b)
	return SrgbAlpha(r, g, b, alpha)
}

// FromRgb8 creates a color from 8-bit sRGB components in linear sRGB.
func FromRgb8(r, g, b uint8) Color {
	return LinearSrgb(convert.Rgb8ToLinearSrgb(r, g, b))
}

// FromRgba8 creates a color from 8-bit sRGB components and alpha in linear sRGB.
func FromRgba8(r, g, b, a uint8) Color {
	c := LinearSrgb(convert.Rgb8ToLinearSrgb(r, g, b))
	c.alpha = float64(a) / 255
	return c
}

// ToRgb8 converts the color to 8-bit sRGB components.
func (c Color) ToRgb8() (r, g, b uint8) {
	switch c.space {
	case space.Srgb:
		r = uint8(clamp01(c.c1)*255 + 0.5)
		g = uint8(clamp01(c.c2)*255 + 0.5)
		b = uint8(clamp01(c.c3)*255 + 0.5)
		return
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint8(clamp01(fr)*255 + 0.5)
		g = uint8(clamp01(fg)*255 + 0.5)
		b = uint8(clamp01(fb)*255 + 0.5)
		return
	case space.LinearSrgb:
		return convert.LinearSrgbToRgb8(c.c1, c.c2, c.c3)
	default:
		return convert.LinearSrgbToRgb8(c.LinearSrgb())
	}
}

// ToRgba8 converts the color to 8-bit sRGB components with alpha.
func (c Color) ToRgba8() (r, g, b, a uint8) {
	a = c.Alpha8()
	switch c.space {
	case space.Srgb:
		r = uint8(clamp01(c.c1)*255 + 0.5)
		g = uint8(clamp01(c.c2)*255 + 0.5)
		b = uint8(clamp01(c.c3)*255 + 0.5)
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.Srgb()
		r = uint8(clamp01(fr)*255 + 0.5)
		g = uint8(clamp01(fg)*255 + 0.5)
		b = uint8(clamp01(fb)*255 + 0.5)
	case space.LinearSrgb:
		r, g, b = convert.LinearSrgbToRgb8(c.c1, c.c2, c.c3)
	default:
		r, g, b = convert.LinearSrgbToRgb8(c.LinearSrgb())
	}
	return
}

// linearSrgbToSrgb converts a linear sRGB component to sRGB using
// a Log/Exp power approximation for improved performance.
func linearSrgbToSrgb(x float64) float64 {
	const inv24 = 1 / 2.4

	neg := x < 0
	x = math.Abs(x)

	if x <= 0.0031308 {
		x *= 12.92
	} else {
		x = 1.055*math.Exp(math.Log(x)*inv24) - 0.055
	}

	if neg {
		return -x
	}
	return x
}
