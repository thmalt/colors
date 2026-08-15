package gradient

import "math"

type Conic struct {
	dx0, dx1, dx2 float64
	dy0, dy1, dy2 float64

	startAngle float64
}

func (c Conic) PositionAt(x, y float64) float64 {
	dx := c.dx0*x + c.dx1*y + c.dx2
	dy := c.dy0*x + c.dy1*y + c.dy2

	angle := math.Atan2(dx, -dy) * invTau
	t := angle - c.startAngle

	return t - math.Floor(t)
}
