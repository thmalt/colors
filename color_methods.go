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
		r, g, b = c.srgb()
	case space.LinearSrgb:
		r = convert.SrgbEncodeExp(c.c1)
		g = convert.SrgbEncodeExp(c.c2)
		b = convert.SrgbEncodeExp(c.c3)
	default:
		r, g, b = c.linearSrgb()
		r = convert.SrgbEncodeExp(r)
		g = convert.SrgbEncodeExp(g)
		b = convert.SrgbEncodeExp(b)
	}

	d := dither.Offset(x, y) * invMaxUint8

	c.space = space.Srgb
	c.c1 = clamp01(r + d)
	c.c2 = clamp01(g + d)
	c.c3 = clamp01(b + d)
	c.c4 = 0

	return c
}
