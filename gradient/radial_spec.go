package gradient

import (
	"math"
)

// RadialShape specifies the shape of a radial gradient.
type RadialShape uint8

const (
	RadialEllipse RadialShape = iota
	RadialCircle
)

// RadialSizeMode specifies how the radius of a radial gradient is determined.
type RadialSizeMode uint8

const (
	RadialFarthestCorner RadialSizeMode = iota
	RadialFarthestSide
	RadialClosestCorner
	RadialClosestSide
	RadialExplicit
)

// RadialSpec specifies the parameters for a radial gradient.
type RadialSpec struct {
	Transform

	width  float64
	height float64

	centerX float64
	centerY float64

	radiusX float64
	radiusY float64

	focalX float64
	focalY float64

	shape    RadialShape
	sizeMode RadialSizeMode

	radiusFixed bool
	hasFocal    bool

	radial *Radial
}

// NewRadialSpec creates a new radial gradient specification.
func NewRadialSpec() RadialSpec {
	return RadialSpec{
		Transform: Identity(),
		centerX:   half,
		centerY:   half,
		radiusX:   1,
		radiusY:   1,
	}
}

// SetSize sets the dimensions of the gradient in pixels.
// Width and height are clamped to a minimum of 1.
func (s *RadialSpec) SetSize(width, height float64) {
	s.width = max(1, width)
	s.height = max(1, height)
}

// SetCenter sets the normalized center point of the gradient.
func (s *RadialSpec) SetCenter(x, y float64) {
	s.centerX = x
	s.centerY = y
}

// SetRadius sets the horizontal and vertical radii of the gradient in normalized units.
func (s *RadialSpec) SetRadius(x, y float64) {
	s.radiusX = x
	s.radiusY = y
	s.radiusFixed = false

	s.shape = RadialEllipse
	s.sizeMode = RadialExplicit
}

// SetRadiusPixel sets the horizontal and vertical radii of the gradient in pixels.
func (s *RadialSpec) SetRadiusPixel(rx, ry float64) {
	s.radiusX = rx
	s.radiusY = ry
	s.radiusFixed = true

	s.shape = RadialEllipse
	s.sizeMode = RadialExplicit
}

// SetCircleRadius sets the radius of a circular gradient in normalized units.
func (s *RadialSpec) SetCircleRadius(radius float64) {
	r := radius * min(s.width, s.height) * half
	s.radiusX = r
	s.radiusY = r
	s.radiusFixed = true

	s.shape = RadialCircle
	s.sizeMode = RadialExplicit
}

// SetCircleRadiusPixel sets the radius of a circular gradient in pixels.
func (s *RadialSpec) SetCircleRadiusPixel(radius float64) {
	s.radiusX = radius
	s.radiusY = radius
	s.radiusFixed = true

	s.shape = RadialCircle
	s.sizeMode = RadialExplicit
}

// SetShape sets the shape of the gradient.
func (s *RadialSpec) SetShape(shape RadialShape) {
	s.shape = shape
}

// SetSizeMode sets the size mode of the gradient.
func (s *RadialSpec) SetSizeMode(size RadialSizeMode) {
	s.sizeMode = size
}

// SetFocal sets the focal point of the gradient.
func (s *RadialSpec) SetFocal(x, y float64) {
	s.focalX = x
	s.focalY = y
	s.hasFocal = true
}

// ClearFocal clears the focal point and restores the default focal point.
func (s *RadialSpec) ClearFocal() {
	s.hasFocal = false
}

// Size returns the dimensions of the gradient in pixels.
func (s *RadialSpec) Size() (width, height float64) {
	return s.width, s.height
}

// Center returns the normalized center point of the gradient.
func (s *RadialSpec) Center() (x, y float64) {
	return s.centerX, s.centerY
}

// Radius returns the horizontal and vertical radii of the gradient in normalized units.
func (s *RadialSpec) Radius() (x, y float64) {
	return s.radiusX, s.radiusY
}

// RadiusPixel returns the horizontal and vertical radii of the gradient in pixels.
func (s *RadialSpec) RadiusPixel() (x, y float64) {
	if s.radiusFixed {
		return s.radiusX, s.radiusY
	}

	return s.radiusX * s.width, s.radiusY * s.height
}

// Shape returns the shape of the gradient.
func (s *RadialSpec) Shape() RadialShape {
	return s.shape
}

// SizeMode returns the size mode of the gradient.
func (s *RadialSpec) SizeMode() RadialSizeMode {
	return s.sizeMode
}

// Focal returns the focal point of the gradient.
func (s *RadialSpec) Focal() (x, y float64, ok bool) {
	return s.focalX, s.focalY, s.hasFocal
}

// Build creates a [Radial] gradient from the specification.
func (s *RadialSpec) Build() (Radial, error) {
	if s.radial == nil {
		s.radial = new(Radial)
	}

	r := s.radial

	if !isFinite(s.width) || !isFinite(s.height) ||
		s.width <= 0 || s.height <= 0 {
		return Radial{}, ErrInvalidSize
	}

	if !isFinite(s.centerX) || !isFinite(s.centerY) {
		return Radial{}, ErrInvalidCenter
	}

	cx := s.centerX * s.width
	cy := s.centerY * s.height

	rx, ry := s.radiusX, s.radiusY

	if s.sizeMode != RadialExplicit {
		rx, ry = resolveRadialRadius(s.sizeMode, s.shape, cx, cy, s.width, s.height)
	} else if !s.radiusFixed {
		rx *= s.width * half
		ry *= s.height * half
	}

	if !isFinite(rx) || !isFinite(ry) || rx <= 0 || ry <= 0 {
		return Radial{}, ErrInvalidRadius
	}

	invRX := 1 / rx
	invRY := 1 / ry

	inv, ok := s.Transform.Inverse()
	if !ok {
		return Radial{}, ErrInvalidTransform
	}

	// dx = (inverseTransform(x,y) - center) / radius
	r.dx0 = inv.m00 * invRX
	r.dx1 = inv.m01 * invRX
	r.dx2 = (inv.tx - cx) * invRX

	r.dy0 = inv.m10 * invRY
	r.dy1 = inv.m11 * invRY
	r.dy2 = (inv.ty - cy) * invRY

	r.focal = false
	r.focalDX = 0
	r.focalDY = 0
	r.focalC = 0

	if s.hasFocal {
		if !isFinite(s.focalX) || !isFinite(s.focalY) {
			return Radial{}, ErrInvalidFocal
		}

		fx := s.focalX * s.width
		fy := s.focalY * s.height

		// Focal relative to center, normalized by radius.
		fdx := (fx - cx) * invRX
		fdy := (fy - cy) * invRY

		if d := math.Hypot(fdx, fdy); d > focalEpsilon {
			if d > 1 {
				scale := 1 / d
				fdx *= scale
				fdy *= scale
			}

			r.focalDX = fdx
			r.focalDY = fdy
			r.focalC = fdx*fdx + fdy*fdy - 1

			// Change the runtime mapping from center-relative
			// to focal-relative coordinates.
			r.dx2 -= fdx
			r.dy2 -= fdy

			r.focal = true
		}
	}

	return *r, nil
}

const focalEpsilon = 1e-12

func resolveRadialRadius(size RadialSizeMode, shape RadialShape, cx, cy float64, width, height float64) (rx, ry float64) {
	switch shape {
	case RadialCircle:
		radius := circleRadius(size, cx, cy, width, height)
		rx, ry = radius, radius
	case RadialEllipse:
		rx, ry = ellipseRadii(size, cx, cy, width, height)
	}
	return
}

func circleRadius(size RadialSizeMode, cx, cy float64, width, height float64) (radius float64) {
	switch size {
	case RadialFarthestCorner:
		dx := max(cx, width-cx)
		dy := max(cy, height-cy)
		radius = math.Hypot(dx, dy)
	case RadialClosestCorner:
		dx := min(cx, width-cx)
		dy := min(cy, height-cy)
		radius = math.Hypot(dx, dy)
	case RadialFarthestSide:
		radius = max(cx, width-cx, cy, height-cy)
	case RadialClosestSide:
		radius = min(cx, width-cx, cy, height-cy)
	}
	return
}

func ellipseRadii(size RadialSizeMode, cx, cy float64, width, height float64) (rx, ry float64) {
	left := cx
	right := width - cx
	top := cy
	bottom := height - cy

	switch size {
	case RadialFarthestCorner:
		aspect := width / height

		dx := max(cx, width-cx)
		dy := max(cy, height-cy)

		ry = math.Hypot(dx/aspect, dy)
		rx = ry * aspect
	case RadialClosestCorner:
		aspect := width / height

		dx := min(cx, width-cx)
		dy := min(cy, height-cy)

		ry = math.Hypot(dx/aspect, dy)
		rx = ry * aspect
	case RadialFarthestSide:
		rx, ry = max(left, right), max(top, bottom)
	case RadialClosestSide:
		rx, ry = min(left, right), min(top, bottom)
	}

	return
}
