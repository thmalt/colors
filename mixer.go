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

func NewMixer(opts MixOptions) Mixer {
	if opts.Space == space.InvalidSpace {
		opts.Space = space.Oklab
	}

	return Mixer{
		space:    opts.Space,
		channels: uint8(opts.Space.ChannelCount()),
		unsafe:   mixer.NewUnsafeMixer(opts.Space.HueIndex(), !opts.Unpremultiplied, opts.Hue),
	}
}
