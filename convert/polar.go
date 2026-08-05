package convert

import "math"

const (
	degToRad = math.Pi / 180.0
	radToDeg = 180.0 / math.Pi
)

func polarToCartesian(l, c, h float64) (float64, float64, float64) {
	rad := h * degToRad
	a := c * math.Cos(rad)
	b := c * math.Sin(rad)
	return l, a, b
}

func cartesianToPolar(l, x, y float64) (float64, float64, float64) {
	h := math.Atan2(y, x) * radToDeg
	if h < 0 {
		h += 360
	}
	c := math.Hypot(x, y)
	return l, c, h
}
