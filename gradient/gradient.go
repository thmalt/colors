package gradient

import "math"

const (
	half   = 0.5
	tau    = 2 * math.Pi
	invTau = 1 / tau
)

func isFinite(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0)
}
