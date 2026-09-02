package convert

import "math"

// A98Decode converts an Adobe RGB (1998) component to a linear A98 RGB component.
func A98Decode(x float64) float64 {
	if x < 0 {
		return -math.Pow(-x, 563.0/256)
	}
	return math.Pow(x, 563.0/256)
}

// A98Encode converts a linear A98 RGB component to an Adobe RGB (1998) component.
func A98Encode(x float64) float64 {
	if x < 0 {
		return -math.Pow(-x, 256.0/563)
	}
	return math.Pow(x, 256.0/563)
}
