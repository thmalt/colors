package data

const XyzAbsD65Scale = 203.0

var (
	D50 = Chromaticity{0.34570, 0.35850} // old: 0.34567, 0.35850, css4: 0.34570, 0.35850
	D65 = Chromaticity{0.31270, 0.32900}

	Bradford = [9]float64{
		0.8951, 0.2664, -0.1614,
		-0.7502, 1.7135, 0.0367,
		0.0389, -0.0685, 1.0296,
	}

	CAT02 = [9]float64{
		0.7328, 0.4296, -0.1624,
		-0.7036, 1.6975, 0.0061,
		0.0030, 0.0136, 0.9834,
	}

	CAT16 = [9]float64{
		0.401288, 0.650173, -0.051461,
		-0.250268, 1.204414, 0.045854,
		-0.002079, 0.048952, 0.953127,
	}

	SrgbPrimaries = [3][2]float64{
		{0.64, 0.33},
		{0.30, 0.60},
		{0.15, 0.06},
	}

	DisplayP3RgbPrimaries = [3][2]float64{
		{0.680, 0.320},
		{0.265, 0.690},
		{0.150, 0.060},
	}

	A98RgbPrimaries = [3][2]float64{
		{0.6400, 0.3300},
		{0.2100, 0.7100},
		{0.1500, 0.0600},
	}

	ProPhotoRgbPrimaries = [3][2]float64{
		{0.7347, 0.2653},
		{0.1596, 0.8404},
		{0.0366, 0.0001},
	}

	Rec2020RgbPrimaries = [3][2]float64{
		{0.708, 0.292},
		{0.170, 0.797},
		{0.131, 0.046},
	}

	OklabLmsToLab = [9]float64{
		0.2104542683093140, 0.7936177747023054, -0.0040720430116193,
		1.9779985324311684, -2.4285922420485799, 0.4505937096174110,
		0.0259040424655478, 0.7827717124575296, -0.8086757549230774,
	}

	OklabXyzD65ToLms = [9]float64{
		0.8190224379967030, 0.3619062600528904, -0.1288737815209879,
		0.0329836539323885, 0.9292868615863434, 0.0361446663506424,
		0.0481771893596242, 0.2642395317527308, 0.6335478284694309,
	}
)

var (
	D50Xyz = ChromaticityToXyz(D50.X, D50.Y)
	D65Xyz = ChromaticityToXyz(D65.X, D65.Y)

	OklabLabToLms    = Mat3InvertFMA(OklabLmsToLab)
	OklabLmsToXyzD65 = Mat3InvertFMA(OklabXyzD65ToLms)

	XyzD50ToXyzD65 = ChromaticAdaptationMatrixFMA(D50Xyz, D65Xyz, Bradford)
	XyzD65ToXyzD50 = ChromaticAdaptationMatrixFMA(D65Xyz, D50Xyz, Bradford)

	LinearSrgbToXyzD65 = RgbToXyzMatrixFMA(SrgbPrimaries, D65Xyz)
	XyzD65ToLinearSrgb = Mat3InvertFMA(LinearSrgbToXyzD65)

	LinearDisplayP3ToXyzD65 = RgbToXyzMatrixFMA(DisplayP3RgbPrimaries, D65Xyz)
	XyzD65ToLinearDisplayP3 = Mat3InvertFMA(LinearDisplayP3ToXyzD65)

	LinearA98ToXyzD65 = RgbToXyzMatrixFMA(A98RgbPrimaries, D65Xyz)
	XyzD65ToLinearA98 = Mat3InvertFMA(LinearA98ToXyzD65)

	LinearProPhotoToXyzD50 = RgbToXyzMatrixFMA(ProPhotoRgbPrimaries, D50Xyz)
	XyzD50ToLinearProPhoto = Mat3InvertFMA(LinearProPhotoToXyzD50)

	LinearRec2020ToXyzD65 = RgbToXyzMatrixFMA(Rec2020RgbPrimaries, D65Xyz)
	XyzD65ToLinearRec2020 = Mat3InvertFMA(LinearRec2020ToXyzD65)
)

type Chromaticity struct {
	X, Y float64
}
