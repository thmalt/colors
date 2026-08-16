package colors

import (
	"math"
	"strconv"
	"strings"
)

func normalizeFloat(x float64) float64 {
	const eps = 1e-15

	if math.Abs(x) < eps {
		return 0
	}

	if math.Abs(x-1) < eps {
		return 1
	}

	if math.Abs(x+1) < eps {
		return -1
	}

	return x
}

func appendFormatFloatPrec(dst []byte, x float64, precision int) []byte {
	switch x {
	case 0:
		return append(dst, '0')
	case 1:
		return append(dst, '1')
	case -1:
		return append(dst, "-1"...)
	}

	start := len(dst)
	dst = strconv.AppendFloat(dst, x, 'f', precision, 64)

	n := len(dst)
	for n > start && dst[n-1] == '0' {
		n--
	}
	if n > start && dst[n-1] == '.' {
		n--
	}

	return dst[:n]
}

func appendFormatNormalizedFloatPrec(dst []byte, x float64, precision int) []byte {
	return appendFormatFloatPrec(dst, normalizeFloat(x), precision)
}

func formatFloatPrec(x float64, precision int) string {
	switch x {
	case 0:
		return "0"
	case 1:
		return "1"
	case -1:
		return "-1"
	}

	s := strconv.FormatFloat(x, 'f', precision, 64)
	if precision >= 0 {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}

	return s
}

func formatFloat(x float64) string {
	return formatFloatPrec(x, -1)
}

func formatNormalizedFloatPrec(x float64, precision int) string {
	return formatFloatPrec(normalizeFloat(x), precision)
}

func formatNormalizedFloat(x float64) string {
	return formatNormalizedFloatPrec(x, -1)
}

func clamp(x, lo, hi float64) float64 {
	return min(hi, max(lo, x))
}

func clamp01(x float64) float64 {
	return clamp(x, 0, 1)
}

func wrap(v, min, max float64) float64 {
	return math.Mod(v-min, max-min) + (min)
}
