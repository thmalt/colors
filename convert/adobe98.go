package convert

import "math"

// Transfer function:
//	Linear Adobe® 98 RGB
func A98ToLinear(r, g, b float64) (float64, float64, float64) {
	return adobeRgbToLinear(r), adobeRgbToLinear(g), adobeRgbToLinear(b)
}

// Transfer function:
//	Adobe® 98 RGB
func LinearToA98(r, g, b float64) (float64, float64, float64) {
	return linearToAdobeRgb(r), linearToAdobeRgb(g), linearToAdobeRgb(b)
}

// Gamma -> Linear
func adobeRgbToLinear(x float64) float64 {
	return math.Pow(x, 563.0/256.0)
}

// Linear -> Gamma
func linearToAdobeRgb(x float64) float64 {
	return math.Pow(x, 256.0/563.0)
}
