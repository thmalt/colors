package convert

func LabD50ToXyzD50(l, a, b float64) (x, y, z float64) {
	fy := (l + 16) * labInv116
	fx := fy + a*labInv500
	fz := fy - b*labInv200

	x = d50X * labInvF(fx)
	y = d50Y * labInvF(fy)
	z = d50Z * labInvF(fz)

	return
}

func XyzD50ToLabD50(x, y, z float64) (l, a, b float64) {
	fx := labF(x * d50InvX)
	fy := labF(y * d50InvY)
	fz := labF(z * d50InvZ)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)

	return
}
