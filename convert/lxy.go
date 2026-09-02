package convert

import "math"

// LxyToLch converts Lxy Cartesian coordinates to LCh polar coordinates.
// This conversion only changes the coordinate representation; it does not
// involve chromatic adaptation or any other color-space transformation.
func LxyToLch(l, x, y float64) (float64, float64, float64) {
	const radToDeg = 180.0 / math.Pi

	h := math.Atan2(y, x) * radToDeg

	if h < 0 {
		h += 360
	}

	c := math.Hypot(x, y)

	return l, c, h
}

// LchToLxy converts LCh polar coordinates to Lxy Cartesian coordinates.
// This conversion only changes the coordinate representation; it does not
// involve chromatic adaptation or any other color-space transformation.
func LchToLxy(l, c, h float64) (float64, float64, float64) {
	const degToRad = math.Pi / 180.0

	rad := h * degToRad
	sin, cos := math.Sincos(rad)

	x := c * cos
	y := c * sin

	return l, x, y
}
