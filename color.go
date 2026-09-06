package colors

import (
	"math"

	"github.com/thmalt/colors/space"
)

const (
	maxUint8  = math.MaxUint8
	maxUint16 = math.MaxUint16

	invMaxUint8  = 1.0 / maxUint8
	invMaxUint16 = 1.0 / maxUint16
)

// New creates a [Color] from a color space and its channel values.
// The number of values must match the color space's channel count,
// optionally followed by an alpha value. Alpha defaults to 1.
//
// It returns an invalid [Color] if the color space or number of values is invalid.
func New(space space.Space, channels ...float64) Color {
	c, _ := TryNew(space, channels...)
	return c
}

// TryNew creates a [Color] from a color space and its channel values.
// The number of values must match the color space's channel count,
// optionally followed by an alpha value. Alpha defaults to 1.
//
// It returns an invalid [Color] and [false] if the color space or number of values is invalid.
func TryNew(space space.Space, channels ...float64) (Color, bool) {
	sc := space.ChannelCount()

	if sc <= 0 || (len(channels) != sc && len(channels) != sc+1) {
		return Color{}, false
	}

	alpha := 1.0
	if len(channels) > sc {
		alpha = channels[sc]
	}

	switch sc {
	case 3:
		_ = channels[2] // Help the compiler eliminate redundant bounds checks.
		return Color{
			space: space,
			c1:    channels[0],
			c2:    channels[1],
			c3:    channels[2],
			alpha: alpha,
		}, true
	case 4:
		_ = channels[3] // Help the compiler eliminate redundant bounds checks.
		return Color{
			space: space,
			c1:    channels[0],
			c2:    channels[1],
			c3:    channels[2],
			c4:    channels[3],
			alpha: alpha,
		}, true
	default:
		return Color{}, false
	}
}

// IsValid reports whether the color's color space is valid.
func (c Color) IsValid() bool {
	return c.space.IsValid()
}

// Space returns the color space of the color.
func (c Color) Space() space.Space {
	return c.space
}

// Alpha returns the alpha channel value without clamping.
func (c Color) Alpha() float64 {
	return c.alpha
}

// Alpha8 returns the alpha channel in the range [0, 255].
func (c Color) Alpha8() uint8 {
	return uint8(clamp01(c.alpha)*maxUint8 + 0.5)
}

// WithAlpha returns a copy of [Color] with the specified alpha value.
func (c Color) WithAlpha(alpha float64) Color {
	c.alpha = alpha
	return c
}

// WithAlpha8 returns a copy of [Color] with the specified alpha value
// in the range [0, 255].
func (c Color) WithAlpha8(alpha uint8) Color {
	c.alpha = float64(alpha) * invMaxUint8
	return c
}

// ChannelCount returns the number of color channels.
func (c Color) ChannelCount() int {
	return c.space.ChannelCount()
}

// CoordinateSystem returns the coordinate system of the color space.
func (c Color) CoordinateSystem() space.CoordinateSystem {
	return c.space.CoordinateSystem()
}
