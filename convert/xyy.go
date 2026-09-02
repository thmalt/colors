package convert

// XyzToXyY converts CIE XYZ tristimulus values to xyY coordinates.
// This conversion only changes the coordinate representation; it does not
// involve chromatic adaptation or any other color-space transformation.
func XyzToXyY(x, y, z float64) (float64, float64, float64) {
	sum := x + y + z
	if sum == 0 {
		return 0, 0, 0
	}

	return x / sum, y / sum, y
}

// XyYToXyz converts xyY coordinates to CIE XYZ tristimulus values.
// This conversion only changes the coordinate representation; it does not
// involve chromatic adaptation or any other color-space transformation.
func XyYToXyz(x, y, Y float64) (float64, float64, float64) {
	if y == 0 {
		return 0, 0, 0
	}

	return x * Y / y, Y, (1 - x - y) * Y / y
}
