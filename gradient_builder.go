package colors

import "math"

// GradientBuilder builds a gradient from a sequence of color stops and hints.
type GradientBuilder struct {
	stops []GradientStop
}

// NewGradientBuilder returns a new GradientBuilder.
func NewGradientBuilder() *GradientBuilder {
	return &GradientBuilder{}
}

// Reset removes all gradient stops and retains the allocated storage for reuse.
func (b *GradientBuilder) Reset() *GradientBuilder {
	b.stops = b.stops[:0]
	return b
}

// AddStop adds a color stop to the builder.
//
// If no offsets are provided, the stop has no explicit position.
// If multiple offsets are provided, a stop is added at each offset.
func (b *GradientBuilder) AddStop(color Color, offsets ...float64) *GradientBuilder {
	if len(offsets) == 0 {
		b.stops = append(b.stops, GradientStop{Color: color, Offset: math.NaN()})
		return b
	}

	for _, offset := range offsets {
		b.stops = append(b.stops, GradientStop{Color: color, Offset: offset})
	}

	return b
}

// AddHint adds a color interpolation hint at the given offset.
func (b *GradientBuilder) AddHint(offset float64) *GradientBuilder {
	b.stops = append(b.stops, GradientStop{Offset: offset})
	return b
}

// Stops returns the gradient stops currently added to the builder
func (b *GradientBuilder) Stops() []GradientStop {
	return b.stops
}

// Build builds a Gradient using the builder's stops and default interpolation options.
func (b *GradientBuilder) Build() Gradient {
	return newGradient(b.stops)
}

// BuildWithOptions builds a Gradient using the builder's stops and the given interpolation options.
func (b *GradientBuilder) BuildWithOptions(opts InterpOptions) Gradient {
	return newGradientWithOptions(opts, b.stops)
}
