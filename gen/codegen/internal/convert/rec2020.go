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

const (
	rec2020Alpha = 1.09929682680944
	rec2020Beta  = 0.018053968510807
)

// decode
func rec2020ToLinearRec2020(x float64) float64 {
	neg := x < 0
	if neg {
		x -= x
	}

	if x < 4.5*rec2020Beta {
		x = x / 4.5
	} else {
		x = math.Pow((x+rec2020Alpha-1.0)/rec2020Alpha, 1.0/0.45)
	}

	if neg {
		return -x
	}

	return x
}

// encode
func linearRec2020ToRec2020(x float64) float64 {
	neg := x < 0
	if neg {
		x -= x
	}

	if x < rec2020Beta {
		x = 4.5 * x
	} else {
		x = rec2020Alpha*math.Pow(x, 0.45) - (rec2020Alpha - 1.0)
	}

	if neg {
		return -x
	}

	return x
}
