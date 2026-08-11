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
