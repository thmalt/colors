package codegen

import (
	"fmt"
	"math"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

const (
	// for UnsafeMixer, Color.ChannelX
	MinGeneratedChannelCount = 4

	FloatType = "float64"

	DefaultPrecision = 6
	AlphaPrecision   = DefaultPrecision

	AppendFloatFormatNormalizedPrecFuncName = "appendFormatNormalizedFloatPrec"
	AppendFloatFormatPrecFuncName           = "appendFormatFloatPrec"

	PanicUnreachable = `panic("unreachable")`

	defaultHubD50 = "XyzD50"
	defaultHubD65 = "XyzD65"
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
	for _, space := range ctx.BuiltSpaces {
		counts[space.ChannelCount()] = true
	}
	counts[ctx.MaxChannelCount] = true
	return counts
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

func wrapEvery(w *writer.GoWriter, n int) func() bool {
	count := 0
	return func() bool {
		count++
		if count%n == 0 {
			w.Indent()
			return true
		}
		return false
	}
}
