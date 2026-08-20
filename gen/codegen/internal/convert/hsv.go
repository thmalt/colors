package convert

// Conversion path (1 steps):
//
//	HSV
//	-> sRGB
func HsvToSrgb(h, s, v float64) (r, g, b float64) {
	if s == 0 {
		return v, v, v
	}

	if h < 0 || h >= 360 {
		h = wrap360(h)
	}

	h *= 1 / 60.0
	sector := int(h)
	f := h - float64(sector)

	c := v * s
	m := v - c

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
//	HSV
//	-> HSL
func HsvToHsl(h, s, v float64) (float64, float64, float64) {
	l := v * (1 - s*0.5)

	if l == 0 || l == 1 {
		return h, 0, l
	}

	return h, (v - l) / min(l, 1-l), l
}

// Conversion path (1 steps):
//
//	HSV
//	-> HWB
func HsvToHwb(h, s, v float64) (float64, float64, float64) {
	w := v * (1 - s)
	b := 1 - v
	return h, w, b
}
