package colors

// Mix is shorthand for [Mix](c, other, t, opts).
func (c Color) Mix(other Color, t float64, opts MixOptions) Color {
	return Mix(c, other, t, opts)
}
