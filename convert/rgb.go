package convert

import "math"

func RgbToSrgb(r, g, b float64) (float64, float64, float64) {
	return r / 255.0, g / 255.0, b / 255.0
}

func SrgbToRgb(r, g, b float64) (float64, float64, float64) {
	return clamp255(math.Round(r * 255)), clamp255(math.Round(g * 255)), clamp255(math.Round(b * 255))
}

func clamp(x, lo, hi float64) float64 {
	return min(hi, max(lo, x))
}

func clamp01(x float64) float64 {
	return clamp(x, 0, 1)
}

func clamp255(x float64) float64 {
	return clamp(x, 0, 255)
}
