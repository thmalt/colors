// Copied from ./gen/codegen/internal/interp/lerp.go.
// DO NOT EDIT. Changes will be overwritten.

package interp

import (
	"math"
	"strconv"
)

// HueInterpolation specifies how hue angles are interpolated.
type HueInterpolation uint8

const (
	HueShorter HueInterpolation = iota
	HueLonger
	HueIncreasing
	HueDecreasing
)

func (h HueInterpolation) String() string {
	switch h {
	case HueShorter:
		return "Shorter"
	case HueLonger:
		return "Longer"
	case HueIncreasing:
		return "Increasing"
	case HueDecreasing:
		return "Decreasing"
	default:
		return "HueInterpolation(" + strconv.FormatUint(uint64(h), 10) + ")"
	}
}

// Lerp linearly interpolates between a and b by t.
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t // or a*(1-t) + b*t
}

// LerpFMA linearly interpolates between a and b by t using fused multiply-add.
func LerpFMA(a, b, t float64) float64 {
	return math.FMA(b-a, t, a)
}

// LerpHue interpolates between two hue angles using the specified interpolation method.
func LerpHue(h1, h2 float64, t float64, hue HueInterpolation) float64 {
	switch hue {
	case HueLonger:
		return LerpHueLonger(h1, h2, t)
	case HueIncreasing:
		return LerpHueIncreasing(h1, h2, t)
	case HueDecreasing:
		return LerpHueDecreasing(h1, h2, t)
	default:
		return LerpHueShorter(h1, h2, t)
	}
}

// Hue Shorter interpolation.
//
//	h1 ∈ [0, 360)
//	h2 ∈ [0, 360)
//	t  ∈ [0, 1]
func LerpHueShorter(h1, h2, t float64) float64 {
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
func LerpHueLonger(h1, h2, t float64) float64 {
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
func LerpHueIncreasing(h1, h2, t float64) float64 {
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
func LerpHueDecreasing(h1, h2, t float64) float64 {
	if h1 < h2 {
		h1 += 360
	}

	h := h1 + (h2-h1)*t

	if h >= 360 {
		h -= 360
	}

	return h
}
