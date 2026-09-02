package colors

import (
	"bytes"
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

	// Trim trailing zeros and the decimal point.
	if bytes.IndexByte(dst[start:], '.') >= 0 {
		for n > start && dst[n-1] == '0' {
			n--
		}
		if n > start && dst[n-1] == '.' {
			n--
		}

		dst = dst[:n]
	}

	// Normalize negative zero.
	if n-start == 2 && dst[start+1] == '0' && dst[start] == '-' {
		dst[start] = '0'
		dst = dst[:start+1]
	}

	return dst
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

	// Trim trailing zeros and the decimal point.
	if strings.IndexByte(s, '.') >= 0 {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}

	// Normalize negative zero.
	if len(s) == 2 && s[1] == '0' && s[0] == '-' {
		return "0"
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

func clamp01(x float64) float64 {
	return min(1, max(0, x))
}

func clamp(x, lo, hi float64) float64 {
	return min(hi, max(lo, x))
}

func wrap01(v float64) float64 {
	return v - math.Floor(v)
}

func wrap360(x float64) float64 {
	const inv360 = 1 / 360.0
	return x - math.Floor(x*inv360)*360
}

func wrap(v, min, max float64) float64 {
	r := max - min
	v -= min
	v -= math.Floor(v/r) * r
	return v + min
}

func wrapMod(v, min, max float64) float64 {
	return math.Mod(v-min, max-min) + (min)
}
