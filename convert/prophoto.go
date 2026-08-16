package convert

import "math"

// Applies the inverse of the ProPhoto transfer function.
func ProPhotoToLinearProPhoto(r, g, b float64) (float64, float64, float64) {
	return proPhotoToLinearProPhoto(r), proPhotoToLinearProPhoto(g), proPhotoToLinearProPhoto(b)
}

// Applies the ProPhoto transfer function.
func LinearProPhotoToProPhoto(r, g, b float64) (float64, float64, float64) {
	return linearProPhotoToProPhoto(r), linearProPhotoToProPhoto(g), linearProPhotoToProPhoto(b)
}

// decode
func proPhotoToLinearProPhoto(x float64) float64 {
	neg := x < 0
	x = math.Abs(x)

	if x < 1/32.0 { // old: 0.031248
		x /= 16
	} else {
		x = math.Pow(x, 1.8)
	}

	if neg {
		return -x
	}
	return x
}

// encode
func linearProPhotoToProPhoto(x float64) float64 {
	neg := x < 0
	x = math.Abs(x)

	if x > 1/512.0 {
		x = math.Pow(x, 1/1.8)
	} else {
		x *= 16
	}

	if neg {
		return -x
	}
	return x
}
