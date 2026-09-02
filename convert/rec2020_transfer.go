package convert

import "math"

// Rec2020Decode converts a Rec. 2020 component to a linear Rec. 2020 component.
func Rec2020Decode(x float64) float64 {
	if x < 0 {
		return -math.Pow(-x, 2.4)
	}
	return math.Pow(x, 2.4)
}

// Rec2020Encode converts a linear Rec. 2020 component to a Rec. 2020 component.
func Rec2020Encode(x float64) float64 {
	if x < 0 {
		return -math.Pow(-x, 1/2.4)
	}
	return math.Pow(x, 1/2.4)
}

const (
	rec2020Alpha = 1.09929682680944
	rec2020Beta  = 0.018053968510807
)

// Rec2020OETFDecode converts a Rec. 2020 OETF component to a linear Rec. 2020 component.
func Rec2020OETFDecode(x float64) float64 {
	const (
		invGamma = 1 / 0.45
		invSlope = 1 / 4.5
		invAlpha = 1 / rec2020Alpha
	)

	abs := math.Abs(x)

	if abs < 4.5*rec2020Beta {
		abs *= invSlope
	} else {
		abs = math.Pow((abs+rec2020Alpha-1.0)*invAlpha, invGamma)
	}

	if x < 0 {
		return -abs
	}
	return abs
}

// Rec2020OETFEncode converts a linear Rec. 2020 component to a Rec. 2020 OETF component.
func Rec2020OETFEncode(x float64) float64 {
	abs := math.Abs(x)

	if abs < rec2020Beta {
		abs *= 4.5
	} else {
		abs = rec2020Alpha*math.Pow(abs, 0.45) - (rec2020Alpha - 1.0)
	}

	if x < 0 {
		return -abs
	}
	return abs
}
