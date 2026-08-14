package colors

import (
	"math"
	"slices"
	"sort"
)

type Gradient struct {
	mixer Mixer
	stops []GradientStop
}

type GradientStop struct {
	Color  Color
	Offset float64
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

func (g Gradient) At(t float64) Color {
	i := g.findStop(t)

	if i == 0 {
		return g.stops[0].Color
	}

	if i == len(g.stops) {
		return g.stops[len(g.stops)-1].Color
	}

	a := g.stops[i-1]
	b := g.stops[i]

	seg := (t - a.Offset) / (b.Offset - a.Offset)

	return g.mixer.Mix(a.Color, b.Color, seg)
}

func (g Gradient) findStop(t float64) int {
	return sort.Search(len(g.stops), func(i int) bool {
		return g.stops[i].Offset > t // >= t
	})
}

// func (g Gradient) findStop(t float64) int {
// 	lo, hi := 0, len(g.stops)

// 	for lo < hi {
// 		mid := lo + (hi-lo)/2

// 		if g.stops[mid].Offset <= t {
// 			lo = mid + 1
// 		} else {
// 			hi = mid
// 		}
// 	}

// 	return lo
// }

func (g Gradient) Stops() []GradientStop {
	return slices.Clone(g.stops)
}
