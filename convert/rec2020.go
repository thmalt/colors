package convert

import "math"

// Transfer function:
//	Linear Rec.2020
func Rec2020ToLinear(r, g, b float64) (float64, float64, float64) {
	return rec2020ToLinear(r), rec2020ToLinear(g), rec2020ToLinear(b)
}

// Transfer function:
//	Rec. 2020
func LinearToRec2020(r, g, b float64) (float64, float64, float64) {
	return linearToRec2020(r), linearToRec2020(g), linearToRec2020(b)
}

// Gamma -> Linear
func rec2020ToLinear(x float64) float64 {
	if x < 0 {
		return -math.Pow(-x, 2.4)
	}

	return math.Pow(x, 2.4)
}

// Linear -> Gamma
func linearToRec2020(x float64) float64 {
	if x < 0 {
		return -math.Pow(-x, 1/2.4)
	}

	return math.Pow(x, 1/2.4)
}

// Gamma -> Linear
func rec2020ToLinear_(x float64) float64 {
	const (
		alpha = 1.09929682680944
		beta  = 0.018053968510807
	)

	if x < 0 {
		x = -x
		if x < beta*4.5 {
			return -x / 4.5
		}
		return -math.Pow((x+alpha-1)/alpha, 2.4)
	}

	if x < beta*4.5 {
		return x / 4.5
	}
	return math.Pow((x+alpha-1)/alpha, 2.4)
}

// Linear -> Gamma
func linearToRec2020_(x float64) float64 {
	const (
		alpha = 1.09929682680944
		beta  = 0.018053968510807
	)

	if x < 0 {
		x = -x
		if x > beta {
			return -alpha*math.Pow(x, 1/2.4) - (alpha - 1)
		}
		return -4.5 * x
	}

	if x > beta {
		return alpha*math.Pow(x, 1/2.4) - (alpha - 1)
	}
	return 4.5 * x
}
