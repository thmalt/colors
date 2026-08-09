package colors

import (
	"github.com/thmalt/colors/interp"
	"github.com/thmalt/colors/space"
)

type MixOptions struct {
	Space space.Space

	Unpremultiplied bool

	Hue HueInterpolation
}

type HueInterpolation = interp.HueInterpolation

const (
	HueShorter    = interp.HueShorter
	HueLonger     = interp.HueLonger
	HueIncreasing = interp.HueIncreasing
	HueDecreasing = interp.HueDecreasing
)

// Mix interpolates two colors using the default mix options.
// It is equivalent to MixWith(c1, c2, t, MixOptions{}).
func Mix(c1, c2 Color, t float64) Color {
	return MixWith(c1, c2, t, MixOptions{})
}
