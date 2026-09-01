package colors

import "github.com/thmalt/colors/space"

// MapToGamut maps the color to the gamut of dst.
// It returns [false] if the color cannot be mapped.
func MapToGamut(c Color, dst space.Space) (Color, bool) {
	return mapToGamutOklch(c, dst)
}

func mapToGamutOklch(c Color, dst space.Space) (Color, bool) {
	if c.space == dst {
		if c.InGamut() {
			return c, true
		}
	} else if d, ok := c.to(dst); ok && d.InGamut() {
		return d, true
	}

	lightness, chroma, hue := c.Oklch()
	alpha := c.alpha

	neutral, ok := OklchAlpha(lightness, 0, hue, alpha).to(dst)
	if !ok {
		return c, false
	}
	if !neutral.InGamut() {
		return neutral.Clamp(), true
	}

	low, high := 0.0, chroma
	for range 20 {
		mid := (low + high) * 0.5

		candidate, ok := Oklch(lightness, mid, hue).to(dst)
		if !ok {
			return c, false
		}

		if candidate.InGamut() {
			low = mid
		} else {
			high = mid
		}
	}

	result, ok := OklchAlpha(lightness, low, hue, alpha).to(dst)
	if !ok {
		return c, false
	}

	return result.Clamp(), true
}
