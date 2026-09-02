package colors

import (
	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/space"
)

// Luminance returns the relative luminance of the color.
func Luminance(c Color) float64 {
	return luminance(c)
}

// Contrast returns the WCAG 2 contrast ratio between two colors.
func Contrast(c1, c2 Color) float64 {
	l1 := luminance(c1)
	l2 := luminance(c2)

	if l1 < l2 {
		l1, l2 = l2, l1
	}

	return (l1 + 0.05) / (l2 + 0.05)
}

// Luminance returns the relative luminance of the color.
func (c Color) Luminance() float64 {
	return luminance(c)
}

// Contrast returns the WCAG 2 contrast ratio between the color and another color.
func (c Color) Contrast(other Color) float64 {
	return Contrast(c, other)
}

func luminance(c Color) float64 {
	var r, g, b float64

	switch c.space {
	case space.XyzD65:
		return c.c2
	case space.XyYD65:
		return c.c3
	case space.XyzAbsD65:
		return c.c2 * (1 / 203.0)
	case space.LinearSrgb:
		r, g, b = c.c1, c.c2, c.c3
	case space.Hsl, space.Hsv, space.Hwb:
		r, g, b = c.LinearSrgb()
	case space.Srgb:
		r = convert.SrgbDecodeExp(c.c1)
		g = convert.SrgbDecodeExp(c.c2)
		b = convert.SrgbDecodeExp(c.c3)
	default:
		_, y, _ := c.XyzD65()
		return y
	}

	return 0.21263900587151024*r + 0.715168678767756*g + 0.07219231536073371*b
}
