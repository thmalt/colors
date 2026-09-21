package convert

// Srgb8Decode converts an 8-bit sRGB value to linear sRGB.
func Srgb8Decode(x uint8) float64 {
	return srgb8DecLUT[x]
}

// Srgb8Encode converts a linear sRGB value to an 8-bit sRGB value.
func Srgb8Encode(x float64) uint8 {
	const (
		size     = 1 << 8
		maxValue = size - 1
		scale    = 12.92 * maxValue
	)

	if x <= 0 {
		return 0
	}

	if x <= 0.0031308 {
		return uint8(x*scale + 0.5)
	}

	if x >= 1 {
		return maxValue
	}

	n := srgb8EncCLUT[uint8(x*size)]

	for n < maxValue && x >= srgb8EncTLUT[n+1] {
		n++
	}

	return n
}

// Srgb16Decode converts a 16-bit sRGB value to linear sRGB.
func Srgb16Decode(x uint16) float64 {
	return srgb16DecLUT[x]
}

// Srgb16Encode converts a linear sRGB value to a 16-bit sRGB value.
func Srgb16Encode(x float64) uint16 {
	const (
		size     = 1 << 16
		maxValue = size - 1
		scale    = 12.92 * maxValue
	)

	if x <= 0 {
		return 0
	}

	if x <= 0.0031308 {
		return uint16(x*scale + 0.5)
	}

	if x >= 1 {
		return maxValue
	}

	n := srgb16EncCLUT[uint16(x*size)]

	for n < maxValue && x >= srgb16EncTLUT[n+1] {
		n++
	}

	return n
}

// Srgb8ToLinearSrgb converts 8-bit sRGB components to linear sRGB.
func Srgb8ToLinearSrgb(r, g, b uint8) (float64, float64, float64) {
	return Srgb8Decode(r), Srgb8Decode(g), Srgb8Decode(b)
}

// LinearSrgbToSrgb8 converts linear sRGB components to 8-bit sRGB.
func LinearSrgbToSrgb8(r, g, b float64) (uint8, uint8, uint8) {
	return Srgb8Encode(r), Srgb8Encode(g), Srgb8Encode(b)
}

// Srgb16ToLinearSrgb converts 16-bit sRGB components to linear sRGB.
func Srgb16ToLinearSrgb(r, g, b uint16) (float64, float64, float64) {
	return Srgb16Decode(r), Srgb16Decode(g), Srgb16Decode(b)
}

// LinearSrgbToSrgb16 converts linear sRGB components to 16-bit sRGB.
func LinearSrgbToSrgb16(r, g, b float64) (uint16, uint16, uint16) {
	return Srgb16Encode(r), Srgb16Encode(g), Srgb16Encode(b)
}
