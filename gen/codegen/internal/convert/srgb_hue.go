package convert

import "math"

// Conversion path (1 steps):
//
//	HSL
//	-> HWB
func SrgbToHsl(r, g, b float64) (h, s, l float64) {
	max := max(r, g, b)
	min := min(r, g, b)

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

	return srgbHue(r, g, b, delta), s, l
}

// Conversion path (1 steps):
//
//	sRGB
//	-> HSV
func SrgbToHsv(r, g, b float64) (h, s, v float64) {
	max := max(r, g, b)
	min := min(r, g, b)

	delta := max - min
	if delta == 0 {
		return 0, 0, max
	}

	return srgbHue(r, g, b, delta), delta / max, max
}

// Conversion path (1 steps):
//
//	sRGB
//	-> HWB
func SrgbToHwb(r, g, b float64) (float64, float64, float64) {
	max := max(r, g, b)
	min := min(r, g, b)

	delta := max - min
	if delta == 0 {
		return 0, min, 1 - max
	}

	return srgbHue(r, g, b, delta), min, 1 - max
}

func srgbHue(r, g, b, delta float64) (h float64) {
	if r >= g && r >= b {
		h = (g - b) / delta
		if g < b {
			h += 6
		}
	} else if g >= b {
		h = (b-r)/delta + 2
	} else {
		h = (r-g)/delta + 4
	}

	return h * 60
}

func wrap360(x float64) float64 {
	const inv360 = 1 / 360.0
	return x - math.Floor(x*inv360)*360
}
