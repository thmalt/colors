package colors

import (
	"github.com/thmalt/colors/interp"
	"github.com/thmalt/colors/space"
)

type HueInterpolation = interp.HueInterpolation

const (
	HueShorter    = interp.HueShorter
	HueLonger     = interp.HueLonger
	HueIncreasing = interp.HueIncreasing
	HueDecreasing = interp.HueDecreasing
)

type InterpOptions struct {
	// Space specifies the color space used for interpolation.
	Space space.Space

	// Premultiplied specifies whether alpha-premultiplied interpolation is used.
	Premultiplied bool

	Hue HueInterpolation
}

// DefaultInterpOptions returns a Mixer configured with the default options.
// The defaults are:
//   - Space: [space.Oklab]
//   - Premultiplied: [true]
//   - Hue: [interp.HueShorter]
func DefaultInterpOptions() InterpOptions {
	return InterpOptions{
		Space:         space.Oklab,
		Premultiplied: true,
		Hue:           HueShorter,
	}
}
