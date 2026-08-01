package convert

import "math"

// Transfer function:
//
//	Linear sRGB
func SrgbToLinear(r, g, b float64) (float64, float64, float64) {
	return srgbToLinear(r), srgbToLinear(g), srgbToLinear(b)
}

// Transfer function:
//
//	sRGB
func LinearToSrgb(r, g, b float64) (float64, float64, float64) {
	return linearToSrgb(r), linearToSrgb(g), linearToSrgb(b)
}

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

func SrgbToHsv(r, g, b float64) (float64, float64, float64) {
	h, max, min := srgbToHueMaxMin(r, g, b)

	delta := max - min
	if delta == 0 {
		return 0, 0, max
	}

	return h, delta / max, max
}

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

func SrgbToHwb(r, g, b float64) (float64, float64, float64) {
	h, max, min := srgbToHueMaxMin(r, g, b)

	return h, min, 1 - max
}

func HwbToSrgb(h, w, b float64) (red, green, blue float64) {
	if sum := w + b; sum >= 1 {
		gray := w / sum
		return gray, gray, gray
	}

	scale := 1 - w - b

	red, green, blue = HslToSrgb(h, 1, 0.5)

	red = red*scale + w
	green = green*scale + w
	blue = blue*scale + w

	return
}

// Gamma -> Linear
func srgbToLinear(x float64) float64 {
	neg := x < 0
	if neg {
		x = -x
	}

	if x > 0.04045 {
		x = math.Pow((x+0.055)/1.055, 2.4)
	} else {
		x /= 12.92
	}

	if neg {
		return -x
	}

	return x
}

// Linear -> Gamma
func linearToSrgb(x float64) float64 {
	neg := x < 0
	if neg {
		x = -x
	}

	if x > 0.0031308 {
		x = (1.055*math.Pow(x, 1.0/2.4) - 0.055)
	} else {
		x *= 12.92
	}

	if neg {
		return -x
	}

	return x
}

// func srgbToLinear(x float64) float64 {
// 	if x < 0 {
// 		x = -x
// 		if x > 0.04045 {
// 			return -math.Pow((x+0.055)/1.055, 2.4)
// 		}
// 		return -x / 12.92
// 	}
// 	if x > 0.04045 {
// 		return math.Pow((x+0.055)/1.055, 2.4)
// 	}
// 	return x / 12.92
// }

// func linearToSrgb(x float64) float64 {
// 	if x < 0 {
// 		x = -x
// 		if x > 0.0031308 {
// 			return -(1.055*math.Pow(x, 1.0/2.4) - 0.055)
// 		}
// 		return -x * 12.92
// 	}
// 	if x > 0.0031308 {
// 		return (1.055*math.Pow(x, 1.0/2.4) - 0.055)
// 	}
// 	return x * 12.92
// }

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
