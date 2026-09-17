package colors

import (
	"github.com/thmalt/colors/mixer"
	"github.com/thmalt/colors/space"
)

type Mixer struct {
	space    space.Space
	channels uint8
	raw      mixer.RawMixer
}

// NewMixerWithOptions creates a [Mixer] that interpolates colors in opts.Space.
// An invalid space defaults to [space.Oklab].
func NewMixerWithOptions(opts InterpOptions) Mixer {
	if !opts.Space.IsValid() {
		opts.Space = space.Oklab
	}

	return Mixer{
		space:    opts.Space,
		channels: uint8(opts.Space.ChannelCount()),
		raw:      mixer.NewRawMixer(opts.Premultiplied, opts.Space.HueIndex(), opts.Hue),
	}
}

// NewMixer returns a new [Mixer] using [DefaultInterpOptions].
func NewMixer() Mixer {
	return NewMixerWithOptions(DefaultInterpOptions())
}

var defaultMixer = NewMixer()

// Space returns the Mixer's color space.
func (m Mixer) Space() space.Space {
	return m.space
}

// Raw returns the [mixer.RawMixer].
func (m Mixer) Raw() mixer.RawMixer {
	return m.raw
}

// Mix interpolates between two colors using the default mix options.
// See [NewMixer] for the default mix options.
func Mix(c1, c2 Color, t float64) Color {
	return defaultMixer.Mix(c1, c2, t)
}

// MixWith interpolates between two colors using the specified interpolation options.
// The interpolation is performed in opts.Space.
// An invalid space defaults to [space.Oklab].
func MixWith(c1, c2 Color, t float64, opts InterpOptions) Color {
	return NewMixerWithOptions(opts).Mix(c1, c2, t)
}
