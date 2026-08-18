package colors

import (
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
	r, g, b := c.Srgb()
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
