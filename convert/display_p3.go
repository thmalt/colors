package convert

// Applies the inverse of the Display P3 transfer function.
func DisplayP3ToLinearDisplayP3(r, g, b float64) (float64, float64, float64) {
	return srgbToLinearSrgb(r), srgbToLinearSrgb(g), srgbToLinearSrgb(b)
}

// Applies the Display P3 transfer function.
func LinearDisplayP3ToDisplayP3(r, g, b float64) (float64, float64, float64) {
	return linearSrgbToSrgb(r), linearSrgbToSrgb(g), linearSrgbToSrgb(b)
}
