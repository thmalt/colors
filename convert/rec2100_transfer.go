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

// Rec2100PQDecode converts a normalized Rec. 2100 PQ component to normalized absolute luminance.
func Rec2100PQDecode(x float64) float64 {
	const (
		invM1    = 1 / rec2100PQM1
		invM2    = 1 / rec2100PQM2
		invScale = 1 / rec2100PQScale

		threshold = 7.309838247667101e-07 // pow(c1, m2)
	)

	abs := math.Abs(x)

	if abs <= threshold {
		return 0
	}

	p := math.Pow(abs, invM2)
	abs = math.Pow(
		(p-rec2100PQC1)/(rec2100PQC2-rec2100PQC3*p),
		invM1,
	) * invScale

	if x < 0 {
		return -abs
	}
	return abs
}

// Rec2100PQEncode converts normalized absolute luminance to a normalized Rec. 2100 PQ component.
func Rec2100PQEncode(x float64) float64 {
	if x == 0 {
		return 0
	}

	abs := math.Abs(x)

	p := math.Pow(abs*rec2100PQScale, rec2100PQM1)
	abs = math.Pow(
		(rec2100PQC1+rec2100PQC2*p)/(1+rec2100PQC3*p),
		rec2100PQM2,
	)

	if x < 0 {
		return -abs
	}
	return abs
}

const (
	rec2100HLGA = 0.17883277
	rec2100HLGB = 1.0 - 4.0*rec2100HLGA // 0.28466892
	rec2100HLGC = 0.559910729529562     // 0.5-a*log(4*a)

	rec2100HLGScale = 0.26496255978640015 // (exp((E-c)/a)+b)/12        E=0.75
)

// Rec2100HLGDecode converts a Rec. 2100 HLG signal component to linear-light component.
func Rec2100HLGDecode(x float64) float64 {
	const (
		invA         = 1 / rec2100HLGA
		invLowScale  = 1 / (rec2100HLGScale * 3.0)
		invHighScale = 1 / (rec2100HLGScale * 12.0)
	)

	abs := math.Abs(x)

	if abs <= 0.5 {
		abs *= abs * invLowScale
	} else {
		abs = (math.Exp((abs-rec2100HLGC)*invA) + rec2100HLGB) * invHighScale
	}

	if x < 0 {
		return -abs
	}
	return abs
}

// Rec2100HLGEncode converts a linear-light component to a Rec. 2100 HLG signal component.
func Rec2100HLGEncode(x float64) float64 {
	const (
		lowThreshold = 1 / (rec2100HLGScale * 12.0)
		lowScale     = rec2100HLGScale * 3.0
		highScale    = rec2100HLGScale * 12.0
	)

	abs := math.Abs(x)

	if abs <= lowThreshold {
		abs = math.Sqrt(abs * lowScale)
	} else {
		abs = rec2100HLGA*math.Log(abs*highScale-rec2100HLGB) + rec2100HLGC
	}

	if x < 0 {
		return -abs
	}
	return abs
}
