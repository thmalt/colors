package codegen

import (
	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genConvertPkgWhitePoint(ctx *Context, w *writer.GoWriter) bool {
	if len(ctx.WhitePoints) == 0 {
		return false
	}

	for _, whitepoint := range ctx.WhitePoints {
		w.BeginGroup("const ")

		xyz := data.ChromaToXyz(whitepoint.X, whitepoint.Y)
		name := whitepoint.Name
		privateName := toLowerCaseFirstWord(name)

		w.LineWriteln(privateName, "X = ", formatNormalizedFloat(xyz[0]))
		w.LineWriteln(privateName, "Y = ", formatNormalizedFloat(xyz[1]))
		w.LineWriteln(privateName, "Z = ", formatNormalizedFloat(xyz[2]))

		w.Separate()

		w.LineWriteln("inv", name, "X = 1 / ", privateName, "X")
		w.LineWriteln("inv", name, "Y = 1 / ", privateName, "Y")
		w.LineWriteln("inv", name, "Z = 1 / ", privateName, "Z")

		w.End()

		w.Separate()
	}

	return true
}
