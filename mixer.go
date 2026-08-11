package colors

import (
	"github.com/thmalt/colors/mixer"
	"github.com/thmalt/colors/space"
)

type Mixer struct {
	space    space.Space
	channels uint8
	unsafe   mixer.UnsafeMixer
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
		unsafe:   mixer.NewUnsafeMixer(opts.Space.HueIndex(), opts.Premultiplied, opts.Hue),
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

// Unsafe returns the mixer without safety checks.
func (m Mixer) Unsafe() mixer.UnsafeMixer {
	return m.unsafe
}

// Mix linearly interpolates c1 and c2 using the default mixer.
// See [NewMixer] for the default mix options.
func Mix(c1, c2 Color, t float64) Color {
	return defaultMixer.Mix(c1, c2, t)
}

// MixWith returns the [Color] interpolation of c1 and c2.
//
// The interpolation is performed in opts.Space. If opts.Space is
// [space.InvalidSpace], [space.Oklab] is used.
//
// The interpolation behavior can be customized through opts, including
// premultiplied alpha and hue interpolation for polar color spaces.
func MixWith(c1, c2 Color, t float64, opts InterpOptions) Color {
	return NewMixerWithOptions(opts).Mix(c1, c2, t)
}
