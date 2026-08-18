package gradient

import "math"

// Radial represents a radial gradient.
type Radial struct {
	dx0, dx1, dx2 float64
	dy0, dy1, dy2 float64

	focalDX, focalDY float64
	focalC           float64

	focal bool
}

// PositionAt returns the normalized gradient position at the specified point.
func (r *Radial) PositionAt(x, y float64) float64 {
	dx := r.dx0*x + r.dx1*y + r.dx2
	dy := r.dy0*x + r.dy1*y + r.dy2

	if !r.focal {
		return math.Hypot(dx, dy)
	}

	a := dx*dx + dy*dy
	if a == 0 {
		return 0
	}

	b := dx*r.focalDX + dy*r.focalDY

	return a / (-b + math.Sqrt(b*b-a*r.focalC))
}

// ScaleRadius scales the radius uniformly.
func (r *Radial) ScaleRadius(s float64) {
	r.ScaleRadiusXY(s, s)
}

// ScaleRadiusXY scales the horizontal and vertical radii independently.
func (r *Radial) ScaleRadiusXY(sx, sy float64) {
	invX := 1 / sx
	invY := 1 / sy

	r.dx0 *= invX
	r.dx1 *= invX
	r.dx2 *= invX

	r.dy0 *= invY
	r.dy1 *= invY
	r.dy2 *= invY
}
