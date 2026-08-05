package convert

func LuvToLchuv(l, u, v float64) (float64, float64, float64) {
	return cartesianToPolar(l, u, v)
}

func LchuvToLuv(l, c, h float64) (float64, float64, float64) {
	return polarToCartesian(l, c, h)
}
