package convert

// Transfer function:
//
//	Linear Display P3
func DisplayP3ToLinear(r, g, b float64) (float64, float64, float64) {
	return srgbToLinear(r), srgbToLinear(g), srgbToLinear(b)
}

// Transfer function:
//
//	Display P3
func LinearToDisplayP3(r, g, b float64) (float64, float64, float64) {
	return linearToSrgb(r), linearToSrgb(g), linearToSrgb(b)
}
