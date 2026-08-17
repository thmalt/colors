package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgColorStringMethod(ctx *Context, w *writer.GoWriter) {
	sub := w.SubWriter()
	groups := rootPkgColorStringsSpace(ctx, sub)

	w.Separate()
	w.Method("c Color", "String")
	w.FuncResults("string")
	w.FuncBody()

	invalidReturn := `"Color(<invalid space: " + strconv.FormatUint(uint64(c.space), 10) + ">)"`
	unhandledReturn := `"Color(<unhandled space: " + strconv.FormatUint(uint64(c.space), 10) + ">)"`

	w.If("!c.space.IsValid()")
	w.Return(invalidReturn)
	w.End()

	w.Separate()
	w.LineWriteln("buf := make([]byte, 0, 64)")

	w.Separate()
	w.Switch("c.space")

	funcNames := make([]string, len(ctx.BuiltSpaces))
	for _, group := range groups {
		for _, index := range group.Indexes {
			funcNames[index] = "appendFormat" + group.Spaces[0].Name
		}
	}

	vars := appendVars(nil, "c.c", ctx.MaxChannelCount)
	for i, space := range ctx.BuiltSpaces {
		w.Case(ctx.SpacePkg.Join(space.Name))

		if space.UseGenericColorFunction {
			w.LineWriteln(`buf = append(buf, "color(`, space.CssName, ` "...)`)
		} else {
			w.LineWriteln(`buf = append(buf, "`, space.CssName, `("...)`)
		}

		w.LineWrite("buf = ", funcNames[i], '(')
		w.Write("buf, ")
		w.WriteJoin(vars[:space.ChannelCount()], ", ")
		w.Writeln(')')

	}

	w.Default()
	w.Return(unhandledReturn)
	w.End()

	w.Separate()

	w.If("alpha := normalizeFloat(c.alpha); alpha != 1")
	w.LineWriteln(`buf = append(buf, " / "...)`)
	w.LineWrite("buf = ")
	w.Write(AppendFloatFormatPrecFuncName, "(buf, ")
	w.Write("alpha, ", AlphaPrecision)
	w.Writeln(")")
	w.End()

	w.Separate()
	w.LineWriteln("buf = append(buf, ')')")

	w.Separate()
	w.Return("unsafe.String(unsafe.SliceData(buf), len(buf))")

	w.End()

	vars = appendVars(vars[:0], "c", ctx.MaxChannelCount)
	for _, group := range groups {
		sub.Reset()
		space := group.Spaces[0]
		w.Separate()
		w.Func(funcNames[group.Indexes[0]])
		for i := range space.Channels {
			if i > 0 {
				sub.Write(", ")
			}
			sub.Write("c", i+1)
		}
		w.FuncParams("dst []byte, ", joinIdentsWithType(FloatType, vars[:space.ChannelCount()]...))
		w.FuncResults("[]byte")
		w.FuncBody()
		w.Write(group.Key)
		w.Return("dst")
		w.End()
	}
}

func rootPkgColorStringsSpace(ctx *Context, w *writer.GoWriter) []groupSpaceValue {
	defer w.Reset()

	gs := newGroupSpace()
	for idx, space := range ctx.BuiltSpaces {
		w.Reset()
		w.In()

		for i, ch := range space.Channels {
			if i > 0 {
				w.LineWriteln("dst = append(dst, ' ')")
			}

			w.LineWrite("dst = ", AppendFloatFormatNormalizedPrecFuncName, "(dst, ")
			percent := ch.Unit == model.UnitPercent
			w.Write("c", i+1)

			if percent {
				w.Write("*100")
			}

			w.Write(", ")

			prec := ch.Precision
			if prec == 0 {
				prec = DefaultPrecision
			}

			if percent {
				prec = max(0, prec-2)
			}

			w.Write(prec)

			w.Writeln(")")

			if percent {
				w.LineWriteln("dst = append(dst, '%')")
			}
		}

		gs.Append(string(w.Bytes()), space.ChannelCount(), idx, space)
	}

	return gs.Slice()
}
