package convert

import (
	"math"
)

const (
	labDelta  = 6.0 / 29
	labDelta2 = labDelta * labDelta

	labEpsilon = labDelta2 * labDelta // 216.0 / 24389.0
	labInv3D2  = 1.0 / (3.0 * labDelta2)
	lab4Over29 = 4.0 / 29

	labKappa    = 24389.0 / 27.0
	labInvKappa = 1 / labKappa

	labInv116 = 1.0 / 116
	labInv500 = 1.0 / 500
	labInv200 = 1.0 / 200
)

func labF(t float64) float64 {
	if t > labEpsilon {
		return math.Cbrt(t)
	}
	return t*labInv3D2 + lab4Over29
}

func labInvF(t float64) float64 {
	if t > labDelta {
		return t * t * t
	}
	return 3 * labDelta2 * (t - lab4Over29)
}
