package colors

// GradientBuilder builds a gradient from a sequence of color stops and hints.
type GradientBuilder struct {
	stops []GradientStop
}

// NewGradientBuilder returns a new GradientBuilder.
func NewGradientBuilder() *GradientBuilder {
	return &GradientBuilder{}
}

// AddStop adds a color stop to the builder.
//
// If no offsets are provided, the stop has no explicit position.
// If multiple offsets are provided, a stop is added at each offset.
func (b *GradientBuilder) AddStop(color Color, offsets ...float64) *GradientBuilder {
	if len(offsets) == 0 {
		b.stops = append(b.stops, NewStop(color))
		return b
	}

	for _, offset := range offsets {
		b.stops = append(b.stops, NewStopAt(color, offset))
	}

	return b
}

// AddHint adds a color interpolation hint at the given offset.
func (b *GradientBuilder) AddHint(offset float64) *GradientBuilder {
	b.stops = append(b.stops, NewHint(offset))
	return b
}

// Stops returns the gradient stops currently added to the builder
func (b *GradientBuilder) Stops() []GradientStop {
	return b.stops
}

// Build builds a Gradient using the builder's stops and default interpolation options.
func (b *GradientBuilder) Build() Gradient {
	return NewGradient(b.stops...)
}

// BuildWithOptions builds a Gradient using the builder's stops and the given interpolation options.
func (b *GradientBuilder) BuildWithOptions(opts InterpOptions) Gradient {
	return NewGradientWithOptions(opts, b.stops...)
}
