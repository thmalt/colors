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

	AlphaPrecision = 3

	FloatFormatPrecFuncName = "formatNormalizedFloatPrec"
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

func valueTypeRepeat(n int) string {
	return strings.Repeat(FloatType+", ", n-1) + " " + FloatType
}

func varJoinWithType(vars ...string) string {
	return strings.Join(vars, ", ") + " " + FloatType
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
	for _, d := range info.Deps {
		if strings.HasPrefix(pkgPath, d.Path) {
			fmt.Println()
			module = d.Path
		}
	}

	path = "." + strings.TrimPrefix(pkgPath, module)

	return module, path
}
