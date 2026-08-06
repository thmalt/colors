package codegen

import (
	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genSpacePkgWhitePoint(ctx *Context, w *writer.GoWriter) {
	w.Begin("type WhitePoint struct ")
	w.LineWriteln("Name string")

	w.Separate()

	w.LineWriteln("X ", FloatType)
	w.LineWriteln("Y ", FloatType)
	w.LineWriteln("Z ", FloatType)
	w.End()

	w.Separate()

	w.BeginGroup("var ")

	for i, whitePoint := range ctx.WhitePoints {
		if i > 0 {
			w.Separate()
		}

		xyz := data.ChromaToXyz(whitePoint.X, whitePoint.Y)

		w.Begin(whitePoint.Name, " = WhitePoint")

		w.LineWritef("Name: %q,\n", whitePoint.Name)
		w.LineWritef("X:    %s,\n", formatNormalizedFloat(xyz[0]))
		w.LineWritef("Y:    %s,\n", formatNormalizedFloat(xyz[1]))
		w.LineWritef("Z:    %s,\n", formatNormalizedFloat(xyz[2]))

		w.End()
	}

	w.End()
}
