package codegen

import "math"

func srgbToLinearSrgb(x float64) float64 {
	const (
		invSlope = 1 / 12.92
		invAlpha = 1 / 1.055
	)

	neg := x < 0
	x = math.Abs(x)

	if x <= 0.04045 {
		x *= invSlope
	} else {
		x = math.Pow((x+0.055)*invAlpha, 2.4)
	}

	if neg {
		return -x
	}
	return x
}
