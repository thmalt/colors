package colors

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/thmalt/colors/mixer"
	"github.com/thmalt/colors/space"
)

type Gradient struct {
	space       space.Space
	channels    uint8
	unsafeMixer mixer.UnsafeMixer
	stops       []gradientStop
}

type GradientStop struct {
	Position float64
	Color    Color
}

func NewGradient(s space.Space, premultiplied bool, hue HueInterpolation, stops ...GradientStop) (Gradient, error) {
	if s == space.InvalidSpace {
		s = space.Oklab
	}

	if !s.IsValid() {
		return Gradient{}, fmt.Errorf("invalid gradient space: %s", s)
	}

	stops = slices.Clone(stops)
	slices.SortFunc(stops, func(a, b GradientStop) int {
		return cmp.Compare(a.Position, b.Position)
	})

	channels := uint8(s.ChannelCount())
	gradientStops := make([]gradientStop, len(stops))
	for i, stop := range stops {
		color, err := stop.Color.To(s)
		if err != nil {
			return Gradient{}, fmt.Errorf("convert stop color at position %g: %w", stop.Position, err)
		}

		gradientStops[i] = newGradientStop(stop.Position, color)
	}

	return Gradient{
		space:       s,
		unsafeMixer: mixer.NewUnsafeMixer(s.HueIndex(), premultiplied, hue),
		stops:       gradientStops,
		channels:    channels,
	}, nil
}
