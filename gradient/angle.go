package gradient

func DegToTurn(deg float64) float64 {
	return deg / 360
}

func RadToTurn(rad float64) float64 {
	return rad * invTau
}

func GradToTurn(grad float64) float64 {
	return grad / 400
}

func TurnToDeg(turn float64) float64 {
	return turn * 360
}

func TurnToRad(turn float64) float64 {
	return turn * tau
}

func TurnToGrad(turn float64) float64 {
	return turn * 400
}

const (
	Turn        = 1.0      // 360°	1   turn
	HalfTurn    = Turn / 2 // 180°  1/2 turn
	QuarterTurn = Turn / 4 // 90°   1/4 turn
	EighthTurn  = Turn / 8 // 45°   1/8 turn

	ToTop         = 0.0            // 0°   0   turn
	ToTopRight    = EighthTurn     // 45°  1/8 turn
	ToRight       = QuarterTurn    // 90°  1/4 turn
	ToBottomRight = EighthTurn * 3 // 135° 3/8 turn
	ToBottom      = HalfTurn       // 180° 1/2 turn
	ToBottomLeft  = EighthTurn * 5 // 225° 5/8 turn
	ToLeft        = EighthTurn * 6 // 270° 3/4 turn
	ToTopLeft     = EighthTurn * 7 // 315° 7/8 turn
)
