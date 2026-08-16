package codegen

import (
	"fmt"
	"math"
	"reflect"
	"runtime/debug"
	"strings"
)

const (
	FloatType = "float64"

	DefaultPrecision = 6
	AlphaPrecision   = DefaultPrecision

	AppendFloatFormatNormalizedPrecFuncName = "appendFormatNormalizedFloatPrec"
	AppendFloatFormatPrecFuncName           = "appendFormatFloatPrec"

	PanicUnreachable = `panic("unreachable")`
)

func smallestUintType(n int) int {
	switch {
	case n <= math.MaxUint8:
		return 8
	case n <= math.MaxUint16:
		return 16
	case n <= math.MaxUint32:
		return 32
	default:
		return 64
	}
}

func joinRepeatN(v string, n int) string {
	return strings.Repeat(v+", ", n-1) + v
}

func joinIdentsWithType(typ string, vars ...string) string {
	return strings.Join(vars, ", ") + " " + typ
}

func buildChannelCounts(ctx *Context) []bool {
	counts := make([]bool, ctx.MaxChannelCount+1)
	for _, space := range ctx.BuildSpaces {
		counts[space.ChannelCount()] = true
	}
	return counts
}

func buildHueIndexes(ctx *Context) [][]bool {
	indexes := make([][]bool, ctx.MaxChannelCount+1)
	for _, space := range ctx.BuildSpaces {
		count := space.ChannelCount()
		if index := space.HueIndex(); index >= 0 {
			if len(indexes[count]) == 0 {
				indexes[count] = make([]bool, count)
			}
			indexes[count][index] = true
		}
	}
	return indexes
}

func ModuleAndPathByType(a any) (module string, path string) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", ""
	}

	t := reflect.TypeOf(a)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	pkgPath := t.PkgPath()
	for _, m := range info.Deps {
		if strings.HasPrefix(pkgPath, m.Path) {
			fmt.Println()
			module = m.Path
		}
	}

	path = "." + strings.TrimPrefix(pkgPath, module)

	return module, path
}
