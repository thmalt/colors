package colors

import (
	"errors"

	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/space"
)

var (
	ErrUnknownSpace = errors.New("unknown space")
)

// Implement [image/color.Color] interface
func (c Color) RGBA() (r, g, b, a uint32) {
	red, green, blue := c.Srgb()

	alpha := c.alpha

	a = uint32(alpha*65535 + 0.5)

	r = uint32(red*65535*alpha + 0.5)
	g = uint32(green*65535*alpha + 0.5)
	b = uint32(blue*65535*alpha + 0.5)

	return
}

func (c Color) WithAlpha(alpha float64) Color {
	c.alpha = clamp01(alpha)
	return c
}

func (c Color) Alpha() float64 {
	return c.alpha
}

func (c Color) Space() space.Space {
	return c.space
}

func (c Color) ChannelCount() (int, error) {
	info := c.space.Info()
	if info == nil {
		return 0, ErrUnknownSpace
	}

	return info.ChannelCount(), nil
}

func Rgb(r, g, b float64) Color {
	return Srgb(convert.RgbToSrgb(r, g, b))
}
