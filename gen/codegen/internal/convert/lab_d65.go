package convert

func LabD65ToXyzD65(l, a, b float64) (x, y, z float64) {
	fy := (l + 16) * labInv116
	fx := fy + a*labInv500
	fz := fy - b*labInv200

	x = d65X * labInvF(fx)
	y = d65Y * labInvF(fy)
	z = d65Z * labInvF(fz)

	return
}

func XyzD65ToLabD65(x, y, z float64) (l, a, b float64) {
	fx := labF(x * d65InvX)
	fy := labF(y * d65InvY)
	fz := labF(z * d65InvZ)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)

	return
}
