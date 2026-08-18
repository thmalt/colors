package colors

import (
	"math"
	"slices"
	"sort"
)

type GradientStop struct {
	Color  Color
	Offset float64

	invRange float64
}

type Spread uint8

const (
	SpreadPad Spread = iota
	SpreadRepeat
	SpreadReflect
)

type Gradient struct {
	mixer Mixer

	stops []GradientStop
}

func NewGradientWithOptions(opts InterpOptions, stops ...GradientStop) Gradient {
	mixer := NewMixerWithOptions(opts)

	resolved := resolveStops(slices.Clone(stops), mixer)

	return Gradient{
		mixer: mixer,
		stops: resolved,
	}
}

// NewGradient returns a new [Gradient] using [DefaultInterpOptions].
func NewGradient(stops ...GradientStop) Gradient {
	return NewGradientWithOptions(DefaultInterpOptions(), stops...)
}

func NewHint(offset float64) GradientStop {
	return GradientStop{
		Color:  Color{},
		Offset: offset,
	}
}

func NewStop(color Color) GradientStop {
	return GradientStop{
		Color:  color,
		Offset: math.NaN(),
	}
}

func NewStopAt(color Color, offset float64) GradientStop {
	return GradientStop{
		Color:  color,
		Offset: offset,
	}
}

func NewStopsAt(color Color, offsets ...float64) []GradientStop {
	stops := make([]GradientStop, len(offsets))
	for i, offset := range offsets {
		stops[i] = GradientStop{Color: color, Offset: offset}
	}
	return stops
}

func (s GradientStop) HasOffset() bool {
	return !math.IsNaN(s.Offset)
}

func (s GradientStop) IsHint() bool {
	return !s.Color.space.IsValid()
}

func (g Gradient) Stops() []GradientStop {
	return slices.Clone(g.stops)
}

func (g Gradient) findStop(t float64) int {
	return sort.Search(len(g.stops), func(i int) bool {
		return g.stops[i].Offset > t
	})
}

func (g Gradient) at(t float64) Color {
	i := g.findStop(t)

	if i == 0 {
		return g.stops[0].Color
	}

	if i == len(g.stops) {
		return g.stops[i-1].Color
	}

	a := &g.stops[i-1]
	b := &g.stops[i]

	seg := (t - a.Offset) * a.invRange

	return g.mixer.UnsafeMix(a.Color, b.Color, seg)
}

func (g Gradient) At(t float64) Color {
	return g.at(t)
}

func (g Gradient) AtRepeat(t float64) Color {
	t = spreadRepeat(t)
	return g.at(t)
}

func (g Gradient) AtReflect(t float64) Color {
	t = spreadReflect(t)
	return g.at(t)
}

func (g Gradient) AtSpread(t float64, spread Spread) Color {
	switch spread {
	case SpreadRepeat:
		return g.AtRepeat(t)
	case SpreadReflect:
		return g.AtReflect(t)
	default:
		return g.At(t)
	}
}

func spreadRepeat(t float64) float64 {
	return t - math.Floor(t)
}

func spreadReflect(t float64) float64 {
	n := math.Floor(t)
	t -= n

	if int64(n)&1 != 0 {
		return 1 - t
	}

	return t
}
