package colors

import "github.com/thmalt/colors/convert"

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
