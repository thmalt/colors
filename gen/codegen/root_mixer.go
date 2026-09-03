package codegen

import (
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genRootPkgMixerMethod(ctx *Context, w *writer.GoWriter) {
	channelCounts := buildChannelCounts(ctx)

	sub := w.SubWriter()
	sub.In()
	sub.Switch("m.channels")
	for channelCount, ok := range channelCounts {
		if !ok {
			continue
		}

		sub.Case(channelCount)
		genRootPkgMixerMethodCase(sub, "m.unsafe.Mix", "m.space", channelCount)
	}
	sub.Default()
	sub.Return("Color{}")
	sub.End()

	w.Separate()
	// func (m Mixer) Mix(c1, c2 Color, t float64) Color
	w.Comment("Mix converts c1 and c2 to the mixer's color space and linearly interpolates them.")
	w.Comment("The result is returned in the mixer's color space.")
	w.Method("m Mixer", "Mix")
	w.FuncParams("c1, c2 Color, t ", FloatType)
	w.FuncResults("Color")
	w.FuncBody()

	w.If("c1.space != m.space")
	w.LineWriteln("c1, _ = c1.to(m.space)")
	w.End()
	w.If("c2.space != m.space")
	w.LineWriteln("c2, _ = c2.to(m.space)")
	w.End()

	w.Separate()
	w.Write(sub.Bytes())
	w.End()

	w.Separate()
	// func (m Mixer) UnsafeMix(c1, c2 Color, t float64) Color
	w.Comment("UnsafeMix linearly interpolates c1 and c2 in the mixer's color space.")
	w.Comment("It assumes both colors are already in the mixer's color space.")
	w.Method("m Mixer", "UnsafeMix")
	w.FuncParams("c1, c2 Color, t ", FloatType)
	w.FuncResults("Color")
	w.FuncBody()
	w.Drain(sub)
	w.End()
}

func genRootPkgMixerMethodCase(w *writer.GoWriter, name, space string, channelCount int) {
	temp := make([]string, channelCount+1)

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
