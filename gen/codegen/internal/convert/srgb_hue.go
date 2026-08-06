package convert

import "math"

func srgbHue(r, g, b, delta float64) (h float64) {
	if delta == 0 {
		return 0
	}

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
	h *= 60

	return h
}

func srgbToHueMaxMin(r, g, b float64) (h, max, min float64) {
	max = math.Max(r, math.Max(g, b))
	min = math.Min(r, math.Min(g, b))

	delta := max - min
	if delta == 0 {
		return 0, max, min
	}

	h = srgbHue(r, g, b, delta)
	return
}
