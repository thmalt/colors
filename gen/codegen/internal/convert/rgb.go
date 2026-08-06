package convert

// RgbToSrgb converts sRGB components from the range [0, 255] to [0, 1].
//
//	r: [0, 255]
//	g: [0, 255]
//	b: [0, 255]
func RgbToSrgb(r, g, b float64) (float64, float64, float64) {
	return r / 255.0, g / 255.0, b / 255.0
}

// SrgbToRgb converts sRGB components from the range [0, 1] to [0, 255].
//
//	r: [0, 1]
//	g: [0, 1]
//	b: [0, 1]
func SrgbToRgb(r, g, b float64) (float64, float64, float64) {
	return r * 255, g * 255, b * 255
}
