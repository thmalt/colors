package convert

import "math"

// decode
func rec2020ToLinearRec2020(x float64) float64 {
	if x < 0 {
		return -math.Pow(-x, 2.4)
	}
	return math.Pow(x, 2.4)
}

// encode
func linearRec2020ToRec2020(x float64) float64 {
	if x < 0 {
		return -math.Pow(-x, 1/2.4)
	}
	return math.Pow(x, 1/2.4)
}
