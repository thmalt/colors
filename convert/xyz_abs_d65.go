package convert

const (
	xyzAbsD65Scale    = 203.0
	xyzAbsD65InvScale = 1 / xyzAbsD65Scale
)

func XyzAbsD65ToXyzD65(x, y, z float64) (float64, float64, float64) {
	return x * xyzAbsD65InvScale, y * xyzAbsD65InvScale, z * xyzAbsD65InvScale
}

func XyzD65ToXyzAbsD65(x, y, z float64) (float64, float64, float64) {
	return x * xyzAbsD65Scale, y * xyzAbsD65Scale, z * xyzAbsD65Scale
}
