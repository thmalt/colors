package colors

func (c Color) Mix(other Color, t float64, opts MixOptions) Color {
	return Mix(c, other, t, opts)
}
