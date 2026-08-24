package convert

import "math"

// decode
func a98ToLinearA98(x float64) float64 {
	return math.Pow(x, 563.0/256)
}

// encode
func linearA98ToA98(x float64) float64 {
	return math.Pow(x, 256.0/563)
}
