package colors

import (
	"github.com/thmalt/colors/space"
)

// Adjuster provides chained adjustments for a color in [space.Oklch].
type Adjuster struct {
	color    Color
	original space.Space
}

// NewAdjuster returns an [Adjuster] initialized with the specified color.
// The color is converted to [space.Oklch] for adjustment.
func NewAdjuster(c Color) *Adjuster {
	original := c.space
	c.mutTo(space.Oklch)
	return &Adjuster{
		color:    c,
		original: original,
	}
}

// Adjuster returns an [Adjuster] for the color.
func (c Color) Adjuster() *Adjuster {
	return NewAdjuster(c)
}

// Lightness returns the current lightness in [space.Oklch].
func (a *Adjuster) Lightness() float64 {
	return a.color.c1
}

// Chroma returns the current chroma in [space.Oklch].
func (a *Adjuster) Chroma() float64 {
	return a.color.c2
}

// Hue returns the current hue in [space.Oklch] in degrees.
func (a *Adjuster) Hue() float64 {
	return a.color.c3
}

// SetLightness sets the lightness in [space.Oklch] to v.
// The value is clamped to [0, 1].
func (a *Adjuster) SetLightness(v float64) *Adjuster {
	a.color.c1 = clamp01(v)
	return a
}

// SetChroma sets the chroma in [space.Oklch] to v.
// Negative values are clamped to 0.
func (a *Adjuster) SetChroma(v float64) *Adjuster {
	a.color.c2 = max(0, v)
	return a
}

// SetHue sets the hue in [space.Oklch] to v degrees.
// The value is normalized to [0, 360).
func (a *Adjuster) SetHue(v float64) *Adjuster {
	a.color.c3 = wrap360(v)
	return a
}

// Color returns the current color in [space.Oklch].
func (a *Adjuster) Color() Color {
	return a.color
}

// Result returns the adjusted color in its original color space.
func (a *Adjuster) Result() Color {
	c := a.color
	c.mutTo(a.original)
	return c
}

// Darken decreases the lightness in [space.Oklch] by amount.
func (a *Adjuster) Darken(amount float64) *Adjuster {
	a.color.c1 = clamp01(a.color.c1 - amount)
	return a
}

// Lighten increases the lightness in [space.Oklch] by amount.
func (a *Adjuster) Lighten(amount float64) *Adjuster {
	a.color.c1 = clamp01(a.color.c1 + amount)
	return a
}

// Saturate increases the chroma in [space.Oklch] by amount.
func (a *Adjuster) Saturate(amount float64) *Adjuster {
	a.color.c2 = max(0, a.color.c2+amount)
	return a
}

// Desaturate decreases the chroma in [space.Oklch] by amount.
func (a *Adjuster) Desaturate(amount float64) *Adjuster {
	a.color.c2 = max(0, a.color.c2-amount)
	return a
}

// AdjustHue adjusts the hue in [space.Oklch] by amount.
func (a *Adjuster) AdjustHue(amount float64) *Adjuster {
	a.color.c3 = wrap360(a.color.c3 + amount)
	return a
}

// Grayscale sets the chroma in [space.Oklch] to 0.
func (a *Adjuster) Grayscale() *Adjuster {
	a.color.c2 = 0
	return a
}

// Darken returns a copy of the color with its lightness in [space.Oklch] decreased by amount.
//
// The resulting lightness is clamped to [0, 1].
func (c Color) Darken(amount float64) Color {
	if s := c.space; s != space.Oklch {
		c.mutTo(space.Oklch)
		c.c1 = clamp01(c.c1 - amount)
		c.mutTo(s)
		return c
	}

	c.c1 = clamp01(c.c1 - amount)
	return c
}

// Lighten returns a copy of the color with its lightness in [space.Oklch] increased by amount.
//
// The resulting lightness is clamped to [0, 1].
func (c Color) Lighten(amount float64) Color {
	if s := c.space; s != space.Oklch {
		c.mutTo(space.Oklch)
		c.c1 = clamp01(c.c1 + amount)
		c.mutTo(s)
		return c
	}

	c.c1 = clamp01(c.c1 + amount)
	return c
}

// Saturate returns a copy of the color with its chroma in [space.Oklch] increased by amount.
//
// The resulting chroma is clamped to [0, +Inf).
func (c Color) Saturate(amount float64) Color {
	if s := c.space; s != space.Oklch {
		c.mutTo(space.Oklch)
		c.c2 = max(0, c.c2+amount)
		c.mutTo(s)
		return c
	}

	c.c2 = max(0, c.c2+amount)
	return c
}

// Desaturate returns a copy of the color with its chroma in [space.Oklch] decreased by amount.
//
// The resulting chroma is clamped to [0, +Inf).
func (c Color) Desaturate(amount float64) Color {
	if s := c.space; s != space.Oklch {
		c.mutTo(space.Oklch)
		c.c2 = max(0, c.c2-amount)
		c.mutTo(s)
		return c
	}

	c.c2 = max(0, c.c2-amount)
	return c
}

// AdjustHue returns a copy of the color with its hue in [space.Oklch] adjusted by amount.
//
// The resulting hue is normalized to [0, 360).
func (c Color) AdjustHue(amount float64) Color {
	if s := c.space; s != space.Oklch {
		c.mutTo(space.Oklch)
		c.c3 = wrap360(c.c3 + amount)
		c.mutTo(s)
		return c
	}

	c.c3 = wrap360(c.c3 + amount)
	return c
}

// Grayscale returns a copy of the color with its chroma in [space.Oklch] set to 0.
func (c Color) Grayscale() Color {
	if s := c.space; s != space.Oklch {
		c.mutTo(space.Oklch)
		c.c2 = 0
		c.mutTo(s)
		return c
	}

	c.c2 = 0
	return c
}
