package convert

import "math"

// SrgbDecode converts a sRGB component to linear sRGB component.
func SrgbDecode(x float64) float64 {
	const (
		invSlope = 1 / 12.92
		invAlpha = 1 / 1.055
	)

	abs := math.Abs(x)

	if abs <= 0.04045 {
		abs *= invSlope
	} else {
		abs = math.Pow((abs+0.055)*invAlpha, 2.4)
	}

	if x < 0 {
		return -abs
	}
	return abs
}

// SrgbEncode converts a linear sRGB component to sRGB component.
func SrgbEncode(x float64) float64 {
	abs := math.Abs(x)

	if abs <= 0.0031308 {
		abs *= 12.92
	} else {
		abs = 1.055*math.Pow(abs, 1/2.4) - 0.055
	}

	if x < 0 {
		return -abs
	}
	return abs
}

// SrgbDecodeExp converts a sRGB component to linear sRGB component
// using a Log/Exp power approximation for improved performance.
func SrgbDecodeExp(x float64) float64 {
	const (
		invSlope = 1 / 12.92
		invAlpha = 1 / 1.055
	)

	abs := math.Abs(x)

	if abs <= 0.04045 {
		abs *= invSlope
	} else if abs == 1 {
		return x
	} else {
		abs = math.Exp(math.Log((abs+0.055)*invAlpha) * 2.4)
	}

	if x < 0 {
		return -abs
	}
	return abs
}

// SrgbEncodeExp converts a linear sRGB component to sRGB component
// using a Log/Exp power approximation for improved performance.
func SrgbEncodeExp(x float64) float64 {
	const inv24 = 1 / 2.4

	abs := math.Abs(x)

	if abs <= 0.0031308 {
		abs *= 12.92
	} else if abs == 1 {
		return x
	} else {
		abs = 1.055*math.Exp(math.Log(abs)*inv24) - 0.055
	}

	if x < 0 {
		return -abs
	}
	return abs
}
