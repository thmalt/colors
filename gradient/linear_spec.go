package gradient

import (
	"math"
)

type LinearSpec struct {
	Transform

	width  float64
	height float64

	angle     float64
	direction bool

	linear *Linear
}

func NewLinearSpec() LinearSpec {
	return LinearSpec{
		Transform: Identity(),
	}
}

// SetSize sets the gradient bounds in pixels.
// Width and height are clamped to a minimum of 1.
func (s *LinearSpec) SetSize(width, height float64) {
	s.width = max(1, width)
	s.height = max(1, height)
}

// SetAngle sets the gradient angle in turns.
func (s *LinearSpec) SetAngle(turn float64) {
	s.angle = turn
	s.direction = false
}

// SetDirection sets the gradient direction relative to the gradient size.
// The direction is resolved during Build.
func (s *LinearSpec) SetDirection(angle float64) {
	s.angle = angle
	s.direction = true
}

func (s *LinearSpec) Size() (width, height float64) {
	return s.width, s.height
}

func (s *LinearSpec) Angle() (turn float64) {
	return s.angle
}

func (s *LinearSpec) Build() (Linear, error) {
	if s.linear == nil {
		s.linear = new(Linear)
	}

	l := s.linear

	if !isFinite(s.width) || !isFinite(s.height) ||
		s.width <= 0 || s.height <= 0 {
		return Linear{}, ErrInvalidSize
	}

	angle := s.angle - QuarterTurn
	if s.direction {
		angle = resolveDirectionAngle(angle, s.width, s.height)
	}

	transform := s.Transform
	transform.Rotate(angle * tau)

	inv, ok := transform.Inverse()
	if !ok {
		return Linear{}, ErrInvalidTransform
	}

	x0 := 0.0
	x1 := inv.m00 * s.width
	x2 := inv.m01 * s.height
	x3 := x1 + x2

	minX := min(x0, x1, x2, x3)
	maxX := max(x0, x1, x2, x3)

	d := maxX - minX
	if !isFinite(d) || d <= 0 {
		return Linear{}, ErrInvalidTransform
	}

	invD := 1 / d

	l.ax = inv.m00 * invD
	l.ay = inv.m01 * invD
	l.bias = (inv.tx - minX) * invD

	return *l, nil
}

func resolveDirectionAngle(angle, width, height float64) float64 {
	a := math.Mod(angle, 1)
	if a < 0 {
		a++
	}

	d := math.Atan2(width, height) / tau
	switch a {
	case ToTop, ToRight, ToBottom, ToLeft:
		return a
	case ToTopRight:
		return d
	case ToBottomRight:
		return 0.5 - d
	case ToBottomLeft:
		return 0.5 + d
	case ToTopLeft:
		return 1 - d
	}

	return a
}
