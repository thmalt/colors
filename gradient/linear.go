package gradient

import "math"

// Linear defines the geometry of a linear gradient.
type Linear struct {
	width  float64
	height float64
	angle  float64

	// Cached transform values.
	dx, dy float64
	start  float64
	bias   float64
}

// NewLinear creates a linear gradient.
//
//   - width, height: gradient bounds in pixels.
//   - angle: gradient angle in turns.
func NewLinear(width, height, angle float64) Linear {
	l := Linear{
		angle: angle,
	}

	l.SetBounds(width, height)
	return l
}

// PositionAt returns the normalized gradient position at (x, y).
func (l *Linear) PositionAt(x, y float64) float64 {
	return l.dx*x + l.dy*y + l.bias
}

func (l *Linear) updateTransform() {
	rad := l.angle * tau
	sin, cos := math.Sincos(rad)

	dx, dy := sin, -cos
	hw, hh := l.width/2, l.height/2

	L := (math.Abs(dx)*hw + math.Abs(dy)*hh) * 2

	l.dx = dx / L
	l.dy = dy / L
	l.bias = 0.5 - hw*l.dx - hh*l.dy
}

// Bounds returns the width and height of the gradient bounds in pixels.
func (l *Linear) Bounds() (width, height float64) {
	return l.width, l.height
}

// SetBounds sets the gradient bounds in pixels.
// Width and height are clamped to a minimum of 1.
func (l *Linear) SetBounds(width, height float64) {
	l.width = max(1, width)
	l.height = max(1, height)

	l.updateTransform()
}

// Angle returns the gradient angle in turns.
func (l *Linear) Angle() float32 {
	return l.Angle()
}

// SetAngle sets the gradient angle in turns.
func (l *Linear) SetAngle(angle float64) {
	l.angle = angle

	l.updateTransform()
}
