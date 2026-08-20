package convert

import "math"

// Conversion path (1 steps):
//
//	HSL
//	-> sRGB
func HslToSrgb(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l
	}

	if h < 0 || h >= 360 {
		h = wrap360(h)
	}

	h *= 1 / 60.0
	sector := int(h)
	f := h - float64(sector)

	c := (1 - math.Abs(2*l-1)) * s
	m := l - c*0.5

	x := c * f
	if sector&1 != 0 {
		x = c - x
	}

	switch sector {
	case 0:
		r, g = c, x
	case 1:
		r, g = x, c
	case 2:
		g, b = c, x
	case 3:
		g, b = x, c
	case 4:
		r, b = x, c
	default:
		r, b = c, x
	}

	return r + m, g + m, b + m
}

// Conversion path (1 steps):
//
//	HSL
//	-> HSV
func HslToHsv(h, s, l float64) (float64, float64, float64) {
	v := l + s*min(l, 1-l)

	if v == 0 {
		return h, 0, v
	}

	return h, 2 * (v - l) / v, v
}

// Conversion path (1 steps):
//
//	HSL
//	-> HWB
func HslToHwb(h, s, l float64) (float64, float64, float64) {
	c := (1 - math.Abs(2*l-1)) * s
	w := l - c*0.5

	return h, w, 1 - w - c
}
