package convert

import "math"

// Applies the inverse of the Rec. 2020 transfer function.
func Rec2020ToLinearRec2020(r, g, b float64) (float64, float64, float64) {
	return rec2020ToLinearRec2020(r), rec2020ToLinearRec2020(g), rec2020ToLinearRec2020(b)
}

// Applies the Rec. 2020 transfer function.
func LinearRec2020ToRec2020(r, g, b float64) (float64, float64, float64) {
	return linearRec2020ToRec2020(r), linearRec2020ToRec2020(g), linearRec2020ToRec2020(b)
}

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
