package colors

import (
	"cmp"
	"slices"

	"github.com/thmalt/colors/mixer"
	"github.com/thmalt/colors/space"
)

type Gradient struct {
	space       space.Space
	unsafeMixer mixer.UnsafeMixer
	stops       []gradientStop
	channels    int
}

type Stop struct {
	Position float64
	Color    Color
}

func NewGradient(s space.Space, premultiplied bool, hue HueInterpolation, stops ...Stop) Gradient {
	if s == space.InvalidSpace {
		s = space.Oklab
	}

	stops = slices.Clone(stops)
	slices.SortFunc(stops, func(a, b Stop) int {
		return cmp.Compare(a.Position, b.Position)
	})

	channels := s.ChannelCount()
	gradientStops := make([]gradientStop, len(stops))
	for i, stop := range stops {
		gradientStops[i] = newGradientStop(stop.Position, stop.Color.MustTo(s))
	}

	return Gradient{
		space:       s,
		unsafeMixer: mixer.NewUnsafeMixer(hueIndex(s), premultiplied, hue),
		stops:       gradientStops,
		channels:    channels,
	}
}
