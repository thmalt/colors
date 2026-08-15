package gradient

import "math"

type Transform struct {
	m00, m01 float64
	m10, m11 float64
	tx, ty   float64
}

func Identity() Transform {
	return Transform{
		m00: 1,
		m11: 1,
	}
}

func (t *Transform) Reset() {
	t.m00 = 1
	t.m01 = 0
	t.m10 = 0
	t.m11 = 1
	t.tx = 0
	t.ty = 0
}

func (t *Transform) Translate(x, y float64) {
	t.tx += t.m00*x + t.m01*y
	t.ty += t.m10*x + t.m11*y
}

func (t *Transform) Scale(sx, sy float64) {
	t.m00 *= sx
	t.m01 *= sy
	t.m10 *= sx
	t.m11 *= sy
}

func (t *Transform) ScaleUniform(s float64) {
	t.m00 *= s
	t.m01 *= s
	t.m10 *= s
	t.m11 *= s
}

// radian
func (t *Transform) Rotate(angle float64) {
	c := math.Cos(angle)
	s := math.Sin(angle)

	m00 := t.m00*c + t.m01*s
	m01 := -t.m00*s + t.m01*c
	m10 := t.m10*c + t.m11*s
	m11 := -t.m10*s + t.m11*c

	t.m00 = m00
	t.m01 = m01
	t.m10 = m10
	t.m11 = m11
}

// radian
func (t *Transform) SkewX(angle float64) {
	k := math.Tan(angle)

	t.m01 += t.m00 * k
	t.m11 += t.m10 * k
}

// radian
func (t *Transform) SkewY(angle float64) {
	k := math.Tan(angle)

	t.m00 += t.m01 * k
	t.m10 += t.m11 * k
}

func (t *Transform) Matrix(m Transform) {
	m00 := t.m00*m.m00 + t.m01*m.m10
	m01 := t.m00*m.m01 + t.m01*m.m11
	m10 := t.m10*m.m00 + t.m11*m.m10
	m11 := t.m10*m.m01 + t.m11*m.m11

	tx := t.m00*m.tx + t.m01*m.ty + t.tx
	ty := t.m10*m.tx + t.m11*m.ty + t.ty

	t.m00 = m00
	t.m01 = m01
	t.m10 = m10
	t.m11 = m11
	t.tx = tx
	t.ty = ty
}

func (t *Transform) PositionAt(x, y float64) (float64, float64) {
	return t.m00*x + t.m01*y + t.tx, t.m10*x + t.m11*y + t.ty
}

func (t *Transform) VectorAt(x, y float64) (float64, float64) {
	return t.m00*x + t.m01*y, t.m10*x + t.m11*y
}

func (t *Transform) Determinant() float64 {
	return t.m00*t.m11 - t.m01*t.m10
}

func (t Transform) Inverse() (Transform, bool) {
	det := t.Determinant()
	if det == 0 {
		return Transform{}, false
	}

	invDet := 1 / det

	return Transform{
		m00: t.m11 * invDet,
		m01: -t.m01 * invDet,
		m10: -t.m10 * invDet,
		m11: t.m00 * invDet,

		tx: (t.m01*t.ty - t.m11*t.tx) * invDet,
		ty: (t.m10*t.tx - t.m00*t.ty) * invDet,
	}, true
}
