package convert

import (
	"math"
)

const (
	labDelta  = 6.0 / 29
	labDelta2 = labDelta * labDelta

	labEpsilon = labDelta2 * labDelta // 216.0 / 24389.0
	labInv3D2  = 1.0 / (3.0 * labDelta2)
	lab4Over29 = 4.0 / 29

	labKappa    = 24389.0 / 27.0
	labInvKappa = 1 / labKappa

	labInv116 = 1.0 / 116
	labInv500 = 1.0 / 500
	labInv200 = 1.0 / 200
)

func labF(t float64) float64 {
	if t > labEpsilon {
		return math.Cbrt(t)
	}
	return t*labInv3D2 + lab4Over29
}

func labInvF(t float64) float64 {
	if t > labDelta {
		return t * t * t
	}
	return 3 * labDelta2 * (t - lab4Over29)
}

func LabToXyzD50(l, a, b float64) (x, y, z float64) {
	fy := (l + 16) * labInv116
	fx := fy + a*labInv500
	fz := fy - b*labInv200

	x = d50X * labInvF(fx)
	y = d50Y * labInvF(fy)
	z = d50Z * labInvF(fz)

	return
}

func XyzD50ToLab(x, y, z float64) (l, a, b float64) {
	fx := labF(x * invD50X)
	fy := labF(y * invD50Y)
	fz := labF(z * invD50Z)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)

	return
}
