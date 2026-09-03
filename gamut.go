package colors

import "github.com/thmalt/colors/space"

// Clamp clamps the color channels to the valid range of the color space.
func (c Color) Clamp() Color {
	return Clamp(c)
}

// InGamut reports whether the color is within the gamut of its color space.
func (c Color) InGamut() bool {
	return InGamut(c)
}

// InGamutSpace reports whether the color is within the gamut of the specified color space.
func (c Color) InGamutSpace(dst space.Space) bool {
	return InGamutSpace(c, dst)
}

// MapToGamut maps the color to the gamut of the specified color space.
// It returns the original color if mapping fails.
func MapToGamut(c Color, dst space.Space) Color {
	c, _ = tryMapToGamutOklch(c, dst)
	return c
}

// MapToGamut maps the color to the gamut of the specified color space.
// It returns the original color if mapping fails.
func (c Color) MapToGamut(dst space.Space) Color {
	return MapToGamut(c, dst)
}

// TryMapToGamut maps the color to the gamut of the specified color space
// and reports whether the mapping succeeded.
func TryMapToGamut(c Color, dst space.Space) (Color, bool) {
	return tryMapToGamutOklch(c, dst)
}

// TryMapToGamut maps the color to the gamut of the specified color space
// and reports whether the mapping succeeded.
func (c Color) TryMapToGamut(dst space.Space) (Color, bool) {
	return TryMapToGamut(c, dst)
}

func tryMapToGamutOklch(c Color, dst space.Space) (Color, bool) {
	needConvert := c.space != dst

	if !needConvert && InGamut(c) {
		return c, true
	} else if d, ok := c.TryTo(dst); ok && InGamut(d) {
		return d, true
	}

	lightness, chroma, hue := c.Oklch()

	color := OklchAlpha(lightness, 0, hue, c.alpha)

	if needConvert && !color.mutTo(dst) {
		return c, false
	}

	if !InGamut(color) {
		return Clamp(color), true
	}

	low, high := 0.0, chroma

	for range 20 {
		mid := (low + high) * 0.5

		color.space = space.Oklch
		color.c1, color.c2, color.c3 = lightness, mid, hue

		if needConvert && !color.mutTo(dst) {
			return c, false
		}

		if InGamut(color) {
			low = mid
		} else {
			high = mid
		}
	}

	color.space = space.Oklch
	color.c1, color.c2, color.c3 = lightness, low, hue

	if needConvert && !color.mutTo(dst) {
		return c, false
	}

	return Clamp(color), true
}
