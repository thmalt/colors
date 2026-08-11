package gradient

import "math"

// Conic defines the geometry of a conic gradient.
type Conic struct {
	width  float64
	height float64

	centerX float64
	centerY float64

	startAngle float64

	// Cached transform values.
	cx, cy float64
}

// NewConic creates a conic gradient.
//
// Parameters:
//   - width, height: gradient bounds in pixels.
//   - centerX, centerY: normalized center coordinates in [0, 1].
//   - startAngle: starting angle in turns, normalized to [0, 1).
func NewConic(width, height float64, centerX, centerY, startAngle float64) Conic {
	c := Conic{
		centerX: centerX,
		centerY: centerY,
	}

	c.SetStartAngle(startAngle)
	c.SetBounds(width, height)
	return c
}

// PositionAt returns the normalized gradient position at (x, y).
func (c Conic) PositionAt(x, y float64) float64 {
	dx := x - c.cx
	dy := y - c.cy

	angle := math.Atan2(dx, -dy) * invTau
	t := angle - c.startAngle

	return t - math.Floor(t)
}

// updateTransform updates the cached pixel-space center from the normalized center
// and gradient bounds.
func (c *Conic) updateTransform() {
	c.cx = c.centerX * c.width
	c.cy = c.centerY * c.height
}

// Bounds returns the width and height of the gradient bounds in pixels.
func (c *Conic) Bounds() (width, height float64) {
	return c.width, c.height
}

// SetBounds sets the gradient bounds in pixels.
// Width and height are clamped to a minimum of 1.
func (c *Conic) SetBounds(width, height float64) {
	c.width = max(1, width)
	c.height = max(1, height)

	c.updateTransform()
}

// Center returns the normalized center coordinates in the range [0, 1].
func (c *Conic) Center() (x, y float64) {
	return c.width, c.height
}

// SetCenter sets the normalized center coordinates in the range [0, 1].
func (c *Conic) SetCenter(x, y float64) {
	c.centerX = x
	c.centerY = y

	c.updateTransform()
}

// StartAngle returns the starting angle in turns.
func (c *Conic) StartAngle() float64 {
	return c.startAngle
}

// SetStartAngle sets the starting angle in turns.
func (c *Conic) SetStartAngle(angle float64) {
	c.startAngle = angle - math.Floor(angle)
}
