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
