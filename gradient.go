package colors

import (
	"fmt"
	"sort"

	"github.com/thmalt/colors/mixer"
	"github.com/thmalt/colors/space"
)

// Gradient interpolates colors between a sequence of stops.
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

	channels := uint8(s.ChannelCount())
	gradientStops := make([]gradientStop, len(stops))
	for i, stop := range stops {
		color, err := stop.Color.To(s)
		if err != nil {
			return Gradient{}, fmt.Errorf(
				"convert stop color at position %g: %w",
				stop.Position, err,
			)
		}

		position := stop.Position
		if i > 0 && position < gradientStops[i-1].position {
			position = gradientStops[i-1].position
		}

		gradientStops[i] = newGradientStop(stop.Position, position, color)
	}

	return Gradient{
		space:       s,
		unsafeMixer: mixer.NewUnsafeMixer(s.HueIndex(), premultiplied, hue),
		stops:       gradientStops,
		channels:    channels,
	}, nil
}

// NewStop returns a new [GradientStop].
func NewStop(position float64, color Color) GradientStop {
	return GradientStop{
		Position: position,
		Color:    color,
	}
}

func (g Gradient) findStop(t float64) int {
	return sort.Search(len(g.stops), func(i int) bool {
		return g.stops[i].position >= t
	})
}
