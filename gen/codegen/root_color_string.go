package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgColorStringMethod(ctx *Context, w *writer.GoWriter) {
	w.Method("c Color", "String")
	w.FuncResults("string")
	w.FuncBody()

	invalidReturn := `"Color(<invalid space: " + strconv.FormatUint(uint64(c.space), 10) + ">)"`
	unhandledReturn := `"Color(<unhandled space: " + strconv.FormatUint(uint64(c.space), 10) + ">)"`
	w.If("!c.space.IsValid()")
	w.Return(invalidReturn)
	w.End()

	w.Separate()
	w.LineWriteln("var b strings.Builder")
	w.LineWriteln("b.Grow(64)")

	w.Separate()
	w.Switch("c.space")

	for _, space := range ctx.BuildSpaces {
		w.Case(ctx.SpacePkg.Join(space.Name))
		if space.UseGenericColorFunction {
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
				w.Write("*100")
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
	w.Return(unhandledReturn)
	w.End()

	w.Separate()

	w.Return("b.String()")

	w.End()
}
