package colors

import (
	"math"

	"github.com/thmalt/colors/space"
)

func resolveStops(stops []GradientStop, mixer Mixer) []GradientStop {
	resolveStopOffsets(stops)
	convertStopColors(stops, mixer.Space())
	return resolveHints(stops, mixer)
}

func convertStopColors(stops []GradientStop, dst space.Space) {
	for i := range stops {
		if stops[i].IsHint() {
			continue
		}

		stops[i].Color, _ = stops[i].Color.To(dst)
	}
}

func resolveStopOffsets(stops []GradientStop) {
	n := len(stops)

	prev := -1
	for i := range stops {
		if !stops[i].HasOffset() {
			switch i {
			case 0:
				stops[i].Offset = 0
			case n - 1:
				stops[i].Offset = 1
			default:
				continue
			}
		}

		if prev >= 0 && stops[i].Offset < stops[prev].Offset {
			stops[i].Offset = stops[prev].Offset
		}

		prev = i
	}
}

// resolveHints returns new value
func resolveHints(stops []GradientStop, mixer Mixer) []GradientStop {
	resolved := make([]GradientStop, 0, len(stops))
	resolved = append(resolved, stops[0])

	last := len(stops) - 1
	for i := 1; i < last; i++ {
		if !stops[i].IsHint() {
			resolved = append(resolved, stops[i])
			continue
		}

		if stops[i-1].IsHint() || stops[i+1].IsHint() {
			continue
		}

		resolved = resolveHint(resolved, stops[i-1], stops[i], stops[i+1], mixer)
	}

	resolved = append(resolved, stops[last])
	return resolved
}

func resolveHint(dst []GradientStop, left, hint, right GradientStop, mixer Mixer) []GradientStop {
	offsetLeft := left.Offset
	offset := hint.Offset
	offsetRight := right.Offset

	leftDist := offset - offsetLeft
	rightDist := offsetRight - offset

	if nearlyEqual(leftDist, rightDist) {
		return dst
	}

	if nearlyEqual(leftDist, 0.0) {
		hint.Color = right.Color
		return append(dst, hint)
	}

	if nearlyEqual(rightDist, 0.0) {
		hint.Color = left.Color
		return append(dst, hint)
	}

	var stops [9]GradientStop

	if leftDist > rightDist {
		for y := 0; y < 7; y++ {
			stops[y].Offset = offsetLeft + leftDist*(float64(7.0+y)/13.0)
		}

		stops[7].Offset = offset + rightDist*(1.0/3.0)
		stops[8].Offset = offset + rightDist*(2.0/3.0)
	} else {
		stops[0].Offset = offsetLeft + leftDist*(1.0/3.0)
		stops[1].Offset = offsetLeft + leftDist*(2.0/3.0)

		for y := 0; y < 7; y++ {
			stops[y+2].Offset = offset + rightDist*(float64(y)/13.0)
		}

	}

	totalDist := offsetRight - offsetLeft
	hintRelativeOffset := leftDist / totalDist
	exponent := math.Log(0.5) / math.Log(hintRelativeOffset)

	for i := range stops {
		t := (stops[i].Offset - offsetLeft) / totalDist
		weighting := math.Pow(t, exponent)

		if math.IsNaN(weighting) || math.IsInf(weighting, 0) {
			continue
		}

		stops[i].Color = mixer.Mix(left.Color, right.Color, weighting)

		dst = append(dst, stops[i])
	}

	return dst
}

const nearlyEqualTolerance = 1.0 / 4096

func nearlyEqual(a, b float64) bool {
	return math.Abs(a-b) <= nearlyEqualTolerance
}
