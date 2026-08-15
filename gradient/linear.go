package gradient

type Linear struct {
	ax   float64
	ay   float64
	bias float64
}

func (l *Linear) PositionAt(x, y float64) float64 {
	return l.ax*x + l.ay*y + l.bias
}
