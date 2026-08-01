package convert

// Conversion path (1 steps):
//
//	sRGB
//	-> HWB
func SrgbToHwb(r, g, b float64) (float64, float64, float64) {
	h, max, min := srgbToHueMaxMin(r, g, b)

	return h, min, 1 - max
}

// Conversion path (1 steps):
//
//	HWB
//	-> sRGB
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
