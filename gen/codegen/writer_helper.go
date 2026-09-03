package codegen

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func pkgInfo(ctx *Context, pkg Pkg) (name, path string) {
	return pkg.Name, filepath.Join(ctx.Directory, pkg.Path)
}

func newWriter(ctx *Context) *writer.GoWriter {
	w := writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.Opts.FormatSource)

	return w
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

func joinRepeatN(v string, n int) string {
	return strings.Repeat(v+", ", n-1) + v
}

func joinIdentsWithType(typ string, vars ...string) string {
	return strings.Join(vars, ", ") + " " + typ
}

func appendVars(dst []string, ident string, count int, extra ...string) []string {
	for i := 1; i <= count; i++ {
		dst = append(dst, ident+strconv.Itoa(i))
	}

	if len(extra) > 0 {
		dst = append(dst, extra...)
	}

	return dst
}

func writeColorLiteral(w *writer.GoWriter, space string, channelCount int, vars ...string) {
	if len(vars) < channelCount {
		panic("writeColorLiteral: len(vars) < channelCount")
	}

	w.Write("Color{space: ", space, ", ")
	for i := range channelCount {
		w.Write("c", i+1, ": ", vars[i], ", ")
	}
	alpha := "1"
	if len(vars) > channelCount {
		alpha = vars[len(vars)-1]
	}
	w.Write("alpha: ", alpha, "}")
}

func writeColorLiteralMultiLine(w *writer.GoWriter, space string, channelCount int, vars ...string) {
	if len(vars) < channelCount {
		panic("writeColorLiteralMultiLine: len(vars) < channelCount")
	}

	w.Writeln("Color{")
	w.In()
	w.LineWriteln("space: ", space, ",")
	for i := range channelCount {
		w.LineWriteln("c", i+1, ":    ", vars[i], ",")
	}
	alpha := "1"
	if len(vars) > channelCount {
		alpha = vars[len(vars)-1]
	}
	w.LineWriteln("alpha: ", alpha, ",")
	w.Out()
	w.LineWriteln('}')
}
