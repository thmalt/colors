package data

import "math"

func ChromaToXyz(x, y float64) [3]float64 {
	return [3]float64{x / y, 1, (1 - x - y) / y}
}

func ChromaticAdaptationMatrix(from, to [3]float64, cat [9]float64) [9]float64 {
	lmsFrom := Mat3MulVec3(cat, from)
	lmsTo := Mat3MulVec3(cat, to)

	scale := [3]float64{
		lmsTo[0] / lmsFrom[0],
		lmsTo[1] / lmsFrom[1],
		lmsTo[2] / lmsFrom[2],
	}

	scaled := Mat3ScaleRows(cat, scale)
	inv := Mat3Invert(cat)

	return Mat3Mul(inv, scaled)
}

func RgbToXyzMatrix(primaries [3][2]float64, white [3]float64) [9]float64 {
	m := [9]float64{}

	r := ChromaToXyz(primaries[0][0], primaries[0][1])
	g := ChromaToXyz(primaries[1][0], primaries[1][1])
	b := ChromaToXyz(primaries[2][0], primaries[2][1])

	m[0] = r[0]
	m[1] = g[0]
	m[2] = b[0]

	m[3] = r[1]
	m[4] = g[1]
	m[5] = b[1]

	m[6] = r[2]
	m[7] = g[2]
	m[8] = b[2]

	inv := Mat3Invert(m)
	scale := Mat3MulVec3(inv, white)

	return Mat3ScaleCols(m, scale)
}

func Mat3MulVec3(m [9]float64, v [3]float64) [3]float64 {
	x0 := m[0]*v[0] + m[1]*v[1] + m[2]*v[2]
	x1 := m[3]*v[0] + m[4]*v[1] + m[5]*v[2]
	x2 := m[6]*v[0] + m[7]*v[1] + m[8]*v[2]

	return [3]float64{x0, x1, x2}
}

func Mat3Mul(a, b [9]float64) [9]float64 {
	x0 := a[0]*b[0] + a[1]*b[3] + a[2]*b[6]
	x1 := a[0]*b[1] + a[1]*b[4] + a[2]*b[7]
	x2 := a[0]*b[2] + a[1]*b[5] + a[2]*b[8]

	x3 := a[3]*b[0] + a[4]*b[3] + a[5]*b[6]
	x4 := a[3]*b[1] + a[4]*b[4] + a[5]*b[7]
	x5 := a[3]*b[2] + a[4]*b[5] + a[5]*b[8]

	x6 := a[6]*b[0] + a[7]*b[3] + a[8]*b[6]
	x7 := a[6]*b[1] + a[7]*b[4] + a[8]*b[7]
	x8 := a[6]*b[2] + a[7]*b[5] + a[8]*b[8]

	return [9]float64{x0, x1, x2, x3, x4, x5, x6, x7, x8}
}

func Mat3Det(m [9]float64) float64 {
	return m[0]*(m[4]*m[8]-m[5]*m[7]) - m[1]*(m[3]*m[8]-m[5]*m[6]) + m[2]*(m[3]*m[7]-m[4]*m[6])
}

func Mat3Invert(m [9]float64) [9]float64 {
	det := Mat3Det(m)
	inv := 1 / det

	x0 := (m[4]*m[8] - m[5]*m[7]) * inv
	x1 := (m[2]*m[7] - m[1]*m[8]) * inv
	x2 := (m[1]*m[5] - m[2]*m[4]) * inv

	x3 := (m[5]*m[6] - m[3]*m[8]) * inv
	x4 := (m[0]*m[8] - m[2]*m[6]) * inv
	x5 := (m[2]*m[3] - m[0]*m[5]) * inv

	x6 := (m[3]*m[7] - m[4]*m[6]) * inv
	x7 := (m[1]*m[6] - m[0]*m[7]) * inv
	x8 := (m[0]*m[4] - m[1]*m[3]) * inv

	return [9]float64{x0, x1, x2, x3, x4, x5, x6, x7, x8}
}

func Mat3ScaleRows(m [9]float64, v [3]float64) [9]float64 {
	x0 := m[0] * v[0]
	x1 := m[1] * v[0]
	x2 := m[2] * v[0]

	x3 := m[3] * v[1]
	x4 := m[4] * v[1]
	x5 := m[5] * v[1]

	x6 := m[6] * v[2]
	x7 := m[7] * v[2]
	x8 := m[8] * v[2]

	return [9]float64{x0, x1, x2, x3, x4, x5, x6, x7, x8}

}

func Mat3ScaleCols(m [9]float64, v [3]float64) [9]float64 {
	x0 := m[0] * v[0]
	x1 := m[1] * v[1]
	x2 := m[2] * v[2]

	x3 := m[3] * v[0]
	x4 := m[4] * v[1]
	x5 := m[5] * v[2]

	x6 := m[6] * v[0]
	x7 := m[7] * v[1]
	x8 := m[8] * v[2]

	return [9]float64{x0, x1, x2, x3, x4, x5, x6, x7, x8}
}

// FMA

func ChromaticAdaptationMatrixFMA(from, to [3]float64, cat [9]float64) [9]float64 {
	lmsFrom := Mat3MulVec3FMA(cat, from)
	lmsTo := Mat3MulVec3FMA(cat, to)

	scale := [3]float64{
		lmsTo[0] / lmsFrom[0],
		lmsTo[1] / lmsFrom[1],
		lmsTo[2] / lmsFrom[2],
	}

	scaled := Mat3ScaleRows(cat, scale)
	inv := Mat3InvertFMA(cat)

	return Mat3MulFMA(inv, scaled)
}

func RgbToXyzMatrixFMA(primaries [3][2]float64, white [3]float64) [9]float64 {
	m := [9]float64{}

	r := ChromaToXyz(primaries[0][0], primaries[0][1])
	g := ChromaToXyz(primaries[1][0], primaries[1][1])
	b := ChromaToXyz(primaries[2][0], primaries[2][1])

	m[0] = r[0]
	m[1] = g[0]
	m[2] = b[0]

	m[3] = r[1]
	m[4] = g[1]
	m[5] = b[1]

	m[6] = r[2]
	m[7] = g[2]
	m[8] = b[2]

	inv := Mat3InvertFMA(m)
	scale := Mat3MulVec3FMA(inv, white)

	return Mat3ScaleCols(m, scale)
}

func Mat3MulVec3FMA(m [9]float64, v [3]float64) [3]float64 {
	return [3]float64{
		math.FMA(m[2], v[2], math.FMA(m[1], v[1], m[0]*v[0])),
		math.FMA(m[5], v[2], math.FMA(m[4], v[1], m[3]*v[0])),
		math.FMA(m[8], v[2], math.FMA(m[7], v[1], m[6]*v[0])),
	}
}

func Mat3MulFMA(a, b [9]float64) [9]float64 {
	var out [9]float64
	x0 := math.FMA(a[1], b[3], a[0]*b[0])
	out[0] = math.FMA(a[2], b[6], x0)

	x1 := math.FMA(a[1], b[4], a[0]*b[1])
	out[1] = math.FMA(a[2], b[7], x1)

	x2 := math.FMA(a[1], b[5], a[0]*b[2])
	out[2] = math.FMA(a[2], b[8], x2)

	x3 := math.FMA(a[4], b[3], a[3]*b[0])
	out[3] = math.FMA(a[5], b[6], x3)

	x4 := math.FMA(a[4], b[4], a[3]*b[1])
	out[4] = math.FMA(a[5], b[7], x4)

	x5 := math.FMA(a[4], b[5], a[3]*b[2])
	out[5] = math.FMA(a[5], b[8], x5)

	x6 := math.FMA(a[7], b[3], a[6]*b[0])
	out[6] = math.FMA(a[8], b[6], x6)

	x7 := math.FMA(a[7], b[4], a[6]*b[1])
	out[7] = math.FMA(a[8], b[7], x7)

	x8 := math.FMA(a[7], b[5], a[6]*b[2])
	out[8] = math.FMA(a[8], b[8], x8)

	return out
}

func Mat3DetFMA(m [9]float64) float64 {
	t0 := math.FMA(-m[5], m[7], m[4]*m[8])
	t1 := math.FMA(-m[5], m[6], m[3]*m[8])
	t2 := math.FMA(-m[4], m[6], m[3]*m[7])

	det := m[0] * t0
	det = math.FMA(-m[1], t1, det)
	det = math.FMA(m[2], t2, det)

	return det
}

func Mat3InvertFMA(m [9]float64) [9]float64 {
	det := Mat3DetFMA(m)
	inv := 1 / det

	x0 := math.FMA(-m[5], m[7], m[4]*m[8]) * inv
	x1 := math.FMA(-m[1], m[8], m[2]*m[7]) * inv
	x2 := math.FMA(-m[2], m[4], m[1]*m[5]) * inv

	x3 := math.FMA(-m[3], m[8], m[5]*m[6]) * inv
	x4 := math.FMA(-m[2], m[6], m[0]*m[8]) * inv
	x5 := math.FMA(-m[0], m[5], m[2]*m[3]) * inv

	x6 := math.FMA(-m[4], m[6], m[3]*m[7]) * inv
	x7 := math.FMA(-m[0], m[7], m[1]*m[6]) * inv
	x8 := math.FMA(-m[1], m[3], m[0]*m[4]) * inv

	return [9]float64{
		x0, x1, x2,
		x3, x4, x5,
		x6, x7, x8,
	}
}

func NewMat3() [9]float64 {
	return [9]float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}
