package colors

import (
	"github.com/thmalt/colors/space"
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

// CoordinateSystem returns the coordinate system of the color space.
func (c Color) CoordinateSystem() space.CoordinateSystem {
	return c.space.CoordinateSystem()
}

// ChannelCount returns the number of color channels.
func (c Color) ChannelCount() int {
	return c.space.ChannelCount()
}

// Alpha returns the alpha value.
func (c Color) Alpha() float64 {
	return c.alpha
}

// Alpha8 returns the alpha value as an 8-bit value in the range [0, 255].
func (c Color) Alpha8() uint8 {
	return uint8(clamp01(c.alpha)*maxUint8 + 0.5)
}

// WithAlpha returns a copy of [Color] with the specified alpha value.
// Alpha values in the range [0, 1] are typical.
func (c Color) WithAlpha(alpha float64) Color {
	c.alpha = alpha
	return c
}

// WithAlpha8 returns a copy of [Color] with the specified 8-bit alpha value
// in the range [0, 255].
func (c Color) WithAlpha8(alpha uint8) Color {
	c.alpha = float64(alpha) * invMaxUint8
	return c
}

// IsOpaque reports whether the color has an alpha value of at least 1.
func (c Color) IsOpaque() bool {
	return c.alpha >= 1
}

// Opaque returns a copy of [Color] with its alpha set to 1.
func (c Color) Opaque() Color {
	c.alpha = 1
	return c
}

// Rgb returns an sRGB [Color] from RGB components.
//
//	r, g, b: [0, 255]
func Rgb(r, g, b float64) Color {
	return Srgb(
		r*invMaxUint8,
		g*invMaxUint8,
		b*invMaxUint8,
	)
}

// Rgb returns an sRGB [Color] from RGB components and alpha.
//
//	r, g, b: [0, 255]
//	a:       [0, 1]
func Rgba(r, g, b, a float64) Color {
	return SrgbAlpha(
		r*invMaxUint8,
		g*invMaxUint8,
		b*invMaxUint8,
		a,
	)
}

// Rgb returns the sRGB components scaled by 255.
func (c Color) Rgb() (r, g, b float64) {
	r, g, b = c.srgb()
	return r * maxUint8, g * maxUint8, b * maxUint8
}

// Rgba returns the sRGB components scaled by 255 and alpha in [0, 1].
func (c Color) Rgba() (r, g, b, a float64) {
	r, g, b = c.srgb()
	return r * maxUint8, g * maxUint8, b * maxUint8, c.alpha
}
