package codegen

import (
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgMixerMethod(ctx *Context, w *writer.GoWriter) {
	channelCounts := buildChannelCounts(ctx)

	w.Method("m Mixer", "Mix")
	w.FuncParams("c1, c2 Color, t ", FloatType)
	w.FuncResults("Color")
	w.FuncBody()
	w.Switch("m.channels")
	for channelCount, ok := range channelCounts {
		if !ok {
			continue
		}

		w.Case(channelCount)
		genRootPkgMixerMethodCase(w, "m.unsafe.Mix", "m.space", channelCount)
	}
	w.Default()
	w.Return("Color{}")
	w.End()
	w.End()
	w.Separate()
}

func genRootPkgMixerMethodCase(w *writer.GoWriter, name, space string, channelCount int) {
	vars := make([]string, 0, channelCount+1)

	vars = toVars(vars, "x", channelCount)
	vars = append(vars, "alpha")

	w.LineWriteJoin(vars, ", ")
	w.Write(" := ", name, channelCount, "(")
	w.In()

	w.Indent()
	writeColorChannels(w, "c1", channelCount)
	w.Write(",")
	w.Indent()
	writeColorChannels(w, "c2", channelCount)
	w.Write(",")
	w.LineWriteln("t,")

	w.Out()
	w.LineWriteln(")")

	w.ReturnInline()
	writeColorLiteral(w, space, channelCount, vars...)
}
