package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateRootPkg(ctx *Context) {
	if ctx.RootPkg.Name == "" {
		return
	}

	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.FormatSource)

	pkgPath := filepath.Join(ctx.Directory, ctx.RootPkg.Path)
	pkg := ctx.RootPkg.Name

	emitGoFile(w, pkg, pkgPath, "color", func(w *writer.GoWriter) {
		w.Import(
			ctx.SpacePkg.Path,
		)

		genRootPkgColorType(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "color_convert", func(w *writer.GoWriter) {
		w.Import(
			ctx.ConvertPkg.Path,
			ctx.SpacePkg.Path,
		)

		genRootPkgColorConversionMethods(ctx, w)
		genRootPkgColorConvertMethods(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "color_constructors", func(w *writer.GoWriter) {
		w.Import(
			ctx.SpacePkg.Path,
		)

		genRootPkgColorConstructors(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "color_string", func(w *writer.GoWriter) {
		w.Import(
			"strconv",
			"strings",
			ctx.SpacePkg.Path,
		)

		genRootPkgColorStringMethod(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "mix", func(w *writer.GoWriter) {
		w.Import(
			ctx.SpacePkg.Path,
		)

		genRootPkgMix(ctx, w)
	})

	emitGoFile(w, pkg, pkgPath, "clamp", func(w *writer.GoWriter) {
		w.Import(
			ctx.SpacePkg.Path,
		)

		genRootPkgClamp(ctx, w)
	})
}

func genRootPkgColorType(ctx *Context, w *writer.GoWriter) {
	maxChannelCnt := 0
	for _, space := range ctx.BuildSpaces {
		maxChannelCnt = max(maxChannelCnt, space.ChannelCount())
	}

	w.Comment("Color represents a color in a [space.Space].")
	w.Begin("type Color struct ")
	w.LineWriteln("space ", ctx.SpacePkg.Join("Space"))
	w.Comment("channels")

	w.Indent()

	for i := range maxChannelCnt {
		if i > 0 {
			w.Write(", ")
		}
		w.Writef("c%d", i+1)
	}

	w.Writeln(" ", FloatType)

	w.Separate()
	w.LineWrite("alpha ", FloatType)
	w.End()

	w.Separate()

	w.Comment("Channel returns the channel value at index in the order defined by the")
	w.Comment("color's [space.Space]. The returned boolean reports whether index is valid.")
	w.Method("c Color", "Channel")
	w.FuncParams("index int")
	w.FuncResults(FloatType, ", ", "bool")
	w.FuncBody()

	w.If("index < 0 || index >= c.ChannelCount()")
	w.Return("0, false")
	w.End()

	w.Separate()
	w.Switch("index ")

	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i - 1)
		w.Return("c.c", i, ", true")
	}

	w.End()

	w.Separate()
	w.LineWritef("return 0, false") // panic(%q)", "Color.Channel unreachable

	w.End()

	var b strings.Builder
	var cases []string

	for i := range maxChannelCnt {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "c.c%d", i+1)
		cases = append(cases, b.String())
	}

	w.Separate()

	w.Comment("Channels returns the channel values of c in the order defined by its [space.Space].")
	w.Method("c Color", "Channels")
	w.FuncResults("[]", FloatType)
	w.FuncBody()

	w.Switch("c.ChannelCount() ")
	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i)
		w.Return("[]", FloatType, "{", cases[i-1], "}")
	}
	w.Default()
	w.Return("nil")
	w.End()

	w.End()

	w.Separate()

	// AppendChannels
	w.Method("c Color", "AppendChannels")
	w.FuncParams("dst []", FloatType)
	w.FuncResults("[]", FloatType)
	w.FuncBody()
	w.Switch("c.ChannelCount() ")
	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i)
		w.Return("append(dst, ", cases[i-1], ")")
	}
	w.Default()
	w.Return("dst")
	w.End()

	w.End()
}

func genRootPkgColorConversionMethods(ctx *Context, w *writer.GoWriter) {
	var spacePkg = ctx.SpacePkg

	w.Separate()

	w.Method("c Color", "To")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("Color, error")
	w.FuncBody()
	w.If("!c.space.IsValid() || !dst.IsValid()")
	w.Return("Color{}, ErrInvalidSpace")
	w.End()
	w.Separate()
	w.If("c.space == dst")
	w.Return("c, nil")
	w.End()

	var scope VariableScope

	w.Separate()

	w.Switch("dst")
	var b strings.Builder
	for _, space := range ctx.BuildSpaces {
		scope.Reset()
		scope.Reserve("c")

		w.Case(spacePkg.Join(space.Name))
		names := scope.ReserveUniqueAll(space.ChannelIdent()...)
		w.LineWriteln(strings.Join(names, ", "), " := ", "c.", space.Name, "()")

		b.Reset()
		for i, name := range names {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "c%d: %s", i+1, name)
		}
		w.Return("Color{space: space.", space.Name, ", ", b.String(), ", alpha: c.alpha}, nil")
	}
	w.Default()
	w.Return("Color{}, ErrInvalidSpace")
	w.End()

	w.End()
}

func genRootPkgColorConstructors(ctx *Context, w *writer.GoWriter) {
	var b strings.Builder

	for _, space := range ctx.BuildSpaces {
		w.CommentFunc(func(gw *writer.GoWriter) {
			gw.Writeln(space.Name, " returns a [Color] from ", space.DisplayName, " components.")
			gw.Newline()

			for i, c := range space.Channels {
				if i > 0 {
					gw.Newline()
				}

				gw.Write('\t')
				gw.Write(c.Ident)
				gw.Write(": [")
				gw.Write(formatNormalizedFloat(c.Min))
				gw.Write(", ")
				gw.Write(formatNormalizedFloat(c.Max))

				if c.Circular {
					gw.Write(')')
				} else {
					gw.Write(']')
				}

				if c.Unrestricted {
					gw.Write(" (typical)")
				}
			}
		})

		w.Func(space.Name)
		params := space.ChannelIdent()
		w.FuncParams(strings.Join(params, ", "), " ", FloatType)
		w.FuncResults("Color")
		w.FuncBody()

		b.Reset()
		for i, name := range params {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "c%d: %s", i+1, name)
		}
		w.Return("Color{space: space.", space.Name, ", ", b.String(), ", alpha: 1}")

		w.End()

		w.Separate()
	}
}
