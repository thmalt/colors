package convert

import "math"

// ProPhotoDecode converts a ProPhoto RGB component to linear ProPhoto RGB component.
func ProPhotoDecode(x float64) float64 {
	abs := math.Abs(x)

	if abs < 1/32.0 { // old: 0.031248
		abs /= 16
	} else {
		abs = math.Pow(abs, 1.8)
	}

	if x < 0 {
		return -abs
	}
	return abs
}

// ProPhotoEncode converts a linear ProPhoto RGB component to ProPhoto RGB component.
func ProPhotoEncode(x float64) float64 {
	abs := math.Abs(x)

	if abs > 1/512.0 {
		abs = math.Pow(abs, 1/1.8)
	} else {
		abs *= 16
	}

	if x < 0 {
		return -abs
	}
	return abs
}
