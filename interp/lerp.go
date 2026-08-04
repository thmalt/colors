package interp

import "math"

type HueInterpolation uint8

const (
	HueShorter HueInterpolation = iota
	HueLonger
	HueIncreasing
	HueDecreasing
)

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t // a*(1-t) + b*t
}

func lerpFMA(a, b, t float64) float64 {
	return math.FMA(b-a, t, a)
}

// Hue Shorter interpolation.
//
//	h1 ∈ [0, 360)
//	h2 ∈ [0, 360)
//	t  ∈ [0, 1]
func lerpHueShorter(h1, h2, t float64) float64 {
	if d := h2 - h1; d > 180 {
		h1 += 360
	} else if d < -180 {
		h2 += 360
	}

	h := h1 + (h2-h1)*t

	if h >= 360 {
		h -= 360
	}

	return h
}

// Hue Longer interpolation.
//
//	h1 ∈ [0, 360)
//	h2 ∈ [0, 360)
//	t  ∈ [0, 1]
func lerpHueLonger(h1, h2, t float64) float64 {
	if d := h2 - h1; 0 < d && d < 180 {
		h1 += 360
	} else if -180 < d && d <= 0 {
		h2 += 360
	}

	h := h1 + (h2-h1)*t

	if h >= 360 {
		h -= 360
	}

	return h
}

// Hue Increasing interpolation.
//
//	h1 ∈ [0, 360)
//	h2 ∈ [0, 360)
//	t  ∈ [0, 1]
func lerpHueIncreasing(h1, h2, t float64) float64 {
	if h2 < h1 {
		h2 += 360
	}

	h := h1 + (h2-h1)*t

	if h >= 360 {
		h -= 360
	}

	return h
}

// Hue Decreasing interpolation.
//
//	h1 ∈ [0, 360)
//	h2 ∈ [0, 360)
//	t  ∈ [0, 1]
func lerpHueDecreasing(h1, h2, t float64) float64 {
	if h1 < h2 {
		h1 += 360
	}

	h := h1 + (h2-h1)*t

	if h >= 360 {
		h -= 360
	}

	return h
}
