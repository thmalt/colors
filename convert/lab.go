package convert

import (
	"math"
)

const (
	labDelta   = 6.0 / 29.0
	labDelta2  = labDelta * labDelta
	labDelta3  = labDelta2 * labDelta
	labInv3D2  = 1.0 / (3.0 * labDelta2) // 841 / 108
	lab4Over29 = 4.0 / 29.0
)

func LabToLch(l, a, b float64) (float64, float64, float64) {
	h := math.Atan2(b, a) * (180. / math.Pi)
	if h < 0 {
		h += 360
	}

	c := math.Hypot(a, b)

	return l, c, h
}

func LchToLab(l, c, h float64) (float64, float64, float64) {
	rad := h * (math.Pi / 180.)
	a := c * math.Cos(rad)
	b := c * math.Sin(rad)

	return l, a, b
}

func labF(t float64) float64 {
	if t > labDelta3 {
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
	fy := (l + 16) / 116
	fx := fy + a/500
	fz := fy - b/200

	x = D50[0] * labInvF(fx)
	y = D50[1] * labInvF(fy)
	z = D50[2] * labInvF(fz)

	return
}

func XyzD50ToLab(x, y, z float64) (l, a, b float64) {
	fx := labF(x / D50[0])
	fy := labF(y / D50[1])
	fz := labF(z / D50[2])

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)

	return
}

func LabToXyz(l, a, b float64, white [3]float64) (x, y, z float64) {
	fy := (l + 16) / 116
	fx := fy + a/500
	fz := fy - b/200

	x = white[0] * labInvF(fx)
	y = white[1] * labInvF(fy)
	z = white[2] * labInvF(fz)

	return
}

func XyzToLab(x, y, z float64, white [3]float64) (l, a, b float64) {
	fx := labF(x / white[0])
	fy := labF(y / white[1])
	fz := labF(z / white[2])

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)

	return
}
