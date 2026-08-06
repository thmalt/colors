package convert

import "math"

const (
	degToRad = math.Pi / 180.0
	radToDeg = 180.0 / math.Pi
)

// Cartesian -> Polar
func LxyToLch(l, x, y float64) (float64, float64, float64) {
	h := math.Atan2(y, x) * radToDeg

	if h < 0 {
		h += 360
	}

	c := math.Hypot(x, y)

	return l, c, h
}

// Polar -> Cartesian
func LchToLxy(l, c, h float64) (float64, float64, float64) {
	rad := h * degToRad
	sin, cos := math.Sincos(rad)

	x := c * cos
	y := c * sin

	return l, x, y
}
