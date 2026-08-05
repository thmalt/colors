package convert

func OklabToOklch(l, a, b float64) (float64, float64, float64) {
	return cartesianToPolar(l, a, b)
}

func OklchToOklab(l, c, h float64) (float64, float64, float64) {
	return polarToCartesian(l, c, h)
}
