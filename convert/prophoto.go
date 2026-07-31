package convert

import "math"

// Transfer function:
//	Linear ProPhoto
func ProPhotoToLinear(r, g, b float64) (float64, float64, float64) {
	return proPhotoToLinear(r), proPhotoToLinear(g), proPhotoToLinear(b)
}

// Transfer function:
//	ProPhoto
func LinearToProPhoto(r, g, b float64) (float64, float64, float64) {
	return linearToProPhoto(r), linearToProPhoto(g), linearToProPhoto(b)
}

// Gamma -> Linear
func proPhotoToLinear(x float64) float64 {
	neg := x < 0
	if neg {
		x = -x
	}

	if x < 1.0/32.0 { // old: 0.031248
		x /= 16.0
	} else {
		x = math.Pow(x, 1.8)
	}

	if neg {
		return -x
	}

	return x
}

// Linear -> Gamma
func linearToProPhoto(x float64) float64 {
	neg := x < 0
	if neg {
		x = -x
	}

	if x > 1.0/512.0 {
		x = math.Pow(x, 1/1.8)
	} else {
		x *= 16
	}

	if neg {
		return -x
	}

	return x
}
