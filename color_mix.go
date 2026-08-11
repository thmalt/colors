package colors

// Mix is shorthand for [Mix](c, other, t).
func (c Color) Mix(other Color, t float64) Color {
	return Mix(c, other, t)
}

// MixWith is shorthand for [MixWith](c, other, t, opts).
func (c Color) MixWith(other Color, t float64, opts InterpOptions) Color {
	return MixWith(c, other, t, opts)
}
