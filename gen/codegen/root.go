package codegen

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func GenerateRootPkg(ctx *Context) {
	var pkgPath = filepath.Join(ctx.Directory, ctx.RootPkg.Path)
	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))

	// generate file color_gen.go
	fileName := "color_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(
		ctx.ConvertPkg.Path,
		ctx.SpacePkg.Path,
	)

	genRootPkgColorType(ctx, w)
	genRootPkgColorMethods(ctx, w)
	genRootPkgColorConversionMethods(ctx, w)
	w.WriteGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.RootPkg.Name,
	)

	// generate file color_constructors_gen.go

	fileName = "color_constructors_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(ctx.SpacePkg.Path)

	genRootPkgColorConstructors(ctx, w)
	w.WriteGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.RootPkg.Name,
	)

	// generate file color_stringify_gen.go
	fileName = "color_stringify_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(
		"fmt",
		"strings",
		ctx.SpacePkg.Path,
	)

	genRootPkgColorStringify(ctx, w)
	w.WriteGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.RootPkg.Name,
	)
}

func genRootPkgColorType(ctx *Context, w *writer.GoWriter) {
	maxChannelCnt := 0
	for _, space := range ctx.BuildSpaces {
		maxChannelCnt = max(maxChannelCnt, space.ChannelCount())
	}

	w.Begin("type Color struct")
	w.LineWriteln("space ", ctx.SpacePkg.Join("Space"))
	w.LineCommentln("channels")

	w.LineWrite()

	for i := range maxChannelCnt {
		if i > 0 {
			w.Write(", ")
		}
		w.Writef("c%d", i+1)
	}

	w.Writeln(" ", FloatType)
	w.LineWrite("alpha ", FloatType)
	w.End()

	w.Writeln()

	w.LineCommentln("Get channel in current [space.Space]")
	w.Method("c Color", "Channel")
	w.FuncParams("index int")
	w.FuncResults(FloatType, ", ", "bool")
	w.FuncBody()

	w.Switch("index ")

	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i - 1)
		w.Return("c.c", i, ", true")
	}

	w.End()
	w.LineWritef("return 0, false") // panic(%q)", "Color.Channel unreachable

	w.End()

	w.Writeln()

	var b strings.Builder
	var cases []string

	for i := range maxChannelCnt {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "c.c%d", i+1)
		cases = append(cases, b.String())
	}

	// Channels
	w.Method("c Color", "Channels")
	w.FuncResults("[]", FloatType)
	w.FuncBody()

	w.LineWriteln("info := c.space.Info()")
	w.If("info == nil")
	w.Return("nil")
	w.End()

	w.Writeln()

	w.Switch("info.ChannelCount()")
	for i := maxChannelCnt; i > 0; i-- {
		w.Case(i)
		w.Return("[]", FloatType, "{", cases[i-1], "}")
	}
	w.Default()
	w.Return("nil")
	w.End()

	w.End()

	w.Writeln()

	// AppendChannels
	w.Method("c Color", "AppendChannels")
	w.FuncParams("dst []", FloatType)
	w.FuncResults("[]", FloatType)
	w.FuncBody()

	w.LineWriteln("info := c.space.Info()")
	w.If("info == nil")
	w.Return("dst")
	w.End()

	w.Writeln()

	w.LineWriteln("ch := [...]", FloatType, "{", b.String(), "}")
	w.Return("append(dst, ch[:info.ChannelCount()]...)")
	w.End()

	w.Writeln()
}

func genRootPkgColorMethods(ctx *Context, w *writer.GoWriter) {
	var scope VariableScope
	var convertPkg = ctx.ConvertPkg
	var spacePkg = ctx.SpacePkg
	var b strings.Builder
	for _, space := range ctx.BuildSpaces {
		b.Reset()

		scope.Reset()
		scope.Reserve("c")

		names := space.ChannelSymbols()

		if scope.ContainsAny(names...) {
			names = scope.ReserveUniqueAll(space.ChannelIdent()...)
		}

		w.LineCommentln(space.Name, " returns the color components in the [", ctx.SpacePkg.Join(space.Name), "] color space.")
		w.Method("c Color", space.Name)
		w.FuncResults(varJoinWithType(names...))
		w.FuncBody()

		w.If("c.space == ", spacePkg.Join(space.Name))
		for i := range names {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "c.c%d", i+1)
		}
		w.Return(b.String())
		w.End()

		w.Writeln()

		wsw := w.NewTemp()

		wsw.Switch("c.space ")
		var args strings.Builder
		var foundPath = false
		for _, src := range ctx.BuildSpaces {
			args.Reset()
			if space == src {
				continue
			}

			path := ctx.Graph.FindPath(src, space)
			if len(path) == 0 {
				log.Printf("no conversion path found: %q -> %q\n", src.Name, space.Name)
				continue
			}
			foundPath = true
			wsw.Case(spacePkg.Join(src.Name))
			for i := range names {
				if i > 0 {
					args.WriteString(", ")
				}
				fmt.Fprintf(&args, "c.c%d", i+1)
			}

			wsw.Return(convertPkg.Join(FuncName(src.Name, space.Name)), "(", args.String(), ")")
		}
		wsw.Default()
		wsw.Return()
		wsw.End()

		if foundPath {
			w.Write(wsw.Bytes())
		} else {
			w.Return()
		}

		w.End()
		w.Writeln()
	}
}

func genRootPkgColorConversionMethods(ctx *Context, w *writer.GoWriter) {
	var spacePkg = ctx.SpacePkg

	w.Method("c Color", "To")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("Color, error")
	w.FuncBody()
	w.If("c.space == dst")
	w.Return("c, nil")
	w.End()

	w.Writeln()
	var scope VariableScope

	w.Switch("dst")
	var b strings.Builder
	for _, space := range ctx.BuildSpaces {
		scope.Reset()
		scope.Reserve("c")

		w.Case(spacePkg.Join(space.Name))
		names := scope.ReserveUniqueAll(space.ChannelSymbols()...)
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
	w.Return("Color{}, ErrUnknownSpace")
	w.End()

	w.End()

	w.Writeln()

	w.Method("c Color", "MustTo")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("Color")
	w.FuncBody()
	w.LineWriteln("to, _ := c.To(dst)")
	w.Return("to")
	w.End()
}

func genRootPkgColorConstructors(ctx *Context, w *writer.GoWriter) {
	var b strings.Builder

	for _, space := range ctx.BuildSpaces {
		w.LineCommentln(space.Name, " returns a [Color] from ", space.DisplayName, " components.")
		w.LineCommentln()

		for _, c := range space.Channels {
			w.LineComment('\t')
			w.Write(c.Symbol)
			w.Write(": [")
			w.Write(formatNormalizedFloat(c.Min))
			w.Write(", ")
			w.Write(formatNormalizedFloat(c.Max))

			if c.Circular {
				w.Write(')')
			} else {
				w.Write(']')
			}

			if c.Unrestricted {
				w.Write(" (typical)")
			}

			w.Writeln()
		}

		w.Func(space.Name)
		params := space.ChannelSymbols()
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
		w.Writeln()
	}
}

func genRootPkgColorStringify(ctx *Context, w *writer.GoWriter) {
	w.Method("c Color", "String")
	w.FuncResults("string")
	w.FuncBody()

	w.LineWriteln("var b strings.Builder")
	w.LineWriteln("b.Grow(64)")

	w.Switch("c.space")

	for _, space := range ctx.BuildSpaces {
		w.Case(ctx.SpacePkg.Join(space.Name))
		if space.UseColorFunction {
			w.LineWriteln(`b.WriteString("color(`, space.CssName, ` ")`)
		} else {
			w.LineWriteln(`b.WriteString("`, space.CssName, `(")`)
		}

		for i, c := range space.Channels {
			if i > 0 {
				w.LineWriteln("b.WriteByte(' ')")
			}

			w.LineWrite("b.WriteString(")
			w.Write(FloatFormatNormalizedPrecFuncName, "(")
			percent := c.Unit == model.UnitPercent
			w.Write("c.c", i+1)

			if percent {
				w.Write(" * 100")
			}

			w.Write(", ")

			prec := c.Precision
			if prec == 0 {
				prec = DefaultPrecision
			}

			if percent {
				prec = max(0, prec-2)
			}

			w.Write(prec)

			w.Writeln("))")

			if percent {
				w.LineWriteln("b.WriteByte('%')")
			}
		}

		w.Writeln()

		w.If("alpha := normalizeFloat(c.alpha); alpha != 1")
		w.LineWriteln(`b.WriteString(" / ")`)
		w.LineWrite("b.WriteString(")
		w.Write(FloatFormatPrecFuncName, "(")
		w.Write("alpha, ", AlphaPrecision)
		w.Writeln("))")
		w.End()

		w.Writeln()

		w.LineWriteln(`b.WriteString(")")`)

	}

	w.Default()
	w.Return(`fmt.Sprintf("Color(<invalid space: %d>)", c.space)`)
	w.End()

	w.Writeln()

	w.Return("b.String()")

	w.End()
}
