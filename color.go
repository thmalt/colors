package colors

import (
	"errors"
	"math"

	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/space"
)

var (
	ErrUnknownSpace = errors.New("unknown space")
)

// RGBA implements the [image/color.Color] interface
func (c Color) RGBA() (r, g, b, a uint32) {
	red, green, blue := c.Srgb()

	red = clamp01(red)
	green = clamp01(green)
	blue = clamp01(blue)
	alpha := clamp01(c.alpha)

	const max = 65535.0

	r = uint32(red*alpha*max + 0.5)
	g = uint32(green*alpha*max + 0.5)
	b = uint32(blue*alpha*max + 0.5)
	a = uint32(alpha*max + 0.5)

	return
}

// WithAlpha returns a copy of [Color] with the specified alpha value.
func (c Color) WithAlpha(alpha float64) Color {
	c.alpha = alpha
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

// Rgb returns the color components in the RGB color space.
// Components are in the range [0, 255].
func (c Color) Rgb() (r, g, b float64) {
	return convert.SrgbToRgb(c.Srgb())
}

// Rgb returns a [Color] from 8-bit RGB components in [0, 255].
//
//	r: [0, 255]
//	g: [0, 255]
//	b: [0, 255]
func Rgb(r, g, b float64) Color {
	return Srgb(convert.RgbToSrgb(r, g, b))
}

func (c Color) Hex() string {
	r, g, b := c.Rgb()
	alpha := byte(math.Round(c.alpha * 255.0))

	var out []byte = make([]byte, 9)
	out[0] = '#'

	encodeHexByte(out[1:], byte(math.Round(r)))
	encodeHexByte(out[3:], byte(math.Round(g)))
	encodeHexByte(out[5:], byte(math.Round(b)))

	if alpha == math.MaxUint8 {
		return string(out[:7])
	}

	encodeHexByte(out[7:], alpha)

	return string(out[:])
}
