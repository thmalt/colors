package gradient

// Linear represents a linear gradient.
type Linear struct {
	ax   float64
	ay   float64
	bias float64
}

// PositionAt returns the normalized gradient position at the specified point.
func (l *Linear) PositionAt(x, y float64) float64 {
	return l.ax*x + l.ay*y + l.bias
}
