package convert

const (
	xyzAbsD65Scale    = 203.0
	invXyzAbsD65Scale = 1 / xyzAbsD65Scale
)

func XyzAbsD65ToXyzD65(x, y, z float64) (float64, float64, float64) {
	return x * invXyzAbsD65Scale, y * invXyzAbsD65Scale, z * invXyzAbsD65Scale
}

func XyzD65ToXyzAbsD65(x, y, z float64) (float64, float64, float64) {
	return x * xyzAbsD65Scale, y * xyzAbsD65Scale, z * xyzAbsD65Scale
}
