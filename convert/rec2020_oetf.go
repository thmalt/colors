package convert

import "math"

// Applies the inverse of the Rec. 2020 OETF transfer function.
func Rec2020OETFToLinearRec2020(r, g, b float64) (float64, float64, float64) {
	return rec2020OETFToLinearRec2020(r), rec2020OETFToLinearRec2020(g), rec2020OETFToLinearRec2020(b)
}

// Applies the Rec. 2020 OETF transfer function.
func LinearRec2020ToRec2020OETF(r, g, b float64) (float64, float64, float64) {
	return linearRec2020ToRec2020OETF(r), linearRec2020ToRec2020OETF(g), linearRec2020ToRec2020OETF(b)
}

const (
	rec2020Alpha = 1.09929682680944
	rec2020Beta  = 0.018053968510807
)

// decode
func rec2020OETFToLinearRec2020(x float64) float64 {
	const (
		invGamma = 1 / 0.45
		invSlope = 1 / 4.5
		invAlpha = 1 / rec2020Alpha
	)

	neg := x < 0
	x = math.Abs(x)

	if x < 4.5*rec2020Beta {
		x *= invSlope
	} else {
		x = math.Pow((x+rec2020Alpha-1.0)*invAlpha, invGamma)
	}

	if neg {
		return -x
	}
	return x
}

// encode
func linearRec2020ToRec2020OETF(x float64) float64 {
	neg := x < 0
	x = math.Abs(x)

	if x < rec2020Beta {
		x *= 4.5
	} else {
		x = rec2020Alpha*math.Pow(x, 0.45) - (rec2020Alpha - 1.0)
	}

	if neg {
		return -x
	}
	return x
}
