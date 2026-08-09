package codegen

import "github.com/thmalt/colors/gen/codegen/writer"

func genRootPkgGradient(ctx *Context, w *writer.GoWriter) {
	hueIndexes := buildHueIndexes(ctx)

	genRootPkgGradientAtMethod(w, hueIndexes)

	w.Separate()
	for channelCount, indexes := range hueIndexes {
		if len(indexes) == 0 {
			continue
		}

		genRootPkgGradientAtNMethod(w, channelCount)
	}

	w.Separate()
	genRootPkgGradientGradientStopType(ctx, w)
	w.Separate()
	genRootPkgGradientNewGradientStop(ctx, w)
}

func genRootPkgGradientAtMethod(w *writer.GoWriter, hueIndexes [][]bool) {
	w.Method("g *Gradient", "At")
	w.FuncParams("t ", FloatType)
	w.FuncResults("Color")
	w.FuncBody()

	w.Switch("g.channels")

	rest := toVars(nil, "c", len(hueIndexes)+1)
	for channelCount, indexes := range hueIndexes {
		if len(indexes) == 0 {
			continue
		}

		rest = toVars(rest[:0], "c", channelCount)
		rest = append(rest, "alpha")

		w.Case(channelCount)
		w.LineWriteJoin(rest, ", ")
		w.Writeln(" := g.at", channelCount, "(t)")

		w.ReturnInline("Color{space: g.space, ")
		for i := range channelCount {
			w.Write(rest[i], ": ", rest[i], ", ")
		}
		w.Writeln(rest[channelCount], ": ", rest[channelCount], "}")

	}
	w.End()

	w.Separate()
	w.Return("Color{}")
	w.End()
}

func genRootPkgGradientAtNMethod(w *writer.GoWriter, channelCount int) {
	rest := toVars(nil, "c", channelCount)
	rest = append(rest, "alpha")

	w.Method("g *Gradient", "at", channelCount)
	w.FuncParams("t ", FloatType)
	w.FuncResults(joinIdentsWithType(FloatType, rest...))
	w.FuncBody()
	w.Begin("i := sort.Search(len(g.stops), func(i int) bool ")
	w.Return("g.stops[i].position >= t")
	w.End(")")

	w.Separate()
	w.If("i == 0")
	w.LineWriteln("s := g.stops[0]")
	w.ReturnInline()
	for i := range channelCount {
		if i > 0 {
			w.Write(", ")
		}
		w.Write("s.c", i+1)
	}
	w.Writeln(", s.alpha")
	w.End()

	w.Separate()
	w.If("i == len(g.stops)")
	w.LineWriteln("s := g.stops[len(g.stops)-1]")
	w.ReturnInline()
	for i := range channelCount {
		if i > 0 {
			w.Write(", ")
		}
		w.Write("s.c", i+1)
	}
	w.Writeln(", s.alpha")
	w.End()

	w.Separate()
	w.LineWriteln("a := g.stops[i-1]")
	w.LineWriteln("b := g.stops[i]")

	w.Separate()
	w.LineWriteln("seg := (t - a.position) / (b.position - a.position)")

	w.Separate()
	w.ReturnInline()
	w.Writeln("g.unsafeMixer.Mix", channelCount, "(")
	w.In()

	w.Indent()
	writeColorChannels(w, "a", channelCount)
	w.Write(",")
	w.Indent()
	writeColorChannels(w, "b", channelCount)
	w.Write(",")

	w.LineWriteln("seg,")
	w.Out()
	w.LineWriteln(")")

	w.End()
}

func genRootPkgGradientGradientStopType(ctx *Context, w *writer.GoWriter) {
	w.Begin("type gradientStop struct ")
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
	w.FuncParams("position float64, color Color")
	w.FuncResults("g gradientStop")
	w.FuncBody()
	w.LineWriteln("g.position = position")
	w.Separate()

	for i := range ctx.MaxChannelCount {
		w.LineWriteln("g.c", i+1, " = color.c", i+1)
	}
	w.Separate()
	w.LineWriteln("g.alpha = color.alpha")

	w.Separate()
	w.Return()
	w.End()
}
