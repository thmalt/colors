package colors

import (
	"math"

	"github.com/thmalt/colors/space"
)

// DeltaEOK returns the Euclidean distance between two [Color]s in the [space.Oklab] color space.
func DeltaEOK(c1, c2 Color) float64 {
	var l1, a1, b1 float64
	var l2, a2, b2 float64

	if c1.space == space.Oklab {
		l1, a1, b1 = c1.c1, c1.c2, c1.c3
	} else {
		l1, a1, b1 = c1.Oklab()
	}

	if c2.space == space.Oklab {
		l2, a2, b2 = c2.c1, c2.c2, c2.c3
	} else {
		l2, a2, b2 = c2.Oklab()
	}

	dL := l1 - l2
	da := a1 - a2
	db := b1 - b2

	return math.Sqrt(dL*dL + da*da + db*db)
}

// DeltaEOK2 returns the Euclidean distance between two [Color]s in the [space.Oklab] color space,
// with the a and b components scaled by 2.
func DeltaEOK2(c1, c2 Color) float64 {
	var l1, a1, b1 float64
	var l2, a2, b2 float64

	if c1.space == space.Oklab {
		l1, a1, b1 = c1.c1, c1.c2, c1.c3
	} else {
		l1, a1, b1 = c1.Oklab()
	}

	if c2.space == space.Oklab {
		l2, a2, b2 = c2.c1, c2.c2, c2.c3
	} else {
		l2, a2, b2 = c2.Oklab()
	}

	const scale = 2
	dL := l1 - l2
	da := (a1 - a2) * scale
	db := (b1 - b2) * scale

	return math.Sqrt(dL*dL + da*da + db*db)
}

// DeltaE76 returns the CIE76 color difference between two [Color]s in [space.LabD50].
func DeltaE76(c1, c2 Color) float64 {
	return deltaE76(c1, c2, false)
}

// DeltaE76LabD65 returns the CIE76 color difference between two [Color]s in [space.LabD65].
func DeltaE76LabD65(c1, c2 Color) float64 {
	return deltaE76(c1, c2, true)
}

// deltaE76 returns the CIE76 color difference between two [Color]s in [space.LabD50] or [space.LabD65].
//
// If d65 is [true], [space.LabD65] is used; otherwise, [space.LabD50] is used.
func deltaE76(c1, c2 Color, d65 bool) float64 {
	var l1, a1, b1 float64
	var l2, a2, b2 float64

	if d65 {
		if c1.space == space.LabD65 {
			l1, a1, b1 = c1.c1, c1.c2, c1.c3
		} else {
			l1, a1, b1 = c1.LabD65()
		}

		if c2.space == space.LabD65 {
			l2, a2, b2 = c2.c1, c2.c2, c2.c3
		} else {
			l2, a2, b2 = c2.LabD65()
		}
	} else {
		if c1.space == space.LabD50 {
			l1, a1, b1 = c1.c1, c1.c2, c1.c3
		} else {
			l1, a1, b1 = c1.LabD50()
		}

		if c2.space == space.LabD50 {
			l2, a2, b2 = c2.c1, c2.c2, c2.c3
		} else {
			l2, a2, b2 = c2.LabD50()
		}
	}

	dL := l1 - l2
	da := a1 - a2
	db := b1 - b2

	return math.Sqrt(dL*dL + da*da + db*db)
}

// DeltaE94 returns the CIE94 color difference between two [Color]s in [space.LabD50].
func DeltaE94(c1, c2 Color) float64 {
	return deltaE94(c1, c2, false)
}

// DeltaE94LabD65 returns the CIE94 color difference between two [Color]s in [space.LabD65].
func DeltaE94LabD65(c1, c2 Color) float64 {
	return deltaE94(c1, c2, true)
}

// deltaE94 returns the CIE94 color difference between two [Color]s in [space.LabD50] or [space.LabD65].
//
// If d65 is [true], [space.LabD65] is used; otherwise, [space.LabD50] is used.
func deltaE94(c1, c2 Color, d65 bool) float64 {
	const (
		kL, kC, kH = 1.0, 0.045, 0.015 // Graphic arts
		// kL, kC, kH = 2.0, 0.048, 0.014 // Textile
	)

	var l1, a1, b1 float64
	var l2, a2, b2 float64

	if d65 {
		if c1.space == space.LabD65 {
			l1, a1, b1 = c1.c1, c1.c2, c1.c3
		} else {
			l1, a1, b1 = c1.LabD65()
		}

		if c2.space == space.LabD65 {
			l2, a2, b2 = c2.c1, c2.c2, c2.c3
		} else {
			l2, a2, b2 = c2.LabD65()
		}
	} else {
		if c1.space == space.LabD50 {
			l1, a1, b1 = c1.c1, c1.c2, c1.c3
		} else {
			l1, a1, b1 = c1.LabD50()
		}

		if c2.space == space.LabD50 {
			l2, a2, b2 = c2.c1, c2.c2, c2.c3
		} else {
			l2, a2, b2 = c2.LabD50()
		}
	}

	dL := l2 - l1
	da := a2 - a1
	db := b2 - b1

	dLs := dL / kL

	ch1 := math.Sqrt(a1*a1 + b1*b1)
	ch2 := math.Sqrt(a2*a2 + b2*b2)
	dC := ch2 - ch1
	dCs := dC / (1 + kC*ch1)

	dH2 := da*da + db*db - dC*dC
	dHs := math.Sqrt(dH2) / (1 + kH*ch1)

	return math.Sqrt(dLs*dLs + dCs*dCs + dHs*dHs)
}

// DeltaE2000 returns the CIEDE2000 color difference between two [Color]s in [space.LabD50].
func DeltaE2000(c1, c2 Color) float64 {
	return deltaE2000(c1, c2, false)
}

// DeltaE2000LabD65 returns the CIEDE2000 color difference between two [Color]s in [space.LabD65].
func DeltaE2000LabD65(c1, c2 Color) float64 {
	return deltaE2000(c1, c2, true)
}

// deltaE2000 returns the CIEDE2000 color difference between two [Color]s in [space.LabD50] or [space.LabD65].
//
// If d65 is [true], [space.LabD65] is used; otherwise, [space.LabD50] is used.
func deltaE2000(c1, c2 Color, d65 bool) float64 {
	const (
		factor = 25 * 25 * 25 * 25 * 25 * 25 * 25

		piMul2   = math.Pi * 2
		degToRad = math.Pi / 180.0

		radInv25    = 180 / (25.0 * math.Pi)
		rad275Div25 = 275 / 25.0

		rad30Mul2 = 2 * (30 * degToRad)

		kL, kC, kH = 1.0, 1.0, 1.0
	)

	var l1, a1, b1 float64
	var l2, a2, b2 float64

	if d65 {
		if c1.space == space.LabD65 {
			l1, a1, b1 = c1.c1, c1.c2, c1.c3
		} else {
			l1, a1, b1 = c1.LabD65()
		}

		if c2.space == space.LabD65 {
			l2, a2, b2 = c2.c1, c2.c2, c2.c3
		} else {
			l2, a2, b2 = c2.LabD65()
		}
	} else {
		if c1.space == space.LabD50 {
			l1, a1, b1 = c1.c1, c1.c2, c1.c3
		} else {
			l1, a1, b1 = c1.LabD50()
		}

		if c2.space == space.LabD50 {
			l2, a2, b2 = c2.c1, c2.c2, c2.c3
		} else {
			l2, a2, b2 = c2.LabD50()
		}
	}

	dLPrime := l2 - l1
	lPrimeAvg := (l1 + l2) * 0.5

	lsq := lPrimeAvg - 50
	lsq *= lsq

	// Scale lightness difference.
	sl := 1 + (0.015*lsq)/math.Sqrt(20+lsq)
	dLs := dLPrime / (kL * sl)

	x1 := (math.Sqrt(a1*a1+b1*b1) + math.Sqrt(a2*a2+b2*b2)) * 0.5
	x2 := x1 * x1
	x4 := x2 * x2
	x7 := x4 * x2 * x1
	g := 0.5 * (1 - math.Sqrt(x7/(x7+factor)))

	a1Prime := a1 * (1 + g)
	a2Prime := a2 * (1 + g)

	c1Prime := math.Sqrt(a1Prime*a1Prime + b1*b1)
	c2Prime := math.Sqrt(a2Prime*a2Prime + b2*b2)

	dCPrime := c2Prime - c1Prime
	cPrimeAvg := (c1Prime + c2Prime) * 0.5
	cPrimeProduct := c1Prime * c2Prime

	// Scale chroma difference.
	sc := 1 + 0.045*cPrimeAvg
	dCs := dCPrime / (kC * sc)

	h1Prime := 0.0
	if b1 != 0 || a1Prime != 0 {
		h1Prime = math.Atan2(b1, a1Prime)
		if h1Prime < 0 {
			h1Prime += piMul2
		}
	}

	h2Prime := 0.0
	if b2 != 0 || a2Prime != 0 {
		h2Prime = math.Atan2(b2, a2Prime)
		if h2Prime < 0 {
			h2Prime += piMul2
		}
	}

	var (
		dHPrime float64
		rt, t   float64
	)

	if cPrimeProduct != 0 {
		dHPrime = h2Prime - h1Prime
		hSum := h1Prime + h2Prime

		hPrimeAvg := hSum * 0.5
		if math.Abs(dHPrime) > math.Pi {
			if hSum < piMul2 {
				hPrimeAvg += math.Pi
			} else {
				hPrimeAvg -= math.Pi
			}
		}

		if dHPrime > math.Pi {
			dHPrime -= piMul2
		} else if dHPrime < -math.Pi {
			dHPrime += piMul2
		}
		dHPrime = 2 * math.Sqrt(cPrimeProduct) * math.Sin(dHPrime*0.5)

		x2 := cPrimeAvg * cPrimeAvg
		x4 := x2 * x2
		x7 := x4 * x2 * cPrimeAvg
		rc := 2 * math.Sqrt(x7/(x7+factor))

		hx := hPrimeAvg*radInv25 - rad275Div25
		dTheta := rad30Mul2 * math.Exp(-hx*hx)
		rt = -rc * math.Sin(dTheta)

		const (
			sin6, cos6   = 0.10452846326765347, 0.9945218953682734 // math.Sincos(6/180.0*math.Pi)
			sin30, cos30 = 0.5, 0.8660254037844386                 // math.Sincos(30/180.0*math.Pi)
			sin63, cos63 = 0.8910065241883678, 0.4539904997395468  // math.Sincos(63/180.0*math.Pi)
		)

		sin, cos := math.Sincos(hPrimeAvg)

		cos2 := 2*cos*cos - 1
		sin2 := 2 * sin * cos

		cos2Mul2 := 2 * cos2

		cos3 := cos * (cos2Mul2 - 1)
		sin3 := sin * (cos2Mul2 + 1)

		cos4 := cos2Mul2*cos2 - 1
		sin4 := 2 * sin2 * cos2

		t = 1 -
			0.17*(cos*cos30+sin*sin30) +
			0.24*cos2 +
			0.32*(cos3*cos6-sin3*sin6) -
			0.20*(cos4*cos63+sin4*sin63)
	}

	// Scale hue difference.
	sh := 1 + 0.015*cPrimeAvg*t
	dHs := dHPrime / (kH * sh)

	return math.Sqrt(dLs*dLs + dCs*dCs + dHs*dHs + rt*dCs*dHs)
}

// DeltaEOK returns the Euclidean distance between two [Color]s in the [space.Oklab] color space.
func (c Color) DeltaEOK(other Color) float64 {
	return DeltaEOK(c, other)
}

// DeltaEOK2 returns the Euclidean distance between two [Color]s in the [space.Oklab] color space,
// with the a and b components scaled by 2.
func (c Color) DeltaEOK2(other Color) float64 {
	return DeltaEOK2(c, other)
}

// DeltaE76 returns the CIE76 color difference between two [Color]s in [space.LabD50].
func (c Color) DeltaE76(other Color) float64 {
	return deltaE76(c, other, false)
}

// DeltaE76LabD65 returns the CIE76 color difference between two [Color]s in [space.LabD65].
func (c Color) DeltaE76LabD65(other Color) float64 {
	return deltaE76(c, other, true)
}

// DeltaE94 returns the CIE94 color difference between two [Color]s in [space.LabD50].
func (c Color) DeltaE94(other Color) float64 {
	return deltaE94(c, other, false)
}

// DeltaE94LabD65 returns the CIE94 color difference between two [Color]s in [space.LabD65].
func (c Color) DeltaE94LabD65(other Color) float64 {
	return deltaE94(c, other, true)
}

// DeltaE2000 returns the CIEDE2000 color difference between two [Color]s in [space.LabD50].
func (c Color) DeltaE2000(other Color) float64 {
	return deltaE2000(c, other, false)
}

// DeltaE2000LabD65 returns the CIEDE2000 color difference between two [Color]s in [space.LabD65].
func (c Color) DeltaE2000LabD65(other Color) float64 {
	return deltaE2000(c, other, true)
}
