package codegen

import (
	_ "embed"
	"math"
	"strconv"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateNamedPkg(ctx *Context) {
	pkg := ctx.NamedPkg
	if pkg.Name == "" {
		return
	}

	w := newWriter(ctx)

	emitGoFile(ctx, pkg, w, "colors", func(w *writer.GoWriter) {
		w.Import(ctx.RootPkg.Path)

		genNamedPkgNamedVar(ctx, w)
	})

	emitGoFile(ctx, pkg, w, "lookup", func(w *writer.GoWriter) {
		w.Import(
			"strings",
			ctx.RootPkg.Path,
		)

		genNamedPkgLookup(ctx, w)
	})
}

func genNamedPkgNamedVar(ctx *Context, w *writer.GoWriter) {
	pkgJoin := func(ident string) string {
		if ctx.NamedPkg == ctx.RootPkg {
			return ident
		}
		return ctx.RootPkg.Join(ident)
	}

	var (
		temp []byte
		rgb  string
		fn   string
	)

	w.BeginGroup("var ")
	for i, nc := range ctx.NamedColors {
		temp = temp[:0]
		temp = strconv.AppendUint(temp, uint64(nc.RGBA[0]), 10)
		temp = append(temp, ',', ' ')
		temp = strconv.AppendUint(temp, uint64(nc.RGBA[1]), 10)
		temp = append(temp, ',', ' ')
		temp = strconv.AppendUint(temp, uint64(nc.RGBA[2]), 10)

		opaque := nc.RGBA[3] == math.MaxUint8
		if opaque {
			rgb = "rgb"
			fn = "Rgb"
		} else {
			rgb = "rgba"
			fn = "Rgba"

			temp = append(temp, ',', ' ')
			temp = appendFormatFloatPrec(temp, float64(nc.RGBA[3])/math.MaxUint8, AlphaPrecision)
		}

		if i > 0 {
			w.Separate()
		}

		w.Comment(nc.Name, " is the CSS named color ", '"', nc.Lower(), '"')
		w.Comment('\t', rgb, '(', temp, ')')
		w.LineWriteln(nc.Name, " = ", pkgJoin(fn), '(', temp, ')')
	}

	w.End()
}

func genNamedPkgLookup(ctx *Context, w *writer.GoWriter) {
	pkgJoin := func(ident string) string {
		if ctx.NamedPkg.Name == ctx.RootPkg.Name {
			return ident
		}
		return ctx.RootPkg.Join(ident)
	}

	w.Separate()
	w.Comment("Lookup returns the named color with the given name.")
	w.Comment()
	w.Comment("The name is case-insensitive.")
	w.Func("Lookup")
	w.FuncParams("name string")
	w.FuncResults(pkgJoin("Color"), ", bool")
	w.FuncBody()
	w.LineWriteln("c, ok := lookup[strings.ToLower(name)]")
	w.Return("c, ok")
	w.End()

	w.Separate()
	w.Begin("var lookup = map[string]", pkgJoin("Color"))
	for _, nc := range ctx.NamedColors {
		w.LineWriteln('"', nc.Lower(), '"', ": ", nc.Name, ',')
	}
	w.End()
}
