package convert

import "math"

// Applies the inverse of the sRGB transfer function.
func SrgbToLinearSrgb(r, g, b float64) (float64, float64, float64) {
	return srgbToLinearSrgb(r), srgbToLinearSrgb(g), srgbToLinearSrgb(b)
}

// Applies the sRGB transfer function.
func LinearSrgbToSrgb(r, g, b float64) (float64, float64, float64) {
	return linearSrgbToSrgb(r), linearSrgbToSrgb(g), linearSrgbToSrgb(b)
}

// decode
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

// encode
func linearSrgbToSrgb(x float64) float64 {
	neg := x < 0
	x = math.Abs(x)

	if x <= 0.0031308 {
		x *= 12.92
	} else {
		x = 1.055*math.Pow(x, 1/2.4) - 0.055
	}

	if neg {
		return -x
	}
	return x
}
