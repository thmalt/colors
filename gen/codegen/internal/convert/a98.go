package convert

import "math"

// Applies the inverse of the Adobe RGB (1998) transfer function.
func A98ToLinearA98(r, g, b float64) (float64, float64, float64) {
	return a98ToLinearA98(r), a98ToLinearA98(g), a98ToLinearA98(b)
}

// Applies the Adobe RGB (1998) transfer function.
func LinearA98ToA98(r, g, b float64) (float64, float64, float64) {
	return linearA98ToA98(r), linearA98ToA98(g), linearA98ToA98(b)
}

// decode
func a98ToLinearA98(x float64) float64 {
	return math.Pow(x, 563.0/256)
}

// encode
func linearA98ToA98(x float64) float64 {
	return math.Pow(x, 256.0/563)
}
