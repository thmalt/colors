package colors

// Luminance returns the relative luminance of the color.
func Luminance(c Color) float64 {
	r, g, b := c.LinearSrgb()
	return 0.2126*r + 0.7152*g + 0.0722*b
}

// Contrast returns the WCAG 2 contrast ratio between two colors.
func Contrast(c1, c2 Color) float64 {
	l1 := Luminance(c1)
	l2 := Luminance(c2)

	darkest := min(l1, l2)
	lightest := max(l1, l2)

	return (lightest + 0.05) / (darkest + 0.05)
}

// Luminance returns the relative luminance of the color.
func (c Color) Luminance() float64 {
	return Luminance(c)
}

// Contrast returns the WCAG 2 contrast ratio between the color and another color.
func (c Color) Contrast(other Color) float64 {
	return Contrast(c, other)
}
