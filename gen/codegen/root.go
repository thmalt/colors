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
	if ctx.RootPkg.Name == "" {
		return
	}

	var pkgPath = filepath.Join(ctx.Directory, ctx.RootPkg.Path)
	var w = writer.NewGoWriter()
	w.SetGeneratedBy(ctx.Module, "./"+filepath.Dir(ctx.Path))
	w.SetFormatSource(ctx.FormatSource)

	// generate file color_gen.go
	fileName := "color_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(ctx.SpacePkg.Path)

	genRootPkgColorType(ctx, w)
	w.SaveGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.RootPkg.Name,
	)

	// generate file color_gen.go
	fileName = "color_convert_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(
		ctx.ConvertPkg.Path,
		ctx.SpacePkg.Path,
	)

	genRootPkgColorConversionMethods(ctx, w)
	genRootPkgColorConvertMethods(ctx, w)
	w.SaveGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.RootPkg.Name,
	)

	// generate file color_constructors_gen.go

	fileName = "color_constructors_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(ctx.SpacePkg.Path)

	genRootPkgColorConstructors(ctx, w)
	w.SaveGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.RootPkg.Name,
	)

	// generate file color_string_gen.go
	fileName = "color_string_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(
		"fmt",
		"strings",
		ctx.SpacePkg.Path,
	)

	genRootPkgColorStringify(ctx, w)
	w.SaveGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.RootPkg.Name,
	)

	// generate file mix_gen.go
	fileName = "mix_gen.go"
	fmt.Println("  Generate file", fileName)

	w.Import(
		ctx.SpacePkg.Path,
	)

	genRootPkgMix(ctx, w)
	w.SaveGoFile(
		filepath.Join(pkgPath, fileName),
		ctx.RootPkg.Name,
	)
}

func genRootPkgColorType(ctx *Context, w *writer.GoWriter) {
	maxChannelCnt := 0
	for _, space := range ctx.BuildSpaces {
		maxChannelCnt = max(maxChannelCnt, space.ChannelCount())
	}

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
	w.LineWrite("alpha ", FloatType)
	w.End()

	w.Separate()

	w.Comment("Get channel in current [space.Space]")
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

	// Channels
	w.Method("c Color", "Channels")
	w.FuncResults("[]", FloatType)
	w.FuncBody()

	w.LineWriteln("info := c.space.Info()")
	w.If("info == nil")
	w.Return("nil")
	w.End()

	w.Separate()

	w.Switch("info.ChannelCount() ")
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

	w.LineWriteln("info := c.space.Info()")
	w.If("info == nil")
	w.Return("dst")
	w.End()

	w.Separate()

	w.LineWriteln("ch := [...]", FloatType, "{", b.String(), "}")
	w.Return("append(dst, ch[:info.ChannelCount()]...)")
	w.End()
}

func genRootPkgColorConvertMethods(ctx *Context, w *writer.GoWriter) {
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

		w.Separate()

		w.Comment(space.Name, " returns the color components in the [", ctx.SpacePkg.Join(space.Name), "] color space.")
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

		w.Separate()

		wsw := w.SubWriter()

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
			w.Drain(wsw)
		} else {
			w.Return()
		}

		w.End()
	}
}

func genRootPkgColorConversionMethods(ctx *Context, w *writer.GoWriter) {
	var spacePkg = ctx.SpacePkg

	w.Separate()

	w.Method("c Color", "To")
	w.FuncParams("dst ", spacePkg.Join("Space"))
	w.FuncResults("Color, error")
	w.FuncBody()
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
				gw.Write(c.Symbol)
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

		w.Separate()
	}
}

func genRootPkgColorStringify(ctx *Context, w *writer.GoWriter) {
	w.Method("c Color", "String")
	w.FuncResults("string")
	w.FuncBody()

	w.LineWriteln("var b strings.Builder")
	w.LineWriteln("b.Grow(64)")

	w.Separate()

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

		w.Separate()

		w.If("alpha := normalizeFloat(c.alpha); alpha != 1")
		w.LineWriteln(`b.WriteString(" / ")`)
		w.LineWrite("b.WriteString(")
		w.Write(FloatFormatPrecFuncName, "(")
		w.Write("alpha, ", AlphaPrecision)
		w.Writeln("))")
		w.End()

		w.Separate()

		w.LineWriteln(`b.WriteString(")")`)

	}

	w.Default()
	w.Return(`fmt.Sprintf("Color(<invalid space: %d>)", c.space)`)
	w.End()

	w.Separate()

	w.Return("b.String()")

	w.End()
}

func genRootPkgMix(ctx *Context, w *writer.GoWriter) {
	w.Func("Mix")
	w.FuncParams("c1, c2 Color, t float64, opts MixOptions")
	w.FuncResults("Color")
	w.FuncBody()

	w.If("c1.space.IsValid() && c2.space.IsValid() && !opts.Space.IsValid()")
	w.Return("Color{}")
	w.End()

	w.Separate()
	w.LineWrite("c1 = c1.MustTo(opts.Space)")
	w.LineWrite("c2 = c2.MustTo(opts.Space)")

	w.Separate()
	w.Switch("opts.Space ")
	for _, space := range ctx.BuildSpaces {
		w.Case(ctx.SpacePkg.Join(space.Name))
		w.Return("mix", space.Name, "(c1, c2, t, opts)")
	}
	w.End()

	w.Separate()
	w.LineWrite(`panic("unreachable")`)
	w.End()

	for _, space := range ctx.BuildSpaces {
		w.Separate()
		genMixSpaceFunction(ctx, w, space)
	}
}

func genMixSpaceFunction(ctx *Context, w *writer.GoWriter, space *model.Space) {
	c1 := "cl1"
	c2 := "cl2"

	w.Func("mix", space.Name)
	w.FuncParams(c1, ", ", c2, " Color, t float64, opts MixOptions")
	w.FuncResults("Color")
	w.FuncBody()

	channels := space.Channels
	v0 := space.ChannelSymbols()
	v1 := make([]string, 0, len(v0))
	v2 := make([]string, 0, len(v0))

	for _, sym := range v0 {
		v1 = append(v1, sym+"1")
		v2 = append(v2, sym+"2")
	}

	w.LineWriteJoin(v1, ", ")
	w.Write(" := ", c1, ".", space.Name, "()")

	w.LineWriteJoin(v2, ", ")
	w.Write(" := ", c2, ".", space.Name, "()")

	w.Separate()
	w.LineWrite("alpha1, alpha2 := ", c1, ".Alpha(), ", c2, ".Alpha()")

	w.Separate()
	w.LineWrite("alpha := lerp(alpha1, alpha2, t)")

	w.Separate()
	w.If("opts.Premultiplied")

	for i := range v1 {
		if !channels[i].Circular {
			w.LineWrite(v1[i], " *= ", "alpha1")
		}
	}

	w.Separate()
	for i := range v2 {
		if !channels[i].Circular {
			w.LineWrite(v2[i], " *= ", "alpha2")
		}
	}

	w.End()

	w.Separate()
	for i := range v0 {
		if channels[i].Circular {
			if i > 0 {
				w.Separate()
			}

			w.LineWrite("var ", v0[i], " ", FloatType)
			w.Switch("opts.Hue ")
			for _, hue := range HueInterpolation {
				w.Case(hue)
				w.LineWrite(v0[i], " = lerp", hue, "(", v1[i], ", ", v2[i], ", t)")
			}
			w.End()

			w.Separate()
		} else {
			w.LineWrite(v0[i], " := lerp(", v1[i], ", ", v2[i], ", t)")
		}
	}

	w.Separate()
	w.If("opts.Premultiplied")
	w.If("alpha != 0")
	w.LineWrite("inv := 1 / alpha")

	w.Separate()
	for i := range v0 {
		if !channels[i].Circular {
			w.LineWrite(v0[i], " *= inv")
		}
	}
	w.End()
	w.End()

	w.Separate()

	var b strings.Builder
	for i, v := range v0 {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString("c")
		b.WriteByte(byte('1' + i))
		b.WriteString(": ")
		b.WriteString(v)
	}
	w.Return("Color{space: ", ctx.SpacePkg.Join(space.Name), ", ", b.String(), ", alpha: alpha}")
	w.End()
}
