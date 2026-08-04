package colors

import (
	"github.com/thmalt/colors/space"
)

func (c Color) Space() space.Space {
	return c.space
}

func (c Color) Alpha() float64 {
	return c.alpha
}

// WithAlpha returns a copy of [Color] with the specified alpha value.
func (c Color) WithAlpha(alpha float64) Color {
	c.alpha = alpha
	return c
}

func (c Color) ChannelCount() (int, error) {
	info := c.space.Info()
	if info == nil {
		return 0, ErrUnknownSpace
	}

	return info.ChannelCount(), nil
}
