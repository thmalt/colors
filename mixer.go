package colors

import (
	"github.com/thmalt/colors/mixer"
	"github.com/thmalt/colors/space"
)

type Mixer struct {
	space    space.Space
	channels int
	unsafe   mixer.UnsafeMixer
}

func (m Mixer) Mixx(c1, c2 Color, t float64) Color {
	c1 = c1.MustTo(m.space)
	c2 = c2.MustTo(m.space)

	switch m.channels {
	case 3:
		a1, a2, a3, aa := c1.c1, c1.c2, c1.c3, c1.alpha
		b1, b2, b3, ba := c2.c1, c2.c2, c2.c3, c2.alpha
		c1, c2, c3, alpha := m.unsafe.Mix3(a1, a2, a3, aa, b1, b2, b3, ba, t)
		return Color{space: m.space, c1: c1, c2: c2, c3: c3, alpha: alpha}
	}

	return Color{}
}

func NewMixer(opts MixOptions) Mixer {
	if opts.Space == space.InvalidSpace {
		opts.Space = space.Oklab
	}

	return Mixer{
		space:    opts.Space,
		channels: opts.Space.ChannelCount(),
		unsafe:   mixer.NewUnsafeMixer(hueIndex(opts.Space), !opts.Unpremultiplied, opts.Hue),
	}
}
