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

	FloatFormatNormalizedPrecFuncName = "formatNormalizedFloatPrec"
	FloatFormatPrecFuncName           = "formatFloatPrec"

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

func repeatType(typ string, n int) string {
	return strings.Repeat(typ+", ", n-1) + typ
}

func joinIdentsWithType(typ string, vars ...string) string {
	return strings.Join(vars, ", ") + " " + typ
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
