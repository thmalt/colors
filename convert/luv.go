package convert

import "math"

const (
	luvWhiteDenInv = 1 / (d50X + 15*d50Y + 3*d50Z) // D65 whitepoint
	luvUn          = 4 * d50X * luvWhiteDenInv
	luvVn          = 9 * d50Y * luvWhiteDenInv

	luvLThreshold = labKappa * labEpsilon
)

func XyzD50ToLuv(x, y, z float64) (l, u, v float64) {
	if yr := y * invD50Y; yr > labEpsilon {
		l = 116*math.Cbrt(yr) - 16
	} else {
		l = labKappa * yr
	}

	den := x + 15*y + 3*z
	if den == 0 {
		return l, 0, 0
	}

	inv := 1 / den
	up := 4 * x * inv
	vp := 9 * y * inv

	f := 13 * l
	u = f * (up - luvUn)
	v = f * (vp - luvVn)

	return
}

func LuvToXyzD50(l, u, v float64) (x, y, z float64) {
	if l <= 0 {
		return 0, 0, 0
	}

	var yr float64
	if l > luvLThreshold {
		f := (l + 16) * labInv116
		yr = f * f * f
	} else {
		yr = l * labInvKappa
	}

	y = yr * d50Y

	inv13L := 1 / (13 * l)
	up := u*inv13L + luvUn
	vp := v*inv13L + luvVn

	if vp == 0 {
		return 0, 0, 0
	}

	t := y / (4 * vp)

	x = 9 * up * t
	z = (12 - 3*up - 20*vp) * t

	return
}
