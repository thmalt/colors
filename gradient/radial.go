package gradient

import "math"

type RadialShape uint8

const (
	RadialEllipse RadialShape = iota
	RadialCircle
)

type RadialSize uint8

const (
	RadialFarthestCorner RadialSize = iota
	RadialFarthestSide
	RadialClosestCorner
	RadialClosestSide
	RadialExplicit
)

// Radial defines the geometry of a radial gradient.
type Radial struct {
	width  float64
	height float64

	centerX float64
	centerY float64

	shape RadialShape
	size  RadialSize

	rx, ry float64 // [0,1]

	// Cached transform values.
	cx, cy float64
	invRX  float64
	invRY  float64
}

// NewRadial creates a radial gradient.
//
//   - width, height: gradient bounds in pixels.
//   - centerX, centerY: normalized center coordinates in the range [0, 1].
//   - The shape defaults to [RadialEllipse].
//   - The size defaults to [RadialFarthestCorner].
func NewRadial(width, height float64, centerX, centerY float64) Radial {
	r := Radial{
		centerX: centerX,
		centerY: centerY,
		shape:   RadialEllipse,
		size:    RadialFarthestCorner,
	}

	r.SetBounds(width, height)
	return r
}

// PositionAt returns the normalized gradient position at (x, y).
func (r Radial) PositionAt(x, y float64) float64 {
	dx := (x - r.cx) * r.invRX
	dy := (y - r.cy) * r.invRY

	return math.Sqrt(dx*dx + dy*dy) // cost 36
}

// Bounds returns the width and height of the gradient bounds in pixels.
func (r *Radial) Bounds() (width, height float64) {
	return r.width, r.height
}

// SetBounds sets the gradient bounds in pixels.
// Width and height are clamped to a minimum of 1.
func (r *Radial) SetBounds(width, height float64) {
	r.width = max(1, width)
	r.height = max(1, height)

	r.updateTransform()
}

// Center returns the normalized center coordinates in the range [0, 1].
func (r *Radial) Center() (x, y float64) {
	return r.width, r.height
}

// SetCenter sets the normalized center coordinates in the range [0, 1].
func (r *Radial) SetCenter(x, y float64) {
	r.centerX = x
	r.centerY = y

	r.updateTransform()
}

// Radii returns the horizontal and vertical radii.
func (r Radial) Radii() (rx, ry float64) {
	return r.rx, r.ry
}

// SetRadius sets the radius as a fraction of half the minimum bounds dimension
// and changes the radial gradient shape to a circle.
func (r *Radial) SetRadius(radius float64) {
	r.shape = RadialCircle
	r.size = RadialExplicit

	radius *= (min(r.width, r.height)) / 2
	r.rx = radius
	r.ry = radius

	r.updateTransform()
}

// SetRadii sets the horizontal and vertical radii as fractions of half the
// corresponding bounds dimensions and changes the radial gradient shape to an ellipse.
func (r *Radial) SetRadii(rx, ry float64) {
	r.shape = RadialEllipse
	r.size = RadialExplicit

	r.rx = rx * r.width / 2
	r.ry = ry * r.height / 2

	r.updateTransform()
}

// Shape returns the radial gradient shape.
func (r *Radial) Shape() RadialShape {
	return r.shape
}

// SetShape sets the radial gradient shape.
func (r *Radial) SetShape(shape RadialShape) {
	r.shape = shape

	r.updateTransform()
}

// Size returns the radial gradient size.
func (r *Radial) Size() RadialSize {
	return r.size
}

// SetSize sets the radial gradient size.
func (r *Radial) SetSize(size RadialSize) {
	r.size = size

	r.updateTransform()
}

// updateTransform updates the cached pixel-space center from the normalized center
// and gradient bounds.
func (r *Radial) updateTransform() {
	r.cx = r.centerX * r.width
	r.cy = r.centerY * r.height

	if r.size != RadialExplicit {
		switch r.shape {
		case RadialCircle:
			radius := r.circleRadius()
			r.rx, r.ry = radius, radius
		case RadialEllipse:
			r.rx, r.ry = r.ellipseRadii()
		}
	}

	r.invRX = 1 / r.rx
	r.invRY = 1 / r.ry
}

func (r *Radial) circleRadius() (radius float64) {
	switch r.size {
	case RadialFarthestCorner:
		dx := max(r.cx, r.width-r.cx)
		dy := max(r.cy, r.height-r.cy)
		radius = math.Hypot(dx, dy)
	case RadialClosestCorner:
		dx := min(r.cx, r.width-r.cx)
		dy := min(r.cy, r.height-r.cy)
		radius = math.Hypot(dx, dy)
	case RadialFarthestSide:
		radius = max(r.cx, r.width-r.cx, r.cy, r.height-r.cy)
	case RadialClosestSide:
		radius = min(r.cx, r.width-r.cx, r.cy, r.height-r.cy)
	case RadialExplicit:
		radius = math.NaN()
	default:
		dx := max(r.cx, r.width-r.cx)
		dy := max(r.cy, r.height-r.cy)
		radius = math.Hypot(dx, dy)
	}

	return
}

func (r *Radial) ellipseRadii() (rx, ry float64) {
	left := r.cx
	right := r.width - r.cx
	top := r.cy
	bottom := r.height - r.cy

	switch r.size {
	case RadialFarthestCorner:
		aspect := r.width / r.height

		dx := max(r.cx, r.width-r.cx)
		dy := max(r.cy, r.height-r.cy)

		ry = math.Hypot(dx/aspect, dy)
		rx = ry * aspect
	case RadialClosestCorner:
		aspect := r.width / r.height

		dx := min(r.cx, r.width-r.cx)
		dy := min(r.cy, r.height-r.cy)

		ry = math.Hypot(dx/aspect, dy)
		rx = ry * aspect
	case RadialFarthestSide:
		rx, ry = max(left, right), max(top, bottom)
	case RadialClosestSide:
		rx, ry = min(left, right), min(top, bottom)
	case RadialExplicit:
		rx, ry = math.NaN(), math.NaN()
	default:
		aspect := r.width / r.height

		dx := max(r.cx, r.width-r.cx)
		dy := max(r.cy, r.height-r.cy)

		ry = math.Hypot(dx/aspect, dy)
		rx = ry * aspect
	}

	return
}
