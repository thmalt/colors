package convert

import "math"

// Conversion path (1 steps):
//
//	sRGB
//	-> HSV
func SrgbToHsv(r, g, b float64) (float64, float64, float64) {
	h, max, min := srgbToHueMaxMin(r, g, b)

	delta := max - min
	if delta == 0 {
		return 0, 0, max
	}

	return h, delta / max, max
}

// Conversion path (1 steps):
//
//	HSV
//	-> sRGB
func HsvToSrgb(h, s, v float64) (r, g, b float64) {
	if s == 0 {
		return v, v, v
	}

	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	r += m
	g += m
	b += m

	return
}
