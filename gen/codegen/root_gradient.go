package codegen

import "github.com/thmalt/colors/gen/codegen/writer"

func genRootPkgGradient(ctx *Context, w *writer.GoWriter) {

	genRootPkgGradientGradientStopType(ctx, w)
	w.Separate()
	genRootPkgGradientNewGradientStop(ctx, w)
	w.Separate()

	genRootPkgGradientAtMethod(ctx, w)
}

func genRootPkgGradientAtMethod(ctx *Context, w *writer.GoWriter) {
	var temp []string

	w.Comment("At returns the interpolated color at position t.")
	w.Method("g Gradient", "At")
	w.FuncParams("t ", FloatType)
	w.FuncResults("Color")
	w.FuncBody()

	w.LineWriteln("i := g.findStop(t)")

	w.Separate()
	w.If("i == 0")
	w.LineWriteln("s := g.stops[0]")

	temp = appendVars(temp[:0], "s.c", ctx.MaxChannelCount, "s.alpha")
	w.ReturnInline()
	writeColorLiteral(w, "g.space", ctx.MaxChannelCount, temp...)

	w.End()

	w.Separate()
	w.If("i == len(g.stops)")
	w.LineWriteln("s := g.stops[len(g.stops)-1]")

	w.ReturnInline()
	writeColorLiteral(w, "g.space", ctx.MaxChannelCount, temp...)

	w.End()

	w.Separate()
	w.LineWriteln("a := g.stops[i-1]")
	w.LineWriteln("b := g.stops[i]")

	w.Separate()
	w.LineWriteln("seg := (t - a.position) / (b.position - a.position)")

	w.Separate()
	w.Switch("g.channels")

	hueIndexes := buildHueIndexes(ctx)

	for channelCount, indexes := range hueIndexes {
		if len(indexes) == 0 {
			continue
		}

		w.Case(channelCount)

		w.LineWriteJoin(appendVars(temp[:0], "c", channelCount, "alpha"), ", ")
		w.Writeln(" := g.unsafeMixer.Mix", channelCount, '(')
		w.In()

		w.LineWriteJoin(appendVars(temp[:0], "a.c", channelCount, "a.alpha"), ", ")
		w.Write(",")
		w.LineWriteJoin(appendVars(temp[:0], "b.c", channelCount, "b.alpha"), ", ")
		w.Write(",")
		w.LineWriteln("seg,")

		w.Out()
		w.LineWriteln(')')

		w.ReturnInline()
		writeColorLiteral(w, "g.space", ctx.MaxChannelCount, appendVars(temp[:0], "c", channelCount, "alpha")...)
	}
	w.End()

	w.Separate()
	w.Return("Color{}")
	w.End()
}

func genRootPkgGradientGradientStopType(ctx *Context, w *writer.GoWriter) {
	w.Begin("type gradientStop struct ")
	w.LineWriteln("original ", FloatType)
	w.LineWriteln("position ", FloatType)
	w.Separate()
	for i := range ctx.MaxChannelCount {
		w.LineWriteln("c", i+1, " ", FloatType)
	}
	w.Separate()
	w.LineWriteln("alpha ", FloatType)
	w.End()
}

func genRootPkgGradientNewGradientStop(ctx *Context, w *writer.GoWriter) {
	w.Func("newGradientStop")
	w.FuncParams("original, position float64, color Color")
	w.FuncResults("gradientStop")
	w.FuncBody()

	w.ReturnInline("gradientStop", '{')
	w.In()
	w.LineWriteln("original: original,")
	w.LineWriteln("position: position,")
	for i := range ctx.MaxChannelCount {
		w.LineWriteln("c", i+1, ": color.c", i+1, ",")
	}
	w.LineWriteln("alpha: color.alpha,")
	w.Out()
	w.LineWriteln('}')

	w.End()
}
