package convert

import "math"

const (
	luvD65WhiteDenInv = 1 / (d65X + 15*d65Y + 3*d65Z)
	luvD65WhiteUPrime = 4 * d65X * luvD65WhiteDenInv
	luvD65WhiteVPrime = 9 * d65Y * luvD65WhiteDenInv
)

func XyzD65ToLuvD65(x, y, z float64) (l, u, v float64) {
	if yr := y * invD65Y; yr > labEpsilon {
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
	u = f * (up - luvD65WhiteUPrime)
	v = f * (vp - luvD65WhiteVPrime)

	return
}

func LuvD65ToXyzD65(l, u, v float64) (x, y, z float64) {
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

	y = yr * d65Y

	inv13L := 1 / (13 * l)
	up := u*inv13L + luvD65WhiteUPrime
	vp := v*inv13L + luvD65WhiteVPrime

	if vp == 0 {
		return 0, 0, 0
	}

	t := y / (4 * vp)

	x = 9 * up * t
	z = (12 - 3*up - 20*vp) * t

	return
}
