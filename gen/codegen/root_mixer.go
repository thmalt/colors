package codegen

import (
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgMixerMethod(ctx *Context, w *writer.GoWriter) {
	channelCounts := buildChannelCounts(ctx)

	// func (m Mixer) Mix(c1, c2 Color, t float64) Color
	w.Comment("Mix converts c1 and c2 to the mixer's color space and linearly interpolates them.")
	w.Comment("The result is returned in the mixer's color space.")
	w.Method("m Mixer", "Mix")
	w.FuncParams("c1, c2 Color, t ", FloatType)
	w.FuncResults("Color")
	w.FuncBody()

	w.LineWriteln("c1, _ = c1.To(m.space)")
	w.LineWriteln("c2, _ = c2.To(m.space)")

	w.Separate()
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
	var temp []string

	w.LineWriteJoin(appendVars(temp[:0], "x", channelCount, "alpha"), ", ")
	w.Write(" := ", name, channelCount, "(")
	w.In()

	w.LineWriteJoin(appendVars(temp[:0], "c1.c", channelCount, "c1.alpha"), ", ")
	w.Write(",")
	w.LineWriteJoin(appendVars(temp[:0], "c2.c", channelCount, "c2.alpha"), ", ")
	w.Write(",")

	w.LineWriteln("t,")

	w.Out()
	w.LineWriteln(")")

	w.ReturnInline()
	writeColorLiteral(w, space, channelCount, appendVars(temp[:0], "x", channelCount, "alpha")...)
}
