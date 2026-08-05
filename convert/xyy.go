package convert

// XyY is the CIE xyY color space using the D65 reference white.
// Conversions to color spaces with a different reference white automatically
// perform chromatic adaptation.

func XyzD65ToXyY(x, y, z float64) (float64, float64, float64) {
	sum := x + y + z
	if sum == 0 {
		return 0, 0, 0
	}

	return x / sum, y / sum, y
}

func XyYToXyzD65(x, y, Y float64) (float64, float64, float64) {
	if y == 0 {
		return 0, 0, 0
	}

	return x * Y / y, Y, (1 - x - y) * Y / y
}
