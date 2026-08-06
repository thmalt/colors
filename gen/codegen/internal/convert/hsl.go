package convert

import "math"

// Conversion path (1 steps):
//
//	HSL
//	-> HWB
func SrgbToHsl(r, g, b float64) (h, s, l float64) {
	h, max, min := srgbToHueMaxMin(r, g, b)

	l = (min + max) * 0.5
	delta := max - min
	if delta == 0 {
		return 0, 0, l
	}

	if l > 0.5 {
		s = delta / (2 - max - min)
	} else {
		s = delta / (max + min)
	}

	return
}

// Conversion path (1 steps):
//
//	HSL
//	-> sRGB
func HslToSrgb(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l
	}

	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

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

	return r + m, g + m, b + m
}
