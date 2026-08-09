package codegen

import (
	"github.com/thmalt/colors/gen/codegen/model"
	"github.com/thmalt/colors/gen/codegen/writer"
)

func genSpacePkgSpaceInfo(ctx *Context, w *writer.GoWriter) {
	for _, space := range ctx.BuildSpaces {
		w.Begin("var ", spaceInfoName(space), " = spaceInfo")

		w.LineWritef("name:        %q,\n", space.Name)
		w.LineWritef("displayName: %q,\n", space.DisplayName)
		w.LineWritef("cssName:     %q,\n", space.CssName)

		if whitePoint := LookupWhitePoint(space.WhitePoint); whitePoint != nil {
			w.LineWritef("whitePoint: %s,\n", whitePoint.Name)
		}

		w.LineWritef("coordinate: %s,\n", space.Coordinate)

		w.Begin("channels: []Channel")
		for _, c := range space.Channels {
			w.Begin()

			w.LineWritef("Name:        %q,\n", c.Name)
			w.LineWritef("Symbol:      %q,\n", c.Symbol)
			w.LineWritef("DisplayName: %q,\n", c.DisplayName)

			unit := c.Unit
			switch c.Unit {
			case model.UnitRadian, model.UnitGradian, model.UnitTurn:
				w.LineWritef("Min: %g,\n", model.AngleToDegree(c.Min, c.Unit))
				w.LineWritef("Max: %g,\n", model.AngleToDegree(c.Max, c.Unit))
				unit = model.UnitDegree
			default:
				w.LineWritef("Min: %g,\n", c.Min)
				w.LineWritef("Max: %g,\n", c.Max)
			}

			if c.Circular {
				w.LineWritef("Circular: %t,\n", c.Circular)
			}

			if c.Unrestricted {
				w.LineWritef("Unrestricted: %t,\n", c.Unrestricted)
			}

			if c.Unit != model.UnitNumber {
				w.LineWritef("Unit: %s,\n", unit)
			}

			prec := c.Precision
			if prec == 0 {
				prec = DefaultPrecision
			}

			w.LineWritef("Precision: %d,\n", prec)

			w.End(',')
		}
		w.End(',')
		w.LineWritef("hueIndex: %d,\n", space.HueIndex())

		if space.UseGenericColorFunction {
			w.LineWriteln("useColorFunction: true,")
		}

		w.End()

		w.Separate()
	}
}

func spaceInfoName(space *model.Space) string {
	return toLowerCaseFirstWord(space.Name) + "Info"
}
