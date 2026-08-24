package convert

import (
	"math"
)

const (
	rec2100PQReferenceWhite = 203.0 // cd/m²

	rec2100PQScale = rec2100PQReferenceWhite / 10000.0

	rec2100PQM1 = 2610.0 / 16384
	rec2100PQM2 = 2523.0 / 32

	rec2100PQC1 = 3424.0 / 4096
	rec2100PQC2 = 2413.0 / 128
	rec2100PQC3 = 2392.0 / 128
)

// pq -> linear
func rec2100PQDecode(x float64) float64 {
	const (
		invM1    = 1 / rec2100PQM1
		invM2    = 1 / rec2100PQM2
		invScale = 1 / rec2100PQScale

		threshold = 7.309838247667101e-07 // pow(c1, m2)
	)

	neg := x < 0
	x = math.Abs(x)

	if x <= threshold {
		return 0
	}

	p := math.Pow(x, invM2)
	x = math.Pow(
		(p-rec2100PQC1)/(rec2100PQC2-rec2100PQC3*p),
		invM1,
	) * invScale

	if neg {
		return -x
	}
	return x
}

// linear -> pq
func rec2100PQEncode(x float64) float64 {
	if x == 0 {
		return 0
	}

	neg := x < 0
	x = math.Abs(x)

	p := math.Pow(x*rec2100PQScale, rec2100PQM1)
	x = math.Pow(
		(rec2100PQC1+rec2100PQC2*p)/(1+rec2100PQC3*p),
		rec2100PQM2,
	)

	if neg {
		return -x
	}
	return x
}

const (
	rec2100HLGA = 0.17883277
	rec2100HLGB = 1.0 - 4.0*rec2100HLGA // 0.28466892
	rec2100HLGC = 0.559910729529562     // a*log(4*a)

	rec2100HLGScale = 0.26496255978640015 // L=(exp((E-c)/a)+b)/12		E=0.75 (1000 nit = 0.7518)
)

// linear -> hlg
func rec2100HLGEncode(x float64) float64 {
	const (
		lowThreshold = 1 / (rec2100HLGScale * 12.0)
		lowScale     = rec2100HLGScale * 3.0
		highScale    = rec2100HLGScale * 12.0
	)

	neg := x < 0
	x = math.Abs(x)

	if x <= lowThreshold {
		x = math.Sqrt(x * lowScale)
	} else {
		x = rec2100HLGA*math.Log(x*highScale-rec2100HLGB) + rec2100HLGC
	}

	if neg {
		return -x
	}
	return x
}

// hlg -> linear
func rec2100HLGDecode(x float64) float64 {
	const (
		invA         = 1 / rec2100HLGA
		invLowScale  = 1 / (rec2100HLGScale * 3.0)
		invHighScale = 1 / (rec2100HLGScale * 12.0)
	)

	neg := x < 0
	x = math.Abs(x)

	if x <= 0.5 {
		x *= x * invLowScale
	} else {
		x = (math.Exp((x-rec2100HLGC)*invA) + rec2100HLGB) * invHighScale
	}

	if neg {
		return -x
	}
	return x
}
