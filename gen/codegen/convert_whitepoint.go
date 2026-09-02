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

		xyz := data.ChromaticityToXyz(whitepoint.X, whitepoint.Y)
		name := whitepoint.Name
		privateName := toLowerCaseFirstWord(name)

		w.LineWriteln(privateName, "X = ", formatNormalizedFloat(xyz[0]))
		w.LineWriteln(privateName, "Y = ", formatNormalizedFloat(xyz[1]))
		w.LineWriteln(privateName, "Z = ", formatNormalizedFloat(xyz[2]))

		w.Separate()

		w.LineWriteln(privateName, "InvX = 1 / ", privateName, "X")
		w.LineWriteln(privateName, "InvY = 1 / ", privateName, "Y")
		w.LineWriteln(privateName, "InvZ = 1 / ", privateName, "Z")

		w.End()

		w.Separate()
	}

	return true
}
